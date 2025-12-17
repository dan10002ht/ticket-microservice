package com.ticketing.booking.service.saga;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

import org.redisson.api.RLock;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.ticketing.booking.entity.Booking;
import com.ticketing.booking.entity.FailedCompensation.CompensationType;
import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.entity.enums.PaymentStatus;
import com.ticketing.booking.exception.BookingLockException;
import com.ticketing.booking.grpcclient.PaymentServiceClient;
import com.ticketing.booking.grpcclient.TicketServiceClient;
import com.ticketing.booking.metrics.BookingMetricsService;
import com.ticketing.booking.repository.BookingRepository;
import com.ticketing.booking.service.BookingLockService;
import com.ticketing.booking.service.CompensationRetryService;
import com.ticketing.booking.service.OutboxService;
import com.ticketing.booking.service.dto.BookingCreateCommand;
import com.ticketing.booking.service.dto.BookingResult;
import com.ticketing.booking.service.mapper.BookingMapper;
import com.ticketing.booking.statemachine.BookingEvent;
import com.ticketing.booking.statemachine.BookingStateMachine;
import com.ticketing.booking.statemachine.StateTransition;
import com.ticketing.booking.util.ReferenceGenerator;
import com.ticketing.booking.grpc.TicketProto;
import com.ticketing.payment.grpc.PaymentProto;

import io.micrometer.core.instrument.Timer;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

/**
 * Saga Orchestrator for Booking Flow
 * 
 * Orchestrates the distributed transaction:
 * 1. Reserve seats (Ticket Service)
 * 2. Process payment (Payment Service)
 * 3. Confirm booking or compensate on failure
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class BookingSagaOrchestrator {

    private final BookingRepository bookingRepository;
    private final BookingMapper bookingMapper;
    private final BookingLockService lockService;
    private final OutboxService outboxService;
    private final TicketServiceClient ticketServiceClient;
    private final PaymentServiceClient paymentServiceClient;
    private final BookingMetricsService metricsService;
    private final CompensationRetryService compensationRetryService;
    private final BookingStateMachine stateMachine;

    /**
     * Execute booking saga
     *
     * @param command Booking creation command
     * @return Booking result
     */
    @Transactional
    public BookingResult executeBookingSaga(BookingCreateCommand command) {
        Timer.Sample sagaTimer = metricsService.startSagaExecution();
        RLock lock = null;
        Booking booking = null;
        String reservationId = null;
        String paymentId = null;

        try {
            // Step 0: Check idempotency - return existing booking if duplicate request
            if (command.getIdempotencyKey() != null && !command.getIdempotencyKey().isEmpty()) {
                var existingBooking = bookingRepository.findByIdempotencyKey(command.getIdempotencyKey());
                if (existingBooking.isPresent()) {
                    log.info("Duplicate booking request detected, returning existing booking. idempotencyKey={}",
                            command.getIdempotencyKey());
                    metricsService.recordSagaStep("idempotency_check", "duplicate");
                    return bookingMapper.toResult(existingBooking.get());
                }
            }

            // Step 1: Acquire distributed lock
            lock = acquireLock(command.getEventId());
            metricsService.recordSagaStep("acquire_lock", "success");

            // Step 2: Create booking record (PENDING state)
            booking = createBookingRecord(command);
            booking = bookingRepository.save(booking);
            metricsService.recordSagaStep("create_booking", "success");
            log.info("Created booking record: {}, idempotencyKey={}", booking.getBookingId(), command.getIdempotencyKey());

            // Step 3: Reserve seats via Ticket Service
            StateTransition toReserving = stateMachine.transition(
                    booking, BookingEvent.RESERVE_SEATS, "saga", "Starting seat reservation");
            bookingRepository.save(booking);
            log.debug("State transition: {}", toReserving);

            reservationId = reserveSeats(booking, command);

            // Transition to SEATS_RESERVED
            StateTransition toSeatsReserved = stateMachine.transition(
                    booking, BookingEvent.SEATS_RESERVED, "saga", "Seats reserved successfully");
            bookingRepository.save(booking);
            log.debug("State transition: {}", toSeatsReserved);

            metricsService.recordSagaStep("reserve_seats", "success");
            log.info("Reserved seats for booking: {}, reservationId: {}", booking.getBookingId(), reservationId);

            // Step 4: Process payment if required
            if (requiresPayment(command)) {
                // Transition to PAYMENT_PENDING
                StateTransition toPaymentPending = stateMachine.transition(
                        booking, BookingEvent.REQUEST_PAYMENT, "saga", "Starting payment processing");
                bookingRepository.save(booking);
                log.debug("State transition: {}", toPaymentPending);

                paymentId = processPayment(booking, command);

                // Transition to PAYMENT_PROCESSING
                StateTransition toPaymentProcessing = stateMachine.transition(
                        booking, BookingEvent.PAYMENT_AUTHORIZED, "saga", "Payment authorized");
                bookingRepository.save(booking);
                log.debug("State transition: {}", toPaymentProcessing);

                metricsService.recordSagaStep("process_payment", "success");
                log.info("Processed payment for booking: {}, paymentId: {}", booking.getBookingId(), paymentId);
            } else {
                metricsService.recordSagaStep("process_payment", "skipped");
            }

            // Step 5: Confirm booking
            if (requiresPayment(command)) {
                // Transition via PAYMENT_CAPTURED event
                StateTransition toConfirmed = stateMachine.transition(
                        booking, BookingEvent.PAYMENT_CAPTURED, "saga", "Payment captured successfully");
                log.debug("State transition: {}", toConfirmed);
            } else {
                // Free event - transition via CONFIRM event
                StateTransition toConfirmed = stateMachine.transition(
                        booking, BookingEvent.CONFIRM, "saga", "Free event - confirmed directly");
                log.debug("State transition: {}", toConfirmed);
            }
            booking.setPaymentStatus(PaymentStatus.CAPTURED);
            if (paymentId != null) {
                booking.setPaymentReference(paymentId);
            }
            booking = bookingRepository.save(booking);

            // Save event to outbox (same transaction - guarantees exactly-once delivery)
            outboxService.saveBookingConfirmedEvent(booking);
            metricsService.recordSagaStep("confirm_booking", "success");
            log.info("Booking saga completed successfully: {}", booking.getBookingId());

            metricsService.recordSagaExecutionDuration(sagaTimer, "success");
            return bookingMapper.toResult(booking);

        } catch (Exception e) {
            log.error("Booking saga failed, initiating compensation", e);
            String errorType = e.getClass().getSimpleName();
            metricsService.recordBookingFailed(command.getEventId(), errorType);
            metricsService.recordSagaExecutionDuration(sagaTimer, "failed");
            compensate(booking, reservationId, paymentId, e.getMessage());
            throw new RuntimeException("Booking saga failed: " + e.getMessage(), e);
        } finally {
            releaseLock(lock);
        }
    }

    /**
     * Compensate for failed saga steps.
     * Failed compensations are saved to DLQ for automatic retry.
     */
    private void compensate(Booking booking, String reservationId, String paymentId, String reason) {
        Long bookingId = booking != null ? booking.getId() : null;
        String bookingRef = booking != null ? booking.getBookingReference() : null;

        try {
            // Compensate payment if processed
            if (paymentId != null) {
                try {
                    paymentServiceClient.cancelPayment(paymentId);
                    metricsService.recordSagaCompensation("cancel_payment", "success");
                    log.info("Compensated payment: {}", paymentId);
                } catch (Exception e) {
                    metricsService.recordSagaCompensation("cancel_payment", "failed");
                    log.error("Failed to compensate payment, saving to DLQ: {}", paymentId, e);
                    // Save to DLQ for retry
                    compensationRetryService.saveFailedCompensation(
                            CompensationType.CANCEL_PAYMENT,
                            paymentId,
                            bookingId,
                            bookingRef,
                            e.getMessage(),
                            null);
                }
            }

            // Release seats if reserved
            if (reservationId != null && booking != null) {
                try {
                    ticketServiceClient.releaseTickets(reservationId, null);
                    metricsService.recordSagaCompensation("release_seats", "success");
                    log.info("Released seats for reservation: {}", reservationId);
                } catch (Exception e) {
                    metricsService.recordSagaCompensation("release_seats", "failed");
                    log.error("Failed to release seats, saving to DLQ: {}", reservationId, e);
                    // Save to DLQ for retry
                    compensationRetryService.saveFailedCompensation(
                            CompensationType.RELEASE_SEATS,
                            reservationId,
                            bookingId,
                            bookingRef,
                            e.getMessage(),
                            null);
                }
            }

            // Update booking status via state machine
            if (booking != null) {
                try {
                    StateTransition toFailed = stateMachine.transition(
                            booking, BookingEvent.FAIL, "saga-compensation", reason);
                    log.debug("State transition during compensation: {}", toFailed);
                } catch (Exception e) {
                    // If state machine fails, force the status (terminal state)
                    log.warn("State machine transition failed, forcing FAILED status", e);
                    booking.setStatus(BookingStatus.FAILED);
                }
                booking.setPaymentStatus(PaymentStatus.FAILED);
                bookingRepository.save(booking);
                // Save cancellation event to outbox (same transaction)
                outboxService.saveBookingCancelledEvent(booking);
            }
        } catch (Exception e) {
            log.error("Compensation failed critically", e);
        }
    }

    /**
     * Step 1: Acquire distributed lock
     */
    private RLock acquireLock(String eventId) {
        try {
            RLock lock = lockService.acquireLock(eventId);
            if (lock == null) {
                throw new BookingLockException("Failed to acquire lock for event: " + eventId);
            }
            return lock;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new BookingLockException("Interrupted while acquiring lock", e);
        }
    }

    /**
     * Step 2: Create booking record
     */
    private Booking createBookingRecord(BookingCreateCommand command) {
        BigDecimal totalAmount = command.getTotalAmount();
        return Booking.builder()
                .bookingReference(ReferenceGenerator.bookingReference())
                .userId(command.getUserId())
                .eventId(command.getEventId())
                .seatCount(command.getSeatCount())
                .seatNumbers(command.getSeatNumbers())
                .totalAmount(totalAmount)
                .currency(command.getCurrency())
                .expiresAt(ReferenceGenerator.bookingExpiry())
                .status(BookingStatus.PENDING)
                .paymentStatus(totalAmount.compareTo(BigDecimal.ZERO) > 0
                        ? PaymentStatus.PENDING
                        : PaymentStatus.NOT_REQUIRED)
                .idempotencyKey(command.getIdempotencyKey()) // Set idempotency key for duplicate detection
                .build();
    }

    /**
     * Step 3: Reserve seats via Ticket Service
     */
    private String reserveSeats(Booking booking, BookingCreateCommand command) {
        try {
            TicketProto.ReserveTicketsResponse response = ticketServiceClient.reserveTickets(
                    booking.getEventId(),
                    command.getSeatNumbers(),
                    booking.getUserId(),
                    900); // 15 minutes timeout

            if (!response.getSuccess()) {
                metricsService.recordSagaStep("reserve_seats", "failed");
                throw new RuntimeException("Failed to reserve seats: " + response.getMessage());
            }

            return response.getReservationId();
        } catch (Exception e) {
            metricsService.recordSagaStep("reserve_seats", "failed");
            log.error("Failed to reserve seats", e);
            throw new RuntimeException("Seat reservation failed", e);
        }
    }

    /**
     * Step 4: Process payment via Payment Service
     */
    private String processPayment(Booking booking, BookingCreateCommand command) {
        try {
            PaymentProto.Payment payment = paymentServiceClient.createPayment(
                    booking.getBookingId().toString(),
                    booking.getUserId(),
                    booking.getTotalAmount(),
                    booking.getCurrency(),
                    "CARD", // TODO: Get from command
                    "STRIPE", // TODO: Get from command
                    booking.getBookingReference(), // Use booking reference as idempotency key
                    null);

            if (payment == null || payment.getStatus().equals("FAILED")) {
                metricsService.recordSagaStep("process_payment", "failed");
                throw new RuntimeException("Payment creation failed");
            }

            // Capture payment immediately (authorize + capture in one step)
            PaymentProto.Payment capturedPayment = paymentServiceClient.capturePayment(payment.getId());

            if (capturedPayment == null || !capturedPayment.getStatus().equals("CAPTURED")) {
                metricsService.recordSagaStep("process_payment", "failed");
                throw new RuntimeException("Payment capture failed");
            }

            return capturedPayment.getId();
        } catch (Exception e) {
            metricsService.recordSagaStep("process_payment", "failed");
            log.error("Failed to process payment", e);
            throw new RuntimeException("Payment processing failed", e);
        }
    }

    /**
     * Check if payment is required
     */
    private boolean requiresPayment(BookingCreateCommand command) {
        return command.getTotalAmount() != null
                && command.getTotalAmount().compareTo(BigDecimal.ZERO) > 0;
    }

    /**
     * Release lock
     */
    private void releaseLock(RLock lock) {
        if (lock != null) {
            try {
                lockService.releaseLock(lock);
            } catch (Exception e) {
                log.error("Failed to release lock", e);
            }
        }
    }
}

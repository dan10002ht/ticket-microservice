package com.ticketing.booking.kafka;

import java.util.UUID;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.booking.entity.Booking;
import com.ticketing.booking.entity.FailedCompensation.CompensationType;
import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.entity.enums.PaymentStatus;
import com.ticketing.booking.exception.InvalidStateTransitionException;
import com.ticketing.booking.grpcclient.TicketServiceClient;
import com.ticketing.booking.repository.BookingRepository;
import com.ticketing.booking.service.CompensationRetryService;
import com.ticketing.booking.service.OutboxService;
import com.ticketing.booking.statemachine.BookingEvent;
import com.ticketing.booking.statemachine.BookingStateMachine;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import lombok.extern.slf4j.Slf4j;

/**
 * Kafka Consumer for Payment Events.
 *
 * Listens to payment-events topic and updates booking status based on payment status.
 * This enables async payment flow where saga doesn't wait for payment capture.
 *
 * Events handled:
 * - PAYMENT_AUTHORIZED: Payment authorized, funds held
 * - PAYMENT_CAPTURED: Payment captured successfully
 * - PAYMENT_FAILED: Payment failed
 * - PAYMENT_REFUNDED: Payment refunded
 */
@Component
@Slf4j
public class PaymentEventConsumer {

    private final BookingRepository bookingRepository;
    private final OutboxService outboxService;
    private final TicketServiceClient ticketServiceClient;
    private final CompensationRetryService compensationRetryService;
    private final BookingStateMachine stateMachine;
    private final ObjectMapper objectMapper;
    private final Counter paymentCapturedCounter;
    private final Counter paymentFailedCounter;

    public PaymentEventConsumer(
            BookingRepository bookingRepository,
            OutboxService outboxService,
            TicketServiceClient ticketServiceClient,
            CompensationRetryService compensationRetryService,
            BookingStateMachine stateMachine,
            ObjectMapper objectMapper,
            MeterRegistry meterRegistry) {
        this.bookingRepository = bookingRepository;
        this.outboxService = outboxService;
        this.ticketServiceClient = ticketServiceClient;
        this.compensationRetryService = compensationRetryService;
        this.stateMachine = stateMachine;
        this.objectMapper = objectMapper;

        this.paymentCapturedCounter = meterRegistry.counter("booking.payment.captured");
        this.paymentFailedCounter = meterRegistry.counter("booking.payment.failed");
    }

    @KafkaListener(
            topics = "${kafka.topics.payment-events:payment-events}",
            groupId = "${kafka.consumer.group-id:booking-service}",
            containerFactory = "kafkaListenerContainerFactory"
    )
    public void onPaymentEvent(String message) {
        try {
            JsonNode event = objectMapper.readTree(message);
            String eventType = event.get("type").asText();
            String paymentId = event.get("paymentId").asText();
            String bookingId = event.has("bookingId") ? event.get("bookingId").asText() : null;

            log.info("Received payment event: type={}, paymentId={}, bookingId={}",
                    eventType, paymentId, bookingId);

            switch (eventType) {
                case "PAYMENT_AUTHORIZED" -> handlePaymentAuthorized(event);
                case "PAYMENT_CAPTURED" -> handlePaymentCaptured(event);
                case "PAYMENT_FAILED" -> handlePaymentFailed(event);
                case "PAYMENT_REFUNDED" -> handlePaymentRefunded(event);
                default -> log.warn("Unknown payment event type: {}", eventType);
            }

        } catch (Exception e) {
            log.error("Error processing payment event: {}", message, e);
            // Don't throw - let Kafka commit offset to avoid infinite retry
            // Failed events should be handled through DLQ or manual reconciliation
        }
    }

    /**
     * Handle payment authorized event.
     * Payment authorized means funds are held but not yet captured.
     */
    @Transactional
    public void handlePaymentAuthorized(JsonNode event) {
        String bookingId = event.get("bookingId").asText();
        String paymentId = event.get("paymentId").asText();

        Booking booking = findBookingById(bookingId);
        if (booking == null) {
            log.warn("Booking not found for payment authorized event: bookingId={}", bookingId);
            return;
        }

        try {
            // Use state machine for transition
            stateMachine.transition(booking, BookingEvent.PAYMENT_AUTHORIZED,
                    "payment-webhook", "Payment authorized via webhook");
            booking.setPaymentStatus(PaymentStatus.PROCESSING);
            booking.setPaymentReference(paymentId);
            bookingRepository.save(booking);

            log.info("Booking status updated to PAYMENT_PROCESSING: bookingId={}, paymentId={}",
                    bookingId, paymentId);
        } catch (InvalidStateTransitionException e) {
            log.warn("Invalid state transition for payment authorized: bookingId={}, currentStatus={}",
                    bookingId, booking.getStatus());
        }
    }

    /**
     * Handle payment captured event.
     * This confirms the booking - the critical success path.
     */
    @Transactional
    public void handlePaymentCaptured(JsonNode event) {
        String bookingId = event.get("bookingId").asText();
        String paymentId = event.get("paymentId").asText();

        Booking booking = findBookingById(bookingId);
        if (booking == null) {
            log.warn("Booking not found for payment captured event: bookingId={}", bookingId);
            return;
        }

        // Check if already confirmed (idempotent)
        if (booking.getStatus() == BookingStatus.CONFIRMED) {
            log.debug("Booking already confirmed, ignoring duplicate event: bookingId={}", bookingId);
            return;
        }

        try {
            // Use state machine for transition to CONFIRMED
            stateMachine.transition(booking, BookingEvent.PAYMENT_CAPTURED,
                    "payment-webhook", "Payment captured via webhook");
            booking.setPaymentStatus(PaymentStatus.CAPTURED);
            booking.setPaymentReference(paymentId);
            bookingRepository.save(booking);

            // Publish booking confirmed event through outbox
            outboxService.saveBookingConfirmedEvent(booking);

            paymentCapturedCounter.increment();
            log.info("Booking confirmed via payment webhook: bookingId={}, paymentId={}",
                    bookingId, paymentId);

        } catch (InvalidStateTransitionException e) {
            log.warn("Invalid state transition for payment captured: bookingId={}, currentStatus={}",
                    bookingId, booking.getStatus());
        }
    }

    /**
     * Handle payment failed event.
     * This triggers compensation - release seats and cancel booking.
     */
    @Transactional
    public void handlePaymentFailed(JsonNode event) {
        String bookingId = event.get("bookingId").asText();
        String paymentId = event.get("paymentId").asText();
        String failureReason = event.has("failureReason")
                ? event.get("failureReason").asText()
                : "Payment processing failed";

        Booking booking = findBookingById(bookingId);
        if (booking == null) {
            log.warn("Booking not found for payment failed event: bookingId={}", bookingId);
            return;
        }

        // Check if already terminal state (idempotent)
        if (stateMachine.isTerminalState(booking.getStatus())) {
            log.debug("Booking already in terminal state, ignoring event: bookingId={}, status={}",
                    bookingId, booking.getStatus());
            return;
        }

        try {
            // Use state machine for transition to PAYMENT_FAILED
            stateMachine.transition(booking, BookingEvent.PAYMENT_FAILED,
                    "payment-webhook", "Payment failed via webhook: " + failureReason);
            booking.setPaymentStatus(PaymentStatus.FAILED);
            booking.setPaymentReference(paymentId);
            bookingRepository.save(booking);

            // Release seats (compensation)
            releaseSeatsSafely(booking, failureReason);

            // Use state machine to transition to CANCELLED
            stateMachine.transition(booking, BookingEvent.SYSTEM_CANCEL,
                    "payment-webhook", "Auto-cancel after payment failure");
            bookingRepository.save(booking);

            // Publish booking cancelled event through outbox
            outboxService.saveBookingCancelledEvent(booking);

            paymentFailedCounter.increment();
            log.info("Booking cancelled due to payment failure: bookingId={}, reason={}",
                    bookingId, failureReason);

        } catch (InvalidStateTransitionException e) {
            log.warn("Invalid state transition for payment failed: bookingId={}, currentStatus={}",
                    bookingId, booking.getStatus());
        }
    }

    /**
     * Handle payment refunded event.
     * Update booking status to reflect refund.
     */
    @Transactional
    public void handlePaymentRefunded(JsonNode event) {
        String bookingId = event.get("bookingId").asText();
        String paymentId = event.get("paymentId").asText();

        Booking booking = findBookingById(bookingId);
        if (booking == null) {
            log.warn("Booking not found for payment refunded event: bookingId={}", bookingId);
            return;
        }

        // Update payment status to refunded
        booking.setPaymentStatus(PaymentStatus.REFUNDED);
        bookingRepository.save(booking);

        log.info("Booking payment refunded: bookingId={}, paymentId={}", bookingId, paymentId);
    }

    /**
     * Find booking by UUID string
     */
    private Booking findBookingById(String bookingIdStr) {
        try {
            UUID bookingUuid = UUID.fromString(bookingIdStr);
            return bookingRepository.findByBookingId(bookingUuid).orElse(null);
        } catch (IllegalArgumentException e) {
            log.error("Invalid booking ID format: {}", bookingIdStr);
            return null;
        }
    }

    /**
     * Safely release seats with DLQ fallback
     */
    private void releaseSeatsSafely(Booking booking, String reason) {
        try {
            // Use eventId as reservation ID (or store actual reservation ID in booking)
            ticketServiceClient.releaseTickets(booking.getEventId(), booking.getSeatNumbers());
            log.info("Released seats for booking: bookingId={}", booking.getBookingId());

        } catch (Exception e) {
            log.error("Failed to release seats for booking: bookingId={}", booking.getBookingId(), e);

            // Save to DLQ for retry
            compensationRetryService.saveFailedCompensation(
                    CompensationType.RELEASE_SEATS,
                    booking.getEventId(),
                    booking.getId(),
                    booking.getBookingReference(),
                    e.getMessage(),
                    null
            );
        }
    }
}

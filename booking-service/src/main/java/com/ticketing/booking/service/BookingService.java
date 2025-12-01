package com.ticketing.booking.service;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.UUID;

import org.redisson.api.RLock;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import com.ticketing.booking.entity.Booking;
import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.entity.enums.PaymentStatus;
import com.ticketing.booking.exception.BookingNotFoundException;
import com.ticketing.booking.metrics.BookingMetricsService;
import com.ticketing.booking.repository.BookingRepository;
import com.ticketing.booking.service.dto.BookingCreateCommand;
import com.ticketing.booking.service.dto.BookingResult;
import com.ticketing.booking.service.mapper.BookingMapper;
import com.ticketing.booking.service.saga.BookingSagaOrchestrator;
import com.ticketing.booking.util.ReferenceGenerator;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class BookingService {

    private final BookingRepository bookingRepository;
    private final BookingMapper bookingMapper;
    private final BookingLockService lockService;
    private final BookingEventPublisher eventPublisher;
    private final BookingSagaOrchestrator sagaOrchestrator;
    private final BookingMetricsService metricsService;

    /**
     * Create booking using saga orchestrator
     * This method orchestrates the full booking flow:
     * 1. Reserve seats (Ticket Service)
     * 2. Process payment (Payment Service)
     * 3. Confirm booking or compensate on failure
     */
    @Transactional
    public BookingResult createBooking(BookingCreateCommand command) {
        validate(command);
        BookingResult result = sagaOrchestrator.executeBookingSaga(command);
        metricsService.recordBookingCreated(command.getEventId());
        return result;
    }

    @Transactional(readOnly = true)
    public Booking getBooking(UUID bookingId) {
        return bookingRepository.findByBookingId(bookingId)
                .orElseThrow(() -> new BookingNotFoundException(bookingId));
    }

    @Transactional
    public BookingResult confirmBooking(UUID bookingId, String paymentReference) {
        Booking booking = getBooking(bookingId);
        booking.confirm(paymentReference);
        bookingRepository.save(booking);
        eventPublisher.publishBookingConfirmed(booking);
        metricsService.recordBookingConfirmed(booking.getEventId());
        return bookingMapper.toResult(booking);
    }

    @Transactional
    public BookingResult cancelBooking(UUID bookingId, String reason) {
        Booking booking = getBooking(bookingId);
        booking.cancel(reason);
        bookingRepository.save(booking);
        eventPublisher.publishBookingCancelled(booking);
        metricsService.recordBookingCancelled(booking.getEventId(), reason);
        return bookingMapper.toResult(booking);
    }

    private Booking buildBooking(BookingCreateCommand command) {
        BigDecimal totalAmount = command.getTotalAmount();
        return Booking.builder()
                .bookingReference(ReferenceGenerator.bookingReference())
                .userId(command.getUserId())
                .eventId(command.getEventId())
                .seatCount(command.getSeatCount())
                .seatNumbers(command.getSeatNumbers())
                .totalAmount(totalAmount)
                .currency(command.getCurrency())
                .metadata(command.getMetadata() == null ? new HashMap<>() : new HashMap<>(command.getMetadata()))
                .expiresAt(ReferenceGenerator.bookingExpiry())
                .status(BookingStatus.PENDING)
                .paymentStatus(totalAmount.compareTo(BigDecimal.ZERO) > 0 ? PaymentStatus.PENDING
                        : PaymentStatus.NOT_REQUIRED)
                .build();
    }

    private void validate(BookingCreateCommand command) {
        if (!StringUtils.hasText(command.getUserId()) || !StringUtils.hasText(command.getEventId())) {
            throw new IllegalArgumentException("User and Event are required");
        }
        if (command.getSeatNumbers().isEmpty()) {
            throw new IllegalArgumentException("Seat selection is required");
        }
    }
}


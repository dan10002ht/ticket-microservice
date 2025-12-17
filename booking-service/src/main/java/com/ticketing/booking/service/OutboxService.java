package com.ticketing.booking.service;

import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.HashMap;
import java.util.Map;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.booking.entity.Booking;
import com.ticketing.booking.entity.OutboxEvent;
import com.ticketing.booking.entity.OutboxEvent.AggregateType;
import com.ticketing.booking.repository.OutboxEventRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

/**
 * Service for saving events to the outbox table.
 *
 * IMPORTANT: All methods in this service should be called within an existing transaction
 * to ensure atomicity between business data and event persistence.
 *
 * Usage:
 * <pre>
 * {@code
 * @Transactional
 * public void createBooking(BookingCreateCommand command) {
 *     Booking booking = bookingRepository.save(newBooking);
 *     outboxService.saveBookingCreatedEvent(booking);  // Same transaction!
 * }
 * }
 * </pre>
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class OutboxService {

    private final OutboxEventRepository outboxEventRepository;
    private final ObjectMapper objectMapper;

    // Event types
    public static final String BOOKING_CREATED = "BOOKING_CREATED";
    public static final String BOOKING_CONFIRMED = "BOOKING_CONFIRMED";
    public static final String BOOKING_CANCELLED = "BOOKING_CANCELLED";
    public static final String BOOKING_FAILED = "BOOKING_FAILED";
    public static final String PAYMENT_REQUESTED = "PAYMENT_REQUESTED";
    public static final String SEATS_RESERVED = "SEATS_RESERVED";
    public static final String SEATS_RELEASED = "SEATS_RELEASED";

    /**
     * Save a booking created event to the outbox.
     * Must be called within an existing transaction.
     */
    @Transactional(propagation = Propagation.MANDATORY)
    public void saveBookingCreatedEvent(Booking booking) {
        saveBookingEvent(booking, BOOKING_CREATED);
    }

    /**
     * Save a booking confirmed event to the outbox.
     * Must be called within an existing transaction.
     */
    @Transactional(propagation = Propagation.MANDATORY)
    public void saveBookingConfirmedEvent(Booking booking) {
        saveBookingEvent(booking, BOOKING_CONFIRMED);
    }

    /**
     * Save a booking cancelled event to the outbox.
     * Must be called within an existing transaction.
     */
    @Transactional(propagation = Propagation.MANDATORY)
    public void saveBookingCancelledEvent(Booking booking) {
        saveBookingEvent(booking, BOOKING_CANCELLED);
    }

    /**
     * Save a booking failed event to the outbox.
     * Must be called within an existing transaction.
     */
    @Transactional(propagation = Propagation.MANDATORY)
    public void saveBookingFailedEvent(Booking booking, String reason) {
        Map<String, Object> payload = buildBookingPayload(booking);
        payload.put("failureReason", reason);
        saveOutboxEvent(AggregateType.BOOKING, booking.getBookingId().toString(), BOOKING_FAILED, payload);
    }

    /**
     * Internal method to save booking events
     */
    private void saveBookingEvent(Booking booking, String eventType) {
        Map<String, Object> payload = buildBookingPayload(booking);
        saveOutboxEvent(AggregateType.BOOKING, booking.getBookingId().toString(), eventType, payload);
    }

    /**
     * Build booking payload with all relevant data
     */
    private Map<String, Object> buildBookingPayload(Booking booking) {
        Map<String, Object> payload = new HashMap<>();
        payload.put("bookingId", booking.getBookingId().toString());
        payload.put("bookingReference", booking.getBookingReference());
        payload.put("userId", booking.getUserId());
        payload.put("eventId", booking.getEventId());
        payload.put("status", booking.getStatus().name());
        payload.put("paymentStatus", booking.getPaymentStatus().name());
        payload.put("totalAmount", booking.getTotalAmount());
        payload.put("currency", booking.getCurrency());
        payload.put("seatCount", booking.getSeatCount());
        payload.put("seatNumbers", booking.getSeatNumbers());

        if (booking.getPaymentReference() != null) {
            payload.put("paymentReference", booking.getPaymentReference());
        }
        if (booking.getConfirmedAt() != null) {
            payload.put("confirmedAt", booking.getConfirmedAt().toString());
        }
        if (booking.getCancelledAt() != null) {
            payload.put("cancelledAt", booking.getCancelledAt().toString());
        }
        if (booking.getCancellationReason() != null) {
            payload.put("cancellationReason", booking.getCancellationReason());
        }

        payload.put("timestamp", OffsetDateTime.now().toString());
        return payload;
    }

    /**
     * Generic method to save outbox events
     */
    private void saveOutboxEvent(AggregateType aggregateType, String aggregateId,
                                  String eventType, Map<String, Object> payload) {
        try {
            String payloadJson = objectMapper.writeValueAsString(payload);

            OutboxEvent event = OutboxEvent.builder()
                    .aggregateType(aggregateType)
                    .aggregateId(aggregateId)
                    .eventType(eventType)
                    .payload(payloadJson)
                    .build();

            outboxEventRepository.save(event);

            log.debug("Saved outbox event: type={}, aggregateId={}, eventType={}",
                    aggregateType, aggregateId, eventType);

        } catch (JsonProcessingException e) {
            log.error("Failed to serialize outbox event payload: type={}, aggregateId={}, eventType={}",
                    aggregateType, aggregateId, eventType, e);
            throw new RuntimeException("Failed to serialize outbox event", e);
        }
    }

    /**
     * Get pending event count (for monitoring)
     */
    public long getPendingEventCount() {
        return outboxEventRepository.countByStatus(OutboxEvent.OutboxStatus.PENDING);
    }

    /**
     * Get failed event count (for monitoring)
     */
    public long getFailedEventCount() {
        return outboxEventRepository.countByStatus(OutboxEvent.OutboxStatus.FAILED);
    }
}

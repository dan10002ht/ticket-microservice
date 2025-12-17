package com.ticketing.booking.statemachine;

import java.time.OffsetDateTime;
import java.util.Map;
import java.util.UUID;

import com.ticketing.booking.entity.enums.BookingStatus;

import lombok.Builder;
import lombok.Value;

/**
 * Represents a state transition in the booking lifecycle.
 *
 * Used for:
 * - Audit logging
 * - Debugging booking flow
 * - Analytics and metrics
 */
@Value
@Builder
public class StateTransition {

    /**
     * Unique ID for this transition
     */
    UUID transitionId;

    /**
     * Booking ID this transition belongs to
     */
    UUID bookingId;

    /**
     * State before the transition
     */
    BookingStatus fromState;

    /**
     * State after the transition
     */
    BookingStatus toState;

    /**
     * Event that triggered this transition
     */
    BookingEvent event;

    /**
     * When this transition occurred
     */
    OffsetDateTime timestamp;

    /**
     * Who/what triggered this transition (user ID, system, webhook, etc.)
     */
    String triggeredBy;

    /**
     * Optional reason for the transition
     */
    String reason;

    /**
     * Additional metadata about the transition
     */
    Map<String, String> metadata;

    /**
     * Create a new transition
     */
    public static StateTransition create(
            UUID bookingId,
            BookingStatus fromState,
            BookingStatus toState,
            BookingEvent event,
            String triggeredBy,
            String reason) {
        return StateTransition.builder()
                .transitionId(UUID.randomUUID())
                .bookingId(bookingId)
                .fromState(fromState)
                .toState(toState)
                .event(event)
                .timestamp(OffsetDateTime.now())
                .triggeredBy(triggeredBy)
                .reason(reason)
                .build();
    }

    /**
     * Create a new transition with metadata
     */
    public static StateTransition create(
            UUID bookingId,
            BookingStatus fromState,
            BookingStatus toState,
            BookingEvent event,
            String triggeredBy,
            String reason,
            Map<String, String> metadata) {
        return StateTransition.builder()
                .transitionId(UUID.randomUUID())
                .bookingId(bookingId)
                .fromState(fromState)
                .toState(toState)
                .event(event)
                .timestamp(OffsetDateTime.now())
                .triggeredBy(triggeredBy)
                .reason(reason)
                .metadata(metadata)
                .build();
    }
}

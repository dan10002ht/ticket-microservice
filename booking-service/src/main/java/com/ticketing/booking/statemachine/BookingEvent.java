package com.ticketing.booking.statemachine;

/**
 * Events that trigger state transitions in the booking lifecycle.
 *
 * Each event represents an action or occurrence that moves a booking
 * from one state to another.
 */
public enum BookingEvent {
    // ==================== Creation Events ====================
    /**
     * Booking request received and validated
     */
    CREATE,

    // ==================== Seat Reservation Events ====================
    /**
     * Start reserving seats
     */
    RESERVE_SEATS,

    /**
     * Seats reserved successfully
     */
    SEATS_RESERVED,

    /**
     * Seat reservation failed (no seats available, timeout, etc.)
     */
    SEATS_RESERVATION_FAILED,

    // ==================== Payment Events ====================
    /**
     * Payment request initiated
     */
    REQUEST_PAYMENT,

    /**
     * Payment authorized (funds held)
     */
    PAYMENT_AUTHORIZED,

    /**
     * Payment captured successfully
     */
    PAYMENT_CAPTURED,

    /**
     * Payment failed (declined, error, etc.)
     */
    PAYMENT_FAILED,

    /**
     * Payment refunded
     */
    PAYMENT_REFUNDED,

    // ==================== Completion Events ====================
    /**
     * Booking confirmed and finalized
     */
    CONFIRM,

    /**
     * Booking cancelled by user
     */
    CANCEL,

    /**
     * Booking cancelled by system (compensation)
     */
    SYSTEM_CANCEL,

    /**
     * Booking expired (timeout)
     */
    EXPIRE,

    /**
     * Generic failure
     */
    FAIL
}

package com.ticketing.booking.entity.enums;

/**
 * Booking lifecycle states
 *
 * Synchronous Flow (current):
 *   PENDING → RESERVING → AWAITING_PAYMENT → PROCESSING_PAYMENT → CONFIRMED
 *   Or: ... → FAILED/CANCELLED
 *
 * Async Flow (recommended):
 *   PENDING → SEATS_RESERVED → PAYMENT_PENDING → PAYMENT_PROCESSING → CONFIRMED
 *   Or: ... → PAYMENT_FAILED → (seats released) → CANCELLED
 */
public enum BookingStatus {
    /** Initial state - booking request received */
    PENDING,

    /** Processing - seat reservation in progress */
    RESERVING,

    /** Seats successfully reserved, ready for payment */
    SEATS_RESERVED,

    /** Seats reserved, waiting for payment (legacy - same as PAYMENT_PENDING) */
    AWAITING_PAYMENT,

    /** Payment intent created, waiting for capture/webhook */
    PAYMENT_PENDING,

    /** Payment processing in progress */
    PROCESSING_PAYMENT,

    /** Payment is being processed by gateway (async webhook pending) */
    PAYMENT_PROCESSING,

    /** Payment failed - seats will be released */
    PAYMENT_FAILED,

    /** Booking confirmed and finalized */
    CONFIRMED,

    /** Booking cancelled by user or system */
    CANCELLED,

    /** Booking failed due to error */
    FAILED,

    /** Booking expired (timeout) */
    EXPIRED
}


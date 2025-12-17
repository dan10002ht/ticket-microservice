package com.ticketing.booking.entity.enums;

/**
 * Payment lifecycle states
 *
 * Flow: NOT_REQUIRED (free events)
 *   Or: PENDING → PROCESSING → AUTHORIZED → CAPTURED
 *   Or: PENDING → PROCESSING → FAILED
 *   Or: CAPTURED → REFUNDING → REFUNDED
 */
public enum PaymentStatus {
    /** Payment not required (free event) */
    NOT_REQUIRED,

    /** Payment created but not yet processed */
    PENDING,

    /** Payment is being processed by gateway */
    PROCESSING,

    /** Payment authorized (funds held) */
    AUTHORIZED,

    /** Payment captured (funds transferred) */
    CAPTURED,

    /** Payment failed */
    FAILED,

    /** Refund in progress */
    REFUNDING,

    /** Payment refunded */
    REFUNDED
}


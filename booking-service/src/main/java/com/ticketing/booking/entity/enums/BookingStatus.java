package com.ticketing.booking.entity.enums;

/**
 * Booking lifecycle states
 * 
 * Flow: PENDING → RESERVING → AWAITING_PAYMENT → CONFIRMED
 *       Or: PENDING → RESERVING → AWAITING_PAYMENT → FAILED/CANCELLED
 */
public enum BookingStatus {
    /** Initial state - booking request received */
    PENDING,
    
    /** Processing - seat reservation in progress */
    RESERVING,
    
    /** Seats reserved, waiting for payment */
    AWAITING_PAYMENT,
    
    /** Payment processing in progress */
    PROCESSING_PAYMENT,
    
    /** Booking confirmed and finalized */
    CONFIRMED,
    
    /** Booking cancelled by user or system */
    CANCELLED,
    
    /** Booking failed due to error */
    FAILED,
    
    /** Booking expired (timeout) */
    EXPIRED
}


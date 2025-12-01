package com.ticketing.booking.service.saga;

/**
 * Saga step enumeration
 * Tracks which step of the saga is currently executing
 */
public enum SagaStep {
    /** Initial step - booking record created */
    CREATED,
    
    /** Step 1 - Acquiring distributed lock */
    LOCK_ACQUIRED,
    
    /** Step 2 - Reserving seats via Ticket Service */
    SEATS_RESERVED,
    
    /** Step 3 - Processing payment via Payment Service */
    PAYMENT_PROCESSED,
    
    /** Step 4 - Confirming booking */
    CONFIRMED,
    
    /** Compensation step - releasing resources */
    COMPENSATING,
    
    /** Saga completed (success or failure) */
    COMPLETED
}


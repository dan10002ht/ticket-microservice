package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Payment Status Enum
 * 
 * Represents the various states a payment can be in.
 */
@Getter
public enum PaymentStatus {
    PENDING("pending", "Pending"),
    PROCESSING("processing", "Processing"),
    SUCCESS("success", "Success"),
    FAILED("failed", "Failed"),
    CANCELLED("cancelled", "Cancelled"),
    REFUNDED("refunded", "Refunded"),
    PARTIALLY_REFUNDED("partially_refunded", "Partially Refunded");

    private final String value;
    private final String displayName;

    PaymentStatus(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static PaymentStatus fromValue(String value) {
        for (PaymentStatus status : PaymentStatus.values()) {
            if (status.value.equalsIgnoreCase(value)) {
                return status;
            }
        }
        throw new IllegalArgumentException("Invalid payment status: " + value);
    }

    /**
     * Check if status is terminal (cannot be changed)
     */
    public boolean isTerminal() {
        return this == SUCCESS || this == FAILED || this == CANCELLED || this == REFUNDED;
    }

    /**
     * Check if status allows refund
     */
    public boolean allowsRefund() {
        return this == SUCCESS || this == PARTIALLY_REFUNDED;
    }

    /**
     * Check if status allows cancellation
     */
    public boolean allowsCancellation() {
        return this == PENDING || this == PROCESSING;
    }
}



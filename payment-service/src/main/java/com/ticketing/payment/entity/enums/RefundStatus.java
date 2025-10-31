package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Refund Status Enum
 * 
 * Represents the various states a refund can be in.
 */
@Getter
public enum RefundStatus {
    PENDING("pending", "Pending"),
    PROCESSING("processing", "Processing"),
    SUCCESS("success", "Success"),
    FAILED("failed", "Failed"),
    CANCELLED("cancelled", "Cancelled");

    private final String value;
    private final String displayName;

    RefundStatus(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static RefundStatus fromValue(String value) {
        for (RefundStatus status : RefundStatus.values()) {
            if (status.value.equalsIgnoreCase(value)) {
                return status;
            }
        }
        throw new IllegalArgumentException("Invalid refund status: " + value);
    }

    /**
     * Check if status is terminal (cannot be changed)
     */
    public boolean isTerminal() {
        return this == SUCCESS || this == FAILED || this == CANCELLED;
    }

    /**
     * Check if status allows cancellation
     */
    public boolean allowsCancellation() {
        return this == PENDING || this == PROCESSING;
    }

    /**
     * Check if status is successful
     */
    public boolean isSuccessful() {
        return this == SUCCESS;
    }

    /**
     * Check if status indicates failure
     */
    public boolean hasFailled() {
        return this == FAILED || this == CANCELLED;
    }
}



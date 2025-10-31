package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Idempotency Status Enum
 * 
 * Represents the status of an idempotency key.
 */
@Getter
public enum IdempotencyStatus {
    PROCESSING("processing", "Processing"),
    COMPLETED("completed", "Completed"),
    FAILED("failed", "Failed");

    private final String value;
    private final String displayName;

    IdempotencyStatus(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static IdempotencyStatus fromValue(String value) {
        for (IdempotencyStatus status : IdempotencyStatus.values()) {
            if (status.value.equalsIgnoreCase(value)) {
                return status;
            }
        }
        throw new IllegalArgumentException("Invalid idempotency status: " + value);
    }

    /**
     * Check if status is processing
     */
    public boolean isProcessing() {
        return this == PROCESSING;
    }

    /**
     * Check if status is completed
     */
    public boolean isCompleted() {
        return this == COMPLETED;
    }

    /**
     * Check if status is failed
     */
    public boolean isFailed() {
        return this == FAILED;
    }

    /**
     * Check if status is terminal (cannot be changed)
     */
    public boolean isTerminal() {
        return this == COMPLETED || this == FAILED;
    }
}



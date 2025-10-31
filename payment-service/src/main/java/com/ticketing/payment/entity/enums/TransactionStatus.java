package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Transaction Status Enum
 * 
 * Represents the status of a transaction log entry.
 */
@Getter
public enum TransactionStatus {
    SUCCESS("success", "Success"),
    FAILED("failed", "Failed"),
    PENDING("pending", "Pending");

    private final String value;
    private final String displayName;

    TransactionStatus(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static TransactionStatus fromValue(String value) {
        for (TransactionStatus status : TransactionStatus.values()) {
            if (status.value.equalsIgnoreCase(value)) {
                return status;
            }
        }
        throw new IllegalArgumentException("Invalid transaction status: " + value);
    }

    /**
     * Check if status is successful
     */
    public boolean isSuccessful() {
        return this == SUCCESS;
    }

    /**
     * Check if status is failed
     */
    public boolean hasFailed() {
        return this == FAILED;
    }

    /**
     * Check if status is pending
     */
    public boolean isPending() {
        return this == PENDING;
    }
}



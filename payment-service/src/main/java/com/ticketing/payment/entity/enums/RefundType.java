package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Refund Type Enum
 * 
 * Represents the type of refund (full or partial).
 */
@Getter
public enum RefundType {
    FULL("full", "Full Refund"),
    PARTIAL("partial", "Partial Refund");

    private final String value;
    private final String displayName;

    RefundType(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static RefundType fromValue(String value) {
        for (RefundType type : RefundType.values()) {
            if (type.value.equalsIgnoreCase(value)) {
                return type;
            }
        }
        throw new IllegalArgumentException("Invalid refund type: " + value);
    }

    /**
     * Check if refund type is full
     */
    public boolean isFull() {
        return this == FULL;
    }

    /**
     * Check if refund type is partial
     */
    public boolean isPartial() {
        return this == PARTIAL;
    }
}



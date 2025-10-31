package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Payment Method Enum
 * 
 * Represents the various payment methods supported by the system.
 */
@Getter
public enum PaymentMethod {
    CREDIT_CARD("credit_card", "Credit Card"),
    DEBIT_CARD("debit_card", "Debit Card"),
    BANK_TRANSFER("bank_transfer", "Bank Transfer"),
    E_WALLET("e_wallet", "E-Wallet"),
    CASH("cash", "Cash");

    private final String value;
    private final String displayName;

    PaymentMethod(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static PaymentMethod fromValue(String value) {
        for (PaymentMethod method : PaymentMethod.values()) {
            if (method.value.equalsIgnoreCase(value)) {
                return method;
            }
        }
        throw new IllegalArgumentException("Invalid payment method: " + value);
    }

    /**
     * Check if payment method requires online gateway
     */
    public boolean requiresGateway() {
        return this == CREDIT_CARD || this == DEBIT_CARD || this == E_WALLET;
    }

    /**
     * Check if payment method is card-based
     */
    public boolean isCardBased() {
        return this == CREDIT_CARD || this == DEBIT_CARD;
    }

    /**
     * Check if payment method is instant
     */
    public boolean isInstant() {
        return this == CREDIT_CARD || this == DEBIT_CARD || this == E_WALLET;
    }
}



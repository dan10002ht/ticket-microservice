package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Payment Gateway Enum
 * 
 * Represents the various payment gateways supported by the system.
 */
@Getter
public enum PaymentGateway {
    STRIPE("stripe", "Stripe"),
    PAYPAL("paypal", "PayPal"),
    VNPAY("vnpay", "VNPay"),
    MOMO("momo", "Momo");

    private final String value;
    private final String displayName;

    PaymentGateway(String value, String displayName) {
        this.value = value;
        this.displayName = displayName;
    }

    /**
     * Get enum from string value
     */
    public static PaymentGateway fromValue(String value) {
        for (PaymentGateway gateway : PaymentGateway.values()) {
            if (gateway.value.equalsIgnoreCase(value)) {
                return gateway;
            }
        }
        throw new IllegalArgumentException("Invalid payment gateway: " + value);
    }

    /**
     * Check if gateway supports refunds
     */
    public boolean supportsRefunds() {
        return true; // All gateways support refunds
    }

    /**
     * Check if gateway supports webhooks
     */
    public boolean supportsWebhooks() {
        return true; // All gateways support webhooks
    }

    /**
     * Check if gateway is international
     */
    public boolean isInternational() {
        return this == STRIPE || this == PAYPAL;
    }

    /**
     * Check if gateway is domestic (Vietnam)
     */
    public boolean isDomestic() {
        return this == VNPAY || this == MOMO;
    }
}



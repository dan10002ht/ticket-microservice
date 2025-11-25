package com.ticketing.payment.entity.enums;

import lombok.Getter;

/**
 * Transaction Type Enum
 * 
 * Represents the various types of transactions that can be logged.
 */
@Getter
public enum TransactionType {
    PAYMENT_INITIATED("payment_initiated", "Payment Initiated", true, false),
    PAYMENT_SUCCESS("payment_success", "Payment Success", true, false),
    PAYMENT_FAILED("payment_failed", "Payment Failed", true, false),
    PAYMENT_CAPTURED("payment_captured", "Payment Captured", true, false),
    PAYMENT_CANCELLED("payment_cancelled", "Payment Cancelled", true, false),
    REFUND_INITIATED("refund_initiated", "Refund Initiated", false, true),
    REFUND_SUCCESS("refund_success", "Refund Success", false, true),
    REFUND_FAILED("refund_failed", "Refund Failed", false, true),
    WEBHOOK_RECEIVED("webhook_received", "Webhook Received", false, false);

    private final String value;
    private final String displayName;
    private final boolean paymentRelated;
    private final boolean refundRelated;

    TransactionType(String value, String displayName, boolean paymentRelated, boolean refundRelated) {
        this.value = value;
        this.displayName = displayName;
        this.paymentRelated = paymentRelated;
        this.refundRelated = refundRelated;
    }

    /**
     * Get enum from string value
     */
    public static TransactionType fromValue(String value) {
        for (TransactionType type : TransactionType.values()) {
            if (type.value.equalsIgnoreCase(value)) {
                return type;
            }
        }
        throw new IllegalArgumentException("Invalid transaction type: " + value);
    }

    /**
     * Check if type is payment related
     */
    public boolean isPaymentRelated() {
        return paymentRelated;
    }

    /**
     * Check if type is refund related
     */
    public boolean isRefundRelated() {
        return refundRelated;
    }

    /**
     * Check if type is webhook
     */
    public boolean isWebhook() {
        return this == WEBHOOK_RECEIVED;
    }

    /**
     * Check if type is success event
     */
    public boolean isSuccessEvent() {
        return this == PAYMENT_SUCCESS || this == REFUND_SUCCESS;
    }

    /**
     * Check if type is failure event
     */
    public boolean isFailureEvent() {
        return this == PAYMENT_FAILED || this == REFUND_FAILED;
    }

    /**
     * Check if type is initiation event
     */
    public boolean isInitiationEvent() {
        return this == PAYMENT_INITIATED || this == REFUND_INITIATED;
    }
}



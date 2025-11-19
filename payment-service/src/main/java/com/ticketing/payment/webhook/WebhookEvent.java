package com.ticketing.payment.webhook;

import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class WebhookEvent {
    EventType type;
    String providerReference;
    String paymentId;
    String refundId;
    String status;
    String rawPayload;

    public enum EventType {
        PAYMENT_SUCCEEDED,
        PAYMENT_FAILED,
        REFUND_SUCCEEDED,
        REFUND_FAILED
    }
}

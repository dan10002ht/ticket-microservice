package com.ticketing.payment.webhook;

import com.ticketing.payment.entity.enums.PaymentGateway;

public interface PaymentWebhookHandler {

    PaymentGateway getGateway();

    WebhookEvent handle(String payload, java.util.Map<String, String> headers);
}

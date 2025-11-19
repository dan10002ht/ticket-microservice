package com.ticketing.payment.webhook;

import com.ticketing.payment.entity.enums.PaymentGateway;

import java.util.EnumMap;
import java.util.Map;

public class WebhookHandlerRegistry {

    private final Map<PaymentGateway, PaymentWebhookHandler> handlers = new EnumMap<>(PaymentGateway.class);

    public void register(PaymentWebhookHandler handler) {
        handlers.put(handler.getGateway(), handler);
    }

    public PaymentWebhookHandler resolve(PaymentGateway gateway) {
        PaymentWebhookHandler handler = handlers.get(gateway);
        if (handler == null) {
            throw new IllegalArgumentException("No webhook handler registered for " + gateway);
        }
        return handler;
    }
}

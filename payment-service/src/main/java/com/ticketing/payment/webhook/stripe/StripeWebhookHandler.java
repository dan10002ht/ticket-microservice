package com.ticketing.payment.webhook.stripe;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.stripe.exception.SignatureVerificationException;
import com.stripe.net.Webhook;
import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.repository.PaymentRepository;
import com.ticketing.payment.repository.RefundRepository;
import com.ticketing.payment.webhook.PaymentWebhookHandler;
import com.ticketing.payment.webhook.WebhookEvent;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.util.Map;

@Component
@RequiredArgsConstructor
@Slf4j
public class StripeWebhookHandler implements PaymentWebhookHandler {

    private final StripeProperties properties;
    private final PaymentRepository paymentRepository;
    private final RefundRepository refundRepository;

    @Override
    public PaymentGateway getGateway() {
        return PaymentGateway.STRIPE;
    }

    @Override
    public WebhookEvent handle(String payload, Map<String, String> headers) {
        verifySignature(payload, headers);

        JsonObject eventJson = JsonParser.parseString(payload).getAsJsonObject();
        String type = eventJson.get("type").getAsString();
        JsonObject dataObject = eventJson.getAsJsonObject("data").getAsJsonObject("object");

        WebhookEvent.WebhookEventBuilder builder = WebhookEvent.builder()
                .rawPayload(payload);

        String providerRef = dataObject.get("id").getAsString();
        String paymentIntentId = dataObject.has("payment_intent") && !dataObject.get("payment_intent").isJsonNull()
                ? dataObject.get("payment_intent").getAsString()
                : providerRef;

        paymentRepository.findByProviderReference(paymentIntentId)
                .ifPresent(payment -> builder.paymentId(payment.getPaymentId().toString()));

        switch (type) {
            case "payment_intent.succeeded" -> builder
                    .type(WebhookEvent.EventType.PAYMENT_SUCCEEDED)
                    .providerReference(paymentIntentId)
                    .status("SUCCESS");
            case "payment_intent.payment_failed" -> builder
                    .type(WebhookEvent.EventType.PAYMENT_FAILED)
                    .providerReference(paymentIntentId)
                    .status("FAILED");
            case "charge.refunded", "charge.refund.updated" -> {
                builder.type(WebhookEvent.EventType.REFUND_SUCCEEDED)
                        .providerReference(paymentIntentId)
                        .status("SUCCESS");
                refundRepository.findByProviderReference(providerRef)
                        .ifPresent(refund -> builder.refundId(refund.getRefundId().toString()));
            }
            case "charge.refund.failed" -> {
                builder.type(WebhookEvent.EventType.REFUND_FAILED)
                        .providerReference(paymentIntentId)
                        .status("FAILED");
                refundRepository.findByProviderReference(providerRef)
                        .ifPresent(refund -> builder.refundId(refund.getRefundId().toString()));
            }
            default -> {
                log.info("Unhandled Stripe event type {}", type);
                builder.type(null);
            }
        }
        return builder.build();
    }

    private void verifySignature(String payload, Map<String, String> headers) {
        if (!StringUtils.hasText(properties.getWebhookSecret())) {
            return;
        }

        String signature = headers.entrySet().stream()
                .filter(entry -> entry.getKey() != null && entry.getKey().equalsIgnoreCase("stripe-signature"))
                .map(Map.Entry::getValue)
                .findFirst()
                .orElse(null);

        if (!StringUtils.hasText(signature)) {
            throw new IllegalArgumentException("Missing Stripe-Signature header");
        }

        try {
            Webhook.constructEvent(payload, signature, properties.getWebhookSecret());
        } catch (SignatureVerificationException e) {
            log.error("Invalid Stripe webhook signature", e);
            throw new IllegalArgumentException("Invalid Stripe webhook signature");
        }
    }
}

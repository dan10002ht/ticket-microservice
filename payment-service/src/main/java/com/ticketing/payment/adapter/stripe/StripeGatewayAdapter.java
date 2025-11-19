package com.ticketing.payment.adapter.stripe;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.Map;

import org.springframework.stereotype.Component;

import com.stripe.Stripe;
import com.stripe.exception.StripeException;
import com.stripe.model.PaymentIntent;
import com.stripe.model.Refund;
import com.stripe.param.PaymentIntentCreateParams;
import com.stripe.param.PaymentIntentCreateParams.CaptureMethod;
import com.stripe.param.PaymentIntentCreateParams.ConfirmationMethod;
import com.stripe.param.PaymentIntentCreateParams.PaymentMethodOptions;
import com.stripe.param.PaymentIntentCaptureParams;
import com.stripe.param.PaymentIntentCancelParams;
import com.stripe.param.RefundCreateParams;
import com.ticketing.payment.adapter.GatewayRequestContext;
import com.ticketing.payment.adapter.GatewayResponse;
import com.ticketing.payment.adapter.PaymentGatewayAdapter;
import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.service.dto.PaymentCreateCommand;
import com.ticketing.payment.service.dto.RefundCreateCommand;

import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Component
@RequiredArgsConstructor
@Slf4j
public class StripeGatewayAdapter implements PaymentGatewayAdapter {

    private final StripeProperties properties;

    @PostConstruct
    void init() {
        Stripe.apiKey = properties.getApiKey();
    }

    @Override
    public PaymentGateway getGateway() {
        return PaymentGateway.STRIPE;
    }

    @Override
    public GatewayResponse authorize(PaymentCreateCommand command, GatewayRequestContext context) {
        try {
            PaymentIntentCreateParams.Builder builder = PaymentIntentCreateParams.builder()
                    .setAmount(toStripeAmount(command.getAmount()))
                    .setCurrency(command.getCurrency().toLowerCase())
                    .setCaptureMethod(CaptureMethod.AUTOMATIC)
                    .setConfirmationMethod(ConfirmationMethod.AUTOMATIC)
                    .setConfirm(true)
                    .setDescription(command.getBookingId())
                    .putAllMetadata(toMetadata(command.getMetadata()))
                    .putMetadata("booking_id", command.getBookingId())
                    .putMetadata("user_id", command.getUserId())
                    .setPaymentMethodTypes(java.util.List.of("card"));

            if (command.getIdempotencyKey() != null) {
                builder.setIdempotencyKey(command.getIdempotencyKey());
            }

            PaymentIntent intent = PaymentIntent.create(builder.build());
            return GatewayResponse.builder()
                    .successful(true)
                    .providerReference(intent.getId())
                    .rawPayload(intent.toJson())
                    .build();
        } catch (StripeException e) {
            log.error("Stripe authorize failed", e);
            return GatewayResponse.builder()
                    .successful(false)
                    .errorCode(e.getCode())
                    .errorMessage(e.getMessage())
                    .build();
        }
    }

    @Override
    public GatewayResponse capture(PaymentCreateCommand command, GatewayRequestContext context) {
        String providerReference = context != null ? context.getProviderReference() : null;
        if (providerReference == null) {
            return GatewayResponse.builder()
                    .successful(false)
                    .errorCode("MISSING_REFERENCE")
                    .errorMessage("Provider reference required for capture")
                    .build();
        }
        try {
            PaymentIntent intent = PaymentIntent.retrieve(providerReference);
            intent.capture(PaymentIntentCaptureParams.builder().build());
            return GatewayResponse.builder()
                    .successful(true)
                    .providerReference(intent.getId())
                    .rawPayload(intent.toJson())
                    .build();
        } catch (StripeException e) {
            log.error("Stripe capture failed", e);
            return GatewayResponse.builder()
                    .successful(false)
                    .errorCode(e.getCode())
                    .errorMessage(e.getMessage())
                    .build();
        }
    }

    @Override
    public GatewayResponse refund(RefundCreateCommand command, GatewayRequestContext context) {
        String providerReference = context != null ? context.getProviderReference() : null;
        if (providerReference == null) {
            return GatewayResponse.builder()
                    .successful(false)
                    .errorCode("MISSING_REFERENCE")
                    .errorMessage("Provider reference required for refund")
                    .build();
        }
        try {
            RefundCreateParams params = RefundCreateParams.builder()
                    .setAmount(toStripeAmount(command.getAmount()))
                    .setPaymentIntent(providerReference)
                    .setReason(RefundCreateParams.Reason.REQUESTED_BY_CUSTOMER)
                    .putMetadata("reason", command.getReason())
                    .build();
            Refund refund = Refund.create(params);
            return GatewayResponse.builder()
                    .successful(true)
                    .providerReference(refund.getId())
                    .rawPayload(refund.toJson())
                    .build();
        } catch (StripeException e) {
            log.error("Stripe refund failed", e);
            return GatewayResponse.builder()
                    .successful(false)
                    .errorCode(e.getCode())
                    .errorMessage(e.getMessage())
                    .build();
        }
    }

    @Override
    public GatewayResponse cancelPayment(String providerReference, String reason, GatewayRequestContext context) {
        try {
            PaymentIntent intent = PaymentIntent.retrieve(providerReference);
            intent.cancel(PaymentIntentCancelParams.builder()
                    .setCancellationReason(PaymentIntentCancelParams.CancellationReason.REQUESTED_BY_CUSTOMER)
                    .build());
            return GatewayResponse.builder()
                    .successful(true)
                    .providerReference(intent.getId())
                    .rawPayload(intent.toJson())
                    .build();
        } catch (StripeException e) {
            log.error("Stripe cancel failed", e);
            return GatewayResponse.builder()
                    .successful(false)
                    .errorCode(e.getCode())
                    .errorMessage(e.getMessage())
                    .build();
        }
    }

    private long toStripeAmount(BigDecimal amount) {
        return amount.multiply(BigDecimal.valueOf(100)).longValue();
    }

    private Map<String, String> toMetadata(Map<String, Object> metadata) {
        Map<String, String> result = new HashMap<>();
        if (metadata == null) {
            return result;
        }
        metadata.forEach((key, value) -> result.put(key, value != null ? value.toString() : ""));
        return result;
    }
}

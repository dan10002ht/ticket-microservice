package com.ticketing.payment.adapter;

import com.ticketing.payment.service.dto.PaymentCreateCommand;
import com.ticketing.payment.service.dto.RefundCreateCommand;
import lombok.Builder;
import lombok.Value;

import java.util.Map;

/**
 * Provides additional context when invoking gateway adapters (headers, tracing
 * info, etc.).
 */
@Value
@Builder
public class GatewayRequestContext {
    Map<String, String> headers;
    String correlationId;
    String providerReference;
    PaymentCreateCommand paymentCommand;
    RefundCreateCommand refundCommand;
}

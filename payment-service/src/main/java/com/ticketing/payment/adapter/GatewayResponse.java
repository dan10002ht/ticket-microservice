package com.ticketing.payment.adapter;

import lombok.Builder;
import lombok.Value;

/**
 * Standardized response from payment gateway adapters.
 */
@Value
@Builder
public class GatewayResponse {
    boolean successful;
    String providerReference;
    String rawPayload;
    String errorCode;
    String errorMessage;
}

package com.ticketing.payment.service.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Value;

import java.util.Map;

/**
 * Represents request metadata needed for idempotency enforcement. Controllers
 * (REST/gRPC) can populate this from headers/body without leaking HTTP details
 * into the service layer.
 */
@Value
@Builder
public class IdempotencyKeyContext {

    @NotBlank
    String key;

    @NotBlank
    @Size(max = 500)
    String requestPath;

    @NotBlank
    @Size(max = 10)
    String requestMethod;

    String userId;
    String ipAddress;
    String userAgent;

    Map<String, Object> requestHeaders;
    Map<String, Object> requestBody;
}

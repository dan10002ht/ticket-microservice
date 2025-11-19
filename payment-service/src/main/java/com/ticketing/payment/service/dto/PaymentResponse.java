package com.ticketing.payment.service.dto;

import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.entity.enums.PaymentMethod;
import com.ticketing.payment.entity.enums.PaymentStatus;
import lombok.Builder;
import lombok.Value;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.Map;
import java.util.UUID;

/**
 * DTO returned by service layer representing current state of a payment.
 */
@Value
@Builder
public class PaymentResponse {
    Long id;
    UUID paymentId;
    String bookingId;
    String ticketId;
    String userId;
    BigDecimal amount;
    String currency;
    PaymentMethod paymentMethod;
    PaymentStatus status;
    PaymentGateway gatewayProvider;
    String externalReference;
    String failureReason;
    String description;
    Map<String, Object> metadata;
    String idempotencyKey;
    LocalDateTime createdAt;
    LocalDateTime updatedAt;
    LocalDateTime paidAt;
    LocalDateTime cancelledAt;
}

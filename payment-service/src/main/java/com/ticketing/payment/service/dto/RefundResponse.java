package com.ticketing.payment.service.dto;

import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.entity.enums.RefundStatus;
import com.ticketing.payment.entity.enums.RefundType;
import lombok.Builder;
import lombok.Value;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.Map;
import java.util.UUID;

/**
 * DTO returned for refund operations.
 */
@Value
@Builder
public class RefundResponse {
    Long id;
    UUID refundId;
    UUID paymentId;
    BigDecimal amount;
    String currency;
    RefundType refundType;
    RefundStatus status;
    PaymentGateway gatewayProvider;
    String reason;
    String description;
    String failureReason;
    Map<String, Object> metadata;
    LocalDateTime createdAt;
    LocalDateTime updatedAt;
    LocalDateTime refundedAt;
    LocalDateTime cancelledAt;
}

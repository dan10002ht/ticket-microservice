package com.ticketing.payment.service.mapper;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.Refund;
import com.ticketing.payment.service.dto.RefundCreateCommand;
import com.ticketing.payment.service.dto.RefundResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.util.Map;

/**
 * Mapper for refund entities and DTOs.
 */
@Component
@RequiredArgsConstructor
public class RefundMapper {

    private final ObjectMapper objectMapper;

    public Refund toEntity(Payment payment, RefundCreateCommand command) {
        Refund refund = Refund.builder()
                .payment(payment)
                .paymentUuid(payment.getPaymentId())
                .amount(command.getAmount())
                .currency(payment.getCurrency())
                .refundType(command.getRefundType())
                .reason(command.getReason())
                .description(command.getDescription())
                .build();

        if (command.getMetadata() != null) {
            refund.setMetadata(writeJson(command.getMetadata()));
        }
        return refund;
    }

    public RefundResponse toResponse(Refund refund) {
        return RefundResponse.builder()
                .id(refund.getId())
                .refundId(refund.getRefundId())
                .paymentId(refund.getPaymentUuid())
                .amount(refund.getAmount())
                .currency(refund.getCurrency())
                .refundType(refund.getRefundType())
                .status(refund.getStatus())
                .gatewayProvider(refund.getGatewayProvider())
                .reason(refund.getReason())
                .description(refund.getDescription())
                .failureReason(refund.getFailureReason())
                .metadata(readJson(refund.getMetadata()))
                .createdAt(refund.getCreatedAt())
                .updatedAt(refund.getUpdatedAt())
                .refundedAt(refund.getRefundedAt())
                .cancelledAt(refund.getCancelledAt())
                .build();
    }

    private String writeJson(Map<String, Object> metadata) {
        try {
            return objectMapper.writeValueAsString(metadata);
        } catch (JsonProcessingException e) {
            throw new IllegalArgumentException("Failed to serialize refund metadata", e);
        }
    }

    @SuppressWarnings("unchecked")
    private Map<String, Object> readJson(String metadata) {
        if (!StringUtils.hasText(metadata)) {
            return null;
        }
        try {
            return objectMapper.readValue(metadata, Map.class);
        } catch (JsonProcessingException e) {
            throw new IllegalStateException("Failed to deserialize refund metadata", e);
        }
    }
}

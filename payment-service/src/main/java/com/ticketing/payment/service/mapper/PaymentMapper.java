package com.ticketing.payment.service.mapper;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.service.dto.PaymentCreateCommand;
import com.ticketing.payment.service.dto.PaymentResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.util.Map;

/**
 * Mapper responsible for converting between entities and service DTOs.
 */
@Component
@RequiredArgsConstructor
public class PaymentMapper {

    private final ObjectMapper objectMapper;

    public Payment toEntity(PaymentCreateCommand command) {
        Payment payment = Payment.builder()
                .bookingId(command.getBookingId())
                .ticketId(command.getTicketId())
                .userId(command.getUserId())
                .amount(command.getAmount())
                .currency(command.getCurrency())
                .paymentMethod(command.getPaymentMethod())
                .gatewayProvider(command.getGatewayProvider())
                .idempotencyKey(command.getIdempotencyKey())
                .build();

        if (command.getMetadata() != null) {
            payment.setMetadata(writeJson(command.getMetadata()));
        }

        return payment;
    }

    public PaymentResponse toResponse(Payment payment) {
        return PaymentResponse.builder()
                .id(payment.getId())
                .paymentId(payment.getPaymentId())
                .bookingId(payment.getBookingId())
                .ticketId(payment.getTicketId())
                .userId(payment.getUserId())
                .amount(payment.getAmount())
                .currency(payment.getCurrency())
                .paymentMethod(payment.getPaymentMethod())
                .status(payment.getStatus())
                .gatewayProvider(payment.getGatewayProvider())
                .externalReference(payment.getExternalReference())
                .failureReason(payment.getFailureReason())
                .description(payment.getDescription())
                .metadata(readJson(payment.getMetadata()))
                .idempotencyKey(payment.getIdempotencyKey())
                .createdAt(payment.getCreatedAt())
                .updatedAt(payment.getUpdatedAt())
                .paidAt(payment.getPaidAt())
                .cancelledAt(payment.getCancelledAt())
                .build();
    }

    private String writeJson(Map<String, Object> metadata) {
        try {
            return objectMapper.writeValueAsString(metadata);
        } catch (JsonProcessingException e) {
            throw new IllegalArgumentException("Failed to serialize metadata", e);
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
            throw new IllegalStateException("Failed to deserialize metadata", e);
        }
    }
}

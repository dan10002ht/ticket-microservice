package com.ticketing.payment.service.dto;

import com.ticketing.payment.entity.enums.PaymentMethod;
import com.ticketing.payment.entity.enums.PaymentGateway;
import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.Digits;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Value;

import java.math.BigDecimal;
import java.util.Map;

/**
 * Command object representing input required to create a payment.
 *
 * Using service-level DTOs decouples controllers/gRPC adapters from the core
 * business logic and allows reuse across REST + gRPC interfaces.
 */
@Value
@Builder
public class PaymentCreateCommand {

    @NotBlank
    @Size(max = 255)
    String bookingId;

    @Size(max = 255)
    String ticketId;

    @NotBlank
    @Size(max = 255)
    String userId;

    @NotNull
    @DecimalMin(value = "0.01")
    @Digits(integer = 10, fraction = 2)
    BigDecimal amount;

    @NotBlank
    @Size(min = 3, max = 3)
    String currency;

    @NotNull
    PaymentMethod paymentMethod;

    @NotNull
    PaymentGateway gatewayProvider;

    @Size(max = 255)
    String idempotencyKey;

    Map<String, Object> metadata;
}

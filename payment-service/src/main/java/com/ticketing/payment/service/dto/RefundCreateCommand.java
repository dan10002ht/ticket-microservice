package com.ticketing.payment.service.dto;

import com.ticketing.payment.entity.enums.RefundType;
import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.Digits;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Value;

import java.math.BigDecimal;
import java.util.Map;
import java.util.UUID;

/**
 * Command representing data required to initiate a refund.
 */
@Value
@Builder
public class RefundCreateCommand {

    @NotNull
    UUID paymentId;

    @NotNull
    @DecimalMin(value = "0.01")
    @Digits(integer = 10, fraction = 2)
    BigDecimal amount;

    @NotBlank
    @Size(max = 255)
    String reason;

    @Size(max = 255)
    String description;

    @NotNull
    RefundType refundType;

    Map<String, Object> metadata;
}

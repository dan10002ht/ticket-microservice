package com.ticketing.booking.service.dto;

import java.math.BigDecimal;
import java.util.List;
import java.util.Map;

import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class BookingCreateCommand {

    @NotBlank
    String userId;

    @NotBlank
    String eventId;

    @NotEmpty
    List<String> seatNumbers;

    @DecimalMin("0.01")
    BigDecimal totalAmount;

    @NotBlank
    @Size(min = 3, max = 3)
    String currency;

    @NotNull
    Integer seatCount;

    Map<String, String> metadata;

    /**
     * Idempotency key to prevent duplicate bookings on retry.
     * Typically the queue item ID from booking-worker.
     */
    String idempotencyKey;
}


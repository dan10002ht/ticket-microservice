package com.ticketing.booking.dto.request;

import java.math.BigDecimal;
import java.util.List;
import java.util.Map;

import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class BookingRequest {

    @NotBlank
    private String userId;

    @NotBlank
    private String eventId;

    @NotEmpty
    private List<String> seatNumbers;

    @NotNull
    @DecimalMin("0.01")
    private BigDecimal totalAmount;

    @NotBlank
    @Size(min = 3, max = 3)
    private String currency;

    private Map<String, String> metadata;
}


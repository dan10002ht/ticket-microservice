package com.ticketing.booking.dto.request;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class BookingCancelRequest {

    @NotBlank
    private String reason;
}


package com.ticketing.booking.service.dto;

import java.util.UUID;

import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class BookingResult {

    UUID bookingId;
    String bookingReference;
    String status;
    String paymentStatus;
}


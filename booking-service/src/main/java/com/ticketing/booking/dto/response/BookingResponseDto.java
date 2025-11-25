package com.ticketing.booking.dto.response;

import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import lombok.Builder;
import lombok.Value;

@Value
@Builder
public class BookingResponseDto {

    UUID bookingId;
    String bookingReference;
    String userId;
    String eventId;
    int seatCount;
    List<String> seatNumbers;
    BigDecimal totalAmount;
    String currency;
    String status;
    String paymentStatus;
    String paymentReference;
    Map<String, String> metadata;
    OffsetDateTime expiresAt;
    OffsetDateTime confirmedAt;
    OffsetDateTime cancelledAt;
}


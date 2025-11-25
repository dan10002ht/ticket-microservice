package com.ticketing.booking.service.mapper;

import java.util.Map;

import org.springframework.stereotype.Component;

import com.ticketing.booking.dto.response.BookingResponseDto;
import com.ticketing.booking.entity.Booking;
import com.ticketing.booking.service.dto.BookingResult;

@Component
public class BookingMapper {

    public BookingResult toResult(Booking booking) {
        return BookingResult.builder()
                .bookingId(booking.getBookingId())
                .bookingReference(booking.getBookingReference())
                .status(booking.getStatus().name())
                .paymentStatus(booking.getPaymentStatus().name())
                .build();
    }

    public BookingResponseDto toResponse(Booking booking) {
        return BookingResponseDto.builder()
                .bookingId(booking.getBookingId())
                .bookingReference(booking.getBookingReference())
                .userId(booking.getUserId())
                .eventId(booking.getEventId())
                .seatCount(booking.getSeatCount())
                .seatNumbers(booking.getSeatNumbers())
                .totalAmount(booking.getTotalAmount())
                .currency(booking.getCurrency())
                .status(booking.getStatus().name())
                .paymentStatus(booking.getPaymentStatus().name())
                .paymentReference(booking.getPaymentReference())
                .metadata(booking.getMetadata() == null ? Map.of() : booking.getMetadata())
                .expiresAt(booking.getExpiresAt())
                .confirmedAt(booking.getConfirmedAt())
                .cancelledAt(booking.getCancelledAt())
                .build();
    }
}


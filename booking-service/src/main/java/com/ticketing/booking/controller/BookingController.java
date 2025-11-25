package com.ticketing.booking.controller;

import java.util.UUID;

import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.ticketing.booking.dto.request.BookingCancelRequest;
import com.ticketing.booking.dto.request.BookingRequest;
import com.ticketing.booking.dto.response.BookingResponseDto;
import com.ticketing.booking.service.BookingService;
import com.ticketing.booking.service.dto.BookingCreateCommand;
import com.ticketing.booking.service.mapper.BookingMapper;

import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/api/bookings")
@RequiredArgsConstructor
public class BookingController {

    private final BookingService bookingService;
    private final BookingMapper bookingMapper;

    @PostMapping
    public ResponseEntity<BookingResponseDto> createBooking(@Validated @RequestBody BookingRequest request) {
        BookingCreateCommand command = BookingCreateCommand.builder()
                .userId(request.getUserId())
                .eventId(request.getEventId())
                .seatNumbers(request.getSeatNumbers())
                .totalAmount(request.getTotalAmount())
                .currency(request.getCurrency())
                .seatCount(request.getSeatNumbers().size())
                .metadata(request.getMetadata())
                .build();
        var result = bookingService.createBooking(command);
        return ResponseEntity.ok(bookingMapper.toResponse(bookingService.getBooking(result.getBookingId())));
    }

    @GetMapping("/{bookingId}")
    public ResponseEntity<BookingResponseDto> getBooking(@PathVariable UUID bookingId) {
        return ResponseEntity.ok(bookingMapper.toResponse(bookingService.getBooking(bookingId)));
    }

    @PostMapping("/{bookingId}/cancel")
    public ResponseEntity<BookingResponseDto> cancelBooking(@PathVariable UUID bookingId,
            @Validated @RequestBody BookingCancelRequest request) {
        bookingService.cancelBooking(bookingId, request.getReason());
        return ResponseEntity.ok(bookingMapper.toResponse(bookingService.getBooking(bookingId)));
    }
}


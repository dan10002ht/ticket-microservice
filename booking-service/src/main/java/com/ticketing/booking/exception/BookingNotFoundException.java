package com.ticketing.booking.exception;

import java.util.UUID;

public class BookingNotFoundException extends BookingException {

    public BookingNotFoundException(UUID bookingId) {
        super("Booking not found: " + bookingId);
    }

    public BookingNotFoundException(String bookingReference) {
        super("Booking not found with reference: " + bookingReference);
    }
}


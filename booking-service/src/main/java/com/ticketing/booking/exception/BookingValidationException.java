package com.ticketing.booking.exception;

public class BookingValidationException extends BookingException {

    public BookingValidationException(String message) {
        super("Validation failed: " + message);
    }
}


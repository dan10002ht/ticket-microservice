package com.ticketing.booking.exception;

public class BookingLockException extends BookingException {

    public BookingLockException(String message) {
        super("Failed to acquire booking lock: " + message);
    }

    public BookingLockException(String message, Throwable cause) {
        super("Failed to acquire booking lock: " + message, cause);
    }
}


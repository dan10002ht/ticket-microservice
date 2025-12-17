package com.ticketing.booking.exception;

import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.statemachine.BookingEvent;

/**
 * Exception thrown when an invalid state transition is attempted.
 */
public class InvalidStateTransitionException extends BookingException {

    private final BookingStatus fromState;
    private final BookingStatus toState;
    private final BookingEvent event;

    public InvalidStateTransitionException(String message) {
        super(message);
        this.fromState = null;
        this.toState = null;
        this.event = null;
    }

    public InvalidStateTransitionException(BookingStatus fromState, BookingStatus toState) {
        super(String.format("Invalid state transition from %s to %s", fromState, toState));
        this.fromState = fromState;
        this.toState = toState;
        this.event = null;
    }

    public InvalidStateTransitionException(BookingStatus fromState, BookingEvent event) {
        super(String.format("Cannot process event %s in state %s", event, fromState));
        this.fromState = fromState;
        this.toState = null;
        this.event = event;
    }

    public InvalidStateTransitionException(BookingStatus fromState, BookingStatus toState, BookingEvent event) {
        super(String.format("Invalid transition from %s to %s via event %s", fromState, toState, event));
        this.fromState = fromState;
        this.toState = toState;
        this.event = event;
    }

    public BookingStatus getFromState() {
        return fromState;
    }

    public BookingStatus getToState() {
        return toState;
    }

    public BookingEvent getEvent() {
        return event;
    }
}

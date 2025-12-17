package com.ticketing.booking.statemachine;

import java.util.ArrayList;
import java.util.EnumMap;
import java.util.EnumSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.UUID;

import org.springframework.stereotype.Component;

import com.ticketing.booking.entity.Booking;
import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.exception.InvalidStateTransitionException;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;

/**
 * Booking State Machine - manages valid state transitions for bookings.
 *
 * State Diagram:
 * <pre>
 *                     ┌────────────────────────────────────────┐
 *                     │                                        │
 *                     ▼                                        │
 * PENDING ──► RESERVING ──► SEATS_RESERVED ──► PAYMENT_PENDING │
 *    │            │              │                   │         │
 *    │            │              │                   ▼         │
 *    │            │              │         PAYMENT_PROCESSING ─┘
 *    │            │              │                   │
 *    ▼            ▼              ▼                   ▼
 * FAILED ◄──── FAILED ◄───── FAILED ◄────── PAYMENT_FAILED
 *                                                    │
 *    ┌───────────────────────────────────────────────┘
 *    │
 *    ▼
 * CANCELLED ◄────────────────── CONFIRMED ◄─── PAYMENT_PROCESSING
 *    │                              │
 *    │                              │
 *    └──────────────────────────────┘
 *              (user cancel)
 * </pre>
 */
@Component
@Slf4j
public class BookingStateMachine {

    private final Map<BookingStatus, Set<BookingStatus>> validTransitions;
    private final Map<BookingStatus, Map<BookingEvent, BookingStatus>> eventTransitions;
    private final Counter transitionCounter;
    private final Counter invalidTransitionCounter;
    private final List<StateTransitionListener> listeners;

    public BookingStateMachine(MeterRegistry meterRegistry, List<StateTransitionListener> listeners) {
        this.validTransitions = buildValidTransitions();
        this.eventTransitions = buildEventTransitions();
        this.transitionCounter = meterRegistry.counter("booking.state.transitions");
        this.invalidTransitionCounter = meterRegistry.counter("booking.state.invalid_transitions");
        this.listeners = listeners != null ? listeners : new ArrayList<>();
    }

    @PostConstruct
    public void init() {
        log.info("BookingStateMachine initialized with {} listeners", listeners.size());
    }

    /**
     * Add a listener for state transitions.
     */
    public void addListener(StateTransitionListener listener) {
        listeners.add(listener);
    }

    /**
     * Notify all listeners of a state transition.
     */
    private void notifyListeners(StateTransition transition) {
        for (StateTransitionListener listener : listeners) {
            try {
                listener.onTransition(transition);
            } catch (Exception e) {
                log.error("Error notifying listener of state transition", e);
            }
        }
    }

    /**
     * Build the map of valid state transitions.
     */
    private Map<BookingStatus, Set<BookingStatus>> buildValidTransitions() {
        Map<BookingStatus, Set<BookingStatus>> transitions = new EnumMap<>(BookingStatus.class);

        // From PENDING
        transitions.put(BookingStatus.PENDING, EnumSet.of(
                BookingStatus.RESERVING,
                BookingStatus.FAILED,
                BookingStatus.CANCELLED,
                BookingStatus.EXPIRED
        ));

        // From RESERVING
        transitions.put(BookingStatus.RESERVING, EnumSet.of(
                BookingStatus.SEATS_RESERVED,
                BookingStatus.AWAITING_PAYMENT, // Legacy compatibility
                BookingStatus.FAILED,
                BookingStatus.CANCELLED
        ));

        // From SEATS_RESERVED
        transitions.put(BookingStatus.SEATS_RESERVED, EnumSet.of(
                BookingStatus.PAYMENT_PENDING,
                BookingStatus.AWAITING_PAYMENT, // Legacy compatibility
                BookingStatus.CONFIRMED, // For free events
                BookingStatus.FAILED,
                BookingStatus.CANCELLED,
                BookingStatus.EXPIRED
        ));

        // From AWAITING_PAYMENT (legacy)
        transitions.put(BookingStatus.AWAITING_PAYMENT, EnumSet.of(
                BookingStatus.PROCESSING_PAYMENT,
                BookingStatus.PAYMENT_PENDING,
                BookingStatus.PAYMENT_PROCESSING,
                BookingStatus.CONFIRMED,
                BookingStatus.FAILED,
                BookingStatus.CANCELLED,
                BookingStatus.EXPIRED
        ));

        // From PAYMENT_PENDING
        transitions.put(BookingStatus.PAYMENT_PENDING, EnumSet.of(
                BookingStatus.PAYMENT_PROCESSING,
                BookingStatus.CONFIRMED,
                BookingStatus.PAYMENT_FAILED,
                BookingStatus.FAILED,
                BookingStatus.CANCELLED,
                BookingStatus.EXPIRED
        ));

        // From PROCESSING_PAYMENT (legacy)
        transitions.put(BookingStatus.PROCESSING_PAYMENT, EnumSet.of(
                BookingStatus.CONFIRMED,
                BookingStatus.PAYMENT_FAILED,
                BookingStatus.FAILED,
                BookingStatus.CANCELLED
        ));

        // From PAYMENT_PROCESSING
        transitions.put(BookingStatus.PAYMENT_PROCESSING, EnumSet.of(
                BookingStatus.CONFIRMED,
                BookingStatus.PAYMENT_FAILED,
                BookingStatus.FAILED,
                BookingStatus.CANCELLED
        ));

        // From PAYMENT_FAILED
        transitions.put(BookingStatus.PAYMENT_FAILED, EnumSet.of(
                BookingStatus.CANCELLED,
                BookingStatus.FAILED
        ));

        // From CONFIRMED (terminal state, limited transitions)
        transitions.put(BookingStatus.CONFIRMED, EnumSet.of(
                BookingStatus.CANCELLED // User can cancel confirmed booking
        ));

        // Terminal states - no valid transitions out
        transitions.put(BookingStatus.CANCELLED, EnumSet.noneOf(BookingStatus.class));
        transitions.put(BookingStatus.FAILED, EnumSet.noneOf(BookingStatus.class));
        transitions.put(BookingStatus.EXPIRED, EnumSet.noneOf(BookingStatus.class));

        return transitions;
    }

    /**
     * Build the map of event-driven state transitions.
     */
    private Map<BookingStatus, Map<BookingEvent, BookingStatus>> buildEventTransitions() {
        Map<BookingStatus, Map<BookingEvent, BookingStatus>> transitions = new EnumMap<>(BookingStatus.class);

        // PENDING state transitions
        Map<BookingEvent, BookingStatus> pendingTransitions = new EnumMap<>(BookingEvent.class);
        pendingTransitions.put(BookingEvent.RESERVE_SEATS, BookingStatus.RESERVING);
        pendingTransitions.put(BookingEvent.FAIL, BookingStatus.FAILED);
        pendingTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        pendingTransitions.put(BookingEvent.EXPIRE, BookingStatus.EXPIRED);
        transitions.put(BookingStatus.PENDING, pendingTransitions);

        // RESERVING state transitions
        Map<BookingEvent, BookingStatus> reservingTransitions = new EnumMap<>(BookingEvent.class);
        reservingTransitions.put(BookingEvent.SEATS_RESERVED, BookingStatus.SEATS_RESERVED);
        reservingTransitions.put(BookingEvent.SEATS_RESERVATION_FAILED, BookingStatus.FAILED);
        reservingTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        transitions.put(BookingStatus.RESERVING, reservingTransitions);

        // SEATS_RESERVED state transitions
        Map<BookingEvent, BookingStatus> seatsReservedTransitions = new EnumMap<>(BookingEvent.class);
        seatsReservedTransitions.put(BookingEvent.REQUEST_PAYMENT, BookingStatus.PAYMENT_PENDING);
        seatsReservedTransitions.put(BookingEvent.CONFIRM, BookingStatus.CONFIRMED); // Free events
        seatsReservedTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        seatsReservedTransitions.put(BookingEvent.EXPIRE, BookingStatus.EXPIRED);
        transitions.put(BookingStatus.SEATS_RESERVED, seatsReservedTransitions);

        // PAYMENT_PENDING state transitions
        Map<BookingEvent, BookingStatus> paymentPendingTransitions = new EnumMap<>(BookingEvent.class);
        paymentPendingTransitions.put(BookingEvent.PAYMENT_AUTHORIZED, BookingStatus.PAYMENT_PROCESSING);
        paymentPendingTransitions.put(BookingEvent.PAYMENT_CAPTURED, BookingStatus.CONFIRMED);
        paymentPendingTransitions.put(BookingEvent.PAYMENT_FAILED, BookingStatus.PAYMENT_FAILED);
        paymentPendingTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        paymentPendingTransitions.put(BookingEvent.EXPIRE, BookingStatus.EXPIRED);
        transitions.put(BookingStatus.PAYMENT_PENDING, paymentPendingTransitions);

        // PAYMENT_PROCESSING state transitions
        Map<BookingEvent, BookingStatus> paymentProcessingTransitions = new EnumMap<>(BookingEvent.class);
        paymentProcessingTransitions.put(BookingEvent.PAYMENT_CAPTURED, BookingStatus.CONFIRMED);
        paymentProcessingTransitions.put(BookingEvent.PAYMENT_FAILED, BookingStatus.PAYMENT_FAILED);
        paymentProcessingTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        transitions.put(BookingStatus.PAYMENT_PROCESSING, paymentProcessingTransitions);

        // PAYMENT_FAILED state transitions
        Map<BookingEvent, BookingStatus> paymentFailedTransitions = new EnumMap<>(BookingEvent.class);
        paymentFailedTransitions.put(BookingEvent.SYSTEM_CANCEL, BookingStatus.CANCELLED);
        transitions.put(BookingStatus.PAYMENT_FAILED, paymentFailedTransitions);

        // CONFIRMED state transitions
        Map<BookingEvent, BookingStatus> confirmedTransitions = new EnumMap<>(BookingEvent.class);
        confirmedTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        transitions.put(BookingStatus.CONFIRMED, confirmedTransitions);

        // Legacy: AWAITING_PAYMENT state transitions
        Map<BookingEvent, BookingStatus> awaitingPaymentTransitions = new EnumMap<>(BookingEvent.class);
        awaitingPaymentTransitions.put(BookingEvent.PAYMENT_AUTHORIZED, BookingStatus.PAYMENT_PROCESSING);
        awaitingPaymentTransitions.put(BookingEvent.PAYMENT_CAPTURED, BookingStatus.CONFIRMED);
        awaitingPaymentTransitions.put(BookingEvent.PAYMENT_FAILED, BookingStatus.PAYMENT_FAILED);
        awaitingPaymentTransitions.put(BookingEvent.CANCEL, BookingStatus.CANCELLED);
        awaitingPaymentTransitions.put(BookingEvent.EXPIRE, BookingStatus.EXPIRED);
        transitions.put(BookingStatus.AWAITING_PAYMENT, awaitingPaymentTransitions);

        // Legacy: PROCESSING_PAYMENT state transitions
        Map<BookingEvent, BookingStatus> processingPaymentTransitions = new EnumMap<>(BookingEvent.class);
        processingPaymentTransitions.put(BookingEvent.PAYMENT_CAPTURED, BookingStatus.CONFIRMED);
        processingPaymentTransitions.put(BookingEvent.PAYMENT_FAILED, BookingStatus.PAYMENT_FAILED);
        transitions.put(BookingStatus.PROCESSING_PAYMENT, processingPaymentTransitions);

        return transitions;
    }

    /**
     * Check if a transition from one state to another is valid.
     */
    public boolean isValidTransition(BookingStatus from, BookingStatus to) {
        Set<BookingStatus> validTargets = validTransitions.get(from);
        return validTargets != null && validTargets.contains(to);
    }

    /**
     * Get the target state for a given event from a given state.
     * Returns null if the event is not valid for the current state.
     */
    public BookingStatus getTargetState(BookingStatus currentState, BookingEvent event) {
        Map<BookingEvent, BookingStatus> stateEvents = eventTransitions.get(currentState);
        if (stateEvents == null) {
            return null;
        }
        return stateEvents.get(event);
    }

    /**
     * Transition a booking to a new state based on an event.
     *
     * @param booking The booking to transition
     * @param event The event triggering the transition
     * @param triggeredBy Who/what triggered this transition
     * @param reason Optional reason for the transition
     * @return StateTransition record of what happened
     * @throws InvalidStateTransitionException if the transition is not valid
     */
    public StateTransition transition(
            Booking booking,
            BookingEvent event,
            String triggeredBy,
            String reason) {

        BookingStatus currentState = booking.getStatus();
        BookingStatus targetState = getTargetState(currentState, event);

        if (targetState == null) {
            invalidTransitionCounter.increment();
            log.warn("Invalid state transition attempted: booking={}, from={}, event={}",
                    booking.getBookingId(), currentState, event);
            throw new InvalidStateTransitionException(
                    String.format("Cannot transition from %s with event %s", currentState, event));
        }

        // Perform the transition
        BookingStatus previousState = booking.getStatus();
        booking.setStatus(targetState);

        // Record the transition
        StateTransition transition = StateTransition.create(
                booking.getBookingId(),
                previousState,
                targetState,
                event,
                triggeredBy,
                reason
        );

        transitionCounter.increment();
        log.info("State transition: booking={}, from={}, to={}, event={}, triggeredBy={}",
                booking.getBookingId(), previousState, targetState, event, triggeredBy);

        // Notify listeners (async persistence)
        notifyListeners(transition);

        return transition;
    }

    /**
     * Transition a booking directly to a new state (without event).
     * Use this for backwards compatibility with existing code.
     *
     * @param booking The booking to transition
     * @param targetState The target state
     * @param triggeredBy Who/what triggered this transition
     * @param reason Optional reason for the transition
     * @return StateTransition record of what happened
     * @throws InvalidStateTransitionException if the transition is not valid
     */
    public StateTransition transitionTo(
            Booking booking,
            BookingStatus targetState,
            String triggeredBy,
            String reason) {

        BookingStatus currentState = booking.getStatus();

        if (!isValidTransition(currentState, targetState)) {
            invalidTransitionCounter.increment();
            log.warn("Invalid state transition attempted: booking={}, from={}, to={}",
                    booking.getBookingId(), currentState, targetState);
            throw new InvalidStateTransitionException(
                    String.format("Cannot transition from %s to %s", currentState, targetState));
        }

        // Infer the event from the transition
        BookingEvent event = inferEvent(currentState, targetState);

        // Perform the transition
        booking.setStatus(targetState);

        // Record the transition
        StateTransition transition = StateTransition.create(
                booking.getBookingId(),
                currentState,
                targetState,
                event,
                triggeredBy,
                reason
        );

        transitionCounter.increment();
        log.info("State transition: booking={}, from={}, to={}, triggeredBy={}",
                booking.getBookingId(), currentState, targetState, triggeredBy);

        // Notify listeners (async persistence)
        notifyListeners(transition);

        return transition;
    }

    /**
     * Infer the event from a state transition.
     */
    private BookingEvent inferEvent(BookingStatus from, BookingStatus to) {
        return switch (to) {
            case RESERVING -> BookingEvent.RESERVE_SEATS;
            case SEATS_RESERVED -> BookingEvent.SEATS_RESERVED;
            case AWAITING_PAYMENT, PAYMENT_PENDING -> BookingEvent.REQUEST_PAYMENT;
            case PROCESSING_PAYMENT, PAYMENT_PROCESSING -> BookingEvent.PAYMENT_AUTHORIZED;
            case CONFIRMED -> BookingEvent.PAYMENT_CAPTURED;
            case PAYMENT_FAILED -> BookingEvent.PAYMENT_FAILED;
            case CANCELLED -> BookingEvent.CANCEL;
            case FAILED -> BookingEvent.FAIL;
            case EXPIRED -> BookingEvent.EXPIRE;
            default -> BookingEvent.FAIL;
        };
    }

    /**
     * Get all valid next states from the current state.
     */
    public Set<BookingStatus> getValidNextStates(BookingStatus currentState) {
        return validTransitions.getOrDefault(currentState, EnumSet.noneOf(BookingStatus.class));
    }

    /**
     * Get all valid events for the current state.
     */
    public Set<BookingEvent> getValidEvents(BookingStatus currentState) {
        Map<BookingEvent, BookingStatus> stateEvents = eventTransitions.get(currentState);
        if (stateEvents == null) {
            return EnumSet.noneOf(BookingEvent.class);
        }
        return stateEvents.keySet();
    }

    /**
     * Check if the state is terminal (no further transitions possible).
     */
    public boolean isTerminalState(BookingStatus state) {
        return state == BookingStatus.CANCELLED
                || state == BookingStatus.FAILED
                || state == BookingStatus.EXPIRED;
    }

    /**
     * Check if the state is a success state.
     */
    public boolean isSuccessState(BookingStatus state) {
        return state == BookingStatus.CONFIRMED;
    }
}

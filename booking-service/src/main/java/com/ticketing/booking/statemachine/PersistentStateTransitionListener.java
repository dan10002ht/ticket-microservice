package com.ticketing.booking.statemachine;

import org.springframework.stereotype.Component;

import com.ticketing.booking.service.StateTransitionService;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

/**
 * Listener that persists state transitions to the database.
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class PersistentStateTransitionListener implements StateTransitionListener {

    private final StateTransitionService stateTransitionService;

    @Override
    public void onTransition(StateTransition transition) {
        try {
            stateTransitionService.recordTransition(transition);
        } catch (Exception e) {
            // Log but don't fail the transition
            log.error("Failed to persist state transition: {} -> {} for booking {}",
                    transition.getFromState(), transition.getToState(), transition.getBookingId(), e);
        }
    }
}

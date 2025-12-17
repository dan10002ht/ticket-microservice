package com.ticketing.booking.statemachine;

/**
 * Listener interface for state transitions.
 * Allows external components to react to state changes.
 */
public interface StateTransitionListener {

    /**
     * Called after a successful state transition.
     *
     * @param transition The transition that occurred
     */
    void onTransition(StateTransition transition);
}

package com.ticketing.booking.service;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import com.ticketing.booking.entity.BookingStateTransition;
import com.ticketing.booking.repository.BookingStateTransitionRepository;
import com.ticketing.booking.statemachine.StateTransition;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

/**
 * Service for managing booking state transitions.
 * Persists transitions for audit trail and debugging.
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class StateTransitionService {

    private final BookingStateTransitionRepository transitionRepository;

    /**
     * Record a state transition.
     * Uses REQUIRES_NEW to ensure transition is saved even if outer transaction fails.
     */
    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public void recordTransition(StateTransition transition) {
        BookingStateTransition entity = BookingStateTransition.builder()
                .transitionId(transition.getTransitionId())
                .bookingId(transition.getBookingId())
                .fromState(transition.getFromState())
                .toState(transition.getToState())
                .event(transition.getEvent())
                .triggeredBy(transition.getTriggeredBy())
                .reason(transition.getReason())
                .createdAt(transition.getTimestamp())
                .build();

        transitionRepository.save(entity);
        log.debug("Recorded state transition: {} -> {} for booking {}",
                transition.getFromState(), transition.getToState(), transition.getBookingId());
    }

    /**
     * Get all transitions for a booking.
     */
    @Transactional(readOnly = true)
    public List<BookingStateTransition> getTransitionHistory(UUID bookingId) {
        return transitionRepository.findByBookingIdOrderByCreatedAtAsc(bookingId);
    }

    /**
     * Get the latest transition for a booking.
     */
    @Transactional(readOnly = true)
    public BookingStateTransition getLatestTransition(UUID bookingId) {
        return transitionRepository.findLatestByBookingId(bookingId);
    }

    /**
     * Cleanup old transitions (keep 90 days).
     * Runs daily at 3 AM.
     */
    @Scheduled(cron = "0 0 3 * * ?")
    @Transactional
    public void cleanupOldTransitions() {
        Instant cutoff = Instant.now().minusSeconds(90 * 24 * 60 * 60L); // 90 days
        int deleted = transitionRepository.deleteOlderThan(cutoff);
        log.info("Cleaned up {} old state transitions (older than 90 days)", deleted);
    }
}

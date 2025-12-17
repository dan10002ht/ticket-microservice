package com.ticketing.booking.repository;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.ticketing.booking.entity.BookingStateTransition;
import com.ticketing.booking.entity.enums.BookingStatus;

/**
 * Repository for booking state transitions.
 */
@Repository
public interface BookingStateTransitionRepository extends JpaRepository<BookingStateTransition, Long> {

    /**
     * Find all transitions for a booking ordered by time.
     */
    List<BookingStateTransition> findByBookingIdOrderByCreatedAtAsc(UUID bookingId);

    /**
     * Find the latest transition for a booking.
     */
    @Query("SELECT t FROM BookingStateTransition t WHERE t.bookingId = :bookingId ORDER BY t.createdAt DESC LIMIT 1")
    BookingStateTransition findLatestByBookingId(@Param("bookingId") UUID bookingId);

    /**
     * Count transitions by target state (for metrics).
     */
    long countByToState(BookingStatus toState);

    /**
     * Find transitions within time range (for audit/debugging).
     */
    @Query("SELECT t FROM BookingStateTransition t WHERE t.createdAt BETWEEN :start AND :end ORDER BY t.createdAt ASC")
    List<BookingStateTransition> findByTimeRange(
            @Param("start") Instant start,
            @Param("end") Instant end);

    /**
     * Count transitions by from and to state (for metrics).
     */
    @Query("SELECT COUNT(t) FROM BookingStateTransition t WHERE t.fromState = :fromState AND t.toState = :toState")
    long countByTransition(
            @Param("fromState") BookingStatus fromState,
            @Param("toState") BookingStatus toState);

    /**
     * Delete old transitions (for cleanup).
     */
    @Modifying
    @Query("DELETE FROM BookingStateTransition t WHERE t.createdAt < :cutoff")
    int deleteOlderThan(@Param("cutoff") Instant cutoff);
}

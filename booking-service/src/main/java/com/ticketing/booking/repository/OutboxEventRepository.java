package com.ticketing.booking.repository;

import java.time.OffsetDateTime;
import java.util.List;
import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.ticketing.booking.entity.OutboxEvent;
import com.ticketing.booking.entity.OutboxEvent.AggregateType;
import com.ticketing.booking.entity.OutboxEvent.OutboxStatus;

/**
 * Repository for OutboxEvent entity.
 * Provides methods for polling, publishing, and cleanup operations.
 */
@Repository
public interface OutboxEventRepository extends JpaRepository<OutboxEvent, UUID> {

    /**
     * Find pending events ordered by creation time (FIFO processing)
     * Uses SELECT FOR UPDATE SKIP LOCKED for concurrent processor support
     */
    @Query(value = """
        SELECT * FROM outbox_events
        WHERE status = 'PENDING'
        ORDER BY created_at ASC
        LIMIT :limit
        FOR UPDATE SKIP LOCKED
        """, nativeQuery = true)
    List<OutboxEvent> findPendingEventsForProcessing(@Param("limit") int limit);

    /**
     * Find pending events (without lock, for monitoring)
     */
    List<OutboxEvent> findByStatusOrderByCreatedAtAsc(OutboxStatus status);

    /**
     * Find failed events that can be retried
     */
    @Query("SELECT e FROM OutboxEvent e WHERE e.status = 'FAILED' AND e.retryCount < 5 ORDER BY e.createdAt ASC")
    List<OutboxEvent> findRetryableFailedEvents();

    /**
     * Count pending events (for metrics)
     */
    long countByStatus(OutboxStatus status);

    /**
     * Find events by aggregate (for debugging)
     */
    List<OutboxEvent> findByAggregateTypeAndAggregateIdOrderByCreatedAtDesc(
            AggregateType aggregateType, String aggregateId);

    /**
     * Mark event as published
     */
    @Modifying
    @Query("UPDATE OutboxEvent e SET e.status = 'PUBLISHED', e.publishedAt = :publishedAt WHERE e.id = :id")
    int markAsPublished(@Param("id") UUID id, @Param("publishedAt") OffsetDateTime publishedAt);

    /**
     * Mark event as failed
     */
    @Modifying
    @Query("UPDATE OutboxEvent e SET e.status = 'FAILED', e.lastError = :error, e.retryCount = e.retryCount + 1 WHERE e.id = :id")
    int markAsFailed(@Param("id") UUID id, @Param("error") String error);

    /**
     * Reset failed events for retry
     */
    @Modifying
    @Query("UPDATE OutboxEvent e SET e.status = 'PENDING' WHERE e.status = 'FAILED' AND e.retryCount < 5")
    int resetFailedEventsForRetry();

    /**
     * Delete old published events (cleanup)
     * Keep events for 7 days after publishing for debugging
     */
    @Modifying
    @Query("DELETE FROM OutboxEvent e WHERE e.status = 'PUBLISHED' AND e.publishedAt < :cutoff")
    int deletePublishedEventsBefore(@Param("cutoff") OffsetDateTime cutoff);

    /**
     * Get oldest pending event age (for monitoring)
     */
    @Query("SELECT MIN(e.createdAt) FROM OutboxEvent e WHERE e.status = 'PENDING'")
    OffsetDateTime findOldestPendingEventTime();
}

package com.ticketing.booking.entity;

import java.time.OffsetDateTime;
import java.util.UUID;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.PrePersist;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

/**
 * Outbox Event entity for Transactional Outbox Pattern.
 *
 * Events are written to this table in the same transaction as business data,
 * then a background processor publishes them to Kafka and marks them as published.
 *
 * This ensures exactly-once delivery semantics and eliminates dual-write problems.
 */
@Entity
@Table(name = "outbox_events", indexes = {
    @Index(name = "idx_outbox_pending", columnList = "status, created_at"),
    @Index(name = "idx_outbox_aggregate", columnList = "aggregate_type, aggregate_id")
})
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OutboxEvent {

    public enum OutboxStatus {
        PENDING,
        PUBLISHED,
        FAILED
    }

    public enum AggregateType {
        BOOKING,
        PAYMENT,
        TICKET
    }

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Enumerated(EnumType.STRING)
    @Column(name = "aggregate_type", nullable = false, length = 50)
    private AggregateType aggregateType;

    @Column(name = "aggregate_id", nullable = false, length = 100)
    private String aggregateId;

    @Column(name = "event_type", nullable = false, length = 50)
    private String eventType;

    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "payload", nullable = false, columnDefinition = "jsonb")
    private String payload;

    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 20)
    @Builder.Default
    private OutboxStatus status = OutboxStatus.PENDING;

    @Column(name = "created_at", nullable = false)
    private OffsetDateTime createdAt;

    @Column(name = "published_at")
    private OffsetDateTime publishedAt;

    @Column(name = "retry_count", nullable = false)
    @Builder.Default
    private int retryCount = 0;

    @Column(name = "last_error")
    private String lastError;

    @PrePersist
    public void onCreate() {
        if (this.createdAt == null) {
            this.createdAt = OffsetDateTime.now();
        }
        if (this.status == null) {
            this.status = OutboxStatus.PENDING;
        }
    }

    /**
     * Mark event as published
     */
    public void markPublished() {
        this.status = OutboxStatus.PUBLISHED;
        this.publishedAt = OffsetDateTime.now();
    }

    /**
     * Mark event as failed with error message
     */
    public void markFailed(String error) {
        this.status = OutboxStatus.FAILED;
        this.lastError = error;
        this.retryCount++;
    }

    /**
     * Reset to pending for retry
     */
    public void resetForRetry() {
        this.status = OutboxStatus.PENDING;
    }

    /**
     * Check if event should be retried (max 5 retries)
     */
    public boolean canRetry() {
        return this.retryCount < 5;
    }

    /**
     * Factory method for booking events
     */
    public static OutboxEvent forBooking(String bookingId, String eventType, String payload) {
        return OutboxEvent.builder()
                .aggregateType(AggregateType.BOOKING)
                .aggregateId(bookingId)
                .eventType(eventType)
                .payload(payload)
                .status(OutboxStatus.PENDING)
                .build();
    }
}

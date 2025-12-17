package com.ticketing.booking.entity;

import java.time.OffsetDateTime;
import java.util.UUID;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

/**
 * Dead Letter Queue entity for failed compensation operations.
 * Failed compensations are stored here for retry and manual intervention.
 */
@Entity
@Table(name = "failed_compensations", indexes = {
        @Index(name = "idx_fc_status", columnList = "status"),
        @Index(name = "idx_fc_type", columnList = "compensation_type"),
        @Index(name = "idx_fc_next_retry", columnList = "next_retry_at")
})
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class FailedCompensation {

    public enum CompensationType {
        CANCEL_PAYMENT,
        RELEASE_SEATS,
        REFUND_PAYMENT
    }

    public enum CompensationStatus {
        PENDING,      // Waiting for retry
        RETRYING,     // Currently being retried
        SUCCEEDED,    // Retry succeeded
        FAILED,       // All retries exhausted
        MANUAL        // Requires manual intervention
    }

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "compensation_id", nullable = false, unique = true, updatable = false)
    private UUID compensationId;

    @Column(name = "booking_id")
    private Long bookingId;

    @Column(name = "booking_reference", length = 50)
    private String bookingReference;

    @Enumerated(EnumType.STRING)
    @Column(name = "compensation_type", nullable = false, length = 30)
    private CompensationType compensationType;

    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 20)
    private CompensationStatus status;

    @Column(name = "reference_id", length = 100)
    private String referenceId; // paymentId or reservationId

    @Column(name = "error_message", length = 1000)
    private String errorMessage;

    @Column(name = "retry_count", nullable = false)
    @Builder.Default
    private int retryCount = 0;

    @Column(name = "max_retries", nullable = false)
    @Builder.Default
    private int maxRetries = 5;

    @Column(name = "next_retry_at")
    private OffsetDateTime nextRetryAt;

    @Column(name = "last_retry_at")
    private OffsetDateTime lastRetryAt;

    @Column(name = "payload", columnDefinition = "TEXT")
    private String payload; // JSON payload for retry

    @Column(name = "created_at", nullable = false)
    private OffsetDateTime createdAt;

    @Column(name = "updated_at", nullable = false)
    private OffsetDateTime updatedAt;

    @PrePersist
    public void onCreate() {
        OffsetDateTime now = OffsetDateTime.now();
        this.createdAt = now;
        this.updatedAt = now;
        this.compensationId = compensationId == null ? UUID.randomUUID() : compensationId;
        this.status = status == null ? CompensationStatus.PENDING : status;
        // Set first retry with exponential backoff (start at 30 seconds)
        if (this.nextRetryAt == null) {
            this.nextRetryAt = now.plusSeconds(30);
        }
    }

    @PreUpdate
    public void onUpdate() {
        this.updatedAt = OffsetDateTime.now();
    }

    /**
     * Calculate next retry time using exponential backoff
     * Delays: 30s, 1m, 2m, 4m, 8m (capped)
     */
    public void scheduleNextRetry() {
        this.retryCount++;
        this.lastRetryAt = OffsetDateTime.now();

        if (this.retryCount >= this.maxRetries) {
            this.status = CompensationStatus.FAILED;
            this.nextRetryAt = null;
        } else {
            // Exponential backoff: 30s * 2^retryCount, max 8 minutes
            long delaySeconds = Math.min(30L * (1L << this.retryCount), 480L);
            this.nextRetryAt = OffsetDateTime.now().plusSeconds(delaySeconds);
            this.status = CompensationStatus.PENDING;
        }
    }

    public void markSucceeded() {
        this.status = CompensationStatus.SUCCEEDED;
        this.nextRetryAt = null;
    }

    public void markForManualIntervention(String reason) {
        this.status = CompensationStatus.MANUAL;
        this.errorMessage = reason;
        this.nextRetryAt = null;
    }

    public boolean canRetry() {
        return this.retryCount < this.maxRetries &&
                (this.status == CompensationStatus.PENDING || this.status == CompensationStatus.RETRYING);
    }
}

package com.ticketing.payment.entity;

import com.ticketing.payment.entity.enums.IdempotencyStatus;
import jakarta.persistence.*;
import jakarta.validation.constraints.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;
import org.hibernate.annotations.Type;
import io.hypersistence.utils.hibernate.type.json.JsonBinaryType;

import java.time.LocalDateTime;
import java.util.UUID;

/**
 * Idempotency Key Entity
 * 
 * Tracks idempotency keys to prevent duplicate payment requests.
 * Maps to the 'idempotency_keys' table in the database.
 */
@Entity
@Table(name = "idempotency_keys", indexes = {
        @Index(name = "idx_idempotency_keys_key", columnList = "idempotency_key"),
        @Index(name = "idx_idempotency_keys_payment_id", columnList = "payment_id"),
        @Index(name = "idx_idempotency_keys_payment_uuid", columnList = "payment_uuid"),
        @Index(name = "idx_idempotency_keys_refund_id", columnList = "refund_id"),
        @Index(name = "idx_idempotency_keys_refund_uuid", columnList = "refund_uuid"),
        @Index(name = "idx_idempotency_keys_user_id", columnList = "user_id"),
        @Index(name = "idx_idempotency_keys_status", columnList = "status"),
        @Index(name = "idx_idempotency_keys_expires_at", columnList = "expires_at"),
        @Index(name = "idx_idempotency_keys_created_at", columnList = "created_at")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@ToString(exclude = { "requestBody", "requestHeaders", "responseBody", "responseHeaders", "metadata" })
@EqualsAndHashCode(of = "id")
public class IdempotencyKey {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Long id;

    @NotBlank(message = "Idempotency key is required")
    @Size(max = 255, message = "Idempotency key must not exceed 255 characters")
    @Column(name = "idempotency_key", unique = true, nullable = false, length = 255)
    private String idempotencyKey;

    // Request tracking
    @NotBlank(message = "Request path is required")
    @Size(max = 500, message = "Request path must not exceed 500 characters")
    @Column(name = "request_path", nullable = false, length = 500)
    private String requestPath;

    @NotBlank(message = "Request method is required")
    @Size(max = 10, message = "Request method must not exceed 10 characters")
    @Column(name = "request_method", nullable = false, length = 10)
    private String requestMethod;

    @Type(JsonBinaryType.class)
    @Column(name = "request_body", columnDefinition = "jsonb")
    private String requestBody;

    @Type(JsonBinaryType.class)
    @Column(name = "request_headers", columnDefinition = "jsonb")
    private String requestHeaders;

    // Response tracking
    @Column(name = "response_status")
    private Integer responseStatus;

    @Type(JsonBinaryType.class)
    @Column(name = "response_body", columnDefinition = "jsonb")
    private String responseBody;

    @Type(JsonBinaryType.class)
    @Column(name = "response_headers", columnDefinition = "jsonb")
    private String responseHeaders;

    // Associated entities
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "payment_id", foreignKey = @ForeignKey(name = "fk_idempotency_payment"))
    private Payment payment;

    @Column(name = "payment_uuid")
    private UUID paymentUuid;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "refund_id", foreignKey = @ForeignKey(name = "fk_idempotency_refund"))
    private Refund refund;

    @Column(name = "refund_uuid")
    private UUID refundUuid;

    // User context
    @Size(max = 255, message = "User ID must not exceed 255 characters")
    @Column(name = "user_id", length = 255)
    private String userId;

    @Size(max = 45, message = "IP address must not exceed 45 characters")
    @Column(name = "ip_address", length = 45)
    private String ipAddress;

    @Column(name = "user_agent", columnDefinition = "TEXT")
    private String userAgent;

    // Status tracking
    @NotNull(message = "Status is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 50)
    private IdempotencyStatus status = IdempotencyStatus.PROCESSING;

    // Expiration
    @NotNull(message = "Expiration time is required")
    @Future(message = "Expiration time must be in the future")
    @Column(name = "expires_at", nullable = false)
    private LocalDateTime expiresAt;

    // Metadata
    @Type(JsonBinaryType.class)
    @Column(name = "metadata", columnDefinition = "jsonb")
    private String metadata;

    // Audit fields
    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(name = "updated_at", nullable = false)
    private LocalDateTime updatedAt;

    @Column(name = "completed_at")
    private LocalDateTime completedAt;

    /**
     * Pre-persist callback
     */
    @PrePersist
    protected void onCreate() {
        if (status == null) {
            status = IdempotencyStatus.PROCESSING;
        }
        // Copy payment UUID if not set
        if (paymentUuid == null && payment != null) {
            paymentUuid = payment.getPaymentId();
        }
        // Copy refund UUID if not set
        if (refundUuid == null && refund != null) {
            refundUuid = refund.getRefundId();
        }
    }

    /**
     * Pre-update callback
     */
    @PreUpdate
    protected void onUpdate() {
        // Set completed_at when status changes to completed
        if (status == IdempotencyStatus.COMPLETED && completedAt == null) {
            completedAt = LocalDateTime.now();
        }
    }

    // Business logic methods

    /**
     * Check if key is processing
     */
    public boolean isProcessing() {
        return status == IdempotencyStatus.PROCESSING;
    }

    /**
     * Check if key is completed
     */
    public boolean isCompleted() {
        return status == IdempotencyStatus.COMPLETED;
    }

    /**
     * Check if key is failed
     */
    public boolean isFailed() {
        return status == IdempotencyStatus.FAILED;
    }

    /**
     * Check if key is expired
     */
    public boolean isExpired() {
        return LocalDateTime.now().isAfter(expiresAt);
    }

    /**
     * Check if key is still valid (not expired)
     */
    public boolean isValid() {
        return !isExpired();
    }

    /**
     * Check if key can be reused (completed and not expired)
     */
    public boolean canBeReused() {
        return isCompleted() && isValid();
    }

    /**
     * Mark as completed
     */
    public void markAsCompleted(Integer responseStatus, String responseBody) {
        this.status = IdempotencyStatus.COMPLETED;
        this.responseStatus = responseStatus;
        this.responseBody = responseBody;
        this.completedAt = LocalDateTime.now();
    }

    /**
     * Mark as failed
     */
    public void markAsFailed(Integer responseStatus, String responseBody) {
        this.status = IdempotencyStatus.FAILED;
        this.responseStatus = responseStatus;
        this.responseBody = responseBody;
        this.completedAt = LocalDateTime.now();
    }

    /**
     * Get remaining time until expiration in minutes
     */
    public long getRemainingMinutes() {
        if (isExpired()) {
            return 0;
        }
        return java.time.Duration.between(LocalDateTime.now(), expiresAt).toMinutes();
    }

    /**
     * Get display status
     */
    public String getDisplayStatus() {
        return status.getDisplayName();
    }

    /**
     * Check if key is for payment request
     */
    public boolean isPaymentRequest() {
        return payment != null || paymentUuid != null;
    }

    /**
     * Check if key is for refund request
     */
    public boolean isRefundRequest() {
        return refund != null || refundUuid != null;
    }

    /**
     * Create idempotency key with default expiration (24 hours)
     */
    public static IdempotencyKey create(String key, String path, String method, String userId, String ipAddress) {
        return IdempotencyKey.builder()
                .idempotencyKey(key)
                .requestPath(path)
                .requestMethod(method)
                .userId(userId)
                .ipAddress(ipAddress)
                .status(IdempotencyStatus.PROCESSING)
                .expiresAt(LocalDateTime.now().plusHours(24))
                .build();
    }

    /**
     * Create idempotency key with custom expiration
     */
    public static IdempotencyKey create(String key, String path, String method, String userId, String ipAddress,
            int expirationHours) {
        return IdempotencyKey.builder()
                .idempotencyKey(key)
                .requestPath(path)
                .requestMethod(method)
                .userId(userId)
                .ipAddress(ipAddress)
                .status(IdempotencyStatus.PROCESSING)
                .expiresAt(LocalDateTime.now().plusHours(expirationHours))
                .build();
    }
}



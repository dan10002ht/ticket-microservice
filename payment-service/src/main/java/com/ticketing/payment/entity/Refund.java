package com.ticketing.payment.entity;

import com.ticketing.payment.entity.enums.RefundStatus;
import com.ticketing.payment.entity.enums.RefundType;
import com.ticketing.payment.entity.enums.PaymentGateway;
import jakarta.persistence.*;
import jakarta.validation.constraints.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;
import org.hibernate.annotations.Type;
import io.hypersistence.utils.hibernate.type.json.JsonBinaryType;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

/**
 * Refund Entity
 * 
 * Represents a refund transaction in the system.
 * Maps to the 'refunds' table in the database.
 */
@Entity
@Table(name = "refunds", indexes = {
        @Index(name = "idx_refunds_payment_id", columnList = "payment_id"),
        @Index(name = "idx_refunds_payment_uuid", columnList = "payment_uuid"),
        @Index(name = "idx_refunds_status", columnList = "status"),
        @Index(name = "idx_refunds_refund_id", columnList = "refund_id"),
        @Index(name = "idx_refunds_external_reference", columnList = "external_reference"),
        @Index(name = "idx_refunds_gateway_provider", columnList = "gateway_provider"),
        @Index(name = "idx_refunds_created_at", columnList = "created_at"),
        @Index(name = "idx_refunds_payment_status", columnList = "payment_id, status")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@ToString(exclude = { "gatewayResponse", "metadata", "payment" })
@EqualsAndHashCode(of = "id")
public class Refund {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Long id;

    @Column(name = "refund_id", unique = true, nullable = false, updatable = false)
    private UUID refundId;

    // Payment reference
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "payment_id", nullable = false, foreignKey = @ForeignKey(name = "fk_refund_payment"))
    @NotNull(message = "Payment is required")
    private Payment payment;

    @NotNull(message = "Payment UUID is required")
    @Column(name = "payment_uuid", nullable = false)
    private UUID paymentUuid;

    // Refund details
    @NotNull(message = "Amount is required")
    @DecimalMin(value = "0.01", message = "Refund amount must be greater than 0")
    @Digits(integer = 10, fraction = 2, message = "Amount must have at most 10 integer digits and 2 decimal places")
    @Column(name = "amount", nullable = false, precision = 10, scale = 2)
    private BigDecimal amount;

    @NotBlank(message = "Currency is required")
    @Size(min = 3, max = 3, message = "Currency must be 3 characters (ISO 4217)")
    @Column(name = "currency", nullable = false, length = 3)
    private String currency = "USD";

    @NotNull(message = "Refund type is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "refund_type", nullable = false, length = 50)
    private RefundType refundType;

    // Status tracking
    @NotNull(message = "Status is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 50)
    private RefundStatus status = RefundStatus.PENDING;

    // External references
    @Size(max = 255, message = "External reference must not exceed 255 characters")
    @Column(name = "external_reference", length = 255)
    private String externalReference;

    @Enumerated(EnumType.STRING)
    @Column(name = "gateway_provider", length = 50)
    private PaymentGateway gatewayProvider;

    // Gateway response (stored as JSONB)
    @Type(JsonBinaryType.class)
    @Column(name = "gateway_response", columnDefinition = "jsonb")
    private String gatewayResponse;

    @Column(name = "failure_reason", columnDefinition = "TEXT")
    private String failureReason;

    // Refund reason
    @NotBlank(message = "Reason is required")
    @Size(max = 255, message = "Reason must not exceed 255 characters")
    @Column(name = "reason", nullable = false, length = 255)
    private String reason;

    @Column(name = "description", columnDefinition = "TEXT")
    private String description;

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

    @Size(max = 255, message = "Created by must not exceed 255 characters")
    @Column(name = "created_by", length = 255)
    private String createdBy;

    @Size(max = 255, message = "Updated by must not exceed 255 characters")
    @Column(name = "updated_by", length = 255)
    private String updatedBy;

    // Refund completion
    @Column(name = "refunded_at")
    private LocalDateTime refundedAt;

    @Column(name = "cancelled_at")
    private LocalDateTime cancelledAt;

    @Size(max = 255, message = "Cancelled by must not exceed 255 characters")
    @Column(name = "cancelled_by", length = 255)
    private String cancelledBy;

    @Column(name = "cancellation_reason", columnDefinition = "TEXT")
    private String cancellationReason;

    /**
     * Pre-persist callback to set refund ID
     */
    @PrePersist
    protected void onCreate() {
        if (refundId == null) {
            refundId = UUID.randomUUID();
        }
        if (status == null) {
            status = RefundStatus.PENDING;
        }
        if (currency == null || currency.isEmpty()) {
            currency = "USD";
        }
        // Copy payment UUID if not set
        if (paymentUuid == null && payment != null) {
            paymentUuid = payment.getPaymentId();
        }
        // Copy gateway provider from payment if not set
        if (gatewayProvider == null && payment != null) {
            gatewayProvider = payment.getGatewayProvider();
        }
    }

    /**
     * Pre-update callback
     */
    @PreUpdate
    protected void onUpdate() {
        // Additional business logic can be added here
    }

    // Business logic methods

    /**
     * Check if refund is pending
     */
    public boolean isPending() {
        return status == RefundStatus.PENDING;
    }

    /**
     * Check if refund is processing
     */
    public boolean isProcessing() {
        return status == RefundStatus.PROCESSING;
    }

    /**
     * Check if refund is successful
     */
    public boolean isSuccess() {
        return status == RefundStatus.SUCCESS;
    }

    /**
     * Check if refund has failed
     */
    public boolean isFailed() {
        return status == RefundStatus.FAILED;
    }

    /**
     * Check if refund is cancelled
     */
    public boolean isCancelled() {
        return status == RefundStatus.CANCELLED;
    }

    /**
     * Check if refund is full refund
     */
    public boolean isFullRefund() {
        return refundType == RefundType.FULL;
    }

    /**
     * Check if refund is partial refund
     */
    public boolean isPartialRefund() {
        return refundType == RefundType.PARTIAL;
    }

    /**
     * Check if refund can be cancelled
     */
    public boolean canBeCancelled() {
        return isPending() || isProcessing();
    }

    /**
     * Mark refund as processing
     */
    public void markAsProcessing() {
        this.status = RefundStatus.PROCESSING;
    }

    /**
     * Mark refund as successful
     */
    public void markAsSuccess(LocalDateTime refundedAt) {
        this.status = RefundStatus.SUCCESS;
        this.refundedAt = refundedAt;
    }

    /**
     * Mark refund as failed
     */
    public void markAsFailed(String failureReason) {
        this.status = RefundStatus.FAILED;
        this.failureReason = failureReason;
    }

    /**
     * Mark refund as cancelled
     */
    public void markAsCancelled(String cancelledBy, String reason) {
        this.status = RefundStatus.CANCELLED;
        this.cancelledAt = LocalDateTime.now();
        this.cancelledBy = cancelledBy;
        this.cancellationReason = reason;
    }

    /**
     * Check if refund is in terminal state (cannot be changed)
     */
    public boolean isTerminalState() {
        return isSuccess() || isFailed() || isCancelled();
    }

    /**
     * Get display status for UI
     */
    public String getDisplayStatus() {
        return status.getDisplayName();
    }

    /**
     * Get refund percentage of payment amount
     */
    public BigDecimal getRefundPercentage() {
        if (payment == null || payment.getAmount().compareTo(BigDecimal.ZERO) == 0) {
            return BigDecimal.ZERO;
        }
        return amount.divide(payment.getAmount(), 4, java.math.RoundingMode.HALF_UP)
                .multiply(BigDecimal.valueOf(100));
    }

    /**
     * Validate refund amount against payment amount
     */
    public boolean isValidRefundAmount() {
        if (payment == null) {
            return false;
        }
        // Refund amount should not exceed payment amount
        return amount.compareTo(payment.getAmount()) <= 0 && amount.compareTo(BigDecimal.ZERO) > 0;
    }

    /**
     * Get remaining refundable amount for the payment
     */
    public BigDecimal getRemainingRefundableAmount() {
        if (payment == null) {
            return BigDecimal.ZERO;
        }
        // This should be calculated from all successful refunds for the payment
        // For now, return payment amount minus this refund amount
        return payment.getAmount().subtract(amount);
    }

    /**
     * Check if this refund will fully refund the payment
     */
    public boolean willFullyRefundPayment() {
        if (payment == null) {
            return false;
        }
        return amount.compareTo(payment.getAmount()) == 0;
    }
}

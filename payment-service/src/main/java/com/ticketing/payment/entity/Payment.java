package com.ticketing.payment.entity;

import com.ticketing.payment.entity.enums.PaymentMethod;
import com.ticketing.payment.entity.enums.PaymentStatus;
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
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

/**
 * Payment Entity
 * 
 * Represents a payment transaction in the system.
 * Maps to the 'payments' table in the database.
 */
@Entity
@Table(name = "payments", indexes = {
        @Index(name = "idx_payments_booking_id", columnList = "booking_id"),
        @Index(name = "idx_payments_user_id", columnList = "user_id"),
        @Index(name = "idx_payments_status", columnList = "status"),
        @Index(name = "idx_payments_payment_id", columnList = "payment_id"),
        @Index(name = "idx_payments_external_reference", columnList = "external_reference"),
        @Index(name = "idx_payments_gateway_provider", columnList = "gateway_provider"),
        @Index(name = "idx_payments_idempotency_key", columnList = "idempotency_key"),
        @Index(name = "idx_payments_created_at", columnList = "created_at"),
        @Index(name = "idx_payments_user_status", columnList = "user_id, status"),
        @Index(name = "idx_payments_booking_status", columnList = "booking_id, status")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@ToString(exclude = { "gatewayResponse", "metadata" })
@EqualsAndHashCode(of = "id")
public class Payment {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Long id;

    @Column(name = "payment_id", unique = true, nullable = false, updatable = false)
    private UUID paymentId;

    // Booking reference
    @NotBlank(message = "Booking ID is required")
    @Size(max = 255, message = "Booking ID must not exceed 255 characters")
    @Column(name = "booking_id", nullable = false, length = 255)
    private String bookingId;

    @Size(max = 255, message = "Ticket ID must not exceed 255 characters")
    @Column(name = "ticket_id", length = 255)
    private String ticketId;

    @NotBlank(message = "User ID is required")
    @Size(max = 255, message = "User ID must not exceed 255 characters")
    @Column(name = "user_id", nullable = false, length = 255)
    private String userId;

    // Payment details
    @NotNull(message = "Amount is required")
    @DecimalMin(value = "0.0", inclusive = false, message = "Amount must be greater than 0")
    @Digits(integer = 10, fraction = 2, message = "Amount must have at most 10 integer digits and 2 decimal places")
    @Column(name = "amount", nullable = false, precision = 10, scale = 2)
    private BigDecimal amount;

    @NotBlank(message = "Currency is required")
    @Size(min = 3, max = 3, message = "Currency must be 3 characters (ISO 4217)")
    @Column(name = "currency", nullable = false, length = 3)
    private String currency = "USD";

    @NotNull(message = "Payment method is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "payment_method", nullable = false, length = 50)
    private PaymentMethod paymentMethod;

    // Status tracking
    @NotNull(message = "Status is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 50)
    private PaymentStatus status = PaymentStatus.PENDING;

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

    // Metadata
    @Column(name = "description", columnDefinition = "TEXT")
    private String description;

    @Type(JsonBinaryType.class)
    @Column(name = "metadata", columnDefinition = "jsonb")
    private String metadata;

    // Idempotency
    @Size(max = 255, message = "Idempotency key must not exceed 255 characters")
    @Column(name = "idempotency_key", unique = true, length = 255)
    private String idempotencyKey;

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

    // Payment completion
    @Column(name = "paid_at")
    private LocalDateTime paidAt;

    @Column(name = "cancelled_at")
    private LocalDateTime cancelledAt;

    @Size(max = 255, message = "Cancelled by must not exceed 255 characters")
    @Column(name = "cancelled_by", length = 255)
    private String cancelledBy;

    @Column(name = "cancellation_reason", columnDefinition = "TEXT")
    private String cancellationReason;

    // Relationships
    @OneToMany(mappedBy = "payment", cascade = CascadeType.ALL, orphanRemoval = true, fetch = FetchType.LAZY)
    @Builder.Default
    private List<Refund> refunds = new ArrayList<>();

    /**
     * Pre-persist callback to set payment ID
     */
    @PrePersist
    protected void onCreate() {
        if (paymentId == null) {
            paymentId = UUID.randomUUID();
        }
        if (status == null) {
            status = PaymentStatus.PENDING;
        }
        if (currency == null || currency.isEmpty()) {
            currency = "USD";
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
     * Check if payment is pending
     */
    public boolean isPending() {
        return status == PaymentStatus.PENDING;
    }

    /**
     * Check if payment is processing
     */
    public boolean isProcessing() {
        return status == PaymentStatus.PROCESSING;
    }

    /**
     * Check if payment is successful
     */
    public boolean isSuccess() {
        return status == PaymentStatus.SUCCESS;
    }

    /**
     * Check if payment has failed
     */
    public boolean isFailed() {
        return status == PaymentStatus.FAILED;
    }

    /**
     * Check if payment is cancelled
     */
    public boolean isCancelled() {
        return status == PaymentStatus.CANCELLED;
    }

    /**
     * Check if payment is refunded
     */
    public boolean isRefunded() {
        return status == PaymentStatus.REFUNDED || status == PaymentStatus.PARTIALLY_REFUNDED;
    }

    /**
     * Check if payment is fully refunded
     */
    public boolean isFullyRefunded() {
        return status == PaymentStatus.REFUNDED;
    }

    /**
     * Check if payment is partially refunded
     */
    public boolean isPartiallyRefunded() {
        return status == PaymentStatus.PARTIALLY_REFUNDED;
    }

    /**
     * Check if payment can be refunded
     */
    public boolean canBeRefunded() {
        return isSuccess() || isPartiallyRefunded();
    }

    /**
     * Check if payment can be cancelled
     */
    public boolean canBeCancelled() {
        return isPending() || isProcessing();
    }

    /**
     * Mark payment as processing
     */
    public void markAsProcessing() {
        this.status = PaymentStatus.PROCESSING;
    }

    /**
     * Mark payment as successful
     */
    public void markAsSuccess(LocalDateTime paidAt) {
        this.status = PaymentStatus.SUCCESS;
        this.paidAt = paidAt;
    }

    /**
     * Mark payment as failed
     */
    public void markAsFailed(String failureReason) {
        this.status = PaymentStatus.FAILED;
        this.failureReason = failureReason;
    }

    /**
     * Mark payment as cancelled
     */
    public void markAsCancelled(String cancelledBy, String reason) {
        this.status = PaymentStatus.CANCELLED;
        this.cancelledAt = LocalDateTime.now();
        this.cancelledBy = cancelledBy;
        this.cancellationReason = reason;
    }

    /**
     * Mark payment as refunded
     */
    public void markAsRefunded() {
        this.status = PaymentStatus.REFUNDED;
    }

    /**
     * Mark payment as partially refunded
     */
    public void markAsPartiallyRefunded() {
        this.status = PaymentStatus.PARTIALLY_REFUNDED;
    }

    /**
     * Check if payment is in terminal state (cannot be changed)
     */
    public boolean isTerminalState() {
        return isSuccess() || isFailed() || isCancelled() || isFullyRefunded();
    }

    /**
     * Get display status for UI
     */
    public String getDisplayStatus() {
        return status.getDisplayName();
    }

    // Refund-related methods

    /**
     * Add refund to payment
     */
    public void addRefund(Refund refund) {
        refunds.add(refund);
        refund.setPayment(this);
        refund.setPaymentUuid(this.paymentId);
    }

    /**
     * Remove refund from payment
     */
    public void removeRefund(Refund refund) {
        refunds.remove(refund);
        refund.setPayment(null);
    }

    /**
     * Get total refunded amount
     */
    public BigDecimal getTotalRefundedAmount() {
        return refunds.stream()
                .filter(Refund::isSuccess)
                .map(Refund::getAmount)
                .reduce(BigDecimal.ZERO, BigDecimal::add);
    }

    /**
     * Get remaining refundable amount
     */
    public BigDecimal getRemainingRefundableAmount() {
        return amount.subtract(getTotalRefundedAmount());
    }

    /**
     * Check if payment has been refunded
     */
    public boolean hasRefunds() {
        return !refunds.isEmpty() && refunds.stream().anyMatch(Refund::isSuccess);
    }

    /**
     * Get successful refunds count
     */
    public long getSuccessfulRefundsCount() {
        return refunds.stream().filter(Refund::isSuccess).count();
    }
}

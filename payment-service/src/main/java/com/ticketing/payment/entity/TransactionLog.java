package com.ticketing.payment.entity;

import com.ticketing.payment.entity.enums.TransactionType;
import com.ticketing.payment.entity.enums.TransactionStatus;
import com.ticketing.payment.entity.enums.PaymentGateway;
import jakarta.persistence.*;
import jakarta.validation.constraints.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.Type;
import io.hypersistence.utils.hibernate.type.json.JsonBinaryType;

import java.time.LocalDateTime;
import java.util.UUID;

/**
 * Transaction Log Entity
 * 
 * Audit trail for all payment and refund transactions.
 * Maps to the 'transaction_logs' table in the database.
 */
@Entity
@Table(name = "transaction_logs", indexes = {
        @Index(name = "idx_transaction_logs_payment_id", columnList = "payment_id"),
        @Index(name = "idx_transaction_logs_payment_uuid", columnList = "payment_uuid"),
        @Index(name = "idx_transaction_logs_refund_id", columnList = "refund_id"),
        @Index(name = "idx_transaction_logs_refund_uuid", columnList = "refund_uuid"),
        @Index(name = "idx_transaction_logs_transaction_type", columnList = "transaction_type"),
        @Index(name = "idx_transaction_logs_event_name", columnList = "event_name"),
        @Index(name = "idx_transaction_logs_gateway_provider", columnList = "gateway_provider"),
        @Index(name = "idx_transaction_logs_external_reference", columnList = "external_reference"),
        @Index(name = "idx_transaction_logs_status", columnList = "status"),
        @Index(name = "idx_transaction_logs_user_id", columnList = "user_id"),
        @Index(name = "idx_transaction_logs_correlation_id", columnList = "correlation_id"),
        @Index(name = "idx_transaction_logs_created_at", columnList = "created_at")
})
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@ToString(exclude = { "requestData", "responseData", "headers", "metadata" })
@EqualsAndHashCode(of = "id")
public class TransactionLog {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Long id;

    @Column(name = "log_id", unique = true, nullable = false, updatable = false)
    private UUID logId;

    // Transaction references
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "payment_id", foreignKey = @ForeignKey(name = "fk_transaction_log_payment"))
    private Payment payment;

    @Column(name = "payment_uuid")
    private UUID paymentUuid;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "refund_id", foreignKey = @ForeignKey(name = "fk_transaction_log_refund"))
    private Refund refund;

    @Column(name = "refund_uuid")
    private UUID refundUuid;

    // Transaction details
    @NotNull(message = "Transaction type is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "transaction_type", nullable = false, length = 50)
    private TransactionType transactionType;

    @NotBlank(message = "Event name is required")
    @Size(max = 100, message = "Event name must not exceed 100 characters")
    @Column(name = "event_name", nullable = false, length = 100)
    private String eventName;

    // Gateway information
    @Enumerated(EnumType.STRING)
    @Column(name = "gateway_provider", length = 50)
    private PaymentGateway gatewayProvider;

    @Size(max = 255, message = "External reference must not exceed 255 characters")
    @Column(name = "external_reference", length = 255)
    private String externalReference;

    // Request/Response data (stored as JSONB)
    @Type(JsonBinaryType.class)
    @Column(name = "request_data", columnDefinition = "jsonb")
    private String requestData;

    @Type(JsonBinaryType.class)
    @Column(name = "response_data", columnDefinition = "jsonb")
    private String responseData;

    @Type(JsonBinaryType.class)
    @Column(name = "headers", columnDefinition = "jsonb")
    private String headers;

    // Status and error tracking
    @NotNull(message = "Status is required")
    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 50)
    private TransactionStatus status;

    @Size(max = 100, message = "Error code must not exceed 100 characters")
    @Column(name = "error_code", length = 100)
    private String errorCode;

    @Column(name = "error_message", columnDefinition = "TEXT")
    private String errorMessage;

    // Performance tracking
    @Min(value = 0, message = "Duration must be non-negative")
    @Column(name = "duration_ms")
    private Integer durationMs;

    // Context information
    @Size(max = 255, message = "User ID must not exceed 255 characters")
    @Column(name = "user_id", length = 255)
    private String userId;

    @Size(max = 45, message = "IP address must not exceed 45 characters")
    @Column(name = "ip_address", length = 45)
    private String ipAddress;

    @Column(name = "user_agent", columnDefinition = "TEXT")
    private String userAgent;

    @Column(name = "correlation_id")
    private UUID correlationId;

    // Metadata
    @Type(JsonBinaryType.class)
    @Column(name = "metadata", columnDefinition = "jsonb")
    private String metadata;

    // Audit fields
    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @Size(max = 255, message = "Created by must not exceed 255 characters")
    @Column(name = "created_by", length = 255)
    private String createdBy;

    /**
     * Pre-persist callback to set log ID
     */
    @PrePersist
    protected void onCreate() {
        if (logId == null) {
            logId = UUID.randomUUID();
        }
        if (correlationId == null) {
            correlationId = UUID.randomUUID();
        }
        // Copy payment UUID if not set
        if (paymentUuid == null && payment != null) {
            paymentUuid = payment.getPaymentId();
        }
        // Copy refund UUID if not set
        if (refundUuid == null && refund != null) {
            refundUuid = refund.getRefundId();
        }
        // Copy gateway provider if not set
        if (gatewayProvider == null) {
            if (payment != null) {
                gatewayProvider = payment.getGatewayProvider();
            } else if (refund != null) {
                gatewayProvider = refund.getGatewayProvider();
            }
        }
    }

    // Business logic methods

    /**
     * Check if transaction was successful
     */
    public boolean isSuccess() {
        return status == TransactionStatus.SUCCESS;
    }

    /**
     * Check if transaction failed
     */
    public boolean isFailed() {
        return status == TransactionStatus.FAILED;
    }

    /**
     * Check if transaction is pending
     */
    public boolean isPending() {
        return status == TransactionStatus.PENDING;
    }

    /**
     * Check if log is for payment transaction
     */
    public boolean isPaymentTransaction() {
        return transactionType.isPaymentRelated();
    }

    /**
     * Check if log is for refund transaction
     */
    public boolean isRefundTransaction() {
        return transactionType.isRefundRelated();
    }

    /**
     * Check if log is for webhook event
     */
    public boolean isWebhookTransaction() {
        return transactionType == TransactionType.WEBHOOK_RECEIVED;
    }

    /**
     * Check if transaction has error
     */
    public boolean hasError() {
        return errorCode != null || errorMessage != null;
    }

    /**
     * Get duration in seconds
     */
    public Double getDurationInSeconds() {
        return durationMs != null ? durationMs / 1000.0 : null;
    }

    /**
     * Check if transaction was slow (> 5 seconds)
     */
    public boolean isSlow() {
        return durationMs != null && durationMs > 5000;
    }

    /**
     * Check if transaction was fast (< 1 second)
     */
    public boolean isFast() {
        return durationMs != null && durationMs < 1000;
    }

    /**
     * Get display status
     */
    public String getDisplayStatus() {
        return status.getDisplayName();
    }

    /**
     * Get display transaction type
     */
    public String getDisplayTransactionType() {
        return transactionType.getDisplayName();
    }

    /**
     * Create log for payment initiated
     */
    public static TransactionLog createPaymentInitiated(Payment payment, String userId, String ipAddress) {
        return TransactionLog.builder()
                .payment(payment)
                .paymentUuid(payment.getPaymentId())
                .transactionType(TransactionType.PAYMENT_INITIATED)
                .eventName("payment.initiated")
                .status(TransactionStatus.SUCCESS)
                .userId(userId)
                .ipAddress(ipAddress)
                .gatewayProvider(payment.getGatewayProvider())
                .build();
    }

    /**
     * Create log for payment success
     */
    public static TransactionLog createPaymentSuccess(Payment payment, String externalReference) {
        return TransactionLog.builder()
                .payment(payment)
                .paymentUuid(payment.getPaymentId())
                .transactionType(TransactionType.PAYMENT_SUCCESS)
                .eventName("payment.success")
                .status(TransactionStatus.SUCCESS)
                .externalReference(externalReference)
                .gatewayProvider(payment.getGatewayProvider())
                .build();
    }

    /**
     * Create log for payment failed
     */
    public static TransactionLog createPaymentFailed(Payment payment, String errorCode, String errorMessage) {
        return TransactionLog.builder()
                .payment(payment)
                .paymentUuid(payment.getPaymentId())
                .transactionType(TransactionType.PAYMENT_FAILED)
                .eventName("payment.failed")
                .status(TransactionStatus.FAILED)
                .errorCode(errorCode)
                .errorMessage(errorMessage)
                .gatewayProvider(payment.getGatewayProvider())
                .build();
    }

    /**
     * Create log for refund initiated
     */
    public static TransactionLog createRefundInitiated(Refund refund, String userId) {
        return TransactionLog.builder()
                .refund(refund)
                .refundUuid(refund.getRefundId())
                .payment(refund.getPayment())
                .paymentUuid(refund.getPaymentUuid())
                .transactionType(TransactionType.REFUND_INITIATED)
                .eventName("refund.initiated")
                .status(TransactionStatus.SUCCESS)
                .userId(userId)
                .gatewayProvider(refund.getGatewayProvider())
                .build();
    }
}



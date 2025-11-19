package com.ticketing.payment.repository;

import com.ticketing.payment.entity.TransactionLog;
import com.ticketing.payment.entity.enums.TransactionStatus;
import com.ticketing.payment.entity.enums.TransactionType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Repository interface for {@link TransactionLog}.
 *
 * Provides access patterns for auditing, monitoring, and reconciliation jobs.
 */
@Repository
public interface TransactionLogRepository extends JpaRepository<TransactionLog, Long> {

    /**
     * Lookup by public UUID.
     */
    Optional<TransactionLog> findByLogId(UUID logId);

    /**
     * Retrieve all logs for a payment, newest first.
     */
    List<TransactionLog> findByPayment_PaymentIdOrderByCreatedAtDesc(UUID paymentId);

    /**
     * Retrieve all logs for a refund, newest first.
     */
    List<TransactionLog> findByRefund_RefundIdOrderByCreatedAtDesc(UUID refundId);

    /**
     * Paginate logs filtered by transaction type.
     */
    Page<TransactionLog> findByTransactionTypeOrderByCreatedAtDesc(TransactionType type, Pageable pageable);

    /**
     * Paginate logs filtered by status.
     */
    Page<TransactionLog> findByStatusOrderByCreatedAtDesc(TransactionStatus status, Pageable pageable);

    /**
     * Find logs that contain specific correlation id (distributed tracing).
     */
    List<TransactionLog> findByCorrelationIdOrderByCreatedAtAsc(UUID correlationId);

    /**
     * Retrieve log entry by gateway reference (useful for webhook reconciliation).
     */
    Optional<TransactionLog> findByExternalReference(String externalReference);

    /**
     * Count logs for a given payment and type.
     */
    long countByPayment_PaymentIdAndTransactionType(UUID paymentId, TransactionType type);

    /**
     * Fetch slow transactions beyond threshold for monitoring alerts.
     */
    List<TransactionLog> findByDurationMsGreaterThan(Integer durationMs);

    /**
     * Delete old logs (retention policy).
     */
    long deleteByCreatedAtBefore(LocalDateTime cutoff);

    /**
     * Aggregate number of failures per gateway within period (for dashboards).
     */
    @Query("""
            SELECT tl.gatewayProvider, COUNT(tl)
            FROM TransactionLog tl
            WHERE tl.status = 'FAILED'
              AND tl.createdAt BETWEEN :start AND :end
            GROUP BY tl.gatewayProvider
            """)
    List<Object[]> countFailuresByGatewayBetween(@Param("start") LocalDateTime start,
            @Param("end") LocalDateTime end);
}

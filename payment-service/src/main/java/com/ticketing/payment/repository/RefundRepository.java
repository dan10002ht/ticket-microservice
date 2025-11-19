package com.ticketing.payment.repository;

import com.ticketing.payment.entity.Refund;
import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.entity.enums.RefundStatus;
import com.ticketing.payment.entity.enums.RefundType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Repository interface for {@link Refund} aggregate.
 *
 * Supplies high-level query methods used by refund workflow, reconciliation,
 * and monitoring jobs.
 */
@Repository
public interface RefundRepository extends JpaRepository<Refund, Long> {

    /**
     * Lookup by public UUID.
     */
    Optional<Refund> findByRefundId(UUID refundId);

    /**
     * Retrieve refund belonging to a payment by amount and type (ensures
     * idempotency for partial vs full refunds).
     */
    Optional<Refund> findByPayment_PaymentIdAndAmountAndRefundType(UUID paymentId, BigDecimal amount,
            RefundType refundType);

    /**
     * Fetch all refunds for a payment ordered by creation time desc.
     */
    List<Refund> findByPayment_PaymentIdOrderByCreatedAtDesc(UUID paymentId);

    /**
     * Paginate refunds filtered by status.
     */
    Page<Refund> findByStatusOrderByCreatedAtDesc(RefundStatus status, Pageable pageable);

    /**
     * Retrieve all pending/processing refunds that were last updated before a
     * threshold (for retry job).
     */
    List<Refund> findByStatusInAndUpdatedAtBefore(List<RefundStatus> statuses, LocalDateTime updatedBefore);

    /**
     * Count successful refunds for a payment (used to determine partial/fully
     * refunded state).
     */
    long countByPayment_PaymentIdAndStatus(UUID paymentId, RefundStatus status);

    /**
     * Find refunds by gateway external reference.
     */
    Optional<Refund> findByGatewayProviderAndExternalReference(PaymentGateway gateway, String externalReference);

    /**
     * Sum total successful refunds for a payment (used for reconciliation).
     */
    @Query("SELECT COALESCE(SUM(r.amount), 0) FROM Refund r WHERE r.payment.paymentId = :paymentId AND r.status = 'SUCCESS'")
    BigDecimal sumSuccessfulRefundsByPayment(@Param("paymentId") UUID paymentId);

    Optional<Refund> findByProviderReference(String providerReference);
}

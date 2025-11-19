package com.ticketing.payment.repository;

import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.enums.PaymentStatus;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Repository interface for {@link Payment} aggregate.
 *
 * Centralized place for frequently used read/write operations so the service
 * layer does not need to craft queries directly.
 */
@Repository
public interface PaymentRepository extends JpaRepository<Payment, Long> {

    /**
     * Lookup by public UUID (indexed column).
     */
    Optional<Payment> findByPaymentId(UUID paymentId);

    /**
     * Lookup by unique idempotency key.
     */
    Optional<Payment> findByIdempotencyKey(String idempotencyKey);

    /**
     * Existence check for idempotency enforcement.
     */
    boolean existsByIdempotencyKey(String idempotencyKey);

    /**
     * Retrieve latest payment for a booking (covers booking queue workflow).
     */
    Optional<Payment> findTopByBookingIdOrderByCreatedAtDesc(String bookingId);

    /**
     * List payments by booking with newest first.
     */
    List<Payment> findByBookingIdOrderByCreatedAtDesc(String bookingId);

    /**
     * Fetch paginated history for a user filtered by status if provided.
     */
    Page<Payment> findByUserIdAndStatusOrderByCreatedAtDesc(String userId, PaymentStatus status, Pageable pageable);

    /**
     * Retrieve payments older than provided timestamp for cleanup tasks.
     */
    List<Payment> findByStatusAndCreatedAtBefore(PaymentStatus status, LocalDateTime createdBefore);

    /**
     * Retrieve payments stuck in processing beyond a threshold.
     */
    List<Payment> findByStatusAndUpdatedAtBefore(PaymentStatus status, LocalDateTime updatedBefore);

    /**
     * Find payments by gateway transaction reference.
     */
    Optional<Payment> findByGatewayProviderAndExternalReference(String gatewayProvider, String externalReference);

    /**
     * Paginate payments by status.
     */
    Page<Payment> findByStatusOrderByCreatedAtDesc(PaymentStatus status, Pageable pageable);

    Optional<Payment> findByProviderReference(String providerReference);
}

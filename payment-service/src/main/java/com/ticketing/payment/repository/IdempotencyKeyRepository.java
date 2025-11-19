package com.ticketing.payment.repository;

import com.ticketing.payment.entity.IdempotencyKey;
import com.ticketing.payment.entity.enums.IdempotencyStatus;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

/**
 * Repository interface for {@link IdempotencyKey}.
 *
 * Used by gateway layer to enforce idempotent requests and cleanup expired
 * records.
 */
@Repository
public interface IdempotencyKeyRepository extends JpaRepository<IdempotencyKey, Long> {

    /**
     * Lookup by unique idempotency key value.
     */
    Optional<IdempotencyKey> findByIdempotencyKey(String key);

    /**
     * Check existence for fast fail path.
     */
    boolean existsByIdempotencyKey(String key);

    /**
     * Retrieve keys by user and status (debugging/analytics).
     */
    List<IdempotencyKey> findByUserIdAndStatusOrderByCreatedAtDesc(String userId, IdempotencyStatus status);

    /**
     * Find keys that have expired but not marked completed or failed (cleanup).
     */
    List<IdempotencyKey> findByExpiresAtBefore(LocalDateTime cutoff);

    /**
     * Delete keys older than provided timestamp.
     */
    long deleteByCreatedAtBefore(LocalDateTime cutoff);

    /**
     * Update status in-place to avoid loading entity for simple state transitions.
     */
    @Modifying(clearAutomatically = true, flushAutomatically = true)
    @Transactional
    @Query("UPDATE IdempotencyKey k SET k.status = :status, k.completedAt = :completedAt WHERE k.id = :id")
    int updateStatus(@Param("id") Long id, @Param("status") IdempotencyStatus status,
            @Param("completedAt") LocalDateTime completedAt);
}

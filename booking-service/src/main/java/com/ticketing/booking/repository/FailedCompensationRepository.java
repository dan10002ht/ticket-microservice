package com.ticketing.booking.repository;

import java.time.OffsetDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.ticketing.booking.entity.FailedCompensation;
import com.ticketing.booking.entity.FailedCompensation.CompensationStatus;

@Repository
public interface FailedCompensationRepository extends JpaRepository<FailedCompensation, Long> {

    Optional<FailedCompensation> findByCompensationId(UUID compensationId);

    List<FailedCompensation> findByStatus(CompensationStatus status);

    /**
     * Find compensations ready for retry (status=PENDING and next_retry_at <= now)
     */
    @Query("SELECT fc FROM FailedCompensation fc " +
            "WHERE fc.status = :status AND fc.nextRetryAt <= :now " +
            "ORDER BY fc.nextRetryAt ASC")
    List<FailedCompensation> findReadyForRetry(
            @Param("status") CompensationStatus status,
            @Param("now") OffsetDateTime now);

    /**
     * Find compensations that need manual intervention
     */
    List<FailedCompensation> findByStatusIn(List<CompensationStatus> statuses);

    /**
     * Count pending compensations
     */
    long countByStatus(CompensationStatus status);

    /**
     * Find by booking reference for debugging
     */
    List<FailedCompensation> findByBookingReference(String bookingReference);
}

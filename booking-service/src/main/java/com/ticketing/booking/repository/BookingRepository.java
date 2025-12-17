package com.ticketing.booking.repository;

import java.util.Optional;
import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.ticketing.booking.entity.Booking;

@Repository
public interface BookingRepository extends JpaRepository<Booking, Long> {

    Optional<Booking> findByBookingId(UUID bookingId);

    Optional<Booking> findByBookingReference(String bookingReference);

    /**
     * Find booking by idempotency key for duplicate detection
     */
    Optional<Booking> findByIdempotencyKey(String idempotencyKey);

    /**
     * Check if a booking with this idempotency key already exists
     */
    boolean existsByIdempotencyKey(String idempotencyKey);
}


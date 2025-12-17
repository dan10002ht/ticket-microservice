package com.ticketing.booking.entity;

import java.time.Instant;
import java.util.UUID;

import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.statemachine.BookingEvent;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.PrePersist;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Entity for persisting booking state transitions.
 * Provides audit trail and debugging capability.
 */
@Entity
@Table(name = "booking_state_transitions", indexes = {
    @Index(name = "idx_transition_booking_id", columnList = "booking_id"),
    @Index(name = "idx_transition_created_at", columnList = "created_at")
})
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class BookingStateTransition {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "transition_id", nullable = false, unique = true)
    private UUID transitionId;

    @Column(name = "booking_id", nullable = false)
    private UUID bookingId;

    @Enumerated(EnumType.STRING)
    @Column(name = "from_state", nullable = false, length = 30)
    private BookingStatus fromState;

    @Enumerated(EnumType.STRING)
    @Column(name = "to_state", nullable = false, length = 30)
    private BookingStatus toState;

    @Enumerated(EnumType.STRING)
    @Column(name = "event", length = 50)
    private BookingEvent event;

    @Column(name = "triggered_by", length = 100)
    private String triggeredBy;

    @Column(name = "reason", length = 500)
    private String reason;

    @Column(name = "created_at", nullable = false)
    private Instant createdAt;

    @PrePersist
    protected void onCreate() {
        if (transitionId == null) {
            transitionId = UUID.randomUUID();
        }
        if (createdAt == null) {
            createdAt = Instant.now();
        }
    }
}

package com.ticketing.booking.entity;

import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import com.ticketing.booking.entity.enums.BookingStatus;
import com.ticketing.booking.entity.enums.PaymentStatus;

import jakarta.persistence.CascadeType;
import jakarta.persistence.CollectionTable;
import jakarta.persistence.Column;
import jakarta.persistence.ElementCollection;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToMany;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Table(name = "bookings", indexes = {
        @Index(name = "idx_booking_reference", columnList = "booking_reference", unique = true),
        @Index(name = "idx_booking_user", columnList = "user_id"),
        @Index(name = "idx_booking_event", columnList = "event_id"),
        @Index(name = "idx_booking_status", columnList = "status")
})
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Booking {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "booking_id", nullable = false, unique = true, updatable = false)
    private UUID bookingId;

    @Column(name = "booking_reference", nullable = false, length = 50, unique = true)
    private String bookingReference;

    @Column(name = "user_id", nullable = false, length = 36)
    private String userId;

    @Column(name = "event_id", nullable = false, length = 36)
    private String eventId;

    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false, length = 20)
    private BookingStatus status;

    @Enumerated(EnumType.STRING)
    @Column(name = "payment_status", nullable = false, length = 20)
    private PaymentStatus paymentStatus;

    @Column(name = "payment_reference", length = 100)
    private String paymentReference;

    @Column(name = "total_amount", nullable = false, precision = 12, scale = 2)
    private BigDecimal totalAmount;

    @Column(name = "currency", nullable = false, length = 3)
    private String currency;

    @Column(name = "seat_count", nullable = false)
    private int seatCount;

    @ElementCollection
    @CollectionTable(name = "booking_seats", joinColumns = @JoinColumn(name = "booking_id"))
    @Column(name = "seat_number", length = 32)
    private List<String> seatNumbers = new ArrayList<>();

    @ElementCollection
    @CollectionTable(name = "booking_metadata", joinColumns = @JoinColumn(name = "booking_id"))
    @Column(name = "value")
    private Map<String, String> metadata = new HashMap<>();

    @Column(name = "expires_at", nullable = false)
    private OffsetDateTime expiresAt;

    @Column(name = "confirmed_at")
    private OffsetDateTime confirmedAt;

    @Column(name = "cancelled_at")
    private OffsetDateTime cancelledAt;

    @Column(name = "cancellation_reason")
    private String cancellationReason;

    @Column(name = "created_at", nullable = false)
    private OffsetDateTime createdAt;

    @Column(name = "updated_at", nullable = false)
    private OffsetDateTime updatedAt;

    @OneToMany(mappedBy = "booking", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<BookingItem> items = new ArrayList<>();

    @PrePersist
    public void onCreate() {
        OffsetDateTime now = OffsetDateTime.now();
        this.createdAt = now;
        this.updatedAt = now;
        this.bookingId = bookingId == null ? UUID.randomUUID() : bookingId;
        this.status = status == null ? BookingStatus.PENDING : status;
        this.paymentStatus = paymentStatus == null ? PaymentStatus.NOT_REQUIRED : paymentStatus;
    }

    @PreUpdate
    public void onUpdate() {
        this.updatedAt = OffsetDateTime.now();
    }

    public void confirm(String paymentReference) {
        this.status = BookingStatus.CONFIRMED;
        this.paymentStatus = PaymentStatus.CAPTURED;
        this.paymentReference = paymentReference;
        this.confirmedAt = OffsetDateTime.now();
    }

    public void fail(String reason) {
        this.status = BookingStatus.FAILED;
        this.cancellationReason = reason;
    }

    public void cancel(String reason) {
        this.status = BookingStatus.CANCELLED;
        this.paymentStatus = PaymentStatus.FAILED;
        this.cancellationReason = reason;
        this.cancelledAt = OffsetDateTime.now();
    }
}


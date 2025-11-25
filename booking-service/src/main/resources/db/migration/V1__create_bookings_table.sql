CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    booking_id UUID NOT NULL UNIQUE,
    booking_reference VARCHAR(50) NOT NULL UNIQUE,
    user_id VARCHAR(36) NOT NULL,
    event_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'NOT_REQUIRED',
    payment_reference VARCHAR(100),
    total_amount NUMERIC(12, 2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    seat_count INTEGER NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancellation_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE booking_seats (
    booking_id BIGINT REFERENCES bookings(id) ON DELETE CASCADE,
    seat_number VARCHAR(32) NOT NULL
);

CREATE TABLE booking_metadata (
    booking_id BIGINT REFERENCES bookings(id) ON DELETE CASCADE,
    metadata_key VARCHAR(64) NOT NULL,
    value TEXT,
    CONSTRAINT booking_metadata_pk PRIMARY KEY (booking_id, metadata_key)
);

CREATE UNIQUE INDEX idx_booking_reference ON bookings (booking_reference);
CREATE INDEX idx_booking_user ON bookings (user_id);
CREATE INDEX idx_booking_event ON bookings (event_id);
CREATE INDEX idx_booking_status ON bookings (status);


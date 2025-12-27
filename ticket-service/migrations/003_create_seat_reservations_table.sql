-- Migration: Create seat reservations table
-- Description: Temporary seat reservations during booking process

CREATE TABLE IF NOT EXISTS seat_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_session_id UUID NOT NULL REFERENCES booking_sessions(id) ON DELETE CASCADE,
    event_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    zone_id UUID NOT NULL,
    reservation_token VARCHAR(255) UNIQUE NOT NULL, -- Unique reservation identifier
    status VARCHAR(20) DEFAULT 'reserved', -- 'reserved', 'confirmed', 'released', 'expired'
    reserved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- Reservation timeout
    released_at TIMESTAMP WITH TIME ZONE,
    released_reason VARCHAR(100),
    pricing_category VARCHAR(50) NOT NULL,
    base_price DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    metadata JSONB, -- Additional reservation metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    -- Ensure expires_at is in the future
    CONSTRAINT check_reservation_expires_at CHECK (expires_at > reserved_at),
    -- Ensure final_price is not negative
    CONSTRAINT check_reservation_final_price CHECK (final_price >= 0)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_seat_reservations_booking_session_id ON seat_reservations(booking_session_id);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_event_id ON seat_reservations(event_id);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_seat_id ON seat_reservations(seat_id);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_zone_id ON seat_reservations(zone_id);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_reservation_token ON seat_reservations(reservation_token);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_status ON seat_reservations(status);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_expires_at ON seat_reservations(expires_at);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_reserved_at ON seat_reservations(reserved_at);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_seat_reservations_event_seat ON seat_reservations(event_id, seat_id);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_session_status ON seat_reservations(booking_session_id, status);
CREATE INDEX IF NOT EXISTS idx_seat_reservations_expires_status ON seat_reservations(expires_at, status);

-- Trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_seat_reservations_updated_at ON seat_reservations;
CREATE TRIGGER update_seat_reservations_updated_at
    BEFORE UPDATE ON seat_reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE seat_reservations IS 'Temporary seat reservations during booking process';
COMMENT ON COLUMN seat_reservations.reservation_token IS 'Unique reservation identifier';
COMMENT ON COLUMN seat_reservations.metadata IS 'JSON object with additional reservation metadata'; 
-- Create event seats table (seats belong to event, defined from canvas_config)
CREATE TABLE IF NOT EXISTS event_seats (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id VARCHAR(36) NOT NULL, -- References event_seating_zones.public_id
    seat_number VARCHAR(20) NOT NULL,
    row_number VARCHAR(20),
    coordinates JSONB DEFAULT '{}'::jsonb, -- Canvas coordinates for rendering
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'reserved', 'booked', 'blocked', 'maintenance')),
    pricing_category VARCHAR(50) CHECK (pricing_category IN ('premium', 'standard', 'economy', 'vip')),
    base_price DECIMAL(10,2) CHECK (base_price >= 0),
    final_price DECIMAL(10,2) CHECK (final_price >= 0),
    currency VARCHAR(3) DEFAULT 'USD',
    version INTEGER DEFAULT 1, -- Optimistic locking
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Ensure unique seat per event
    UNIQUE(event_id, zone_id, seat_number)
);

-- Create event seat availability table for real-time tracking
CREATE TABLE IF NOT EXISTS event_seat_availability (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    seat_id VARCHAR(36) NOT NULL, -- References event_seats.public_id
    zone_id VARCHAR(36) NOT NULL, -- References event_seating_zones.public_id
    availability_status VARCHAR(20) DEFAULT 'available' CHECK (availability_status IN ('available', 'reserved', 'booked', 'blocked', 'maintenance')),
    reservation_id UUID, -- Reference to booking session or reservation
    blocked_reason VARCHAR(100), -- Reason if seat is blocked
    blocked_until TIMESTAMP WITH TIME ZONE, -- When block expires
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Ensure unique seat per event
    UNIQUE(event_id, seat_id)
);

-- Indexes for event_seats
CREATE INDEX IF NOT EXISTS idx_event_seats_public_id ON event_seats(public_id);
CREATE INDEX IF NOT EXISTS idx_event_seats_event_id ON event_seats(event_id);
CREATE INDEX IF NOT EXISTS idx_event_seats_zone_id ON event_seats(zone_id);
CREATE INDEX IF NOT EXISTS idx_event_seats_status ON event_seats(status);
CREATE INDEX IF NOT EXISTS idx_event_seats_event_zone ON event_seats(event_id, zone_id);

-- Indexes for event_seat_availability
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_public_id ON event_seat_availability(public_id);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_event_id ON event_seat_availability(event_id);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_seat_id ON event_seat_availability(seat_id);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_zone_id ON event_seat_availability(zone_id);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_status ON event_seat_availability(availability_status);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_reservation_id ON event_seat_availability(reservation_id);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_blocked_until ON event_seat_availability(blocked_until);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_event_status ON event_seat_availability(event_id, availability_status);
CREATE INDEX IF NOT EXISTS idx_event_seat_availability_event_zone_status ON event_seat_availability(event_id, zone_id, availability_status);

-- Triggers for event_seats
DROP TRIGGER IF EXISTS update_event_seats_updated_at ON event_seats;
CREATE TRIGGER update_event_seats_updated_at
    BEFORE UPDATE ON event_seats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_event_seat_availability_updated_at ON event_seat_availability;
CREATE TRIGGER update_event_seat_availability_updated_at
    BEFORE UPDATE ON event_seat_availability
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update last_updated timestamp
CREATE OR REPLACE FUNCTION update_last_updated_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_event_seat_availability_last_updated ON event_seat_availability;
CREATE TRIGGER update_event_seat_availability_last_updated
    BEFORE UPDATE ON event_seat_availability
    FOR EACH ROW
    EXECUTE FUNCTION update_last_updated_column();

-- Add comments
COMMENT ON TABLE event_seats IS 'Seats for events (defined from canvas_config)';
COMMENT ON TABLE event_seat_availability IS 'Real-time seat availability tracking for events';
COMMENT ON COLUMN event_seats.version IS 'Optimistic locking version for concurrent updates';
COMMENT ON COLUMN event_seat_availability.availability_status IS 'Current availability status of the seat';
COMMENT ON COLUMN event_seat_availability.blocked_reason IS 'Reason why seat is blocked (maintenance, vip, etc.)'; 
-- EVENT (venue info nhúng, layout riêng, hybrid id)
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    -- Venue info nhúng
    venue_name VARCHAR(255) NOT NULL,
    venue_address TEXT,
    venue_city VARCHAR(100),
    venue_country VARCHAR(100),
    venue_capacity INTEGER,
    -- Layout
    canvas_config JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ZONE cho event (hybrid id)
CREATE TABLE event_seating_zones (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(50),
    coordinates JSONB NOT NULL,
    seat_count INTEGER NOT NULL,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SEAT cho event (hybrid id)
CREATE TABLE event_seats (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id BIGINT NOT NULL REFERENCES event_seating_zones(id) ON DELETE CASCADE,
    seat_number VARCHAR(20) NOT NULL,
    row_number VARCHAR(10),
    coordinates JSONB NOT NULL,
    -- Status tracking for concurrency control
    status VARCHAR(20) DEFAULT 'available',
    reserved_by BIGINT,
    reserved_until TIMESTAMP,
    booked_by BIGINT,
    booking_id BIGINT,
    version INTEGER DEFAULT 1,
    -- Pricing information
    pricing_category VARCHAR(50),
    base_price DECIMAL(10,2),
    final_price DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'USD',
    -- Metadata and audit
    metadata JSONB,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_event_seats_status ON event_seats(status);
CREATE INDEX idx_event_seats_event_status ON event_seats(event_id, status);
CREATE INDEX idx_event_seats_reserved_by ON event_seats(reserved_by);
CREATE INDEX idx_event_seats_booked_by ON event_seats(booked_by);
CREATE INDEX idx_event_seats_pricing_category ON event_seats(pricing_category);
CREATE INDEX idx_event_seats_zone_id ON event_seats(zone_id);
CREATE INDEX idx_event_seats_event_id ON event_seats(event_id);

-- Constraints
ALTER TABLE event_seats 
ADD CONSTRAINT check_status CHECK (status IN ('available', 'reserved', 'sold', 'blocked')),
ADD CONSTRAINT check_final_price CHECK (final_price >= 0),
ADD CONSTRAINT check_base_price CHECK (base_price >= 0);

-- Comments
COMMENT ON COLUMN event_seats.status IS 'Seat status: available, reserved, sold, blocked';
COMMENT ON COLUMN event_seats.reserved_by IS 'User ID who reserved this seat';
COMMENT ON COLUMN event_seats.reserved_until IS 'Reservation expiration time';
COMMENT ON COLUMN event_seats.booked_by IS 'User ID who booked this seat';
COMMENT ON COLUMN event_seats.booking_id IS 'Booking ID for this seat';
COMMENT ON COLUMN event_seats.version IS 'Version for optimistic locking';
COMMENT ON COLUMN event_seats.pricing_category IS 'Pricing category: premium, standard, economy, vip';
COMMENT ON COLUMN event_seats.base_price IS 'Base price for this seat';
COMMENT ON COLUMN event_seats.final_price IS 'Final price after discounts';
COMMENT ON COLUMN event_seats.metadata IS 'Additional seat metadata'; 
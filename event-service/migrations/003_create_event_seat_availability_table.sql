-- Migration: Create event seat availability table
-- Description: Track seat availability for events

CREATE TABLE IF NOT EXISTS event_seat_availability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    seat_id UUID NOT NULL,
    zone_id UUID NOT NULL,
    availability_status VARCHAR(20) DEFAULT 'available', -- 'available', 'reserved', 'booked', 'blocked'
    reservation_id UUID, -- Reference to booking session or reservation
    blocked_reason VARCHAR(100), -- Reason if seat is blocked
    blocked_until TIMESTAMP WITH TIME ZONE, -- When block expires
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    -- Ensure unique seat per event
    UNIQUE(event_id, seat_id)
);

-- Indexes for performance
CREATE INDEX idx_event_seat_availability_event_id ON event_seat_availability(event_id);
CREATE INDEX idx_event_seat_availability_seat_id ON event_seat_availability(seat_id);
CREATE INDEX idx_event_seat_availability_zone_id ON event_seat_availability(zone_id);
CREATE INDEX idx_event_seat_availability_status ON event_seat_availability(availability_status);
CREATE INDEX idx_event_seat_availability_reservation_id ON event_seat_availability(reservation_id);
CREATE INDEX idx_event_seat_availability_blocked_until ON event_seat_availability(blocked_until);
CREATE INDEX idx_event_seat_availability_last_updated ON event_seat_availability(last_updated);

-- Composite indexes for common queries
CREATE INDEX idx_event_seat_availability_event_status ON event_seat_availability(event_id, availability_status);
CREATE INDEX idx_event_seat_availability_event_zone_status ON event_seat_availability(event_id, zone_id, availability_status);

-- Trigger to update updated_at timestamp
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

CREATE TRIGGER update_event_seat_availability_last_updated
    BEFORE UPDATE ON event_seat_availability
    FOR EACH ROW
    EXECUTE FUNCTION update_last_updated_column();

-- Add comments
COMMENT ON TABLE event_seat_availability IS 'Track seat availability for events';
COMMENT ON COLUMN event_seat_availability.availability_status IS 'Current availability status of the seat';
COMMENT ON COLUMN event_seat_availability.blocked_reason IS 'Reason why seat is blocked (maintenance, vip, etc.)'; 
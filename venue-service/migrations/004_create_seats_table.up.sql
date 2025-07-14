-- Create seats table with Hybrid ID Pattern
CREATE TABLE seats (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    zone_id BIGINT NOT NULL REFERENCES seating_zones(id) ON DELETE CASCADE,
    seat_number VARCHAR(20) NOT NULL,
    row_number VARCHAR(20) NOT NULL,
    seat_type VARCHAR(50) DEFAULT 'standard' CHECK (seat_type IN ('standard', 'wheelchair', 'companion', 'aisle', 'premium', 'vip')),
    coordinates JSONB NOT NULL, -- Canvas coordinates for seat position
    properties JSONB, -- Additional seat properties (width, height, angle, etc.)
    is_active BOOLEAN DEFAULT true,
    is_available BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique seat within zone
    UNIQUE(zone_id, row_number, seat_number)
);

-- Indexes for performance
CREATE INDEX idx_seats_public_id ON seats(public_id);
CREATE INDEX idx_seats_zone_id ON seats(zone_id);
CREATE INDEX idx_seats_row_number ON seats(row_number);
CREATE INDEX idx_seats_seat_number ON seats(seat_number);
CREATE INDEX idx_seats_seat_type ON seats(seat_type);
CREATE INDEX idx_seats_is_active ON seats(is_active);
CREATE INDEX idx_seats_is_available ON seats(is_available);
CREATE INDEX idx_seats_display_order ON seats(display_order);

-- Composite index for seat lookup
CREATE INDEX idx_seats_zone_row_seat ON seats(zone_id, row_number, seat_number);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_seats_updated_at 
    BEFORE UPDATE ON seats 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update zone seat count
CREATE OR REPLACE FUNCTION update_zone_seat_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE seating_zones 
        SET seat_count = seat_count + 1 
        WHERE id = NEW.zone_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE seating_zones 
        SET seat_count = seat_count - 1 
        WHERE id = OLD.zone_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_zone_seat_count_trigger
    AFTER INSERT OR DELETE ON seats
    FOR EACH ROW
    EXECUTE FUNCTION update_zone_seat_count();

-- Add comments
COMMENT ON TABLE seats IS 'Individual seats within seating zones';
COMMENT ON COLUMN seats.coordinates IS 'JSON object with seat position coordinates on canvas';
COMMENT ON COLUMN seats.properties IS 'JSON object with additional seat properties'; 
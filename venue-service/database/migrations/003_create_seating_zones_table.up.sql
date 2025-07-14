-- Create seating zones table with Hybrid ID Pattern
CREATE TABLE seating_zones (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    layout_id BIGINT NOT NULL REFERENCES venue_layouts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    zone_type VARCHAR(50) NOT NULL CHECK (zone_type IN ('orchestra', 'balcony', 'mezzanine', 'vip', 'general', 'premium', 'standard', 'economy')),
    color VARCHAR(7), -- Hex color code for zone visualization
    coordinates JSONB NOT NULL, -- Canvas coordinates for zone boundaries
    seat_count INTEGER NOT NULL DEFAULT 0 CHECK (seat_count >= 0),
    row_count INTEGER,
    seats_per_row INTEGER,
    pricing_category VARCHAR(50) CHECK (pricing_category IN ('premium', 'standard', 'economy', 'vip')),
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_seating_zones_public_id ON seating_zones(public_id);
CREATE INDEX idx_seating_zones_layout_id ON seating_zones(layout_id);
CREATE INDEX idx_seating_zones_zone_type ON seating_zones(zone_type);
CREATE INDEX idx_seating_zones_pricing_category ON seating_zones(pricing_category);
CREATE INDEX idx_seating_zones_is_active ON seating_zones(is_active);
CREATE INDEX idx_seating_zones_display_order ON seating_zones(display_order);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_seating_zones_updated_at 
    BEFORE UPDATE ON seating_zones 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE seating_zones IS 'Seating zones within venue layouts';
COMMENT ON COLUMN seating_zones.coordinates IS 'JSON object with zone boundary coordinates on canvas';
COMMENT ON COLUMN seating_zones.color IS 'Hex color code for zone visualization'; 
-- Create event seating zones table (zones belong to event, not external venue)
CREATE TABLE IF NOT EXISTS event_seating_zones (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(50) NOT NULL CHECK (zone_type IN ('seated', 'standing', 'vip', 'accessible')),
    coordinates JSONB DEFAULT '{}'::jsonb, -- Canvas coordinates for rendering
    seat_count INTEGER DEFAULT 0 CHECK (seat_count >= 0),
    color VARCHAR(20), -- Zone color for UI
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(event_id, name)
);

-- Create event pricing table with Hybrid ID Pattern
CREATE TABLE IF NOT EXISTS event_pricing (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id VARCHAR(36) NOT NULL, -- References event_seating_zones.public_id (string for flexibility)
    pricing_category VARCHAR(50) NOT NULL CHECK (pricing_category IN ('premium', 'standard', 'economy', 'vip')),
    base_price DECIMAL(10,2) NOT NULL CHECK (base_price >= 0),
    currency VARCHAR(3) DEFAULT 'USD' CHECK (currency IN ('USD', 'EUR', 'GBP', 'VND')),
    pricing_rules JSONB DEFAULT '{}'::jsonb, -- Dynamic pricing rules (early bird, group discounts, etc.)
    discount_rules JSONB DEFAULT '{}'::jsonb, -- Discount rules and conditions
    is_active BOOLEAN DEFAULT true,
    valid_from TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36),
    updated_by VARCHAR(36),

    -- Ensure unique pricing per event-zone-category combination
    UNIQUE(event_id, zone_id, pricing_category)
);

-- Indexes for event_seating_zones
CREATE INDEX IF NOT EXISTS idx_event_seating_zones_public_id ON event_seating_zones(public_id);
CREATE INDEX IF NOT EXISTS idx_event_seating_zones_event_id ON event_seating_zones(event_id);

-- Indexes for event_pricing
CREATE INDEX IF NOT EXISTS idx_event_pricing_public_id ON event_pricing(public_id);
CREATE INDEX IF NOT EXISTS idx_event_pricing_event_id ON event_pricing(event_id);
CREATE INDEX IF NOT EXISTS idx_event_pricing_zone_id ON event_pricing(zone_id);
CREATE INDEX IF NOT EXISTS idx_event_pricing_pricing_category ON event_pricing(pricing_category);
CREATE INDEX IF NOT EXISTS idx_event_pricing_is_active ON event_pricing(is_active);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_event_pricing_event_active ON event_pricing(event_id, is_active);
CREATE INDEX IF NOT EXISTS idx_event_pricing_event_zone ON event_pricing(event_id, zone_id);

-- Triggers for event_seating_zones
DROP TRIGGER IF EXISTS update_event_seating_zones_updated_at ON event_seating_zones;
CREATE TRIGGER update_event_seating_zones_updated_at
    BEFORE UPDATE ON event_seating_zones
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_event_pricing_updated_at ON event_pricing;
CREATE TRIGGER update_event_pricing_updated_at
    BEFORE UPDATE ON event_pricing
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE event_seating_zones IS 'Seating zones for events (embedded in event, not from external venue service)';
COMMENT ON TABLE event_pricing IS 'Dynamic pricing rules for events';
COMMENT ON COLUMN event_pricing.zone_id IS 'Reference to event_seating_zones.public_id';
COMMENT ON COLUMN event_pricing.pricing_rules IS 'JSON object with dynamic pricing rules';
COMMENT ON COLUMN event_pricing.discount_rules IS 'JSON object with discount rules and conditions'; 
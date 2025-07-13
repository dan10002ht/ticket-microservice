-- Migration: Create event pricing table
-- Description: Dynamic pricing rules for events

CREATE TABLE IF NOT EXISTS event_pricing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id UUID NOT NULL,
    pricing_category VARCHAR(50) NOT NULL, -- 'premium', 'standard', 'economy', 'vip'
    base_price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    pricing_rules JSONB, -- Dynamic pricing rules (early bird, group discounts, etc.)
    discount_rules JSONB, -- Discount rules and conditions
    is_active BOOLEAN DEFAULT true,
    valid_from TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    -- Ensure unique pricing per event-zone-category combination
    UNIQUE(event_id, zone_id, pricing_category)
);

-- Indexes for performance
CREATE INDEX idx_event_pricing_event_id ON event_pricing(event_id);
CREATE INDEX idx_event_pricing_zone_id ON event_pricing(zone_id);
CREATE INDEX idx_event_pricing_pricing_category ON event_pricing(pricing_category);
CREATE INDEX idx_event_pricing_is_active ON event_pricing(is_active);
CREATE INDEX idx_event_pricing_valid_from ON event_pricing(valid_from);
CREATE INDEX idx_event_pricing_valid_until ON event_pricing(valid_until);

-- Composite indexes for common queries
CREATE INDEX idx_event_pricing_event_active ON event_pricing(event_id, is_active);
CREATE INDEX idx_event_pricing_event_zone ON event_pricing(event_id, zone_id);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_event_pricing_updated_at 
    BEFORE UPDATE ON event_pricing 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE event_pricing IS 'Dynamic pricing rules for events';
COMMENT ON COLUMN event_pricing.pricing_rules IS 'JSON object with dynamic pricing rules';
COMMENT ON COLUMN event_pricing.discount_rules IS 'JSON object with discount rules and conditions'; 
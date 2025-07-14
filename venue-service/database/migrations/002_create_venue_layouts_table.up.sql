-- Create venue layouts table with Hybrid ID Pattern
CREATE TABLE venue_layouts (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    venue_id BIGINT NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    layout_type VARCHAR(50) NOT NULL CHECK (layout_type IN ('theater', 'stadium', 'conference', 'custom')),
    canvas_config JSONB NOT NULL, -- Canvas configuration with dimensions, background, etc.
    seating_config JSONB NOT NULL, -- Seating arrangement configuration
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure only one default layout per venue
    UNIQUE(venue_id, is_default) DEFERRABLE INITIALLY DEFERRED
);

-- Indexes for performance
CREATE INDEX idx_venue_layouts_public_id ON venue_layouts(public_id);
CREATE INDEX idx_venue_layouts_venue_id ON venue_layouts(venue_id);
CREATE INDEX idx_venue_layouts_is_active ON venue_layouts(is_active);
CREATE INDEX idx_venue_layouts_is_default ON venue_layouts(is_default);
CREATE INDEX idx_venue_layouts_layout_type ON venue_layouts(layout_type);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_venue_layouts_updated_at 
    BEFORE UPDATE ON venue_layouts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to ensure only one default layout per venue
CREATE OR REPLACE FUNCTION ensure_single_default_layout()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = true THEN
        -- Set all other layouts for this venue to non-default
        UPDATE venue_layouts 
        SET is_default = false 
        WHERE venue_id = NEW.venue_id AND id != NEW.id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER ensure_single_default_layout_trigger
    BEFORE INSERT OR UPDATE ON venue_layouts
    FOR EACH ROW
    EXECUTE FUNCTION ensure_single_default_layout();

-- Add comments
COMMENT ON TABLE venue_layouts IS 'Canvas-based seating layouts for venues';
COMMENT ON COLUMN venue_layouts.canvas_config IS 'JSON configuration for canvas dimensions, background, etc.';
COMMENT ON COLUMN venue_layouts.seating_config IS 'JSON configuration for seating arrangement'; 
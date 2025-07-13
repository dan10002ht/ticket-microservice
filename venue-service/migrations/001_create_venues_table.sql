-- Migration: Create venues table
-- Description: Core venue information table

CREATE TABLE IF NOT EXISTS venues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    capacity INTEGER NOT NULL,
    venue_type VARCHAR(50) NOT NULL, -- 'stadium', 'theater', 'conference_center', 'arena', 'outdoor'
    amenities JSONB, -- Store amenities as JSON array
    images JSONB, -- Store image URLs as JSON array
    coordinates JSONB, -- Store lat/lng as JSON object
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'maintenance'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Indexes for performance
CREATE INDEX idx_venues_city ON venues(city);
CREATE INDEX idx_venues_country ON venues(country);
CREATE INDEX idx_venues_venue_type ON venues(venue_type);
CREATE INDEX idx_venues_status ON venues(status);
CREATE INDEX idx_venues_created_at ON venues(created_at);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_venues_updated_at 
    BEFORE UPDATE ON venues 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE venues IS 'Core venue information table';
COMMENT ON COLUMN venues.amenities IS 'JSON array of venue amenities';
COMMENT ON COLUMN venues.images IS 'JSON array of image URLs';
COMMENT ON COLUMN venues.coordinates IS 'JSON object with lat/lng coordinates'; 
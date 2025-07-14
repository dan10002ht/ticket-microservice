-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create venues table with Hybrid ID Pattern
CREATE TABLE venues (
    id BIGSERIAL PRIMARY KEY,                           -- Internal ID for performance
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(), -- Public ID for API
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
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    venue_type VARCHAR(50) NOT NULL CHECK (venue_type IN ('stadium', 'theater', 'conference_center', 'arena', 'auditorium', 'outdoor', 'other')),
    amenities JSONB DEFAULT '[]'::jsonb, -- Store amenities as JSON array
    images JSONB DEFAULT '[]'::jsonb, -- Store image URLs as JSON array
    coordinates JSONB, -- Store lat/lng as JSON object
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_venues_public_id ON venues(public_id);
CREATE INDEX idx_venues_city ON venues(city);
CREATE INDEX idx_venues_country ON venues(country);
CREATE INDEX idx_venues_venue_type ON venues(venue_type);
CREATE INDEX idx_venues_status ON venues(status);
CREATE INDEX idx_venues_created_at ON venues(created_at);

-- Add updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_venues_updated_at 
    BEFORE UPDATE ON venues 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE venues IS 'Core venue information table';
COMMENT ON COLUMN venues.amenities IS 'JSON array of venue amenities';
COMMENT ON COLUMN venues.images IS 'JSON array of image URLs';
COMMENT ON COLUMN venues.coordinates IS 'JSON object with lat/lng coordinates'; 
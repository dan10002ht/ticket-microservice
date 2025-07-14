-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create events table with Hybrid ID Pattern
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    venue_id BIGINT NOT NULL, -- Will reference venues.id from venue-service
    layout_id BIGINT, -- Will reference venue_layouts.id from venue-service
    name VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('concert', 'theater', 'sports', 'conference', 'exhibition', 'comedy', 'workshop', 'other')),
    category VARCHAR(100),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    doors_open TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'active', 'cancelled', 'completed', 'postponed')),
    max_capacity INTEGER CHECK (max_capacity > 0),
    current_capacity INTEGER DEFAULT 0 CHECK (current_capacity >= 0),
    min_age INTEGER,
    is_featured BOOLEAN DEFAULT false,
    images JSONB DEFAULT '[]'::jsonb, -- Store image URLs as JSON array
    tags JSONB DEFAULT '[]'::jsonb, -- Store tags as JSON array
    metadata JSONB DEFAULT '{}'::jsonb, -- Additional event metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure end_date is after start_date
    CONSTRAINT check_event_dates CHECK (end_date > start_date),
    CONSTRAINT check_doors_open CHECK (doors_open IS NULL OR doors_open <= start_date),
    CONSTRAINT check_capacity CHECK (current_capacity <= max_capacity OR max_capacity IS NULL)
);

-- Indexes for performance
CREATE INDEX idx_events_public_id ON events(public_id);
CREATE INDEX idx_events_venue_id ON events(venue_id);
CREATE INDEX idx_events_layout_id ON events(layout_id);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_date ON events(start_date);
CREATE INDEX idx_events_end_date ON events(end_date);
CREATE INDEX idx_events_is_featured ON events(is_featured);
CREATE INDEX idx_events_created_at ON events(created_at);

-- Composite indexes for common queries
CREATE INDEX idx_events_status_start_date ON events(status, start_date);
CREATE INDEX idx_events_venue_status ON events(venue_id, status);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE events IS 'Core event information table';
COMMENT ON COLUMN events.images IS 'JSON array of image URLs';
COMMENT ON COLUMN events.tags IS 'JSON array of event tags';
COMMENT ON COLUMN events.metadata IS 'JSON object with additional event metadata'; 
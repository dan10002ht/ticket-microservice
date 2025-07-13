-- Migration: Create events table
-- Description: Core event information table

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id UUID NOT NULL,
    layout_id UUID,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) NOT NULL, -- 'concert', 'theater', 'sports', 'conference', 'exhibition'
    category VARCHAR(100),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    doors_open TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'draft', -- 'draft', 'published', 'active', 'cancelled', 'completed'
    max_capacity INTEGER,
    current_capacity INTEGER DEFAULT 0,
    min_age INTEGER,
    is_featured BOOLEAN DEFAULT false,
    images JSONB, -- Store image URLs as JSON array
    tags JSONB, -- Store tags as JSON array
    metadata JSONB, -- Additional event metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    -- Ensure end_date is after start_date
    CONSTRAINT check_event_dates CHECK (end_date > start_date),
    CONSTRAINT check_doors_open CHECK (doors_open IS NULL OR doors_open <= start_date)
);

-- Indexes for performance
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
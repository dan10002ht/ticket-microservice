-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create events table with Hybrid ID Pattern
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    organization_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Embedded venue information (no FK to external venue service)
    venue_name VARCHAR(255) NOT NULL,
    venue_address TEXT,
    venue_city VARCHAR(100),
    venue_country VARCHAR(100),
    venue_capacity INTEGER CHECK (venue_capacity > 0),

    -- Canvas/seating configuration stored as JSON
    canvas_config JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Ensure end_date is after start_date
    CONSTRAINT check_event_dates CHECK (end_date > start_date)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_events_public_id ON events(public_id);
CREATE INDEX IF NOT EXISTS idx_events_organization_id ON events(organization_id);
CREATE INDEX IF NOT EXISTS idx_events_start_date ON events(start_date);
CREATE INDEX IF NOT EXISTS idx_events_end_date ON events(end_date);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);
CREATE INDEX IF NOT EXISTS idx_events_venue_city ON events(venue_city);

-- Trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
CREATE TRIGGER update_events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE events IS 'Core event information table with embedded venue data';
COMMENT ON COLUMN events.canvas_config IS 'JSON object containing seating/canvas configuration'; 
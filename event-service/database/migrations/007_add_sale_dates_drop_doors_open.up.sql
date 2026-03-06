-- Add sale window dates for controlling when tickets are available for purchase
ALTER TABLE events ADD COLUMN IF NOT EXISTS sale_start_date TIMESTAMP WITH TIME ZONE;
ALTER TABLE events ADD COLUMN IF NOT EXISTS sale_end_date TIMESTAMP WITH TIME ZONE;

-- Drop doors_open (display-only info, moved to metadata)
ALTER TABLE events DROP COLUMN IF EXISTS doors_open;

-- Indexes for sale window queries
CREATE INDEX IF NOT EXISTS idx_events_sale_start_date ON events(sale_start_date);
CREATE INDEX IF NOT EXISTS idx_events_sale_end_date ON events(sale_end_date);

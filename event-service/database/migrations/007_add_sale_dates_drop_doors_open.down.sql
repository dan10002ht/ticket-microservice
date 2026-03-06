-- Re-add doors_open
ALTER TABLE events ADD COLUMN IF NOT EXISTS doors_open TIMESTAMP WITH TIME ZONE;

-- Drop sale window columns
DROP INDEX IF EXISTS idx_events_sale_start_date;
DROP INDEX IF EXISTS idx_events_sale_end_date;
ALTER TABLE events DROP COLUMN IF EXISTS sale_start_date;
ALTER TABLE events DROP COLUMN IF EXISTS sale_end_date;

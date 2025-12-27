-- Drop trigger first
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop table
DROP TABLE IF EXISTS events; 
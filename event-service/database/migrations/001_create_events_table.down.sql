-- Drop trigger first
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

-- Drop table
DROP TABLE IF EXISTS events; 
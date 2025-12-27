-- Drop trigger first
DROP TRIGGER IF EXISTS update_event_schedules_updated_at ON event_schedules;

-- Drop table
DROP TABLE IF EXISTS event_schedules;

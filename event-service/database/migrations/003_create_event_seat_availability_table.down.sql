-- Drop triggers first
DROP TRIGGER IF EXISTS update_event_seat_availability_last_updated ON event_seat_availability;
DROP TRIGGER IF EXISTS update_event_seat_availability_updated_at ON event_seat_availability;

-- Drop functions
DROP FUNCTION IF EXISTS update_last_updated_column();

-- Drop table
DROP TABLE IF EXISTS event_seat_availability; 
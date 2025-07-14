-- Drop triggers first
DROP TRIGGER IF EXISTS update_zone_seat_count_trigger ON seats;
DROP TRIGGER IF EXISTS update_seats_updated_at ON seats;

-- Drop functions
DROP FUNCTION IF EXISTS update_zone_seat_count();

-- Drop table
DROP TABLE IF EXISTS seats; 
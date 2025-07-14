-- Drop trigger first
DROP TRIGGER IF EXISTS update_seating_zones_updated_at ON seating_zones;

-- Drop table
DROP TABLE IF EXISTS seating_zones; 
-- Drop trigger first
DROP TRIGGER IF EXISTS update_venues_updated_at ON venues;

-- Drop table
DROP TABLE IF EXISTS venues; 
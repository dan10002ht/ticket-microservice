-- Drop trigger first
DROP TRIGGER IF EXISTS update_event_pricing_updated_at ON event_pricing;

-- Drop table
DROP TABLE IF EXISTS event_pricing; 
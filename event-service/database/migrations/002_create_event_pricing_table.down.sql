-- Drop triggers first
DROP TRIGGER IF EXISTS update_event_pricing_updated_at ON event_pricing;
DROP TRIGGER IF EXISTS update_event_seating_zones_updated_at ON event_seating_zones;

-- Drop tables (in correct order due to dependencies)
DROP TABLE IF EXISTS event_pricing;
DROP TABLE IF EXISTS event_seating_zones; 
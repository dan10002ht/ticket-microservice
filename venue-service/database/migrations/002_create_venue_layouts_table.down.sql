-- Drop triggers first
DROP TRIGGER IF EXISTS ensure_single_default_layout_trigger ON venue_layouts;
DROP TRIGGER IF EXISTS update_venue_layouts_updated_at ON venue_layouts;

-- Drop functions
DROP FUNCTION IF EXISTS ensure_single_default_layout();

-- Drop table
DROP TABLE IF EXISTS venue_layouts; 
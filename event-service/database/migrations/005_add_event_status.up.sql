ALTER TABLE events ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'draft';
CREATE INDEX IF NOT EXISTS idx_events_status ON events (status);

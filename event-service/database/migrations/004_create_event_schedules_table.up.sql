-- Create event schedules table for recurring events
CREATE TABLE IF NOT EXISTS event_schedules (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    schedule_type VARCHAR(50) NOT NULL CHECK (schedule_type IN ('one_time', 'daily', 'weekly', 'monthly', 'custom')),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    recurrence_rule TEXT, -- iCalendar RRULE format for complex recurrence
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_event_schedules_public_id ON event_schedules(public_id);
CREATE INDEX IF NOT EXISTS idx_event_schedules_event_id ON event_schedules(event_id);
CREATE INDEX IF NOT EXISTS idx_event_schedules_schedule_type ON event_schedules(schedule_type);
CREATE INDEX IF NOT EXISTS idx_event_schedules_is_active ON event_schedules(is_active);
CREATE INDEX IF NOT EXISTS idx_event_schedules_start_date ON event_schedules(start_date);

-- Trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_event_schedules_updated_at ON event_schedules;
CREATE TRIGGER update_event_schedules_updated_at
    BEFORE UPDATE ON event_schedules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE event_schedules IS 'Schedules for recurring or multi-date events';
COMMENT ON COLUMN event_schedules.schedule_type IS 'Type of schedule: one_time, daily, weekly, monthly, custom';
COMMENT ON COLUMN event_schedules.recurrence_rule IS 'iCalendar RRULE format for complex recurrence patterns';

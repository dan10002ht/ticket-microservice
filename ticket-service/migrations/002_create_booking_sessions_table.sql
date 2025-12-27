-- Migration: Create booking sessions table
-- Description: Manage booking sessions with timeout logic

CREATE TABLE IF NOT EXISTS booking_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    session_token VARCHAR(255) UNIQUE NOT NULL, -- Unique session identifier
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'completed', 'expired', 'cancelled'
    seat_count INTEGER NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- Session timeout
    completed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancelled_reason VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    metadata JSONB, -- Additional session metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    -- Ensure expires_at is in the future
    CONSTRAINT check_expires_at CHECK (expires_at > created_at)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_booking_sessions_user_id ON booking_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_event_id ON booking_sessions(event_id);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_session_token ON booking_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_status ON booking_sessions(status);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_expires_at ON booking_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_created_at ON booking_sessions(created_at);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_booking_sessions_user_status ON booking_sessions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_event_status ON booking_sessions(event_id, status);
CREATE INDEX IF NOT EXISTS idx_booking_sessions_expires_status ON booking_sessions(expires_at, status);

-- Trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_booking_sessions_updated_at ON booking_sessions;
CREATE TRIGGER update_booking_sessions_updated_at
    BEFORE UPDATE ON booking_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE booking_sessions IS 'Manage booking sessions with timeout logic';
COMMENT ON COLUMN booking_sessions.session_token IS 'Unique session identifier';
COMMENT ON COLUMN booking_sessions.metadata IS 'JSON object with additional session metadata'; 
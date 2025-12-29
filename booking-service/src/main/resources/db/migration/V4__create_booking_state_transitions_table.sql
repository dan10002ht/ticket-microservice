-- V4: Create booking_state_transitions table for audit trail
-- This table records all state transitions for debugging and compliance

CREATE TABLE IF NOT EXISTS booking_state_transitions (
    id BIGSERIAL PRIMARY KEY,
    transition_id UUID NOT NULL UNIQUE,
    booking_id UUID NOT NULL,
    from_state VARCHAR(30) NOT NULL,
    to_state VARCHAR(30) NOT NULL,
    event VARCHAR(50),
    triggered_by VARCHAR(100),
    reason VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_state_from CHECK (from_state IN (
        'PENDING', 'RESERVING', 'SEATS_RESERVED',
        'AWAITING_PAYMENT', 'PAYMENT_PENDING',
        'PROCESSING_PAYMENT', 'PAYMENT_PROCESSING',
        'PAYMENT_FAILED', 'CONFIRMED', 'CANCELLED', 'FAILED', 'EXPIRED'
    )),
    CONSTRAINT chk_state_to CHECK (to_state IN (
        'PENDING', 'RESERVING', 'SEATS_RESERVED',
        'AWAITING_PAYMENT', 'PAYMENT_PENDING',
        'PROCESSING_PAYMENT', 'PAYMENT_PROCESSING',
        'PAYMENT_FAILED', 'CONFIRMED', 'CANCELLED', 'FAILED', 'EXPIRED'
    ))
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_transition_booking_id ON booking_state_transitions(booking_id);
CREATE INDEX IF NOT EXISTS idx_transition_created_at ON booking_state_transitions(created_at);
CREATE INDEX IF NOT EXISTS idx_transition_to_state ON booking_state_transitions(to_state);

-- Partitioning by month for better cleanup (optional - can be added later)
-- For now, using simple table with cleanup job

COMMENT ON TABLE booking_state_transitions IS 'Audit trail of all booking state transitions';
COMMENT ON COLUMN booking_state_transitions.transition_id IS 'Unique ID for this transition';
COMMENT ON COLUMN booking_state_transitions.booking_id IS 'Reference to the booking';
COMMENT ON COLUMN booking_state_transitions.from_state IS 'State before transition';
COMMENT ON COLUMN booking_state_transitions.to_state IS 'State after transition';
COMMENT ON COLUMN booking_state_transitions.event IS 'Event that triggered the transition';
COMMENT ON COLUMN booking_state_transitions.triggered_by IS 'Who/what triggered this transition';
COMMENT ON COLUMN booking_state_transitions.reason IS 'Optional reason for the transition';

-- Failed Compensations table (Dead Letter Queue for saga compensations)
CREATE TABLE IF NOT EXISTS failed_compensations (
    id BIGSERIAL PRIMARY KEY,
    compensation_id UUID NOT NULL UNIQUE,
    booking_id BIGINT,
    booking_reference VARCHAR(50),
    compensation_type VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    reference_id VARCHAR(100),
    error_message VARCHAR(1000),
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 5,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    last_retry_at TIMESTAMP WITH TIME ZONE,
    payload TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fc_status ON failed_compensations (status);
CREATE INDEX IF NOT EXISTS idx_fc_type ON failed_compensations (compensation_type);
CREATE INDEX IF NOT EXISTS idx_fc_next_retry ON failed_compensations (next_retry_at);

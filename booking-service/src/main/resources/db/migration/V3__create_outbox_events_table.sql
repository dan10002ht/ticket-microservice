-- Outbox Events Table for Transactional Outbox Pattern
-- Ensures exactly-once delivery between DB commit and Kafka publish

CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP WITH TIME ZONE,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT,
    CONSTRAINT chk_outbox_status CHECK (status IN ('PENDING', 'PUBLISHED', 'FAILED'))
);

-- Index for polling pending events (most frequent query)
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox_events (status, created_at)
WHERE status = 'PENDING';

-- Index for cleanup of old published events
CREATE INDEX IF NOT EXISTS idx_outbox_published ON outbox_events (status, published_at)
WHERE status = 'PUBLISHED';

-- Index for aggregate lookup (debugging/admin)
CREATE INDEX IF NOT EXISTS idx_outbox_aggregate ON outbox_events (aggregate_type, aggregate_id);

-- Comment for documentation
COMMENT ON TABLE outbox_events IS 'Transactional outbox for reliable event publishing to Kafka';
COMMENT ON COLUMN outbox_events.aggregate_type IS 'Type of aggregate (BOOKING, PAYMENT, etc.)';
COMMENT ON COLUMN outbox_events.aggregate_id IS 'ID of the aggregate (booking_id, payment_id, etc.)';
COMMENT ON COLUMN outbox_events.event_type IS 'Event type (BOOKING_CREATED, BOOKING_CONFIRMED, etc.)';
COMMENT ON COLUMN outbox_events.payload IS 'Full event payload as JSON';

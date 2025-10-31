-- V4: Create idempotency_keys table
-- Payment Service - Idempotency tracking to prevent duplicate payments

CREATE TABLE idempotency_keys (
    -- Primary keys
    id BIGSERIAL PRIMARY KEY,
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,

    -- Request tracking
    request_path VARCHAR(500) NOT NULL, -- API endpoint
    request_method VARCHAR(10) NOT NULL, -- GET, POST, PUT, DELETE
    request_body JSONB, -- Original request body
    request_headers JSONB, -- Original request headers
    
    -- Response tracking
    response_status INTEGER, -- HTTP status code
    response_body JSONB, -- Response body
    response_headers JSONB, -- Response headers
    
    -- Associated entities
    payment_id BIGINT REFERENCES payments(id) ON DELETE SET NULL,
    payment_uuid UUID, -- Denormalized for quick lookup
    refund_id BIGINT REFERENCES refunds(id) ON DELETE SET NULL,
    refund_uuid UUID, -- Denormalized for quick lookup
    
    -- User context
    user_id VARCHAR(255),
    ip_address VARCHAR(45), -- IPv6 support
    user_agent TEXT,
    
    -- Status tracking
    status VARCHAR(50) NOT NULL DEFAULT 'processing', -- processing, completed, failed
    
    -- Expiration
    expires_at TIMESTAMP NOT NULL,
    
    -- Metadata
    metadata JSONB, -- Additional flexible data
    
    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_idempotency_keys_key ON idempotency_keys(idempotency_key);
CREATE INDEX idx_idempotency_keys_payment_id ON idempotency_keys(payment_id);
CREATE INDEX idx_idempotency_keys_payment_uuid ON idempotency_keys(payment_uuid);
CREATE INDEX idx_idempotency_keys_refund_id ON idempotency_keys(refund_id);
CREATE INDEX idx_idempotency_keys_refund_uuid ON idempotency_keys(refund_uuid);
CREATE INDEX idx_idempotency_keys_user_id ON idempotency_keys(user_id);
CREATE INDEX idx_idempotency_keys_status ON idempotency_keys(status);
CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);
CREATE INDEX idx_idempotency_keys_created_at ON idempotency_keys(created_at);

-- Composite indexes for common queries
CREATE INDEX idx_idempotency_keys_user_created ON idempotency_keys(user_id, created_at DESC);
CREATE INDEX idx_idempotency_keys_status_expires ON idempotency_keys(status, expires_at);

-- JSONB GIN indexes for flexible querying
CREATE INDEX idx_idempotency_keys_request_body ON idempotency_keys USING GIN (request_body);
CREATE INDEX idx_idempotency_keys_response_body ON idempotency_keys USING GIN (response_body);
CREATE INDEX idx_idempotency_keys_metadata ON idempotency_keys USING GIN (metadata);

-- Comments for documentation
COMMENT ON TABLE idempotency_keys IS 'Idempotency tracking to prevent duplicate payment requests';
COMMENT ON COLUMN idempotency_keys.idempotency_key IS 'Unique key provided by client to ensure idempotency';
COMMENT ON COLUMN idempotency_keys.request_path IS 'API endpoint path';
COMMENT ON COLUMN idempotency_keys.request_method IS 'HTTP method: GET, POST, PUT, DELETE';
COMMENT ON COLUMN idempotency_keys.request_body IS 'Original request body in JSON format';
COMMENT ON COLUMN idempotency_keys.response_status IS 'HTTP response status code';
COMMENT ON COLUMN idempotency_keys.response_body IS 'Response body in JSON format';
COMMENT ON COLUMN idempotency_keys.status IS 'Status: processing, completed, failed';
COMMENT ON COLUMN idempotency_keys.expires_at IS 'Expiration timestamp for cleanup';
COMMENT ON COLUMN idempotency_keys.metadata IS 'Additional flexible data in JSON format';

-- Trigger to auto-update updated_at
CREATE TRIGGER update_idempotency_keys_updated_at
    BEFORE UPDATE ON idempotency_keys
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to cleanup expired idempotency keys
CREATE OR REPLACE FUNCTION cleanup_expired_idempotency_keys()
RETURNS void AS $$
BEGIN
    DELETE FROM idempotency_keys
    WHERE expires_at < CURRENT_TIMESTAMP
    AND status = 'completed';
END;
$$ LANGUAGE plpgsql;

-- Comment on cleanup function
COMMENT ON FUNCTION cleanup_expired_idempotency_keys() IS 'Cleanup expired idempotency keys to prevent table bloat';




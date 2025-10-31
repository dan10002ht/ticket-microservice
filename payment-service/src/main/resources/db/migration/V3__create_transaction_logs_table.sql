-- V3: Create transaction_logs table
-- Payment Service - Audit trail for all payment transactions

CREATE TABLE transaction_logs (
    -- Primary keys
    id BIGSERIAL PRIMARY KEY,
    log_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Transaction reference
    payment_id BIGINT REFERENCES payments(id) ON DELETE SET NULL,
    payment_uuid UUID, -- Denormalized for quick lookup
    refund_id BIGINT REFERENCES refunds(id) ON DELETE SET NULL,
    refund_uuid UUID, -- Denormalized for quick lookup

    -- Transaction details
    transaction_type VARCHAR(50) NOT NULL, -- payment_initiated, payment_success, payment_failed, refund_initiated, refund_success, refund_failed, webhook_received
    event_name VARCHAR(100) NOT NULL, -- Specific event name
    
    -- Gateway information
    gateway_provider VARCHAR(50), -- stripe, paypal, vnpay, momo
    external_reference VARCHAR(255), -- Gateway transaction ID

    -- Request/Response data
    request_data JSONB, -- Request payload
    response_data JSONB, -- Response payload
    headers JSONB, -- HTTP headers if applicable

    -- Status and error tracking
    status VARCHAR(50) NOT NULL, -- success, failed, pending
    error_code VARCHAR(100),
    error_message TEXT,
    
    -- Performance tracking
    duration_ms INTEGER, -- Request duration in milliseconds
    
    -- Context information
    user_id VARCHAR(255),
    ip_address VARCHAR(45), -- IPv6 support
    user_agent TEXT,
    correlation_id UUID, -- For distributed tracing
    
    -- Metadata
    metadata JSONB, -- Additional flexible data
    
    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)
);

-- Indexes for performance
CREATE INDEX idx_transaction_logs_payment_id ON transaction_logs(payment_id);
CREATE INDEX idx_transaction_logs_payment_uuid ON transaction_logs(payment_uuid);
CREATE INDEX idx_transaction_logs_refund_id ON transaction_logs(refund_id);
CREATE INDEX idx_transaction_logs_refund_uuid ON transaction_logs(refund_uuid);
CREATE INDEX idx_transaction_logs_transaction_type ON transaction_logs(transaction_type);
CREATE INDEX idx_transaction_logs_event_name ON transaction_logs(event_name);
CREATE INDEX idx_transaction_logs_gateway_provider ON transaction_logs(gateway_provider);
CREATE INDEX idx_transaction_logs_external_reference ON transaction_logs(external_reference);
CREATE INDEX idx_transaction_logs_status ON transaction_logs(status);
CREATE INDEX idx_transaction_logs_user_id ON transaction_logs(user_id);
CREATE INDEX idx_transaction_logs_correlation_id ON transaction_logs(correlation_id);
CREATE INDEX idx_transaction_logs_created_at ON transaction_logs(created_at);

-- Composite indexes for common queries
CREATE INDEX idx_transaction_logs_payment_type ON transaction_logs(payment_id, transaction_type);
CREATE INDEX idx_transaction_logs_gateway_external ON transaction_logs(gateway_provider, external_reference);
CREATE INDEX idx_transaction_logs_user_created ON transaction_logs(user_id, created_at DESC);

-- JSONB GIN indexes for flexible querying
CREATE INDEX idx_transaction_logs_request_data ON transaction_logs USING GIN (request_data);
CREATE INDEX idx_transaction_logs_response_data ON transaction_logs USING GIN (response_data);
CREATE INDEX idx_transaction_logs_metadata ON transaction_logs USING GIN (metadata);

-- Comments for documentation
COMMENT ON TABLE transaction_logs IS 'Audit trail for all payment and refund transactions';
COMMENT ON COLUMN transaction_logs.log_id IS 'Public-facing UUID for log identification';
COMMENT ON COLUMN transaction_logs.transaction_type IS 'Type: payment_initiated, payment_success, payment_failed, refund_initiated, refund_success, refund_failed, webhook_received';
COMMENT ON COLUMN transaction_logs.event_name IS 'Specific event name from gateway or internal';
COMMENT ON COLUMN transaction_logs.gateway_provider IS 'Payment gateway: stripe, paypal, vnpay, momo';
COMMENT ON COLUMN transaction_logs.request_data IS 'Full request payload in JSON format';
COMMENT ON COLUMN transaction_logs.response_data IS 'Full response payload in JSON format';
COMMENT ON COLUMN transaction_logs.duration_ms IS 'Request duration in milliseconds for performance tracking';
COMMENT ON COLUMN transaction_logs.correlation_id IS 'UUID for distributed tracing across services';
COMMENT ON COLUMN transaction_logs.metadata IS 'Additional flexible data in JSON format';




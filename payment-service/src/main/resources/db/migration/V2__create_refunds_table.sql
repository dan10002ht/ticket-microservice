-- V2: Create refunds table
-- Payment Service - Refund transactions table

CREATE TABLE refunds (
    -- Primary keys
    id BIGSERIAL PRIMARY KEY,
    refund_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Payment reference
    payment_id BIGINT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    payment_uuid UUID NOT NULL, -- Denormalized for quick lookup

    -- Refund details
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    refund_type VARCHAR(50) NOT NULL, -- full, partial
    
    -- Status tracking
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, success, failed, cancelled

    -- External references
    external_reference VARCHAR(255), -- Gateway refund transaction ID
    gateway_provider VARCHAR(50), -- stripe, paypal, vnpay, momo
    provider_reference VARCHAR(255), -- Gateway-specific refund identifier

    -- Gateway response
    gateway_response JSONB, -- Full response from payment gateway
    failure_reason TEXT, -- Reason for failure if status = failed

    -- Refund reason
    reason VARCHAR(255) NOT NULL,
    description TEXT,

    -- Metadata
    metadata JSONB, -- Additional flexible data

    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    
    -- Refund completion
    refunded_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancelled_by VARCHAR(255),
    cancellation_reason TEXT
);

-- Indexes for performance
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_payment_uuid ON refunds(payment_uuid);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_refund_id ON refunds(refund_id);
CREATE INDEX idx_refunds_external_reference ON refunds(external_reference);
CREATE INDEX idx_refunds_gateway_provider ON refunds(gateway_provider);
CREATE INDEX idx_refunds_provider_reference ON refunds(provider_reference);
CREATE INDEX idx_refunds_created_at ON refunds(created_at);

-- Composite indexes for common queries
CREATE INDEX idx_refunds_payment_status ON refunds(payment_id, status);
CREATE INDEX idx_refunds_gateway_external ON refunds(gateway_provider, external_reference);

-- Comments for documentation
COMMENT ON TABLE refunds IS 'Refund transactions table';
COMMENT ON COLUMN refunds.refund_id IS 'Public-facing UUID for refund identification';
COMMENT ON COLUMN refunds.payment_id IS 'Foreign key to payments table';
COMMENT ON COLUMN refunds.payment_uuid IS 'Denormalized payment UUID for quick lookup';
COMMENT ON COLUMN refunds.refund_type IS 'Refund type: full, partial';
COMMENT ON COLUMN refunds.status IS 'Refund status: pending, processing, success, failed, cancelled';
COMMENT ON COLUMN refunds.gateway_provider IS 'Payment gateway: stripe, paypal, vnpay, momo';
COMMENT ON COLUMN refunds.provider_reference IS 'Gateway refund identifier for webhook mapping';
COMMENT ON COLUMN refunds.gateway_response IS 'Full JSON response from payment gateway';
COMMENT ON COLUMN refunds.reason IS 'Short reason for refund';
COMMENT ON COLUMN refunds.description IS 'Detailed description of refund';
COMMENT ON COLUMN refunds.metadata IS 'Additional flexible data in JSON format';

-- Trigger to auto-update updated_at
CREATE TRIGGER update_refunds_updated_at
    BEFORE UPDATE ON refunds
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();




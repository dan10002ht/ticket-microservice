-- V1: Create payments table
-- Payment Service - Core payment transactions table

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS payments (
    -- Primary keys
    id BIGSERIAL PRIMARY KEY,
    payment_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Booking reference
    booking_id VARCHAR(255) NOT NULL,
    ticket_id VARCHAR(255),
    user_id VARCHAR(255) NOT NULL,

    -- Payment details
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    payment_method VARCHAR(50) NOT NULL, -- credit_card, debit_card, bank_transfer, e_wallet, cash

    -- Status tracking
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, success, failed, cancelled, refunded

    -- External references
    external_reference VARCHAR(255), -- Gateway transaction ID
    gateway_provider VARCHAR(50), -- stripe, paypal, vnpay, momo

    -- Gateway response
    gateway_response JSONB, -- Full response from payment gateway
    failure_reason TEXT, -- Reason for failure if status = failed

    -- Metadata
    description TEXT,
    metadata JSONB, -- Additional flexible data

    -- Idempotency
    idempotency_key VARCHAR(255) UNIQUE,

    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),

    -- Payment completion
    paid_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancelled_by VARCHAR(255),
    cancellation_reason TEXT
);

-- Indexes for performance
CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_payment_id ON payments(payment_id);
CREATE INDEX idx_payments_external_reference ON payments(external_reference);
CREATE INDEX idx_payments_gateway_provider ON payments(gateway_provider);
CREATE INDEX idx_payments_idempotency_key ON payments(idempotency_key);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- Composite indexes for common queries
CREATE INDEX idx_payments_user_status ON payments(user_id, status);
CREATE INDEX idx_payments_booking_status ON payments(booking_id, status);
CREATE INDEX idx_payments_gateway_external ON payments(gateway_provider, external_reference);

-- Comments for documentation
COMMENT ON TABLE payments IS 'Core payment transactions table';
COMMENT ON COLUMN payments.payment_id IS 'Public-facing UUID for payment identification';
COMMENT ON COLUMN payments.booking_id IS 'Reference to booking in booking service';
COMMENT ON COLUMN payments.ticket_id IS 'Reference to ticket in ticket service';
COMMENT ON COLUMN payments.status IS 'Payment status: pending, processing, success, failed, cancelled, refunded';
COMMENT ON COLUMN payments.payment_method IS 'Payment method: credit_card, debit_card, bank_transfer, e_wallet, cash';
COMMENT ON COLUMN payments.gateway_provider IS 'Payment gateway: stripe, paypal, vnpay, momo';
COMMENT ON COLUMN payments.gateway_response IS 'Full JSON response from payment gateway';
COMMENT ON COLUMN payments.idempotency_key IS 'Unique key to prevent duplicate payments';
COMMENT ON COLUMN payments.metadata IS 'Additional flexible data in JSON format';

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();




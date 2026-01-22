-- V5: Add additional constraints, views, and helper functions
-- Payment Service - Enhanced database features

-- ================================================
-- CONSTRAINTS
-- ================================================

-- Ensure payment amount matches total refunds
ALTER TABLE payments
ADD CONSTRAINT check_payment_amount_positive CHECK (amount > 0);

-- Ensure refund amount doesn't exceed payment amount
-- (This will be enforced in application logic, but document here)

-- Add check constraints for status values
ALTER TABLE payments
ADD CONSTRAINT check_payment_status CHECK (
    status IN ('pending', 'processing', 'success', 'failed', 'cancelled', 'refunded', 'partially_refunded')
);

ALTER TABLE refunds
ADD CONSTRAINT check_refund_status CHECK (
    status IN ('pending', 'processing', 'success', 'failed', 'cancelled')
);

ALTER TABLE refunds
ADD CONSTRAINT check_refund_type CHECK (
    refund_type IN ('full', 'partial')
);

ALTER TABLE transaction_logs
ADD CONSTRAINT check_transaction_log_status CHECK (
    status IN ('success', 'failed', 'pending')
);

ALTER TABLE idempotency_keys
ADD CONSTRAINT check_idempotency_status CHECK (
    status IN ('processing', 'completed', 'failed')
);

-- ================================================
-- VIEWS
-- ================================================

-- View: Payment Summary by User
CREATE OR REPLACE VIEW v_user_payment_summary AS
SELECT
    user_id,
    COUNT(*) as total_payments,
    COUNT(*) FILTER (WHERE status = 'success') as successful_payments,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_payments,
    COUNT(*) FILTER (WHERE status = 'refunded' OR status = 'partially_refunded') as refunded_payments,
    SUM(amount) FILTER (WHERE status = 'success') as total_amount_paid,
    SUM(amount) FILTER (WHERE status = 'refunded') as total_amount_refunded,
    AVG(amount) FILTER (WHERE status = 'success') as average_payment_amount,
    MIN(created_at) as first_payment_date,
    MAX(created_at) as last_payment_date
FROM payments
GROUP BY user_id;

COMMENT ON VIEW v_user_payment_summary IS 'Payment summary statistics by user';

-- View: Payment Summary by Gateway
CREATE OR REPLACE VIEW v_gateway_payment_summary AS
SELECT
    gateway_provider,
    COUNT(*) as total_payments,
    COUNT(*) FILTER (WHERE status = 'success') as successful_payments,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_payments,
    SUM(amount) FILTER (WHERE status = 'success') as total_amount_processed,
    AVG(amount) FILTER (WHERE status = 'success') as average_payment_amount,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'success')::NUMERIC / NULLIF(COUNT(*), 0) * 100), 
        2
    ) as success_rate_percentage
FROM payments
WHERE gateway_provider IS NOT NULL
GROUP BY gateway_provider;

COMMENT ON VIEW v_gateway_payment_summary IS 'Payment summary statistics by gateway provider';

-- View: Daily Payment Statistics
CREATE OR REPLACE VIEW v_daily_payment_statistics AS
SELECT
    DATE(created_at) as payment_date,
    COUNT(*) as total_payments,
    COUNT(*) FILTER (WHERE status = 'success') as successful_payments,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_payments,
    SUM(amount) FILTER (WHERE status = 'success') as daily_revenue,
    AVG(amount) FILTER (WHERE status = 'success') as average_payment_amount,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT booking_id) as unique_bookings
FROM payments
GROUP BY DATE(created_at)
ORDER BY payment_date DESC;

COMMENT ON VIEW v_daily_payment_statistics IS 'Daily payment statistics for reporting';

-- View: Refund Statistics
CREATE OR REPLACE VIEW v_refund_statistics AS
SELECT
    DATE(r.created_at) as refund_date,
    COUNT(*) as total_refunds,
    COUNT(*) FILTER (WHERE r.status = 'success') as successful_refunds,
    COUNT(*) FILTER (WHERE r.refund_type = 'full') as full_refunds,
    COUNT(*) FILTER (WHERE r.refund_type = 'partial') as partial_refunds,
    SUM(r.amount) FILTER (WHERE r.status = 'success') as total_refunded_amount,
    AVG(r.amount) FILTER (WHERE r.status = 'success') as average_refund_amount
FROM refunds r
GROUP BY DATE(r.created_at)
ORDER BY refund_date DESC;

COMMENT ON VIEW v_refund_statistics IS 'Daily refund statistics for reporting';

-- ================================================
-- HELPER FUNCTIONS
-- ================================================

-- Function: Get total refunded amount for a payment
CREATE OR REPLACE FUNCTION get_total_refunded_amount(p_payment_id BIGINT)
RETURNS DECIMAL(10, 2) AS $$
DECLARE
    total_refunded DECIMAL(10, 2);
BEGIN
    SELECT COALESCE(SUM(amount), 0)
    INTO total_refunded
    FROM refunds
    WHERE payment_id = p_payment_id
    AND status = 'success';
    
    RETURN total_refunded;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_total_refunded_amount(BIGINT) IS 'Calculate total refunded amount for a payment';

-- Function: Check if payment can be refunded
CREATE OR REPLACE FUNCTION can_payment_be_refunded(p_payment_id BIGINT)
RETURNS BOOLEAN AS $$
DECLARE
    payment_status VARCHAR(50);
    payment_amount DECIMAL(10, 2);
    total_refunded DECIMAL(10, 2);
BEGIN
    -- Get payment details
    SELECT status, amount
    INTO payment_status, payment_amount
    FROM payments
    WHERE id = p_payment_id;
    
    -- Check if payment exists and is successful
    IF payment_status IS NULL OR payment_status != 'success' THEN
        RETURN FALSE;
    END IF;
    
    -- Get total refunded amount
    total_refunded := get_total_refunded_amount(p_payment_id);
    
    -- Check if there's remaining amount to refund
    RETURN (payment_amount - total_refunded) > 0;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION can_payment_be_refunded(BIGINT) IS 'Check if a payment can be refunded (has remaining amount)';

-- Function: Get remaining refundable amount
CREATE OR REPLACE FUNCTION get_remaining_refundable_amount(p_payment_id BIGINT)
RETURNS DECIMAL(10, 2) AS $$
DECLARE
    payment_amount DECIMAL(10, 2);
    total_refunded DECIMAL(10, 2);
BEGIN
    -- Get payment amount
    SELECT amount
    INTO payment_amount
    FROM payments
    WHERE id = p_payment_id;
    
    -- Get total refunded
    total_refunded := get_total_refunded_amount(p_payment_id);
    
    -- Return remaining amount
    RETURN GREATEST(payment_amount - total_refunded, 0);
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_remaining_refundable_amount(BIGINT) IS 'Calculate remaining refundable amount for a payment';

-- ================================================
-- PERFORMANCE OPTIMIZATION
-- ================================================

-- Partial index for active payments (most queried)
CREATE INDEX idx_payments_active ON payments(id, status, created_at)
WHERE status IN ('pending', 'processing');

-- Partial index for active refunds
CREATE INDEX idx_refunds_active ON refunds(id, status, created_at)
WHERE status IN ('pending', 'processing');

-- Index for recent transaction logs (ordered by created_at for efficient queries)
CREATE INDEX idx_transaction_logs_recent ON transaction_logs(created_at DESC, transaction_type);

-- Index for idempotency keys by expiration and status
CREATE INDEX idx_idempotency_keys_active ON idempotency_keys(expires_at, idempotency_key, status);




-- V6: Add provider_reference column to payments table
-- This column stores the reference ID from the payment provider (e.g., Stripe charge ID)

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS provider_reference VARCHAR(255);

-- Add index for provider reference lookups
CREATE INDEX IF NOT EXISTS idx_payments_provider_reference ON payments(provider_reference)
WHERE provider_reference IS NOT NULL;

COMMENT ON COLUMN payments.provider_reference IS 'Reference ID from payment provider (e.g., Stripe charge ID, PayPal transaction ID)';

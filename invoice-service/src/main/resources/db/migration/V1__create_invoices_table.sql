-- Invoice tables for the invoice-service.
-- invoice.booking_id has a UNIQUE constraint — idempotent invoice generation.

CREATE TABLE invoices (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50)   NOT NULL UNIQUE,
    booking_id     UUID          NOT NULL UNIQUE,
    payment_id     UUID,
    user_id        UUID          NOT NULL,
    event_id       UUID,
    subtotal       DECIMAL(15,2) NOT NULL,
    tax_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount   DECIMAL(15,2) NOT NULL,
    currency       VARCHAR(3)    NOT NULL DEFAULT 'VND',
    status         VARCHAR(50)   NOT NULL DEFAULT 'PENDING',
    issued_at      TIMESTAMPTZ,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_invoice_amounts CHECK (
        subtotal >= 0 AND tax_amount >= 0 AND total_amount >= 0
    )
);

COMMENT ON TABLE  invoices                IS 'One invoice per completed booking';
COMMENT ON COLUMN invoices.invoice_number IS 'Human-readable number: INV-YYYYMM-XXXXXX';
COMMENT ON COLUMN invoices.status         IS 'PENDING | GENERATED | CANCELLED';

CREATE TABLE invoice_items (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id  UUID          NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    description VARCHAR(500)  NOT NULL,
    quantity    INT           NOT NULL DEFAULT 1,
    unit_price  DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,

    CONSTRAINT chk_item_prices CHECK (quantity > 0 AND unit_price >= 0 AND total_price >= 0)
);

COMMENT ON TABLE invoice_items IS 'Line items for each invoice (one row per ticket type)';

-- Indexes
CREATE INDEX idx_invoices_booking_id ON invoices (booking_id);
CREATE INDEX idx_invoices_user_id    ON invoices (user_id);
CREATE INDEX idx_invoices_event_id   ON invoices (event_id);
CREATE INDEX idx_invoices_status     ON invoices (status);

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_invoice_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_invoices_updated_at
    BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION update_invoice_updated_at();

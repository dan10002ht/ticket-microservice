-- Migration: Create tickets table
-- Description: Core ticket information table

CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    zone_id UUID NOT NULL,
    user_id UUID NOT NULL,
    booking_session_id UUID,
    ticket_number VARCHAR(50) UNIQUE NOT NULL, -- Unique ticket identifier
    ticket_type VARCHAR(50) DEFAULT 'standard', -- 'standard', 'vip', 'wheelchair', 'companion'
    pricing_category VARCHAR(50) NOT NULL, -- 'premium', 'standard', 'economy', 'vip'
    base_price DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    discount_amount DECIMAL(10,2) DEFAULT 0,
    discount_reason VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'confirmed', 'cancelled', 'refunded', 'used'
    payment_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'failed', 'refunded'
    payment_method VARCHAR(50),
    payment_reference VARCHAR(100),
    qr_code VARCHAR(255), -- QR code data or URL
    barcode VARCHAR(100), -- Barcode data
    valid_from TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE,
    used_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancelled_reason VARCHAR(100),
    refunded_at TIMESTAMP WITH TIME ZONE,
    refunded_amount DECIMAL(10,2),
    metadata JSONB, -- Additional ticket metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    
    -- Ensure final_price is not negative
    CONSTRAINT check_final_price CHECK (final_price >= 0),
    CONSTRAINT check_discount_amount CHECK (discount_amount >= 0)
);

-- Indexes for performance
CREATE INDEX idx_tickets_event_id ON tickets(event_id);
CREATE INDEX idx_tickets_seat_id ON tickets(seat_id);
CREATE INDEX idx_tickets_zone_id ON tickets(zone_id);
CREATE INDEX idx_tickets_user_id ON tickets(user_id);
CREATE INDEX idx_tickets_booking_session_id ON tickets(booking_session_id);
CREATE INDEX idx_tickets_ticket_number ON tickets(ticket_number);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_payment_status ON tickets(payment_status);
CREATE INDEX idx_tickets_pricing_category ON tickets(pricing_category);
CREATE INDEX idx_tickets_created_at ON tickets(created_at);
CREATE INDEX idx_tickets_valid_until ON tickets(valid_until);

-- Composite indexes for common queries
CREATE INDEX idx_tickets_event_status ON tickets(event_id, status);
CREATE INDEX idx_tickets_user_status ON tickets(user_id, status);
CREATE INDEX idx_tickets_event_user ON tickets(event_id, user_id);
CREATE INDEX idx_tickets_booking_session ON tickets(booking_session_id, status);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_tickets_updated_at 
    BEFORE UPDATE ON tickets 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE tickets IS 'Core ticket information table';
COMMENT ON COLUMN tickets.ticket_number IS 'Unique ticket identifier';
COMMENT ON COLUMN tickets.metadata IS 'JSON object with additional ticket metadata'; 
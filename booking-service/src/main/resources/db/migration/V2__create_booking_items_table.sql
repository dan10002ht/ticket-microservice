CREATE TABLE booking_items (
    id BIGSERIAL PRIMARY KEY,
    booking_fk BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    ticket_type_id VARCHAR(36) NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price NUMERIC(12, 2) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    seat_numbers TEXT
);

CREATE INDEX idx_booking_items_booking_fk ON booking_items (booking_fk);


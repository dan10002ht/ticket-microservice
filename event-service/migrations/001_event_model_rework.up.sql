-- EVENT (venue info nhúng, layout riêng, hybrid id)
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    -- Venue info nhúng
    venue_name VARCHAR(255) NOT NULL,
    venue_address TEXT,
    venue_city VARCHAR(100),
    venue_country VARCHAR(100),
    venue_capacity INTEGER,
    -- Layout
    canvas_config JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ZONE cho event (hybrid id)
CREATE TABLE event_seating_zones (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(50),
    coordinates JSONB NOT NULL,
    seat_count INTEGER NOT NULL,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SEAT cho event (hybrid id)
CREATE TABLE event_seats (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id BIGINT NOT NULL REFERENCES event_seating_zones(id) ON DELETE CASCADE,
    seat_number VARCHAR(20) NOT NULL,
    row_number VARCHAR(10),
    coordinates JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 
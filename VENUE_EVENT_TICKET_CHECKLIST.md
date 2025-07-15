# 🎫 Venue, Event & Ticket System Implementation Checklist

## 📋 Tổng quan

Triển khai hệ thống quản lý venue, event và ticket booking theo kiến trúc microservices với 3 service chính:

- **venue-service**: Quản lý venue và layout
- **event-service**: Quản lý event và pricing
- **ticket-service**: Quản lý ticket booking và seat availability

---

## 🏗️ Phase 1: Database Schema & Infrastructure

### 1.1 Tạo Database Migrations

#### **Hybrid ID Pattern (Tham khảo auth-service)**

Sử dụng pattern **hybrid approach** cho optimal performance và security:

- **`id`** (BIGSERIAL): Auto-incrementing primary key cho internal operations và performance
- **`public_id`** (UUID): Globally unique identifier cho API exposure và security

#### **venue-service Migrations**

```sql
-- 001_create_venue_tables.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- venues table
CREATE TABLE venues (
    id BIGSERIAL PRIMARY KEY,                           -- Internal ID for performance
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(), -- Public ID for API
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20),
    phone VARCHAR(20),
    email VARCHAR(255),
    website_url TEXT,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    venue_type VARCHAR(50) NOT NULL CHECK (venue_type IN ('stadium', 'theater', 'conference_center', 'arena', 'auditorium', 'outdoor', 'other')),
    amenities JSONB DEFAULT '[]'::jsonb, -- ['parking', 'wifi', 'catering', 'accessibility', etc.]
    images JSONB DEFAULT '[]'::jsonb, -- Array of image URLs
    coordinates POINT, -- Latitude/Longitude
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- venue_layouts table
CREATE TABLE venue_layouts (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    venue_id BIGINT NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    canvas_config JSONB NOT NULL, -- Canvas layout configuration
    seat_count INTEGER NOT NULL CHECK (seat_count >= 0),
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- seating_zones table
CREATE TABLE seating_zones (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    layout_id BIGINT NOT NULL REFERENCES venue_layouts(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    zone_type VARCHAR(50) NOT NULL CHECK (zone_type IN ('vip', 'general', 'premium', 'standard', 'wheelchair', 'aisle')),
    coordinates JSONB NOT NULL, -- Zone boundary coordinates
    seat_count INTEGER NOT NULL CHECK (seat_count >= 0),
    color VARCHAR(7), -- Hex color code
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- seats table
CREATE TABLE seats (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    zone_id BIGINT NOT NULL REFERENCES seating_zones(id) ON DELETE CASCADE,
    seat_number VARCHAR(20) NOT NULL,
    row_number VARCHAR(10) NOT NULL,
    seat_type VARCHAR(50) NOT NULL CHECK (seat_type IN ('standard', 'wheelchair', 'aisle', 'premium', 'vip')),
    coordinates JSONB NOT NULL, -- Seat position in canvas
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(zone_id, row_number, seat_number)
);

-- Add indexes for performance
CREATE INDEX idx_venues_public_id ON venues(public_id);
CREATE INDEX idx_venues_city ON venues(city);
CREATE INDEX idx_venues_venue_type ON venues(venue_type);
CREATE INDEX idx_venues_is_active ON venues(is_active);
CREATE INDEX idx_venues_country ON venues(country);

CREATE INDEX idx_venue_layouts_public_id ON venue_layouts(public_id);
CREATE INDEX idx_venue_layouts_venue_id ON venue_layouts(venue_id);
CREATE INDEX idx_venue_layouts_is_default ON venue_layouts(is_default);
CREATE INDEX idx_venue_layouts_is_active ON venue_layouts(is_active);

CREATE INDEX idx_seating_zones_public_id ON seating_zones(public_id);
CREATE INDEX idx_seating_zones_layout_id ON seating_zones(layout_id);
CREATE INDEX idx_seating_zones_zone_type ON seating_zones(zone_type);
CREATE INDEX idx_seating_zones_is_active ON seating_zones(is_active);

CREATE INDEX idx_seats_public_id ON seats(public_id);
CREATE INDEX idx_seats_zone_id ON seats(zone_id);
CREATE INDEX idx_seats_row_number ON seats(row_number);
CREATE INDEX idx_seats_seat_type ON seats(seat_type);
CREATE INDEX idx_seats_is_active ON seats(is_active);

-- Add updated_at triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_venues_updated_at BEFORE UPDATE ON venues FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_venue_layouts_updated_at BEFORE UPDATE ON venue_layouts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_seating_zones_updated_at BEFORE UPDATE ON seating_zones FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_seats_updated_at BEFORE UPDATE ON seats FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- 001_create_venue_tables.down.sql
DROP TRIGGER IF EXISTS update_seats_updated_at ON seats;
DROP TRIGGER IF EXISTS update_seating_zones_updated_at ON seating_zones;
DROP TRIGGER IF EXISTS update_venue_layouts_updated_at ON venue_layouts;
DROP TRIGGER IF EXISTS update_venues_updated_at ON venues;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS seating_zones;
DROP TABLE IF EXISTS venue_layouts;
DROP TABLE IF EXISTS venues;
```

#### **event-service Migrations**

```sql
-- 001_create_event_tables.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- events table
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    venue_id BIGINT NOT NULL, -- Will reference venues.id from venue-service
    layout_id BIGINT NOT NULL, -- Will reference venue_layouts.id from venue-service
    organizer_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('concert', 'sports', 'conference', 'theater', 'comedy', 'exhibition', 'workshop', 'other')),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    doors_open TIMESTAMP,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'cancelled', 'completed', 'postponed')),
    max_capacity INTEGER NOT NULL CHECK (max_capacity > 0),
    current_capacity INTEGER DEFAULT 0 CHECK (current_capacity >= 0),
    images JSONB DEFAULT '[]'::jsonb, -- Array of image URLs
    tags JSONB DEFAULT '[]'::jsonb, -- Array of tags
    metadata JSONB DEFAULT '{}'::jsonb, -- Additional event data
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT event_date_check CHECK (end_date > start_date),
    CONSTRAINT event_capacity_check CHECK (current_capacity <= max_capacity)
);

-- event_pricing table
CREATE TABLE event_pricing (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id BIGINT NOT NULL, -- Will reference seating_zones.id from venue-service
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) DEFAULT 'USD' CHECK (currency IN ('USD', 'EUR', 'GBP', 'VND')),
    pricing_type VARCHAR(50) NOT NULL CHECK (pricing_type IN ('fixed', 'dynamic', 'tiered')),
    pricing_rules JSONB DEFAULT '{}'::jsonb, -- Dynamic pricing rules
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, zone_id)
);

-- event_seat_availability table
CREATE TABLE event_seat_availability (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    seat_id BIGINT NOT NULL, -- Will reference seats.id from venue-service
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'reserved', 'sold', 'blocked', 'maintenance')),
    reserved_until TIMESTAMP,
    booking_id BIGINT, -- Will reference bookings.id from booking-service
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, seat_id)
);

-- event_schedules table (for recurring events)
CREATE TABLE event_schedules (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('one_time', 'daily', 'weekly', 'monthly')),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    recurrence_rule JSONB DEFAULT '{}'::jsonb, -- RRULE format for recurring events
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX idx_events_public_id ON events(public_id);
CREATE INDEX idx_events_venue_id ON events(venue_id);
CREATE INDEX idx_events_organizer_id ON events(organizer_id);
CREATE INDEX idx_events_start_date ON events(start_date);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_is_active ON events(is_active);
CREATE INDEX idx_events_city ON events(city);

CREATE INDEX idx_event_pricing_public_id ON event_pricing(public_id);
CREATE INDEX idx_event_pricing_event_id ON event_pricing(event_id);
CREATE INDEX idx_event_pricing_zone_id ON event_pricing(zone_id);
CREATE INDEX idx_event_pricing_is_active ON event_pricing(is_active);

CREATE INDEX idx_event_seat_availability_public_id ON event_seat_availability(public_id);
CREATE INDEX idx_event_seat_availability_event_id ON event_seat_availability(event_id);
CREATE INDEX idx_event_seat_availability_seat_id ON event_seat_availability(seat_id);
CREATE INDEX idx_event_seat_availability_status ON event_seat_availability(status);
CREATE INDEX idx_event_seat_availability_reserved_until ON event_seat_availability(reserved_until);
CREATE INDEX idx_event_seat_availability_booking_id ON event_seat_availability(booking_id);

CREATE INDEX idx_event_schedules_public_id ON event_schedules(public_id);
CREATE INDEX idx_event_schedules_event_id ON event_schedules(event_id);
CREATE INDEX idx_event_schedules_start_date ON event_schedules(start_date);
CREATE INDEX idx_event_schedules_is_active ON event_schedules(is_active);

-- Add updated_at triggers
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_event_pricing_updated_at BEFORE UPDATE ON event_pricing FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_event_seat_availability_updated_at BEFORE UPDATE ON event_seat_availability FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_event_schedules_updated_at BEFORE UPDATE ON event_schedules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- 001_create_event_tables.down.sql
DROP TRIGGER IF EXISTS update_event_schedules_updated_at ON event_schedules;
DROP TRIGGER IF EXISTS update_event_seat_availability_updated_at ON event_seat_availability;
DROP TRIGGER IF EXISTS update_event_pricing_updated_at ON event_pricing;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

DROP TABLE IF EXISTS event_schedules;
DROP TABLE IF EXISTS event_seat_availability;
DROP TABLE IF EXISTS event_pricing;
DROP TABLE IF EXISTS events;
```

#### **ticket-service Migrations**

```sql
-- 001_create_ticket_tables.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- tickets table
CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    booking_id BIGINT NOT NULL, -- Will reference bookings.id from booking-service
    event_id BIGINT NOT NULL, -- Will reference events.id from event-service
    seat_id BIGINT NOT NULL, -- Will reference seats.id from venue-service
    user_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    ticket_number VARCHAR(50) UNIQUE NOT NULL, -- Human-readable ticket number
    ticket_type VARCHAR(50) NOT NULL CHECK (ticket_type IN ('standard', 'vip', 'premium', 'student', 'senior', 'child')),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) DEFAULT 'USD' CHECK (currency IN ('USD', 'EUR', 'GBP', 'VND')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'cancelled', 'refunded', 'expired', 'used')),
    qr_code TEXT, -- QR code data
    barcode TEXT, -- Barcode data
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb, -- Additional ticket data
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ticket_validity_check CHECK (valid_until > valid_from)
);

-- booking_sessions table
CREATE TABLE booking_sessions (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    event_id BIGINT NOT NULL, -- Will reference events.id from event-service
    session_token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'expired', 'cancelled')),
    expires_at TIMESTAMP NOT NULL,
    seat_selections JSONB DEFAULT '[]'::jsonb, -- Selected seats data
    total_amount DECIMAL(10,2) CHECK (total_amount >= 0),
    currency VARCHAR(3) DEFAULT 'USD' CHECK (currency IN ('USD', 'EUR', 'GBP', 'VND')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- seat_reservations table
CREATE TABLE seat_reservations (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    session_id BIGINT NOT NULL REFERENCES booking_sessions(id) ON DELETE CASCADE,
    event_id BIGINT NOT NULL, -- Will reference events.id from event-service
    seat_id BIGINT NOT NULL, -- Will reference seats.id from venue-service
    status VARCHAR(20) DEFAULT 'reserved' CHECK (status IN ('reserved', 'confirmed', 'released', 'expired')),
    reserved_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, seat_id)
);

-- ticket_transfers table (for ticket transfers between users)
CREATE TABLE ticket_transfers (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    ticket_id BIGINT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    from_user_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    to_user_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    transfer_reason VARCHAR(100),
    transfer_fee DECIMAL(10,2) DEFAULT 0 CHECK (transfer_fee >= 0),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    approved_at TIMESTAMP,
    approved_by BIGINT, -- Will reference users.id from auth-service
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ticket_refunds table
CREATE TABLE ticket_refunds (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    ticket_id BIGINT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    refund_amount DECIMAL(10,2) NOT NULL CHECK (refund_amount >= 0),
    refund_reason VARCHAR(255),
    refund_method VARCHAR(50) CHECK (refund_method IN ('original_payment', 'credit', 'voucher')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'processed', 'rejected')),
    processed_at TIMESTAMP,
    processed_by BIGINT, -- Will reference users.id from auth-service
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX idx_tickets_public_id ON tickets(public_id);
CREATE INDEX idx_tickets_booking_id ON tickets(booking_id);
CREATE INDEX idx_tickets_event_id ON tickets(event_id);
CREATE INDEX idx_tickets_user_id ON tickets(user_id);
CREATE INDEX idx_tickets_ticket_number ON tickets(ticket_number);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_valid_until ON tickets(valid_until);

CREATE INDEX idx_booking_sessions_public_id ON booking_sessions(public_id);
CREATE INDEX idx_booking_sessions_user_id ON booking_sessions(user_id);
CREATE INDEX idx_booking_sessions_event_id ON booking_sessions(event_id);
CREATE INDEX idx_booking_sessions_session_token ON booking_sessions(session_token);
CREATE INDEX idx_booking_sessions_expires_at ON booking_sessions(expires_at);
CREATE INDEX idx_booking_sessions_status ON booking_sessions(status);

CREATE INDEX idx_seat_reservations_public_id ON seat_reservations(public_id);
CREATE INDEX idx_seat_reservations_session_id ON seat_reservations(session_id);
CREATE INDEX idx_seat_reservations_event_id ON seat_reservations(event_id);
CREATE INDEX idx_seat_reservations_seat_id ON seat_reservations(seat_id);
CREATE INDEX idx_seat_reservations_reserved_until ON seat_reservations(reserved_until);
CREATE INDEX idx_seat_reservations_status ON seat_reservations(status);

CREATE INDEX idx_ticket_transfers_public_id ON ticket_transfers(public_id);
CREATE INDEX idx_ticket_transfers_ticket_id ON ticket_transfers(ticket_id);
CREATE INDEX idx_ticket_transfers_from_user_id ON ticket_transfers(from_user_id);
CREATE INDEX idx_ticket_transfers_to_user_id ON ticket_transfers(to_user_id);
CREATE INDEX idx_ticket_transfers_status ON ticket_transfers(status);

CREATE INDEX idx_ticket_refunds_public_id ON ticket_refunds(public_id);
CREATE INDEX idx_ticket_refunds_ticket_id ON ticket_refunds(ticket_id);
CREATE INDEX idx_ticket_refunds_user_id ON ticket_refunds(user_id);
CREATE INDEX idx_ticket_refunds_status ON ticket_refunds(status);

-- Add updated_at triggers
CREATE TRIGGER update_tickets_updated_at BEFORE UPDATE ON tickets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_booking_sessions_updated_at BEFORE UPDATE ON booking_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_seat_reservations_updated_at BEFORE UPDATE ON seat_reservations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ticket_transfers_updated_at BEFORE UPDATE ON ticket_transfers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ticket_refunds_updated_at BEFORE UPDATE ON ticket_refunds FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- 001_create_ticket_tables.down.sql
DROP TRIGGER IF EXISTS update_ticket_refunds_updated_at ON ticket_refunds;
DROP TRIGGER IF EXISTS update_ticket_transfers_updated_at ON ticket_transfers;
DROP TRIGGER IF EXISTS update_seat_reservations_updated_at ON seat_reservations;
DROP TRIGGER IF EXISTS update_booking_sessions_updated_at ON booking_sessions;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;

DROP TABLE IF EXISTS ticket_refunds;
DROP TABLE IF EXISTS ticket_transfers;
DROP TABLE IF EXISTS seat_reservations;
DROP TABLE IF EXISTS booking_sessions;
DROP TABLE IF EXISTS tickets;
```

#### **Additional Migration Files**

##### **002_add_venue_analytics.up.sql**

```sql
-- Venue analytics table for tracking venue performance
CREATE TABLE venue_analytics (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    venue_id BIGINT NOT NULL, -- Will reference venues.id from venue-service
    event_id BIGINT NOT NULL, -- Will reference events.id from event-service
    total_tickets_sold INTEGER DEFAULT 0 CHECK (total_tickets_sold >= 0),
    total_revenue DECIMAL(12,2) DEFAULT 0 CHECK (total_revenue >= 0),
    average_ticket_price DECIMAL(10,2) DEFAULT 0 CHECK (average_ticket_price >= 0),
    occupancy_rate DECIMAL(5,2) DEFAULT 0 CHECK (occupancy_rate >= 0 AND occupancy_rate <= 100),
    event_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_venue_analytics_venue_id ON venue_analytics(venue_id);
CREATE INDEX idx_venue_analytics_event_id ON venue_analytics(event_id);
CREATE INDEX idx_venue_analytics_event_date ON venue_analytics(event_date);

CREATE TRIGGER update_venue_analytics_updated_at BEFORE UPDATE ON venue_analytics FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

##### **003_add_event_notifications.up.sql**

```sql
-- Event notification preferences
CREATE TABLE event_notifications (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    event_id BIGINT NOT NULL, -- Will reference events.id from event-service
    user_id BIGINT NOT NULL, -- Will reference users.id from auth-service
    notification_type VARCHAR(50) NOT NULL CHECK (notification_type IN ('reminder', 'update', 'cancellation', 'promotion')),
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('email', 'sms', 'push', 'in_app')),
    is_enabled BOOLEAN DEFAULT true,
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, user_id, notification_type, channel)
);

CREATE INDEX idx_event_notifications_event_id ON event_notifications(event_id);
CREATE INDEX idx_event_notifications_user_id ON event_notifications(user_id);
CREATE INDEX idx_event_notifications_scheduled_at ON event_notifications(scheduled_at);
CREATE INDEX idx_event_notifications_is_enabled ON event_notifications(is_enabled);

CREATE TRIGGER update_event_notifications_updated_at BEFORE UPDATE ON event_notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

#### **Migration Management Scripts**

##### **Makefile for Migration Management**

```makefile
# Makefile for migration management
.PHONY: migrate-up migrate-down migrate-create migrate-status migrate-force

# Database URLs (update these for your environment)
VENUE_DB_URL ?= postgres://user:password@localhost:5432/venue_db?sslmode=disable
EVENT_DB_URL ?= postgres://user:password@localhost:5432/event_db?sslmode=disable
TICKET_DB_URL ?= postgres://user:password@localhost:5432/ticket_db?sslmode=disable

# Run migrations up for all services
migrate-up-all:
	@echo "Running migrations for venue-service..."
	@cd venue-service && migrate -path migrations -database "$(VENUE_DB_URL)" up
	@echo "Running migrations for event-service..."
	@cd event-service && migrate -path migrations -database "$(EVENT_DB_URL)" up
	@echo "Running migrations for ticket-service..."
	@cd ticket-service && migrate -path migrations -database "$(TICKET_DB_URL)" up

# Rollback migrations for all services
migrate-down-all:
	@echo "Rolling back migrations for venue-service..."
	@cd venue-service && migrate -path migrations -database "$(VENUE_DB_URL)" down
	@echo "Rolling back migrations for event-service..."
	@cd event-service && migrate -path migrations -database "$(EVENT_DB_URL)" down
	@echo "Rolling back migrations for ticket-service..."
	@cd ticket-service && migrate -path migrations -database "$(TICKET_DB_URL)" down

# Create new migration for specific service
migrate-create:
	@if [ -z "$(service)" ]; then \
		echo "Usage: make migrate-create service=<service_name> name=<migration_name>"; \
		exit 1; \
	fi
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create service=<service_name> name=<migration_name>"; \
		exit 1; \
	fi
	@cd $(service) && migrate create -ext sql -dir migrations -seq $(name)

# Check migration status for all services
migrate-status-all:
	@echo "Migration status for venue-service:"
	@cd venue-service && migrate -path migrations -database "$(VENUE_DB_URL)" version
	@echo "Migration status for event-service:"
	@cd event-service && migrate -path migrations -database "$(EVENT_DB_URL)" version
	@echo "Migration status for ticket-service:"
	@cd ticket-service && migrate -path migrations -database "$(TICKET_DB_URL)" version

# Force migration version (use with caution)
migrate-force:
	@if [ -z "$(service)" ]; then \
		echo "Usage: make migrate-force service=<service_name> version=<version>"; \
		exit 1; \
	fi
	@if [ -z "$(version)" ]; then \
		echo "Usage: make migrate-force service=<service_name> version=<version>"; \
		exit 1; \
	fi
	@cd $(service) && migrate -path migrations -database "$(DB_URL)" force $(version)

# Reset database (drop all tables and run migrations)
reset-db:
	@echo "Resetting venue-service database..."
	@cd venue-service && migrate -path migrations -database "$(VENUE_DB_URL)" drop
	@cd venue-service && migrate -path migrations -database "$(VENUE_DB_URL)" up
	@echo "Resetting event-service database..."
	@cd event-service && migrate -path migrations -database "$(EVENT_DB_URL)" drop
	@cd event-service && migrate -path migrations -database "$(EVENT_DB_URL)" up
	@echo "Resetting ticket-service database..."
	@cd ticket-service && migrate -path migrations -database "$(TICKET_DB_URL)" drop
	@cd ticket-service && migrate -path migrations -database "$(TICKET_DB_URL)" up
```

##### **Go Migration Runner**

```go
// cmd/migrate/main.go
package main

import (
    "database/sql"
    "flag"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    var (
        dbURL    = flag.String("database", "", "Database URL")
        path     = flag.String("path", "migrations", "Path to migration files")
        action   = flag.String("action", "up", "Migration action: up, down, version, force")
        version  = flag.Int("version", 0, "Version for force action")
    )
    flag.Parse()

    if *dbURL == "" {
        log.Fatal("Database URL is required")
    }

    db, err := sql.Open("postgres", *dbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+*path,
        "postgres", driver)
    if err != nil {
        log.Fatal(err)
    }

    switch *action {
    case "up":
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
        log.Println("Migrations completed successfully")
    case "down":
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
        log.Println("Migrations rolled back successfully")
    case "version":
        version, dirty, err := m.Version()
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Current version: %d, Dirty: %t", version, dirty)
    case "force":
        if err := m.Force(*version); err != nil {
            log.Fatal(err)
        }
        log.Printf("Forced migration to version %d", *version)
    default:
        log.Fatal("Invalid action. Use: up, down, version, or force")
    }
}
```

#### **Database Connection Setup**

```go
// internal/database/connection.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

type Config struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

func NewConnection(config *Config) (*sqlx.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    log.Println("Database connection established successfully")
    return db, nil
}

// CloseConnection closes the database connection
func CloseConnection(db *sqlx.DB) error {
    if db != nil {
        return db.Close()
    }
    return nil
}

// HealthCheck checks if the database is healthy
func HealthCheck(db *sqlx.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return db.PingContext(ctx)
}
```

#### **Migration Status Tracking**

```go
// internal/database/migration.go
package database

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationStatus struct {
    Version int
    Dirty   bool
    Error   error
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("failed to create migration driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsPath,
        "postgres", driver)
    if err != nil {
        return fmt.Errorf("failed to create migration instance: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    log.Println("Migrations completed successfully")
    return nil
}

func GetMigrationStatus(db *sql.DB, migrationsPath string) (*MigrationStatus, error) {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to create migration driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsPath,
        "postgres", driver)
    if err != nil {
        return nil, fmt.Errorf("failed to create migration instance: %w", err)
    }

    version, dirty, err := m.Version()
    if err != nil {
        return &MigrationStatus{Error: err}, nil
    }

    return &MigrationStatus{
        Version: int(version),
        Dirty:   dirty,
    }, nil
}
```

### 1.2 Protobuf Definitions

- [ ] Tạo `shared-lib/protos/venue.proto` với services:
  - `VenueService` (CRUD venues)
  - `LayoutService` (CRUD layouts, zones, seats)
- [ ] Tạo `shared-lib/protos/event.proto` với services:
  - `EventService` (CRUD events)
  - `PricingService` (CRUD pricing rules)
- [ ] Tạo `shared-lib/protos/ticket.proto` với services:
  - `TicketService` (booking, availability)
  - `SeatService` (seat management)

### 1.3 Service Infrastructure

- [ ] **venue-service**: Setup Go server (Gin/Fiber), gRPC server, database connection
- [ ] **event-service**: Setup Go server (Gin/Fiber), gRPC server, database connection
- [ ] **ticket-service**: Setup Go server (Gin/Fiber), gRPC server, database connection
- [ ] Tất cả services: Setup Redis connection cho caching
- [ ] Tất cả services: Setup Prometheus metrics

---

## 🏢 Phase 2: Venue Service Implementation (Go)

### 2.1 Core Models & Repositories

- [x] Tạo `Venue` struct với validation tags
- [x] Tạo `VenueLayout` struct với JSON canvas config
- [x] Tạo `SeatingZone` struct với coordinates
- [x] Tạo `Seat` struct với seat properties
- [x] Tạo repositories cho tất cả models (SQLx/GORM)

### 2.2 Business Logic Services

- [x] `venueService.go`: CRUD operations cho venues
- [x] `layoutService.go`: Layout management, canvas operations
- [x] `zoneService.go`: Zone management, coordinate calculations
- [x] `seatService.go`: Seat management, bulk operations

### 2.3 Canvas Layout System

- [x] Implement canvas coordinate system
- [x] Tạo seat positioning algorithms
- [x] Implement zone boundary calculations
- [x] Tạo layout validation logic

### 2.4 gRPC Controllers

- [x] `venueController.go`: Handle venue CRUD requests
- [x] `layoutController.go`: Handle layout operations
- [x] `zoneController.go`: Handle zone management
- [x] `seatController.go`: Handle seat operations

### 2.5 gRPC Services (Internal Communication)

- [x] `VenueService`: CRUD operations cho venues
- [x] `LayoutService`: Layout management operations
- [x] `ZoneService`: Zone management operations
- [x] `SeatService`: Seat management operations

### 2.6 Caching & Performance

- [x] Implement Redis caching cho venue data
- [x] Implement layout caching
- [x] Setup cache invalidation strategies
- [x] Implement bulk seat operations với goroutines

---

## 🎭 Phase 3: Event Service Implementation (Go)

### 3.1 Core Models & Repositories

- [ ] Tạo `Event` struct với validation tags
- [ ] Tạo `EventPricing` struct với pricing rules
- [ ] Tạo `EventSeatAvailability` struct
- [ ] Tạo repositories cho tất cả models (SQLx/GORM)

### 3.2 Business Logic Services

- [ ] `eventService.go`: CRUD operations cho events
- [ ] `pricingService.go`: Dynamic pricing logic
- [ ] `availabilityService.go`: Seat availability management
- [ ] `eventValidationService.go`: Event validation rules

### 3.3 Pricing System

- [ ] Implement dynamic pricing algorithms
- [ ] Tạo pricing rule engine
- [ ] Implement discount calculations
- [ ] Tạo pricing validation logic

### 3.4 gRPC Controllers

- [ ] `eventController.go`: Handle event CRUD requests
- [ ] `pricingController.go`: Handle pricing operations
- [ ] `availabilityController.go`: Handle availability checks

### 3.5 gRPC Services (Internal Communication)

- [ ] `EventService`: CRUD operations cho events
- [ ] `PricingService`: Pricing operations
- [ ] `AvailabilityService`: Availability checks

### 3.6 Integration với Venue Service

- [ ] gRPC client để gọi venue service
- [ ] Validate venue existence khi tạo event
- [ ] Fetch venue layout cho event
- [ ] Cache venue data locally

---

## 🎫 Phase 4: Ticket Service Implementation (Go)

### 4.1 Core Models & Repositories

- [ ] Tạo `Ticket` struct với validation tags
- [ ] Tạo `BookingSession` struct với timeout logic
- [ ] Tạo `SeatReservation` struct
- [ ] Tạo repositories cho tất cả models (SQLx/GORM)

### 4.2 Business Logic Services

- [ ] `ticketService.go`: Ticket booking operations
- [ ] `bookingSessionService.go`: Session management
- [ ] `seatReservationService.go`: Seat reservation logic
- [ ] `availabilityService.go`: Real-time availability

### 4.3 Booking Flow Logic

- [ ] Implement seat selection algorithm
- [ ] Tạo booking session timeout với goroutines
- [ ] Implement seat reservation logic
- [ ] Tạo booking validation rules

### 4.4 gRPC Controllers

- [ ] `ticketController.go`: Handle ticket operations
- [ ] `bookingController.go`: Handle booking requests
- [ ] `availabilityController.go`: Handle availability queries

### 4.5 gRPC Services (Internal Communication)

- [ ] `TicketService`: Ticket operations
- [ ] `BookingService`: Booking operations
- [ ] `AvailabilityService`: Real-time availability

### 4.6 Integration với Event & Venue Services

- [ ] gRPC clients để gọi event và venue services
- [ ] Fetch event details và pricing
- [ ] Fetch venue layout và seat info
- [ ] Real-time availability updates

---

## 🔗 Phase 5: Gateway Integration

### 5.1 gRPC Clients

- [ ] Tạo venue service client
- [ ] Tạo event service client
- [ ] Tạo ticket service client
- [ ] Implement connection pooling

### 5.2 Authorization Integration

- [ ] Add venue management permissions
- [ ] Add event management permissions
- [ ] Add ticket booking permissions
- [ ] Implement role-based access control

### 5.3 API Routes

- [ ] `/api/venues/*` - Venue management endpoints
- [ ] `/api/events/*` - Event management endpoints
- [ ] `/api/tickets/*` - Ticket booking endpoints
- [ ] Implement proper error handling

### 5.4 Rate Limiting

- [ ] Add rate limits cho booking endpoints
- [ ] Add rate limits cho venue/event management
- [ ] Implement different limits cho different user roles

---

## 🧪 Phase 6: Testing & Quality Assurance

### 6.1 Unit Tests (Go)

- [ ] **venue-service**: Test tất cả services và controllers với `testing` package
- [ ] **event-service**: Test tất cả services và controllers với `testing` package
- [ ] **ticket-service**: Test tất cả services và controllers với `testing` package
- [ ] Test error handling và edge cases
- [ ] Use `testify` package cho assertions và mocking
- [ ] Test database operations với test database

### 6.2 Integration Tests

- [ ] Test gRPC communication giữa services
- [ ] Test booking flow end-to-end
- [ ] Test venue-event-ticket integration
- [ ] Test caching behavior
- [ ] Use `httptest` package cho HTTP testing

### 6.3 Performance Tests

- [ ] Load testing cho booking endpoints với `vegeta`
- [ ] Stress testing cho seat reservation
- [ ] Test Redis caching performance
- [ ] Test concurrent booking scenarios với goroutines
- [ ] Benchmark critical functions với `testing.B`

### 6.4 API Documentation

- [ ] Generate OpenAPI specs cho tất cả endpoints
- [ ] Create API documentation với examples
- [ ] Document error codes và responses
- [ ] Create integration guides
- [ ] Use `swaggo/swag` cho Go API documentation

---

## 🚀 Phase 7: Deployment & Monitoring

### 7.1 Docker Configuration (Go)

- [ ] Tạo Dockerfile cho venue-service (multi-stage build)
- [ ] Tạo Dockerfile cho event-service (multi-stage build)
- [ ] Tạo Dockerfile cho ticket-service (multi-stage build)
- [ ] Setup Go modules và dependency caching
- [ ] Optimize binary size với Alpine Linux

#### **Dockerfile Example (Multi-stage Build)**

```dockerfile
# Dockerfile for venue-service
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
```

#### **Docker Compose Example**

```yaml
# docker-compose.yml
version: "3.8"

services:
  venue-service:
    build:
      context: ./venue-service
      dockerfile: Dockerfile
    ports:
      - "8081:8080"
      - "50051:50051"
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/venue_db?sslmode=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - ENV=development
    depends_on:
      - postgres
      - redis
    networks:
      - booking-network
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--no-verbose",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  event-service:
    build:
      context: ./event-service
      dockerfile: Dockerfile
    ports:
      - "8082:8080"
      - "50052:50051"
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/event_db?sslmode=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - ENV=development
    depends_on:
      - postgres
      - redis
    networks:
      - booking-network

  ticket-service:
    build:
      context: ./ticket-service
      dockerfile: Dockerfile
    ports:
      - "8083:8080"
      - "50053:50051"
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/ticket_db?sslmode=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - ENV=development
    depends_on:
      - postgres
      - redis
    networks:
      - booking-network

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=booking_system
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - booking-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - booking-network

volumes:
  postgres_data:
  redis_data:

networks:
  booking-network:
    driver: bridge
```

### 7.2 Docker Compose

- [ ] Update `deploy/docker-compose.yml` với 3 services mới
- [ ] Setup service dependencies
- [ ] Configure networking
- [ ] Setup volume mounts
- [ ] Configure health checks cho Go services

### 7.3 Environment Configuration

- [ ] Tạo `.env.example` cho mỗi service
- [ ] Setup database connection strings
- [ ] Configure Redis connections
- [ ] Setup gRPC service addresses
- [ ] Use `viper` package cho configuration management

### 7.4 Monitoring & Logging

- [ ] Setup Prometheus metrics cho tất cả services với `prometheus/client_golang`
- [ ] Implement structured logging với `zap` hoặc `logrus`
- [ ] Setup health check endpoints
- [ ] Configure alerting rules
- [ ] Use `opentelemetry` cho distributed tracing

---

## 📊 Phase 8: Advanced Features

### 8.1 Real-time Updates

- [ ] Integrate với realtime-service
- [ ] Send seat availability updates
- [ ] Send booking confirmations
- [ ] Send event updates

### 8.2 Analytics Integration

- [ ] Send booking events đến analytics-service
- [ ] Track venue performance metrics
- [ ] Track event popularity
- [ ] Generate booking reports

### 8.3 Notification Integration

- [ ] Send booking confirmations
- [ ] Send event reminders
- [ ] Send venue updates
- [ ] Integrate với email-worker

---

## ✅ Completion Criteria

### Phase 1-2: Venue Service (Go)

- [x] Venue CRUD operations working với Go structs
- [x] Canvas layout system functional với JSONB
- [x] Seat management working với goroutines
- [x] gRPC endpoints responding với protobuf
- [x] gRPC services working cho internal communication
- [x] Caching implemented với Redis
- [x] Unit tests passing với `testing` package

### Phase 3: Event Service (Go)

- [ ] Event CRUD operations working với Go structs
- [ ] Pricing system functional với dynamic algorithms
- [ ] Integration với venue service working qua gRPC
- [ ] gRPC endpoints responding với protobuf
- [ ] gRPC services working cho internal communication
- [ ] Unit tests passing với `testing` package

### Phase 4: Ticket Service (Go)

- [ ] Ticket booking flow working với goroutines
- [ ] Seat reservation system functional với concurrency
- [ ] Integration với event/venue services working qua gRPC
- [ ] Booking session timeout working với goroutines
- [ ] Real-time availability updates working
- [ ] gRPC services working cho internal communication
- [ ] Unit tests passing với `testing` package

### Phase 5-6: Integration & Testing

- [ ] Gateway integration complete (Node.js) - REST API cho clients
- [ ] All Go tests passing với `go test`
- [ ] Performance benchmarks met với `testing.B`
- [ ] gRPC API documentation complete
- [ ] Integration tests passing với gRPC testing

### Phase 7-8: Deployment & Features

- [ ] Go services deployed và running với Docker
- [ ] Monitoring setup complete với Prometheus
- [ ] Real-time features working
- [ ] Analytics integration complete

---

## 🎯 Priority Order

1. **Phase 1**: Database schema và infrastructure (Foundation)
2. **Phase 2**: Venue service (Go) - Core venue management
3. **Phase 3**: Event service (Go) - Event management
4. **Phase 4**: Ticket service (Go) - Booking system
5. **Phase 5**: Gateway integration (Node.js) - API layer
6. **Phase 6**: Testing (Go testing packages) - Quality assurance
7. **Phase 7**: Deployment (Docker + Go) - Production readiness
8. **Phase 8**: Advanced features (Enhancement)

---

## 📝 Notes

- **Database**: Sử dụng PostgreSQL với JSONB fields cho canvas config
- **Caching**: Redis cho performance optimization
- **Communication**: gRPC cho inter-service communication
- **Authorization**: JWT + role-based permissions
- **Monitoring**: Prometheus metrics + structured logging (zap/logrus)
- **Testing**: Go testing package + testify + gRPC testing
- **Language**: Go cho venue/event/ticket services, Node.js cho gateway
- **Framework**: gRPC cho inter-service communication, REST API chỉ ở gateway
- **ORM**: SQLx/GORM cho database operations
- **Configuration**: Viper cho environment management
- **Architecture**: Pure gRPC microservices, gateway exposes REST API cho clients

Bắt đầu với Phase 1 và làm từng bước một cách có hệ thống! 🚀

Ghi chú: Phase 2 đã hoàn thành toàn bộ venue-service với models, repository, service, controller, caching, gRPC endpoints, unit test.

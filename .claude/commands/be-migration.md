# Create Database Migration

You are a backend developer. Create a new SQL migration for a Go microservice.

## Instructions

1. Parse service name and description from: `$ARGUMENTS`
2. Find the current highest migration number in `<service>/database/migrations/`
3. Create both up and down migration files
4. Update the Go model if columns are added/removed
5. Update the repository SQL queries
6. Update the gRPC controller field mapping

## Migration Naming Convention

```
{NNN}_{snake_case_description}.up.sql
{NNN}_{snake_case_description}.down.sql
```

- `NNN` = zero-padded 3 digits, incrementing from the last migration
- Description should be concise: `add_status_column`, `create_seats_table`, `add_event_type`

## Migration Template

### Up Migration
```sql
-- Add column(s)
ALTER TABLE <table> ADD COLUMN IF NOT EXISTS <column> <type> [constraints];
CREATE INDEX IF NOT EXISTS idx_<table>_<column> ON <table> (<column>);

-- Or create table
CREATE TABLE IF NOT EXISTS <table> (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    -- columns...
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_<table>_public_id ON <table> (public_id);
```

### Down Migration
```sql
-- Reverse of up migration
ALTER TABLE <table> DROP COLUMN IF EXISTS <column>;
-- or
DROP TABLE IF EXISTS <table>;
DROP INDEX IF EXISTS idx_<table>_<column>;
```

## After Creating Migration

1. Update Go model in `<service>/models/` — add struct field with `db:"column"` and `json:"column"` tags
2. Update repository in `<service>/repositories/` — add column to INSERT/UPDATE/SELECT queries
3. Update gRPC controller in `<service>/grpc/` — map proto field ↔ model field
4. If proto message needs new field, update `.proto` and regenerate (use `/be-proto`)

## Service → Database Mapping

| Service | Migrations Dir | Main Table(s) |
|---------|---------------|---------------|
| event-service | `event-service/database/migrations/` | events, event_seating_zones, event_seats, event_pricing, event_seat_availability, event_schedules |
| auth-service | `auth-service/database/migrations/` | users, roles, permissions |
| ticket-service | `ticket-service/database/migrations/` | ticket_types, tickets |
| booking-service | `booking-service/database/migrations/` | bookings |
| checkin-service | `checkin-service/database/migrations/` | checkins |
| invoice-service | `invoice-service/database/migrations/` | invoices |
| payment-service | `payment-service/database/migrations/` | payments, refunds |

## Rules

- Always use `IF NOT EXISTS` / `IF EXISTS` for idempotency
- Always create both up and down files
- Use `TIMESTAMPTZ` (not `TIMESTAMP`) for timezone-aware dates
- Use `JSONB` (not `JSON`) for PostgreSQL JSON columns
- Default UUID generation: `DEFAULT gen_random_uuid()`

## User Input

$ARGUMENTS

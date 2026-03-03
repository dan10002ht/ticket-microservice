# ✅ Phase 1 - Database Setup (Planned)

Documentation placeholder for booking-service database work.

## 🎯 Scope
- Primary tables: `bookings`, `booking_items`, `booking_actions`
- Supporting tables: saga audit, lock/watchdog states
- Flyway migration strategy (V1 bootstrap, V2 indexes, V3 partitions)
- Connection pools (Hikari + PgPool)

## 📌 Planned Tasks
- [ ] Design ERD aligning with ticket-service + payment-service references
- [ ] Create Flyway V1__create_bookings_table.sql & V2__create_booking_items_table.sql
- [ ] Add constraints/indexes for lookup (booking_uuid, session_id, user_id)
- [ ] Document retention policy + archive strategy
- [ ] Provide seed/test data helpers (optional)

## 🧪 Testing Strategy
- TODO: Add Testcontainers plan for integration tests
- TODO: Outline rollback scenarios for failed migrations

## 📁 File Map
```
docs/booking-service/02_DATABASE_SETUP_COMPLETE.md
booking-service/src/main/resources/db/migration/V1__create_bookings_table.sql
booking-service/src/main/resources/db/migration/V2__create_booking_items_table.sql
```

_Last updated: Planning stage (2024)_

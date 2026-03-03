# Services Documentation

Per-service detailed documentation for the ticket booking system.

---

## Services Overview

| Service | Language | gRPC Port | HTTP Port | Status | Docs |
|---------|----------|-----------|-----------|--------|------|
| Auth Service | Node.js | 50051 | - | 90% | [docs](./auth/) |
| User Service | Go | 50052 | - | 100% | [docs](./user/) |
| Event Service | Go | 50053 | - | 80% | [docs](./event/) |
| Ticket Service | Go | 50054 | - | 70% | [docs](./ticket/) |
| Payment Service | Java | 50056 | 8081 | 20% | [docs](./payment/) |
| Realtime Service | Go | 50057 | 3003 | 100% | [docs](./realtime/) |
| Booking Service | Java | 50058 | 8084 | 100% | [docs](./booking/) |
| Booking Worker | Go | 50059 | - | 100% | [docs](./booking-worker/) |
| Email Worker | Go | 50060 | 8080 | 85% | [docs](./email-worker/) |
| API Gateway | Node.js | - | 53000 | 70% | [docs](./gateway/) |

---

## Auth Service

**Language**: Node.js
**Port**: gRPC: 50051
**Status**: 90% Complete

### Features

- JWT authentication
- Role-based access control (RBAC)
- Session management with Redis
- User registration/login
- Password reset
- Email verification

### Documentation

- [CACHE_IMPLEMENTATION.md](./auth/CACHE_IMPLEMENTATION.md) - Redis caching
- [database-design.md](./auth/database-design.md) - Database schema
- [IMPLEMENTATION_CHECKLIST.md](./auth/IMPLEMENTATION_CHECKLIST.md) - Implementation checklist
- [INTEGRATION_FLOWS_README.md](./auth/INTEGRATION_FLOWS_README.md) - Integration flows
- [REGISTRATION_FLOWS_README.md](./auth/REGISTRATION_FLOWS_README.md) - Registration flows
- [REPOSITORY_PATTERN.md](./auth/REPOSITORY_PATTERN.md) - Repository pattern

---

## User Service

**Language**: Go 1.22+
**Port**: gRPC: 50052
**Status**: 100% Complete

### Features

- User profile management (CRUD)
- Address management
- User preferences (JSONB)

### Documentation

- [README.md](./user/README.md) - Complete documentation

---

## Event Service

**Language**: Go 1.22+
**Port**: gRPC: 50053
**Status**: 80% Complete

### Features

- Event management (CRUD)
- Venue management
- Seating layout management
- Event scheduling
- Availability tracking

### Documentation

- [EVENT_NEW.md](./event/EVENT_NEW.md) - Event model design
- [README_EVENT_MODEL.md](./event/README_EVENT_MODEL.md) - Event model reference

---

## Ticket Service

**Language**: Go 1.22+
**Port**: gRPC: 50054
**Status**: 70% Complete

### Features

- Ticket type management
- Seat reservations
- Ticket availability checking
- Ticket CRUD operations

### Documentation

- [DATABASE_SETUP.md](./ticket/DATABASE_SETUP.md) - Database setup guide
- [IMPLEMENTATION_STATUS.md](./ticket/IMPLEMENTATION_STATUS.md) - Implementation status

---

## Payment Service

**Language**: Java 17 + Spring Boot 3.2
**Ports**: gRPC: 50056, REST: 8081
**Status**: 20% Complete

### Features

- Stripe integration
- Idempotency guarantee
- Transaction logging
- Refund management
- Webhook handling
- Async payment flow with Kafka

### Documentation

- [README.md](./payment/README.md) - Documentation index
- [01_SETUP_COMPLETE.md](./payment/01_SETUP_COMPLETE.md) - Project setup
- [02_DATABASE_SETUP_COMPLETE.md](./payment/02_DATABASE_SETUP_COMPLETE.md) - Database setup
- [03_DATABASE_SCHEMA.md](./payment/03_DATABASE_SCHEMA.md) - Schema reference

### Design Document

- [PAYMENT_SERVICE.md](../architecture/PAYMENT_SERVICE.md) - Complete design document

---

## Realtime Service

**Language**: Go 1.22+ + Gorilla WebSocket
**Ports**: HTTP/WebSocket: 3003, gRPC: 50057
**Status**: 100% Complete

### Features

- WebSocket connections with JWT auth
- Real-time notifications (booking, payment, tickets)
- Room-based broadcasting
- Redis Pub/Sub for horizontal scaling
- gRPC API for internal services

### Documentation

- [README.md](./realtime/README.md) - Complete documentation

---

## Booking Service

**Language**: Java 17 + Spring Boot 3.2
**Ports**: REST: 8084, gRPC: 50058
**Status**: 100% Complete

### Features

- Saga Pattern orchestration
- Distributed locking (Redisson)
- Transactional Outbox pattern
- State machine for booking status
- Compensation handling (DLQ)
- Kafka event publishing

### Documentation

- [README.md](./booking/README.md) - Main documentation
- [SAGA_PATTERN_EXPLAINED.md](./booking/SAGA_PATTERN_EXPLAINED.md) - Saga pattern details
- [IMPLEMENTATION_SUMMARY.md](./booking/IMPLEMENTATION_SUMMARY.md) - Implementation summary
- [REDIS_KEY_SCHEMA.md](./booking/REDIS_KEY_SCHEMA.md) - Redis key schema

---

## Booking Worker

**Language**: Go 1.22+
**Port**: gRPC: 50059
**Status**: 100% Complete

### Features

- Redis queue processing
- Worker pool pattern
- Distributed booking queue per event
- Position tracking
- Timeout handling
- gRPC calls to Booking Service and Realtime Service

### Documentation

- [README.md](./booking-worker/README.md) - Main documentation
- [01_SETUP_COMPLETE.md](./booking-worker/01_SETUP_COMPLETE.md) - Setup guide
- [02_QUEUE_SETUP_COMPLETE.md](./booking-worker/02_QUEUE_SETUP_COMPLETE.md) - Queue setup
- [03_ARCHITECTURE.md](./booking-worker/03_ARCHITECTURE.md) - Architecture details

---

## Email Worker

**Language**: Go 1.22+
**Ports**: HTTP: 8080, gRPC: 50060
**Status**: 85% Complete

### Features

- Email queue processing via Kafka
- Template management
- Provider abstraction (SendGrid, SES, SMTP)
- Retry mechanism with exponential backoff
- Email tracking

### Documentation

- [API.md](./email-worker/API.md) - API documentation
- [folder-structure.md](./email-worker/folder-structure.md) - Folder structure
- [STEP.md](./email-worker/STEP.md) - Implementation steps

---

## API Gateway

**Language**: Node.js + Express
**Port**: HTTP: 53000
**Status**: 70% Complete

### Features

- Request routing to all microservices via gRPC
- JWT authentication middleware
- Role-based authorization
- Request validation
- Swagger/OpenAPI documentation at `/api/docs`

### Documentation

- [API_CHECKLIST.md](./gateway/API_CHECKLIST.md) - API implementation checklist
- [AUTHORIZATION_GUIDE.md](./gateway/AUTHORIZATION_GUIDE.md) - Authorization guide
- [DEVELOPMENT.md](./gateway/DEVELOPMENT.md) - Development guide
- [SWAGGER_README.md](./gateway/SWAGGER_README.md) - Swagger/OpenAPI docs

---

## Service Communication

### gRPC (internal)

| Service | Port | Proto |
|---------|------|-------|
| Auth Service | 50051 | auth.proto |
| User Service | 50052 | user.proto |
| Event Service | 50053 | event.proto |
| Ticket Service | 50054 | ticket.proto |
| Payment Service | 50056 | payment.proto |
| Realtime Service | 50057 | realtime.proto |
| Booking Service | 50058 | booking.proto |
| Booking Worker | 50059 | booking_worker.proto |
| Email Worker | 50060 | email.proto |

### REST/HTTP (external or actuator)

| Service | Port | Description |
|---------|------|-------------|
| API Gateway | 53000 | Public REST API |
| Realtime Service | 3003 | WebSocket + health |
| Booking Service | 8084 | Spring Actuator |
| Payment Service | 8081 | Spring Actuator |
| Email Worker | 8080 | HTTP health |

---

## Database Architecture

Each service owns its own database (Database-per-Service pattern):

| Service | Database Name | Engine |
|---------|---------------|--------|
| Auth Service | booking_system_auth | PostgreSQL (50432) |
| User Service | user_service | PostgreSQL (50433) |
| Event Service | booking_system_event | PostgreSQL (50433) |
| Ticket Service | booking_system_ticket | PostgreSQL (50433) |
| Booking Service | booking_service | PostgreSQL (50433) |
| Payment Service | payment_db | PostgreSQL (50433) |
| Email Worker | booking_system | PostgreSQL (50433) |

### Shared Infrastructure

- **PostgreSQL Auth** (50432) - Dedicated instance for auth service
- **PostgreSQL Main** (50433) - Shared instance, separate databases per service
- **Redis** (50379) - Caching, session storage, distributed locks, queues
- **Kafka** (50092) - Event streaming between services

---

## Development

### Start All Services

```bash
./scripts/dev-all.sh
```

### Health Check Dashboard

```bash
./scripts/health-check.sh          # One-shot status
./scripts/health-check.sh --watch  # Live dashboard
```

### View Service Logs

```bash
./scripts/dev-logs.sh --list          # List available logs
./scripts/dev-logs.sh --tail auth     # Tail auth service log
./scripts/dev-logs.sh --all           # Tail all services
```

---

## Related Documentation

- [Main Documentation Index](../README.md)
- [Architecture Documentation](../architecture/)
- [Service Connections](../architecture/SERVICE_CONNECTIONS.md)
- [Implementation Checklists](../checklists/)
- [Setup Guides](../guides/)

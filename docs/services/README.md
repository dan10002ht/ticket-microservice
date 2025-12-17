# Services Documentation

Per-service detailed documentation for the ticket booking system.

---

## Services Overview

| Service | Language | gRPC Port | Status | Documentation |
|---------|----------|-----------|--------|---------------|
| [Auth Service](#auth-service) | Node.js | 50051 | 90% | [Docs](./auth/) |
| [User Service](#user-service) | Go | 50052 | 100% | [Docs](../user-service/) |
| [Event Service](#event-service) | Node.js | 50054 | 80% | [Docs](./event/) |
| [Ticket Service](#ticket-service) | Node.js | 50055 | 70% | [Docs](./ticket/) |
| [Booking Service](#booking-service) | Java | 50056 | 100% | [Docs](../booking-service/) |
| [Realtime Service](#realtime-service) | Go | 50057 | 100% | [Docs](../realtime-service/) |
| [Payment Service](#payment-service) | Java | 50058 | 100% | [Docs](./payment-service/) |
| [API Gateway](#api-gateway) | Node.js | - | 70% | [Docs](./gateway/) |
| [Booking Worker](#booking-worker) | Go | - | 100% | [Docs](../booking-worker/) |
| [Email Worker](#email-worker) | Go | - | 85% | [Docs](./email-worker/) |

---

## Auth Service

**Language**: Node.js
**Status**: 90% Complete
**Port**: gRPC: 50051

### Features

- JWT authentication
- Role-based access control (RBAC)
- Session management
- User registration/login
- Password reset
- Email verification

### Documentation

- **[CACHE_IMPLEMENTATION.md](./auth/CACHE_IMPLEMENTATION.md)** - Redis caching
- **[database-design.md](./auth/database-design.md)** - Database schema
- **[IMPLEMENTATION_CHECKLIST.md](./auth/IMPLEMENTATION_CHECKLIST.md)** - Implementation checklist

---

## User Service

**Language**: Go 1.22+
**Status**: 100% Complete
**Port**: gRPC: 50052, Metrics: 9092

### Features

- User profile management (CRUD)
- Address management
- User preferences (JSONB)
- Legacy API compatibility

### Documentation

- **[README.md](../user-service/README.md)** - Complete documentation

---

## Event Service

**Language**: Node.js
**Status**: 80% Complete
**Port**: gRPC: 50054

### Features

- Event management (CRUD)
- Venue management
- Seating layout management
- Event scheduling
- Availability tracking

### Documentation

- **[EVENT_NEW.md](./event/EVENT_NEW.md)** - Event model design
- **[README_EVENT_MODEL.md](./event/README_EVENT_MODEL.md)** - Event model reference

---

## Ticket Service

**Language**: Node.js
**Status**: 70% Complete
**Port**: gRPC: 50055

### Features

- Ticket generation
- Seat reservations
- Ticket validation
- QR code generation
- Cancellation handling

### Documentation

- **[DATABASE_SETUP.md](./ticket/DATABASE_SETUP.md)** - Database setup guide
- **[IMPLEMENTATION_STATUS.md](./ticket/IMPLEMENTATION_STATUS.md)** - Implementation status

---

## Booking Service

**Language**: Java 17 + Spring Boot 3.2
**Status**: 100% Complete
**Ports**: REST: 8080, gRPC: 50056

### Features

- Saga Pattern orchestration
- Distributed locking (Redisson)
- Transactional Outbox pattern
- State machine for booking status
- Compensation handling (DLQ)
- Kafka event publishing

### Documentation

- **[README.md](../booking-service/README.md)** - Main documentation
- **[SAGA_PATTERN_EXPLAINED.md](../booking-service/SAGA_PATTERN_EXPLAINED.md)** - Saga pattern details
- **[IMPLEMENTATION_SUMMARY.md](../booking-service/IMPLEMENTATION_SUMMARY.md)** - Implementation summary
- **[REDIS_KEY_SCHEMA.md](../booking-service/REDIS_KEY_SCHEMA.md)** - Redis key schema

---

## Realtime Service

**Language**: Go 1.22+ + Gorilla WebSocket
**Status**: 100% Complete
**Ports**: HTTP/WebSocket: 3003, gRPC: 50057, Metrics: 9057

### Features

- WebSocket connections with JWT auth
- Real-time notifications (booking, payment, tickets)
- Room-based broadcasting
- Redis Pub/Sub for horizontal scaling
- gRPC API for internal services

### Documentation

- **[README.md](../realtime-service/README.md)** - Complete documentation

---

## Payment Service

**Language**: Java 17 + Spring Boot 3.2
**Status**: 100% Complete
**Ports**: REST: 8081, gRPC: 50058

### Features

- Stripe integration
- Idempotency guarantee
- Transaction logging
- Refund management
- Webhook handling
- Async payment flow with Kafka

### Documentation

- **[README.md](./payment-service/README.md)** - Documentation index
- **[01_SETUP_COMPLETE.md](./payment-service/01_SETUP_COMPLETE.md)** - Project setup
- **[02_DATABASE_SETUP_COMPLETE.md](./payment-service/02_DATABASE_SETUP_COMPLETE.md)** - Database setup
- **[03_DATABASE_SCHEMA.md](./payment-service/03_DATABASE_SCHEMA.md)** - Schema reference

### Design Document

- **[PAYMENT_SERVICE.md](../architecture/PAYMENT_SERVICE.md)** - Complete design document

---

## API Gateway

**Language**: Node.js + Express
**Status**: 70% Complete
**Port**: REST: 3000

### Features

- Request routing to all services
- Authentication middleware
- Rate limiting
- Request validation
- Response formatting
- Swagger/OpenAPI documentation

### Documentation

- **[API_CHECKLIST.md](./gateway/API_CHECKLIST.md)** - API implementation checklist
- **[AUTHORIZATION_GUIDE.md](./gateway/AUTHORIZATION_GUIDE.md)** - Authorization guide
- **[DEVELOPMENT.md](./gateway/DEVELOPMENT.md)** - Development guide
- **[SWAGGER_README.md](./gateway/SWAGGER_README.md)** - Swagger/OpenAPI docs

---

## Booking Worker

**Language**: Go 1.22+
**Status**: 100% Complete
**Port**: Metrics: 9093

### Features

- Redis queue processing
- Worker pool pattern
- Distributed booking queue per event
- Position tracking
- Timeout handling
- gRPC calls to Booking Service and Realtime Service

### Documentation

- **[README.md](../booking-worker/README.md)** - Main documentation
- **[01_SETUP_COMPLETE.md](../booking-worker/01_SETUP_COMPLETE.md)** - Setup guide
- **[02_QUEUE_SETUP_COMPLETE.md](../booking-worker/02_QUEUE_SETUP_COMPLETE.md)** - Queue setup
- **[03_ARCHITECTURE.md](../booking-worker/03_ARCHITECTURE.md)** - Architecture details

---

## Email Worker

**Language**: Go 1.22+
**Status**: 85% Complete
**Ports**: HTTP: 8080, Metrics: 9094

### Features

- Email queue processing
- Template management
- Provider abstraction (SendGrid, SES, SMTP)
- Retry mechanism with exponential backoff
- Email tracking
- Scheduled jobs

### Documentation

- **[API.md](./email-worker/API.md)** - API documentation
- **[folder-structure.md](./email-worker/folder-structure.md)** - Folder structure
- **[STEP.md](./email-worker/STEP.md)** - Implementation steps

---

## Service Communication

### gRPC Services

All internal service-to-service communication uses gRPC:

| Service | Port | Description |
|---------|------|-------------|
| Auth Service | 50051 | Authentication |
| User Service | 50052 | User profiles |
| Event Service | 50054 | Event management |
| Ticket Service | 50055 | Ticket management |
| Booking Service | 50056 | Booking orchestration |
| Realtime Service | 50057 | Real-time notifications |
| Payment Service | 50058 | Payment processing |

### REST/HTTP APIs

| Service | Port | Description |
|---------|------|-------------|
| API Gateway | 3000 | External API |
| Realtime Service | 3003 | WebSocket |
| Booking Service | 8080 | REST + Actuator |
| Payment Service | 8081 | REST + Actuator |

---

## Database Architecture

### Per-Service Databases

| Service | Database Name |
|---------|---------------|
| Auth Service | booking_system_auth |
| User Service | user_service |
| Event Service | booking_system_event |
| Ticket Service | booking_system_ticket |
| Booking Service | booking_service |
| Payment Service | payment_db |
| Email Worker | booking_system |

### Infrastructure

- PostgreSQL with Master-Slave replication
- Redis for caching and queues
- Kafka for event streaming

---

## Development

### Starting Individual Services

```bash
# Auth Service
cd auth-service && yarn dev:local

# User Service
cd user-service && go run main.go

# Event Service
cd event-service && yarn dev:local

# Ticket Service
cd ticket-service && yarn dev:local

# Booking Service
cd booking-service && mvn spring-boot:run -Dspring-boot.run.profiles=dev

# Booking Worker
cd booking-worker && go run main.go

# Realtime Service
cd realtime-service && go run main.go

# Payment Service
cd payment-service && mvn spring-boot:run -Dspring-boot.run.profiles=dev

# Gateway
cd gateway && yarn dev:local

# Email Worker
cd email-worker && go run main.go
```

---

## Related Documentation

- [Main Documentation Index](../README.md)
- [Architecture Documentation](../architecture/)
- [Service Connections](../architecture/SERVICE_CONNECTIONS.md)
- [Implementation Checklists](../checklists/)
- [Setup Guides](../guides/)

---

**Last Updated**: December 2024
**Total Services**: 10
**Average Completion**: 88%

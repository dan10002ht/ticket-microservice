# 🎯 Services Documentation

Per-service detailed documentation for the ticket booking system.

---

## 📦 Services Overview

| Service                              | Language | Status | Documentation              |
| ------------------------------------ | -------- | ------ | -------------------------- |
| [Payment Service](#-payment-service) | Java     | 20%    | [Docs](./payment-service/) |
| [Auth Service](#-auth-service)       | Node.js  | 90%    | [Docs](./auth/)            |
| [Event Service](#-event-service)     | Go       | 80%    | [Docs](./event/)           |
| [Ticket Service](#-ticket-service)   | Go       | 70%    | [Docs](./ticket/)          |
| [API Gateway](#-api-gateway)         | Node.js  | 70%    | [Docs](./gateway/)         |
| [Email Worker](#-email-worker)       | Go       | 85%    | [Docs](./email-worker/)    |

---

## 💳 Payment Service

**Language**: Java 17 + Spring Boot 3.2.0  
**Status**: ✅ Phase 1 - Database Setup Complete (20%)  
**Ports**: REST: 8081, gRPC: 50056

### Features

- Multi-gateway support (Stripe, PayPal, VNPay, Momo)
- Idempotency guarantee
- Transaction logging
- Refund management
- Webhook handling

### Documentation

- **[README.md](./payment-service/README.md)** - Documentation index
- **[01_SETUP_COMPLETE.md](./payment-service/01_SETUP_COMPLETE.md)** - Project setup
- **[02_DATABASE_SETUP_COMPLETE.md](./payment-service/02_DATABASE_SETUP_COMPLETE.md)** - Database setup
- **[03_DATABASE_SCHEMA.md](./payment-service/03_DATABASE_SCHEMA.md)** - Schema reference

### Design Document

- **[PAYMENT_SERVICE.md](../architecture/PAYMENT_SERVICE.md)** - Complete design document

---

## 🔐 Auth Service

**Language**: Node.js  
**Status**: ✅ Core Implementation Complete (90%)  
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
- **[INTEGRATION_FLOWS_README.md](./auth/INTEGRATION_FLOWS_README.md)** - Integration flows
- **[INTEGRATION_TODO.md](./auth/INTEGRATION_TODO.md)** - Integration tasks
- **[PGPOOL_COMPATIBILITY.md](./auth/PGPOOL_COMPATIBILITY.md)** - PgPool integration
- **[REGISTRATION_FLOWS_README.md](./auth/REGISTRATION_FLOWS_README.md)** - Registration flows
- **[REPOSITORY_PATTERN.md](./auth/REPOSITORY_PATTERN.md)** - Repository pattern

---

## 🎪 Event Service

**Language**: Go  
**Status**: ✅ Core Features Complete (80%)  
**Port**: gRPC: 50053

### Features

- Event management (CRUD)
- Venue management
- Seating layout management
- Zone and seat management
- Event scheduling
- Availability tracking

### Documentation

- **[EVENT_NEW.md](./event/EVENT_NEW.md)** - Event model design
- **[README_EVENT_MODEL.md](./event/README_EVENT_MODEL.md)** - Event model reference

---

## 🎫 Ticket Service

**Language**: Go  
**Status**: ✅ Models, Repos, Services Complete (70%)  
**Port**: gRPC: 50057

### Features

- Ticket generation
- Booking sessions
- Seat reservations
- Ticket validation
- QR code generation
- Cancellation & refunds

### Documentation

- **[DATABASE_SETUP.md](./ticket/DATABASE_SETUP.md)** - Database setup guide
- **[IMPLEMENTATION_STATUS.md](./ticket/IMPLEMENTATION_STATUS.md)** - Implementation status

---

## 🌐 API Gateway

**Language**: Node.js + Express  
**Status**: ✅ Core Routes Complete (70%)  
**Port**: REST: 53000

### Features

- Request routing
- Authentication middleware
- Rate limiting
- Request validation
- Response formatting
- Error handling
- Swagger/OpenAPI documentation

### Documentation

- **[API_CHECKLIST.md](./gateway/API_CHECKLIST.md)** - API implementation checklist
- **[API_EVENT_CREATION_CHECKLIST.md](./gateway/API_EVENT_CREATION_CHECKLIST.md)** - Event creation API
- **[AUTHORIZATION_GUIDE.md](./gateway/AUTHORIZATION_GUIDE.md)** - Authorization guide
- **[DEVELOPMENT.md](./gateway/DEVELOPMENT.md)** - Development guide
- **[EVENT_CREATION_CHECKLIST.md](./gateway/EVENT_CREATION_CHECKLIST.md)** - Event creation flow
- **[EVENT_DISPLAY_CHECKLIST.md](./gateway/EVENT_DISPLAY_CHECKLIST.md)** - Event display flow
- **[QUICK_CHECKLIST.md](./gateway/QUICK_CHECKLIST.md)** - Quick reference
- **[SWAGGER_README.md](./gateway/SWAGGER_README.md)** - Swagger/OpenAPI docs
- **[YARN_MIGRATION.md](./gateway/YARN_MIGRATION.md)** - Yarn migration guide

---

## 📧 Email Worker

**Language**: Go  
**Status**: ✅ Core Implementation Complete (85%)  
**Ports**: HTTP: 8080, gRPC: 50060

### Features

- Email queue processing
- Template management
- Provider abstraction (SendGrid, etc.)
- Retry mechanism
- Email tracking
- Bulk email support

### Documentation

- **[API.md](./email-worker/API.md)** - API documentation
- **[folder-structure.md](./email-worker/folder-structure.md)** - Folder structure
- **[STEP.md](./email-worker/STEP.md)** - Implementation steps

---

## 🔗 Service Communication

### gRPC Services

All internal service-to-service communication uses gRPC:

- Auth Service: `localhost:50051`
- Event Service: `localhost:50053`
- Ticket Service: `localhost:50057`
- Payment Service: `localhost:50056`
- Email Worker: `localhost:50060`

### REST APIs

External-facing REST APIs:

- API Gateway: `http://localhost:53000`
- Payment Service: `http://localhost:8081`
- Email Worker: `http://localhost:8080`

---

## 📊 Database Architecture

### Per-Service Databases

Each service has its own database:

- `booking_system_auth` - Auth Service
- `booking_system_event` - Event Service
- `booking_system_ticket` - Ticket Service
- `payment_db` - Payment Service
- `booking_system` - Email Worker

### Replication

- PostgreSQL Master-Slave replication
- PgPool-II for connection pooling
- Read replicas for heavy read operations

---

## 🔧 Development

### Starting Individual Services

```bash
# Auth Service
cd auth-service
yarn dev:local

# Event Service
cd event-service
go run main.go

# Ticket Service
cd ticket-service
go run main.go

# Payment Service
cd payment-service
mvn spring-boot:run -Dspring-boot.run.profiles=dev

# Gateway
cd gateway
yarn dev:local

# Email Worker
cd email-worker
go run main.go
```

### Starting All Services

```bash
# From project root
cd scripts
./dev-all.sh
```

---

## 📖 Related Documentation

- [Main Documentation Index](../README.md)
- [Architecture Documentation](../architecture/)
- [Implementation Checklists](../checklists/)
- [Setup Guides](../guides/)

---

**Last Updated**: 2024  
**Total Services**: 6  
**Average Completion**: 70%



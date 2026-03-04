# Architecture Documentation

High-level architecture design and system documentation.

---

## Documents

### Core Architecture

#### [AI_README.md](./AI_README.md)

- AI agent guidance for understanding the project
- Project structure conventions
- Service discovery
- API discovery
- Integration points

#### [MICROSERVICE_BEST_PRACTICES.md](./MICROSERVICE_BEST_PRACTICES.md)

- Microservices design patterns
- Best practices for distributed systems
- Communication patterns
- Error handling strategies
- Testing approaches

#### [SERVICE_CONNECTIONS.md](./SERVICE_CONNECTIONS.md)

- Service dependency map
- Communication protocols (gRPC, REST, WebSocket)
- Integration flows
- Complete architecture diagram

### Service Design Documents

#### [PAYMENT_SERVICE.md](./PAYMENT_SERVICE.md)

- Complete Payment Service design document
- Stripe integration
- Database schema design
- Async payment flow with Kafka

---

## System Architecture Overview

```
                                    ┌─────────────────────┐
                                    │      Frontend       │
                                    │   (React/Next.js)   │
                                    └─────────┬───────────┘
                                              │ HTTP/WebSocket
                                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            API Gateway (Node.js)                             │
│                              Port: 53000                                     │
└──┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬───────────────┘
   │      │      │      │      │      │      │      │      │
   │gRPC  │gRPC  │gRPC  │gRPC  │gRPC  │gRPC  │gRPC  │gRPC  │gRPC
   ▼      ▼      ▼      ▼      ▼      ▼      ▼      ▼      ▼
┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐
│ Auth ││ User ││Event ││Ticket││Book- ││Pay-  ││Check-││Invoi-││Real- │
│:50051││:50052││:50053││:50054││ing   ││ment  ││in    ││ce    ││time  │
│Node  ││ Go   ││ Go   ││ Go   ││:50058││:50062││:50059││:50060││:50057│
│      ││      ││      ││      ││ Java ││ Java ││ Go   ││ Java ││ Go   │
└──────┘└──────┘└──────┘└──────┘└──┬───┘└──┬───┘└──────┘└──────┘└──┬───┘
                                   │       │                        │
                            ┌──────┴───────┘                        │ WebSocket
                            │ gRPC calls (Saga)                     ▼
                            ▼                                ┌─────────────┐
                       ┌─────────────┐                       │  Frontend   │
                       │   Booking   │──────────────────────▶│ (Real-time) │
                       │   Worker    │ gRPC (realtime)       └─────────────┘
                       │  :50056 Go  │
                       └──────┬──────┘
                              │ gRPC
                              ▼
                       ┌─────────────┐
                       │   Email     │
                       │   Worker    │
                       │  :50061 Go  │
                       └─────────────┘

Infrastructure:
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ PostgreSQL  │  │    Redis    │  │    Kafka    │
│ (Master-    │  │ (Cache,     │  │ (Event      │
│  Slave)     │  │  Queue)     │  │  Stream)    │
└─────────────┘  └─────────────┘  └─────────────┘
```

---

## Service Inventory

| Service | Language | gRPC Port | HTTP Port | Description |
|---------|----------|-----------|-----------|-------------|
| Gateway | Node.js | - | 53000 | API Gateway, Swagger docs |
| Auth Service | Node.js | 50051 | - | Authentication & authorization |
| User Service | Go | 50052 | - | User profiles & addresses |
| Event Service | Go | 50053 | - | Event & venue management |
| Ticket Service | Go | 50054 | - | Ticket inventory & types |
| Booking Worker | Go | 50056 | - | Queue-based booking processing |
| Realtime Service | Go | 50057 | 3003 | WebSocket notifications |
| Booking Service | Java | 50058 | 8084 | Booking orchestration (Saga) |
| Checkin Service | Go | 50059 | - | Event check-in & QR scanning |
| Invoice Service | Java | 50060 | 8083 | Invoice generation & PDF |
| Email Worker | Go | 50061 | - | Async email delivery |
| Payment Service | Java | 50062 | 8080 | Payment processing (Stripe) |

---

## Key Integration Patterns

### Synchronous Communication (gRPC)

- Gateway ↔ All Services
- Booking Service ↔ Ticket Service (Saga)
- Booking Service ↔ Payment Service (Saga)
- Booking Worker ↔ Booking Service
- Booking Worker ↔ Realtime Service
- Auth Service ↔ Email Worker

### Asynchronous Communication (Kafka)

- Payment events (booking-service → outbox → kafka)
- Booking confirmed events
- Email notifications
- Audit logs

### Real-time Communication (WebSocket)

- Realtime Service ↔ Frontend
- Redis Pub/Sub for horizontal scaling

### Caching & Queue (Redis)

- Session storage
- Booking queue (per event)
- Distributed locks
- API rate limiting
- Idempotency keys

### Database (PostgreSQL)

- Master-Slave replication
- Per-service database isolation

---

## Technology Stack

| Layer | Technologies |
|-------|--------------|
| **API Gateway** | Node.js, Express, JWT |
| **Services** | Node.js, Go 1.22+, Java 17 + Spring Boot 3.2 |
| **Communication** | gRPC, REST, WebSocket, Kafka |
| **Real-time** | Gorilla WebSocket, Redis Pub/Sub |
| **Databases** | PostgreSQL (Master-Slave) |
| **Caching** | Redis |
| **Monitoring** | Prometheus, Grafana |
| **Logging** | Zap (Go), Logback (Java), Winston (Node.js) |
| **Containers** | Docker, Docker Compose |

---

## Design Patterns Used

| Pattern | Used In | Purpose |
|---------|---------|---------|
| **Saga Pattern** | Booking Service | Distributed transaction orchestration |
| **Transactional Outbox** | Booking Service | Reliable event publishing |
| **Worker Pool** | Booking Worker, Email Worker | Concurrent processing |
| **Circuit Breaker** | All gRPC clients | Fault tolerance |
| **Distributed Lock** | Booking Service | Prevent concurrent modifications |
| **State Machine** | Booking Service | Booking status transitions |
| **Hub Pattern** | Realtime Service | WebSocket connection management |

---

## Design Principles

1. **Microservices First**: Each service is independently deployable
2. **Polyglot Architecture**: Use best language for each service
3. **Database Per Service**: No shared databases
4. **API Gateway Pattern**: Single entry point for clients
5. **Event-Driven**: Async communication for loosely coupled services
6. **Saga Pattern**: Distributed transactions with compensation
7. **Observability**: Comprehensive logging, metrics, tracing
8. **Resilience**: Circuit breakers, retries, fallbacks
9. **Security**: JWT auth, RBAC, encryption at rest/transit

---

## Related Documentation

- [Main Documentation Index](../README.md)
- [Service Documentation](../services/)
- [Implementation Checklists](../checklists/)
- [Setup Guides](../guides/)
- [Booking Service Saga](../booking-service/SAGA_PATTERN_EXPLAINED.md)
- [Realtime Service](../realtime-service/README.md)
- [User Service](../user-service/README.md)

---

**Last Updated**: March 2026
**Architecture Version**: 3.0

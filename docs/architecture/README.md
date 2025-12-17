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
│                              Port: 3000                                      │
└────┬──────────┬──────────┬──────────┬──────────┬──────────┬────────────────┘
     │          │          │          │          │          │
     │ gRPC     │ gRPC     │ gRPC     │ gRPC     │ gRPC     │ gRPC
     ▼          ▼          ▼          ▼          ▼          ▼
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│  Auth   │ │  Event  │ │ Ticket  │ │ Booking │ │ Payment │ │  User   │
│ Service │ │ Service │ │ Service │ │ Service │ │ Service │ │ Service │
│ :50051  │ │ :50054  │ │ :50055  │ │ :50056  │ │ :50058  │ │ :50052  │
│ Node.js │ │ Node.js │ │ Node.js │ │  Java   │ │  Java   │ │   Go    │
└─────────┘ └─────────┘ └────┬────┘ └────┬────┘ └────┬────┘ └─────────┘
                             │           │           │
                             │    ┌──────┴───────────┘
                             │    │ gRPC calls (Saga)
                             ▼    ▼
                        ┌─────────────┐
                        │   Booking   │      ┌─────────────┐
                        │   Worker    │─────▶│  Realtime   │
                        │    (Go)     │ gRPC │  Service    │
                        └──────┬──────┘      │ :50057/:3003│
                               │             │    (Go)     │
                               │             └──────┬──────┘
                               │                    │ WebSocket
                               │                    ▼
                               │             ┌─────────────┐
                               │             │  Frontend   │
                               │             │ (Real-time) │
                               │             └─────────────┘
                               ▼
                        ┌─────────────┐
                        │   Email     │
                        │   Worker    │
                        │    (Go)     │
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
| Gateway | Node.js | - | 3000 | API Gateway |
| Auth Service | Node.js | 50051 | - | Authentication |
| User Service | Go | 50052 | - | User profiles |
| Event Service | Node.js | 50054 | - | Event management |
| Ticket Service | Node.js | 50055 | - | Ticket/seat management |
| Booking Service | Java | 50056 | 8080 | Booking orchestration |
| Realtime Service | Go | 50057 | 3003 | WebSocket notifications |
| Payment Service | Java | 50058 | 8081 | Payment processing |
| Booking Worker | Go | - | - | Queue processing |
| Email Worker | Go | - | - | Email delivery |

---

## Key Integration Patterns

### Synchronous Communication (gRPC)

- Gateway ↔ All Services
- Booking Service ↔ Ticket Service (Saga)
- Booking Service ↔ Payment Service (Saga)
- Booking Worker ↔ Booking Service
- Booking Worker ↔ Realtime Service

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

**Last Updated**: December 2024
**Architecture Version**: 2.0

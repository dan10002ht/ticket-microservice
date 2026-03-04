# Ticket Booking System - Documentation

Central documentation hub for all microservices in the ticket booking system.

---

## Documentation Structure

```
docs/
├── README.md                              # This file - documentation index
├── CONTRIBUTING.md                        # Contributing guidelines
├── architecture/                          # Architecture & design
│   ├── README.md                          # Architecture overview
│   ├── SERVICE_CONNECTIONS.md             # Service integration map
│   ├── MICROSERVICE_BEST_PRACTICES.md     # Best practices guide
│   ├── PAYMENT_SERVICE.md                 # Payment service design doc
│   ├── AI_README.md                       # AI agent guidance
│   └── ORGANIZATION_SUMMARY.md            # Organization model design
├── checklists/                            # Implementation checklists
│   ├── AUTHORIZATION_CHECKLIST.md
│   ├── LOOSE_COUPLING_CHECKLIST.md
│   ├── PHASE1_COMPLETION_SUMMARY.md
│   └── VENUE_EVENT_TICKET_CHECKLIST.md
├── guides/                                # Setup & configuration guides
│   ├── development-setup.md               # Dev environment setup
│   ├── go-service-structure.md            # Go service boilerplate
│   ├── MASTER_SLAVE_SETUP.md              # PostgreSQL replication
│   ├── METRICS_SETUP.md                   # Monitoring setup
│   └── README_PORTS.md                    # Port configuration
└── services/                              # Per-service documentation
    ├── README.md                          # Service index & overview
    ├── auth/                              # Auth Service (Node.js)
    ├── user/                              # User Service (Go)
    ├── event/                             # Event Service (Go)
    ├── ticket/                            # Ticket Service (Go)
    ├── payment/                           # Payment Service (Java)
    ├── realtime/                          # Realtime Service (Go)
    ├── booking/                           # Booking Service (Java)
    ├── booking-worker/                    # Booking Worker (Go)
    ├── email-worker/                      # Email Worker (Go)
    └── gateway/                           # API Gateway (Node.js)
```

---

## Architecture Overview

```
                    ┌─────────────────────────────────┐
                    │     API Gateway (Node.js)        │
                    │         HTTP :53000               │
                    └──────────────┬──────────────────┘
                                   │ gRPC
     ┌──────────┬──────────┬───────┼───────┬──────────┬──────────┐
     ▼          ▼          ▼       ▼       ▼          ▼          ▼
┌─────────┐┌─────────┐┌─────────┐┌──────┐┌─────────┐┌─────────┐┌─────────┐
│  Auth   ││  User   ││  Event  ││Ticket││ Booking ││ Payment ││Realtime │
│ :50051  ││ :50052  ││ :50053  ││:50054││ :50058  ││ :50062  ││ :50057  │
│ Node.js ││   Go    ││   Go    ││  Go  ││  Java   ││  Java   ││   Go    │
└────┬────┘└────┬────┘└────┬────┘└──┬───┘└────┬────┘└────┬────┘└────┬────┘
     │          │          │        │         │          │          │
     └──────────┴──────────┴────────┴─────────┴──────────┴──────────┘
                                    │
          ┌────────────────┬────────┼────────┬────────────────┐
          ▼                ▼        ▼        ▼                ▼
   ┌────────────┐  ┌────────────┐┌──────┐┌───────┐  ┌──────────────┐
   │ PostgreSQL │  │ PostgreSQL ││Redis ││ Kafka │  │Booking Worker│
   │   Auth     │  │    Main    ││:50379││:50092 │  │   Go :50056  │
   │  :50432    │  │   :50433   │└──────┘└───────┘  └──────────────┘
   └────────────┘  └────────────┘                   ┌──────────────┐
                                                    │ Email Worker │
                                                    │  Go :50061   │
                                                    └──────────────┘
```

---

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Node.js 18+
- Go 1.22+
- Java 17+ (for Booking & Payment Services)
- Maven 3.8+ (for Booking & Payment Services)

### Start Development Environment

```bash
# Start infrastructure + all services
./scripts/dev-all.sh

# Check health status
./scripts/health-check.sh

# Live health dashboard
./scripts/health-check.sh --watch

# View logs per service
./scripts/dev-logs.sh --tail auth
```

See [Development Setup Guide](./guides/development-setup.md) for detailed instructions.

### Access Points

| Service | URL |
|---------|-----|
| API Gateway (REST) | http://localhost:53000 |
| API Documentation (Swagger) | http://localhost:53000/api/docs |
| Booking Service (Actuator) | http://localhost:8084/actuator/health |
| Payment Service (Actuator) | http://localhost:8080/actuator/health |
| Realtime (WebSocket) | ws://localhost:3003 |

---

## Services

| Service | Language | gRPC | HTTP | Status |
|---------|----------|------|------|--------|
| Auth Service | Node.js | 50051 | - | 90% |
| User Service | Go | 50052 | - | 100% |
| Event Service | Go | 50053 | - | 80% |
| Ticket Service | Go | 50054 | - | 70% |
| Booking Service | Java | 50058 | 8084 | 100% |
| Payment Service | Java | 50062 | 8080 | 20% |
| Realtime Service | Go | 50057 | 3003 | 100% |
| Checkin Service | Go | 50059 | - | 100% |
| Booking Worker | Go | 50056 | - | 100% |
| Email Worker | Go | 50061 | - | 85% |
| Invoice Service | Java | 50060 | 8083 | 70% |
| API Gateway | Node.js | - | 53000 | 70% |

See [Services Documentation](./services/README.md) for detailed per-service info.

---

## Architecture Documentation

- [Service Connections](./architecture/SERVICE_CONNECTIONS.md) - How services communicate
- [Microservice Best Practices](./architecture/MICROSERVICE_BEST_PRACTICES.md) - Design guidelines
- [Payment Service Design](./architecture/PAYMENT_SERVICE.md) - Payment architecture decisions
- [Organization Model](./architecture/ORGANIZATION_SUMMARY.md) - Organization entity design
- [AI Agent Guide](./architecture/AI_README.md) - AI assistant context

---

## Guides

- [Development Setup](./guides/development-setup.md) - Dev environment, scripts, health checks
- [PostgreSQL Replication](./guides/MASTER_SLAVE_SETUP.md) - Master-slave setup
- [Monitoring Setup](./guides/METRICS_SETUP.md) - Prometheus & Grafana
- [Port Configuration](./guides/README_PORTS.md) - Port allocation reference
- [Go Service Structure](./guides/go-service-structure.md) - Go service boilerplate

---

## Implementation Checklists

- [Authorization Checklist](./checklists/AUTHORIZATION_CHECKLIST.md) - RBAC implementation
- [Loose Coupling Checklist](./checklists/LOOSE_COUPLING_CHECKLIST.md) - Service independence
- [Venue/Event/Ticket Checklist](./checklists/VENUE_EVENT_TICKET_CHECKLIST.md) - Core features
- [Phase 1 Summary](./checklists/PHASE1_COMPLETION_SUMMARY.md) - Phase 1 completion status

---

## Infrastructure

| Component | Port | Purpose |
|-----------|------|---------|
| PostgreSQL Auth | 50432 | Auth service database |
| PostgreSQL Main | 50433 | Shared instance (separate DBs per service) |
| Redis | 50379 | Caching, sessions, distributed locks, queues |
| Kafka | 50092 | Event streaming between services |
| Zookeeper | 50181 | Kafka coordination |

---

## Security

- **Authentication**: JWT tokens issued by Auth Service
- **Authorization**: Role-based access control (user, organization, admin)
- **API Security**: Rate limiting at Gateway level
- **Payment**: Idempotency keys, Stripe webhook signature verification
- **Database**: Encryption at rest, per-service database isolation

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on:
- Code style conventions per language
- Branch naming and PR workflow
- Documentation standards
- Testing requirements

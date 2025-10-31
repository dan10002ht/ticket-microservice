# 📚 Ticket Booking System - Documentation

Central documentation hub for all microservices in the ticket booking system.

---

## 🗂️ Documentation Structure

```
docs/
├── README.md                           # This file - main index
├── architecture/                       # Architecture & Design docs
│   ├── AI_README.md                   # AI agent guidance
│   ├── MICROSERVICE_BEST_PRACTICES.md # Best practices guide
│   ├── PAYMENT_SERVICE.md             # Payment service design
│   └── SERVICE_CONNECTIONS.md         # Service integration map
├── checklists/                        # Implementation checklists
│   ├── AUTHORIZATION_CHECKLIST.md     # Auth implementation
│   ├── PHASE1_COMPLETION_SUMMARY.md   # Phase 1 summary
│   └── VENUE_EVENT_TICKET_CHECKLIST.md # Event/Ticket implementation
├── guides/                            # Setup & configuration guides
│   ├── go-service-structure.md        # Go service boilerplate
│   ├── MASTER_SLAVE_SETUP.md          # PostgreSQL replication
│   ├── METRICS_SETUP.md               # Monitoring setup
│   └── README_PORTS.md                # Port configuration
└── services/                          # Per-service documentation
    ├── auth/                          # Auth Service docs
    ├── email-worker/                  # Email Worker docs
    ├── event/                         # Event Service docs
    ├── gateway/                       # API Gateway docs
    ├── payment-service/               # Payment Service docs
    └── ticket/                        # Ticket Service docs
```

---

## 🏗️ Architecture Documentation

### Core Architecture

- **[AI_README.md](./architecture/AI_README.md)** - AI agent guidance, project structure
- **[MICROSERVICE_BEST_PRACTICES.md](./architecture/MICROSERVICE_BEST_PRACTICES.md)** - Best practices for microservices
- **[SERVICE_CONNECTIONS.md](./architecture/SERVICE_CONNECTIONS.md)** - Service integration map

### Service Design Documents

- **[PAYMENT_SERVICE.md](./architecture/PAYMENT_SERVICE.md)** - Payment Service design & strategy
  - Why Java for Payment Service
  - Architecture overview
  - Database design
  - Implementation roadmap

---

## ✅ Implementation Checklists

### Phase Summaries

- **[PHASE1_COMPLETION_SUMMARY.md](./checklists/PHASE1_COMPLETION_SUMMARY.md)** - Phase 1 completion status

### Feature Checklists

- **[AUTHORIZATION_CHECKLIST.md](./checklists/AUTHORIZATION_CHECKLIST.md)** - Authorization implementation
- **[VENUE_EVENT_TICKET_CHECKLIST.md](./checklists/VENUE_EVENT_TICKET_CHECKLIST.md)** - Event & Ticket features

---

## 📖 Setup & Configuration Guides

### Infrastructure

- **[MASTER_SLAVE_SETUP.md](./guides/MASTER_SLAVE_SETUP.md)** - PostgreSQL Master-Slave replication
- **[METRICS_SETUP.md](./guides/METRICS_SETUP.md)** - Prometheus & Grafana monitoring
- **[README_PORTS.md](./guides/README_PORTS.md)** - Port allocation reference

### Development

- **[go-service-structure.md](./guides/go-service-structure.md)** - Go service boilerplate structure

---

## 🎯 Service-Specific Documentation

### 💳 Payment Service (Java + Spring Boot)

**Location**: [docs/services/payment-service/](./services/payment-service/)

- **[README.md](./payment-service/README.md)** - Documentation index
- **[01_SETUP_COMPLETE.md](./payment-service/01_SETUP_COMPLETE.md)** - Project setup summary
- **[02_DATABASE_SETUP_COMPLETE.md](./payment-service/02_DATABASE_SETUP_COMPLETE.md)** - Database setup
- **[03_DATABASE_SCHEMA.md](./payment-service/03_DATABASE_SCHEMA.md)** - Schema reference

**Status**: ✅ Phase 1 - Database Setup Complete (20%)

---

### 🔐 Auth Service (Node.js)

**Location**: [docs/services/auth/](./services/auth/)

- **[CACHE_IMPLEMENTATION.md](./services/auth/CACHE_IMPLEMENTATION.md)** - Redis caching
- **[database-design.md](./services/auth/database-design.md)** - Database schema
- **[IMPLEMENTATION_CHECKLIST.md](./services/auth/IMPLEMENTATION_CHECKLIST.md)** - Implementation checklist
- **[INTEGRATION_FLOWS_README.md](./services/auth/INTEGRATION_FLOWS_README.md)** - Integration flows
- **[INTEGRATION_TODO.md](./services/auth/INTEGRATION_TODO.md)** - Integration tasks
- **[PGPOOL_COMPATIBILITY.md](./services/auth/PGPOOL_COMPATIBILITY.md)** - PgPool integration
- **[REGISTRATION_FLOWS_README.md](./services/auth/REGISTRATION_FLOWS_README.md)** - Registration flows
- **[REPOSITORY_PATTERN.md](./services/auth/REPOSITORY_PATTERN.md)** - Repository pattern

**Status**: ✅ Core Implementation Complete (90%)

---

### 🎪 Event Service (Go)

**Location**: [docs/services/event/](./services/event/)

- **[EVENT_NEW.md](./services/event/EVENT_NEW.md)** - Event model design
- **[README_EVENT_MODEL.md](./services/event/README_EVENT_MODEL.md)** - Event model reference

**Status**: ✅ Core Features Complete (80%)

---

### 🎫 Ticket Service (Go)

**Location**: [docs/services/ticket/](./services/ticket/)

- **[DATABASE_SETUP.md](./services/ticket/DATABASE_SETUP.md)** - Database setup guide
- **[IMPLEMENTATION_STATUS.md](./services/ticket/IMPLEMENTATION_STATUS.md)** - Implementation status

**Status**: ✅ Models, Repos, Services Complete (70%)

---

### 🌐 API Gateway (Node.js + Express)

**Location**: [docs/services/gateway/](./services/gateway/)

- **[API_CHECKLIST.md](./services/gateway/API_CHECKLIST.md)** - API implementation checklist
- **[API_EVENT_CREATION_CHECKLIST.md](./services/gateway/API_EVENT_CREATION_CHECKLIST.md)** - Event creation API
- **[AUTHORIZATION_GUIDE.md](./services/gateway/AUTHORIZATION_GUIDE.md)** - Authorization guide
- **[DEVELOPMENT.md](./services/gateway/DEVELOPMENT.md)** - Development guide
- **[EVENT_CREATION_CHECKLIST.md](./services/gateway/EVENT_CREATION_CHECKLIST.md)** - Event creation flow
- **[EVENT_DISPLAY_CHECKLIST.md](./services/gateway/EVENT_DISPLAY_CHECKLIST.md)** - Event display flow
- **[QUICK_CHECKLIST.md](./services/gateway/QUICK_CHECKLIST.md)** - Quick reference
- **[SWAGGER_README.md](./services/gateway/SWAGGER_README.md)** - Swagger/OpenAPI docs
- **[YARN_MIGRATION.md](./services/gateway/YARN_MIGRATION.md)** - Yarn migration guide

**Status**: ✅ Core Routes Complete (70%)

---

### 📧 Email Worker (Go)

**Location**: [docs/services/email-worker/](./services/email-worker/)

- **[API.md](./services/email-worker/API.md)** - API documentation
- **[folder-structure.md](./services/email-worker/folder-structure.md)** - Folder structure
- **[STEP.md](./services/email-worker/STEP.md)** - Implementation steps

**Status**: ✅ Core Implementation Complete (85%)

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    API GATEWAY (Node.js)                     │
│                     Port: 53000 (REST)                       │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┬───────────────┐
        │               │               │               │
        ▼               ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Auth Service │ │Event Service │ │Ticket Service│ │Payment Service│
│  (Node.js)   │ │    (Go)      │ │    (Go)      │ │   (Java)     │
│  Port: 50051 │ │  Port: 50053 │ │  Port: 50057 │ │  Port: 50056 │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
        │               │               │               │
        └───────────────┴───────────────┴───────────────┘
                        │
        ┌───────────────┼───────────────┬───────────────┐
        ▼               ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  PostgreSQL  │ │    Redis     │ │    Kafka     │ │  Prometheus  │
│ (Master/Slave)│ │ Port: 56379  │ │ Port: 59092  │ │ Port: 59090  │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
```

---

## 🔧 Quick Start

### Prerequisites

- Docker & Docker Compose
- Node.js 18+
- Go 1.19+
- Java 17+ (for Payment Service)
- PostgreSQL 15+
- Maven 3.8+ (for Payment Service)

### Start Development Environment

```bash
# Start all services
cd scripts
./dev-all.sh
```

### Access Points

| Service         | Type | URL/Port               |
| --------------- | ---- | ---------------------- |
| API Gateway     | REST | http://localhost:53000 |
| Payment Service | REST | http://localhost:8081  |
| Auth Service    | gRPC | localhost:50051        |
| Event Service   | gRPC | localhost:50053        |
| Ticket Service  | gRPC | localhost:50057        |
| Payment Service | gRPC | localhost:50056        |
| Email Worker    | gRPC | localhost:50060        |
| **Monitoring**  |      |                        |
| Grafana         | Web  | http://localhost:53001 |
| Prometheus      | Web  | http://localhost:59090 |
| Kibana          | Web  | http://localhost:55601 |

---

## 📊 Implementation Progress

| Service         | Setup | Database | Code | APIs | Tests | Status |
| --------------- | ----- | -------- | ---- | ---- | ----- | ------ |
| Payment Service | ✅    | ✅       | 🚧   | ⏳   | ⏳    | 20%    |
| Ticket Service  | ✅    | ✅       | ✅   | ✅   | ⏳    | 70%    |
| Event Service   | ✅    | ✅       | ✅   | ✅   | ⏳    | 80%    |
| Auth Service    | ✅    | ✅       | ✅   | ✅   | ⏳    | 90%    |
| Gateway         | ✅    | N/A      | ✅   | ✅   | ⏳    | 70%    |
| Email Worker    | ✅    | ✅       | ✅   | ✅   | ⏳    | 85%    |

**Legend**: ✅ Complete | 🚧 In Progress | ⏳ Not Started

---

## 🎯 Project Conventions

### Code Style

- **JavaScript/Node.js**: Functions over classes, camelCase files
- **Go**: Standard Go conventions, snake_case files
- **Java**: Standard Java conventions, PascalCase classes

### Service Naming

- Suffix with `-service` or `-worker`
- Examples: `payment-service`, `email-worker`

### Documentation

- Main README in service root
- Detailed docs in `docs/services/{service-name}/`
- Architecture docs in `docs/architecture/`
- Guides in `docs/guides/`

---

## 🔒 Security

- JWT Authentication (Auth Service)
- Role-Based Access Control (RBAC)
- API Rate Limiting (Gateway)
- Idempotency Keys (Payment Service)
- Database Encryption at Rest

---

## 📈 Monitoring

- **Metrics**: Prometheus for all services
- **Dashboards**: Grafana
- **Logging**: Elasticsearch + Kibana
- **Tracing**: Correlation IDs

---

## 🤝 Contributing

1. Follow project conventions
2. Update documentation
3. Write tests
4. Use feature branches
5. Submit pull requests

---

## 📝 Quick Links

### Architecture

- [Service Connections](./architecture/SERVICE_CONNECTIONS.md)
- [Best Practices](./architecture/MICROSERVICE_BEST_PRACTICES.md)
- [AI Agent Guide](./architecture/AI_README.md)

### Setup Guides

- [PostgreSQL Replication](./guides/MASTER_SLAVE_SETUP.md)
- [Monitoring Setup](./guides/METRICS_SETUP.md)
- [Port Configuration](./guides/README_PORTS.md)

### Service READMEs

- [Payment Service](./payment-service/README.md)
- [Auth Service](../auth-service/README.md)
- [Event Service](../event-service/)
- [Ticket Service](../ticket-service/)
- [Gateway](../gateway/README.md)
- [Email Worker](../email-worker/README.md)

---

**Last Updated**: 2024  
**Total Services**: 8  
**Documentation Files**: 40+  
**Implementation Status**: In Active Development

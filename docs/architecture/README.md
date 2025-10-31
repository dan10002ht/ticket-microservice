# 🏗️ Architecture Documentation

High-level architecture design and system documentation.

---

## 📋 Documents

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
- Communication protocols (gRPC, REST)
- Integration flows
- Service mesh topology

### Service Design Documents

#### [PAYMENT_SERVICE.md](./PAYMENT_SERVICE.md)

- **Complete Payment Service design document**
- Why Java Spring Boot for payment processing
- Architecture overview
- Database schema design
- Payment gateway integrations (Stripe, PayPal, VNPay, Momo)
- Security considerations
- Implementation roadmap (8 phases)
- Testing strategy
- Monitoring & observability

---

## 🎯 System Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    API GATEWAY                              │
│              (Node.js + Express)                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
      ┌───────────────┼───────────────────────────┐
      │               │                           │
      ▼               ▼                           ▼
┌──────────┐    ┌──────────┐              ┌──────────┐
│   Auth   │    │  Event   │              │ Payment  │
│ Service  │    │ Service  │    ...       │ Service  │
│(Node.js) │    │   (Go)   │              │  (Java)  │
└──────────┘    └──────────┘              └──────────┘
      │               │                           │
      └───────────────┴───────────────────────────┘
                      │
      ┌───────────────┼───────────────────────────┐
      ▼               ▼                           ▼
┌──────────┐    ┌──────────┐              ┌──────────┐
│PostgreSQL│    │  Redis   │              │  Kafka   │
│(Master-  │    │(Cache &  │              │(Event    │
│ Slave)   │    │ Session) │              │ Stream)  │
└──────────┘    └──────────┘              └──────────┘
```

---

## 🔗 Key Integration Patterns

### Synchronous Communication (gRPC)

- Gateway ↔ All Services
- Service ↔ Service (for immediate responses)

### Asynchronous Communication (Kafka)

- Payment events
- Email notifications
- Analytics data
- Audit logs

### Caching (Redis)

- Session storage
- API rate limiting
- Temporary data (idempotency keys)

### Database (PostgreSQL)

- Master-Slave replication
- PgPool-II for connection pooling
- Per-service database isolation

---

## 📊 Technology Stack

| Layer             | Technologies                  |
| ----------------- | ----------------------------- |
| **API Gateway**   | Node.js, Express, JWT         |
| **Services**      | Node.js, Go, Java Spring Boot |
| **Communication** | gRPC, REST, Kafka             |
| **Databases**     | PostgreSQL (Master-Slave)     |
| **Caching**       | Redis                         |
| **Monitoring**    | Prometheus, Grafana           |
| **Logging**       | Elasticsearch, Kibana         |
| **Containers**    | Docker, Docker Compose        |

---

## 🎯 Design Principles

1. **Microservices First**: Each service is independently deployable
2. **Polyglot Architecture**: Use best language for each service
3. **Database Per Service**: No shared databases
4. **API Gateway Pattern**: Single entry point for clients
5. **Event-Driven**: Async communication for loosely coupled services
6. **Observability**: Comprehensive logging, metrics, tracing
7. **Resilience**: Circuit breakers, retries, fallbacks
8. **Security**: JWT auth, RBAC, encryption at rest/transit

---

## 📖 Related Documentation

- [Main Documentation Index](../README.md)
- [Implementation Checklists](../checklists/)
- [Setup Guides](../guides/)
- [Service Documentation](../services/)

---

**Last Updated**: 2024  
**Architecture Version**: 1.0



# 📘 Booking Service - Documentation

Complete documentation hub for the Booking Service microservice implementation.

---

## 📚 Documentation Index

### Phase 1: Core Setup

1. **[01_SETUP_COMPLETE.md](./01_SETUP_COMPLETE.md)** ⏳ Planned
   - Project bootstrap (Spring Boot 3.2, Java 17, Maven)
   - Folder structure, base modules, shared configs
   - Docker and dev profile wiring

2. **[02_DATABASE_SETUP_COMPLETE.md](./02_DATABASE_SETUP_COMPLETE.md)** ⏳ Planned
   - Booking + booking_item tables
   - Flyway migrations, indexes, partitions, retention policies
   - Connection pooling + PgPool compatibility

3. **[03_DOMAIN_MODEL.md](./03_DOMAIN_MODEL.md)** ⏳ Planned
   - Booking lifecycle states + transitions
   - Saga orchestration entities (booking, booking_action, booking_task)
   - Redis lock keys, TTL strategy

---

## 🎯 Implementation Progress

| Phase | Section | Status | Documents |
| --- | --- | --- | --- |
| **Phase 1: Core Setup** | Project Setup | ⏳ Not Started | - |
|  | Database Setup | ⏳ Not Started | - |
|  | Entity Models | ⏳ Not Started | - |
|  | Repository Layer | ⏳ Not Started | - |
| **Phase 2: Booking Flow** | Saga Orchestrator | ⏳ Not Started | - |
|  | gRPC API | ⏳ Not Started | - |
|  | Kafka Integration | ⏳ Not Started | - |
|  | Redis Locking | ⏳ Not Started | - |
| **Phase 3: Integrations** | Ticket Service integration | ⏳ Not Started | - |
|  | Payment Service integration | ⏳ Not Started | - |
|  | Notification / Email hooks | ⏳ Not Started | - |

---

## 🛠️ Technical Stack

- **Language**: Java 17
- **Framework**: Spring Boot 3.2.x
- **Database**: PostgreSQL 15 (Flyway migrations)
- **Cache / Lock**: Redis (Redisson)
- **Messaging**: Kafka (booking lifecycle events)
- **Build Tool**: Maven
- **APIs**: gRPC (internal), REST optional through gateway

---

## 🔄 Responsibilities & Flow

1. Receive booking requests from gateway / booking-worker
2. Acquire distributed lock per booking session / seat cluster
3. Call Ticket Service gRPC to reserve seats
4. Trigger Payment Service gRPC for authorization / capture
5. Coordinate saga completion (confirm or rollback)
6. Publish booking events to Kafka (notification, realtime, analytics)
7. Manage booking lifecycle (pending → processing → confirmed/cancelled)

---

## 🚀 Upcoming Boilerplate Tasks

1. Scaffold Maven module (copy baseline from `payment-service`)
2. Set up Spring Boot entrypoint + basic config (profiles, logging)
3. Add Flyway with initial migration placeholders
4. Wire Redis + Kafka configs (empty beans with TODO markers)
5. Create skeleton packages (controller, service, repository, saga, grpc)

---

## ✅ TODO / Deferred Items

- [ ] Define booking saga specification (per state machine doc)
- [ ] Add grpc stubs referencing `shared-lib/protos/booking.proto`
- [ ] Document Redis key schema + TTL enforcement
- [ ] Write integration plan for Ticket + Payment service calls
- [ ] Prepare sequence diagrams for queue + booking service interplay

---

## 📝 Related Documentation

- [Main Project Documentation](../README.md)
- [Payment Service Docs](../payment-service/README.md)
- [Ticket Service Docs](../services/ticket/IMPLEMENTATION_STATUS.md)
- [Architecture Overview](../architecture/AI_README.md)
- [Service Connections](../architecture/SERVICE_CONNECTIONS.md)

---

**Last Updated**: 2024
**Status**: Planning Phase

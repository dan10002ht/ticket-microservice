# 🔄 Booking Worker Service - Documentation

Complete documentation hub for the Booking Worker Service (Go) implementation.

---

## 📚 Documentation Index

### Phase 1: Core Setup

1. **[01_SETUP_COMPLETE.md](./01_SETUP_COMPLETE.md)** ⏳ Planned
   - Go module setup (go.mod, go.sum)
   - Project structure (cmd, internal, pkg, queue, grpc)
   - Dependencies (go-redis, grpc-go, zap, prometheus)
   - Dockerfile + docker-compose integration
   - Development environment setup

2. **[02_QUEUE_SETUP_COMPLETE.md](./02_QUEUE_SETUP_COMPLETE.md)** ⏳ Planned
   - Redis queue implementation (go-redis)
   - Kafka queue implementation (segmentio/kafka-go)
   - Queue manager interface + implementations
   - Queue position tracking (Redis sorted sets)
   - Timeout handling (Redis TTL + sorted sets)

3. **[03_ARCHITECTURE.md](./03_ARCHITECTURE.md)** ⏳ Planned
   - Component architecture (queue, processor, notifier)
   - Goroutine patterns (worker pools, channels)
   - Distributed locking strategy (go-redsync)
   - Error handling & retry logic
   - Metrics & observability

---

## 🎯 Implementation Progress

| Phase | Section | Status | Documents |
| --- | --- | --- | --- |
| **Phase 1: Core Setup** | Project Setup | ⏳ Not Started | - |
|  | Queue Infrastructure | ⏳ Not Started | - |
|  | gRPC Client Setup | ⏳ Not Started | - |
| **Phase 2: Queue Management** | Enqueue Logic | ⏳ Not Started | - |
|  | Dequeue & Processing | ⏳ Not Started | - |
|  | Position Tracking | ⏳ Not Started | - |
|  | Timeout Handler | ⏳ Not Started | - |
| **Phase 3: Integration** | Booking Service gRPC | ⏳ Not Started | - |
|  | Realtime Service gRPC | ⏳ Not Started | - |
|  | Gateway Integration | ⏳ Not Started | - |
|  | Metrics & Monitoring | ⏳ Not Started | - |

---

## 🛠️ Technical Stack

- **Language**: Go 1.21+
- **Framework**: Go kit / Fiber / Gin (suggested)
- **Queue**: Redis (go-redis) or Kafka (segmentio/kafka-go)
- **gRPC**: google.golang.org/grpc
- **WebSocket**: gorilla/websocket (for realtime notifications)
- **Monitoring**: Prometheus client_golang
- **Logging**: zap / zerolog
- **Distributed Locking**: go-redsync (Redlock)

---

## 🔄 Responsibilities & Flow

1. **Receive booking requests** from Gateway via gRPC
2. **Enqueue requests** into distributed queue (Redis/Kafka)
3. **Track queue position** for each client (Redis sorted sets)
4. **Notify clients** of queue position via Realtime Service
5. **Dequeue and process** requests in order (fairness)
6. **Call Booking Service** to reserve tickets
7. **Handle timeouts** - release tickets if payment not completed
8. **Scale horizontally** - multiple workers process queue in parallel

---

## 🚀 Upcoming Boilerplate Tasks

1. Initialize Go module (`go mod init`)
2. Set up project structure (cmd, internal, pkg, queue, grpc)
3. Add dependencies (go-redis, grpc-go, zap, prometheus)
4. Create queue manager interface + Redis implementation
5. Set up gRPC client stubs (booking-service, realtime-service)
6. Implement worker processor with goroutines
7. Add timeout handler (Redis TTL + cleanup)
8. Wire Prometheus metrics

---

## 🔗 Integration Points

### Dependencies (booking-worker calls these)
- **Booking Service** (Java): Reserve and confirm bookings via gRPC
- **Realtime Service**: Send queue position updates via gRPC
- **Ticket Service** (Go): Check availability (optional, may go through booking-service)

### Dependents (these call booking-worker)
- **Gateway**: Routes booking requests to booking-worker
- **Clients**: Poll queue status or receive WebSocket updates

---

## 📊 Queue Flow Diagram

```
Client → Gateway → booking-worker (enqueue)
                           ↓
                    Redis/Kafka Queue
                           ↓
              booking-worker (dequeue)
                           ↓
                    Booking Service
                    (reserve tickets)
                           ↓
                    Client notified
                    (proceed to payment)
```

---

## ✅ TODO / Deferred Items

- [ ] Define protobuf messages for `BookingWorkerService` in `shared-lib/protos/`
- [ ] Document Redis key schema (`booking-queue`, `booking-timeouts`, `booking-status:{clientId}`)
- [ ] Design queue position calculation algorithm
- [ ] Plan timeout strategy (TTL, cleanup workers)
- [ ] Write integration tests with Testcontainers (Redis, gRPC mocks)
- [ ] Prepare load testing strategy (100k+ concurrent clients)

---

## 📝 Related Documentation

- [Main Project Documentation](../README.md)
- [Booking Service Docs](../booking-service/README.md)
- [Payment Service Docs](../payment-service/README.md)
- [Ticket Service Docs](../services/ticket/IMPLEMENTATION_STATUS.md)
- [Architecture Overview](../architecture/AI_README.md)
- [Service Connections](../architecture/SERVICE_CONNECTIONS.md)

---

## 🎯 Implementation Order Recommendation

**Recommended: Implement `booking-service` FIRST, then `booking-worker`**

### Why booking-service first?
1. **Core dependency**: booking-worker calls booking-service via gRPC
2. **Business logic foundation**: Need booking-service APIs defined before worker can use them
3. **Testing**: Can test booking-service independently, then integrate with worker
4. **Incremental development**: Build core logic first, then add queue layer

### Why booking-worker second?
1. **Depends on booking-service**: Worker needs booking-service gRPC stubs
2. **Queue is infrastructure layer**: Can be added after business logic is stable
3. **Easier testing**: Can mock booking-service calls during worker development

### Alternative: Parallel development
- If teams are separate, can develop in parallel with:
  - booking-service: Define gRPC contracts first (protobuf)
  - booking-worker: Implement queue logic, mock booking-service calls
  - Integration: Connect both when ready

---

**Last Updated**: 2024  
**Status**: Planning Phase


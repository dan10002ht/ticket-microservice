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

### Phase 4: Infrastructure & Operations

4. **[REDIS_KEY_SCHEMA.md](./REDIS_KEY_SCHEMA.md)** ✅ Complete
   - Redis key schema và naming conventions
   - TTL enforcement strategy cho distributed locks
   - Lock management best practices
   - Troubleshooting guide

---

## 🎯 Implementation Progress

| Phase                     | Section                     | Status      | Documents                                                                          |
| ------------------------- | --------------------------- | ----------- | ---------------------------------------------------------------------------------- |
| **Phase 1: Core Setup**   | Project Setup               | ✅ Complete | Source code (`pom.xml`, `BookingServiceApplication.java`)                          |
|                           | Database Setup              | ✅ Complete | Migrations (`V1__create_bookings_table.sql`, `V2__create_booking_items_table.sql`) |
|                           | Entity Models               | ✅ Complete | Source code (`entity/Booking.java`, `entity/BookingItem.java`)                     |
|                           | Repository Layer            | ✅ Complete | Source code (`repository/BookingRepository.java`, `BookingItemRepository.java`)    |
| **Phase 2: Booking Flow** | Saga Orchestrator           | ✅ Complete | `service/saga/BookingSagaOrchestrator.java` - Full saga with compensation          |
|                           | gRPC API                    | ✅ Complete | `grpc/BookingGrpcService.java` - All RPCs with error handling                      |
|                           | Kafka Integration           | ✅ Complete | `service/BookingEventPublisher.java` + `config/KafkaConfig.java`                   |
|                           | Redis Locking               | ✅ Complete | `service/BookingLockService.java` + `config/RedisConfig.java`                      |
| **Phase 3: Integrations** | Ticket Service integration  | ✅ Complete | `grpcclient/TicketServiceClient.java` - Reserve/release tickets                    |
|                           | Payment Service integration | ✅ Complete | `grpcclient/PaymentServiceClient.java` - Create/capture/cancel payment             |
|                           | Error Handling              | ✅ Complete | `exception/GlobalExceptionHandler.java` + exception classes                        |
|                           | Retry Logic                 | ✅ Complete | `config/GrpcRetryConfig.java` + `@Retryable` on gRPC clients                       |
|                           | Metrics & Monitoring        | ✅ Complete | `metrics/BookingMetricsService.java` + Prometheus integration                      |

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

## ✅ Completed Implementation

### Core Components

- ✅ **Project Structure**: Maven module với Spring Boot 3.2, Java 17
- ✅ **Database**: Flyway migrations cho bookings và booking_items tables
- ✅ **Entities**: Booking, BookingItem với lifecycle states
- ✅ **Repositories**: JPA repositories với custom queries
- ✅ **Service Layer**: BookingService với saga orchestration
- ✅ **Saga Orchestrator**: Full booking flow với compensation logic
- ✅ **gRPC Clients**: TicketServiceClient, PaymentServiceClient
- ✅ **gRPC Server**: BookingGrpcService với error handling
- ✅ **Error Handling**: Global exception handler (REST + gRPC)
- ✅ **Configuration**: Redis, Kafka, gRPC, Database configs

### Saga Flow Implementation

1. ✅ Acquire distributed lock (Redis)
2. ✅ Create booking record (PENDING)
3. ✅ Reserve seats via Ticket Service (RESERVING → AWAITING_PAYMENT)
4. ✅ Process payment via Payment Service (PROCESSING_PAYMENT)
5. ✅ Confirm booking (CONFIRMED)
6. ✅ Compensation logic (release seats, cancel payment on failure)

---

## ✅ TODO / Deferred Items

- [x] REST API endpoints (optional - can go through gateway) ✅
- [x] Retry logic cho gRPC calls ✅
- [x] Metrics và monitoring integration ✅
- [x] Document Redis key schema + TTL enforcement ✅
- [ ] Unit tests và integration tests
- [ ] Prepare sequence diagrams for queue + booking service interplay

## 🎯 Implementation Order Recommendation

**Recommended: Implement `booking-service` FIRST, then `booking-worker`**

### Why booking-service first?

1. **Core dependency**: booking-worker calls booking-service via gRPC
2. **Business logic foundation**: Need booking-service APIs defined before worker can use them
3. **Testing**: Can test booking-service independently, then integrate with worker
4. **Incremental development**: Build core logic first, then add queue layer

### Integration Flow

```
Gateway → booking-worker (queue) → booking-service (orchestrate)
                                              ↓
                                    ticket-service + payment-service
```

See [Booking Worker Docs](../booking-worker/README.md) for worker implementation details.

---

## 📝 Related Documentation

- [Main Project Documentation](../README.md)
- [Payment Service Docs](../payment-service/README.md)
- [Ticket Service Docs](../services/ticket/IMPLEMENTATION_STATUS.md)
- [Architecture Overview](../architecture/AI_README.md)
- [Service Connections](../architecture/SERVICE_CONNECTIONS.md)
- [Redis Key Schema & TTL Enforcement](./REDIS_KEY_SCHEMA.md) ⭐

---

**Last Updated**: 2024
**Status**: Phase 1-3 Core Implementation Complete ✅

## 📋 Implementation Summary

### What's Been Implemented

1. **Core Infrastructure**

   - Spring Boot application với full configuration
   - Database migrations (Flyway)
   - Redis và Kafka integration
   - gRPC server và client setup

2. **Business Logic**

   - Booking saga orchestrator với full compensation
   - Integration với Ticket Service (seat reservation)
   - Integration với Payment Service (payment processing)
   - Distributed locking (Redis/Redisson)

3. **API Layer**

   - gRPC service với error handling
   - Exception handling (REST + gRPC)
   - Event publishing (Kafka)

4. **Error Handling**

   - Custom exception classes
   - Global exception handler
   - Proper gRPC status codes

5. **Retry Logic**

   - Spring Retry integration với exponential backoff
   - Retryable gRPC calls (TicketServiceClient, PaymentServiceClient)
   - Smart retry: chỉ retry transient errors (DEADLINE_EXCEEDED, UNAVAILABLE, INTERNAL)
   - Non-retryable errors: INVALID_ARGUMENT, NOT_FOUND, PERMISSION_DENIED, FAILED_PRECONDITION
   - Configurable retry policy (max attempts, backoff intervals)

6. **Metrics & Monitoring**
   - Prometheus metrics integration (Micrometer)
   - Custom metrics cho booking lifecycle events
   - Saga orchestration metrics (duration, step counts, compensation)
   - gRPC call metrics (success/failure counts)
   - Exposed via `/actuator/prometheus` endpoint

### Retry Configuration

Retry logic được cấu hình trong `application.yml`:

```yaml
booking:
  grpc:
    retry:
      max-attempts: 3 # Số lần retry tối đa
      initial-interval-ms: 500 # Khoảng thời gian chờ ban đầu (ms)
      multiplier: 2.0 # Hệ số nhân cho exponential backoff
      max-interval-ms: 5000 # Khoảng thời gian chờ tối đa (ms)
```

**Retry Strategy:**

- **Exponential Backoff**: 500ms → 1000ms → 2000ms (max 5000ms)
- **Retryable Errors**: `DEADLINE_EXCEEDED`, `UNAVAILABLE`, `INTERNAL`, `RESOURCE_EXHAUSTED`
- **Non-Retryable Errors**: `INVALID_ARGUMENT`, `NOT_FOUND`, `PERMISSION_DENIED`, `FAILED_PRECONDITION`

### Metrics & Monitoring

Booking Service exposes Prometheus metrics via `/actuator/prometheus` endpoint.

#### Available Metrics

**Booking Lifecycle Metrics:**

- `booking.created` - Total bookings created (counter)
- `booking.confirmed` - Total bookings confirmed (counter)
- `booking.cancelled` - Total bookings cancelled (counter, tagged with reason)
- `booking.failed` - Total bookings failed (counter, tagged with error_type)

**Saga Orchestration Metrics:**

- `saga.execution.duration` - Saga execution duration in seconds (timer, tagged with status: success/failed)
- `saga.step.executed` - Total saga steps executed (counter, tagged with step name and status)
- `saga.compensation.triggered` - Total compensations triggered (counter, tagged with step and reason)

**gRPC Call Metrics:**

- `grpc.calls.total` - Total gRPC calls (counter, tagged with target_service, method, status)
- `grpc.calls.errors` - Total gRPC call errors (counter, tagged with target_service, method, error_code)

#### Example Prometheus Queries

```promql
# Booking success rate
rate(booking_confirmed_total[5m]) / rate(booking_created_total[5m])

# Saga failure rate
rate(saga_execution_duration_seconds_count{status="failed"}[5m]) / rate(saga_execution_duration_seconds_count[5m])

# Average saga execution time
rate(saga_execution_duration_seconds_sum[5m]) / rate(saga_execution_duration_seconds_count[5m])

# gRPC error rate by service
rate(grpc_calls_errors_total[5m]) / rate(grpc_calls_total[5m])
```

#### Metrics Configuration

Metrics are automatically configured via Spring Boot Actuator:

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health,metrics,prometheus,info
  metrics:
    export:
      prometheus:
        enabled: true
    tags:
      application: booking-service
```

### Next Steps (Optional)

- Metrics và monitoring integration (Prometheus)
- Unit tests và integration tests
- Document Redis key schema + TTL enforcement
- Prepare sequence diagrams for queue + booking service interplay

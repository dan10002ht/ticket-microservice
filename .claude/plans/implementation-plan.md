# Comprehensive Review: Ticket Microservices System

## Executive Summary

Dự án ticketing microservices có kiến trúc tốt với việc sử dụng các patterns chuẩn enterprise: **Saga Pattern**, **Distributed Locking**, **Worker Pool**, **Idempotency**.

**Update (Dec 2024):**
- Đã hoàn thành Loose Coupling optimizations (Transactional Outbox, Async Payment, Kafka Consumers, Worker DLQ, State Machine)
- **COMPLETED:** User Service + Realtime Service implementation
- Phần Scale 1M+ users để **PENDING** cho sau

---

# Service Inventory & Status

## Service Status Overview (Updated)

| Service | Language | Status | Priority |
|---------|----------|--------|----------|
| Gateway | Node.js | ✅ Complete | - |
| Auth Service | Node.js | ✅ Complete | - |
| Event Service | Node.js | ✅ Complete | - |
| Ticket Service | Node.js | ✅ Complete | - |
| Booking Service | Java | ✅ Complete | - |
| Payment Service | Java | ✅ Complete | - |
| Booking Worker | Go | ✅ Complete | - |
| Email Worker | Go | ✅ Complete | - |
| **User Service** | **Go** | ✅ **COMPLETED** | - |
| **Realtime Service** | **Go** | ✅ **COMPLETED** | - |
| Device Service | Node.js | 🟡 Partial | MEDIUM |
| Security Service | Node.js | 🟡 Partial | LOW |
| Notification Service | - | ❌ Stub only | LOW |
| Pricing Service | - | ❌ Stub only | LOW |
| Analytics Service | - | ❌ Stub only | LOW |
| Invoice Service | - | ❌ Stub only | LOW |
| Support Service | - | ❌ Stub only | LOW |
| Checkin Service | - | ❌ Stub only | LOW |

---

# COMPLETED: Phase 1 Implementation

## User Service (Go) - DONE

**Location:** `user-service/`
**Port:** 50052 (gRPC), 9092 (Metrics)
**Database:** PostgreSQL

### Features Implemented:
- Profile CRUD (GetProfile, CreateProfile, UpdateProfile)
- Address management (GetAddresses, AddAddress, UpdateAddress, DeleteAddress)
- Legacy compatibility (GetUser, CreateUser, ListUsers)
- Health checks (`/health`, `/ready`)
- Prometheus metrics (`/metrics`)
- Database migrations

### Structure:
```
user-service/
├── cmd/server/main.go
├── config/config.go
├── internal/
│   ├── grpc/server.go
│   ├── grpc/handlers/profile_handler.go
│   ├── service/profile_service.go, address_service.go
│   ├── repository/profile_repository.go, address_repository.go
│   └── model/profile.go, address.go
├── migrations/001_init.sql
├── Dockerfile
└── README.md
```

---

## Realtime Service (Go + Gorilla WebSocket) - DONE

**Location:** `realtime-service/`
**Ports:** 3003 (HTTP/WebSocket), 50057 (gRPC), 9057 (Metrics)
**Dependencies:** Redis (Pub/Sub)

### Features Implemented:
- WebSocket server with JWT authentication
- Room-based message broadcasting
- User-specific notifications
- Redis Pub/Sub for horizontal scaling
- gRPC API for internal service calls:
  - `NotifyBookingResult` - Booking status updates
  - `NotifyQueuePosition` - Queue position updates
  - `NotifyPaymentStatus` - Payment status updates
  - `BroadcastEvent` - Room broadcasts
  - `SendToUser` - Direct user messages
  - `GetConnectionStats` - Connection statistics
- Prometheus metrics
- Health checks

### Structure:
```
realtime-service/
├── main.go
├── config/config.go
├── internal/
│   ├── grpc/server.go
│   ├── grpc/handlers/notification_handler.go
│   ├── websocket/hub.go, client.go, server.go, message.go
│   ├── pubsub/subscriber.go, publisher.go
│   ├── service/notification_service.go
│   └── middleware/auth.go
├── metrics/metrics.go
├── Dockerfile
└── README.md
```

---

# NEXT STEPS (Priority Order)

## Immediate (Before Production)

### 1. Integration: Gateway → User Service
**Priority:** HIGH
**Effort:** 2-4 hours

Gateway đang có routes đến user-service nhưng chưa kết nối:
```javascript
// gateway/src/routes/users.js - cần update
GET  /api/users/profile       → userService.getProfile()
PUT  /api/users/profile       → userService.updateProfile()
GET  /api/users/addresses     → userService.getAddresses()
POST /api/users/addresses     → userService.addAddress()
```

**Tasks:**
- [ ] Update gateway gRPC client config for user-service (port 50052)
- [ ] Test end-to-end user profile flow
- [ ] Update gateway environment variables

### 2. Integration: Booking Worker → Realtime Service ✅ DONE
**Priority:** HIGH
**Status:** COMPLETED (Dec 18, 2024)

**Completed:**
- [x] Fixed port config: 50060 → 50057
- [x] Full gRPC client implementation (NotifyBookingResult, NotifyQueuePosition, NotifyPaymentStatus, BroadcastEvent, SendToUser)
- [x] Proto generation script updated
- [x] Proto files generated in `booking-worker/internal/protos/realtime/`
- [x] Code review verified integration points:
  - [processor.go:268](booking-worker/internal/worker/processor.go#L268) - NotifyBookingResult on success
  - [processor.go:316](booking-worker/internal/worker/processor.go#L316) - NotifyBookingResult on failure
- [x] Test script created: `scripts/test-e2e-realtime.sh`

**Ready for Testing (requires Docker):**
- [ ] Test booking notification flow end-to-end
- [ ] Verify WebSocket delivery to frontend

### 3. Database Setup
**Priority:** HIGH
**Status:** Script created, ready to run

**Script:** `scripts/setup-databases.sh`

**Manual commands (if script not available):**
```bash
# Using docker-compose postgres (port 55435)
export PGPASSWORD=booking_pass
createdb -h localhost -p 55435 -U booking_user user_service
psql -h localhost -p 55435 -U booking_user -d user_service -f user-service/migrations/001_init.sql

# Verify
psql -h localhost -p 55435 -U booking_user -d user_service -c '\dt'
```

**Tasks:**
- [x] Create setup script: `scripts/setup-databases.sh`
- [ ] Run: `./scripts/setup-databases.sh`
- [ ] Verify connection from user-service

### 4. Redis Setup for Realtime Service ✅ Ready
**Priority:** HIGH
**Status:** Already configured in docker-compose.dev.yml

**Docker Compose config (already set):**
- `REDIS_HOST=redis` (docker network hostname)
- `REDIS_PORT=6379` (internal port)

**Tasks:**
- [x] Redis config in docker-compose.dev.yml
- [x] REDIS_HOST, REDIS_PORT set for realtime-service
- [ ] Test Pub/Sub connectivity when services are running

---

## Short-term (Week 1-2)

### 5. Update Documentation ✅ DONE
**Priority:** MEDIUM
**Status:** COMPLETED (Dec 18, 2024)

**Completed:**
- [x] `docs/user-service/README.md` - Created
- [x] `docs/realtime-service/README.md` - Created
- [x] `docs/architecture/SERVICE_CONNECTIONS.md` - Completely rewritten
- [x] `docs/architecture/README.md` - Updated with new diagram
- [x] `docs/services/README.md` - Added all 10 services
- [x] `docs/guides/README_PORTS.md` - Updated with new ports
- [x] `docs/services/auth/INTEGRATION_FLOWS_README.md` - Cleaned up outdated references
- [x] `docs/services/auth/INTEGRATION_TODO.md` - Updated status

### 6. Docker Compose Update ✅ DONE
**Priority:** MEDIUM
**Status:** COMPLETED (Dec 18, 2024)

**Completed:**
- [x] Added user-service to `deploy/docker-compose.dev.yml` (ports 50052, 9092)
- [x] Added realtime-service to `deploy/docker-compose.dev.yml` (ports 3003, 50057, 9057)
- [x] Added booking-worker to `deploy/docker-compose.dev.yml` (ports 50059, 9091)
- [x] Fixed gateway GRPC_USER_SERVICE_URL: `user-profile` → `user-service`
- [x] Added gateway GRPC_REALTIME_SERVICE_URL

### 7. Frontend WebSocket Integration
**Priority:** MEDIUM
**Effort:** 4-8 hours

**Tasks:**
- [ ] Create WebSocket client in frontend
- [ ] Handle booking status notifications
- [ ] Handle payment status notifications
- [ ] Implement reconnection logic

---

## Medium-term (Week 3-4)

### 8. Complete Gateway Event Handlers ✅ DONE
**Priority:** MEDIUM
**Status:** COMPLETED (Dec 18, 2024)

All endpoints implemented in [eventHandlers.js](gateway/src/handlers/eventHandlers.js):
- [x] `POST /api/events/:id/layout` - saveEventLayoutHandler (stores in canvasConfig)
- [x] `POST /api/events/:id/pricing` - saveEventPricingHandler (stores in metadata)
- [x] `POST /api/events/:id/publish` - publishEventHandler (updates status)
- [x] `GET /api/events/templates` - getEventTemplatesHandler (static templates)
- [x] `POST /api/events/:id/duplicate` - duplicateEventHandler (copy event)

**Notes:**
- Pricing uses simplified metadata storage (full PricingService integration pending)
- Templates are static, could be moved to database in production

### 9. Device Service Implementation
**Priority:** LOW

- Device fingerprinting
- Multi-device session management

### 10. Security Service Implementation
**Priority:** LOW

- Audit logging
- Threat detection

---

## Long-term (Pending)

### Scale to 1M+ Users
**Status:** PENDING (reference old plan section)

- Redis Cluster migration
- Database sharding
- Kafka partition optimization
- Auto-scaling configuration

---

# Documentation Updates Required

## Files to Create/Update

| File | Action | Priority | Status |
|------|--------|----------|--------|
| `docs/user-service/README.md` | Create | HIGH | ✅ Done |
| `docs/realtime-service/README.md` | Create | HIGH | ✅ Done |
| `docs/architecture/SERVICE_CONNECTIONS.md` | Rewrite | HIGH | ✅ Done |
| `docs/architecture/README.md` | Update diagram | MEDIUM | ✅ Done |
| `docs/services/README.md` | Add new services | MEDIUM | ✅ Done |
| `docs/guides/README_PORTS.md` | Add new ports | MEDIUM | ✅ Done |
| `docker-compose.dev.yml` | Add new services | MEDIUM | ✅ Done |

## Service Communication Diagram Update

```
┌─────────────┐     ┌─────────────┐     ┌─────────────────┐
│   Frontend  │────▶│   Gateway   │────▶│  User Service   │
│  (React)    │     │  (Node.js)  │     │     (Go)        │
└──────┬──────┘     └──────┬──────┘     └─────────────────┘
       │                   │
       │ WebSocket         │ gRPC
       ▼                   ▼
┌─────────────────┐  ┌─────────────┐     ┌─────────────────┐
│ Realtime Service│◀─│   Booking   │────▶│ Booking Service │
│ (Go+WebSocket)  │  │   Worker    │     │     (Java)      │
└────────┬────────┘  │    (Go)     │     └─────────────────┘
         │           └─────────────┘
         │
    Redis Pub/Sub
    (Horizontal Scaling)
```

---

# Quick Start Commands

```bash
# 1. Setup User Service Database
createdb user_service
psql -d user_service -f user-service/migrations/001_init.sql

# 2. Start Redis (for Realtime Service)
docker run -d -p 6379:6379 redis:7-alpine

# 3. Run User Service
cd user-service && go run main.go

# 4. Run Realtime Service
cd realtime-service && go run main.go

# 5. Test User Service gRPC
grpcurl -plaintext localhost:50052 list
grpcurl -plaintext -d '{"user_id":"test-uuid"}' localhost:50052 user.UserService/GetProfile

# 6. Test Realtime Service gRPC
grpcurl -plaintext localhost:50057 list
grpcurl -plaintext localhost:50057 realtime.RealtimeService/GetConnectionStats

# 7. Test WebSocket
wscat -c "ws://localhost:3003/ws"
```

---

# Reference: Original Issues (From Code Review)

## P0 - CRITICAL (Fixed)
- [x] Replace Redis KEYS with SCAN
- [x] Add idempotency key
- [x] Implement DLQ for compensations
- [x] Add authorization check for queue operations

## P1 - HIGH (Fixed)
- [x] Configure proper lock lease time
- [x] Implement circuit breaker pattern
- [x] Add fair queue processing
- [x] Implement transactional outbox pattern

## P2 - MEDIUM (Backlog)
- [ ] Add exponential backoff with jitter
- [ ] Optimize Remove operation complexity
- [ ] Configure gRPC connection pooling
- [ ] Optimize timeout handler with Lua scripts

# Comprehensive Review: Ticket Microservices System

## Executive Summary

Dự án ticketing microservices của bạn có kiến trúc tốt với việc sử dụng các patterns chuẩn enterprise: **Saga Pattern**, **Distributed Locking**, **Worker Pool**, **Idempotency**. Tuy nhiên, để handle **100k concurrent users**, có một số điểm cần tối ưu **CRITICAL** và **HIGH priority**.

---

## 1. Tech Stack Overview

| Service | Language | Framework | Purpose |
|---------|----------|-----------|---------|
| Gateway | Node.js | Express | API entry point |
| Booking Service | Java | Spring Boot 3.2 | Core booking logic, Saga orchestrator |
| Booking Worker | Go | Native | Queue-based distributed processing |
| Payment Service | Java | Spring Boot | Payment processing (Stripe) |
| Email Worker | Go | Native | Async email delivery |
| Event Service | Go | Native | Event management |

**Infrastructure:**
- Database: PostgreSQL (Master-Slave replication)
- Cache/Queue: Redis 7
- Event Streaming: Kafka
- Monitoring: Prometheus + Grafana
- Inter-service: gRPC

---

## 2. CRITICAL Issues (P0 - Fix Immediately)

### 2.1 Race Condition: Redis KEYS Command in Production

**File:** `booking-worker/internal/queue/redis_queue.go:104`

```go
// PROBLEM: KEYS is O(N) and BLOCKS Redis
keys, err := r.client.Keys(ctx, "booking-queue:*").Result()
```

**Impact:**
- KEYS command scans ALL keys in Redis - blocks entire Redis instance
- At 100k concurrent users with multiple events, this will cause severe latency spikes
- Redis documentation explicitly warns against using KEYS in production

**Fix:**
- Use Redis `SCAN` command with cursor-based iteration
- Better: Maintain a SET of active event queues (`active-booking-queues`) and use `SMEMBERS`
- Best: Use Redis Streams instead of List + Sorted Set combination

---

### 2.2 Missing Idempotency Key in Booking Worker → Booking Service

**File:** `booking-worker/internal/worker/processor.go:190-198`

```go
bookingID, err := p.bookingClient.CreateBooking(
    ctx,
    item.UserID,
    item.EventID,
    item.SeatNumbers,
    // MISSING: Idempotency key (item.ID)
)
```

**Impact:**
- If network fails AFTER booking created but BEFORE response received, retry creates **duplicate booking**
- User pays twice for same seats
- Financial loss and customer complaints

**Fix:**
- Pass `item.ID` (queue item ID) as idempotency key to booking-service
- Booking-service should check this key before creating booking
- Add unique constraint on idempotency_key in booking table

---

### 2.3 Compensation Logic Silent Failures

**File:** `booking-service/.../BookingSagaOrchestrator.java:128-165`

```java
if (paymentId != null) {
    try {
        paymentServiceClient.cancelPayment(paymentId);
    } catch (Exception e) {
        log.error("Failed to compensate payment: {}", paymentId, e);
        // PROBLEM: Silent failure - no DLQ, no retry, no alerting
    }
}
```

**Impact:**
- Failed compensations leave system in inconsistent state
- Money refunded in payment but seats still reserved
- No automatic recovery mechanism

**Fix:**
- Implement Dead Letter Queue (DLQ) for failed compensations
- Add retry mechanism with exponential backoff
- Create scheduled job to reconcile failed compensations
- Alert on-call team for manual intervention

---

### 2.4 Authorization Missing in Queue Operations

**File:** `booking-worker/grpc/booking_worker_service.go` (based on TODO comment)

```go
// TODO: Verify user owns the queue item before removing
```

**Impact:**
- Any authenticated user can cancel ANY other user's queue position
- Malicious users can sabotage competitors' bookings

**Fix:**
- Store `user_id` in queue item
- Verify `request.user_id == item.user_id` before cancel/update operations
- Add rate limiting per user for queue operations

---

## 3. HIGH Priority Issues (P1)

### 3.1 Lock Timeout vs Saga Duration Mismatch

**File:** `booking-service/.../BookingLockService.java:21-28`

```java
boolean acquired = lock.tryLock(lockTimeoutSeconds, TimeUnit.SECONDS);
// lockTimeoutSeconds = 5 (from config)
// BUT: Saga can take 10-30 seconds (gRPC calls + payment)
```

**Impact:**
- Lock expires while saga still executing
- Two concurrent sagas for same event
- Double booking of same seats

**Fix:**
```java
// Explicitly set lease time to match saga duration
boolean acquired = lock.tryLock(5, 60, TimeUnit.SECONDS); // Wait 5s, hold 60s
```

---

### 3.2 No Circuit Breaker Pattern

**File:** `booking-worker/internal/worker/processor.go:180-207`

```go
for attempt := 1; attempt <= p.config.Worker.MaxRetries; attempt++ {
    time.Sleep(p.config.Worker.RetryInterval) // Fixed 5s delay
    // No circuit breaker - keeps retrying even if service is down
}
```

**Impact:**
- If booking-service crashes, all workers blocked retrying
- Worker pool exhaustion
- Cascading failures across services

**Fix:**
- Implement circuit breaker (e.g., `sony/gobreaker` for Go)
- Add exponential backoff with jitter
- Add health check before processing

---

### 3.3 Queue Dequeue Fairness Issue

**File:** `booking-worker/internal/queue/redis_queue.go:117`

```go
// Uses FIRST available queue only - no round-robin
result, err := r.client.BRPop(ctx, timeout, keys[0]).Result()
```

**Impact:**
- Popular events starve less popular ones
- Users in less popular events wait longer than necessary
- Unfair processing distribution

**Fix:**
- Implement weighted round-robin across event queues
- Use Redis Streams with consumer groups for fair processing
- Or: Separate worker pools per event tier (hot/normal)

---

### 3.4 Missing Transactional Outbox

**File:** `booking-service/.../BookingSagaOrchestrator.java:106`

```java
booking = bookingRepository.save(booking);  // Step 1: DB transaction commits
eventPublisher.publishBookingConfirmed(booking);  // Step 2: Kafka publish (may fail)
```

**Impact:**
- If Kafka publish fails after DB commit, booking confirmed but downstream services never notified
- Email never sent, analytics never updated
- Inconsistent system state

**Fix:**
- Implement Outbox Pattern:
  1. Save event to `outbox` table in SAME transaction
  2. Background job polls outbox and publishes to Kafka
  3. Mark event as published after successful delivery
- Or: Use Debezium CDC for automatic outbox processing

---

## 4. MEDIUM Priority Issues (P2)

### 4.1 Exponential Backoff Missing in Booking Worker

```go
time.Sleep(p.config.Worker.RetryInterval) // Fixed 5s
```

**Fix:** Implement exponential backoff with jitter:
```go
delay := baseDelay * time.Duration(1<<attempt) + randomJitter
```

### 4.2 Remove Operation O(N) Complexity

**File:** `booking-worker/internal/queue/redis_queue.go:192-243`

```go
// LRANGE(0, -1) - reads ENTIRE queue into memory
items, err := r.client.LRange(ctx, key, 0, -1).Result()
for i, itemData := range items { // O(N) iteration
```

**Fix:**
- Store item ID → event ID mapping in separate hash
- Use LREM directly with serialized item data

### 4.3 No Connection Pooling Configuration

**File:** `booking-worker/grpcclient/booking_service_client.go`

Single connection per service, no pool configuration.

**Fix:**
- Configure gRPC connection pool
- Add load balancing configuration
- Implement connection health checks

### 4.4 Timeout Handler Inefficiency

**File:** `booking-worker/internal/queue/timeout_handler.go`

- Uses KEYS pattern (blocking)
- Processes in batches of 100
- LRANGE to find items (O(N))

**Fix:**
- Use Lua script for atomic timeout processing
- Or Redis Streams with XAUTOCLAIM for automatic timeout handling

---

## 5. Flow Logic Assessment

### 5.1 Booking Flow - Current State

```
Client → Gateway → Booking Worker (Queue)
                        ↓
              Wait for position
                        ↓
              Worker processes → Booking Service (Saga)
                                       ↓
                              ├─ Acquire Lock
                              ├─ Create Booking (PENDING)
                              ├─ Reserve Seats (gRPC → Ticket Service)
                              ├─ Process Payment (gRPC → Payment Service)
                              ├─ Confirm Booking
                              ├─ Publish Event (Kafka)
                              └─ Release Lock
                                       ↓
              Notify Client (Realtime Service)
```

**Good:**
- Saga pattern correctly orchestrates distributed transaction
- Distributed locking prevents concurrent modifications
- Compensation logic handles failures

**Missing:**
- [ ] Booking expiry handling (what if user doesn't complete payment?)
- [ ] Partial failure recovery (seats reserved but payment failed)
- [ ] Reconciliation job for orphaned bookings

### 5.2 Payment Flow - Current State

**Good:**
- Idempotency key properly implemented
- Transaction logging for audit trail
- Gateway adapter pattern for multiple providers
- Webhook handling for async confirmation

**Missing:**
- [ ] Payment timeout handling
- [ ] Refund flow not fully integrated with booking saga
- [ ] No fraud detection integration point

### 5.3 Email Flow - Current State

**Good:**
- Worker pool pattern
- Exponential backoff on retries
- Scheduled job support
- Database tracking of job status
- Multiple provider support (SES, SendGrid, SMTP)

**Complete and well-implemented.**

---

## 6. Performance Optimization Recommendations

### 6.1 For 100k Concurrent Users

| Component | Current | Recommended | Reason |
|-----------|---------|-------------|--------|
| Redis Queue | List + ZSet | Redis Streams | Built-in consumer groups, better scaling |
| Lock timeout | 5s | 60s | Match saga duration |
| Worker pool | Fixed | Auto-scaling | Handle load spikes |
| gRPC connections | 1 per service | Pool (10-50) | Better throughput |
| Dequeue strategy | First queue | Round-robin | Fair processing |

### 6.2 Database Optimization

```sql
-- Add indexes for common queries
CREATE INDEX idx_booking_user_status ON bookings(user_id, status);
CREATE INDEX idx_booking_event_created ON bookings(event_id, created_at);
CREATE INDEX idx_payment_booking ON payments(booking_id);

-- Add unique constraint for idempotency
ALTER TABLE bookings ADD CONSTRAINT uk_idempotency_key UNIQUE (idempotency_key);
```

### 6.3 Caching Strategy

```
Event details    → Cache 5 min (high read)
Seat availability → Cache 10s (frequent updates)
User session     → Cache 30 min
Booking status   → No cache (real-time accuracy)
```

---

## 7. Security Recommendations

| Issue | Severity | Fix |
|-------|----------|-----|
| Missing authorization in queue cancel | HIGH | Verify user ownership |
| No rate limiting per user | MEDIUM | Add Redis-based rate limiter |
| gRPC insecure credentials | MEDIUM | Use TLS in production |
| No input validation in worker | LOW | Validate queue item fields |

---

## 8. Summary: Action Items by Priority

### P0 - CRITICAL (Fix before production)
1. Replace Redis KEYS with SCAN or active queue set
2. Add idempotency key to booking worker → booking service calls
3. Implement DLQ and retry for failed compensations
4. Add authorization check for queue operations

### P1 - HIGH (Fix within 1-2 weeks)
5. Configure proper lock lease time
6. Implement circuit breaker pattern
7. Add fair queue processing (round-robin)
8. Implement transactional outbox pattern

### P2 - MEDIUM (Backlog)
9. Add exponential backoff with jitter
10. Optimize Remove operation complexity
11. Configure gRPC connection pooling
12. Optimize timeout handler with Lua scripts

---

## 9. Overall Assessment

| Aspect | Score | Notes |
|--------|-------|-------|
| Architecture | ⭐⭐⭐⭐ | Good use of patterns, clear separation |
| Code Quality | ⭐⭐⭐⭐ | Clean, well-structured |
| Scalability | ⭐⭐⭐ | Needs optimization for 100k users |
| Security | ⭐⭐⭐ | Missing authorization in some areas |
| Resilience | ⭐⭐⭐ | Missing circuit breaker, DLQ |
| Observability | ⭐⭐⭐⭐ | Good metrics, logging |

**Overall: 3.5/5 - Solid foundation, needs optimization for high-scale production**

---

## Files Reviewed

- `booking-service/src/main/java/com/ticketing/booking/service/saga/BookingSagaOrchestrator.java`
- `booking-worker/internal/queue/redis_queue.go`
- `booking-worker/internal/worker/processor.go`
- `payment-service/src/main/java/com/ticketing/payment/service/PaymentService.java`
- `email-worker/processor/processor.go`
- `shared-lib/protos/*.proto`

---

## 10. Scaling to 1M+ Concurrent Users

### 10.1 Infrastructure Changes Required

#### Redis Cluster (CRITICAL)
```yaml
# Current: Single Redis instance
# Required: Redis Cluster with 6+ nodes

redis-cluster:
  nodes: 6  # 3 masters + 3 replicas
  memory_per_node: 32GB
  max_connections_per_node: 65000

# Key sharding strategy
# booking-queue:{eventId} -> hash slot based on eventId
# This distributes load across cluster nodes
```

#### Database Sharding
```sql
-- Current: Single PostgreSQL with replicas
-- Required: Horizontal sharding by event_id

-- Shard 0: events 0-999999
-- Shard 1: events 1000000-1999999
-- etc.

-- Use Citus or manual sharding with PgBouncer
```

#### Kafka Partitioning
```yaml
# Current: Default partitions
# Required: 50+ partitions per topic

booking-events:
  partitions: 50
  replication_factor: 3

# Partition key: event_id
# This ensures all events for same event go to same partition
# Maintains ordering within event
```

### 10.2 Application Architecture Changes

#### 1. Replace Redis List with Redis Streams
```go
// Current: Redis List (LPUSH/RPOP)
// Problem: No consumer groups, no automatic reprocessing

// Required: Redis Streams with Consumer Groups
// Benefits:
// - Built-in consumer groups for load distribution
// - Automatic message acknowledgment
// - Dead letter handling with XAUTOCLAIM
// - Better scaling characteristics

// Example:
XADD booking-stream:{eventId} * user_id "123" seats "A1,A2"
XREADGROUP GROUP workers worker-1 COUNT 10 STREAMS booking-stream:* >
XACK booking-stream:{eventId} workers {messageId}
```

#### 2. Implement Event-Driven Architecture with Kafka Streams
```java
// Current: Synchronous gRPC calls
// Problem: Tight coupling, cascading failures

// Required: Event sourcing for critical operations
// - Booking requested -> Kafka topic
// - Seats reserved -> Kafka topic
// - Payment processed -> Kafka topic
// - Booking confirmed -> Kafka topic

// Benefits:
// - Decoupled services
// - Natural backpressure
// - Replay capability for recovery
```

#### 3. Add API Rate Limiting at Gateway
```javascript
// Per-user rate limits
const rateLimits = {
  'booking.create': { requests: 5, window: '1m' },
  'booking.cancel': { requests: 10, window: '1m' },
  'queue.position': { requests: 60, window: '1m' },
};

// Per-event rate limits (prevent single event from overwhelming system)
const eventRateLimits = {
  'hot_events': { requests: 10000, window: '1s' },
  'normal_events': { requests: 1000, window: '1s' },
};
```

#### 4. Implement Read Replicas with Automatic Failover
```yaml
# Use PgBouncer with read/write splitting
pgbouncer:
  pools:
    - name: write
      host: postgres-master
      mode: transaction
    - name: read
      hosts:
        - postgres-replica-1
        - postgres-replica-2
        - postgres-replica-3
      mode: session
      load_balance: round_robin
```

### 10.3 Caching Strategy for 1M+ Users

```
Layer 1: CDN (CloudFront/Cloudflare)
├── Static assets
├── Event details (5 min TTL)
└── Seat maps (1 min TTL)

Layer 2: Redis Cache
├── User sessions (30 min TTL)
├── Event metadata (5 min TTL)
├── Seat availability (10 sec TTL - short for accuracy)
└── Queue positions (no cache - real-time)

Layer 3: Local Cache (in-process)
├── Configuration
├── Rate limit counters
└── Circuit breaker state
```

### 10.4 Worker Pool Scaling

```yaml
# Current: Fixed worker pool
# Required: Auto-scaling based on queue depth

booking-worker:
  min_replicas: 10
  max_replicas: 100

  scaling_rules:
    - metric: queue_depth
      threshold: 1000
      scale_up_by: 5

    - metric: processing_latency_p99
      threshold: 5s
      scale_up_by: 3

    - metric: circuit_breaker_open
      action: scale_down  # Don't scale if downstream is failing
```

### 10.5 Database Optimization

```sql
-- Partitioning bookings table by created_at
CREATE TABLE bookings (
    id BIGSERIAL,
    booking_id UUID NOT NULL,
    event_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    -- ... other columns
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE bookings_2025_01 PARTITION OF bookings
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- Indexes for hot queries
CREATE INDEX CONCURRENTLY idx_booking_event_status
    ON bookings (event_id, status)
    WHERE status IN ('PENDING', 'CONFIRMED');

-- Partial index for active bookings only
CREATE INDEX CONCURRENTLY idx_booking_active
    ON bookings (user_id, created_at DESC)
    WHERE status NOT IN ('CANCELLED', 'FAILED');
```

### 10.6 Monitoring & Alerting for Scale

```yaml
# Key metrics to monitor at 1M scale

metrics:
  - name: booking_queue_depth
    alert_threshold: 50000
    action: scale_workers

  - name: booking_latency_p99
    alert_threshold: 10s
    action: page_oncall

  - name: circuit_breaker_open_duration
    alert_threshold: 5m
    action: page_oncall

  - name: redis_memory_usage_percent
    alert_threshold: 80%
    action: scale_redis

  - name: db_connection_pool_exhausted
    alert_threshold: 1
    action: page_oncall

  - name: kafka_consumer_lag
    alert_threshold: 10000
    action: scale_consumers
```

### 10.7 Cost Estimation (AWS)

| Component | 100K Users | 1M Users | Notes |
|-----------|------------|----------|-------|
| EKS Cluster | $500/mo | $3000/mo | Auto-scaling nodes |
| RDS PostgreSQL | $800/mo | $5000/mo | Multi-AZ, read replicas |
| ElastiCache Redis | $300/mo | $2000/mo | Cluster mode |
| MSK (Kafka) | $400/mo | $2000/mo | 50+ partitions |
| ALB | $200/mo | $500/mo | Cross-AZ |
| CloudWatch | $100/mo | $500/mo | Custom metrics |
| **Total** | **~$2300/mo** | **~$13000/mo** | |

### 10.8 Migration Path

**Phase 1: Optimize Current Architecture (Weeks 1-2)**
- ✅ Replace Redis KEYS with SET tracking
- ✅ Add idempotency keys
- ✅ Implement circuit breaker
- ✅ Fix lock lease times
- Add proper indexes

**Phase 2: Scale Infrastructure (Weeks 3-4)**
- Deploy Redis Cluster
- Add read replicas
- Increase Kafka partitions
- Implement PgBouncer

**Phase 3: Architecture Evolution (Weeks 5-8)**
- Migrate to Redis Streams
- Implement event sourcing
- Add distributed tracing
- Implement auto-scaling

**Phase 4: Testing & Optimization (Weeks 9-12)**
- Load testing at 1M concurrent
- Chaos engineering
- Performance tuning
- Documentation

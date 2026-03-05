# Redis Key Schema & TTL Enforcement

Documentation về Redis key schema và TTL enforcement strategy cho Booking Service.

---

## 📋 Overview

Booking Service sử dụng Redis (Redisson) cho:
- **Distributed Locking**: Đảm bảo chỉ một booking request được xử lý tại một thời điểm cho mỗi event
- **Concurrency Control**: Ngăn chặn race conditions khi nhiều users cùng book seats cho cùng một event

---

## 🔑 Redis Key Schema

### Lock Keys

#### Pattern: `booking:lock:{eventId}`

**Mục đích**: Distributed lock để serialize booking requests cho cùng một event.

**Format**:
```
booking:lock:{eventId}
```

**Ví dụ**:
```
booking:lock:evt-12345
booking:lock:evt-67890
```

**Key Components**:
- Prefix: `booking:lock:` - Namespace cho booking locks
- `{eventId}`: Event identifier (string) - Unique identifier của event

**Implementation**:
```java
// BookingLockService.java
RLock lock = redissonClient.getLock("booking:lock:" + eventId);
boolean acquired = lock.tryLock(lockTimeoutSeconds, TimeUnit.SECONDS);
```

**Characteristics**:
- **Type**: Redisson RLock (distributed lock)
- **TTL**: Auto-managed by Redisson watchdog (default 30s, auto-renewed)
- **Acquisition Timeout**: Configurable via `booking.lock.timeout-seconds` (default: 5s)
- **Release**: Automatic on unlock() or thread termination

---

## ⏱️ TTL Enforcement Strategy

### 1. Distributed Lock TTL

#### Configuration

```yaml
booking:
  lock:
    ttl-seconds: ${BOOKING_LOCK_TTL_SECONDS:120}  # Default: 120 seconds
```

**Lock Acquisition**:
- **Try Lock Timeout**: `booking.lock.timeout-seconds` (default: 5s)
  - Thời gian chờ tối đa để acquire lock
  - Nếu không acquire được trong thời gian này → throw `BookingLockException`

**Lock Lease Time**:
- Redisson sử dụng **watchdog mechanism** để auto-renew locks
- Default lease time: **30 seconds** (Redisson default)
- Auto-renewed nếu lock vẫn được hold bởi thread
- Lock sẽ expire nếu thread bị terminate hoặc crash

**Lock Release**:
- **Explicit**: `lock.unlock()` được gọi trong `finally` block
- **Automatic**: Lock auto-expires sau 30s nếu thread bị terminate
- **Safety Check**: `lock.isHeldByCurrentThread()` đảm bảo chỉ thread owner mới unlock

#### TTL Enforcement Flow

```
1. Thread A acquires lock: booking:lock:evt-123
   → Redisson sets lease time: 30s
   → Watchdog starts auto-renewal

2. Thread A holds lock during saga execution
   → Watchdog renews lock every 10s (30s / 3)
   → Lock remains active

3. Thread A completes saga
   → lock.unlock() called in finally block
   → Lock released immediately

4. If Thread A crashes/terminates
   → Watchdog stops
   → Lock expires after 30s
   → Thread B can acquire lock
```

### 2. Booking Expiry TTL

#### Configuration

```yaml
booking:
  expiry:
    minutes: ${BOOKING_EXPIRY_MINUTES:15}  # Default: 15 minutes
```

**Booking Expiry**:
- **Database Field**: `expires_at` (TIMESTAMP WITH TIME ZONE)
- **Default TTL**: 15 minutes từ thời điểm booking được tạo
- **Purpose**: 
  - Pending bookings expire sau 15 phút nếu không được confirmed
  - Released seats back to inventory khi booking expires

**Implementation**:
```java
// ReferenceGenerator.java
public static OffsetDateTime bookingExpiry() {
    return OffsetDateTime.now().plus(15, ChronoUnit.MINUTES);
}

// Booking entity
@Column(name = "expires_at", nullable = false)
private OffsetDateTime expiresAt;
```

**Expiry Enforcement**:
- **Application Level**: Booking status checked against `expiresAt` before operations
- **Database Level**: `expires_at` column stored in PostgreSQL
- **Background Job** (Future): Scheduled job để cleanup expired bookings

---

## 🔒 Lock Management

### Lock Acquisition

```java
// BookingSagaOrchestrator.java
private RLock acquireLock(String eventId) {
    try {
        RLock lock = lockService.acquireLock(eventId);
        if (lock == null) {
            throw new BookingLockException("Failed to acquire lock for event: " + eventId);
        }
        return lock;
    } catch (InterruptedException e) {
        Thread.currentThread().interrupt();
        throw new BookingLockException("Interrupted while acquiring lock", e);
    }
}
```

**Lock Acquisition Flow**:
1. Call `lockService.acquireLock(eventId)`
2. Redisson creates lock key: `booking:lock:{eventId}`
3. Try acquire với timeout: `tryLock(timeoutSeconds, TimeUnit.SECONDS)`
4. Nếu acquired → return RLock
5. Nếu timeout → throw `BookingLockException`

### Lock Release

```java
// BookingSagaOrchestrator.java
finally {
    releaseLock(lock);
}

// BookingLockService.java
public void releaseLock(RLock lock) {
    if (lock != null && lock.isHeldByCurrentThread()) {
        lock.unlock();
    }
}
```

**Lock Release Flow**:
1. Check lock is not null
2. Check lock is held by current thread (safety check)
3. Call `lock.unlock()`
4. Lock released immediately

**Safety Mechanisms**:
- `isHeldByCurrentThread()`: Prevents unlocking locks held by other threads
- `finally` block: Ensures lock is always released even on exceptions

---

## 📊 Key Naming Conventions

### Pattern Rules

1. **Namespace Prefix**: Tất cả keys bắt đầu với `booking:`
2. **Resource Type**: Sau prefix là resource type (`lock:`, `cache:`, etc.)
3. **Identifier**: Cuối cùng là unique identifier (`{eventId}`, `{bookingId}`, etc.)

### Examples

```
✅ Good:
booking:lock:evt-12345
booking:cache:booking:uuid-12345

❌ Bad:
lock:evt-12345              # Missing namespace
booking:evt-12345           # Missing resource type
booking:lock:                # Missing identifier
```

---

## 🛡️ Best Practices

### 1. Lock Scope

- **Granularity**: Lock per event, not per booking
  - ✅ `booking:lock:evt-12345` (one lock per event)
  - ❌ `booking:lock:booking-uuid-123` (too granular, causes deadlocks)

### 2. Lock Timeout

- **Acquisition Timeout**: 5 seconds (configurable)
  - Đủ thời gian để wait nếu lock đang được hold
  - Không quá dài để tránh blocking requests

- **Lock Lease Time**: 30 seconds (Redisson default)
  - Auto-renewed bởi watchdog
  - Expires nếu thread crashes

### 3. Lock Release

- **Always release in finally block**:
  ```java
  RLock lock = null;
  try {
      lock = acquireLock(eventId);
      // ... saga execution ...
  } finally {
      releaseLock(lock);  // Always executed
  }
  ```

- **Check ownership before unlock**:
  ```java
  if (lock != null && lock.isHeldByCurrentThread()) {
      lock.unlock();
  }
  ```

### 4. Error Handling

- **Lock Acquisition Failure**: Throw `BookingLockException`
- **Interrupted Exception**: Restore interrupt status
- **Lock Release Failure**: Log error, don't throw (lock will auto-expire)

### 5. Monitoring

- **Metrics**: Track lock acquisition failures
  ```java
  metricsService.recordSagaStep("acquire_lock", "failed");
  ```

- **Logging**: Log lock operations for debugging
  ```java
  log.info("Acquired lock for event: {}", eventId);
  log.warn("Failed to acquire lock for event: {}", eventId);
  ```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. Lock Not Released

**Symptoms**: Lock remains held after saga completion

**Causes**:
- Exception thrown before `finally` block
- Thread terminated unexpectedly
- Network partition between app and Redis

**Solutions**:
- Lock auto-expires after 30s (Redisson watchdog)
- Check `finally` block is always executed
- Monitor lock TTL in Redis

#### 2. Lock Acquisition Timeout

**Symptoms**: `BookingLockException: Unable to acquire booking lock`

**Causes**:
- High concurrency for same event
- Previous lock not released
- Redis connection issues

**Solutions**:
- Increase `booking.lock.timeout-seconds` (not recommended)
- Check lock release logic
- Monitor Redis connection health

#### 3. Deadlock

**Symptoms**: Multiple threads waiting for locks indefinitely

**Causes**:
- Circular lock dependencies
- Locking multiple resources in different order

**Solutions**:
- Always lock in same order (eventId sorted)
- Use single lock per event (current implementation)
- Set appropriate timeout

---

## 📈 Performance Considerations

### Lock Contention

**High Concurrency Scenario**:
- Nhiều users cùng book seats cho cùng event
- Lock serializes requests → sequential processing
- Trade-off: Consistency vs. Throughput

**Optimization Strategies**:
1. **Lock Granularity**: Lock per event (current) vs. per seat (finer granularity)
2. **Lock Timeout**: Balance between wait time and failure rate
3. **Retry Logic**: Retry lock acquisition with backoff

### Redis Performance

**Key Operations**:
- `SETNX` (lock acquisition): O(1)
- `EXPIRE` (TTL setting): O(1)
- `DEL` (lock release): O(1)

**Network Overhead**:
- Each lock operation = 1 Redis round-trip
- Watchdog renewal = periodic Redis calls (every 10s)

---

## 🔐 Security Considerations

### Key Isolation

- **Namespace**: `booking:` prefix isolates booking keys
- **Database**: Use separate Redis database (configurable via `redis.database`)
- **Password**: Redis password protection (configurable via `redis.password`)

### Access Control

- **Lock Ownership**: Only lock owner can unlock
- **Thread Safety**: `isHeldByCurrentThread()` check prevents unauthorized unlock

---

## 📝 Configuration Reference

### application.yml

```yaml
# Redis Configuration
redis:
  host: ${REDIS_HOST:localhost}
  port: ${REDIS_PORT:50379}
  password: ${REDIS_PASSWORD:}
  database: ${REDIS_DATABASE:0}

# Booking Configuration
booking:
  lock:
    timeout-seconds: ${BOOKING_LOCK_TIMEOUT_SECONDS:5}    # Lock acquisition timeout
    ttl-seconds: ${BOOKING_LOCK_TTL_SECONDS:120}          # Lock TTL (Redisson uses 30s default)
  expiry:
    minutes: ${BOOKING_EXPIRY_MINUTES:15}                  # Booking expiry TTL
```

### Environment Variables

```bash
# Redis
REDIS_HOST=localhost
REDIS_PORT=50379
REDIS_PASSWORD=
REDIS_DATABASE=0

# Booking Locks
BOOKING_LOCK_TIMEOUT_SECONDS=5
BOOKING_LOCK_TTL_SECONDS=120

# Booking Expiry
BOOKING_EXPIRY_MINUTES=15
```

---

## 📚 Related Documentation

- [Booking Service README](./README.md)
- [Saga Orchestrator Implementation](./IMPLEMENTATION_SUMMARY.md)
- [Redisson Documentation](https://github.com/redisson/redisson/wiki/8.-Distributed-locks-and-synchronizers)

---

**Last Updated**: 2024  
**Status**: ✅ Complete



# ✅ Phase 1 - Queue Infrastructure Setup (Planned)

Documentation placeholder for booking-worker queue implementation.

## 🎯 Scope
- Redis queue implementation (primary)
- Kafka queue implementation (optional, for high throughput)
- Queue manager interface + implementations
- Queue position tracking (Redis sorted sets)
- Timeout handling (TTL + cleanup workers)

## 📌 Planned Tasks
- [ ] Design queue interface (`QueueManager` interface)
- [ ] Implement Redis queue manager:
  - Enqueue operations (LPUSH to `booking-queue`)
  - Dequeue operations (BRPOP with blocking)
  - Position tracking (ZADD/ZRANK for sorted sets)
  - Status tracking (HSET for `booking-status:{clientId}`)
- [ ] Implement Kafka queue manager (optional):
  - Producer for enqueue
  - Consumer groups for dequeue
  - Offset management
- [ ] Add timeout handler:
  - Redis TTL on queue entries
  - Sorted set for timeout tracking (`booking-timeouts`)
  - Cleanup worker (goroutine)
- [ ] Add distributed locking (go-redsync) for concurrent dequeue

## 🧪 Testing Strategy
- TODO: Unit tests for queue operations
- TODO: Integration tests with Testcontainers (Redis)
- TODO: Load tests for 100k+ concurrent enqueue

## 📁 File Map
```
docs/booking-worker/02_QUEUE_SETUP_COMPLETE.md
booking-worker/internal/queue/
├── interface.go          # QueueManager interface
├── redis_queue.go        # Redis implementation
├── kafka_queue.go        # Kafka implementation (optional)
└── timeout_handler.go    # Timeout cleanup logic
```

## 🔑 Redis Key Schema
- `booking-queue:{eventId}`: Redis List for queue entries
- `booking-timeouts:{eventId}`: Sorted Set (score = expiry timestamp)
- `booking-status:{clientId}`: Hash for client status
- `booking-lock:{eventId}`: Distributed lock key

_Last updated: Planning stage (2024)_


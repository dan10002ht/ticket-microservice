# 🏗️ Booking Worker Architecture (Planning)

Reference document for booking-worker component design, goroutine patterns, and integration flows.

## 🧬 Core Components
- **QueueManager**: Interface for queue operations (Redis/Kafka)
- **WorkerProcessor**: Goroutine pool for processing queue items
- **TimeoutHandler**: Cleanup expired queue entries
- **gRPCClient**: Clients for booking-service, realtime-service
- **MetricsExporter**: Prometheus metrics collection

## 🔄 Goroutine Patterns

### Worker Pool Pattern
```go
// Multiple workers process queue concurrently
for i := 0; i < workerCount; i++ {
    go func() {
        for {
            item := queue.Dequeue()
            process(item)
        }
    }()
}
```

### Channel-based Communication
- Use channels for worker coordination
- Context for cancellation/timeouts
- Select statements for non-blocking operations

## 🔁 Processing Flow

1. **Enqueue**: Gateway → booking-worker → Redis/Kafka
2. **Position Tracking**: Update sorted set, notify client
3. **Dequeue**: Worker goroutine pulls from queue
4. **Processing**: Call booking-service gRPC
5. **Notification**: Send result via realtime-service
6. **Timeout**: Cleanup if payment not completed

## 🔒 Distributed Locking

- Use go-redsync (Redlock algorithm)
- Lock key: `booking-lock:{eventId}`
- TTL: Match booking timeout (120s default)
- Prevents concurrent processing of same event

## 📊 Metrics & Observability

- Queue length (Gauge)
- Processing rate (Counter)
- Processing duration (Histogram)
- Error rate (Counter)
- Timeout count (Counter)

## 🧱 TODOs
- [ ] Design worker pool size calculation (based on queue length)
- [ ] Document error handling strategy (retry, dead letter queue)
- [ ] Define backpressure mechanism (reject when queue full)
- [ ] Add sequence diagrams (queue → worker → booking-service)
- [ ] Document scaling strategy (horizontal scaling with consumer groups)

_Last updated: Planning stage (2024)_


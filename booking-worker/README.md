# Booking Worker Service (Go)

**Language:** Go (Golang)

**Why Go?**

- Handles 100k+ concurrent booking requests efficiently
- Goroutines for lightweight concurrency
- High throughput, low latency, easy to scale horizontally

## Overview

The **Booking Worker Service** is designed to handle high-concurrency booking requests by managing a distributed queue for clients attempting to book tickets. With 5,000 tickets and over 100,000 clients, this service ensures fairness, prevents overselling, and provides a smooth user experience by placing clients in a queue and processing their booking requests in order.

## 🎯 Responsibilities

- **Distributed Queue Management**: Place and manage client booking requests in a distributed queue (Redis/Kafka).
- **Fairness & Concurrency Control**: Ensure first-come, first-served processing and prevent overselling.
- **Queue Notifications**: Notify clients of their queue position and status (via WebSocket or polling).
- **Timeout Handling**: Remove clients who do not complete booking in time.
- **Scalability**: Scale horizontally to handle spikes in demand.
- **gRPC Integration**: Communicate with Booking Service, Ticket Service, Realtime Service, and Gateway via gRPC.

## 🏗️ Architecture

### Technology Stack

- **Language**: Go 1.21+
- **Framework**: Go kit / Fiber / Gin (suggested)
- **Queue**: Redis (go-redis) or Kafka (segmentio/kafka-go)
- **gRPC**: google.golang.org/grpc
- **WebSocket**: gorilla/websocket
- **Monitoring**: Prometheus client_golang, Grafana
- **Concurrency**: Goroutines, channels
- **Distributed Locking**: Redlock (go-redsync)

### Key Components

```
Booking Worker Service (Go)
├── Queue Manager (Redis/Kafka)
├── Worker Processor (goroutines)
├── Timeout Handler
├── gRPC Server/Client (protobuf)
├── Realtime Notifier (WebSocket)
├── Metrics Exporter (Prometheus)
└── Health Checker
```

## 🔄 Booking Queue Flow

1. **Client Request**: Client requests to book a ticket via Gateway.
2. **Queue Placement**: Client is placed in a distributed queue (Redis List or Kafka Topic).
3. **Queue Notification**: Client receives their queue position and estimated wait time (via Realtime Service).
4. **Worker Processing**: Booking Worker dequeues requests and processes them one by one:
   - Checks ticket availability (via Ticket Service)
   - Reserves ticket (via Booking Service)
   - Notifies client to proceed with payment
5. **Timeout/Failure Handling**: If client does not complete payment in time, ticket is released and next client is notified.
6. **Completion**: Successful bookings are confirmed and tickets are issued.

## 📡 gRPC APIs

### BookingWorkerService

```protobuf
service BookingWorkerService {
  rpc EnqueueBooking(EnqueueBookingRequest) returns (EnqueueBookingResponse);
  rpc GetQueuePosition(GetQueuePositionRequest) returns (GetQueuePositionResponse);
  rpc GetQueueStatus(GetQueueStatusRequest) returns (GetQueueStatusResponse);
  rpc CancelQueueItem(CancelQueueItemRequest) returns (CancelQueueItemResponse);
  rpc Health(HealthRequest) returns (HealthResponse);
}
```

See `shared-lib/protos/booking_worker.proto` for full message definitions.

### Integration Points

- **Booking Service**: Reserve and confirm bookings
- **Ticket Service**: Check and update ticket inventory
- **Realtime Service**: Send queue position and status updates to clients
- **Gateway**: Route booking requests and queue status APIs

## ⚙️ Configuration

### Environment Variables

See `env.example` for all available configuration options. Key variables:

```bash
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# gRPC Configuration
GRPC_PORT=50059
GRPC_BOOKING_SERVICE_ENDPOINT=localhost:50058
GRPC_REALTIME_SERVICE_ENDPOINT=localhost:50060

# Queue Configuration
QUEUE_MAX_SIZE=100000
QUEUE_TIMEOUT_SECONDS=900
QUEUE_CLEANUP_INTERVAL=30s

# Worker Configuration
WORKER_POOL_SIZE=10
WORKER_MAX_RETRIES=3
WORKER_RETRY_INTERVAL=5s

# Metrics Configuration
METRICS_ENABLED=true
METRICS_PORT=9091

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

Copy `env.example` to `.env` and update with your configuration.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Redis (for queue backend)
- protoc (Protocol Buffers compiler)
- protoc-gen-go and protoc-gen-go-grpc plugins

### Installation

1. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

2. **Generate protobuf code:**
   ```bash
   chmod +x scripts/generate-proto.sh
   ./scripts/generate-proto.sh
   ```

3. **Setup environment:**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Build the service:**
   ```bash
   chmod +x scripts/build.sh
   ./scripts/build.sh
   ```

5. **Run the service:**
   ```bash
   # Development mode
   chmod +x scripts/dev.sh
   ./scripts/dev.sh
   
   # Or run directly
   go run main.go
   
   # Or run built binary
   ./booking-worker
   ```

### Project Structure

```
booking-worker/
├── main.go                    # Application entry point
├── go.mod                      # Go module definition
├── config/                     # Configuration management
│   └── config.go
├── internal/
│   ├── app/                    # Application initialization
│   │   └── app.go
│   ├── queue/                  # Queue management
│   │   ├── interface.go       # Queue interface
│   │   ├── redis_queue.go      # Redis implementation
│   │   └── timeout_handler.go  # Timeout cleanup
│   └── worker/                 # Worker pool
│       └── processor.go
├── grpc/                       # gRPC server handlers
│   └── booking_worker_service.go
├── grpcclient/                 # gRPC clients
│   ├── booking_service_client.go
│   └── realtime_service_client.go
├── metrics/                     # Prometheus metrics
│   └── metrics.go
├── scripts/                     # Utility scripts
│   ├── generate-proto.sh
│   ├── build.sh
│   └── dev.sh
├── Dockerfile
├── env.example
└── README.md
```

## 🚀 Performance & Scalability

- **Horizontal Scaling**: Multiple worker instances can process the queue in parallel (with distributed locking or consumer groups).
- **Atomic Operations**: Use Redis Lua scripts or Kafka consumer groups for atomic dequeue and processing.
- **Backpressure**: If queue is full, new requests are rejected or delayed.
- **Metrics**: Expose queue length, processing rate, and timeout metrics for monitoring.

## 🔐 Security

- **gRPC Authentication**: Use mTLS or JWT for secure inter-service calls.
- **Rate Limiting**: Integrate with Rate Limiter Service to prevent abuse.
- **Audit Logging**: Log all queue and booking operations.

## 🧪 Testing

- **Unit Tests**: Test queue logic, timeout handling, and gRPC endpoints.
- **Integration Tests**: Simulate high concurrency booking scenarios.
- **Load Tests**: Validate performance under 100,000+ concurrent clients.

## 🛠️ Redis Queue Schema

The service uses the following Redis key patterns:

- `booking-queue:{eventId}`: Redis List for incoming booking requests per event
- `booking-queue-positions:{eventId}`: Sorted Set for tracking queue positions
- `booking-timeouts:{eventId}`: Sorted Set for tracking item expiry (score = expiry timestamp)
- `booking-lock:{eventId}`: Distributed lock key (via Redlock)

### Queue Item Structure

Each queue item contains:
- `ID`: Unique queue item identifier
- `EventID`: Event identifier
- `UserID`: User identifier
- `SeatNumbers`: List of requested seat numbers
- `SeatCount`: Number of seats
- `TotalAmount`: Total booking amount
- `Currency`: Currency code
- `Metadata`: Additional metadata (map)
- `EnqueuedAt`: Timestamp when enqueued
- `ExpiresAt`: Timestamp when item expires

## 📊 Metrics

The service exposes Prometheus metrics at `:9091/metrics`:

- `booking_worker_queue_length`: Current length of the booking queue (Gauge)
- `booking_worker_items_processed_total`: Total number of queue items processed (Counter)
- `booking_worker_processing_duration_seconds`: Duration of queue item processing (Histogram)
- `booking_worker_errors_total`: Total number of processing errors (Counter)

## 📞 Troubleshooting

- **Queue Stuck**: Check Redis health and worker logs
- **Overselling**: Ensure atomic reservation and ticket checks
- **Slow Processing**: Scale out workers, check for bottlenecks
- **Protobuf Errors**: Run `scripts/generate-proto.sh` to regenerate protobuf code
- **Connection Errors**: Verify gRPC service endpoints are accessible

## 🔗 Related Services

- **Booking Service**: Reservation logic
- **Ticket Service**: Inventory management
- **Realtime Service**: Client notifications
- **Gateway**: API routing

## 🧑‍💻 Implementation Status

### ✅ Completed

- [x] Go module setup and dependencies
- [x] Configuration management
- [x] Redis queue implementation (Enqueue, Dequeue, GetPosition, Remove)
- [x] Worker pool with goroutines
- [x] Timeout handler for expired items
- [x] gRPC client stubs (Booking Service, Realtime Service)
- [x] gRPC service handler (BookingWorkerService)
- [x] Prometheus metrics exporter
- [x] Protobuf definitions (booking_worker.proto)
- [x] Application structure with graceful shutdown
- [x] Scripts for proto generation, build, and dev

### ⏳ Pending (Requires Protobuf Generation)

- [ ] Generate Go code from protobuf files
- [ ] Complete gRPC service handler implementation
- [ ] Complete gRPC client implementations
- [ ] Integration with Booking Service
- [ ] Integration with Realtime Service
- [ ] Unit tests
- [ ] Integration tests
- [ ] Load testing

### 📝 Next Steps

1. **Generate Protobuf Code:**
   ```bash
   ./scripts/generate-proto.sh
   ```

2. **Update gRPC Handlers:**
   - Uncomment implementation code in `grpc/booking_worker_service.go`
   - Update method signatures to use generated protobuf types

3. **Update gRPC Clients:**
   - Implement methods in `grpcclient/booking_service_client.go`
   - Implement methods in `grpcclient/realtime_service_client.go`

4. **Testing:**
   - Write unit tests for queue operations
   - Write integration tests with Redis
   - Write gRPC handler tests

## 🧑‍💻 Go Implementation Best Practices

- Use goroutines for concurrent queue processing
- Use channels for internal communication
- Use context for cancellation/timeouts
- Use structured logging (zap)
- Use dependency injection for testability
- Use go modules for dependency management
- Write unit and integration tests (testing, testify)

## 📚 Related Documentation

- [Booking Service Docs](../docs/booking-service/README.md)
- [Booking Worker Docs](../docs/booking-worker/README.md)
- [Architecture Overview](../docs/architecture/AI_README.md)

---

# Realtime Service (Go + Gorilla WebSocket)

**Language:** Go

**Why Go + Gorilla WebSocket?**
- Performance: Go handles 100k+ concurrent connections per instance
- Memory efficiency: ~30MB per 10k connections (vs ~100MB for Node.js)
- Native concurrency: Goroutines are ideal for WebSocket connections
- Consistency: Matches existing Go workers (booking-worker, email-worker)
- Production-proven: Gorilla WebSocket is widely used in production systems

## Overview

The Realtime Service handles real-time notifications and WebSocket broadcasting for the ticketing system. It enables instant updates for booking status, payment confirmations, ticket availability, and queue positions.

## Responsibilities

- **WebSocket Connections**: Manage client WebSocket connections with authentication
- **Real-time Notifications**: Push booking, payment, and ticket updates to clients
- **Room Management**: Group clients by event for targeted broadcasts
- **Redis Pub/Sub**: Horizontal scaling across multiple instances
- **gRPC API**: Internal service communication for sending notifications

## Architecture

### Technology Stack

- **Runtime**: Go 1.22+
- **WebSocket**: Gorilla WebSocket
- **Message Broker**: Redis Pub/Sub
- **Internal API**: gRPC
- **Monitoring**: Prometheus + Grafana

### Key Components

```
realtime-service/
├── config/                  # Configuration management
├── internal/
│   ├── grpc/               # gRPC server and handlers
│   │   └── handlers/       # Notification handlers
│   ├── websocket/          # WebSocket components
│   │   ├── hub.go          # Connection hub (broadcast)
│   │   ├── client.go       # Client connection handler
│   │   ├── server.go       # HTTP upgrade handler
│   │   └── message.go      # Message types
│   ├── pubsub/             # Redis Pub/Sub
│   │   ├── subscriber.go   # Event subscriber
│   │   └── publisher.go    # Event publisher
│   ├── service/            # Business logic
│   └── middleware/         # Auth middleware
├── pkg/logger/             # Zap logger
├── metrics/                # Prometheus metrics
├── main.go                 # Main entry point
├── go.mod
├── Dockerfile
└── env.example
```

## API Endpoints

### WebSocket Endpoint

```
ws://localhost:3003/ws?token=<JWT>
```

### HTTP Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /ws` | WebSocket upgrade endpoint |
| `GET /health` | Basic health check |
| `GET /ready` | Readiness check (includes Redis) |
| `GET /metrics` | Prometheus metrics |

### gRPC Methods (Internal)

| Method | Description |
|--------|-------------|
| `NotifyBookingResult` | Push booking result to user |
| `NotifyQueuePosition` | Update queue position |
| `NotifyPaymentStatus` | Push payment status update |
| `BroadcastEvent` | Broadcast to room |
| `SendToUser` | Send to specific user |
| `GetConnectionStats` | Get connection statistics |

## WebSocket Message Protocol

### Message Format

```json
{
  "type": "booking:confirmed",
  "payload": { ... }
}
```

### Event Types

| Type | Direction | Description |
|------|-----------|-------------|
| `system:connected` | Server->Client | Connection established |
| `system:ping` | Client->Server | Keep-alive ping |
| `system:pong` | Server->Client | Ping response |
| `room:join` | Client->Server | Join event room |
| `room:leave` | Client->Server | Leave event room |
| `booking:queue_position` | Server->Client | Queue position update |
| `booking:processing` | Server->Client | Booking being processed |
| `booking:confirmed` | Server->Client | Booking confirmed |
| `booking:failed` | Server->Client | Booking failed |
| `payment:processing` | Server->Client | Payment processing |
| `payment:success` | Server->Client | Payment successful |
| `payment:failed` | Server->Client | Payment failed |
| `ticket:availability` | Server->Client | Ticket availability changed |

### Example Messages

**Server -> Client (Booking Confirmed)**
```json
{
  "type": "booking:confirmed",
  "payload": {
    "booking_id": "uuid-here",
    "booking_reference": "BK12345",
    "event_id": "event-uuid",
    "seat_numbers": ["A1", "A2"],
    "total_amount": "100.00",
    "currency": "USD"
  }
}
```

**Client -> Server (Join Room)**
```json
{
  "type": "room:join",
  "payload": {
    "room": "event:event-uuid"
  }
}
```

## Configuration

### Environment Variables

```bash
# HTTP Configuration
HTTP_PORT=3003
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=15s

# gRPC Configuration
GRPC_PORT=50057

# WebSocket Configuration
WS_PING_INTERVAL=30s
WS_PONG_TIMEOUT=60s
WS_ALLOW_ANONYMOUS=true

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=100

# JWT Configuration
JWT_SECRET=your-jwt-secret

# Metrics Configuration
METRICS_ENABLED=true
METRICS_PORT=9057

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

## Development

### Prerequisites

- Go 1.22+
- Redis 7+
- protoc (Protocol Buffers compiler)
- protoc-gen-go
- protoc-gen-go-grpc

### Setup

1. Clone the repository
2. Copy environment file:
   ```bash
   cp env.example .env
   ```
3. Start Redis:
   ```bash
   docker run -d -p 6379:6379 redis:7-alpine
   ```
4. Generate protobuf files:
   ```bash
   ./scripts/generate-proto.sh
   ```
5. Run the service:
   ```bash
   go run main.go
   ```

### Building

```bash
go build -o bin/realtime-service .
```

### Testing WebSocket

```bash
# Using wscat
npm install -g wscat
wscat -c "ws://localhost:3003/ws?token=<JWT>"

# Send ping
> {"type":"system:ping"}

# Join room
> {"type":"room:join","payload":{"room":"event:123"}}
```

### Testing gRPC

```bash
# Notify booking result
grpcurl -plaintext -d '{
  "user_id": "user-uuid",
  "booking_id": "booking-uuid",
  "success": true,
  "message": "Booking confirmed",
  "booking_reference": "BK12345"
}' localhost:50057 realtime.RealtimeService/NotifyBookingResult

# Get connection stats
grpcurl -plaintext localhost:50057 realtime.RealtimeService/GetConnectionStats
```

## Docker

### Build

```bash
docker build -t realtime-service .
```

### Run

```bash
docker run -p 3003:3003 -p 50057:50057 -p 9057:9057 \
  -e REDIS_HOST=host.docker.internal \
  -e JWT_SECRET=your-secret \
  realtime-service
```

## Monitoring

### Prometheus Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `realtime_websocket_connections_total` | Gauge | Total active connections |
| `realtime_websocket_connections_authenticated` | Gauge | Authenticated connections |
| `realtime_websocket_connections_per_room` | GaugeVec | Connections per room |
| `realtime_messages_sent_total` | CounterVec | Messages sent by type |
| `realtime_messages_received_total` | CounterVec | Messages received by type |
| `realtime_notifications_delivered_total` | CounterVec | Notifications delivered |
| `realtime_grpc_requests_total` | CounterVec | gRPC requests by method |
| `realtime_grpc_request_duration_seconds` | HistogramVec | gRPC request duration |
| `realtime_errors_total` | CounterVec | Errors by type |

### Health Checks

- `GET /health` - Basic health check (always returns 200 if server is running)
- `GET /ready` - Readiness check (returns 503 if Redis is unavailable)

## Scaling

### Horizontal Scaling with Redis Pub/Sub

The service uses Redis Pub/Sub for horizontal scaling:

1. All instances subscribe to the same Redis channels
2. When a notification is sent via gRPC, it's published to Redis
3. All instances receive the message and deliver to connected clients
4. Only instances with the target user connected will deliver

### Recommended Configuration for Production

```yaml
replicas: 3-10  # Scale based on connection count
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

## Integration Points

### Inbound (gRPC Callers)

- **Booking Worker**: Calls `NotifyBookingResult` when booking completes
- **Payment Service**: Calls `NotifyPaymentStatus` for payment updates
- **Ticket Service**: Calls `BroadcastEvent` for availability changes

### Outbound

- **Redis Pub/Sub**: Subscribe/publish for horizontal scaling
- **No database**: Stateless service, all state in Redis

## Dependencies

### Go Dependencies

- `github.com/gorilla/websocket` - WebSocket implementation
- `github.com/redis/go-redis/v9` - Redis client
- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/golang-jwt/jwt/v5` - JWT validation
- `go.uber.org/zap` - Structured logging
- `github.com/prometheus/client_golang` - Prometheus metrics

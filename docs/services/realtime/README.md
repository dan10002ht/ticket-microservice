# Realtime Service

**Language:** Go 1.22+
**Ports:** 3003 (HTTP/WebSocket), 50057 (gRPC), 9057 (Metrics)
**Dependencies:** Redis (Pub/Sub)

## Overview

Realtime Service xử lý các notifications real-time và WebSocket broadcasting cho hệ thống ticketing. Service này cho phép cập nhật tức thì về trạng thái booking, payment, ticket availability và queue position.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.22+ |
| WebSocket | Gorilla WebSocket |
| Message Broker | Redis Pub/Sub |
| Internal API | gRPC |
| Metrics | Prometheus |
| Logging | Zap |

**Lý do chọn Go + Gorilla WebSocket:**
- Performance: Handle 100k+ concurrent connections per instance
- Memory efficiency: ~30MB per 10k connections (vs ~100MB Node.js)
- Native concurrency: Goroutines ideal cho WebSocket connections
- Consistent: Phù hợp với existing Go workers

## Features

- **WebSocket Connections**: Client connections với JWT authentication
- **Real-time Notifications**: Push booking, payment, ticket updates
- **Room Management**: Group clients by event for targeted broadcasts
- **Redis Pub/Sub**: Horizontal scaling across multiple instances
- **gRPC API**: Internal service communication

## API Endpoints

### HTTP/WebSocket

| Endpoint | Description |
|----------|-------------|
| `GET /ws` | WebSocket upgrade endpoint |
| `GET /health` | Health check |
| `GET /ready` | Readiness check (includes Redis) |
| `GET /metrics` | Prometheus metrics |

### gRPC Methods (Internal)

| Method | Description | Called By |
|--------|-------------|-----------|
| `NotifyBookingResult` | Push booking result to user | Booking Worker |
| `NotifyQueuePosition` | Update queue position | Booking Worker |
| `NotifyPaymentStatus` | Push payment status | Payment Service |
| `BroadcastEvent` | Broadcast to room | Any Service |
| `SendToUser` | Send to specific user | Any Service |
| `GetConnectionStats` | Get connection statistics | Monitoring |

## WebSocket Protocol

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
| `system:connected` | Server→Client | Connection established |
| `system:ping` | Client→Server | Keep-alive ping |
| `system:pong` | Server→Client | Ping response |
| `room:join` | Client→Server | Join event room |
| `room:leave` | Client→Server | Leave event room |
| `booking:queue_position` | Server→Client | Queue position update |
| `booking:confirmed` | Server→Client | Booking confirmed |
| `booking:failed` | Server→Client | Booking failed |
| `payment:success` | Server→Client | Payment successful |
| `payment:failed` | Server→Client | Payment failed |
| `ticket:availability` | Server→Client | Ticket availability changed |

### Example Messages

**Server → Client (Booking Confirmed)**
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

**Client → Server (Join Room)**
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
# HTTP
HTTP_PORT=3003
HTTP_READ_TIMEOUT=15s
HTTP_WRITE_TIMEOUT=15s

# gRPC
GRPC_PORT=50057

# WebSocket
WS_PING_INTERVAL=30s
WS_PONG_TIMEOUT=60s
WS_ALLOW_ANONYMOUS=true

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=100

# JWT
JWT_SECRET=your-jwt-secret

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9057

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## Development

### Prerequisites

- Go 1.22+
- Redis 7+
- protoc + protoc-gen-go + protoc-gen-go-grpc

### Setup

```bash
# 1. Start Redis
docker run -d -p 6379:6379 redis:7-alpine

# 2. Generate proto files
./scripts/generate-proto.sh

# 3. Run service
go run main.go
```

### Testing

```bash
# Test WebSocket
wscat -c "ws://localhost:3003/ws?token=<JWT>"

# Send ping
> {"type":"system:ping"}

# Join room
> {"type":"room:join","payload":{"room":"event:123"}}

# Test gRPC
grpcurl -plaintext localhost:50057 list

grpcurl -plaintext -d '{
  "user_id": "user-uuid",
  "booking_id": "booking-uuid",
  "success": true,
  "message": "Booking confirmed"
}' localhost:50057 realtime.RealtimeService/NotifyBookingResult

grpcurl -plaintext localhost:50057 realtime.RealtimeService/GetConnectionStats
```

## Project Structure

```
realtime-service/
├── config/config.go              # Configuration
├── internal/
│   ├── grpc/
│   │   ├── server.go             # gRPC server setup
│   │   └── handlers/             # gRPC handlers
│   ├── websocket/
│   │   ├── hub.go                # Connection hub
│   │   ├── client.go             # Client handler
│   │   ├── server.go             # HTTP upgrade
│   │   └── message.go            # Message types
│   ├── pubsub/
│   │   ├── subscriber.go         # Redis subscriber
│   │   └── publisher.go          # Redis publisher
│   ├── service/                  # Business logic
│   ├── middleware/               # Auth middleware
│   └── protos/                   # Generated proto files
├── metrics/metrics.go            # Prometheus metrics
├── pkg/logger/                   # Zap logger
├── scripts/                      # Build scripts
├── main.go                       # Entry point
├── Dockerfile
└── env.example
```

## Integration

### Booking Worker → Realtime Service

```go
// booking-worker/internal/worker/processor.go
p.realtimeClient.NotifyBookingResult(ctx, &pb.NotifyBookingResultRequest{
    UserId:           item.UserID,
    BookingId:        bookingID,
    Success:          true,
    Message:          "Booking confirmed",
    BookingReference: reference,
})
```

### Frontend WebSocket Client

```javascript
const ws = new WebSocket('ws://localhost:3003/ws?token=' + jwtToken);

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  switch (msg.type) {
    case 'booking:confirmed':
      showBookingConfirmation(msg.payload);
      break;
    case 'payment:success':
      showPaymentSuccess(msg.payload);
      break;
  }
};

// Join event room
ws.send(JSON.stringify({
  type: 'room:join',
  payload: { room: 'event:' + eventId }
}));
```

## Scaling

### Horizontal Scaling with Redis Pub/Sub

```
┌─────────────────┐     ┌─────────────────┐
│ Realtime Pod 1  │     │ Realtime Pod 2  │
│   (Go)          │     │   (Go)          │
└────────┬────────┘     └────────┬────────┘
         │                       │
         └───────────┬───────────┘
                     │
              ┌──────▼──────┐
              │    Redis    │
              │  Pub/Sub    │
              └─────────────┘
```

1. All instances subscribe to same Redis channels
2. gRPC notification → publish to Redis
3. All instances receive → deliver to connected clients
4. Only instance with user connected will deliver

### Recommended Production Config

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

## Health Checks

```bash
# HTTP
curl http://localhost:3003/health
curl http://localhost:3003/ready

# gRPC
grpcurl -plaintext localhost:50057 grpc.health.v1.Health/Check
```

## Metrics

Prometheus metrics at `http://localhost:9057/metrics`

| Metric | Type | Description |
|--------|------|-------------|
| `realtime_websocket_connections_total` | Gauge | Total connections |
| `realtime_websocket_connections_authenticated` | Gauge | Authenticated connections |
| `realtime_websocket_connections_per_room` | GaugeVec | Connections per room |
| `realtime_messages_sent_total` | CounterVec | Messages sent by type |
| `realtime_messages_received_total` | CounterVec | Messages received by type |
| `realtime_notifications_delivered_total` | CounterVec | Notifications delivered |
| `realtime_grpc_requests_total` | CounterVec | gRPC requests by method |
| `realtime_errors_total` | CounterVec | Errors by type |

## Proto File

Location: `shared-lib/protos/realtime.proto`

# Service Connections Documentation

## Overview

This document describes the inter-service communication setup in the ticket booking microservices architecture.

## Service Inventory

| Service | Language | gRPC Port | HTTP Port | Description |
|---------|----------|-----------|-----------|-------------|
| Gateway | Node.js | - | 3000 | API Gateway, routing |
| Auth Service | Node.js | 50051 | - | Authentication, JWT |
| Event Service | Node.js | 50054 | - | Event management |
| Ticket Service | Node.js | 50055 | - | Ticket/seat management |
| Booking Service | Java | 50056 | 8080 | Booking orchestration (Saga) |
| Payment Service | Java | 50058 | 8081 | Payment processing (Stripe) |
| User Service | Go | 50052 | - | User profiles, addresses |
| Booking Worker | Go | - | - | Queue processing |
| Email Worker | Go | - | - | Email delivery |
| Realtime Service | Go | 50057 | 3003 | WebSocket notifications |

## Architecture Diagram

```
                                    ┌─────────────────────┐
                                    │      Frontend       │
                                    │   (React/Next.js)   │
                                    └─────────┬───────────┘
                                              │ HTTP/WebSocket
                                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            API Gateway (Node.js)                             │
│                              Port: 3000                                      │
└────┬──────────┬──────────┬──────────┬──────────┬──────────┬────────────────┘
     │          │          │          │          │          │
     │ gRPC     │ gRPC     │ gRPC     │ gRPC     │ gRPC     │ gRPC
     ▼          ▼          ▼          ▼          ▼          ▼
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│  Auth   │ │  Event  │ │ Ticket  │ │ Booking │ │ Payment │ │  User   │
│ Service │ │ Service │ │ Service │ │ Service │ │ Service │ │ Service │
│ :50051  │ │ :50054  │ │ :50055  │ │ :50056  │ │ :50058  │ │ :50052  │
│ Node.js │ │ Node.js │ │ Node.js │ │  Java   │ │  Java   │ │   Go    │
└─────────┘ └─────────┘ └────┬────┘ └────┬────┘ └────┬────┘ └─────────┘
                             │           │           │
                             │    ┌──────┴───────────┘
                             │    │ gRPC calls
                             ▼    ▼
                        ┌─────────────┐
                        │   Booking   │      ┌─────────────┐
                        │   Worker    │─────▶│  Realtime   │
                        │    (Go)     │ gRPC │  Service    │
                        └──────┬──────┘      │ :50057/:3003│
                               │             │    (Go)     │
                               │             └──────┬──────┘
                               │                    │ WebSocket
                               │                    ▼
                               │             ┌─────────────┐
                               │             │  Frontend   │
                               │             │ (Real-time) │
                               │             └─────────────┘
                               │
                               ▼
                        ┌─────────────┐
                        │   Email     │
                        │   Worker    │
                        │    (Go)     │
                        └─────────────┘

Infrastructure:
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ PostgreSQL  │  │    Redis    │  │    Kafka    │
│   (DB)      │  │ (Cache/Queue)│  │  (Events)   │
└─────────────┘  └─────────────┘  └─────────────┘
```

## Connection Matrix

### Gateway → Services (gRPC)

| Target Service | Purpose | Proto File |
|----------------|---------|------------|
| Auth Service | Login, register, token validation | `auth.proto` |
| Event Service | Event CRUD, search | `event.proto` |
| Ticket Service | Seat availability, pricing | `ticket.proto` |
| Booking Service | Create/cancel bookings | `booking.proto` |
| Payment Service | Payment processing | `payment.proto` |
| User Service | Profile, addresses | `user.proto` |

### Booking Service → Services (gRPC)

| Target Service | Purpose | When Called |
|----------------|---------|-------------|
| Ticket Service | Reserve/release seats | During Saga execution |
| Payment Service | Create/cancel payments | During Saga execution |

### Booking Worker → Services

| Target Service | Protocol | Purpose |
|----------------|----------|---------|
| Booking Service | gRPC | Create bookings from queue |
| Realtime Service | gRPC | Notify booking results |

### Realtime Service Connections

| Direction | Protocol | Purpose |
|-----------|----------|---------|
| Frontend → Realtime | WebSocket | Real-time updates |
| Booking Worker → Realtime | gRPC | Push booking notifications |
| Payment Service → Realtime | gRPC | Push payment status |
| Internal (scaling) | Redis Pub/Sub | Multi-instance sync |

## Proto Definitions

### User Service (user.proto)

```protobuf
service UserService {
  // Profile operations
  rpc GetProfile (GetProfileRequest) returns (GetProfileResponse);
  rpc CreateProfile (CreateProfileRequest) returns (CreateProfileResponse);
  rpc UpdateProfile (UpdateProfileRequest) returns (UpdateProfileResponse);

  // Address operations
  rpc GetAddresses (GetAddressesRequest) returns (GetAddressesResponse);
  rpc AddAddress (AddAddressRequest) returns (AddAddressResponse);
  rpc UpdateAddress (UpdateAddressRequest) returns (UpdateAddressResponse);
  rpc DeleteAddress (DeleteAddressRequest) returns (DeleteAddressResponse);
}
```

### Realtime Service (realtime.proto)

```protobuf
service RealtimeService {
  // Called by booking-worker when booking processing completes
  rpc NotifyBookingResult (NotifyBookingResultRequest) returns (NotifyBookingResultResponse);

  // Called to update user's position in booking queue
  rpc NotifyQueuePosition (NotifyQueuePositionRequest) returns (NotifyQueuePositionResponse);

  // Called by payment-service for payment status updates
  rpc NotifyPaymentStatus (NotifyPaymentStatusRequest) returns (NotifyPaymentStatusResponse);

  // Broadcast event to a room (e.g., ticket availability changes)
  rpc BroadcastEvent (BroadcastEventRequest) returns (BroadcastEventResponse);

  // Send notification to a specific user
  rpc SendToUser (SendToUserRequest) returns (SendToUserResponse);

  // Get current connection statistics (for monitoring)
  rpc GetConnectionStats (GetConnectionStatsRequest) returns (GetConnectionStatsResponse);
}
```

### Booking Service (booking.proto)

```protobuf
service BookingService {
  rpc CreateBooking (CreateBookingRequest) returns (CreateBookingResponse);
  rpc GetBooking (GetBookingRequest) returns (GetBookingResponse);
  rpc CancelBooking (CancelBookingRequest) returns (CancelBookingResponse);
  rpc GetUserBookings (GetUserBookingsRequest) returns (GetUserBookingsResponse);
}
```

## Business Flows

### 1. Booking Flow (Saga Pattern)

```
1. User selects seats → Gateway
2. Gateway → Booking Worker (add to Redis queue)
3. Booking Worker processes queue
4. Booking Worker → Booking Service (CreateBooking)
5. Booking Service executes Saga:
   a. Create Booking (PENDING)
   b. Ticket Service → Reserve Seats
   c. Payment Service → Process Payment
   d. Confirm Booking (CONFIRMED)
6. Booking Worker → Realtime Service (NotifyBookingResult)
7. Realtime Service → Frontend (WebSocket push)
```

### 2. Real-time Notification Flow

```
1. Service event occurs (booking confirmed, payment processed)
2. Service → Realtime Service (gRPC: NotifyBookingResult/NotifyPaymentStatus)
3. Realtime Service finds user's WebSocket connection
4. Realtime Service → Frontend (WebSocket message)
5. If multiple instances: Redis Pub/Sub broadcasts to all instances
```

### 3. User Profile Flow

```
1. User requests profile → Gateway
2. Gateway → User Service (gRPC: GetProfile)
3. User Service → PostgreSQL
4. Response flows back through Gateway
```

## Environment Configuration

### Service URLs

```bash
# Gateway
AUTH_SERVICE_URL=localhost:50051
EVENT_SERVICE_URL=localhost:50054
TICKET_SERVICE_URL=localhost:50055
BOOKING_SERVICE_URL=localhost:50056
PAYMENT_SERVICE_URL=localhost:50058
USER_SERVICE_URL=localhost:50052

# Booking Worker
BOOKING_SERVICE_URL=localhost:50056
REALTIME_SERVICE_URL=localhost:50057

# Realtime Service
REDIS_HOST=localhost
REDIS_PORT=6379

# All Services
KAFKA_BROKERS=localhost:9092
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

## Health Checks

All services implement health check endpoints:

```bash
# gRPC Health Check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check
grpcurl -plaintext localhost:50057 grpc.health.v1.Health/Check

# HTTP Health Check (services with HTTP endpoints)
curl http://localhost:3000/health    # Gateway
curl http://localhost:3003/health    # Realtime Service
curl http://localhost:8080/actuator/health  # Booking Service
curl http://localhost:8081/actuator/health  # Payment Service
```

## Monitoring

### Prometheus Metrics Ports

| Service | Metrics Port | Endpoint |
|---------|--------------|----------|
| Gateway | 9090 | /metrics |
| Booking Service | 9091 | /actuator/prometheus |
| Payment Service | 9092 | /actuator/prometheus |
| User Service | 9092 | /metrics |
| Realtime Service | 9057 | /metrics |
| Booking Worker | 9093 | /metrics |
| Email Worker | 9094 | /metrics |

## Troubleshooting

### Common Issues

1. **Service Not Found**
   - Check if service is running on correct port
   - Verify proto file paths
   - Check service URL configuration

2. **Connection Refused**
   - Ensure service is started
   - Check firewall settings
   - Verify port availability

3. **gRPC Deadline Exceeded**
   - Check network latency
   - Increase timeout settings
   - Verify service health

### Debug Commands

```bash
# List available gRPC services
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50052 list
grpcurl -plaintext localhost:50057 list

# Test specific methods
grpcurl -plaintext -d '{"user_id":"test"}' localhost:50052 user.UserService/GetProfile
grpcurl -plaintext localhost:50057 realtime.RealtimeService/GetConnectionStats

# Check Kafka topics
kafka-topics.sh --list --bootstrap-server localhost:9092

# Check Redis
redis-cli ping
redis-cli keys "booking-queue:*"
```

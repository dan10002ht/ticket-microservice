# Service Connections Documentation

## Overview

This document describes the inter-service communication setup in the ticket booking microservices architecture.

## Service Inventory

| Service          | Language | gRPC Port | HTTP Port | Description                         |
| ---------------- | -------- | --------- | --------- | ----------------------------------- |
| Gateway          | Node.js  | -         | 53000     | API Gateway, routing                |
| Auth Service     | Node.js  | 50051     | -         | Authentication, JWT                 |
| User Service     | Go       | 50052     | -         | User profiles, addresses            |
| Ticket Service   | Go       | 50054     | -         | Ticket inventory & lifecycle        |
| Event Service    | Go       | 50053     | -         | Event management, zones, seats      |
| Booking Worker   | Go       | 50056     | -         | Queue processing (gRPC server)      |
| Realtime Service | Go       | 50057     | 3003      | Real-time notifications             |
| Booking Service  | Java     | 50058     | 8084      | Booking orchestration (Saga)        |
| Checkin Service  | Go       | 50059     | -         | Event check-in (QR/barcode)         |
| Invoice Service  | Java     | 50060     | 8083      | Invoice generation & PDF            |
| Email Worker     | Go       | 50061     | -         | Email delivery                      |
| Payment Service  | Java     | 50062     | 8080      | Payment processing (multi-gateway)  |

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
│                         HTTP :53000 | Metrics :53000                         │
└────┬──────────┬──────────┬──────────┬──────────┬──────────┬────────────────┘
     │          │          │          │          │          │
     │ gRPC     │ gRPC     │ gRPC     │ gRPC     │ gRPC     │ gRPC
     ▼          ▼          ▼          ▼          ▼          ▼
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│  Auth   │ │  Event  │ │ Ticket  │ │ Booking │ │ Payment │ │  User   │
│ Service │ │ Service │ │ Service │ │ Service │ │ Service │ │ Service │
│ :50051  │ │ :50053  │ │ :50054  │ │ :50058  │ │ :50062  │ │ :50052  │
│ Node.js │ │   Go    │ │   Go    │ │  Java   │ │  Java   │ │   Go    │
└─────────┘ └─────────┘ └────┬────┘ └────┬────┘ └─────────┘ └─────────┘
                             │           │
                             │    ┌──────┘
                             │    │ gRPC calls (Saga)
                             ▼    ▼
                        ┌─────────────┐
                        │   Booking   │      ┌─────────────┐
                        │   Worker    │─────▶│  Realtime   │
                        │    (Go)     │ gRPC │  Service    │
                        │   :50056   │      │:50057/:3003 │
                        └──────┬──────┘      │    (Go)     │
                               │             └──────┬──────┘
                               │                    │ WebSocket
                               ▼                    ▼
                        ┌─────────────┐      ┌─────────────┐
                        │   Email     │      │  Frontend   │
                        │   Worker    │      │ (Real-time) │
                        │    (Go)     │      └─────────────┘
                        │   :50061   │
                        └─────────────┘

Check-in Flow:
                        ┌─────────────┐
                        │  Checkin    │
                        │  Service    │──gRPC──▶ Ticket Service :50054
                        │    (Go)     │         (GetTicket, UpdateStatus)
                        │   :50059   │
                        └─────────────┘

Invoice Flow:
         Kafka (payment-events)
                        ┌─────────────┐
                        │  Invoice    │
                        │  Service    │──gRPC──▶ clients (GetInvoice, GetPDF)
                        │   (Java)    │
                        │   :50060   │
                        └─────────────┘

Infrastructure:
┌──────────────┐  ┌──────────────────────┐  ┌─────────────┐
│  PostgreSQL  │  │        Redis         │  │    Kafka    │
│ auth  :5432  │  │ cache   :6379        │  │    :9092    │
│ main  :5433  │  │ queue   :6380        │  └─────────────┘
└──────────────┘  │ pubsub  :6381        │
                  └──────────────────────┘
```

## Connection Matrix

### Gateway → Services (gRPC)

| Target Service  | gRPC Port | Purpose                             | Proto File      |
| --------------- | --------- | ----------------------------------- | --------------- |
| Auth Service    | 50051     | Login, register, token validation   | `auth.proto`    |
| User Service    | 50052     | Profile, addresses                  | `user.proto`    |
| Ticket Service  | 50054     | Seat availability, pricing          | `ticket.proto`  |
| Event Service   | 50053     | Event CRUD, zones, seats            | `event.proto`   |
| Booking Service | 50058     | Create/cancel bookings              | `booking.proto` |
| Payment Service | 50062     | Payment processing                  | `payment.proto` |

### Booking Service → Services (gRPC — Saga)

| Target Service  | gRPC Port | Purpose                | When Called           |
| --------------- | --------- | ---------------------- | --------------------- |
| Ticket Service  | 50054     | Reserve/release seats  | During Saga execution |
| Payment Service | 50062     | Create/cancel payments | During Saga execution |

### Booking Worker → Services

| Target Service   | Protocol | Port  | Purpose                    |
| ---------------- | -------- | ----- | -------------------------- |
| Booking Service  | gRPC     | 50058 | Create bookings from queue |
| Realtime Service | gRPC     | 50057 | Notify booking results     |

### Checkin Service → Services

| Target Service  | Protocol | Port  | Purpose                                  |
| --------------- | -------- | ----- | ---------------------------------------- |
| Ticket Service  | gRPC     | 50054 | Validate ticket, update status to "used" |

### Invoice Service (Event-driven)

| Source          | Protocol | Topic            | Trigger Event      |
| --------------- | -------- | ---------------- | ------------------ |
| Payment Service | Kafka    | `payment-events` | `PAYMENT_CAPTURED` |

### Realtime Service Connections

| Direction                  | Protocol      | Purpose                    |
| -------------------------- | ------------- | -------------------------- |
| Frontend → Realtime        | WebSocket     | Real-time updates          |
| Booking Worker → Realtime  | gRPC :50057   | Push booking notifications |
| Payment Service → Realtime | gRPC :50057   | Push payment status        |
| Internal (scaling)         | Redis Pub/Sub | Multi-instance sync :6381  |

## Proto Definitions

### User Service (user.proto)

```protobuf
service UserService {
  rpc GetProfile (GetProfileRequest) returns (GetProfileResponse);
  rpc CreateProfile (CreateProfileRequest) returns (CreateProfileResponse);
  rpc UpdateProfile (UpdateProfileRequest) returns (UpdateProfileResponse);
  rpc GetAddresses (GetAddressesRequest) returns (GetAddressesResponse);
  rpc AddAddress (AddAddressRequest) returns (AddAddressResponse);
  rpc UpdateAddress (UpdateAddressRequest) returns (UpdateAddressResponse);
  rpc DeleteAddress (DeleteAddressRequest) returns (DeleteAddressResponse);
}
```

### Realtime Service (realtime.proto)

```protobuf
service RealtimeService {
  rpc NotifyBookingResult (NotifyBookingResultRequest) returns (NotifyBookingResultResponse);
  rpc NotifyQueuePosition (NotifyQueuePositionRequest) returns (NotifyQueuePositionResponse);
  rpc NotifyPaymentStatus (NotifyPaymentStatusRequest) returns (NotifyPaymentStatusResponse);
  rpc BroadcastEvent (BroadcastEventRequest) returns (BroadcastEventResponse);
  rpc SendToUser (SendToUserRequest) returns (SendToUserResponse);
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

### Checkin Service (checkin.proto)

```protobuf
service CheckinService {
  rpc CheckIn (CheckInRequest) returns (CheckInResponse);
  rpc GetCheckIn (GetCheckInRequest) returns (CheckInRecord);
  rpc ListCheckIns (ListCheckInsRequest) returns (ListCheckInsResponse);
  rpc GetEventStats (GetEventStatsRequest) returns (EventStatsResponse);
  rpc Health (HealthRequest) returns (HealthResponse);
}
```

### Invoice Service (invoice.proto)

```protobuf
service InvoiceService {
  rpc GetInvoice (GetInvoiceRequest) returns (Invoice);
  rpc ListInvoices (ListInvoicesRequest) returns (ListInvoicesResponse);
  rpc GetInvoicePdf (GetInvoicePdfRequest) returns (InvoicePdfResponse);
  rpc Health (HealthRequest) returns (HealthResponse);
}
```

## Business Flows

### 1. Booking Flow (Saga Pattern)

```
1. User selects seats → Gateway :53000
2. Gateway → Booking Worker :50056 (add to Redis queue :6380)
3. Booking Worker processes queue
4. Booking Worker → Booking Service :50058 (CreateBooking)
5. Booking Service executes Saga:
   a. Create Booking (PENDING)
   b. Ticket Service :50054 → Reserve Seats
   c. Payment Service :50062 → Process Payment
   d. Confirm Booking (CONFIRMED)
6. Booking Worker → Realtime Service :50057 (NotifyBookingResult)
7. Realtime Service → Frontend (WebSocket push)
```

### 2. Check-in Flow

```
1. Staff scans QR code → Gateway :53000
2. Gateway → Checkin Service :50059 (CheckIn)
3. Checkin Service → Ticket Service :50054 (GetTicket) — validate
4. Checkin Service → Ticket Service :50054 (UpdateTicketStatus → "used")
5. Checkin Service → DB (INSERT checkin record)
6. Response: success / error (ALREADY_CHECKED_IN, INVALID_TICKET, etc.)
```

### 3. Invoice Generation Flow

```
1. Payment Service publishes PAYMENT_CAPTURED → Kafka :9092 (payment-events)
2. Invoice Service consumes event
3. Invoice Service generates invoice (idempotent — skips if already exists)
4. Invoice Service → DB (INSERT invoice + invoice_items)
5. PDF available via gRPC GetInvoicePdf :50060
```

### 4. Real-time Notification Flow

```
1. Service event occurs (booking confirmed, payment processed)
2. Service → Realtime Service :50057 (gRPC)
3. Realtime Service finds user's WebSocket connection
4. Realtime Service → Frontend (WebSocket message)
5. Multi-instance sync: Redis Pub/Sub :6381
```

## Environment Configuration

### Service URLs (Development)

```bash
# Gateway → downstream gRPC services
GRPC_AUTH_SERVICE_URL=auth-service:50051
GRPC_USER_SERVICE_URL=user-service:50052
GRPC_TICKET_SERVICE_URL=ticket-service:50054
GRPC_EVENT_SERVICE_URL=event-service:50053
GRPC_BOOKING_SERVICE_URL=booking-service:50058
GRPC_PAYMENT_SERVICE_URL=payment-service:50062

# Booking Worker
BOOKING_SERVICE_GRPC_ADDR=booking-service:50058
REALTIME_SERVICE_GRPC_ADDR=realtime-service:50057

# Checkin Service → Ticket Service
TICKET_SERVICE_HOST=ticket-service
TICKET_SERVICE_PORT=50054

# All Services — Infrastructure
KAFKA_BROKERS=kafka:9092
REDIS_CACHE_URL=redis-cache:6379
REDIS_QUEUE_URL=redis-queue:6380
REDIS_PUBSUB_URL=redis-pubsub:6381

# Go services — postgres-main (shared)
DB_HOST=postgres-main
DB_PORT=5432          # container port (mapped from host :5433)

# Auth Service — postgres-auth (dedicated)
DB_HOST=postgres-auth
DB_PORT=5432
```

## Health Checks

```bash
# HTTP Health
curl http://localhost:53000/health         # Gateway
curl http://localhost:3003/health          # Realtime Service
curl http://localhost:8084/actuator/health # Booking Service
curl http://localhost:8080/actuator/health # Payment Service
curl http://localhost:8083/actuator/health # Invoice Service

# gRPC Health
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check  # Auth Service
grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check  # User Service
grpcurl -plaintext localhost:50054 list                          # Ticket Service
grpcurl -plaintext localhost:50053 list                          # Event Service
grpcurl -plaintext localhost:50059 list                          # Checkin Service
grpcurl -plaintext localhost:50057 grpc.health.v1.Health/Check  # Realtime Service
```

## Monitoring

### Prometheus Metrics Ports

| Service          | Metrics Port | Endpoint              |
| ---------------- | ------------ | --------------------- |
| Gateway          | 53000        | /metrics              |
| Auth Service     | 53001        | /metrics              |
| User Service     | 9092         | /metrics              |
| Ticket Service   | 9096         | /metrics              |
| Event Service    | 9095         | /metrics              |
| Booking Worker   | 9091         | /metrics              |
| Email Worker     | 9090         | /metrics              |
| Realtime Service | 9057         | /metrics              |
| Checkin Service  | 2112         | /metrics              |
| Booking Service  | 8084         | /actuator/prometheus  |
| Payment Service  | 8085         | /actuator/prometheus  |
| Invoice Service  | 8083         | /actuator/prometheus  |

## Schema Isolation (postgres-main)

Each service on `postgres-main` owns its own PostgreSQL schema:

| Service         | Schema    | DSN / Config                         |
| --------------- | --------- | ------------------------------------ |
| User Service    | `users`   | `search_path=users` in DSN           |
| Ticket Service  | `tickets` | `search_path=tickets` in DSN         |
| Event Service   | `events`  | `search_path=events` in DSN          |
| Checkin Service | `checkin` | `search_path=checkin` in DSN         |
| Booking Service | `booking` | Flyway `schemas: booking`            |
| Payment Service | `payment` | Flyway `schemas: payment`            |
| Invoice Service | `invoice` | Flyway `schemas: invoice`            |

Auth Service uses a **dedicated** `postgres-auth` instance (`:5432` host, `:5432` container).

## gRPC TLS Configuration

All gRPC connections support optional TLS (mTLS) via the `GRPC_TLS_ENABLED` environment variable.

### Environment Variables

| Variable          | Default              | Description                  |
| ----------------- | -------------------- | ---------------------------- |
| `GRPC_TLS_ENABLED`| `false`              | Set to `true` to enable TLS  |
| `GRPC_TLS_CERT`   | `/certs/server.crt`  | Server/client certificate    |
| `GRPC_TLS_KEY`    | `/certs/server.key`  | Private key                  |
| `GRPC_TLS_CA`     | `/certs/ca.crt`      | CA certificate for trust     |

### Certificate Generation

```bash
# Generate self-signed CA + server + client certs
bash scripts/gen-certs.sh

# Output:
# certs/ca.crt, ca.key       — Root CA
# certs/server.crt, server.key — Server cert (SAN: all service DNS names)
# certs/client.crt, client.key — Client cert (for mTLS)
```

### Enabling TLS

**Development** (default): No changes needed. All services run plaintext.

**Staging/Production**: Use the TLS docker-compose overlay:

```bash
cd deploy/environments/staging
docker compose -f docker-compose.yml -f docker-compose.tls.yml up -d
```

### Verify TLS

```bash
# With TLS enabled
grpcurl -cacert certs/ca.crt -cert certs/client.crt -key certs/client.key \
  localhost:50052 grpc.health.v1.Health/Check

# Without TLS (default dev)
grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check
```

### Implementation Per Language

| Language | Server TLS                                    | Client TLS                            |
| -------- | --------------------------------------------- | ------------------------------------- |
| Go       | `grpctls.ServerOption()` (shared-lib/go/grpctls) | `grpctls.DialOption()`             |
| Node.js  | `getServerCredentials()` (utils/grpcCredentials.js) | `getGrpcCredentials()`          |
| Java     | `grpc.server.security.*` in application.yml   | `GrpcClientConfig.java` conditional  |

## Troubleshooting

### Common Issues

1. **Service Not Found**
   - Check if service is running on correct port
   - Verify proto file paths in `shared-lib/protos/`
   - Check gRPC service URL configuration

2. **Connection Refused**
   - Ensure service is started
   - Verify port availability: `lsof -i :<port>`
   - Check Docker network if running in containers

3. **gRPC Deadline Exceeded**
   - Check network latency
   - Increase timeout settings (gateway default: 60s deadline)
   - Verify service health via health check endpoints

### Debug Commands

```bash
# List available gRPC services
grpcurl -plaintext localhost:50051 list   # Auth
grpcurl -plaintext localhost:50052 list   # User
grpcurl -plaintext localhost:50054 list   # Ticket
grpcurl -plaintext localhost:50053 list   # Event
grpcurl -plaintext localhost:50057 list   # Realtime
grpcurl -plaintext localhost:50059 list   # Checkin

# Test specific methods
grpcurl -plaintext -d '{"user_id":"test"}' localhost:50052 user.UserService/GetProfile
grpcurl -plaintext localhost:50057 realtime.RealtimeService/GetConnectionStats

# Check Kafka topics
kafka-topics.sh --list --bootstrap-server localhost:9092

# Check Redis instances
redis-cli -p 6379 ping   # redis-cache
redis-cli -p 6380 ping   # redis-queue
redis-cli -p 6381 ping   # redis-pubsub
redis-cli -p 6380 keys "booking-queue:*"
```

# User Service

**Language:** Go 1.22+
**Port:** 50052 (gRPC), 9092 (Metrics)
**Database:** PostgreSQL

## Overview

User Service quản lý thông tin profile và địa chỉ của người dùng. Service này được gọi bởi Gateway để phục vụ các API liên quan đến user profile.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.22+ |
| Framework | Native Go + gRPC |
| Database | PostgreSQL (pgx/v5) |
| Logging | Zap |
| Metrics | Prometheus |

**Lý do chọn Go:**
- Simple CRUD service → Go's simplicity phù hợp
- Memory footprint thấp (~20MB vs ~100MB Node.js)
- Fast cold start (~100ms) → tốt cho Kubernetes auto-scaling
- Consistent với existing Go workers (booking-worker, email-worker)

## API Endpoints

### gRPC Methods

| Method | Description |
|--------|-------------|
| `GetProfile` | Lấy thông tin profile theo user_id |
| `CreateProfile` | Tạo profile mới |
| `UpdateProfile` | Cập nhật profile |
| `GetAddresses` | Lấy danh sách địa chỉ của user |
| `AddAddress` | Thêm địa chỉ mới |
| `UpdateAddress` | Cập nhật địa chỉ |
| `DeleteAddress` | Xóa địa chỉ |

### Legacy Methods (Backward Compatibility)

| Method | Description |
|--------|-------------|
| `GetUser` | Get user (maps to profile) |
| `CreateUser` | Create user (creates profile) |
| `ListUsers` | List users (paginated) |

## Database Schema

### user_profiles

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | Reference to auth service |
| first_name | VARCHAR(100) | First name |
| last_name | VARCHAR(100) | Last name |
| email | VARCHAR(255) | Email address |
| phone | VARCHAR(20) | Phone number |
| avatar_url | VARCHAR(500) | Avatar image URL |
| date_of_birth | DATE | Birth date |
| preferences | JSONB | User preferences |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

### user_addresses

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | User reference |
| label | VARCHAR(50) | Address label (home/work/other) |
| street | VARCHAR(255) | Street address |
| city | VARCHAR(100) | City |
| state | VARCHAR(100) | State/Province |
| postal_code | VARCHAR(20) | Postal/ZIP code |
| country | VARCHAR(100) | Country |
| is_default | BOOLEAN | Default address flag |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

## Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_service
DB_SSL_MODE=disable
DB_MAX_CONNS=25
DB_MIN_CONNS=5

# gRPC
GRPC_PORT=50052

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9092

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## Development

### Prerequisites

- Go 1.22+
- PostgreSQL 15+
- protoc + protoc-gen-go + protoc-gen-go-grpc

### Setup

```bash
# 1. Create database
createdb user_service

# 2. Run migrations
psql -d user_service -f migrations/001_init.sql

# 3. Generate proto files
./scripts/generate-proto.sh

# 4. Run service
go run main.go
```

### Testing

```bash
# List services
grpcurl -plaintext localhost:50052 list

# Get profile
grpcurl -plaintext -d '{"user_id": "uuid-here"}' \
  localhost:50052 user.UserService/GetProfile

# Create profile
grpcurl -plaintext -d '{
  "user_id": "uuid",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com"
}' localhost:50052 user.UserService/CreateProfile

# Add address
grpcurl -plaintext -d '{
  "user_id": "uuid",
  "label": "home",
  "street": "123 Main St",
  "city": "NYC",
  "country": "USA"
}' localhost:50052 user.UserService/AddAddress
```

## Project Structure

```
user-service/
├── cmd/server/main.go       # Alternative entry point
├── config/config.go         # Configuration
├── internal/
│   ├── grpc/
│   │   ├── server.go        # gRPC server setup
│   │   └── handlers/        # gRPC handlers
│   ├── service/             # Business logic
│   ├── repository/          # Database operations
│   ├── model/               # Data models
│   └── protos/              # Generated proto files
├── pkg/logger/              # Zap logger
├── migrations/              # SQL migrations
├── scripts/                 # Build scripts
├── main.go                  # Entry point
├── Dockerfile
└── env.example
```

## Integration

### Gateway Integration

Gateway gọi User Service qua gRPC cho các endpoints:

```
GET  /api/users/profile       → GetProfile()
PUT  /api/users/profile       → UpdateProfile()
GET  /api/users/addresses     → GetAddresses()
POST /api/users/addresses     → AddAddress()
PUT  /api/users/addresses/:id → UpdateAddress()
DELETE /api/users/addresses/:id → DeleteAddress()
```

### Proto File

Location: `shared-lib/protos/user.proto`

## Health Checks

```bash
# HTTP health check
curl http://localhost:9092/health
curl http://localhost:9092/ready

# gRPC health check
grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check
```

## Metrics

Prometheus metrics available at `http://localhost:9092/metrics`

| Metric | Type | Description |
|--------|------|-------------|
| `user_service_requests_total` | Counter | Total requests by method |
| `user_service_request_duration_seconds` | Histogram | Request latency |
| `user_service_errors_total` | Counter | Total errors |

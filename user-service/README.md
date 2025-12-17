# User Service (Go)

**Language:** Go

**Why Go?**
- Simple CRUD service → Go's simplicity is ideal
- Low memory footprint (~20MB vs ~100MB Node.js)
- Fast cold start (~100ms) → great for Kubernetes auto-scaling
- Consistent with existing Go workers (booking-worker, email-worker)

## Overview

The User Service manages user profiles and addresses for the booking system. It provides gRPC endpoints for profile CRUD operations and address management.

## Responsibilities

- **Profile Management**: Create, read, update user profiles
- **Address Management**: Manage user shipping/billing addresses
- **Preferences Storage**: Store user preferences as JSON

## Architecture

### Technology Stack

- **Runtime**: Go 1.22+
- **Framework**: Native Go + gRPC
- **Database**: PostgreSQL (pgx driver)
- **Cache**: Redis (future use)
- **Monitoring**: Prometheus + Grafana

### Key Components

```
user-service/
├── cmd/server/          # Entry point
├── config/              # Configuration management
├── internal/
│   ├── grpc/           # gRPC server and handlers
│   ├── service/        # Business logic
│   ├── repository/     # Database operations
│   ├── model/          # Data models
│   └── protos/         # Generated protobuf files
├── pkg/logger/         # Zap logger
├── migrations/         # SQL migrations
└── scripts/            # Build and proto generation scripts
```

## API Endpoints (gRPC)

### Profile Operations

| Method | Description |
|--------|-------------|
| `GetProfile` | Get user profile by user_id |
| `CreateProfile` | Create new user profile |
| `UpdateProfile` | Update existing profile |

### Address Operations

| Method | Description |
|--------|-------------|
| `GetAddresses` | Get all addresses for a user |
| `AddAddress` | Add new address |
| `UpdateAddress` | Update existing address |
| `DeleteAddress` | Delete an address |

### Legacy Operations (Backward Compatibility)

| Method | Description |
|--------|-------------|
| `GetUser` | Get user (maps to profile) |
| `CreateUser` | Create user (creates profile) |
| `ListUsers` | List users (paginated) |

## Configuration

### Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_service
DB_SSL_MODE=disable
DB_MAX_CONNS=25
DB_MIN_CONNS=5

# gRPC Configuration
GRPC_PORT=50052

# Metrics Configuration
METRICS_ENABLED=true
METRICS_PORT=9092

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

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

## Development

### Prerequisites

- Go 1.22+
- PostgreSQL 15+
- protoc (Protocol Buffers compiler)
- protoc-gen-go
- protoc-gen-go-grpc

### Setup

1. Clone the repository
2. Copy environment file:
   ```bash
   cp env.example .env
   ```
3. Create database:
   ```bash
   createdb user_service
   ```
4. Run migrations:
   ```bash
   psql -d user_service -f migrations/001_init.sql
   ```
5. Generate protobuf files:
   ```bash
   ./scripts/generate-proto.sh
   ```
6. Run the service:
   ```bash
   go run main.go
   ```

### Building

```bash
go build -o bin/user-service .
```

### Testing with grpcurl

```bash
# Get profile
grpcurl -plaintext -d '{"user_id": "uuid-here"}' localhost:50052 user.UserService/GetProfile

# Create profile
grpcurl -plaintext -d '{"user_id": "uuid", "first_name": "John", "last_name": "Doe", "email": "john@example.com"}' localhost:50052 user.UserService/CreateProfile

# Add address
grpcurl -plaintext -d '{"user_id": "uuid", "label": "home", "street": "123 Main St", "city": "NYC", "country": "USA"}' localhost:50052 user.UserService/AddAddress
```

## Docker

### Build

```bash
docker build -t user-service .
```

### Run

```bash
docker run -p 50052:50052 -p 9092:9092 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=postgres \
  user-service
```

## Monitoring

### Health Endpoints

- `GET /health` - Basic health check
- `GET /ready` - Readiness check (includes DB connection)
- `GET /metrics` - Prometheus metrics

### gRPC Health Check

The service implements the standard gRPC health protocol.

## Dependencies

### External Services

- **Auth Service**: User authentication (user_id reference)
- **PostgreSQL**: Data storage
- **Redis**: Future caching layer

### Go Dependencies

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `go.uber.org/zap` - Structured logging
- `github.com/prometheus/client_golang` - Prometheus metrics

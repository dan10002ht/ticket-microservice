# 🗄️ Ticket Service Database Setup

## Overview

Ticket Service sử dụng PostgreSQL với PgPool-II để quản lý database connections và load balancing. Database được thiết kế theo event-centric model với các bảng chính:

- **tickets** - Thông tin vé đã được tạo
- **booking_sessions** - Phiên đặt vé với timeout
- **seat_reservations** - Đặt chỗ tạm thời trong quá trình booking

## 🏗️ Database Schema

### Tickets Table

```sql
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    zone_id UUID NOT NULL,
    user_id UUID NOT NULL,
    booking_session_id UUID,
    ticket_number VARCHAR(50) UNIQUE NOT NULL,
    ticket_type VARCHAR(50) DEFAULT 'standard',
    pricing_category VARCHAR(50) NOT NULL,
    base_price DECIMAL(10,2) NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'pending',
    payment_status VARCHAR(20) DEFAULT 'pending',
    -- ... other fields
);
```

### Booking Sessions Table

```sql
CREATE TABLE booking_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    -- ... other fields
);
```

### Seat Reservations Table

```sql
CREATE TABLE seat_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_session_id UUID NOT NULL REFERENCES booking_sessions(id),
    event_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    reservation_token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'reserved',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    -- ... other fields
);
```

## 🚀 Quick Start

### 1. Environment Setup

```bash
# Copy environment template
cp env.example .env

# Edit configuration
vim .env
```

### 2. Database Initialization

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Initialize database
./scripts/init-db.sh
```

### 3. Install Dependencies

```bash
# Install Go dependencies
go mod tidy
```

### 4. Run Service

```bash
# Development mode
./scripts/dev.sh

# Or directly
go run main.go
```

## 🔧 Configuration

### Database Configuration

```bash
# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=booking_system_ticket
DB_USER=postgres
DB_PASSWORD=postgres_password
DB_SSL_MODE=disable

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
DB_CONN_MAX_IDLE_TIME=1m
```

### Service Configuration

```bash
# Server Configuration
PORT=3003
GRPC_PORT=50053

# External Services
EVENT_SERVICE_HOST=localhost
EVENT_SERVICE_PORT=50051
PAYMENT_SERVICE_HOST=localhost
PAYMENT_SERVICE_PORT=50054

# Booking Configuration
BOOKING_SESSION_TIMEOUT=15m
SEAT_RESERVATION_TIMEOUT=15m
MAX_SEATS_PER_BOOKING=10
```

## 📊 Database Features

### 1. Connection Pooling

- **PgPool-II** cho load balancing và failover
- **Connection pooling** với configurable limits
- **Health checks** và automatic reconnection

### 2. Migrations

- **Automatic migrations** on startup
- **Version tracking** trong schema_migrations table
- **Rollback support** (down migrations)

### 3. Performance

- **Optimized indexes** cho common queries
- **Composite indexes** cho complex queries
- **Connection pooling** cho high concurrency

### 4. Data Integrity

- **Foreign key constraints** giữa các bảng
- **Check constraints** cho data validation
- **Unique constraints** cho business rules

## 🔍 Monitoring

### Health Checks

```bash
# Database health check
curl http://localhost:3003/health

# Detailed health info
curl http://localhost:3003/health/detailed
```

### Metrics

```bash
# Prometheus metrics
curl http://localhost:9090/metrics
```

### Logs

```bash
# View service logs
tail -f logs/ticket-service.log

# View database logs
tail -f logs/database.log
```

## 🛠️ Development

### Database Operations

```bash
# Connect to database
psql -h localhost -p 5432 -U postgres -d booking_system_ticket

# Run migrations manually
go run main.go migrate

# Reset database
go run main.go migrate:reset
```

### Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Run with coverage
go test -cover ./...
```

## 🔗 Integration

### Event Service Integration

- **gRPC client** để lấy seat/zone information
- **Real-time availability** updates
- **Seat status** synchronization

### Payment Service Integration

- **gRPC client** cho payment processing
- **Payment status** updates
- **Refund handling**

### Redis Integration

- **Distributed locks** cho seat reservations
- **Session management** với timeout
- **Real-time updates** via pub/sub

## 📝 Notes

- **Event-centric model**: Tất cả tickets đều gắn với event cụ thể
- **Hybrid approach**: Pre-generate seats, generate tickets on checkout
- **Timeout management**: Automatic cleanup của expired sessions
- **Concurrency control**: Redis distributed locks + optimistic locking

## 🚨 Troubleshooting

### Common Issues

1. **Database Connection Failed**

   ```bash
   # Check database status
   pg_isready -h localhost -p 5432

   # Check connection string
   echo $DATABASE_URL
   ```

2. **Migration Failed**

   ```bash
   # Check migration status
   psql -d booking_system_ticket -c "SELECT * FROM schema_migrations;"

   # Manual migration
   psql -d booking_system_ticket -f migrations/001_create_tickets_table.sql
   ```

3. **Service Won't Start**

   ```bash
   # Check logs
   tail -f logs/ticket-service.log

   # Check configuration
   go run main.go config:validate
   ```

### Debug Commands

```bash
# Database connection test
go run main.go db:test

# Configuration validation
go run main.go config:validate

# Service health check
go run main.go health:check
```

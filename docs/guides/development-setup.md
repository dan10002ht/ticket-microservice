# Development Environment Setup

## Prerequisites

| Tool | Version | Required For |
|------|---------|-------------|
| Docker & Docker Compose | Latest | Infrastructure (PostgreSQL, Redis, Kafka) |
| Node.js | 18+ | auth-service, gateway |
| Go | 1.22+ | user, event, ticket, realtime, booking-worker, email-worker |
| Java | 17+ | booking-service, payment-service |
| Maven | 3.8+ | booking-service, payment-service |
| Yarn | 1.x | Node.js services |

---

## Quick Start

### 1. Start Everything

```bash
./scripts/dev-all.sh
```

This will:
1. Start infrastructure (PostgreSQL, Redis, Kafka, Zookeeper) via Docker Compose
2. Wait for infrastructure to be healthy
3. Start all 10 application services with colored log prefixes
4. Verify each service is healthy before reporting status

### 2. Check Health

```bash
# One-time status check
./scripts/health-check.sh

# Live dashboard (refreshes every 5s)
./scripts/health-check.sh --watch

# Custom refresh interval (3s)
./scripts/health-check.sh --watch 3

# Check only infrastructure
./scripts/health-check.sh --infra

# Check only app services
./scripts/health-check.sh --services
```

### 3. View Logs

Each service logs to `logs/<service>.log` with colored prefix in terminal.

```bash
# List available log files
./scripts/dev-logs.sh --list

# Tail a specific service
./scripts/dev-logs.sh --tail auth
./scripts/dev-logs.sh --tail gateway

# Tail all services
./scripts/dev-logs.sh --all

# Filter for errors
./scripts/dev-logs.sh --all --filter ERROR
```

---

## Infrastructure

Started automatically by `dev-all.sh` via Docker Compose:

| Service | Port | Container Name |
|---------|------|---------------|
| PostgreSQL (Auth) | 50432 | dev-postgres-auth |
| PostgreSQL (Main) | 50433 | dev-postgres-main |
| Redis | 50379 | dev-redis |
| Kafka | 50092 | dev-kafka |
| Zookeeper | 50181 | dev-zookeeper |

### Start Infrastructure Only

```bash
./scripts/dev-all.sh --infra-only
```

---

## Application Services

| Service | Type | Port(s) | Start Command |
|---------|------|---------|---------------|
| Auth | gRPC | 50051 | `cd auth-service && yarn dev` |
| User | gRPC | 50052 | `cd user-service && go run main.go` |
| Ticket | gRPC | 50054 | `cd ticket-service && go run main.go` |
| Event | gRPC | 50053 | `cd event-service && go run main.go` |
| Booking Worker | gRPC | 50056 | `cd booking-worker && go run main.go` |
| Realtime | gRPC+HTTP | 50057, 3003 | `cd realtime-service && go run main.go` |
| Booking | gRPC+REST | 50058, 8084 | `cd booking-service && mvn spring-boot:run -Dspring-boot.run.profiles=dev` |
| Checkin | gRPC | 50059 | `cd checkin-service && go run main.go` |
| Invoice | gRPC+REST | 50060, 8083 | `cd invoice-service && mvn spring-boot:run -Dspring-boot.run.profiles=dev` |
| Email Worker | gRPC | 50061 | `cd email-worker && go run main.go` |
| Payment | gRPC+REST | 50062, 8080 | `cd payment-service && mvn spring-boot:run -Dspring-boot.run.profiles=dev` |
| Gateway | HTTP | 53000 | `cd gateway && yarn dev` |

---

## Key URLs

| URL | Description |
|-----|-------------|
| http://localhost:53000/api/docs | Swagger API Documentation |
| http://localhost:53000/health | Gateway Health Check |
| http://localhost:8084/actuator/health | Booking Service Health |
| http://localhost:8080/actuator/health | Payment Service Health |
| ws://localhost:3003 | WebSocket (Realtime) |

---

## Troubleshooting

### Service shows DOWN in health check

1. Check if the port is in use: `lsof -i :<port>`
2. Check the service log: `./scripts/dev-logs.sh --tail <service>`
3. Try restarting just that service manually

### Port conflict

If a service fails to start due to port conflict:

```bash
# Find what's using the port
lsof -i :50054

# Kill the process
kill <PID>
```

### Infrastructure not starting

```bash
# Check Docker containers
docker ps -a | grep dev-

# Restart infrastructure
docker compose -f docker-compose.dev.yml down
./scripts/dev-all.sh --infra-only
```

### Logs are empty

Logs are written to `logs/` directory. If empty, the service may not have started. Check terminal output for startup errors.

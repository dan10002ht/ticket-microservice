# Deployment Guide

Hướng dẫn triển khai Booking System Microservices cho các môi trường khác nhau.

## Cấu trúc thư mục

```
deploy/
├── environments/
│   ├── development/     # Dev: Single PostgreSQL, auto migration
│   ├── staging/         # Staging: HAProxy + PostgreSQL Primary/Replica
│   └── production/      # Prod: PgBouncer + HAProxy + PostgreSQL Cluster
├── scripts/
│   ├── start-dev.sh     # Script khởi động Development
│   ├── start-staging.sh # Script khởi động Staging
│   └── start-prod.sh    # Script khởi động Production
└── shared/
    ├── migrations/      # Migration runners
    └── postgres-init/   # PostgreSQL init scripts
```

## Quick Start

### Development Environment

```bash
# Start với auto migration
./deploy/scripts/start-dev.sh up

# Hoặc start nhanh (không chờ migration)
./deploy/scripts/start-dev.sh up-quick

# Stop
./deploy/scripts/start-dev.sh down
```

### Staging Environment

```bash
# Start
./deploy/scripts/start-staging.sh up

# Check replication status
./deploy/scripts/start-staging.sh status

# Stop
./deploy/scripts/start-staging.sh down
```

### Production Environment

```bash
# 1. Copy và configure .env
cp deploy/environments/production/.env.example deploy/environments/production/.env
# Edit .env với passwords mạnh!

# 2. Start
./deploy/scripts/start-prod.sh up

# 3. Check health
./deploy/scripts/start-prod.sh health

# Stop
./deploy/scripts/start-prod.sh down
```

---

## Chi tiết từng môi trường

### 1. Development

**Đặc điểm:**
- Single PostgreSQL (không replica)
- Auto migration khi start
- Đơn giản, dễ debug

**Services:**
| Service | Port | Description |
|---------|------|-------------|
| postgres-auth | 5432 | Auth database |
| postgres-main | 5433 | Main database (users, events, tickets, bookings) |
| redis | 6379 | Cache & session |
| kafka | 9092 | Message queue |
| zookeeper | 2181 | Kafka coordination |

**Connection:**
```bash
# Auth DB
psql -h localhost -p 5432 -U booking_user -d booking_system_auth

# Main DB
psql -h localhost -p 5433 -U booking_user -d booking_system
```

**Commands:**
```bash
./deploy/scripts/start-dev.sh up           # Start với migration
./deploy/scripts/start-dev.sh up-quick     # Start nhanh
./deploy/scripts/start-dev.sh down         # Stop
./deploy/scripts/start-dev.sh logs         # View logs
./deploy/scripts/start-dev.sh logs kafka   # View specific service logs
./deploy/scripts/start-dev.sh ps           # List services
./deploy/scripts/start-dev.sh migrate      # Run migrations only
./deploy/scripts/start-dev.sh migrate-logs # View migration logs
./deploy/scripts/start-dev.sh clean        # Remove all data (WARNING!)
```

---

### 2. Staging

**Đặc điểm:**
- HAProxy load balancer
- PostgreSQL Primary + 1 Replica
- Read/Write splitting
- Test failover scenarios

**Architecture:**
```
                    ┌─────────────┐
                    │   HAProxy   │
                    │ :5432 :5433 │
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
        ┌──────────┐              ┌──────────┐
        │ Primary  │    ──WAL──▶  │ Replica  │
        │  (R/W)   │              │   (RO)   │
        └──────────┘              └──────────┘
```

**Ports:**
| Service | Port | Description |
|---------|------|-------------|
| HAProxy (Auth Write) | 5432 | → Auth Primary |
| HAProxy (Auth Read) | 5433 | → Auth Replicas |
| HAProxy (Main Write) | 5434 | → Main Primary |
| HAProxy (Main Read) | 5435 | → Main Replicas |
| HAProxy Stats | 8404 | Stats page (admin/admin) |
| Auth Primary | 55432 | Direct access |
| Auth Replica | 55433 | Direct access |
| Main Primary | 55434 | Direct access |
| Main Replica | 55435 | Direct access |

**Connection:**
```bash
# Write (via HAProxy → Primary)
psql -h localhost -p 5432 -U booking_user -d booking_system_auth

# Read (via HAProxy → Replicas)
psql -h localhost -p 5433 -U booking_user -d booking_system_auth

# Direct to Primary (debugging)
psql -h localhost -p 55432 -U booking_user -d booking_system_auth
```

**Commands:**
```bash
./deploy/scripts/start-staging.sh up      # Start
./deploy/scripts/start-staging.sh down    # Stop
./deploy/scripts/start-staging.sh status  # Check replication
./deploy/scripts/start-staging.sh logs    # View logs
./deploy/scripts/start-staging.sh clean   # Remove all data
```

**HAProxy Stats:**
- URL: http://localhost:8404/stats
- User: admin
- Password: admin

---

### 3. Production

**Đặc điểm:**
- PgBouncer connection pooling
- HAProxy load balancer
- PostgreSQL Primary + 2 Replicas
- Synchronous replication
- Prometheus + Grafana monitoring

**Architecture:**
```
     ┌─────────────────────────────────────────────────────────┐
     │                      Application                         │
     └─────────────────────────┬───────────────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │     PgBouncer       │
                    │      :6432          │
                    │  (Connection Pool)  │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │      HAProxy        │
                    │   (Load Balancer)   │
                    └──────────┬──────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Primary  │──WAL──▶ │ Replica1 │         │ Replica2 │
   │  (R/W)   │──WAL──▶ │   (RO)   │         │   (RO)   │
   └──────────┘         └──────────┘         └──────────┘
```

**Ports:**
| Service | Port | Description |
|---------|------|-------------|
| PgBouncer | 6432 | Application connects here |
| HAProxy Stats | 8404 | HAProxy stats |
| Prometheus | 9090 | Metrics |
| Grafana | 3001 | Dashboards |
| Redis | 6379 | Cache |
| Kafka | 9092 | Message queue |

**PgBouncer Databases:**
| Database | Description |
|----------|-------------|
| `booking_system_auth` | Auth DB (Read/Write) |
| `booking_system_auth_ro` | Auth DB (Read Only) |
| `booking_system` | Main DB (Read/Write) |
| `booking_system_ro` | Main DB (Read Only) |

**Connection:**
```bash
# Via PgBouncer (recommended)
psql -h localhost -p 6432 -U booking_user -d booking_system_auth

# Read-only connection
psql -h localhost -p 6432 -U booking_user -d booking_system_auth_ro
```

**Commands:**
```bash
./deploy/scripts/start-prod.sh up      # Start
./deploy/scripts/start-prod.sh down    # Stop
./deploy/scripts/start-prod.sh status  # Check replication + PgBouncer
./deploy/scripts/start-prod.sh health  # Quick health check
./deploy/scripts/start-prod.sh logs    # View logs
./deploy/scripts/start-prod.sh clean   # Remove all data (DANGEROUS!)
```

---

## Migration

### Auto Migration (Development)

Migrations chạy tự động khi start:

```bash
./deploy/scripts/start-dev.sh up
```

Flow:
1. PostgreSQL starts
2. `migrate-auth` runs Knex migrations
3. `migrate-user`, `migrate-ticket` run SQL migrations
4. `migrate-event` runs (depends on user tables)
5. Other services start

### Manual Migration

```bash
# Chỉ chạy migration
./deploy/scripts/start-dev.sh migrate

# Xem logs
./deploy/scripts/start-dev.sh migrate-logs
```

### Migration cho từng service

```bash
# Auth Service (Knex)
cd auth-service
npm run migrate:latest
npm run migrate:status
npm run migrate:rollback

# Other services (SQL files)
psql -h localhost -p 5433 -U booking_user -d booking_system -f migrations/001_init.sql
```

---

## Troubleshooting

### Database không start

```bash
# Check logs
docker logs dev-postgres-auth

# Check health
docker exec dev-postgres-auth pg_isready -U booking_user
```

### Migration failed

```bash
# Xem logs chi tiết
./deploy/scripts/start-dev.sh logs migrate-auth

# Chạy lại migration
./deploy/scripts/start-dev.sh migrate
```

### Replication lag (Staging/Production)

```bash
# Check replication status
./deploy/scripts/start-staging.sh status

# Hoặc manual
docker exec staging-auth-postgres-primary psql -U booking_user -d booking_system_auth \
  -c "SELECT client_addr, state, sent_lsn, replay_lsn FROM pg_stat_replication;"
```

### Reset everything

```bash
# Development
./deploy/scripts/start-dev.sh clean

# Staging
./deploy/scripts/start-staging.sh clean

# Production (requires confirmation)
./deploy/scripts/start-prod.sh clean
```

---

## Environment Variables

### Development (.env)

```env
DB_USER=booking_user
DB_PASSWORD=booking_pass
AUTH_DB_NAME=booking_system_auth
AUTH_DB_PORT=5432
MAIN_DB_NAME=booking_system
MAIN_DB_PORT=5433
REDIS_PORT=6379
KAFKA_PORT=9092
```

### Production (.env)

```env
# IMPORTANT: Change these!
DB_USER=booking_user
DB_PASSWORD=STRONG_PASSWORD_HERE
REPLICATION_PASSWORD=STRONG_REPLICATION_PASSWORD
HAPROXY_STATS_PASSWORD=STRONG_HAPROXY_PASSWORD
GRAFANA_PASSWORD=STRONG_GRAFANA_PASSWORD
```

---

## Service Connections

### Application Configuration

**Development:**
```yaml
# Auth DB
DB_HOST: localhost
DB_PORT: 5432
DB_NAME: booking_system_auth

# Main DB
DB_HOST: localhost
DB_PORT: 5433
DB_NAME: booking_system
```

**Staging/Production (Read/Write Splitting):**
```yaml
# Write operations
DB_WRITE_HOST: localhost
DB_WRITE_PORT: 5432  # HAProxy write port

# Read operations
DB_READ_HOST: localhost
DB_READ_PORT: 5433   # HAProxy read port
```

**Production (via PgBouncer):**
```yaml
# All operations via PgBouncer
DB_HOST: localhost
DB_PORT: 6432
DB_NAME: booking_system_auth     # R/W
DB_NAME_RO: booking_system_auth_ro  # Read-only
```

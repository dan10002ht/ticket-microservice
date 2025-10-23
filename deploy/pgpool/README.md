# PgPool-II Infrastructure

This directory contains the PgPool-II infrastructure setup for the booking system.

## Overview

PgPool-II provides:

- **Automatic load balancing** between master/slave databases
- **Connection pooling** for better performance
- **Automatic failover** when master/slave goes down
- **Read/write routing** (SELECT → slave, INSERT/UPDATE/DELETE → master)

## Architecture

```
[Application Services] → [PgPool-II] → [PostgreSQL Master/Slave]
     ↓                      ↓                    ↓
[Auth Service]         [pgpool-auth]      [auth-master/slave]
[Event Service]        [pgpool-event]     [event-master/slave]
[Ticket Service]       [pgpool-ticket]    [ticket-master/slave]
```

## Quick Start

### 1. Start PgPool-II Infrastructure

```bash
# Start PgPool-II and PostgreSQL clusters
cd deploy/pgpool
chmod +x scripts/*.sh
./scripts/start-pgpool.sh
```

### 2. Start Application Services

```bash
# Start all services with PgPool-II
cd ../..
./scripts/start-dev-with-pgpool.sh
```

### 3. Check Health

```bash
# Check PgPool-II health
cd deploy/pgpool
./scripts/check-health.sh
```

## Services

### PgPool-II Instances

| Service       | Port | PCP Port | Backend             |
| ------------- | ---- | -------- | ------------------- |
| pgpool-auth   | 5432 | 9898     | auth-master/slave   |
| pgpool-event  | 5433 | 9899     | event-master/slave  |
| pgpool-ticket | 5434 | 9900     | ticket-master/slave |

### PostgreSQL Clusters

| Database | Master Port | Slave Port | Database Name         |
| -------- | ----------- | ---------- | --------------------- |
| Auth     | 5435        | 5436       | booking_system_auth   |
| Event    | 5437        | 5438       | booking_system_event  |
| Ticket   | 5439        | 5440       | booking_system_ticket |

## Configuration

### PgPool-II Config Files

- `config/pgpool-auth.conf` - Auth database configuration
- `config/pgpool-event.conf` - Event database configuration
- `config/pgpool-ticket.conf` - Ticket database configuration
- `config/pool_passwd` - User passwords
- `config/pool_hba.conf` - Host-based authentication

### Environment Variables

Copy `env.example` to `.env` and update:

```bash
cp env.example .env
# Edit .env with your values
```

## Usage

### Connect to Databases

```bash
# Auth Database (via PgPool-II)
psql -h localhost -p 5432 -U postgres -d booking_system_auth

# Event Database (via PgPool-II)
psql -h localhost -p 5433 -U postgres -d booking_system_event

# Ticket Database (via PgPool-II)
psql -h localhost -p 5434 -U postgres -d booking_system_ticket
```

### Check PgPool-II Status

```bash
# Check pool status
docker exec pgpool-auth pgpool -c "show pool_status"

# Check backend nodes
docker exec pgpool-auth pgpool -c "show pool_nodes"

# Check health
docker exec pgpool-auth pgpool -c "show pool_health"
```

## Service Integration

### Update Service Configs

Update your service database configs to use PgPool-II endpoints:

```javascript
// Before (Direct DB)
const config = {
  connection: {
    host: "postgres-master:5432",
    // ...
  },
};

// After (PgPool-II)
const config = {
  connection: {
    host: "pgpool-auth:5432", // PgPool-II endpoint
    // ...
  },
};
```

### Environment Variables for Services

```bash
# Auth Service
PGPOOL_AUTH_HOST=pgpool-auth
PGPOOL_AUTH_PORT=5432

# Event Service
PGPOOL_EVENT_HOST=pgpool-event
PGPOOL_EVENT_PORT=5432

# Ticket Service
PGPOOL_TICKET_HOST=pgpool-ticket
PGPOOL_TICKET_PORT=5432
```

## Monitoring

### Logs

```bash
# View PgPool-II logs
docker logs pgpool-auth
docker logs pgpool-event
docker logs pgpool-ticket

# View PostgreSQL logs
docker logs auth-master
docker logs auth-slave
```

### Metrics

PgPool-II provides metrics via PCP (PgPool Control Protocol):

```bash
# Check pool status
docker exec pgpool-auth pgpool -c "show pool_status"

# Check backend status
docker exec pgpool-auth pgpool -c "show pool_nodes"

# Check health
docker exec pgpool-auth pgpool -c "show pool_health"
```

## Troubleshooting

### Common Issues

1. **Connection refused**: Check if PgPool-II is running
2. **Authentication failed**: Check pool_passwd file
3. **Backend down**: Check PostgreSQL containers
4. **Load balancing not working**: Check PgPool-II configuration

### Debug Commands

```bash
# Check container status
docker ps | grep pgpool

# Check network connectivity
docker exec pgpool-auth ping auth-master

# Check database connectivity
docker exec pgpool-auth psql -h auth-master -U postgres -d booking_system_auth -c "SELECT 1"
```

## Scripts

- `scripts/start-pgpool.sh` - Start PgPool-II infrastructure
- `scripts/stop-pgpool.sh` - Stop PgPool-II infrastructure
- `scripts/check-health.sh` - Check PgPool-II health

## Benefits

- **Automatic load balancing**: Read queries go to slaves, write queries go to master
- **Connection pooling**: Better performance and resource utilization
- **Automatic failover**: When master/slave goes down, PgPool-II handles it
- **Easy scaling**: Add more slaves without changing application code
- **Monitoring**: Built-in health checks and metrics

## Production Considerations

- Use dedicated servers for PgPool-II
- Setup monitoring and alerting
- Configure SSL/TLS for secure connections
- Setup backup and recovery procedures
- Use connection limits and timeouts
- Monitor performance metrics


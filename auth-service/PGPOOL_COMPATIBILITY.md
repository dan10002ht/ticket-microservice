# PgPool-II Compatibility Guide

## Overview

This document explains how the Auth Service has been updated to work with PgPool-II instead of manual master/slave database connections.

## Changes Made

### 1. Database Configuration (`src/config/databaseConfig.js`)

- **Before**: Manual master/slave connection management
- **After**: Single PgPool-II endpoint that handles routing automatically

```javascript
// Before (Manual Master/Slave)
const masterDb = knex(masterConfig);
const slaveDbs = slaveConfigs.map((config) => knex(config));
const getSlaveDb = () => slaveDbs[currentIndex];

// After (PgPool-II)
const masterDb = knex(pgpoolConfig); // PgPool-II endpoint
const getSlaveDb = () => masterDb; // Same endpoint, PgPool-II routes automatically
```

### 2. BaseRepository (`src/repositories/baseRepository.js`)

- **Before**: Different connections for master/slave operations
- **After**: Single connection, PgPool-II handles routing

```javascript
// Before
getMasterDb() {
  return masterDb(this.tableName);
}
getSlaveDb() {
  return getSlaveDb()(this.tableName);
}

// After
getMasterDb() {
  return masterDb(this.tableName); // PgPool-II routes to master
}
getSlaveDb() {
  return masterDb(this.tableName); // PgPool-II routes to slave
}
```

## How PgPool-II Works

### Automatic Query Routing

- **SELECT queries** → Automatically routed to slave databases
- **INSERT/UPDATE/DELETE queries** → Automatically routed to master database
- **Transactions** → Always use master database

### Benefits

1. **Transparent**: No code changes needed in repositories
2. **Automatic**: PgPool-II handles routing based on query type
3. **Load Balancing**: Multiple slaves for read operations
4. **Failover**: Automatic failover when master/slave goes down

## Repository Usage

### Read Operations (Automatically routed to slaves)

```javascript
// These will be routed to slave databases by PgPool-II
const users = await userRepository.findAll();
const user = await userRepository.findByEmail('user@example.com');
const count = await userRepository.count({ is_active: true });
```

### Write Operations (Automatically routed to master)

```javascript
// These will be routed to master database by PgPool-II
const newUser = await userRepository.create(userData);
const updatedUser = await userRepository.updateById(id, updateData);
const deletedUser = await userRepository.deleteById(id);
```

### Transactions (Always use master)

```javascript
// Transactions always use master database
await userRepository.transaction(async (trx) => {
  await trx('users').insert(userData);
  await trx('user_profiles').insert(profileData);
});
```

## Environment Variables

### PgPool-II Configuration

```bash
# PgPool-II endpoints
PGPOOL_AUTH_HOST=pgpool-auth
PGPOOL_AUTH_PORT=5432

# Database credentials
DB_NAME=booking_system_auth
DB_USER=postgres
DB_PASSWORD=postgres_password
```

### Connection Pool Settings

```bash
# Connection pool configuration
DB_POOL_MIN=2
DB_POOL_MAX=10
```

## Migration Guide

### For Developers

1. **No code changes needed** in repository methods
2. **Same API** - all existing code works unchanged
3. **Better performance** - automatic load balancing
4. **Automatic failover** - no manual handling needed

### For Operations

1. **Start PgPool-II infrastructure** first
2. **Services connect to PgPool-II** instead of direct PostgreSQL
3. **Monitor PgPool-II health** instead of individual databases

## Health Checks

### PgPool-II Health

```bash
# Check PgPool-II status
./deploy/pgpool/scripts/check-health.sh

# Check specific service
curl http://localhost:9898/pcp_node_info
```

### Service Health

```javascript
// Health check in service
const health = await checkDatabaseHealth();
console.log(health); // { pgpool: true, timestamp: '...' }
```

## Troubleshooting

### Common Issues

1. **Connection refused**: Check if PgPool-II is running
2. **Slow queries**: Check PgPool-II load balancing configuration
3. **Transaction errors**: Ensure transactions use master database

### Debug Commands

```bash
# Check PgPool-II logs
docker logs pgpool-auth

# Check database connections
docker exec pgpool-auth psql -h localhost -p 5432 -U postgres -d booking_system_auth -c "SELECT 1"
```

## Performance Benefits

### Before (Manual Master/Slave)

- Manual connection management
- No automatic load balancing
- Manual failover handling
- Complex configuration

### After (PgPool-II)

- Automatic connection pooling
- Automatic load balancing
- Automatic failover
- Simple configuration
- Better performance
- Easier maintenance

## Conclusion

The migration to PgPool-II provides:

- **Better performance** through automatic load balancing
- **Simplified code** with no manual master/slave handling
- **Automatic failover** for high availability
- **Easier maintenance** with centralized connection management

All existing repository code continues to work without changes, but now benefits from PgPool-II's automatic routing and load balancing capabilities.

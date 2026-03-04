import dotenv from 'dotenv';
import knex from 'knex';

dotenv.config();

// PgPool-II configuration - handles master/slave routing automatically
const dbConfig = {
  client: 'postgresql',
  connection: {
    host: process.env.PGPOOL_AUTH_HOST || process.env.DB_MASTER_HOST || 'localhost',
    port: process.env.PGPOOL_AUTH_PORT || process.env.DB_MASTER_PORT || 50432,
    database: process.env.DB_NAME || process.env.DB_MASTER_NAME || 'booking_system_auth',
    user: process.env.DB_USER || process.env.DB_MASTER_USER || 'booking_user',
    password: process.env.DB_PASSWORD || process.env.DB_MASTER_PASSWORD || 'booking_pass',
    ssl: process.env.NODE_ENV === 'production' ? { rejectUnauthorized: false } : false,
  },
  pool: {
    min: parseInt(process.env.DB_POOL_MIN) || 2,
    max: parseInt(process.env.DB_POOL_MAX) || 10,
    acquireTimeoutMillis: 30000,
    createTimeoutMillis: 30000,
    destroyTimeoutMillis: 5000,
    idleTimeoutMillis: 30000,
    reapIntervalMillis: 1000,
    createRetryIntervalMillis: 100,
  },
  migrations: {
    directory: './migrations',
  },
  seeds: {
    directory: './seeds',
  },
};

// Single database connection - PgPool-II handles master/slave routing
const db = knex(dbConfig);

// Health check function
const checkDatabaseHealth = async () => {
  try {
    await db.raw('SELECT 1');
    return { status: 'healthy', timestamp: new Date().toISOString() };
  } catch (error) {
    console.error('Database health check failed:', error.message);
    return { status: 'unhealthy', error: error.message, timestamp: new Date().toISOString() };
  }
};

// Graceful shutdown
const closeConnections = async () => {
  try {
    await db.destroy();
    console.log('Database connection closed');
  } catch (error) {
    console.error('Error closing database connection:', error);
  }
};

export { db, checkDatabaseHealth, closeConnections, dbConfig };

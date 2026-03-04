import { createClient } from 'redis';
import config from '../config/index.js';
import logger from './logger.js';

/**
 * Shared Redis client for the gateway.
 * Used by rate-limit-redis and any other middleware that needs Redis.
 */
const redisClient = createClient({
  socket: {
    host: config.redis.host,
    port: config.redis.port,
  },
  password: config.redis.password || undefined,
  database: config.redis.database,
});

redisClient.on('error', (err) => {
  logger.error('Redis client error', { error: err.message });
});

redisClient.on('connect', () => {
  logger.info('Redis client connected', {
    host: config.redis.host,
    port: config.redis.port,
  });
});

// Connect lazily — errors are non-fatal; rate limiter will fall back to memory store.
redisClient.connect().catch((err) => {
  logger.warn('Redis connection failed — rate limiter will use in-memory store', {
    error: err.message,
  });
});

export default redisClient;

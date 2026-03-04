import rateLimit from 'express-rate-limit';
import slowDown from 'express-slow-down';
import { RedisStore } from 'rate-limit-redis';
import redisClient from '../utils/redisClient.js';
import config from '../config/index.js';

/**
 * Build a Redis-backed store for express-rate-limit.
 * Falls back gracefully to in-memory if Redis is unavailable.
 */
const buildRedisStore = (prefix) => {
  try {
    return new RedisStore({
      sendCommand: (...args) => redisClient.sendCommand(args),
      prefix: `rl:gateway:${prefix}:`,
    });
  } catch {
    // Redis not ready — in-memory fallback (single-instance only)
    return undefined;
  }
};

/**
 * Rate limiting middleware configuration.
 * Uses Redis so limits are shared across multiple gateway instances.
 */
export const rateLimitMiddleware = () => {
  const limiter = rateLimit({
    windowMs: config.rateLimit.windowMs,
    max: config.rateLimit.maxRequests,
    message: 'Too many requests from this IP, please try again later.',
    standardHeaders: true,
    legacyHeaders: false,
    store: buildRedisStore('limit'),
  });

  const speedLimiter = slowDown({
    windowMs: config.rateLimit.windowMs,
    delayAfter: config.rateLimit.delayAfter,
    delayMs: (used) => {
      const delayAfter = config.rateLimit.delayAfter;
      const delayMs = config.rateLimit.delayMs;
      return (used - delayAfter) * delayMs;
    },
    store: buildRedisStore('slow'),
  });

  return { limiter, speedLimiter };
};

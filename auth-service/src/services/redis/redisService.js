import Redis from 'ioredis';
import logger from '../../utils/logger.js';

// Redis client configuration
const redisClient = new Redis({
  host: process.env.REDIS_HOST || 'localhost',
  port: process.env.REDIS_PORT || 6379,
  password: process.env.REDIS_PASSWORD,
  db: process.env.REDIS_DB || 0,
  retryDelayOnFailover: 100,
  maxRetriesPerRequest: 3,
  lazyConnect: true,
});

// Redis configuration for different databases
const REDIS_DBS = {
  CACHE: 0,
  SESSIONS: 1,
  EMAIL_VERIFICATION: 2,
  PASSWORD_RESET: 4,
  BACKGROUND_JOBS: 6,
};

// Cache TTL configuration
const CACHE_TTL = {
  USER_PROFILE: 300, // 5 minutes
  USER_ROLES: 600, // 10 minutes
  TOKEN_VALIDATION: 60, // 1 minute
  PERMISSIONS: 1800, // 30 minutes
  ORGANIZATION: 900, // 15 minutes
  SESSION: 300, // 5 minutes
  EMAIL_VERIFICATION: 15 * 60, // 15 minutes
  PASSWORD_RESET: 60 * 60, // 1 hour
  PASSWORD_RESET_ATTEMPTS: 15 * 60, // 15 minutes
};

// Cache key patterns
const CACHE_KEYS = {
  USER_PROFILE: (userId) => `user:profile:${userId}`,
  USER_ROLES: (userId) => `user:roles:${userId}`,
  TOKEN_VALIDATION: (token) => `token:validation:${token}`,
  PERMISSIONS: (roleId) => `role:permissions:${roleId}`,
  ORGANIZATION: (orgId) => `org:${orgId}`,
  SESSION: (sessionId) => `session:${sessionId}`,
  USER_SESSIONS: (userId) => `user:sessions:${userId}`,
  EMAIL_VERIFICATION: (userId) => `email_verification:${userId}`,
  PASSWORD_RESET: (tokenHash) => `pwd_reset:${tokenHash}`,
  PASSWORD_RESET_ATTEMPTS: (userId) => `pwd_reset:attempts:${userId}`,
};

// ========== CORE REDIS OPERATIONS ==========

/**
 * Set key with TTL
 */
export async function setWithTTL(key, value, ttlSeconds) {
  try {
    await redisClient.setex(
      key,
      ttlSeconds,
      typeof value === 'string' ? value : JSON.stringify(value)
    );
    logger.debug(`Redis SETEX: ${key}, TTL: ${ttlSeconds}s`);
    return true;
  } catch (error) {
    logger.error(`Redis SETEX failed for key ${key}:`, error);
    throw error;
  }
}

/**
 * Get value by key
 */
export async function get(key) {
  try {
    const value = await redisClient.get(key);
    logger.debug(`Redis GET: ${key}, found: ${!!value}`);
    return value;
  } catch (error) {
    logger.error(`Redis GET failed for key ${key}:`, error);
    throw error;
  }
}

/**
 * Get and parse JSON value
 */
export async function getJSON(key) {
  try {
    const value = await redisClient.get(key);
    if (!value) return null;

    const parsed = JSON.parse(value);
    logger.debug(`Redis GET JSON: ${key}, found: ${!!parsed}`);
    return parsed;
  } catch (error) {
    logger.error(`Redis GET JSON failed for key ${key}:`, error);
    throw error;
  }
}

/**
 * Delete key
 */
export async function del(key) {
  try {
    const result = await redisClient.del(key);
    logger.debug(`Redis DEL: ${key}, deleted: ${result > 0}`);
    return result > 0;
  } catch (error) {
    logger.error(`Redis DEL failed for key ${key}:`, error);
    throw error;
  }
}

/**
 * Increment counter
 */
export async function incr(key) {
  try {
    const result = await redisClient.incr(key);
    logger.debug(`Redis INCR: ${key}, new value: ${result}`);
    return result;
  } catch (error) {
    logger.error(`Redis INCR failed for key ${key}:`, error);
    throw error;
  }
}

/**
 * Set TTL for existing key
 */
export async function expire(key, ttlSeconds) {
  try {
    const result = await redisClient.expire(key, ttlSeconds);
    logger.debug(`Redis EXPIRE: ${key}, TTL: ${ttlSeconds}s, success: ${result}`);
    return result;
  } catch (error) {
    logger.error(`Redis EXPIRE failed for key ${key}:`, error);
    throw error;
  }
}

/**
 * Check if key exists
 */
export async function exists(key) {
  try {
    const result = await redisClient.exists(key);
    logger.debug(`Redis EXISTS: ${key}, exists: ${result > 0}`);
    return result > 0;
  } catch (error) {
    logger.error(`Redis EXISTS failed for key ${key}:`, error);
    throw error;
  }
}

// ========== CACHE OPERATIONS ==========

/**
 * Set user profile cache
 */
export async function setUserProfile(userId, profileData) {
  const key = CACHE_KEYS.USER_PROFILE(userId);
  return await setWithTTL(key, profileData, CACHE_TTL.USER_PROFILE);
}

/**
 * Get user profile cache
 */
export async function getUserProfile(userId) {
  const key = CACHE_KEYS.USER_PROFILE(userId);
  return await getJSON(key);
}

/**
 * Invalidate user profile cache
 */
export async function invalidateUserProfile(userId) {
  const key = CACHE_KEYS.USER_PROFILE(userId);
  return await del(key);
}

/**
 * Set user roles cache
 */
export async function setUserRoles(userId, rolesData) {
  const key = CACHE_KEYS.USER_ROLES(userId);
  return await setWithTTL(key, rolesData, CACHE_TTL.USER_ROLES);
}

/**
 * Get user roles cache
 */
export async function getUserRoles(userId) {
  const key = CACHE_KEYS.USER_ROLES(userId);
  return await getJSON(key);
}

/**
 * Invalidate user roles cache
 */
export async function invalidateUserRoles(userId) {
  const key = CACHE_KEYS.USER_ROLES(userId);
  return await del(key);
}

/**
 * Set token validation cache
 */
export async function setTokenValidation(token, validationData) {
  const key = CACHE_KEYS.TOKEN_VALIDATION(token);
  return await setWithTTL(key, validationData, CACHE_TTL.TOKEN_VALIDATION);
}

/**
 * Get token validation cache
 */
export async function getTokenValidation(token) {
  const key = CACHE_KEYS.TOKEN_VALIDATION(token);
  return await getJSON(key);
}

/**
 * Invalidate token validation cache
 */
export async function invalidateTokenValidation(token) {
  const key = CACHE_KEYS.TOKEN_VALIDATION(token);
  return await del(key);
}

// ========== PASSWORD RESET OPERATIONS ==========

/**
 * Set password reset token
 */
export async function setPasswordResetToken(tokenHash, tokenData) {
  const key = CACHE_KEYS.PASSWORD_RESET(tokenHash);
  return await setWithTTL(key, tokenData, CACHE_TTL.PASSWORD_RESET);
}

/**
 * Get password reset token
 */
export async function getPasswordResetToken(tokenHash) {
  const key = CACHE_KEYS.PASSWORD_RESET(tokenHash);
  return await getJSON(key);
}

/**
 * Delete password reset token
 */
export async function deletePasswordResetToken(tokenHash) {
  const key = CACHE_KEYS.PASSWORD_RESET(tokenHash);
  return await del(key);
}

/**
 * Increment password reset attempts
 */
export async function incrementPasswordResetAttempts(userId) {
  const key = CACHE_KEYS.PASSWORD_RESET_ATTEMPTS(userId);
  const attempts = await incr(key);
  await expire(key, CACHE_TTL.PASSWORD_RESET_ATTEMPTS);
  return attempts;
}

/**
 * Get password reset attempts
 */
export async function getPasswordResetAttempts(userId) {
  const key = CACHE_KEYS.PASSWORD_RESET_ATTEMPTS(userId);
  const attempts = await get(key);
  return attempts ? parseInt(attempts) : 0;
}

/**
 * Clear password reset attempts
 */
export async function clearPasswordResetAttempts(userId) {
  const key = CACHE_KEYS.PASSWORD_RESET_ATTEMPTS(userId);
  return await del(key);
}

// ========== EMAIL VERIFICATION OPERATIONS ==========

/**
 * Set email verification PIN
 */
export async function setEmailVerificationPin(userId, pinCode) {
  const key = CACHE_KEYS.EMAIL_VERIFICATION(userId);
  return await setWithTTL(key, pinCode, CACHE_TTL.EMAIL_VERIFICATION);
}

/**
 * Get email verification PIN
 */
export async function getEmailVerificationPin(userId) {
  const key = CACHE_KEYS.EMAIL_VERIFICATION(userId);
  return await get(key);
}

/**
 * Delete email verification PIN
 */
export async function deleteEmailVerificationPin(userId) {
  const key = CACHE_KEYS.EMAIL_VERIFICATION(userId);
  return await del(key);
}

// ========== SESSION OPERATIONS ==========

/**
 * Set session data
 */
export async function setSession(sessionId, sessionData) {
  const key = CACHE_KEYS.SESSION(sessionId);
  return await setWithTTL(key, sessionData, CACHE_TTL.SESSION);
}

/**
 * Get session data
 */
export async function getSession(sessionId) {
  const key = CACHE_KEYS.SESSION(sessionId);
  return await getJSON(key);
}

/**
 * Delete session
 */
export async function deleteSession(sessionId) {
  const key = CACHE_KEYS.SESSION(sessionId);
  return await del(key);
}

/**
 * Set user sessions list
 */
export async function setUserSessions(userId, sessionsData) {
  const key = CACHE_KEYS.USER_SESSIONS(userId);
  return await setWithTTL(key, sessionsData, CACHE_TTL.SESSION);
}

/**
 * Get user sessions list
 */
export async function getUserSessions(userId) {
  const key = CACHE_KEYS.USER_SESSIONS(userId);
  return await getJSON(key);
}

/**
 * Invalidate user sessions cache
 */
export async function invalidateUserSessions(userId) {
  const key = CACHE_KEYS.USER_SESSIONS(userId);
  return await del(key);
}

// ========== ORGANIZATION OPERATIONS ==========

/**
 * Set organization data
 */
export async function setOrganization(orgId, orgData) {
  const key = CACHE_KEYS.ORGANIZATION(orgId);
  return await setWithTTL(key, orgData, CACHE_TTL.ORGANIZATION);
}

/**
 * Get organization data
 */
export async function getOrganization(orgId) {
  const key = CACHE_KEYS.ORGANIZATION(orgId);
  return await getJSON(key);
}

/**
 * Invalidate organization cache
 */
export async function invalidateOrganization(orgId) {
  const key = CACHE_KEYS.ORGANIZATION(orgId);
  return await del(key);
}

// ========== PERMISSIONS OPERATIONS ==========

/**
 * Set role permissions
 */
export async function setRolePermissions(roleId, permissionsData) {
  const key = CACHE_KEYS.PERMISSIONS(roleId);
  return await setWithTTL(key, permissionsData, CACHE_TTL.PERMISSIONS);
}

/**
 * Get role permissions
 */
export async function getRolePermissions(roleId) {
  const key = CACHE_KEYS.PERMISSIONS(roleId);
  return await getJSON(key);
}

/**
 * Invalidate role permissions cache
 */
export async function invalidateRolePermissions(roleId) {
  const key = CACHE_KEYS.PERMISSIONS(roleId);
  return await del(key);
}

// ========== UTILITY OPERATIONS ==========

/**
 * Get all keys matching pattern
 */
export async function getKeys(pattern) {
  try {
    const keys = await redisClient.keys(pattern);
    logger.debug(`Redis KEYS: ${pattern}, found: ${keys.length} keys`);
    return keys;
  } catch (error) {
    logger.error(`Redis KEYS failed for pattern ${pattern}:`, error);
    throw error;
  }
}

/**
 * Delete keys by pattern
 */
export async function deleteByPattern(pattern) {
  try {
    const keys = await redisClient.keys(pattern);
    if (keys.length > 0) {
      const deleted = await redisClient.del(...keys);
      logger.debug(`Redis DEL pattern: ${pattern}, deleted: ${deleted} keys`);
      return deleted;
    }
    return 0;
  } catch (error) {
    logger.error(`Redis DEL pattern failed for ${pattern}:`, error);
    throw error;
  }
}

/**
 * Get Redis info
 */
export async function getRedisInfo() {
  try {
    const info = await redisClient.info();
    return info;
  } catch (error) {
    logger.error('Redis INFO failed:', error);
    throw error;
  }
}

/**
 * Get Redis memory usage
 */
export async function getMemoryUsage() {
  try {
    const memory = await redisClient.info('memory');
    return memory;
  } catch (error) {
    logger.error('Redis memory info failed:', error);
    throw error;
  }
}

/**
 * Flush current database
 */
export async function flushDB() {
  try {
    await redisClient.flushdb();
    logger.warn('Redis database flushed');
    return true;
  } catch (error) {
    logger.error('Redis FLUSHDB failed:', error);
    throw error;
  }
}

// ========== HEALTH CHECK ==========

/**
 * Check Redis health
 */
export async function healthCheck() {
  try {
    await redisClient.ping();
    return { status: 'healthy', service: 'redis' };
  } catch (error) {
    logger.error('Redis health check failed:', error);
    return { status: 'unhealthy', service: 'redis', error: error.message };
  }
}

/**
 * Close Redis connection
 */
export async function closeConnection() {
  try {
    await redisClient.quit();
    logger.info('Redis connection closed');
  } catch (error) {
    logger.error('Error closing Redis connection:', error);
  }
}

// ========== EVENT HANDLERS ==========

redisClient.on('connect', () => {
  logger.info('✅ Redis connection established');
});

redisClient.on('error', (error) => {
  logger.error('❌ Redis connection error:', error);
});

redisClient.on('close', () => {
  logger.warn('⚠️ Redis connection closed');
});

redisClient.on('reconnecting', () => {
  logger.info('🔄 Redis reconnecting...');
});

// Export Redis client for advanced operations
export { redisClient };

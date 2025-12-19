import { grpcClients } from '../grpc/clients.js';
import logger from '../utils/logger.js';

// Cache for user permissions to reduce gRPC calls
const permissionCache = new Map();
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

/**
 * Check if user has a specific permission
 * @param {string} userId - User ID
 * @param {string} permission - Permission to check
 * @param {string} correlationId - Request correlation ID
 * @returns {Promise<boolean>} - Whether user has permission
 */
const checkPermission = async (userId, permission, correlationId) => {
  try {
    const cacheKey = `permission:${userId}:${permission}`;
    const cached = permissionCache.get(cacheKey);

    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
      logger.debug('Permission check cache hit', {
        userId,
        permission,
        correlationId,
        cached: true,
      });
      return cached.hasPermission;
    }

    const response = await grpcClients.authService.CheckPermission({
      user_id: userId,
      permission_name: permission,
    });

    const hasPermission = response.hasPermission;

    // Cache the result
    permissionCache.set(cacheKey, {
      hasPermission,
      timestamp: Date.now(),
    });

    logger.debug('Permission check completed', {
      userId,
      permission,
      hasPermission,
      correlationId,
      cached: false,
    });

    return hasPermission;
  } catch (error) {
    logger.error('Permission check failed', {
      userId,
      permission,
      error: error.message,
      correlationId,
    });

    // On error, deny access for security
    return false;
  }
};

/**
 * Check if user has permission for a specific resource
 * @param {string} userId - User ID
 * @param {string} permission - Permission to check
 * @param {string} resourceType - Type of resource (e.g., 'event', 'booking')
 * @param {string} resourceId - Resource ID
 * @param {string} correlationId - Request correlation ID
 * @returns {Promise<boolean>} - Whether user has permission
 */
const checkResourcePermission = async (
  userId,
  permission,
  resourceType,
  resourceId,
  correlationId
) => {
  try {
    const cacheKey = `resource:${userId}:${permission}:${resourceType}:${resourceId}`;
    const cached = permissionCache.get(cacheKey);

    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
      logger.debug('Resource permission check cache hit', {
        userId,
        permission,
        resourceType,
        resourceId,
        correlationId,
        cached: true,
      });
      return cached.hasPermission;
    }

    const response = await grpcClients.authService.CheckResourcePermission({
      user_id: userId,
      resource: resourceType,
      action: permission,
      context: { resource_id: resourceId },
    });

    const hasPermission = response.hasPermission;

    // Cache the result
    permissionCache.set(cacheKey, {
      hasPermission,
      timestamp: Date.now(),
    });

    logger.debug('Resource permission check completed', {
      userId,
      permission,
      resourceType,
      resourceId,
      hasPermission,
      correlationId,
      cached: false,
    });

    return hasPermission;
  } catch (error) {
    logger.error('Resource permission check failed', {
      userId,
      permission,
      resourceType,
      resourceId,
      error: error.message,
      correlationId,
    });

    // On error, deny access for security
    return false;
  }
};

/**
 * Get all roles for a user
 * @param {string} userId - User ID
 * @param {string} correlationId - Request correlation ID
 * @returns {Promise<Array>} - Array of user roles
 */
const getUserRoles = async (userId, correlationId) => {
  try {
    const cacheKey = `roles:${userId}`;
    const cached = permissionCache.get(cacheKey);

    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
      logger.debug('User roles cache hit', {
        userId,
        correlationId,
        cached: true,
      });
      return cached.roles;
    }

    const response = await grpcClients.authService.GetUserRoles({
      user_id: userId,
    });

    const roles = response.roles || [];

    // Cache the result
    permissionCache.set(cacheKey, {
      roles,
      timestamp: Date.now(),
    });

    logger.debug('User roles retrieved', {
      userId,
      rolesCount: roles.length,
      correlationId,
      cached: false,
    });

    return roles;
  } catch (error) {
    logger.error('Get user roles failed', {
      userId,
      error: error.message,
      correlationId,
    });

    return [];
  }
};

/**
 * Check multiple permissions at once
 * @param {string} userId - User ID
 * @param {Array<string>} permissions - Array of permissions to check
 * @param {string} correlationId - Request correlation ID
 * @returns {Promise<Object>} - Object with permission results
 */
const batchCheckPermissions = async (userId, permissions, correlationId) => {
  try {
    const response = await grpcClients.authService.BatchCheckPermissions({
      user_id: userId,
      permission_names: permissions,
    });

    const results = response.results || {};

    // Cache individual permission results
    Object.entries(results).forEach(([permission, hasPermission]) => {
      const cacheKey = `permission:${userId}:${permission}`;
      permissionCache.set(cacheKey, {
        hasPermission,
        timestamp: Date.now(),
      });
    });

    logger.debug('Batch permission check completed', {
      userId,
      permissionsCount: permissions.length,
      correlationId,
    });

    return results;
  } catch (error) {
    logger.error('Batch permission check failed', {
      userId,
      permissions,
      error: error.message,
      correlationId,
    });

    // Return all false for security
    const results = {};
    permissions.forEach((permission) => {
      results[permission] = false;
    });
    return results;
  }
};

/**
 * Clear cache for a specific user
 * @param {string} userId - User ID to clear cache for
 */
const clearUserCache = (userId) => {
  const keysToDelete = [];

  for (const [key] of permissionCache) {
    if (key.includes(`:${userId}:`)) {
      keysToDelete.push(key);
    }
  }

  keysToDelete.forEach((key) => permissionCache.delete(key));

  logger.debug('Cleared user permission cache', {
    userId,
    clearedKeys: keysToDelete.length,
  });
};

/**
 * Clear all permission cache
 */
const clearAllCache = () => {
  const cacheSize = permissionCache.size;
  permissionCache.clear();

  logger.debug('Cleared all permission cache', {
    clearedEntries: cacheSize,
  });
};

/**
 * Get cache statistics
 * @returns {Object} - Cache statistics
 */
const getCacheStats = () => {
  return {
    size: permissionCache.size,
    entries: Array.from(permissionCache.entries()).map(([key, value]) => ({
      key,
      timestamp: value.timestamp,
      age: Date.now() - value.timestamp,
    })),
  };
};

export {
  checkPermission,
  checkResourcePermission,
  getUserRoles,
  batchCheckPermissions,
  clearUserCache,
  clearAllCache,
  getCacheStats,
};

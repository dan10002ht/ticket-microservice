import { cacheUserProfile, cacheUserRoles } from '../../services/internal/cacheService.js';
import logger from '../logger.js';

/**
 * Cache User Profile Job Handler
 *
 * This job handles:
 * 1. Caching user profile data in Redis
 * 2. Setting appropriate TTL for profile cache
 * 3. Error handling and logging
 */
export async function handleCacheUserProfileJob(data) {
  const { userId, userProfile } = data;

  try {
    logger.info(`Processing cache user profile job for user: ${userId}`);

    await cacheUserProfile(userId, userProfile);

    logger.info(`User profile cached successfully for userId: ${userId}`);

    return {
      success: true,
      message: 'User profile cached successfully',
      userId,
    };
  } catch (error) {
    logger.error(`Cache user profile job failed for user ${userId}:`, error);
    throw error;
  }
}

/**
 * Cache User Roles Job Handler
 *
 * This job handles:
 * 1. Caching user roles data in Redis
 * 2. Setting appropriate TTL for roles cache
 * 3. Error handling and logging
 */
export async function handleCacheUserRolesJob(data) {
  const { userId, roles } = data;

  try {
    logger.info(`Processing cache user roles job for user: ${userId}`);

    await cacheUserRoles(userId, roles || []);

    logger.info(`User roles cached successfully for userId: ${userId}`);

    return {
      success: true,
      message: 'User roles cached successfully',
      userId,
      rolesCount: roles ? roles.length : 0,
    };
  } catch (error) {
    logger.error(`Cache user roles job failed for user ${userId}:`, error);
    throw error;
  }
}

/**
 * Cache User Data Job Handler (Combined)
 *
 * This job handles:
 * 1. Caching both user profile and roles in parallel
 * 2. Setting appropriate TTL for both caches
 * 3. Error handling and logging
 */
export async function handleCacheUserDataJob(data) {
  const { userId, userProfile, userRoles } = data;

  try {
    logger.info(`Processing cache user data job for user: ${userId}`);

    // Cache both profile and roles in parallel for better performance
    await Promise.all([
      cacheUserProfile(userId, userProfile),
      cacheUserRoles(userId, userRoles || []),
    ]);

    logger.info(`User data (profile + roles) cached successfully for userId: ${userId}`);

    return {
      success: true,
      message: 'User data cached successfully',
      userId,
      rolesCount: userRoles ? userRoles.length : 0,
    };
  } catch (error) {
    logger.error(`Cache user data job failed for user ${userId}:`, error);
    throw error;
  }
}

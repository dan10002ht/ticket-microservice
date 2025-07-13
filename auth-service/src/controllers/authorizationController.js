import * as permissionService from '../services/internal/permissionService.js';
import * as roleService from '../services/internal/roleService.js';
import logger from '../utils/logger.js';

/**
 * Check if user has specific permission
 */
export async function checkPermission(call, callback) {
  try {
    const { user_id, permission_name } = call.request;

    // Input validation
    if (!user_id || !permission_name) {
      return callback({
        code: 3, // INVALID_ARGUMENT
        message: 'user_id and permission_name are required',
      });
    }

    logger.info('Checking permission', {
      userId: user_id,
      permissionName: permission_name,
      correlationId: call.metadata.get('correlation-id')[0] || 'unknown',
    });

    // Check permission
    const hasPermission = await permissionService.hasPermission(user_id, permission_name);

    // Get user roles for context
    const userRoles = await roleService.getUserRoles(user_id);
    const roleNames = userRoles.map((role) => role.name);

    const response = {
      allowed: hasPermission,
      reason: hasPermission ? 'Permission granted' : 'Permission denied',
      roles: roleNames,
    };

    logger.info('Permission check completed', {
      userId: user_id,
      permissionName: permission_name,
      allowed: hasPermission,
      roles: roleNames,
    });

    callback(null, response);
  } catch (error) {
    logger.error('Permission check failed', {
      userId: call.request.user_id,
      permissionName: call.request.permission_name,
      error: error.message,
      stack: error.stack,
    });

    callback({
      code: 13, // INTERNAL
      message: `Permission check failed: ${error.message}`,
    });
  }
}

/**
 * Check if user has permission for specific resource and action
 */
export async function checkResourcePermission(call, callback) {
  try {
    const { user_id, resource, action, context = {} } = call.request;

    // Input validation
    if (!user_id || !resource || !action) {
      return callback({
        code: 3, // INVALID_ARGUMENT
        message: 'user_id, resource, and action are required',
      });
    }

    logger.info('Checking resource permission', {
      userId: user_id,
      resource: resource,
      action: action,
      context: context,
      correlationId: call.metadata.get('correlation-id')[0] || 'unknown',
    });

    // Check resource permission
    const hasPermission = await permissionService.hasResourcePermission(user_id, resource, action);

    // Get user permissions for context
    const userPermissions = await permissionService.getUserPermissions(user_id);
    const permissionNames = userPermissions.permissions.map((p) => p.name);

    const response = {
      allowed: hasPermission,
      reason: hasPermission ? 'Resource permission granted' : 'Resource permission denied',
      permissions: permissionNames,
    };

    logger.info('Resource permission check completed', {
      userId: user_id,
      resource: resource,
      action: action,
      allowed: hasPermission,
      permissions: permissionNames,
    });

    callback(null, response);
  } catch (error) {
    logger.error('Resource permission check failed', {
      userId: call.request.user_id,
      resource: call.request.resource,
      action: call.request.action,
      error: error.message,
      stack: error.stack,
    });

    callback({
      code: 13, // INTERNAL
      message: `Resource permission check failed: ${error.message}`,
    });
  }
}

/**
 * Get user permissions
 */
export async function getUserPermissions(call, callback) {
  try {
    const { user_id } = call.request;

    // Input validation
    if (!user_id) {
      return callback({
        code: 3, // INVALID_ARGUMENT
        message: 'user_id is required',
      });
    }

    logger.info('Getting user permissions', {
      userId: user_id,
      correlationId: call.metadata.get('correlation-id')[0] || 'unknown',
    });

    // Get user permissions
    const result = await permissionService.getUserPermissions(user_id);

    const response = {
      success: true,
      message: 'User permissions retrieved successfully',
      permissions: result.permissions,
    };

    logger.info('User permissions retrieved', {
      userId: user_id,
      permissionsCount: result.permissions.length,
    });

    callback(null, response);
  } catch (error) {
    logger.error('Get user permissions failed', {
      userId: call.request.user_id,
      error: error.message,
      stack: error.stack,
    });

    callback({
      code: 13, // INTERNAL
      message: `Failed to get user permissions: ${error.message}`,
    });
  }
}

/**
 * Get user roles
 */
export async function getUserRoles(call, callback) {
  try {
    const { user_id } = call.request;

    // Input validation
    if (!user_id) {
      return callback({
        code: 3, // INVALID_ARGUMENT
        message: 'user_id is required',
      });
    }

    logger.info('Getting user roles', {
      userId: user_id,
      correlationId: call.metadata.get('correlation-id')[0] || 'unknown',
    });

    // Get user roles
    const roles = await roleService.getUserRoles(user_id);

    const response = {
      success: true,
      message: 'User roles retrieved successfully',
      roles: roles,
    };

    logger.info('User roles retrieved', {
      userId: user_id,
      rolesCount: roles.length,
    });

    callback(null, response);
  } catch (error) {
    logger.error('Get user roles failed', {
      userId: call.request.user_id,
      error: error.message,
      stack: error.stack,
    });

    callback({
      code: 13, // INTERNAL
      message: `Failed to get user roles: ${error.message}`,
    });
  }
}

/**
 * Batch check multiple permissions
 */
export async function batchCheckPermissions(call, callback) {
  try {
    const { user_id, permission_names } = call.request;

    // Input validation
    if (!user_id || !permission_names || permission_names.length === 0) {
      return callback({
        code: 3, // INVALID_ARGUMENT
        message: 'user_id and permission_names array are required',
      });
    }

    logger.info('Batch checking permissions', {
      userId: user_id,
      permissionCount: permission_names.length,
      permissions: permission_names,
      correlationId: call.metadata.get('correlation-id')[0] || 'unknown',
    });

    // Get user roles for context
    const userRoles = await roleService.getUserRoles(user_id);
    const roleNames = userRoles.map((role) => role.name);

    // Check each permission
    const results = {};
    const checkPromises = permission_names.map(async (permissionName) => {
      try {
        const hasPermission = await permissionService.hasPermission(user_id, permissionName);
        results[permissionName] = hasPermission;
      } catch (error) {
        logger.error('Individual permission check failed in batch', {
          userId: user_id,
          permissionName: permissionName,
          error: error.message,
        });
        results[permissionName] = false;
      }
    });

    await Promise.all(checkPromises);

    const response = {
      success: true,
      message: 'Batch permission check completed',
      results: results,
      roles: roleNames,
    };

    logger.info('Batch permission check completed', {
      userId: user_id,
      results: results,
      roles: roleNames,
    });

    callback(null, response);
  } catch (error) {
    logger.error('Batch permission check failed', {
      userId: call.request.user_id,
      permissionNames: call.request.permission_names,
      error: error.message,
      stack: error.stack,
    });

    callback({
      code: 13, // INTERNAL
      message: `Batch permission check failed: ${error.message}`,
    });
  }
}

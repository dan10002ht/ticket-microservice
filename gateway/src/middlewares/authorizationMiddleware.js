import {
  checkPermission,
  checkResourcePermission,
  getUserRoles,
  batchCheckPermissions,
} from '../clients/authorizationClient.js';
import logger from '../utils/logger.js';

/**
 * Middleware to check if user has specific permission
 * @param {string|Array<string>} permissions - Permission(s) to check
 * @returns {Function} - Express middleware function
 */
const requirePermission = (permissions) => {
  return async (req, res, next) => {
    try {
      if (!req.user || !req.user.id) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Authentication required',
          correlationId: req.correlationId,
        });
      }

      const permissionList = Array.isArray(permissions) ? permissions : [permissions];

      // Check if user has any of the required permissions
      let hasPermission = false;

      for (const permission of permissionList) {
        const result = await checkPermission(req.user.id, permission, req.correlationId);
        if (result) {
          hasPermission = true;
          break;
        }
      }

      if (!hasPermission) {
        logger.warn('Permission denied', {
          userId: req.user.id,
          requiredPermissions: permissionList,
          correlationId: req.correlationId,
        });

        return res.status(403).json({
          error: 'Forbidden',
          message: 'Insufficient permissions',
          requiredPermissions: permissionList,
          correlationId: req.correlationId,
        });
      }

      logger.debug('Permission check passed', {
        userId: req.user.id,
        permissions: permissionList,
        correlationId: req.correlationId,
      });

      next();
    } catch (error) {
      logger.error('Permission check error', {
        userId: req.user?.id,
        permissions,
        error: error.message,
        correlationId: req.correlationId,
      });

      return res.status(500).json({
        error: 'Internal Server Error',
        message: 'Authorization check failed',
        correlationId: req.correlationId,
      });
    }
  };
};

/**
 * Middleware to check if user has permission for a specific resource
 * @param {string} permission - Permission to check
 * @param {Function} resourceExtractor - Function to extract resource info from request
 * @returns {Function} - Express middleware function
 */
const requireResourcePermission = (permission, resourceExtractor) => {
  return async (req, res, next) => {
    try {
      if (!req.user || !req.user.id) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Authentication required',
          correlationId: req.correlationId,
        });
      }

      // Extract resource information from request
      const resourceInfo = resourceExtractor(req);
      if (!resourceInfo || !resourceInfo.type || !resourceInfo.id) {
        return res.status(400).json({
          error: 'Bad Request',
          message: 'Invalid resource information',
          correlationId: req.correlationId,
        });
      }

      const hasPermission = await checkResourcePermission(
        req.user.id,
        permission,
        resourceInfo.type,
        resourceInfo.id,
        req.correlationId
      );

      if (!hasPermission) {
        logger.warn('Resource permission denied', {
          userId: req.user.id,
          permission,
          resourceType: resourceInfo.type,
          resourceId: resourceInfo.id,
          correlationId: req.correlationId,
        });

        return res.status(403).json({
          error: 'Forbidden',
          message: 'Insufficient permissions for this resource',
          resourceType: resourceInfo.type,
          resourceId: resourceInfo.id,
          correlationId: req.correlationId,
        });
      }

      logger.debug('Resource permission check passed', {
        userId: req.user.id,
        permission,
        resourceType: resourceInfo.type,
        resourceId: resourceInfo.id,
        correlationId: req.correlationId,
      });

      next();
    } catch (error) {
      logger.error('Resource permission check error', {
        userId: req.user?.id,
        permission,
        error: error.message,
        correlationId: req.correlationId,
      });

      return res.status(500).json({
        error: 'Internal Server Error',
        message: 'Resource authorization check failed',
        correlationId: req.correlationId,
      });
    }
  };
};

/**
 * Middleware to check if user has any of the specified roles
 * @param {Array<string>} roles - Roles to check
 * @returns {Function} - Express middleware function
 */
const requireRole = (roles) => {
  return async (req, res, next) => {
    try {
      if (!req.user || !req.user.id) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Authentication required',
          correlationId: req.correlationId,
        });
      }

      const userRoles = await getUserRoles(req.user.id, req.correlationId);
      const hasRole = roles.some((role) => userRoles.includes(role));

      if (!hasRole) {
        logger.warn('Role check failed', {
          userId: req.user.id,
          userRoles,
          requiredRoles: roles,
          correlationId: req.correlationId,
        });

        return res.status(403).json({
          error: 'Forbidden',
          message: 'Insufficient role permissions',
          requiredRoles: roles,
          userRoles,
          correlationId: req.correlationId,
        });
      }

      logger.debug('Role check passed', {
        userId: req.user.id,
        userRoles,
        requiredRoles: roles,
        correlationId: req.correlationId,
      });

      next();
    } catch (error) {
      logger.error('Role check error', {
        userId: req.user?.id,
        roles,
        error: error.message,
        correlationId: req.correlationId,
      });

      return res.status(500).json({
        error: 'Internal Server Error',
        message: 'Role authorization check failed',
        correlationId: req.correlationId,
      });
    }
  };
};

/**
 * Middleware to check multiple permissions and require all of them
 * @param {Array<string>} permissions - Permissions to check
 * @returns {Function} - Express middleware function
 */
const requireAllPermissions = (permissions) => {
  return async (req, res, next) => {
    try {
      if (!req.user || !req.user.id) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Authentication required',
          correlationId: req.correlationId,
        });
      }

      const results = await batchCheckPermissions(req.user.id, permissions, req.correlationId);
      const missingPermissions = permissions.filter((permission) => !results[permission]);

      if (missingPermissions.length > 0) {
        logger.warn('All permissions check failed', {
          userId: req.user.id,
          requiredPermissions: permissions,
          missingPermissions,
          correlationId: req.correlationId,
        });

        return res.status(403).json({
          error: 'Forbidden',
          message: 'Missing required permissions',
          missingPermissions,
          correlationId: req.correlationId,
        });
      }

      logger.debug('All permissions check passed', {
        userId: req.user.id,
        permissions,
        correlationId: req.correlationId,
      });

      next();
    } catch (error) {
      logger.error('All permissions check error', {
        userId: req.user?.id,
        permissions,
        error: error.message,
        correlationId: req.correlationId,
      });

      return res.status(500).json({
        error: 'Internal Server Error',
        message: 'Authorization check failed',
        correlationId: req.correlationId,
      });
    }
  };
};

/**
 * Middleware to check if user is the owner of a resource
 * @param {Function} ownerExtractor - Function to extract owner ID from request
 * @returns {Function} - Express middleware function
 */
const requireOwnership = (ownerExtractor) => {
  return async (req, res, next) => {
    try {
      if (!req.user || !req.user.id) {
        return res.status(401).json({
          error: 'Unauthorized',
          message: 'Authentication required',
          correlationId: req.correlationId,
        });
      }

      const ownerId = ownerExtractor(req);
      if (!ownerId) {
        return res.status(400).json({
          error: 'Bad Request',
          message: 'Could not determine resource owner',
          correlationId: req.correlationId,
        });
      }

      if (req.user.id !== ownerId) {
        logger.warn('Ownership check failed', {
          userId: req.user.id,
          ownerId,
          correlationId: req.correlationId,
        });

        return res.status(403).json({
          error: 'Forbidden',
          message: 'You can only access your own resources',
          correlationId: req.correlationId,
        });
      }

      logger.debug('Ownership check passed', {
        userId: req.user.id,
        ownerId,
        correlationId: req.correlationId,
      });

      next();
    } catch (error) {
      logger.error('Ownership check error', {
        userId: req.user?.id,
        error: error.message,
        correlationId: req.correlationId,
      });

      return res.status(500).json({
        error: 'Internal Server Error',
        message: 'Ownership check failed',
        correlationId: req.correlationId,
      });
    }
  };
};

/**
 * Helper function to extract resource info from URL parameters
 * @param {string} resourceType - Type of resource
 * @param {string} idParam - Parameter name for resource ID
 * @returns {Function} - Resource extractor function
 */
const extractResourceFromParams = (resourceType, idParam = 'id') => {
  return (req) => ({
    type: resourceType,
    id: req.params[idParam],
  });
};

/**
 * Helper function to extract resource info from request body
 * @param {string} resourceType - Type of resource
 * @param {string} idField - Field name for resource ID in body
 * @returns {Function} - Resource extractor function
 */
const extractResourceFromBody = (resourceType, idField = 'id') => {
  return (req) => ({
    type: resourceType,
    id: req.body[idField],
  });
};

/**
 * Helper function to extract owner ID from request body
 * @param {string} ownerField - Field name for owner ID in body
 * @returns {Function} - Owner extractor function
 */
const extractOwnerFromBody = (ownerField = 'userId') => {
  return (req) => req.body[ownerField];
};

/**
 * Helper function to extract owner ID from URL parameters
 * @param {string} ownerParam - Parameter name for owner ID
 * @returns {Function} - Owner extractor function
 */
const extractOwnerFromParams = (ownerParam = 'userId') => {
  return (req) => req.params[ownerParam];
};

export {
  requirePermission,
  requireResourcePermission,
  requireRole,
  requireAllPermissions,
  requireOwnership,
  extractResourceFromParams,
  extractResourceFromBody,
  extractOwnerFromBody,
  extractOwnerFromParams,
};

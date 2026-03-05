import jwt from 'jsonwebtoken';
import config from '../config/index.js';
import logger from '../utils/logger.js';
import redisClient from '../utils/redisClient.js';

const buildMeta = (correlationId) => ({
  correlationId,
  timestamp: new Date().toISOString(),
});

const sendAuthError = (res, statusCode, code, message, correlationId) => {
  return res.status(statusCode).json({
    error: { code, message },
    meta: buildMeta(correlationId),
  });
};

const authMiddleware = async (req, res, next) => {
  try {
    // Get token from header
    const authHeader = req.headers.authorization;
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return sendAuthError(res, 401, 'UNAUTHENTICATED', 'No token provided', req.correlationId);
    }

    const token = authHeader.substring(7); // Remove 'Bearer ' prefix

    // Verify token signature and expiry
    const decoded = jwt.verify(token, config.jwt.secret);

    // Check revocation boundary — auth-service writes token:nbf:{userId} on logout.
    // Any token with iat <= that timestamp is considered revoked.
    const userId = decoded.sub || decoded.userId;
    if (userId) {
      try {
        const nbf = await redisClient.get(`token:nbf:${userId}`);
        if (nbf !== null && decoded.iat <= Number(nbf)) {
          return sendAuthError(res, 401, 'TOKEN_REVOKED', 'Token has been revoked', req.correlationId);
        }
      } catch (redisErr) {
        // Redis unavailable — log and continue (fail-open to preserve availability)
        logger.warn('Token revocation check skipped — Redis unavailable', {
          userId,
          error: redisErr.message,
          correlationId: req.correlationId,
        });
      }
    }

    // Add user info to request
    req.user = {
      id: decoded.sub || decoded.userId,
      email: decoded.email,
      role: decoded.role,
      permissions: decoded.permissions || [],
    };

    // Add correlation ID to response headers
    if (req.correlationId) {
      res.setHeader('X-Correlation-ID', req.correlationId);
    }

    logger.info('Authentication successful', {
      userId: req.user.id,
      email: req.user.email,
      correlationId: req.correlationId
    });

    next();
  } catch (error) {
    logger.error('Authentication failed', {
      error: error.message,
      correlationId: req.correlationId
    });

    if (error.name === 'TokenExpiredError') {
      return sendAuthError(res, 401, 'TOKEN_EXPIRED', 'Token expired', req.correlationId);
    }

    if (error.name === 'JsonWebTokenError') {
      return sendAuthError(res, 401, 'INVALID_TOKEN', 'Invalid token', req.correlationId);
    }

    return sendAuthError(res, 401, 'UNAUTHENTICATED', 'Authentication failed', req.correlationId);
  }
};

// Optional authentication middleware (doesn't fail if no token)
const optionalAuthMiddleware = async (req, res, next) => {
  try {
    const authHeader = req.headers.authorization;
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const token = authHeader.substring(7);
      const decoded = jwt.verify(token, config.jwt.secret);

      req.user = {
        id: decoded.sub,
        email: decoded.email,
        role: decoded.role,
        permissions: decoded.permissions || []
      };
    }
    next();
  } catch (error) {
    // Continue without authentication
    next();
  }
};

// Role-based authorization middleware
const requireRole = (roles) => {
  return (req, res, next) => {
    if (!req.user) {
      return sendAuthError(res, 401, 'UNAUTHENTICATED', 'Authentication required', req.correlationId);
    }

    if (!roles.includes(req.user.role)) {
      return sendAuthError(res, 403, 'FORBIDDEN', 'Insufficient permissions', req.correlationId);
    }

    next();
  };
};

// Permission-based authorization middleware
const requirePermission = (permissions) => {
  return (req, res, next) => {
    if (!req.user) {
      return sendAuthError(res, 401, 'UNAUTHENTICATED', 'Authentication required', req.correlationId);
    }

    const hasPermission = permissions.some(permission =>
      req.user.permissions.includes(permission)
    );

    if (!hasPermission) {
      return sendAuthError(res, 403, 'FORBIDDEN', 'Insufficient permissions', req.correlationId);
    }

    next();
  };
};

// Alias for authMiddleware - used in routes that need explicit auth check
const requireAuth = authMiddleware;

export { authMiddleware, optionalAuthMiddleware, requireRole, requirePermission, requireAuth };
export default authMiddleware;

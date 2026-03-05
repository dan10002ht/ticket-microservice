import { validationResult } from 'express-validator';
import logger from './logger.js';
import { getErrorMapping } from './errorMapping.js';
import { toCamelCase } from './caseTransform.js';

/**
 * Build meta object for responses
 */
const buildMeta = (correlationId) => ({
  correlationId,
  timestamp: new Date().toISOString(),
});

/**
 * Handle validation errors
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @returns {boolean} True if validation passed, false if failed
 */
export const handleValidation = (req, res) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    res.status(400).json({
      error: {
        code: 'VALIDATION_ERROR',
        message: 'Validation failed',
        details: errors.array(),
      },
      meta: buildMeta(req.correlationId),
    });
    return false;
  }
  return true;
};

/**
 * Send success response (HTTP 2xx)
 *
 * Format:
 * {
 *   data: { ... },
 *   meta: { correlationId, timestamp }
 * }
 *
 * @param {Object} res - Express response object
 * @param {number} statusCode - HTTP status code
 * @param {Object} data - Response data
 * @param {string} correlationId - Correlation ID
 */
export const sendSuccessResponse = (res, statusCode, data, correlationId) => {
  res.status(statusCode).json({
    data: toCamelCase(data),
    meta: buildMeta(correlationId),
  });
};

/**
 * Send error response (HTTP 4xx/5xx)
 *
 * Format:
 * {
 *   error: { code, message, details? },
 *   meta: { correlationId, timestamp }
 * }
 *
 * @param {Object} res - Express response object
 * @param {number} statusCode - HTTP status code
 * @param {string} message - Error message
 * @param {string} correlationId - Correlation ID
 * @param {string} code - Error code
 * @param {*} details - Additional error details
 */
export const sendErrorResponse = (
  res,
  statusCode,
  message,
  correlationId,
  code = 'UNKNOWN_ERROR',
  details = null
) => {
  const errorObj = {
    code,
    message,
  };

  if (details) {
    errorObj.details = details;
  }

  res.status(statusCode).json({
    error: errorObj,
    meta: buildMeta(correlationId),
  });
};

/**
 * Handle gRPC errors and map to HTTP status codes
 * @param {Object} res - Express response object
 * @param {Error} error - gRPC error
 * @param {string} correlationId - Correlation ID
 * @param {string} serviceName - Service name for logging
 * @param {string} methodName - Method name for logging
 * @param {Object} customErrorMapping - Custom error mapping (optional)
 */
export const handleGrpcError = (
  res,
  error,
  correlationId,
  serviceName,
  methodName,
  customErrorMapping = {}
) => {
  // Log error
  logger.error(`${serviceName} ${methodName} error`, {
    error: error.message,
    code: error.code,
    correlationId,
  });

  // Get error mapping for the service
  const serviceErrorMapping = getErrorMapping(serviceName);

  // Merge with custom mapping (custom mapping takes precedence)
  const mapping = { ...serviceErrorMapping, ...customErrorMapping };

  // Clean gRPC error message - remove prefix like "9 FAILED_PRECONDITION: "
  let errorMessage = error.message || '';
  if (errorMessage.includes(': ')) {
    errorMessage = errorMessage.split(': ').slice(1).join(': ');
  }

  // Handle specific error messages for auth service
  if (serviceName.toLowerCase() === 'auth') {
    if (errorMessage.includes('PIN code has expired')) {
      return sendErrorResponse(res, 410, errorMessage, correlationId, 'PIN_CODE_EXPIRED');
    }
    if (errorMessage.includes('Email is already verified')) {
      return sendErrorResponse(res, 409, errorMessage, correlationId, 'EMAIL_ALREADY_VERIFIED');
    }
    if (errorMessage.includes('User not found')) {
      return sendErrorResponse(res, 404, errorMessage, correlationId, 'USER_NOT_FOUND');
    }
    if (errorMessage.includes('Invalid PIN code')) {
      return sendErrorResponse(res, 400, errorMessage, correlationId, 'INVALID_PIN_CODE');
    }
  }

  const errorInfo = mapping[error.code] || {
    status: 500,
    message: 'Internal Server Error',
    code: 'INTERNAL_ERROR',
  };

  sendErrorResponse(
    res,
    errorInfo.status,
    errorMessage || errorInfo.message,
    correlationId,
    errorInfo.code
  );
};

/**
 * Create a handler wrapper with common error handling
 * @param {Function} handler - The handler function
 * @param {string} serviceName - Service name for logging
 * @param {string} methodName - Method name for logging
 * @param {Object} customErrorMapping - Custom error mapping (optional)
 * @returns {Function} Wrapped handler function
 */
export const createHandler = (handler, serviceName, methodName, customErrorMapping = {}) => {
  return async (req, res) => {
    try {
      // Handle validation if needed
      if (!handleValidation(req, res)) {
        return;
      }

      // Call the actual handler
      await handler(req, res);
    } catch (error) {
      handleGrpcError(res, error, req.correlationId, serviceName, methodName, customErrorMapping);
    }
  };
};

/**
 * Create a simple handler wrapper (without validation)
 * @param {Function} handler - The handler function
 * @param {string} serviceName - Service name for logging
 * @param {string} methodName - Method name for logging
 * @param {Object} customErrorMapping - Custom error mapping (optional)
 * @returns {Function} Wrapped handler function
 */
export const createSimpleHandler = (handler, serviceName, methodName, customErrorMapping = {}) => {
  return async (req, res) => {
    try {
      await handler(req, res);
    } catch (error) {
      handleGrpcError(res, error, req.correlationId, serviceName, methodName, customErrorMapping);
    }
  };
};

import logger from '../utils/logger.js';

const errorHandler = (err, req, res, next) => {
  // Log error
  logger.logError(err, req);

  // Set correlation ID in response
  if (req.correlationId) {
    res.setHeader('X-Correlation-ID', req.correlationId);
  }

  const meta = {
    correlationId: req.correlationId,
    timestamp: new Date().toISOString(),
  };

  // Handle different types of errors
  if (err.name === 'ValidationError') {
    return res.status(400).json({
      error: {
        code: 'VALIDATION_ERROR',
        message: err.message,
        details: err.details || null,
      },
      meta,
    });
  }

  if (err.name === 'UnauthorizedError') {
    return res.status(401).json({
      error: {
        code: 'UNAUTHENTICATED',
        message: 'Invalid or missing authentication token',
      },
      meta,
    });
  }

  if (err.name === 'ForbiddenError') {
    return res.status(403).json({
      error: {
        code: 'FORBIDDEN',
        message: 'Access denied',
      },
      meta,
    });
  }

  if (err.name === 'NotFoundError') {
    return res.status(404).json({
      error: {
        code: 'NOT_FOUND',
        message: err.message || 'Resource not found',
      },
      meta,
    });
  }

  if (err.name === 'RateLimitError') {
    return res.status(429).json({
      error: {
        code: 'RATE_LIMIT_EXCEEDED',
        message: 'Rate limit exceeded',
      },
      meta,
    });
  }

  // Default error response
  const statusCode = err.statusCode || 500;
  const message = process.env.NODE_ENV === 'production'
    ? 'Something went wrong'
    : (err.message || 'Internal Server Error');

  const errorResponse = {
    error: {
      code: 'INTERNAL_ERROR',
      message,
    },
    meta,
  };

  if (process.env.NODE_ENV === 'development' && err.stack) {
    errorResponse.error.details = { stack: err.stack };
  }

  res.status(statusCode).json(errorResponse);
};

export default errorHandler;

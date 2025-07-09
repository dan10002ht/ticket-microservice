/**
 * Error Codes for Auth Service
 * Using string constants for better readability and maintainability
 */

export const ERROR_CODES = {
  // Validation errors (3xx)
  INVALID_ARGUMENT: 'INVALID_ARGUMENT',
  MISSING_REQUIRED_FIELD: 'MISSING_REQUIRED_FIELD',
  INVALID_EMAIL_FORMAT: 'INVALID_EMAIL_FORMAT',
  INVALID_PASSWORD_FORMAT: 'INVALID_PASSWORD_FORMAT',
  INVALID_PIN_CODE: 'INVALID_PIN_CODE',
  PIN_CODE_EXPIRED: 'PIN_CODE_EXPIRED',

  // Authentication errors (4xx)
  UNAUTHORIZED: 'UNAUTHORIZED',
  INVALID_CREDENTIALS: 'INVALID_CREDENTIALS',
  INVALID_TOKEN: 'INVALID_TOKEN',
  TOKEN_EXPIRED: 'TOKEN_EXPIRED',
  INVALID_REFRESH_TOKEN: 'INVALID_REFRESH_TOKEN',

  // Resource errors (5xx)
  USER_NOT_FOUND: 'USER_NOT_FOUND',
  EMAIL_ALREADY_EXISTS: 'EMAIL_ALREADY_EXISTS',
  EMAIL_ALREADY_VERIFIED: 'EMAIL_ALREADY_VERIFIED',
  EMAIL_NOT_VERIFIED: 'EMAIL_NOT_VERIFIED',
  ACCOUNT_LOCKED: 'ACCOUNT_LOCKED',
  ACCOUNT_DISABLED: 'ACCOUNT_DISABLED',

  // Business logic errors (6xx)
  EMAIL_VERIFICATION_FAILED: 'EMAIL_VERIFICATION_FAILED',
  PASSWORD_RESET_FAILED: 'PASSWORD_RESET_FAILED',
  SESSION_EXPIRED: 'SESSION_EXPIRED',
  TOO_MANY_ATTEMPTS: 'TOO_MANY_ATTEMPTS',

  // System errors (7xx)
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',
  DATABASE_ERROR: 'DATABASE_ERROR',
  EMAIL_SERVICE_ERROR: 'EMAIL_SERVICE_ERROR',
  REDIS_ERROR: 'REDIS_ERROR',
};

/**
 * Error Messages mapping
 */
export const ERROR_MESSAGES = {
  [ERROR_CODES.INVALID_ARGUMENT]: 'Invalid argument provided',
  [ERROR_CODES.MISSING_REQUIRED_FIELD]: 'Required field is missing',
  [ERROR_CODES.INVALID_EMAIL_FORMAT]: 'Invalid email format',
  [ERROR_CODES.INVALID_PASSWORD_FORMAT]: 'Invalid password format',
  [ERROR_CODES.INVALID_PIN_CODE]: 'Invalid PIN code',
  [ERROR_CODES.PIN_CODE_EXPIRED]: 'PIN code has expired',

  [ERROR_CODES.UNAUTHORIZED]: 'Unauthorized access',
  [ERROR_CODES.INVALID_CREDENTIALS]: 'Invalid email or password',
  [ERROR_CODES.INVALID_TOKEN]: 'Invalid token',
  [ERROR_CODES.TOKEN_EXPIRED]: 'Token has expired',
  [ERROR_CODES.INVALID_REFRESH_TOKEN]: 'Invalid refresh token',

  [ERROR_CODES.USER_NOT_FOUND]: 'User not found',
  [ERROR_CODES.EMAIL_ALREADY_EXISTS]: 'Email is already registered',
  [ERROR_CODES.EMAIL_ALREADY_VERIFIED]: 'Email is already verified',
  [ERROR_CODES.EMAIL_NOT_VERIFIED]: 'Email is not verified',
  [ERROR_CODES.ACCOUNT_LOCKED]: 'Account is locked',
  [ERROR_CODES.ACCOUNT_DISABLED]: 'Account is disabled',

  [ERROR_CODES.EMAIL_VERIFICATION_FAILED]: 'Email verification failed',
  [ERROR_CODES.PASSWORD_RESET_FAILED]: 'Password reset failed',
  [ERROR_CODES.SESSION_EXPIRED]: 'Session has expired',
  [ERROR_CODES.TOO_MANY_ATTEMPTS]: 'Too many attempts, please try again later',

  [ERROR_CODES.INTERNAL_ERROR]: 'Internal server error',
  [ERROR_CODES.SERVICE_UNAVAILABLE]: 'Service temporarily unavailable',
  [ERROR_CODES.DATABASE_ERROR]: 'Database error occurred',
  [ERROR_CODES.EMAIL_SERVICE_ERROR]: 'Email service error',
  [ERROR_CODES.REDIS_ERROR]: 'Cache service error',
};

/**
 * Custom Error class for Auth Service
 */
export class AuthError extends Error {
  constructor(errorCode, message = null, details = null) {
    super(message || ERROR_MESSAGES[errorCode] || 'Unknown error');
    this.name = 'AuthError';
    this.errorCode = errorCode;
    this.message = message || ERROR_MESSAGES[errorCode] || 'Unknown error';
    this.details = details;
  }
}

/**
 * Helper function to create AuthError instances
 */
export function createAuthError(errorCode, message = null, details = null) {
  return new AuthError(errorCode, message, details);
}

/**
 * Helper function to get gRPC error response
 */
export function getGrpcErrorResponse(error) {
  if (error instanceof AuthError) {
    return {
      code: getGrpcErrorCode(error.errorCode),
      message: error.message,
    };
  }

  return {
    code: 13, // INTERNAL
    message: error.message || 'Internal server error',
  };
}

/**
 * Map error codes to gRPC error codes
 */
function getGrpcErrorCode(errorCode) {
  const errorCodeMap = {
    // Validation errors -> INVALID_ARGUMENT (3)
    [ERROR_CODES.INVALID_ARGUMENT]: 3,
    [ERROR_CODES.MISSING_REQUIRED_FIELD]: 3,
    [ERROR_CODES.INVALID_EMAIL_FORMAT]: 3,
    [ERROR_CODES.INVALID_PASSWORD_FORMAT]: 3,
    [ERROR_CODES.INVALID_PIN_CODE]: 3,
    [ERROR_CODES.PIN_CODE_EXPIRED]: 3,

    // Authentication errors -> UNAUTHENTICATED (16)
    [ERROR_CODES.UNAUTHORIZED]: 16,
    [ERROR_CODES.INVALID_CREDENTIALS]: 16,
    [ERROR_CODES.INVALID_TOKEN]: 16,
    [ERROR_CODES.TOKEN_EXPIRED]: 16,
    [ERROR_CODES.INVALID_REFRESH_TOKEN]: 16,

    // Resource errors -> NOT_FOUND (5) or ALREADY_EXISTS (6)
    [ERROR_CODES.USER_NOT_FOUND]: 5,
    [ERROR_CODES.EMAIL_ALREADY_EXISTS]: 6,
    [ERROR_CODES.EMAIL_ALREADY_VERIFIED]: 9, // FAILED_PRECONDITION
    [ERROR_CODES.EMAIL_NOT_VERIFIED]: 9, // FAILED_PRECONDITION
    [ERROR_CODES.ACCOUNT_LOCKED]: 9, // FAILED_PRECONDITION
    [ERROR_CODES.ACCOUNT_DISABLED]: 9, // FAILED_PRECONDITION

    // Business logic errors -> FAILED_PRECONDITION (9)
    [ERROR_CODES.EMAIL_VERIFICATION_FAILED]: 9,
    [ERROR_CODES.PASSWORD_RESET_FAILED]: 9,
    [ERROR_CODES.SESSION_EXPIRED]: 9,
    [ERROR_CODES.TOO_MANY_ATTEMPTS]: 9,

    // System errors -> INTERNAL (13) or UNAVAILABLE (14)
    [ERROR_CODES.INTERNAL_ERROR]: 13,
    [ERROR_CODES.SERVICE_UNAVAILABLE]: 14,
    [ERROR_CODES.DATABASE_ERROR]: 13,
    [ERROR_CODES.EMAIL_SERVICE_ERROR]: 13,
    [ERROR_CODES.REDIS_ERROR]: 13,
  };

  return errorCodeMap[errorCode] || 13; // Default to INTERNAL
}

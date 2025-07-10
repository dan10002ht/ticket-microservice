import crypto from 'crypto';

/**
 * Generate a random token with specified length
 * @param {number} length - Length of the token in bytes (default: 32)
 * @param {string} encoding - Encoding format (default: 'hex')
 * @returns {string} Generated token
 */
export function generateToken(length = 32, encoding = 'hex') {
  return crypto.randomBytes(length).toString(encoding);
}

/**
 * Generate a random token and its hash
 * @param {number} length - Length of the token in bytes (default: 32)
 * @param {string} encoding - Encoding format (default: 'hex')
 * @param {string} algorithm - Hash algorithm (default: 'sha256')
 * @returns {Object} Object containing token and tokenHash
 */
export function generateTokenWithHash(length = 32, encoding = 'hex', algorithm = 'sha256') {
  const token = generateToken(length, encoding);
  const tokenHash = crypto.createHash(algorithm).update(token).digest(encoding);

  return {
    token,
    tokenHash,
  };
}

/**
 * Hash a token using specified algorithm
 * @param {string} token - Token to hash
 * @param {string} algorithm - Hash algorithm (default: 'sha256')
 * @param {string} encoding - Encoding format (default: 'hex')
 * @returns {string} Hashed token
 */
export function hashToken(token, algorithm = 'sha256', encoding = 'hex') {
  return crypto.createHash(algorithm).update(token).digest(encoding);
}

/**
 * Generate a secure random string for specific use cases
 * @param {string} type - Type of token ('reset', 'verification', 'api', 'session')
 * @returns {Object} Object containing token and tokenHash
 */
export function generateSecureToken(type = 'reset') {
  const tokenConfigs = {
    reset: { length: 32, algorithm: 'sha256' },
    verification: { length: 24, algorithm: 'sha256' },
    api: { length: 64, algorithm: 'sha256' },
    session: { length: 32, algorithm: 'sha256' },
  };

  const config = tokenConfigs[type] || tokenConfigs.reset;

  return generateTokenWithHash(config.length, 'hex', config.algorithm);
}

/**
 * Verify if a token matches its hash
 * @param {string} token - Original token
 * @param {string} tokenHash - Hashed token to compare against
 * @param {string} algorithm - Hash algorithm (default: 'sha256')
 * @param {string} encoding - Encoding format (default: 'hex')
 * @returns {boolean} True if token matches hash
 */
export function verifyToken(token, tokenHash, algorithm = 'sha256', encoding = 'hex') {
  const computedHash = hashToken(token, algorithm, encoding);
  return crypto.timingSafeEqual(
    Buffer.from(computedHash, encoding),
    Buffer.from(tokenHash, encoding)
  );
}

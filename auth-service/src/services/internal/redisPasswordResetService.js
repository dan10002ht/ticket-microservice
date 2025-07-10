import { generateSecureToken, hashToken } from '../../helpers/tokenHelper.js';
import {
  getUserRepository,
  getUserSessionRepository,
} from '../../repositories/repositoryFactory.js';
import {
  setPasswordResetToken,
  getPasswordResetToken,
  deletePasswordResetToken,
  incrementPasswordResetAttempts,
  getPasswordResetAttempts,
  clearPasswordResetAttempts,
  getKeys,
  healthCheck as redisHealthCheck,
} from '../redis/redisService.js';
import logger from '../../utils/logger.js';
import { getBackgroundService } from '../../background/backgroundService.js';
import { EMAIL_RESET_PASSWORD_JOB } from '../../const/background.js';

const userRepository = getUserRepository();
const userSessionRepository = getUserSessionRepository();
const backgroundService = getBackgroundService();

/**
 * Send forgot password email with Redis-based token storage
 */
export async function forgotPassword(email) {
  try {
    const user = await userRepository.findByEmail(email);

    if (!user) {
      return {
        message: 'If the email exists, a password reset link has been sent',
      };
    }

    backgroundService.enqueueJob(EMAIL_RESET_PASSWORD_JOB, {
      email: user.email,
      userId: user.id,
      resetToken: user.reset_token,
    });

    return {
      message: 'If the email exists, a password reset link has been sent',
    };
  } catch (error) {
    logger.error('Error in forgotPassword:', error);
    throw new Error(`Failed to send password reset email: ${error.message}`);
  }
}

/**
 * Reset password with Redis-based token validation
 */
export async function resetPassword(token, newPassword) {
  try {
    const tokenHash = hashToken(token);

    // Get token data from Redis
    const tokenData = await getPasswordResetToken(tokenHash);

    if (!tokenData) {
      throw new Error('Invalid or expired reset token');
    }

    // Update user password
    await userRepository.updatePassword(tokenData.user_id, newPassword);

    // Delete all user sessions to force re-login
    await userSessionRepository.deleteAllByUserId(tokenData.user_id);

    // Delete the used token
    await deletePasswordResetToken(tokenHash);

    // Clear attempts counter for this user
    await clearPasswordResetAttempts(tokenData.user_id);

    logger.info(`Password reset successful for user ${tokenData.user_id}`, {
      userId: tokenData.user_id,
      email: tokenData.email,
    });

    return {
      message: 'Password reset successfully',
    };
  } catch (error) {
    logger.error('Password reset failed:', error);
    throw new Error(`Password reset failed: ${error.message}`);
  }
}

/**
 * Validate password reset token (for checking if token is valid without resetting)
 */
export async function validatePasswordResetToken(token) {
  try {
    const tokenHash = hashToken(token);
    const tokenData = await getPasswordResetToken(tokenHash);

    if (!tokenData) {
      return { valid: false, message: 'Invalid or expired reset token' };
    }

    return {
      valid: true,
      userId: tokenData.user_id,
      email: tokenData.email,
      createdAt: tokenData.created_at,
    };
  } catch (error) {
    logger.error('Error validating password reset token:', error);
    return { valid: false, message: 'Token validation failed' };
  }
}

/**
 * Clean up expired tokens (Redis handles this automatically, but this can be used for manual cleanup)
 */
export async function cleanupExpiredTokens() {
  try {
    // Redis automatically handles TTL, but we can scan for any orphaned keys
    const pattern = 'pwd_reset:*';
    const keys = await getKeys(pattern);

    logger.info(`Found ${keys.length} password reset tokens in Redis`);

    return {
      totalTokens: keys.length,
      message: 'Redis automatically handles token expiration via TTL',
    };
  } catch (error) {
    logger.error('Error cleaning up expired tokens:', error);
    throw error;
  }
}

/**
 * Get Redis health status
 */
export async function getRedisHealth() {
  return await redisHealthCheck();
}

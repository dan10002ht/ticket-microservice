import * as redisPasswordResetService from './redisPasswordResetService.js';

// ========== PASSWORD RESET ==========

/**
 * Send forgot password email
 */
export async function forgotPassword(email) {
  try {
    return await redisPasswordResetService.forgotPassword(email);
  } catch (error) {
    console.error('❌ Error in forgotPassword:', error);
    throw new Error(`Failed to send password reset email: ${error.message}`);
  }
}

/**
 * Reset password with token
 */
export async function resetPassword(token, newPassword) {
  try {
    return await redisPasswordResetService.resetPassword(token, newPassword);
  } catch (error) {
    throw new Error(`Password reset failed: ${error.message}`);
  }
}

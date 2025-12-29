import { generateSecureToken } from '../../helpers/tokenHelper.js';
import { setPasswordResetToken } from '../../services/redis/redisService.js';
import { getBackgroundService } from '../backgroundService.js';
import { EMAIL_RESET_PASSWORD_JOB } from '../../const/background.js';
import logger from '../../utils/logger.js';
import { createAuthError, ERROR_CODES } from '../../utils/errorCodes.js';

const backgroundService = getBackgroundService();

/**
 * Background job handler for password reset
 * Generates secure token and sends reset email
 */
export async function handlePasswordResetJob(jobData) {
  try {
    const { email, userId } = jobData;

    if (!email || !userId) {
      throw createAuthError(
        ERROR_CODES.MISSING_REQUIRED_FIELD,
        'Email and userId are required for password reset'
      );
    }

    logger.info('Processing password reset job', {
      userId,
      email,
      jobId: jobData.jobId,
    });

    // Generate secure reset token
    const resetToken = generateSecureToken();
    const tokenHash = resetToken.hash;

    // Store token in Redis with user data
    const tokenData = {
      user_id: userId,
      email: email,
      created_at: new Date().toISOString(),
    };

    await setPasswordResetToken(tokenHash, tokenData);

    // Create forgot password URL
    const forgotPasswordUrl = `${process.env.FRONTEND_URL || 'http://localhost:53000'}/forgot-password?user_id=${userId}&token=${resetToken.token}`;

    // Enqueue email job to send reset email
    await backgroundService.enqueueJob(EMAIL_RESET_PASSWORD_JOB, {
      email: email,
      userId: userId,
      forgotPasswordUrl: forgotPasswordUrl,
    });

    logger.info('Password reset job completed successfully', {
      userId,
      email,
      jobId: jobData.jobId,
    });

    return {
      success: true,
      message: 'Password reset email queued successfully',
      userId,
      email,
    };
  } catch (error) {
    logger.error('Password reset job failed:', error, {
      jobData,
      jobId: jobData.jobId,
    });

    throw createAuthError(
      ERROR_CODES.PASSWORD_RESET_FAILED,
      `Password reset job failed: ${error.message}`
    );
  }
}

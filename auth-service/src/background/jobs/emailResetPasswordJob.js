import logger from '../../utils/logger.js';
import { grpcClients } from '../../grpc/clients.js';
import { PASSWORD_RESET_CONFIG } from '../../const/background.js';
import {
  deletePasswordResetToken,
  getPasswordResetAttempts,
  incrementPasswordResetAttempts,
  setPasswordResetToken,
} from '../../services/redis/redisService.js';
import { generateSecureToken } from '../../helpers/tokenHelper.js';

export async function handleEmailResetPasswordJob(jobData) {
  const { email, userId } = jobData;
  const { token: resetToken, tokenHash } = generateSecureToken('reset');
  try {
    logger.info(`Processing email reset password job for email: ${email}`);
    // Check rate limiting for this user
    const currentAttempts = await getPasswordResetAttempts(userId);

    if (currentAttempts >= PASSWORD_RESET_CONFIG.MAX_ATTEMPTS) {
      logger.warn(`Password reset rate limit exceeded for user ${userId}`);
      return {
        message: 'If the email exists, a password reset link has been sent',
      };
    }

    // Store token in Redis with TTL
    const tokenData = {
      user_id: userId,
      email: email,
      created_at: new Date().toISOString(),
    };

    await setPasswordResetToken(tokenHash, tokenData);

    // Increment attempts counter
    await incrementPasswordResetAttempts(userId);

    logger.info(`Password reset token created for user ${userId}`, {
      userId: userId,
      email: email,
      tokenExpiresIn: PASSWORD_RESET_CONFIG.TTL,
    });
    await sendPasswordResetEmailViaGrpc({
      email: email,
      resetToken: resetToken,
    });
  } catch (error) {
    logger.error('Error in handleEmailResetPasswordJob:', error);
    try {
      await deletePasswordResetToken(tokenHash);
    } catch (redisError) {
      logger.error(`Failed to remove password reset token from Redis:`, redisError);
    }
    throw error;
  }
}

export async function sendPasswordResetEmailViaGrpc(data) {
  const { email, resetToken } = data;
  const emailWorkerClient = grpcClients.emailWorkerClient;
  const response = await emailWorkerClient.sendPasswordResetEmail(data);
  // todo @dantt
  return response;
}

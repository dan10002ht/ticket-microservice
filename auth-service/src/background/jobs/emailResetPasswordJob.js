import logger from '../../utils/logger.js';
import { grpcClients } from '../../grpc/clients.js';
import { PASSWORD_RESET_CONFIG } from '../../const/background.js';
import {
  getPasswordResetAttempts,
  incrementPasswordResetAttempts,
} from '../../services/redis/redisService.js';

export async function handleEmailResetPasswordJob(jobData) {
  const { email, userId, forgotPasswordUrl } = jobData;

  try {
    logger.info(`Processing email reset password job for email: ${email}`);

    // Check rate limiting for this user
    const currentAttempts = await getPasswordResetAttempts(userId);

    if (currentAttempts >= PASSWORD_RESET_CONFIG.MAX_ATTEMPTS) {
      logger.warn(`Password reset rate limit exceeded for user ${userId}`);
      return;
    }

    // Increment attempts counter
    await incrementPasswordResetAttempts(userId);

    logger.info(`Password reset email job for user ${userId}`, {
      userId: userId,
      email: email,
    });

    await sendPasswordResetEmailViaGrpc({
      email,
      userId,
      forgotPasswordUrl,
    });
  } catch (error) {
    logger.error('Error in handleEmailResetPasswordJob:', error);
    throw error;
  }
}

export async function sendPasswordResetEmailViaGrpc(data) {
  try {
    const { email, userId, forgotPasswordUrl } = data;
    const passwordResetData = {
      email,
      user_id: userId,
      forgot_password_url: forgotPasswordUrl,
    };
    const response = await grpcClients.emailService.sendPasswordResetEmail(passwordResetData);
    logger.info(`Password reset email sent successfully to ${email}`, response);
  } catch (e) {
    logger.error('Failed to send password reset email via gRPC:', e.message);
    throw new Error(`Email sending failed: ${e.message}`);
  }
}

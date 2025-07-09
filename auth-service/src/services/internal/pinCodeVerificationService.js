import { getUserRepository } from '../../repositories/repositoryFactory.js';
import { sanitizeUserForResponse } from '../../utils/sanitizers.js';
import { validatePinCodeFromRedis } from '../../background/jobs/emailVerificationJob.js';
import { getBackgroundService } from '../../background/backgroundService.js';
import { EMAIL_VERIFICATION_JOB, JOB_RETRY_CONFIGS } from '../../const/background.js';

const userRepository = getUserRepository();

// ========== PIN CODE VERIFICATION ==========

/**
 * Send verification email with PIN code
 */
export async function sendVerificationEmailWithPin(email) {
  try {
    const user = await userRepository.findByEmail(email);
    if (!user) {
      throw new Error('User not found');
    }

    if (user.is_verified) {
      throw new Error('Email is already verified');
    }

    // Enqueue background job to send email (PIN will be generated in background)
    const backgroundService = getBackgroundService();
    await backgroundService.enqueueJob(
      EMAIL_VERIFICATION_JOB,
      {
        userId: user.id,
        userEmail: user.email,
        userName: user.first_name || user.email,
      },
      JOB_RETRY_CONFIGS.EMAIL_OPERATIONS
    );

    // Return data for background processing
    return {
      message: 'Verification email queued successfully',
      userId: user.id,
      userEmail: user.email,
    };
  } catch (error) {
    throw new Error(`Failed to send verification email: ${error.message}`);
  }
}

/**
 * Verify email with PIN code
 */
export async function verifyEmailWithPin(userId, inputPinCode) {
  try {
    // Get user
    const user = await userRepository.findById(userId);
    if (!user) {
      throw new Error('User not found');
    }

    if (user.is_verified) {
      throw new Error('Email is already verified');
    }

    // Check PIN code from Redis
    const validationResult = await validatePinCodeFromRedis(userId, inputPinCode);

    if (!validationResult.valid) {
      throw new Error(validationResult.message);
    }

    // Update user verification status
    const updatedUser = await userRepository.updateUser(userId, {
      is_verified: true,
      email_verified_at: new Date(),
      updated_at: new Date(),
    });

    return {
      message: 'Email verified successfully',
      user: sanitizeUserForResponse(updatedUser),
    };
  } catch (error) {
    throw new Error(`Email verification failed: ${error.message}`);
  }
}

/**
 * Resend verification email with new PIN code
 */
export async function resendVerificationEmail(email) {
  try {
    const user = await userRepository.findByEmail(email);
    if (!user) {
      throw new Error('User not found');
    }

    if (user.is_verified) {
      throw new Error('Email is already verified');
    }

    // Enqueue background job to send email (new PIN will be generated in background)
    const backgroundService = getBackgroundService();
    await backgroundService.enqueueJob(
      EMAIL_VERIFICATION_JOB,
      {
        userId: user.id,
        userEmail: user.email,
        userName: user.first_name || user.email,
        isResend: true,
      },
      JOB_RETRY_CONFIGS.EMAIL_OPERATIONS
    );

    return {
      message: 'Verification email resent successfully',
      userId: user.id,
      userEmail: user.email,
    };
  } catch (error) {
    throw new Error(`Failed to resend verification email: ${error.message}`);
  }
}

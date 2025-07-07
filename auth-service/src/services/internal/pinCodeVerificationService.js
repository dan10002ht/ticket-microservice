import { getUserRepository } from '../../repositories/repositoryFactory.js';
import { sanitizeUserForResponse } from '../../utils/sanitizers.js';
import { validatePinCodeFromRedis } from '../../background/jobs/emailVerificationJob.js';
import { getBackgroundService } from '../../background/backgroundService.js';

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
      'email_verification',
      {
        userId: user.id,
        userEmail: user.email,
        userName: user.first_name || user.email,
      },
      {
        priority: 'high',
        maxRetries: 3,
        timeout: 30000, // 30 seconds
      }
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
      'email_verification',
      {
        userId: user.id,
        userEmail: user.email,
        userName: user.first_name || user.email,
        isResend: true,
      },
      {
        priority: 'high',
        maxRetries: 3,
        timeout: 30000, // 30 seconds
      }
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

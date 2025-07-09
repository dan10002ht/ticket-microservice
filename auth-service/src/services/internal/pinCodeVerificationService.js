import { getUserRepository } from '../../repositories/repositoryFactory.js';
import { sanitizeUserForResponse } from '../../utils/sanitizers.js';
import { validatePinCodeFromRedis } from '../../background/jobs/emailVerificationJob.js';
import { getBackgroundService } from '../../background/backgroundService.js';
import { EMAIL_VERIFICATION_JOB, JOB_RETRY_CONFIGS } from '../../const/background.js';
import { createAuthError, ERROR_CODES } from '../../utils/errorCodes.js';

const userRepository = getUserRepository();

// ========== PIN CODE VERIFICATION ==========

/**
 * Send verification email with PIN code
 */
export async function sendVerificationEmailWithPin(email) {
  try {
    const user = await userRepository.findByEmail(email);
    if (!user) {
      throw createAuthError(ERROR_CODES.USER_NOT_FOUND);
    }

    if (user.is_verified) {
      throw createAuthError(ERROR_CODES.EMAIL_ALREADY_VERIFIED);
    }

    // Enqueue background job to send email (PIN will be generated in background)
    const backgroundService = getBackgroundService();
    await backgroundService.enqueueJob(
      EMAIL_VERIFICATION_JOB,
      {
        userId: user.public_id,
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
    if (error.name === 'AuthError') {
      throw error;
    }

    // Wrap other errors
    throw createAuthError(ERROR_CODES.EMAIL_VERIFICATION_FAILED, error.message);
  }
}

/**
 * Verify email with PIN code
 */
export async function verifyEmailWithPin(userId, inputPinCode) {
  try {
    // Get user
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      throw createAuthError(ERROR_CODES.USER_NOT_FOUND);
    }

    if (user.is_verified) {
      throw createAuthError(ERROR_CODES.EMAIL_ALREADY_VERIFIED);
    }

    // Check PIN code from Redis
    const validationResult = await validatePinCodeFromRedis(userId, inputPinCode);

    if (!validationResult.valid) {
      // Map validation errors to appropriate error codes
      if (validationResult.message.includes('expired')) {
        throw createAuthError(ERROR_CODES.PIN_CODE_EXPIRED);
      }
      throw createAuthError(ERROR_CODES.INVALID_PIN_CODE, validationResult.message);
    }

    // Update user verification status
    const updatedUser = await userRepository.updateUser(user.id, {
      is_verified: true,
      email_verified_at: new Date(),
      updated_at: new Date(),
    });

    return {
      message: 'Email verified successfully',
      user: sanitizeUserForResponse(updatedUser),
    };
  } catch (error) {
    // Re-throw AuthError as is
    if (error.name === 'AuthError') {
      throw error;
    }

    // Wrap other errors
    throw createAuthError(ERROR_CODES.EMAIL_VERIFICATION_FAILED, error.message);
  }
}

/**
 * Resend verification email with new PIN code
 */
export async function resendVerificationEmail(email) {
  try {
    const user = await userRepository.findByEmail(email);
    if (!user) {
      throw createAuthError(ERROR_CODES.USER_NOT_FOUND);
    }

    if (user.is_verified) {
      throw createAuthError(ERROR_CODES.EMAIL_ALREADY_VERIFIED);
    }

    // Enqueue background job to send email (new PIN will be generated in background)
    const backgroundService = getBackgroundService();
    await backgroundService.enqueueJob(
      EMAIL_VERIFICATION_JOB,
      {
        userId: user.public_id,
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
    // Re-throw AuthError as is
    if (error.name === 'AuthError') {
      throw error;
    }

    // Wrap other errors
    throw createAuthError(ERROR_CODES.EMAIL_VERIFICATION_FAILED, error.message);
  }
}

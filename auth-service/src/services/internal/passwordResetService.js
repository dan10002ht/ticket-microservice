import { generateSecureToken, hashToken } from '../../helpers/tokenHelper.js';
import {
  getUserRepository,
  getPasswordResetTokenRepository,
  getUserSessionRepository,
} from '../../repositories/repositoryFactory.js';

const userRepository = getUserRepository();
const passwordResetTokenRepository = getPasswordResetTokenRepository();
const userSessionRepository = getUserSessionRepository();

// ========== PASSWORD RESET ==========

/**
 * Send forgot password email
 */
export async function forgotPassword(email) {
  try {
    console.log('🔍 Looking for user with email:', email);
    const user = await userRepository.findByEmail(email);
    console.log('👤 User found:', user ? 'Yes' : 'No');

    if (!user) {
      return {
        message: 'If the email exists, a password reset link has been sent',
      };
    }

    console.log('🔑 Generating secure token...');
    const { token: resetToken, tokenHash } = generateSecureToken('reset');
    console.log('✅ Token generated successfully');

    console.log('💾 Creating password reset token in database...');
    await passwordResetTokenRepository.createPasswordResetToken({
      user_id: user.id,
      token_hash: tokenHash,
      expires_at: new Date(Date.now() + 1 * 60 * 60 * 1000), // 1 hour
    });
    console.log('✅ Password reset token created');

    // TODO: Send email via email service
    // await emailService.sendPasswordResetEmail(user.email, resetToken);

    return {
      message: 'If the email exists, a password reset link has been sent',
      reset_token: resetToken, // For testing only
    };
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
    const tokenHash = hashToken(token);

    // Find valid reset token
    const resetToken = await passwordResetTokenRepository
      .findValid()
      .where('token_hash', tokenHash)
      .first();

    if (!resetToken) {
      throw new Error('Invalid or expired reset token');
    }

    // Mark token as used
    await passwordResetTokenRepository.markAsUsed(resetToken.id);

    // Update user password
    await userRepository.updatePassword(resetToken.user_id, newPassword);

    // Delete all user sessions to force re-login
    await userSessionRepository.deleteAllByUserId(resetToken.user_id);

    return {
      message: 'Password reset successfully',
    };
  } catch (error) {
    throw new Error(`Password reset failed: ${error.message}`);
  }
}

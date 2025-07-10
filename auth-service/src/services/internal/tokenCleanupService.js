import {
  getRefreshTokenRepository,
  getEmailVerificationTokenRepository,
} from '../../repositories/repositoryFactory.js';

/**
 * Token Cleanup Service
 * Xử lý cleanup tất cả các loại tokens expired
 */
class TokenCleanupService {
  constructor() {
    this.refreshTokenRepo = getRefreshTokenRepository();
    this.emailVerificationTokenRepo = getEmailVerificationTokenRepository();
  }

  /**
   * Cleanup tất cả expired tokens
   */
  async cleanupAllExpiredTokens() {
    try {
      const results = await Promise.all([
        this.refreshTokenRepo.deleteExpired(),
        this.emailVerificationTokenRepo.deleteExpired(),
      ]);

      return {
        refreshTokens: results[0],
        emailVerificationTokens: results[1],
        totalDeleted: results.reduce((sum, count) => sum + count, 0),
      };
    } catch (error) {
      console.error('Error cleaning up expired tokens:', error);
      throw error;
    }
  }

  /**
   * Cleanup chỉ refresh tokens
   */
  async cleanupExpiredRefreshTokens() {
    try {
      return await this.refreshTokenRepo.deleteExpired();
    } catch (error) {
      console.error('Error cleaning up expired refresh tokens:', error);
      throw error;
    }
  }

  /**
   * Cleanup chỉ email verification tokens
   */
  async cleanupExpiredEmailVerificationTokens() {
    try {
      return await this.emailVerificationTokenRepo.deleteExpired();
    } catch (error) {
      console.error('Error cleaning up expired email verification tokens:', error);
      throw error;
    }
  }
}

export default TokenCleanupService;

import { getRefreshTokenRepository } from '../../repositories/repositoryFactory.js';

/**
 * Token Cleanup Service
 * Xử lý cleanup tất cả các loại tokens expired
 */
class TokenCleanupService {
  constructor() {
    this.refreshTokenRepo = getRefreshTokenRepository();
  }

  /**
   * Cleanup tất cả expired tokens
   */
  async cleanupAllExpiredTokens() {
    try {
      const deletedCount = await this.refreshTokenRepo.deleteExpired();

      return {
        refreshTokens: deletedCount,
        totalDeleted: deletedCount,
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
}

export default TokenCleanupService;

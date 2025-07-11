import { getBackgroundService } from './backgroundService.js';
import { handleEmailVerificationJob } from './jobs/emailVerificationJob.js';
import {
  handleCacheUserProfileJob,
  handleCacheUserRolesJob,
  handleCacheUserDataJob,
} from './jobs/cacheJob.js';
import logger from './logger.js';
import {
  EMAIL_VERIFICATION_JOB,
  CACHE_USER_DATA_JOB,
  CACHE_USER_PROFILE_JOB,
  CACHE_USER_ROLES_JOB,
  EMAIL_RESET_PASSWORD_JOB,
} from '../const/background.js';
import { handleEmailResetPasswordJob } from './jobs/emailResetPasswordJob.js';

/**
 * Register all job handlers for the background service
 */
export async function registerJobHandlers() {
  try {
    const backgroundService = getBackgroundService();

    // Register email verification job handler
    backgroundService.registerHandler(EMAIL_VERIFICATION_JOB, handleEmailVerificationJob);
    logger.info('Registered email_verification job handler');
    backgroundService.registerHandler(EMAIL_RESET_PASSWORD_JOB, handleEmailResetPasswordJob);
    logger.info('Registered email_reset_password job handler');

    // Register cache job handlers
    backgroundService.registerHandler(CACHE_USER_PROFILE_JOB, handleCacheUserProfileJob);
    backgroundService.registerHandler(CACHE_USER_ROLES_JOB, handleCacheUserRolesJob);
    backgroundService.registerHandler(CACHE_USER_DATA_JOB, handleCacheUserDataJob);
    logger.info('Registered cache job handlers');

    // Register other job handlers here as needed
    // backgroundService.registerHandler('other_job_type', otherJobHandler);

    logger.info('All job handlers registered successfully');
  } catch (error) {
    logger.error('Failed to register job handlers:', error);
    throw error;
  }
}

/**
 * Initialize background service with job handlers
 */
export async function initializeBackgroundService() {
  try {
    const backgroundService = getBackgroundService();

    // Initialize background service
    await backgroundService.initialize();

    // Register job handlers
    await registerJobHandlers();

    logger.info('Background service initialized with job handlers');
  } catch (error) {
    logger.error('Failed to initialize background service:', error);
    throw error;
  }
}

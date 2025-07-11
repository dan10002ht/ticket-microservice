/**
 * Background Job Constants
 *
 * This file contains all background job names used in the auth service.
 * Centralizing these constants makes the code more maintainable and reduces
 * the risk of typos when referencing job names.
 */

// ========== USER MANAGEMENT JOBS ==========
export const CACHE_USER_DATA_JOB = 'cache_user_data';
export const CACHE_USER_PROFILE_JOB = 'cache_user_profile';
export const CACHE_USER_ROLES_JOB = 'cache_user_roles';
export const EMAIL_VERIFICATION_JOB = 'email_verification';

// ========== PASSWORD RESET JOBS ==========
export const EMAIL_RESET_PASSWORD_JOB = 'email_reset_password';

// ========== SESSION MANAGEMENT JOBS ==========
export const CLEANUP_EXPIRED_SESSIONS_JOB = 'cleanup_expired_sessions';
export const CLEANUP_EXPIRED_TOKENS_JOB = 'cleanup_expired_tokens';

// ========== AUDIT & LOGGING JOBS ==========
export const AUDIT_LOG_JOB = 'audit_log';
export const SECURITY_LOG_JOB = 'security_log';

// ========== NOTIFICATION JOBS ==========
export const SEND_WELCOME_EMAIL_JOB = 'send_welcome_email';
export const SEND_PASSWORD_RESET_JOB = 'send_password_reset';
export const SEND_EMAIL_VERIFICATION_JOB = 'send_email_verification';

// ========== MAINTENANCE JOBS ==========
export const CLEANUP_OLD_AUDIT_LOGS_JOB = 'cleanup_old_audit_logs';
export const DATABASE_MAINTENANCE_JOB = 'database_maintenance';
export const CACHE_CLEANUP_JOB = 'cache_cleanup';

// ========== JOB PRIORITIES ==========
export const JOB_PRIORITIES = {
  LOW: 'low',
  NORMAL: 'normal',
  HIGH: 'high',
  CRITICAL: 'critical',
};

// ========== JOB TIMEOUTS (in milliseconds) ==========
export const JOB_TIMEOUTS = {
  CACHE_OPERATIONS: 15000, // 15 seconds
  EMAIL_OPERATIONS: 30000, // 30 seconds
  DATABASE_OPERATIONS: 60000, // 1 minute
  CLEANUP_OPERATIONS: 120000, // 2 minutes
};

// ========== JOB RETRY CONFIGURATIONS ==========
export const JOB_RETRY_CONFIGS = {
  CACHE_OPERATIONS: {
    maxRetries: 2,
    timeout: JOB_TIMEOUTS.CACHE_OPERATIONS,
    priority: JOB_PRIORITIES.NORMAL,
  },
  EMAIL_OPERATIONS: {
    maxRetries: 3,
    timeout: JOB_TIMEOUTS.EMAIL_OPERATIONS,
    priority: JOB_PRIORITIES.HIGH,
  },
  CLEANUP_OPERATIONS: {
    maxRetries: 1,
    timeout: JOB_TIMEOUTS.CLEANUP_OPERATIONS,
    priority: JOB_PRIORITIES.LOW,
  },
};

// Configuration
export const PASSWORD_RESET_CONFIG = {
  MAX_ATTEMPTS: 3,
  TTL: 15,
};

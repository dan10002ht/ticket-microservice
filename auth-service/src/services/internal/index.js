/**
 * Internal Services Index
 * Business logic services that handle internal operations
 */

// Core Authentication Service
export * as authService from './authService.js';

// User Management Service
export * as userManagementService from './userManagementService.js';

// Admin Service
export * as adminService from './adminService.js';

// Organization Management Service
export * as organizationManagementService from './organizationManagementService.js';

// OAuth Service
export * as oauthService from './oauthService.js';

// Email Verification Service
// emailVerificationService removed - using pinCodeVerificationService instead

// Password Reset Service
export * as passwordResetService from './passwordResetService.js';

// Permission Service
export * as permissionService from './permissionService.js';

// Two-Factor Authentication Service
export * as twoFactorService from './twoFactorService.js';

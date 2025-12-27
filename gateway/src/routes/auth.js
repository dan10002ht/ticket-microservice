import express from 'express';
import * as authHandler from '../handlers/authHandlers.js';
import * as validation from '../middlewares/index.js';
import * as authorizationMiddleware from '../middlewares/authorizationMiddleware.js';
import { requireAuth } from '../middlewares/auth.js';

const router = express.Router();

router.post(
  '/register/email',
  validation.validateRegistration,
  authHandler.registerWithEmailHandler
);
router.post(
  '/register/oauth',
  validation.validateOAuthRegistration,
  authHandler.registerWithOAuthHandler
);
router.post('/login', validation.validateLogin, authHandler.loginHandler);
router.post(
  '/forgot-password',
  validation.validateForgotPassword,
  authHandler.forgotPasswordHandler
);
router.post('/reset-password', validation.validateResetPassword, authHandler.resetPasswordHandler);
router.post(
  '/send-verification-email',
  validation.validateSendVerificationEmail,
  authHandler.sendVerificationEmailHandler
);
router.post(
  '/verify-user',
  validation.validateVerifyEmailWithPin,
  authHandler.verifyEmailWithPinHandler
);
router.post(
  '/resend-verification-email',
  validation.validateResendVerificationEmail,
  authHandler.resendVerificationEmailHandler
);
router.post('/refresh', validation.validateRefreshToken, authHandler.refreshTokenHandler);
router.post('/logout', authHandler.logoutHandler);

// ============================================
// Token Validation & OAuth
// ============================================
router.post('/validate', validation.validateTokenValidation, authHandler.validateTokenHandler);
router.post('/oauth/login', validation.validateOAuthLogin, authHandler.oAuthLoginHandler);
router.post(
  '/verify-email-token',
  validation.validateVerifyEmailToken,
  authHandler.verifyEmailHandler
);
router.post('/register', validation.validateRegistration, authHandler.registerHandler);

// ============================================
// Permissions & Roles (requires auth)
// ============================================
router.get(
  '/permissions',
  requireAuth,
  authHandler.getUserPermissionsHandler
);
router.post(
  '/permissions/check',
  requireAuth,
  validation.validateCheckPermission,
  authHandler.checkPermissionHandler
);
router.post(
  '/permissions/resource',
  requireAuth,
  validation.validateCheckResourcePermission,
  authHandler.checkResourcePermissionHandler
);
router.post(
  '/permissions/batch',
  requireAuth,
  validation.validateBatchCheckPermissions,
  authHandler.batchCheckPermissionsHandler
);
router.get('/roles', requireAuth, authHandler.getUserRolesHandler);

// ============================================
// User Management (Admin only)
// ============================================
router.get(
  '/users/:userId',
  authorizationMiddleware.requireRole(['admin', 'super_admin']),
  authHandler.getUserHandler
);
router.put(
  '/users/:userId',
  authorizationMiddleware.requireRole(['admin', 'super_admin']),
  authHandler.updateUserHandler
);
router.delete(
  '/users/:userId',
  authorizationMiddleware.requireRole(['admin', 'super_admin']),
  authHandler.deleteUserHandler
);

// ============================================
// Health Check
// ============================================
router.get('/health', authHandler.healthCheckHandler);

// Public hello endpoint (no authorization required)
router.get('/hello', (req, res) => {
  res.send('Hello World');
});

// Protected hello endpoint (requires authentication and permission)
router.get(
  '/hello/protected',
  authorizationMiddleware.requirePermission('auth.hello.view'),
  (req, res) => {
    res.json({
      message: 'Hello Protected World',
      user: req.user,
    });
  }
);

// Admin hello endpoint (requires admin role)
router.get(
  '/hello/admin',
  authorizationMiddleware.requireRole(['admin', 'super_admin']),
  (req, res) => {
    res.json({
      message: 'Hello Admin World',
      user: req.user,
    });
  }
);

export default router;

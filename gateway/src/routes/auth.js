import express from 'express';
import {
  registerWithEmailHandler,
  registerWithOAuthHandler,
  loginHandler,
  refreshTokenHandler,
  logoutHandler,
  forgotPasswordHandler,
  resetPasswordHandler,
  sendVerificationEmailHandler,
  verifyEmailWithPinHandler,
  resendVerificationEmailHandler,
} from '../handlers/index.js';
import {
  validateRegistration,
  validateLogin,
  validateRefreshToken,
  validateOAuthRegistration,
  validateForgotPassword,
  validateResetPassword,
  validateSendVerificationEmail,
  validateVerifyEmailWithPin,
  validateResendVerificationEmail,
} from '../middlewares/index.js';
import { requirePermission, requireRole } from '../middlewares/authorizationMiddleware.js';

const router = express.Router();

router.post('/register/email', validateRegistration, registerWithEmailHandler);
router.post('/register/oauth', validateOAuthRegistration, registerWithOAuthHandler);
router.post('/login', validateLogin, loginHandler);
router.post('/forgot-password', validateForgotPassword, forgotPasswordHandler);
router.post('/reset-password', validateResetPassword, resetPasswordHandler);
router.post(
  '/send-verification-email',
  validateSendVerificationEmail,
  sendVerificationEmailHandler
);
router.post('/verify-user', validateVerifyEmailWithPin, verifyEmailWithPinHandler);
router.post(
  '/resend-verification-email',
  validateResendVerificationEmail,
  resendVerificationEmailHandler
);
router.post('/refresh', validateRefreshToken, refreshTokenHandler);
router.post('/logout', logoutHandler);

// Public hello endpoint (no authorization required)
router.get('/hello', (req, res) => {
  res.send('Hello World');
});

// Protected hello endpoint (requires authentication and permission)
router.get('/hello/protected', requirePermission('auth.hello.view'), (req, res) => {
  res.json({
    message: 'Hello Protected World',
    user: req.user,
  });
});

// Admin hello endpoint (requires admin role)
router.get('/hello/admin', requireRole(['admin', 'super_admin']), (req, res) => {
  res.json({
    message: 'Hello Admin World',
    user: req.user,
  });
});

export default router;

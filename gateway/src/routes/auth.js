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

const router = express.Router();

router.post('/register/email', validateRegistration, registerWithEmailHandler);
router.post('/register/oauth', validateOAuthRegistration, registerWithOAuthHandler);
router.post('/register/test/google/callback', (req, res) => {
  return res.status(200).json(req.body);
});
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

router.get('/hello', (req, res) => {
  res.send('Hello World');
});

export default router;

import grpcClients from '../grpc/clients.js';
import {
  sendSuccessResponse,
  createHandler,
  createSimpleHandler,
} from '../utils/responseHandler.js';

/**
 * Register a new user with email and password
 */
const registerUserWithEmail = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.registerWithEmail(requestData);
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Register a new user with OAuth (Google, Facebook, etc.)
 */
const registerUserWithOAuth = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.registerWithOAuth(requestData);
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * User login
 */
const loginUser = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.login(requestData);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Refresh access token
 */
const refreshUserToken = async (req, res) => {
  const result = await grpcClients.authService.refreshToken(req.body);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * User logout
 */
const logoutUser = async (req, res) => {
  await grpcClients.authService.logout({
    refresh_token: req.body.refresh_token,
  });
  sendSuccessResponse(res, 200, { message: 'Logout successful' }, req.correlationId);
};

/**
 * Forgot password - send reset email
 */
const forgotPassword = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.forgotPassword(requestData);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Reset password with token
 */
const resetPassword = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.resetPassword(requestData);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Send verification email with PIN code
 */
const sendVerificationEmail = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.sendVerificationEmail(requestData);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Verify email with PIN code
 */
const verifyEmailWithPin = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.verifyEmailWithPin(requestData);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Resend verification email
 */
const resendVerificationEmail = async (req, res) => {
  const requestData = {
    ...req.body,
    ip_address: req.ip || req.connection.remoteAddress,
    user_agent: req.get('User-Agent'),
  };

  const result = await grpcClients.authService.resendVerificationEmail(requestData);
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Token Validation
// ============================================

/**
 * Validate JWT token
 */
const validateToken = async (req, res) => {
  const token = req.body.token || req.headers.authorization?.replace('Bearer ', '');
  const result = await grpcClients.authService.ValidateToken({ token });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// OAuth Login
// ============================================

/**
 * OAuth login flow
 */
const oAuthLogin = async (req, res) => {
  const result = await grpcClients.authService.OAuthLogin({
    provider: req.body.provider,
    access_token: req.body.access_token,
    provider_user_id: req.body.provider_user_id,
    email: req.body.email,
    first_name: req.body.first_name,
    last_name: req.body.last_name,
    picture: req.body.picture,
    refresh_token: req.body.refresh_token,
    expires_at: req.body.expires_at,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Email Verification (Token-based)
// ============================================

/**
 * Verify email with token
 */
const verifyEmail = async (req, res) => {
  const result = await grpcClients.authService.VerifyEmail({
    token: req.body.token || req.query.token,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Permissions & Roles
// ============================================

/**
 * Get user permissions
 */
const getUserPermissions = async (req, res) => {
  const result = await grpcClients.authService.GetUserPermissions({
    user_id: req.user.id,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Check single permission
 */
const checkPermission = async (req, res) => {
  const result = await grpcClients.authService.CheckPermission({
    user_id: req.user.id,
    permission_name: req.body.permission_name,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Check resource permission
 */
const checkResourcePermission = async (req, res) => {
  const result = await grpcClients.authService.CheckResourcePermission({
    user_id: req.user.id,
    resource: req.body.resource,
    action: req.body.action,
    context: req.body.context || {},
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get user roles
 */
const getUserRoles = async (req, res) => {
  const result = await grpcClients.authService.GetUserRoles({
    user_id: req.user.id,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Batch check permissions
 */
const batchCheckPermissions = async (req, res) => {
  const result = await grpcClients.authService.BatchCheckPermissions({
    user_id: req.user.id,
    permission_names: req.body.permission_names,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// User Management (Admin)
// ============================================

/**
 * Get user by ID (Admin)
 */
const getUser = async (req, res) => {
  const result = await grpcClients.authService.GetUser({
    user_id: req.params.userId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update user (Admin)
 */
const updateUser = async (req, res) => {
  const result = await grpcClients.authService.UpdateUser({
    user_id: req.params.userId,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Delete user (Admin)
 */
const deleteUser = async (req, res) => {
  const result = await grpcClients.authService.DeleteUser({
    user_id: req.params.userId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Utility
// ============================================

/**
 * Basic registration
 */
const register = async (req, res) => {
  const result = await grpcClients.authService.Register(req.body);
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Health check
 */
const healthCheck = async (req, res) => {
  const result = await grpcClients.authService.Health({});
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// Export wrapped handlers
export const registerWithEmailHandler = createHandler(
  registerUserWithEmail,
  'auth',
  'registerWithEmail'
);
export const registerWithOAuthHandler = createHandler(
  registerUserWithOAuth,
  'auth',
  'registerWithOAuth'
);
export const loginHandler = createHandler(loginUser, 'auth', 'login');
export const refreshTokenHandler = createHandler(refreshUserToken, 'auth', 'refreshToken');
export const logoutHandler = createSimpleHandler(logoutUser, 'auth', 'logout');
export const forgotPasswordHandler = createHandler(forgotPassword, 'auth', 'forgotPassword');
export const resetPasswordHandler = createHandler(resetPassword, 'auth', 'resetPassword');
export const sendVerificationEmailHandler = createHandler(
  sendVerificationEmail,
  'auth',
  'sendVerificationEmail'
);
export const verifyEmailWithPinHandler = createHandler(
  verifyEmailWithPin,
  'auth',
  'verifyEmailWithPin'
);
export const resendVerificationEmailHandler = createHandler(
  resendVerificationEmail,
  'auth',
  'resendVerificationEmail'
);

// Token Validation
export const validateTokenHandler = createHandler(validateToken, 'auth', 'validateToken');

// OAuth Login
export const oAuthLoginHandler = createHandler(oAuthLogin, 'auth', 'oAuthLogin');

// Email Verification (Token-based)
export const verifyEmailHandler = createHandler(verifyEmail, 'auth', 'verifyEmail');

// Permissions & Roles
export const getUserPermissionsHandler = createSimpleHandler(getUserPermissions, 'auth', 'getUserPermissions');
export const checkPermissionHandler = createHandler(checkPermission, 'auth', 'checkPermission');
export const checkResourcePermissionHandler = createHandler(checkResourcePermission, 'auth', 'checkResourcePermission');
export const getUserRolesHandler = createSimpleHandler(getUserRoles, 'auth', 'getUserRoles');
export const batchCheckPermissionsHandler = createHandler(batchCheckPermissions, 'auth', 'batchCheckPermissions');

// User Management (Admin)
export const getUserHandler = createSimpleHandler(getUser, 'auth', 'getUser');
export const updateUserHandler = createHandler(updateUser, 'auth', 'updateUser');
export const deleteUserHandler = createHandler(deleteUser, 'auth', 'deleteUser');

// Utility
export const registerHandler = createHandler(register, 'auth', 'register');
export const healthCheckHandler = createSimpleHandler(healthCheck, 'auth', 'health');

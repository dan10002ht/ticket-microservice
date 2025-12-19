/**
 * @swagger
 * components:
 *   schemas:
 *     UserRegistration:
 *       type: object
 *       required:
 *         - email
 *         - password
 *         - first_name
 *         - last_name
 *       properties:
 *         email:
 *           type: string
 *           format: email
 *           description: User email address
 *         password:
 *           type: string
 *           minLength: 8
 *           description: User password (min 8 characters)
 *         first_name:
 *           type: string
 *           description: User first name
 *         last_name:
 *           type: string
 *           description: User last name
 *     OAuthRegistration:
 *       type: object
 *       required:
 *         - provider
 *         - token
 *       properties:
 *         provider:
 *           type: string
 *           enum: [google, facebook, github]
 *           description: OAuth provider
 *         token:
 *           type: string
 *           description: OAuth token from provider
 *         access_token:
 *           type: string
 *           description: OAuth access token
 *         refresh_token:
 *           type: string
 *           description: OAuth refresh token
 *         expires_at:
 *           type: integer
 *           description: Token expiration timestamp
 *     UserLogin:
 *       type: object
 *       required:
 *         - email
 *         - password
 *       properties:
 *         email:
 *           type: string
 *           format: email
 *           description: User email address
 *         password:
 *           type: string
 *           description: User password
 *     RefreshToken:
 *       type: object
 *       required:
 *         - refresh_token
 *       properties:
 *         refresh_token:
 *           type: string
 *           description: Refresh token for getting new access token
 *     AuthResponse:
 *       type: object
 *       properties:
 *         access_token:
 *           type: string
 *           description: JWT access token
 *         refresh_token:
 *           type: string
 *           description: JWT refresh token
 *         expires_in:
 *           type: integer
 *           description: Token expiration time in seconds
 *         auth_type:
 *           type: string
 *           enum: [email, oauth]
 *           description: Authentication type
 *         is_new_user:
 *           type: boolean
 *           description: Whether this is a new user (for OAuth)
 */

/**
 * @swagger
 * /auth/register/email:
 *   post:
 *     summary: Register a new user with email and password
 *     description: Create a new user account using email and password
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UserRegistration'
 *     responses:
 *       201:
 *         description: User registered successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AuthResponse'
 *       400:
 *         description: Validation error
 *       409:
 *         description: User already exists
 */

/**
 * @swagger
 * /auth/register/oauth:
 *   post:
 *     summary: Register a new user with OAuth
 *     description: Create a new user account using OAuth provider (Google, Facebook, GitHub)
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/OAuthRegistration'
 *     responses:
 *       201:
 *         description: User registered successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AuthResponse'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Invalid OAuth token
 */

/**
 * @swagger
 * /auth/login:
 *   post:
 *     summary: User login
 *     description: Authenticate user and return JWT tokens
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UserLogin'
 *     responses:
 *       200:
 *         description: Login successful
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AuthResponse'
 *       401:
 *         description: Invalid credentials
 *       400:
 *         description: Validation error
 */

/**
 * @swagger
 * /auth/refresh:
 *   post:
 *     summary: Refresh access token
 *     description: Get new access token using refresh token
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/RefreshToken'
 *     responses:
 *       200:
 *         description: Token refreshed successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AuthResponse'
 *       401:
 *         description: Invalid refresh token
 */

/**
 * @swagger
 * /auth/logout:
 *   post:
 *     summary: User logout
 *     description: Invalidate refresh token
 *     tags: [Auth]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Logout successful
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /auth/hello:
 *   get:
 *     summary: Hello endpoint
 *     description: Simple hello world endpoint
 *     tags: [Auth]
 *     responses:
 *       200:
 *         description: Hello message
 *         content:
 *           text/plain:
 *             schema:
 *               type: string
 *               example: "Hello World"
 */

/**
 * @swagger
 * components:
 *   schemas:
 *     SendVerificationEmail:
 *       type: object
 *       required:
 *         - email
 *       properties:
 *         email:
 *           type: string
 *           format: email
 *           description: User email address to send verification to
 *     VerifyUserWithPin:
 *       type: object
 *       required:
 *         - user_id
 *         - pin_code
 *       properties:
 *         user_id:
 *           type: string
 *           format: uuid
 *           description: User ID to verify
 *         pin_code:
 *           type: string
 *           minLength: 6
 *           maxLength: 6
 *           pattern: '^[0-9]{6}$'
 *           description: 6-digit PIN code sent to email
 *     ResendVerificationEmail:
 *       type: object
 *       required:
 *         - email
 *       properties:
 *         email:
 *           type: string
 *           format: email
 *           description: User email address to resend verification to
 *     EmailVerificationResponse:
 *       type: object
 *       properties:
 *         success:
 *           type: boolean
 *           description: Whether the operation was successful
 *         message:
 *           type: string
 *           description: Response message
 *         user_id:
 *           type: string
 *           description: User ID (for send/resend operations)
 *         user_email:
 *           type: string
 *           description: User email address
 *         pin_code:
 *           type: string
 *           description: Generated PIN code (only in development mode)
 *         expires_at:
 *           type: string
 *           format: date-time
 *           description: PIN code expiration time
 */

/**
 * @swagger
 * /auth/send-verification-email:
 *   post:
 *     summary: Send verification email
 *     description: Send a verification email with PIN code to the user's email address
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/SendVerificationEmail'
 *     responses:
 *       200:
 *         description: Verification email sent successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/EmailVerificationResponse'
 *       400:
 *         description: Validation error
 *       404:
 *         description: User not found
 *       409:
 *         description: Email already verified
 *       500:
 *         description: Failed to send verification email
 */

/**
 * @swagger
 * /auth/verify-user:
 *   post:
 *     summary: Verify user email with PIN code
 *     description: Verify user's email address using the PIN code sent to their email
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/VerifyUserWithPin'
 *     responses:
 *       200:
 *         description: Email verified successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                   example: true
 *                 message:
 *                   type: string
 *                   example: "Email verified successfully"
 *                 user:
 *                   $ref: '#/components/schemas/User'
 *       400:
 *         description: Validation error or invalid PIN code
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "User ID and PIN code are required"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "MISSING_REQUIRED_FIELD"
 *       404:
 *         description: User not found
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "User not found"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "USER_NOT_FOUND"
 *       409:
 *         description: Email already verified
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Email is already verified"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "EMAIL_ALREADY_VERIFIED"
 *       410:
 *         description: PIN code expired
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "PIN code has expired"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "PIN_CODE_EXPIRED"
 *       500:
 *         description: Internal server error
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Email verification failed"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "EMAIL_VERIFICATION_FAILED"
 */

/**
 * @swagger
 * /auth/resend-verification-email:
 *   post:
 *     summary: Resend verification email
 *     description: Resend a verification email with new PIN code to the user's email address
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/ResendVerificationEmail'
 *     responses:
 *       200:
 *         description: Verification email resent successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/EmailVerificationResponse'
 *       400:
 *         description: Validation error
 *       404:
 *         description: User not found
 *       409:
 *         description: Email already verified
 *       500:
 *         description: Failed to send verification email
 */

/**
 * @swagger
 * components:
 *   schemas:
 *     ForgotPassword:
 *       type: object
 *       required:
 *         - email
 *       properties:
 *         email:
 *           type: string
 *           format: email
 *           description: User email address to send password reset link
 *     ResetPassword:
 *       type: object
 *       required:
 *         - token
 *         - new_password
 *         - confirm_password
 *       properties:
 *         token:
 *           type: string
 *           description: Password reset token received via email
 *         new_password:
 *           type: string
 *           minLength: 8
 *           description: New password (min 8 characters)
 *         confirm_password:
 *           type: string
 *           description: Confirm new password
 *     PasswordResetResponse:
 *       type: object
 *       properties:
 *         success:
 *           type: boolean
 *           description: Whether the operation was successful
 *         message:
 *           type: string
 *           description: Response message
 */

/**
 * @swagger
 * /auth/forgot-password:
 *   post:
 *     summary: Request password reset
 *     description: Send a password reset link to the user's email address
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/ForgotPassword'
 *     responses:
 *       200:
 *         description: Password reset email sent successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/PasswordResetResponse'
 *       400:
 *         description: Validation error
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Email is required"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "MISSING_REQUIRED_FIELD"
 *       404:
 *         description: User not found
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "User not found"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "USER_NOT_FOUND"
 *       429:
 *         description: Too many password reset attempts
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Too many password reset attempts, please try again later"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "TOO_MANY_ATTEMPTS"
 *       500:
 *         description: Failed to send password reset email
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Failed to send password reset email"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "PASSWORD_RESET_FAILED"
 */

/**
 * @swagger
 * /auth/reset-password:
 *   post:
 *     summary: Reset password with token
 *     description: Reset user password using the token received via email
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/ResetPassword'
 *     responses:
 *       200:
 *         description: Password reset successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/PasswordResetResponse'
 *       400:
 *         description: Validation error
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Token and new password are required"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "MISSING_REQUIRED_FIELD"
 *       401:
 *         description: Invalid or expired reset token
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Invalid or expired reset token"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "INVALID_TOKEN"
 *       404:
 *         description: User not found
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "User not found"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "USER_NOT_FOUND"
 *       500:
 *         description: Password reset failed
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 error:
 *                   type: string
 *                   example: "Password reset failed"
 *                 correlationId:
 *                   type: string
 *                   format: uuid
 *                 code:
 *                   type: string
 *                   example: "PASSWORD_RESET_FAILED"
 */

/**
 * @swagger
 * components:
 *   schemas:
 *     ValidateTokenRequest:
 *       type: object
 *       properties:
 *         token:
 *           type: string
 *           description: JWT token to validate (optional if using Authorization header)
 *     ValidateTokenResponse:
 *       type: object
 *       properties:
 *         valid:
 *           type: boolean
 *           description: Whether the token is valid
 *         user:
 *           $ref: '#/components/schemas/User'
 *     OAuthLoginRequest:
 *       type: object
 *       required:
 *         - provider
 *         - access_token
 *       properties:
 *         provider:
 *           type: string
 *           enum: [google, facebook, github]
 *           description: OAuth provider
 *         access_token:
 *           type: string
 *           description: OAuth access token
 *         provider_user_id:
 *           type: string
 *           description: User ID from OAuth provider
 *         email:
 *           type: string
 *           format: email
 *         first_name:
 *           type: string
 *         last_name:
 *           type: string
 *         picture:
 *           type: string
 *           description: Profile picture URL
 *         refresh_token:
 *           type: string
 *         expires_at:
 *           type: integer
 *     PermissionData:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         name:
 *           type: string
 *         description:
 *           type: string
 *         resource:
 *           type: string
 *         action:
 *           type: string
 *     RoleData:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         name:
 *           type: string
 *         description:
 *           type: string
 *         permissions:
 *           type: array
 *           items:
 *             $ref: '#/components/schemas/PermissionData'
 *     CheckPermissionRequest:
 *       type: object
 *       required:
 *         - permission_name
 *       properties:
 *         permission_name:
 *           type: string
 *           description: Permission name to check
 *     CheckResourcePermissionRequest:
 *       type: object
 *       required:
 *         - resource
 *         - action
 *       properties:
 *         resource:
 *           type: string
 *           description: Resource name
 *         action:
 *           type: string
 *           description: Action to check
 *         context:
 *           type: object
 *           description: Additional context for permission check
 *     BatchCheckPermissionsRequest:
 *       type: object
 *       required:
 *         - permission_names
 *       properties:
 *         permission_names:
 *           type: array
 *           items:
 *             type: string
 *           description: List of permission names to check
 *     User:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         email:
 *           type: string
 *         first_name:
 *           type: string
 *         last_name:
 *           type: string
 *         phone:
 *           type: string
 *         role:
 *           type: string
 *         is_active:
 *           type: boolean
 *         is_verified:
 *           type: boolean
 */

/**
 * @swagger
 * /api/auth/validate:
 *   post:
 *     summary: Validate JWT token
 *     description: Validate a JWT token and return user info if valid
 *     tags: [Auth - Token]
 *     requestBody:
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/ValidateTokenRequest'
 *     responses:
 *       200:
 *         description: Token validation result
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/ValidateTokenResponse'
 *       401:
 *         description: Invalid token
 */

/**
 * @swagger
 * /api/auth/oauth/login:
 *   post:
 *     summary: OAuth login
 *     description: Login using OAuth provider
 *     tags: [Auth - OAuth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/OAuthLoginRequest'
 *     responses:
 *       200:
 *         description: Login successful
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AuthResponse'
 *       401:
 *         description: Invalid OAuth credentials
 */

/**
 * @swagger
 * /api/auth/verify-email-token:
 *   post:
 *     summary: Verify email with token
 *     description: Verify user email using token from email link
 *     tags: [Auth - Email Verification]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - token
 *             properties:
 *               token:
 *                 type: string
 *                 description: Email verification token
 *     responses:
 *       200:
 *         description: Email verified successfully
 *       400:
 *         description: Invalid or expired token
 */

/**
 * @swagger
 * /api/auth/register:
 *   post:
 *     summary: Basic registration
 *     description: Register a new user with basic information
 *     tags: [Auth]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - email
 *               - password
 *             properties:
 *               email:
 *                 type: string
 *                 format: email
 *               password:
 *                 type: string
 *                 minLength: 8
 *               first_name:
 *                 type: string
 *               last_name:
 *                 type: string
 *               phone:
 *                 type: string
 *               role:
 *                 type: string
 *                 enum: [individual, organization, admin]
 *     responses:
 *       201:
 *         description: User registered successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AuthResponse'
 *       400:
 *         description: Validation error
 *       409:
 *         description: User already exists
 */

/**
 * @swagger
 * /api/auth/permissions:
 *   get:
 *     summary: Get user permissions
 *     description: Get all permissions for the authenticated user
 *     tags: [Auth - Permissions]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: User permissions
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 permissions:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/PermissionData'
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/auth/permissions/check:
 *   post:
 *     summary: Check single permission
 *     description: Check if user has a specific permission
 *     tags: [Auth - Permissions]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/CheckPermissionRequest'
 *     responses:
 *       200:
 *         description: Permission check result
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 allowed:
 *                   type: boolean
 *                 reason:
 *                   type: string
 *                 roles:
 *                   type: array
 *                   items:
 *                     type: string
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/auth/permissions/resource:
 *   post:
 *     summary: Check resource permission
 *     description: Check if user has permission for a specific resource and action
 *     tags: [Auth - Permissions]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/CheckResourcePermissionRequest'
 *     responses:
 *       200:
 *         description: Resource permission check result
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 allowed:
 *                   type: boolean
 *                 reason:
 *                   type: string
 *                 permissions:
 *                   type: array
 *                   items:
 *                     type: string
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/auth/permissions/batch:
 *   post:
 *     summary: Batch check permissions
 *     description: Check multiple permissions at once
 *     tags: [Auth - Permissions]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/BatchCheckPermissionsRequest'
 *     responses:
 *       200:
 *         description: Batch permission check results
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 results:
 *                   type: object
 *                   additionalProperties:
 *                     type: boolean
 *                 roles:
 *                   type: array
 *                   items:
 *                     type: string
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/auth/roles:
 *   get:
 *     summary: Get user roles
 *     description: Get all roles for the authenticated user
 *     tags: [Auth - Permissions]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: User roles
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 roles:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/RoleData'
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/auth/users/{userId}:
 *   get:
 *     summary: Get user by ID (Admin)
 *     description: Get user details by ID (Admin only)
 *     tags: [Auth - User Management]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: userId
 *         required: true
 *         schema:
 *           type: string
 *         description: User ID
 *     responses:
 *       200:
 *         description: User details
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 user:
 *                   $ref: '#/components/schemas/User'
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 *       404:
 *         description: User not found
 *   put:
 *     summary: Update user (Admin)
 *     description: Update user details (Admin only)
 *     tags: [Auth - User Management]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: userId
 *         required: true
 *         schema:
 *           type: string
 *         description: User ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               first_name:
 *                 type: string
 *               last_name:
 *                 type: string
 *               phone:
 *                 type: string
 *               address:
 *                 type: string
 *               city:
 *                 type: string
 *               state:
 *                 type: string
 *               country:
 *                 type: string
 *               postal_code:
 *                 type: string
 *     responses:
 *       200:
 *         description: User updated successfully
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 *       404:
 *         description: User not found
 *   delete:
 *     summary: Delete user (Admin)
 *     description: Delete a user (Admin only)
 *     tags: [Auth - User Management]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: userId
 *         required: true
 *         schema:
 *           type: string
 *         description: User ID
 *     responses:
 *       200:
 *         description: User deleted successfully
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /api/auth/health:
 *   get:
 *     summary: Health check
 *     description: Check auth service health
 *     tags: [Auth - Health]
 *     responses:
 *       200:
 *         description: Service is healthy
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 status:
 *                   type: string
 *                   example: healthy
 *                 timestamp:
 *                   type: string
 *                   format: date-time
 */

export default {};

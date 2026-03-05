// ── Request types ──

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  role?: "individual" | "organization";
  organization?: { name: string };
}

export interface OAuthLoginRequest {
  provider: "google" | "facebook" | "github";
  access_token: string;
  provider_user_id?: string;
  email?: string;
  first_name?: string;
  last_name?: string;
  picture?: string;
  refresh_token?: string;
  expires_at?: number;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  password: string;
}

export interface SendVerificationEmailRequest {
  email: string;
}

export interface VerifyUserRequest {
  user_id: string;
  pin_code: string;
}

export interface CheckPermissionRequest {
  permission_name: string;
}

export interface CheckResourcePermissionRequest {
  resource: string;
  action: string;
  context?: Record<string, unknown>;
}

export interface BatchCheckPermissionsRequest {
  permission_names: string[];
}

// ── Response types ──
// Note: Gateway toCamelCase transforms all response keys to camelCase

export interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  authType: "email" | "oauth";
  isNewUser?: boolean;
}

export interface AuthUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  role?: string;
  isActive?: boolean;
  isVerified?: boolean;
}

export interface ValidateTokenResponse {
  valid: boolean;
  user?: AuthUser;
}

export interface EmailVerificationResponse {
  success: boolean;
  message: string;
  userId?: string;
  userEmail?: string;
  expiresAt?: string;
}

export interface PasswordResetResponse {
  success: boolean;
  message: string;
}

export interface PermissionData {
  id: string;
  name: string;
  description?: string;
  resource?: string;
  action?: string;
}

export interface RoleData {
  id: string;
  name: string;
  description?: string;
  permissions: PermissionData[];
}

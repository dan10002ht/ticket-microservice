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
  role?: "user" | "organization";
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
  new_password: string;
  confirm_password: string;
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

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  auth_type: "email" | "oauth";
  is_new_user?: boolean;
}

export interface AuthUser {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  role?: string;
  is_active?: boolean;
  is_verified?: boolean;
}

export interface ValidateTokenResponse {
  valid: boolean;
  user?: AuthUser;
}

export interface EmailVerificationResponse {
  success: boolean;
  message: string;
  user_id?: string;
  user_email?: string;
  expires_at?: string;
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

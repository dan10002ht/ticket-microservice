import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import { useAuthStore } from "@/stores/auth-store";
import { setTokens, clearTokens, getAccessToken } from "@/lib/auth/token";
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  AuthUser,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  VerifyUserRequest,
  ValidateTokenResponse,
  EmailVerificationResponse,
  PasswordResetResponse,
} from "../types/auth";
import type { ApiError } from "../types/common";

/**
 * Validate existing token and fetch current user.
 * Runs on app initialization to restore auth state from cookies.
 * Side effects (Zustand store sync) are handled by AuthInitializer.
 */
export function useMe() {
  return useQuery<AuthUser, ApiError>({
    queryKey: queryKeys.auth.user(),
    queryFn: async () => {
      const token = getAccessToken();
      if (!token) throw { error: { code: "UNAUTHENTICATED", message: "No access token" }, statusCode: 401 } as ApiError;

      const { data } = await apiClient.post<ValidateTokenResponse>(
        API_ENDPOINTS.auth.validate,
        { token },
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        { _skipErrorToast: true } as any
      );

      if (!data.valid || !data.user) {
        throw { error: { code: "INVALID_TOKEN", message: "Invalid token" }, statusCode: 401 } as ApiError;
      }
      return data.user;
    },
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: typeof window !== "undefined",
  });
}

/** Login with email/password */
export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation<AuthResponse, ApiError, LoginRequest>({
    mutationFn: async (credentials) => {
      const { data } = await apiClient.post<AuthResponse>(
        API_ENDPOINTS.auth.login,
        credentials,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        { _skipErrorToast: true } as any
      );
      return data;
    },
    onSuccess: (data) => {
      setTokens(data.access_token, data.refresh_token, data.expires_in);
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.user() });
    },
  });
}

/** Register with email */
export function useRegister() {
  const queryClient = useQueryClient();

  return useMutation<AuthResponse, ApiError, RegisterRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post<AuthResponse>(
        API_ENDPOINTS.auth.registerEmail,
        input,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        { _skipErrorToast: true } as any
      );
      return data;
    },
    onSuccess: (data) => {
      setTokens(data.access_token, data.refresh_token, data.expires_in);
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.user() });
    },
  });
}

/** Logout — clears tokens + store + cache, redirects to login */
export function useLogout() {
  const queryClient = useQueryClient();
  const { clearUser } = useAuthStore();

  return useMutation<void, ApiError, void>({
    mutationFn: async () => {
      await apiClient.post(API_ENDPOINTS.auth.logout);
    },
    onSettled: () => {
      clearTokens();
      clearUser();
      queryClient.clear();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    },
  });
}

/** Send forgot password email */
export function useForgotPassword() {
  return useMutation<PasswordResetResponse, ApiError, ForgotPasswordRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post<PasswordResetResponse>(
        API_ENDPOINTS.auth.forgotPassword,
        input
      );
      return data;
    },
  });
}

/** Reset password with token */
export function useResetPassword() {
  return useMutation<PasswordResetResponse, ApiError, ResetPasswordRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post<PasswordResetResponse>(
        API_ENDPOINTS.auth.resetPassword,
        input
      );
      return data;
    },
  });
}

/** Send verification email */
export function useSendVerificationEmail() {
  return useMutation<EmailVerificationResponse, ApiError, { email: string }>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post<EmailVerificationResponse>(
        API_ENDPOINTS.auth.sendVerificationEmail,
        input
      );
      return data;
    },
  });
}

/** Verify user with PIN code */
export function useVerifyUser() {
  return useMutation<{ success: boolean; message: string; user?: AuthUser }, ApiError, VerifyUserRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.auth.verifyUser,
        input
      );
      return data;
    },
  });
}

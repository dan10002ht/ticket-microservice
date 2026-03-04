import type { AxiosInstance, InternalAxiosRequestConfig, AxiosError } from "axios";
import { getAccessToken, getRefreshToken, setTokens, clearTokens } from "@/lib/auth/token";
import { API_ENDPOINTS } from "../endpoints";
import type { AuthResponse } from "../types/auth";

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null): void {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error);
    else if (token) resolve(token);
  });
  failedQueue = [];
}

/** Attach JWT access token to every request */
export function attachAuthToken(
  config: InternalAxiosRequestConfig
): InternalAxiosRequestConfig {
  const token = getAccessToken();
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}

/**
 * Create a 401 response interceptor that atomically refreshes the token.
 *
 * Race condition prevention: when multiple requests fail with 401 simultaneously,
 * only the first triggers a refresh. All others queue and replay once the new
 * token arrives.
 */
export function createRefreshInterceptor(apiClient: AxiosInstance) {
  return async (error: AxiosError) => {
    const originalRequest = error.config;
    if (!originalRequest || error.response?.status !== 401) {
      return Promise.reject(error);
    }

    // Don't refresh for auth endpoints (avoid infinite loop)
    if (originalRequest.url?.includes("/auth/")) {
      return Promise.reject(error);
    }

    // Guest user (never authenticated) — just reject, don't redirect
    if (!getAccessToken() && !getRefreshToken()) {
      return Promise.reject(error);
    }

    // If already refreshing, queue this request
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({
          resolve: (token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(apiClient(originalRequest));
          },
          reject,
        });
      });
    }

    isRefreshing = true;
    const refreshToken = getRefreshToken();

    if (!refreshToken) {
      clearTokens();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
      return Promise.reject(error);
    }

    try {
      const { data } = await apiClient.post<AuthResponse>(
        API_ENDPOINTS.auth.refresh,
        { refresh_token: refreshToken }
      );

      setTokens(data.access_token, data.refresh_token, data.expires_in);
      processQueue(null, data.access_token);

      originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
      return apiClient(originalRequest);
    } catch (refreshError) {
      processQueue(refreshError, null);
      clearTokens();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  };
}

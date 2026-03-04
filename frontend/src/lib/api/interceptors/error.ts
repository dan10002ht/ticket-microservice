import type { AxiosError, AxiosResponse } from "axios";
import { showToast } from "@/lib/toast";
import type { ApiError } from "../types/common";

const ERROR_MESSAGES: Record<number, string> = {
  400: "Invalid request. Please check and try again.",
  403: "You do not have permission to perform this action.",
  404: "The requested resource was not found.",
  409: "Data already exists or conflicts with existing data.",
  422: "Invalid data. Please check and try again.",
  429: "Too many requests. Please try again later.",
  500: "Something went wrong. Please try again.",
  502: "Service temporarily unavailable.",
  503: "Service temporarily unavailable.",
};

/**
 * Normalize API errors and show auto-toast.
 *
 * To suppress auto-toast for a specific request, set `_skipErrorToast: true`
 * on the Axios request config:
 *
 * ```ts
 * apiClient.post(url, data, { _skipErrorToast: true } as any)
 * ```
 */
export function handleResponseError(
  error: AxiosError<ApiError>
): Promise<never> {
  const status = error.response?.status;
  const apiError = error.response?.data;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const skipToast = (error.config as any)?._skipErrorToast;

  // Normalize error for consumers
  const normalizedError: ApiError = {
    error: apiError?.error || error.message || "Network error",
    correlationId: apiError?.correlationId,
    code: apiError?.code,
    statusCode: status,
    details: apiError?.details,
  };

  // Auto-toast for non-401 errors (401 handled by refresh interceptor)
  if (!skipToast && status && status !== 401) {
    // Use showToast.apiError for validation errors with details
    if (apiError?.details && Array.isArray(apiError.details)) {
      showToast.apiError(normalizedError);
    } else {
      const message =
        apiError?.error ||
        ERROR_MESSAGES[status] ||
        "An unexpected error occurred.";
      showToast.error(message);
    }
  }

  return Promise.reject(normalizedError);
}

export function handleResponseSuccess(response: AxiosResponse): AxiosResponse {
  return response;
}

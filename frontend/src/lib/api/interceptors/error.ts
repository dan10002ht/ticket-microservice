import type { AxiosError, AxiosResponse } from "axios";
import { showToast } from "@/lib/toast";
import type { ApiError } from "../types/common";

/** Gateway error response shape: { error: { code, message, details? }, meta } */
interface GatewayErrorResponse {
  error: {
    code: string;
    message: string;
    details?: unknown;
  };
  meta?: {
    correlationId?: string;
    timestamp?: string;
  };
}

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
  error: AxiosError<GatewayErrorResponse>
): Promise<never> {
  const status = error.response?.status;
  const responseData = error.response?.data;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const skipToast = (error.config as any)?._skipErrorToast;

  // Normalize to ApiError for consumers
  const normalizedError: ApiError = {
    error: responseData?.error || {
      code: "NETWORK_ERROR",
      message: error.message || "Network error",
    },
    meta: responseData?.meta,
    statusCode: status,
  };

  // Auto-toast for non-401 errors (401 handled by refresh interceptor)
  if (!skipToast && status && status !== 401) {
    const errorObj = normalizedError.error;

    // Use showToast.apiError for validation errors with details
    if (errorObj.details && Array.isArray(errorObj.details)) {
      showToast.apiError(normalizedError);
    } else {
      const message =
        errorObj.message ||
        (status ? ERROR_MESSAGES[status] : undefined) ||
        "An unexpected error occurred.";
      showToast.error(message);
    }
  }

  return Promise.reject(normalizedError);
}

export function handleResponseSuccess(response: AxiosResponse): AxiosResponse {
  // Unwrap gateway envelope: { data: {...}, meta: {...} } → just the inner data
  if (response.data && typeof response.data === "object" && "data" in response.data && "meta" in response.data) {
    response.data = response.data.data;
  }
  return response;
}

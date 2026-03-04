import { toast, type ExternalToast } from "sonner";
import type { ApiError } from "./api/types/common";

interface ValidationDetail {
  path: string;
  msg: string;
  type?: string;
  value?: unknown;
  location?: string;
}

// ── Core showToast ──

function showToast(
  type: "success" | "error" | "info" | "warning",
  message: string,
  options?: ExternalToast
) {
  toast[type](message, options);
}

// ── Convenience methods ──

showToast.success = (message: string, options?: ExternalToast) =>
  toast.success(message, options);

showToast.error = (message: string, options?: ExternalToast) =>
  toast.error(message, options);

showToast.info = (message: string, options?: ExternalToast) =>
  toast.info(message, options);

showToast.warning = (message: string, options?: ExternalToast) =>
  toast.warning(message, options);

/**
 * Smart API error toast.
 *
 * - If `details[]` exists (validation errors): shows first detail's `msg`
 *   with a description summarizing additional errors
 * - Otherwise: shows the generic `error` message
 *
 * Returns parsed field errors for form integration:
 * ```ts
 * const fieldErrors = showToast.apiError(err);
 * fieldErrors.forEach(({ path, msg }) => setError(path, { message: msg }));
 * ```
 */
showToast.apiError = (
  error: ApiError,
  fallbackMessage = "Something went wrong. Please try again."
): ValidationDetail[] => {
  const details = parseValidationDetails(error.details);

  if (details.length > 0) {
    const firstMsg = details[0].msg;
    const description =
      details.length > 1
        ? `and ${details.length - 1} more error${details.length - 1 > 1 ? "s" : ""}`
        : undefined;

    toast.error(firstMsg, { description });
  } else {
    toast.error(error.error || fallbackMessage);
  }

  return details;
};

// ── Helpers ──

function parseValidationDetails(details: unknown): ValidationDetail[] {
  if (!Array.isArray(details)) return [];
  return details.filter(
    (d): d is ValidationDetail =>
      typeof d === "object" &&
      d !== null &&
      typeof d.path === "string" &&
      typeof d.msg === "string"
  );
}

export { showToast };

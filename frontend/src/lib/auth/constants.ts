export const AUTH_COOKIE_NAMES = {
  ACCESS_TOKEN: "tb_access_token",
  REFRESH_TOKEN: "tb_refresh_token",
} as const;

export const AUTH_CONFIG = {
  /** Seconds before expiry to trigger preemptive refresh */
  ACCESS_TOKEN_EXPIRY_BUFFER: 60,
  /** Access token max age in seconds (15 minutes) */
  ACCESS_TOKEN_MAX_AGE: 15 * 60,
  /** Refresh token max age in seconds (7 days) */
  REFRESH_TOKEN_MAX_AGE: 7 * 24 * 60 * 60,
} as const;

/** Routes accessible without authentication */
export const PUBLIC_ROUTES = [
  "/",
  "/login",
  "/register",
  "/events",
  "/forgot-password",
  "/reset-password",
] as const;

/** Auth routes — redirect to home if already authenticated */
export const AUTH_ROUTES = ["/login", "/register"] as const;

/** Route prefixes that require authentication */
export const PROTECTED_ROUTE_PREFIXES = [
  "/my-bookings",
  "/profile",
  "/org",
  "/admin",
] as const;

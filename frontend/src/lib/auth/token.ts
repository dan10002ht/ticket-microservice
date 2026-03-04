import { AUTH_COOKIE_NAMES, AUTH_CONFIG } from "./constants";

// ── Client-side cookie helpers ──

function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(
    new RegExp(`(?:^|; )${name.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}=([^;]*)`)
  );
  return match ? decodeURIComponent(match[1]) : null;
}

function setCookie(name: string, value: string, maxAge: number): void {
  const isSecure =
    typeof window !== "undefined" && window.location.protocol === "https:";
  const parts = [
    `${name}=${encodeURIComponent(value)}`,
    `path=/`,
    `max-age=${maxAge}`,
    `SameSite=Lax`,
  ];
  if (isSecure) parts.push("Secure");
  document.cookie = parts.join("; ");
}

function deleteCookie(name: string): void {
  document.cookie = `${name}=; path=/; max-age=0`;
}

// ── Public API ──

export function getAccessToken(): string | null {
  return getCookie(AUTH_COOKIE_NAMES.ACCESS_TOKEN);
}

export function getRefreshToken(): string | null {
  return getCookie(AUTH_COOKIE_NAMES.REFRESH_TOKEN);
}

export function setTokens(
  accessToken: string,
  refreshToken: string,
  expiresIn: number
): void {
  setCookie(
    AUTH_COOKIE_NAMES.ACCESS_TOKEN,
    accessToken,
    expiresIn || AUTH_CONFIG.ACCESS_TOKEN_MAX_AGE
  );
  setCookie(
    AUTH_COOKIE_NAMES.REFRESH_TOKEN,
    refreshToken,
    AUTH_CONFIG.REFRESH_TOKEN_MAX_AGE
  );
}

export function clearTokens(): void {
  deleteCookie(AUTH_COOKIE_NAMES.ACCESS_TOKEN);
  deleteCookie(AUTH_COOKIE_NAMES.REFRESH_TOKEN);
}

// ── Server-side / Middleware helpers ──

export function getAccessTokenFromCookies(cookieString: string): string | null {
  return parseCookieValue(cookieString, AUTH_COOKIE_NAMES.ACCESS_TOKEN);
}

export function getRefreshTokenFromCookies(
  cookieString: string
): string | null {
  return parseCookieValue(cookieString, AUTH_COOKIE_NAMES.REFRESH_TOKEN);
}

function parseCookieValue(
  cookieString: string,
  name: string
): string | null {
  const match = cookieString.match(
    new RegExp(`(?:^|; )${name}=([^;]*)`)
  );
  return match ? decodeURIComponent(match[1]) : null;
}

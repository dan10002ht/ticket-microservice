import { NextRequest, NextResponse } from "next/server";

const ACCESS_TOKEN_COOKIE = "tb_access_token";

const AUTH_ROUTES = ["/login", "/register"];
const PROTECTED_ROUTE_PREFIXES = ["/my-bookings", "/profile", "/org", "/admin"];
const PROTECTED_ROUTE_PATTERNS = [/^\/events\/[^/]+\/book/];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get(ACCESS_TOKEN_COOKIE)?.value;

  const isProtectedRoute =
    PROTECTED_ROUTE_PREFIXES.some((prefix) => pathname.startsWith(prefix)) ||
    PROTECTED_ROUTE_PATTERNS.some((re) => re.test(pathname));

  const isAuthRoute = AUTH_ROUTES.some((route) => pathname === route);

  // Redirect unauthenticated users away from protected routes
  if (isProtectedRoute && !accessToken) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("callbackUrl", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // Redirect authenticated users away from auth routes
  if (isAuthRoute && accessToken) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all paths except:
     * - _next/static, _next/image (Next.js internals)
     * - favicon.ico, sitemap.xml, robots.txt (metadata)
     * - api routes, public assets
     */
    "/((?!_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt|api|public).*)",
  ],
};

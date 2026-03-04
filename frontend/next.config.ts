import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Proxy API requests through Next.js to avoid CORS in development
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:53000/api/v1"}/:path*`,
      },
    ];
  },
};

export default nextConfig;

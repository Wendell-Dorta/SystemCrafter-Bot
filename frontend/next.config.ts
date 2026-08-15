import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${process.env.INTERNAL_BACKEND_URL || "http://backend:8080"}/api/:path*`,
      },
      {
        source: "/health",
        destination: `${process.env.INTERNAL_BACKEND_URL || "http://backend:8080"}/health`,
      },
    ];
  },
};

export default nextConfig;

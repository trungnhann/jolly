import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    const baseURL = process.env.API_BASE_URL ?? "http://localhost:8080";
    return [
      {
        source: "/api/:path*",
        destination: `${baseURL}/:path*`,
      },
    ];
  },
};

export default nextConfig;

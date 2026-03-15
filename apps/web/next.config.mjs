/** @type {import('next').NextConfig} */
const apiOrigin = (
  process.env.NEXT_PUBLIC_API_URL || "https://hackersquare-api.up.railway.app"
).replace(/\/+$/, "");

const nextConfig = {
  async rewrites() {
    return {
      fallback: [
        {
          source: '/api/:path*',
          destination: `${apiOrigin}/:path*`,
        },
      ],
    }
  },
}

export default nextConfig;
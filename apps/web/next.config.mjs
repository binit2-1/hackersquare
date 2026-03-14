/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return {
      fallback: [
        {
          source: '/api/:path*',
          destination: 'https://hackersquare-api.up.railway.app/:path*',
        },
      ],
    }
  },
}

export default nextConfig;
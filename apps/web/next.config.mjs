/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'https://hackersquare-api.up.railway.app';
    
    return [
      {
        source: '/api/:path*',
        destination: `${apiUrl}/:path*`, 
      },
    ];
  },
};

export default nextConfig;
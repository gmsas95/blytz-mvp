/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    domains: ['localhost', 'images.unsplash.com'],
  },
  env: {
    MODE: process.env.MODE || 'mock',
    REMOTE_API_BASE: process.env.REMOTE_API_BASE || 'http://localhost:8081',
    NEXT_PUBLIC_API_BASE: process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:8081',
  },
  // experimental: {
  //   appDir: true,
  // },
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/health',
        destination: '/api/health',
      },
    ]
  },
}

module.exports = nextConfig
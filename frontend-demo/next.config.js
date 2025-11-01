/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'https://api.blytz.app',
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL || 'wss://api.blytz.app',
    NEXT_PUBLIC_LIVEKIT_URL: process.env.NEXT_PUBLIC_LIVEKIT_URL || 'wss://blytz-live-u5u72ozx.livekit.cloud',
  },
  async rewrites() {
    return [
      {
        source: '/api/livekit/token',
        destination: 'https://api.blytz.app/api/public/livekit/token',
      },
    ];
  },
}

module.exports = nextConfig
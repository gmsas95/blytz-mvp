import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
  // Add security headers
  const response = NextResponse.next();

  // Remove sensitive headers
  response.headers.delete('x-powered-by');
  response.headers.delete('server');

  // Add additional security headers
  response.headers.set('X-Permitted-Cross-Domain-Policies', 'none');
  response.headers.set('Cross-Origin-Embedder-Policy', 'require-corp');
  response.headers.set('Cross-Origin-Opener-Policy', 'same-origin');
  response.headers.set('Cross-Origin-Resource-Policy', 'same-origin');
  response.headers.set(
    'Permissions-Policy',
    'camera=(), microphone=(), geolocation=(), browsing-topics=(), ' +
      'interest-cohort=(), bluetooth=(), usb=(), payment=(), ' +
      'gyroscope=(), accelerometer=(), magnetometer=(), fullscreen=*'
  );

  // Rate limiting headers
  response.headers.set('X-RateLimit-Limit', '100');
  response.headers.set('X-RateLimit-Remaining', '99');
  response.headers.set('X-RateLimit-Reset', new Date(Date.now() + 3600000).toISOString());

  // Cache control for static assets
  if (
    request.nextUrl.pathname.includes('/_next/static/') ||
    request.nextUrl.pathname.includes('/images/') ||
    request.nextUrl.pathname.includes('/favicon.ico')
  ) {
    response.headers.set('Cache-Control', 'public, max-age=31536000, immutable');
  }

  // No cache for dynamic content
  if (request.nextUrl.pathname.includes('/api/')) {
    response.headers.set('Cache-Control', 'no-store, no-cache, must-revalidate, proxy-revalidate');
    response.headers.set('Pragma', 'no-cache');
    response.headers.set('Expires', '0');
  }

  return response;
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    '/((?!_next/static|_next/image|favicon.ico|public).*)',
  ],
};

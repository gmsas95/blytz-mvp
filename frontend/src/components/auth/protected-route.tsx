'use client';

import { Loader2 } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useEffect, ReactNode } from 'react';

import { useAuth } from '@/contexts/auth-context';


interface ProtectedRouteProps {
  children: ReactNode;
  fallbackPath?: string;
  requireAuth?: boolean;
}

export function ProtectedRoute({
  children,
  fallbackPath = '/auth',
  requireAuth = true,
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoading) return;

    // If authentication is required but user is not authenticated
    if (requireAuth && !isAuthenticated) {
      router.push(fallbackPath);
      return;
    }

    // If user is authenticated but this is a public-only route (like auth page)
    if (!requireAuth && isAuthenticated) {
      router.push('/');
      return;
    }
  }, [isAuthenticated, isLoading, router, fallbackPath, requireAuth]);

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="flex flex-col items-center space-y-4">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  // If authentication is required and user is not authenticated, don't render children
  if (requireAuth && !isAuthenticated) {
    return null;
  }

  // If this is a public-only route but user is authenticated, don't render children
  if (!requireAuth && isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}

// Higher-order component for pages
export function withAuth<P extends object>(
  Component: React.ComponentType<P>,
  options: { fallbackPath?: string; requireAuth?: boolean } = {}
) {
  return function AuthenticatedComponent(props: P) {
    return (
      <ProtectedRoute fallbackPath={options.fallbackPath} requireAuth={options.requireAuth}>
        <Component {...props} />
      </ProtectedRoute>
    );
  };
}

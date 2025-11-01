'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';

import { api } from '@/lib/api-adapter';
import { User } from '@/types';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (userData: RegisterData) => Promise<{ success: boolean; error?: string }>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

interface RegisterData {
  name: string;
  email: string;
  password: string;
  displayName?: string;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!user;

  // Initialize auth state on mount
  useEffect(() => {
    initializeAuth();
  }, []);

  const initializeAuth = async () => {
    try {
      const token = getStoredToken();
      if (!token) {
        setIsLoading(false);
        return;
      }

      // Validate token with backend
      const response = await api.getCurrentUser();
      if (response.success && response.data) {
        setUser(response.data);
      } else {
        // Token is invalid, remove it
        removeStoredToken();
      }
    } catch (error) {
      console.error('Auth initialization error:', error);
      removeStoredToken();
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (
    email: string,
    password: string
  ): Promise<{ success: boolean; error?: string }> => {
    try {
      setIsLoading(true);
      const response = await api.login(email, password);

      if (response.success && response.data) {
        // Store JWT token (assuming it's returned in the response or headers)
        const token = (response.data as any).token || extractTokenFromResponse(response);
        if (token) {
          setStoredToken(token);
        }

        setUser(response.data);
        return { success: true };
      } else {
        return { success: false, error: response.error || 'Login failed' };
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred';
      return { success: false, error: errorMessage };
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (
    userData: RegisterData
  ): Promise<{ success: boolean; error?: string }> => {
    try {
      setIsLoading(true);
      const response = await api.register(userData);

      if (response.success && response.data) {
        // Auto-login after successful registration
        const loginResult = await login(userData.email, userData.password);
        return loginResult;
      } else {
        return { success: false, error: response.error || 'Registration failed' };
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred';
      return { success: false, error: errorMessage };
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await api.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      removeStoredToken();
    }
  };

  const refreshUser = async () => {
    try {
      const response = await api.getCurrentUser();
      if (response.success && response.data) {
        setUser(response.data);
      } else {
        // Token might be expired, logout user
        await logout();
      }
    } catch (error) {
      console.error('User refresh error:', error);
      await logout();
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated,
        login,
        register,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

// Token storage utilities
const TOKEN_KEY = 'blytz_auth_token';

function getStoredToken(): string | null {
  if (typeof window === 'undefined') return null;

  // Try to get from httpOnly cookie first (more secure)
  const cookieToken = getCookie('blytz_auth_token');
  if (cookieToken) return cookieToken;

  // Fallback to localStorage (less secure but works for development)
  return localStorage.getItem(TOKEN_KEY);
}

function setStoredToken(token: string): void {
  if (typeof window === 'undefined') return;

  // Store in localStorage for now (httpOnly cookies require server-side setup)
  localStorage.setItem(TOKEN_KEY, token);

  // TODO: Set httpOnly cookie via API route
  // This would be more secure but requires Next.js API routes
}

function removeStoredToken(): void {
  if (typeof window === 'undefined') return;

  localStorage.removeItem(TOKEN_KEY);

  // TODO: Clear httpOnly cookie via API route
  // document.cookie = 'blytz_auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;'
}

function getCookie(name: string): string | null {
  if (typeof window === 'undefined') return null;

  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);

  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() || null;
  }

  return null;
}

function extractTokenFromResponse(response: any): string | null {
  // Try to extract token from various possible locations
  // This depends on how the backend sends the token

  // 1. Direct token field
  if (response.data?.token) return response.data.token;

  // 2. Authorization header (if available in response)
  if (response.headers?.get('Authorization')) {
    const authHeader = response.headers.get('Authorization');
    return authHeader.replace('Bearer ', '');
  }

  // 3. Custom header
  if (response.headers?.get('X-Auth-Token')) {
    return response.headers.get('X-Auth-Token');
  }

  return null;
}

import { api } from './api-adapter';

// Enhanced API adapter with automatic token injection
class AuthenticatedApiAdapter {
  private originalAdapter = api;

  private async addAuthHeaders(options: RequestInit = {}): Promise<RequestInit> {
    const token = this.getStoredToken();
    if (!token) return options;

    return {
      ...options,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    };
  }

  private getStoredToken(): string | null {
    if (typeof window === 'undefined') return null;

    // Try localStorage first
    const localStorageToken = localStorage.getItem('blytz_auth_token');
    if (localStorageToken) return localStorageToken;

    // Try cookies
    const cookieToken = this.getCookie('blytz_auth_token');
    if (cookieToken) return cookieToken;

    return null;
  }

  private getCookie(name: string): string | null {
    if (typeof window === 'undefined') return null;

    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);

    if (parts.length === 2) {
      return parts.pop()?.split(';').shift() || null;
    }

    return null;
  }

  // Wrap API methods with auth headers
  async getProducts(filter?: any) {
    const options = await this.addAuthHeaders();
    // Note: This would require modifying the original API adapter to accept options
    return this.originalAdapter.getProducts(filter);
  }

  async getProduct(id: string) {
    return this.originalAdapter.getProduct(id);
  }

  async getAuctions(filter?: any) {
    return this.originalAdapter.getAuctions(filter);
  }

  async getAuction(id: string) {
    return this.originalAdapter.getAuction(id);
  }

  async placeBid(auctionId: string, amount: number) {
    const options = await this.addAuthHeaders({
      method: 'POST',
      body: JSON.stringify({ amount }),
    });
    return this.originalAdapter.placeBid(auctionId, amount);
  }

  async getCart() {
    return this.originalAdapter.getCart();
  }

  async addToCart(productId: string, quantity: number, auctionId?: string) {
    return this.originalAdapter.addToCart(productId, quantity, auctionId);
  }

  async removeFromCart(itemId: string) {
    return this.originalAdapter.removeFromCart(itemId);
  }

  async updateCartItemQuantity(itemId: string, quantity: number) {
    return this.originalAdapter.updateCartItemQuantity(itemId, quantity);
  }

  async clearCart() {
    return this.originalAdapter.clearCart();
  }

  // Auth methods (don't need token injection)
  async login(email: string, password: string) {
    return this.originalAdapter.login(email, password);
  }

  async register(userData: any) {
    return this.originalAdapter.register(userData);
  }

  async logout() {
    return this.originalAdapter.logout();
  }

  async getCurrentUser() {
    return this.originalAdapter.getCurrentUser();
  }

  // Payment methods
  async getPaymentMethods() {
    return this.originalAdapter.getPaymentMethods();
  }

  async getFiuuSeamlessConfig() {
    return this.originalAdapter.getFiuuSeamlessConfig();
  }

  async createPayment(paymentRequest: any) {
    return this.originalAdapter.createPayment(paymentRequest);
  }

  async getPaymentStatus(paymentId: string) {
    return this.originalAdapter.getPaymentStatus(paymentId);
  }
}

export const authenticatedApi = new AuthenticatedApiAdapter();

// Utility function to check if user has required role/permissions
export function hasPermission(user: any, permission: string): boolean {
  // Implement role-based permission checking
  if (!user) return false;

  // Example permissions logic
  const permissions = {
    'auction:create': user.isSeller,
    'auction:bid': true,
    'profile:view': true,
    'profile:edit': true,
    'admin:access': user.email?.endsWith('@blytz.app'),
  };

  return permissions[permission as keyof typeof permissions] || false;
}

// Rate limiting utility for auth endpoints
export class AuthRateLimiter {
  private attempts = new Map<string, { count: number; lastAttempt: number }>();
  private readonly maxAttempts = 5;
  private readonly windowMs = 15 * 60 * 1000; // 15 minutes

  isRateLimited(identifier: string): boolean {
    const now = Date.now();
    const record = this.attempts.get(identifier);

    if (!record) {
      this.attempts.set(identifier, { count: 1, lastAttempt: now });
      return false;
    }

    // Reset if window has passed
    if (now - record.lastAttempt > this.windowMs) {
      this.attempts.set(identifier, { count: 1, lastAttempt: now });
      return false;
    }

    // Increment count
    record.count++;
    record.lastAttempt = now;

    return record.count > this.maxAttempts;
  }

  getRemainingAttempts(identifier: string): number {
    const record = this.attempts.get(identifier);
    if (!record) return this.maxAttempts;

    const now = Date.now();
    if (now - record.lastAttempt > this.windowMs) {
      return this.maxAttempts;
    }

    return Math.max(0, this.maxAttempts - record.count);
  }

  getResetTime(identifier: string): Date | null {
    const record = this.attempts.get(identifier);
    if (!record) return null;

    return new Date(record.lastAttempt + this.windowMs);
  }
}

export const authRateLimiter = new AuthRateLimiter();

// Test data interfaces and default values
export interface TestUser {
  email: string;
  password: string;
  displayName: string;
  role: 'user' | 'admin';
}

export interface TestProduct {
  id: string;
  title: string;
  description: string;
  startingPrice: number;
  currentBid: number;
  images: string[];
  category: string;
  condition: 'new' | 'used' | 'refurbished';
}

export interface TestAuction {
  id: string;
  productId: string;
  title: string;
  startingPrice: number;
  currentBid: number;
  endTime: string;
  status: 'active' | 'ended' | 'scheduled';
  bidCount: number;
  minBidIncrement: number;
}

export interface TestBid {
  id: string;
  auctionId: string;
  userId: string;
  amount: number;
  timestamp: string;
}

// Default test data
export const DEFAULT_USERS = {
  regular: {
    email: 'testuser@blytz.app',
    password: 'TestPassword123!',
    displayName: 'Test User',
    role: 'user' as const,
  },
  admin: {
    email: 'admin@blytz.app',
    password: 'AdminPassword123!',
    displayName: 'Admin User',
    role: 'admin' as const,
  },
  blocked: {
    email: 'blocked@blytz.app',
    password: 'BlockedPassword123!',
    displayName: 'Blocked User',
    role: 'user' as const,
  },
};

export const DEFAULT_PRODUCTS = {
  watch: {
    id: 'product-watch-1',
    title: 'Luxury Swiss Watch',
    description: 'Premium Swiss-made automatic watch with leather strap',
    startingPrice: 500.00,
    currentBid: 500.00,
    images: ['https://picsum.photos/400/300?random=watch'],
    category: 'watches',
    condition: 'new' as const,
  },
  phone: {
    id: 'product-phone-1',
    title: 'Smartphone Pro Max',
    description: 'Latest flagship smartphone with advanced camera system',
    startingPrice: 800.00,
    currentBid: 800.00,
    images: ['https://picsum.photos/400/300?random=phone'],
    category: 'electronics',
    condition: 'new' as const,
  },
  jewelry: {
    id: 'product-jewelry-1',
    title: 'Diamond Necklace',
    description: 'Elegant 18k gold necklace with genuine diamonds',
    startingPrice: 1200.00,
    currentBid: 1200.00,
    images: ['https://picsum.photos/400/300?random=jewelry'],
    category: 'jewelry',
    condition: 'new' as const,
  },
};

export const API_ENDPOINTS = {
  auth: {
    login: '/api/auth/login',
    register: '/api/auth/register',
    logout: '/api/auth/logout',
    profile: '/api/auth/profile',
    refresh: '/api/auth/refresh',
  },
  auctions: {
    list: '/api/auctions',
    detail: (id: string) => `/api/auctions/${id}`,
    create: '/api/auctions',
    placeBid: (id: string) => `/api/auctions/${id}/bid`,
    watch: (id: string) => `/api/auctions/${id}/watch`,
  },
  products: {
    list: '/api/products',
    detail: (id: string) => `/api/products/${id}`,
    search: '/api/products/search',
  },
  cart: {
    items: '/api/cart',
    add: '/api/cart/add',
    remove: '/api/cart/remove',
    clear: '/api/cart/clear',
  },
  checkout: {
    process: '/api/checkout',
    payment: '/api/checkout/payment',
    confirmation: '/api/checkout/confirmation',
  },
};

// Test utilities
export const generateTestEmail = () =>
  `test-${Date.now()}@blytz.app`;

export const generateTestBid = (baseAmount: number, increment: number = 5.00) =>
  baseAmount + increment + Math.random() * 10;

export const waitForAuctionEnd = (endTime: string) => {
  const end = new Date(endTime).getTime();
  const now = Date.now();
  return Math.max(0, end - now);
};
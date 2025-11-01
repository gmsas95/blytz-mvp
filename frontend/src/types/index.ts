export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  isSeller: boolean;
  storeName?: string;
  storeDescription?: string;
  rating: number;
  totalSales: number;
  createdAt: string;
}

export interface Product {
  id: string;
  title: string;
  description: string;
  price: number;
  originalPrice?: number;
  images: string[];
  category: string;
  seller: User;
  auction?: Auction;
  specifications: Record<string, string>;
  inStock: boolean;
  stockQuantity: number;
  createdAt: string;
  updatedAt: string;
}

export interface Auction {
  id: string;
  productId: string;
  product: Product;
  startingPrice: number;
  currentBid: number;
  minBidIncrement: number;
  startTime: string;
  endTime: string;
  status: 'scheduled' | 'active' | 'ended';
  totalBids: number;
  participants: number;
  isLive: boolean;
  streamUrl?: string;
  winner?: User;
  bids: Bid[];
}

export interface Bid {
  id: string;
  auctionId: string;
  userId: string;
  amount: number;
  timestamp: string;
  user: User;
}

export interface CartItem {
  id: string;
  product: Product;
  quantity: number;
  selectedAuction?: Auction;
}

export interface Cart {
  id: string;
  userId: string;
  items: CartItem[];
  total: number;
  itemCount: number;
}

export interface Livestream {
  id: string;
  title: string;
  description: string;
  streamUrl: string;
  thumbnail: string;
  seller: User;
  products: Product[];
  viewers: number;
  likes: number;
  isLive: boolean;
  startedAt: string;
  duration: number;
  status: 'scheduled' | 'live' | 'ended';
}

export interface Order {
  id: string;
  userId: string;
  items: CartItem[];
  total: number;
  status: 'pending' | 'confirmed' | 'shipped' | 'delivered' | 'cancelled';
  shippingAddress: Address;
  paymentMethod: PaymentMethod;
  createdAt: string;
  updatedAt: string;
}

export interface Address {
  id: string;
  userId: string;
  name: string;
  street: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  isDefault: boolean;
}

export interface PaymentMethod {
  id: string;
  userId: string;
  type: 'card' | 'paypal' | 'bank';
  provider: string;
  last4?: string;
  expiryMonth?: number;
  expiryYear?: number;
  isDefault: boolean;
}

export interface Notification {
  id: string;
  userId: string;
  title: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
  isRead: boolean;
  createdAt: string;
  actionUrl?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export type ProductFilter = {
  category?: string;
  priceRange?: [number, number];
  seller?: string;
  search?: string;
  sortBy?: 'price' | 'name' | 'created' | 'popularity';
  sortOrder?: 'asc' | 'desc';
};

export type AuctionFilter = {
  status?: 'scheduled' | 'active' | 'ended';
  category?: string;
  priceRange?: [number, number];
  seller?: string;
  search?: string;
};

export type LivestreamFilter = {
  status?: 'scheduled' | 'live' | 'ended';
  seller?: string;
  search?: string;
};

// Payment related types
export interface PaymentMethodInfo {
  id: string;
  name: string;
  type: string;
  description: string;
  available: boolean;
  fee: number;
  icon?: string;
}

export interface FiuuSeamlessConfig {
  version: string;
  actionType: string;
  merchantID: string;
  paymentMethod: string;
  orderNumber: string;
  amount: number;
  currency: string;
  productDescription: string;
  userName: string;
  userEmail: string;
  userContact: string;
  remark: string;
  lang: string;
  vcode: string;
  callbackURL: string;
  returnURL: string;
  backgroundUrl: string;
  sandbox: boolean;
  scriptUrl: string;
}

export interface PaymentRequest {
  amount: number;
  currency: string;
  paymentMethod: string;
  orderNumber: string;
  description: string;
  returnUrl: string;
  cancelUrl: string;
  webhookUrl: string;
}

export interface PaymentResponse {
  paymentId: string;
  orderNumber: string;
  amount: number;
  currency: string;
  status: string;
  paymentMethod: string;
  redirectUrl: string;
  createdAt: string;
  expiresAt: string;
}

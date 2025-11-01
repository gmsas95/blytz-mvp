import { User, Product, Auction, Livestream, Cart } from '@/types';

// Mock Users
export const mockUsers: User[] = [
  {
    id: '1',
    name: 'Sarah Johnson',
    email: 'sarah@jstyle.com',
    avatar:
      'https://images.unsplash.com/photo-1494790108755-2616b612b5bc?w=150&h=150&fit=crop&crop=face',
    isSeller: true,
    storeName: 'JStyle Boutique',
    storeDescription: 'Curated fashion and lifestyle products',
    rating: 4.8,
    totalSales: 1234,
    createdAt: '2023-01-15T10:00:00Z',
  },
  {
    id: '2',
    name: 'Mike Chen',
    email: 'mike@techhub.com',
    avatar:
      'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
    isSeller: true,
    storeName: 'TechHub Deals',
    storeDescription: 'Latest gadgets and electronics',
    rating: 4.6,
    totalSales: 856,
    createdAt: '2023-03-20T14:30:00Z',
  },
  {
    id: '3',
    name: 'Emma Wilson',
    email: 'emma@homecomfort.com',
    avatar:
      'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face',
    isSeller: true,
    storeName: 'Home Comfort',
    storeDescription: 'Beautiful home decor and furniture',
    rating: 4.9,
    totalSales: 2103,
    createdAt: '2022-11-08T09:15:00Z',
  },
];

// Mock Products
export const mockProducts: Product[] = [
  {
    id: '1',
    title: 'Vintage Leather Handbag',
    description:
      'Authentic vintage leather handbag with gold hardware. Perfect condition with original dust bag.',
    price: 299.99,
    originalPrice: 450.0,
    images: [
      'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1590874103328-eac38a683ce7?w=400&h=400&fit=crop',
    ],
    category: 'Fashion',
    seller: mockUsers[0],
    specifications: {
      Material: 'Genuine Leather',
      Color: 'Brown',
      Dimensions: '12" x 8" x 4"',
      Condition: 'Excellent',
    },
    inStock: true,
    stockQuantity: 1,
    createdAt: '2024-01-10T16:00:00Z',
    updatedAt: '2024-01-10T16:00:00Z',
  },
  {
    id: '2',
    title: 'Smart Home Security Camera',
    description:
      '4K wireless security camera with night vision, motion detection, and cloud storage.',
    price: 189.99,
    originalPrice: 249.99,
    images: [
      'https://images.unsplash.com/photo-1558089687-f282ffcbc126?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1584464491033-06628f3a6b7b?w=400&h=400&fit=crop',
    ],
    category: 'Electronics',
    seller: mockUsers[1],
    specifications: {
      Resolution: '4K Ultra HD',
      Connectivity: 'WiFi 6',
      'Night Vision': 'Up to 30ft',
      Storage: 'Cloud + Local SD',
    },
    inStock: true,
    stockQuantity: 15,
    createdAt: '2024-01-12T11:30:00Z',
    updatedAt: '2024-01-12T11:30:00Z',
  },
  {
    id: '3',
    title: 'Mid-Century Modern Coffee Table',
    description:
      'Beautiful teak wood coffee table with clean lines and tapered legs. Iconic mid-century design.',
    price: 599.0,
    images: [
      'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=400&h=400&fit=crop',
    ],
    category: 'Home & Garden',
    seller: mockUsers[2],
    specifications: {
      Material: 'Solid Teak Wood',
      Style: 'Mid-Century Modern',
      Dimensions: '48" x 24" x 16"',
      Finish: 'Natural Oil',
    },
    inStock: true,
    stockQuantity: 3,
    createdAt: '2024-01-08T14:20:00Z',
    updatedAt: '2024-01-08T14:20:00Z',
  },
  {
    id: '4',
    title: 'Wireless Noise-Canceling Headphones',
    description:
      'Premium wireless headphones with active noise cancellation, 30-hour battery life.',
    price: 349.99,
    originalPrice: 399.99,
    images: [
      'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1484704849700-f032a568e944?w=400&h=400&fit=crop',
    ],
    category: 'Electronics',
    seller: mockUsers[1],
    specifications: {
      'Battery Life': '30 hours',
      'Noise Cancellation': 'Active',
      Connectivity: 'Bluetooth 5.2',
      Weight: '250g',
    },
    inStock: true,
    stockQuantity: 8,
    createdAt: '2024-01-14T09:45:00Z',
    updatedAt: '2024-01-14T09:45:00Z',
  },
  {
    id: '5',
    title: 'Ceramic Table Lamp Set',
    description:
      'Pair of handcrafted ceramic table lamps with linen shades. Perfect for bedside tables.',
    price: 179.99,
    images: [
      'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop',
      'https://images.unsplash.com/photo-1524484485831-a92ffc0de03f?w=400&h=400&fit=crop',
    ],
    category: 'Home & Garden',
    seller: mockUsers[2],
    specifications: {
      Material: 'Ceramic + Linen',
      Height: '24 inches',
      'Shade Color': 'Natural Linen',
      'Bulb Type': 'LED Compatible',
    },
    inStock: true,
    stockQuantity: 5,
    createdAt: '2024-01-11T13:15:00Z',
    updatedAt: '2024-01-11T13:15:00Z',
  },
];

// Mock Auctions
export const mockAuctions: Auction[] = [
  {
    id: '1',
    productId: '1',
    product: mockProducts[0],
    startingPrice: 199.99,
    currentBid: 275.0,
    minBidIncrement: 10.0,
    startTime: new Date(Date.now() - 30 * 60 * 1000).toISOString(), // 30 minutes ago
    endTime: new Date(Date.now() + 25 * 60 * 1000).toISOString(), // 25 minutes from now
    status: 'active',
    totalBids: 12,
    participants: 8,
    isLive: true,
    streamUrl: `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/stream/auction-1`,
    bids: [
      {
        id: '1',
        auctionId: '1',
        userId: '4',
        amount: 275.0,
        timestamp: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
        user: {
          id: '4',
          name: 'John Doe',
          email: 'john@example.com',
          isSeller: false,
          rating: 0,
          totalSales: 0,
          createdAt: '2024-01-01T00:00:00Z',
        },
      },
      {
        id: '2',
        auctionId: '1',
        userId: '5',
        amount: 265.0,
        timestamp: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
        user: {
          id: '5',
          name: 'Jane Smith',
          email: 'jane@example.com',
          isSeller: false,
          rating: 0,
          totalSales: 0,
          createdAt: '2024-01-01T00:00:00Z',
        },
      },
    ],
  },
  {
    id: '2',
    productId: '2',
    product: mockProducts[1],
    startingPrice: 149.99,
    currentBid: 185.5,
    minBidIncrement: 5.0,
    startTime: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(), // 2 hours from now
    endTime: new Date(Date.now() + 3 * 60 * 60 * 1000).toISOString(), // 3 hours from now
    status: 'scheduled',
    totalBids: 0,
    participants: 0,
    isLive: false,
    bids: [],
  },
];

// Mock Livestreams
export const mockLivestreams: Livestream[] = [
  {
    id: '1',
    title: 'Vintage Fashion Show - Spring Collection',
    description:
      'Join Sarah for an exclusive showcase of vintage fashion pieces with live auctions!',
    streamUrl: `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/stream/live-1`,
    thumbnail: 'https://images.unsplash.com/photo-1558769132-cb1aea458c5e?w=800&h=450&fit=crop',
    seller: mockUsers[0],
    products: [mockProducts[0]],
    viewers: 234,
    likes: 156,
    isLive: true,
    startedAt: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
    duration: 900,
    status: 'live',
  },
  {
    id: '2',
    title: 'Tech Review & Gadget Auctions',
    description: 'Latest tech reviews with special auction prices on featured products!',
    streamUrl: 'https://stream.example.com/live-2',
    thumbnail: 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=800&h=450&fit=crop',
    seller: mockUsers[1],
    products: [mockProducts[1], mockProducts[3]],
    viewers: 189,
    likes: 98,
    isLive: true,
    startedAt: new Date(Date.now() - 45 * 60 * 1000).toISOString(),
    duration: 2700,
    status: 'live',
  },
];

export const mockCart: Cart = {
  id: '1',
  userId: 'current-user',
  items: [
    {
      id: '1',
      product: mockProducts[0],
      quantity: 1,
      selectedAuction: mockAuctions[0],
    },
    {
      id: '2',
      product: mockProducts[1],
      quantity: 2,
    },
  ],
  total: 675.97,
  itemCount: 3,
};

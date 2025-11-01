import {
  Product,
  Auction,
  Livestream,
  Cart,
  User,
  Bid,
  ProductFilter,
  AuctionFilter,
  LivestreamFilter,
  ApiResponse,
  PaginatedResponse,
  PaymentMethodInfo,
  FiuuSeamlessConfig,
  PaymentRequest,
  PaymentResponse
} from '@/types'
import { mockProducts, mockAuctions, mockLivestreams, mockCart, mockUsers } from '@/data/mock-data'

export interface ApiAdapter {
  // Products
  getProducts(filter?: ProductFilter): Promise<ApiResponse<PaginatedResponse<Product>>>
  getProduct(id: string): Promise<ApiResponse<Product>>
  getFeaturedProducts(): Promise<ApiResponse<Product[]>>

  // Auctions
  getAuctions(filter?: AuctionFilter): Promise<ApiResponse<PaginatedResponse<Auction>>>
  getAuction(id: string): Promise<ApiResponse<Auction>>
  placeBid(auctionId: string, amount: number): Promise<ApiResponse<Bid>>
  getActiveAuctions(): Promise<ApiResponse<Auction[]>>

  // Livestreams
  getLivestreams(filter?: LivestreamFilter): Promise<ApiResponse<PaginatedResponse<Livestream>>>
  getLivestream(id: string): Promise<ApiResponse<Livestream>>
  getActiveLivestreams(): Promise<ApiResponse<Livestream[]>>

  // Cart
  getCart(): Promise<ApiResponse<Cart>>
  addToCart(productId: string, quantity: number, auctionId?: string): Promise<ApiResponse<Cart>>
  removeFromCart(itemId: string): Promise<ApiResponse<Cart>>
  updateCartItemQuantity(itemId: string, quantity: number): Promise<ApiResponse<Cart>>
  clearCart(): Promise<ApiResponse<Cart>>

  // Auth
  login(email: string, password: string): Promise<ApiResponse<User>>
  register(userData: Partial<User>): Promise<ApiResponse<User>>
  logout(): Promise<ApiResponse<void>>
  getCurrentUser(): Promise<ApiResponse<User>>

  // Payments
  getPaymentMethods(): Promise<ApiResponse<PaymentMethodInfo[]>>
  getFiuuSeamlessConfig(): Promise<ApiResponse<FiuuSeamlessConfig>>
  createPayment(paymentRequest: PaymentRequest): Promise<ApiResponse<PaymentResponse>>
  getPaymentStatus(paymentId: string): Promise<ApiResponse<PaymentResponse>>
}

// Mock API Adapter
export class MockApiAdapter implements ApiAdapter {
  private delay(ms: number = 100): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  async getProducts(filter?: ProductFilter): Promise<ApiResponse<PaginatedResponse<Product>>> {
    await this.delay()

    let filteredProducts = [...mockProducts]

    if (filter?.search) {
      const searchLower = filter.search.toLowerCase()
      filteredProducts = filteredProducts.filter(p =>
        p.title.toLowerCase().includes(searchLower) ||
        p.description.toLowerCase().includes(searchLower)
      )
    }

    if (filter?.category) {
      filteredProducts = filteredProducts.filter(p => p.category === filter.category)
    }

    if (filter?.priceRange) {
      filteredProducts = filteredProducts.filter(p =>
        p.price >= filter.priceRange![0] && p.price <= filter.priceRange![1]
      )
    }

    if (filter?.sortBy) {
      filteredProducts.sort((a, b) => {
        const aVal = filter.sortBy === 'price' ? a.price : a.title
        const bVal = filter.sortBy === 'price' ? b.price : b.title
        const comparison = aVal > bVal ? 1 : -1
        return filter.sortOrder === 'desc' ? -comparison : comparison
      })
    }

    return {
      success: true,
      data: {
        items: filteredProducts,
        total: filteredProducts.length,
        page: 1,
        limit: 20,
        totalPages: 1
      }
    }
  }

  async getProduct(id: string): Promise<ApiResponse<Product>> {
    await this.delay()
    const product = mockProducts.find(p => p.id === id)

    if (!product) {
      return { success: false, error: 'Product not found' }
    }

    return { success: true, data: product }
  }

  async getFeaturedProducts(): Promise<ApiResponse<Product[]>> {
    await this.delay()
    return { success: true, data: mockProducts.slice(0, 3) }
  }

  async getAuctions(filter?: AuctionFilter): Promise<ApiResponse<PaginatedResponse<Auction>>> {
    await this.delay()

    let filteredAuctions = [...mockAuctions]

    if (filter?.status) {
      filteredAuctions = filteredAuctions.filter(a => a.status === filter.status)
    }

    return {
      success: true,
      data: {
        items: filteredAuctions,
        total: filteredAuctions.length,
        page: 1,
        limit: 20,
        totalPages: 1
      }
    }
  }

  async getAuction(id: string): Promise<ApiResponse<Auction>> {
    await this.delay()
    const auction = mockAuctions.find(a => a.id === id)

    if (!auction) {
      return { success: false, error: 'Auction not found' }
    }

    return { success: true, data: auction }
  }

  async placeBid(auctionId: string, amount: number): Promise<ApiResponse<Bid>> {
    await this.delay(800)
    const auction = mockAuctions.find(a => a.id === auctionId)

    if (!auction) {
      return { success: false, error: 'Auction not found' }
    }

    if (amount < auction.currentBid + auction.minBidIncrement) {
      return { success: false, error: 'Bid must be at least $' + (auction.currentBid + auction.minBidIncrement).toFixed(2) }
    }

    const newBid: Bid = {
      id: Date.now().toString(),
      auctionId,
      userId: 'current-user',
      amount,
      timestamp: new Date().toISOString(),
      user: mockUsers[0] // Mock current user
    }

    auction.currentBid = amount
    auction.totalBids += 1
    auction.bids.unshift(newBid)

    return { success: true, data: newBid }
  }

  async getActiveAuctions(): Promise<ApiResponse<Auction[]>> {
    await this.delay()
    const activeAuctions = mockAuctions.filter(a => a.status === 'active')
    return { success: true, data: activeAuctions }
  }

  async getLivestreams(filter?: LivestreamFilter): Promise<ApiResponse<PaginatedResponse<Livestream>>> {
    await this.delay()

    let filteredLivestreams = [...mockLivestreams]

    if (filter?.status) {
      filteredLivestreams = filteredLivestreams.filter(l => l.status === filter.status)
    }

    return {
      success: true,
      data: {
        items: filteredLivestreams,
        total: filteredLivestreams.length,
        page: 1,
        limit: 20,
        totalPages: 1
      }
    }
  }

  async getLivestream(id: string): Promise<ApiResponse<Livestream>> {
    await this.delay()
    const livestream = mockLivestreams.find(l => l.id === id)

    if (!livestream) {
      return { success: false, error: 'Livestream not found' }
    }

    return { success: true, data: livestream }
  }

  async getActiveLivestreams(): Promise<ApiResponse<Livestream[]>> {
    await this.delay()
    const activeLivestreams = mockLivestreams.filter(l => l.status === 'live')
    return { success: true, data: activeLivestreams }
  }

  async getCart(): Promise<ApiResponse<Cart>> {
    await this.delay()
    return { success: true, data: mockCart }
  }

  async addToCart(productId: string, quantity: number, auctionId?: string): Promise<ApiResponse<Cart>> {
    await this.delay()

    const product = mockProducts.find(p => p.id === productId)
    if (!product) {
      return { success: false, error: 'Product not found' }
    }

    const newItem = {
      id: Date.now().toString(),
      product,
      quantity,
      selectedAuction: auctionId ? mockAuctions.find(a => a.id === auctionId) : undefined
    }

    mockCart.items.push(newItem)
    mockCart.itemCount += quantity
    mockCart.total += product.price * quantity

    return { success: true, data: mockCart }
  }

  async removeFromCart(itemId: string): Promise<ApiResponse<Cart>> {
    await this.delay()

    const itemIndex = mockCart.items.findIndex(item => item.id === itemId)
    if (itemIndex === -1) {
      return { success: false, error: 'Cart item not found' }
    }

    const item = mockCart.items[itemIndex]
    mockCart.items.splice(itemIndex, 1)
    mockCart.itemCount -= item.quantity
    mockCart.total -= item.product.price * item.quantity

    return { success: true, data: mockCart }
  }

  async updateCartItemQuantity(itemId: string, quantity: number): Promise<ApiResponse<Cart>> {
    await this.delay()

    const item = mockCart.items.find(item => item.id === itemId)
    if (!item) {
      return { success: false, error: 'Cart item not found' }
    }

    const quantityDiff = quantity - item.quantity
    const priceDiff = item.product.price * quantityDiff

    item.quantity = quantity
    mockCart.itemCount += quantityDiff
    mockCart.total += priceDiff

    return { success: true, data: mockCart }
  }

  async clearCart(): Promise<ApiResponse<Cart>> {
    await this.delay()

    mockCart.items = []
    mockCart.itemCount = 0
    mockCart.total = 0

    return { success: true, data: mockCart }
  }

  async login(email: string, password: string): Promise<ApiResponse<User>> {
    await this.delay(1000)

    // Input validation
    if (!email || !password) {
      return { success: false, error: 'Email and password are required' }
    }

    if (!email.includes('@') || email.length < 5) {
      return { success: false, error: 'Invalid email format' }
    }

    if (password.length < 6) {
      return { success: false, error: 'Password must be at least 6 characters' }
    }

    // Development demo credentials - only available in development mode
    if (process.env.NODE_ENV === 'development' && email === 'demo@blytz.app' && password === 'demo123') {
      return { success: true, data: mockUsers[0] }
    }

    return { success: false, error: 'Invalid credentials' }
  }

  async register(userData: Partial<User>): Promise<ApiResponse<User>> {
    await this.delay(1000)

    const newUser: User = {
      id: Date.now().toString(),
      name: userData.name || '',
      email: userData.email || '',
      isSeller: userData.isSeller || false,
      rating: 0,
      totalSales: 0,
      createdAt: new Date().toISOString()
    }

    return { success: true, data: newUser }
  }

  async logout(): Promise<ApiResponse<void>> {
    await this.delay(500)
    return { success: true }
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    await this.delay(300)
    return { success: true, data: mockUsers[0] }
  }

  // Payment methods
  async getPaymentMethods(): Promise<ApiResponse<PaymentMethodInfo[]>> {
    await this.delay()
    return {
      success: true,
      data: [
        {
          id: 'fpx',
          name: 'FPX Online Banking',
          type: 'bank_transfer',
          icon: '🏦',
          description: 'Pay directly from your bank account',
          available: true,
          fee: 0
        },
        {
          id: 'tng',
          name: 'Touch \'n Go eWallet',
          type: 'ewallet',
          icon: '📱',
          description: 'Instant payment with TNG eWallet',
          available: true,
          fee: 0.50
        },
        {
          id: 'grabpay',
          name: 'GrabPay',
          type: 'ewallet',
          icon: '🚗',
          description: 'Pay with GrabPay eWallet',
          available: true,
          fee: 0.50
        },
        {
          id: 'credit_card',
          name: 'Credit/Debit Card',
          type: 'card',
          icon: '💳',
          description: 'Visa, Mastercard, AMEX',
          available: true,
          fee: 1.5
        },
        {
          id: 'shopeepay',
          name: 'ShopeePay',
          type: 'ewallet',
          icon: '🛒',
          description: 'Quick payment with ShopeePay',
          available: true,
          fee: 0.50
        }
      ]
    }
  }

  async getFiuuSeamlessConfig(): Promise<ApiResponse<FiuuSeamlessConfig>> {
    await this.delay()
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/public/seamless/config?order_id=ORDER_${Date.now()}&amount=10000&bill_name=Test%20User&bill_email=test@example.com&bill_mobile=01234567890&bill_desc=Test%20Payment&channel=FPX`)
      const data = await response.json()
      
      if (data.success) {
        return {
          success: true,
          data: data.data
        }
      } else {
        // Fallback to mock if backend fails
        const isSandbox = process.env.NEXT_PUBLIC_FIUU_SANDBOX !== 'false'
        return {
          success: true,
          data: {
            version: '7.5.0',
            actionType: 'Pay',
            merchantID: 'MERCHANT_ID',
            paymentMethod: 'all',
            orderNumber: `ORDER_${Date.now()}`,
            amount: 100.00,
            currency: 'MYR',
            productDescription: 'Blytz Auction Purchase',
            userName: 'John Doe',
            userEmail: 'john@example.com',
            userContact: '+60123456789',
            remark: 'Payment for auction items',
            lang: 'en',
            vcode: 'MOCK_VCODE_123456',
            callbackURL: `${process.env.NEXT_PUBLIC_API_URL}/api/v1/webhooks/fiuu`,
            returnURL: 'https://blytz.app/checkout/success',
            backgroundUrl: `${process.env.NEXT_PUBLIC_API_URL}/api/v1/webhooks/fiuu`,
            sandbox: isSandbox,
            scriptUrl: isSandbox 
              ? 'https://sandbox.merchant.razer.com/RMS2/IPGSeamless/IPGSeamless.js'
              : 'https://api.merchant.razer.com/RMS2/IPGSeamless/IPGSeamless.js'
          }
        }
      }
    } catch (error) {
      // Fallback to mock on network error
      const isSandbox = process.env.NEXT_PUBLIC_FIUU_SANDBOX !== 'false'
      return {
        success: true,
        data: {
          version: '7.5.0',
          actionType: 'Pay',
          merchantID: 'MERCHANT_ID',
          paymentMethod: 'all',
          orderNumber: `ORDER_${Date.now()}`,
          amount: 100.00,
          currency: 'MYR',
          productDescription: 'Blytz Auction Purchase',
          userName: 'John Doe',
          userEmail: 'john@example.com',
          userContact: '+60123456789',
          remark: 'Payment for auction items',
          lang: 'en',
          vcode: 'MOCK_VCODE_123456',
          callbackURL: `${process.env.NEXT_PUBLIC_API_URL}/api/v1/webhooks/fiuu`,
          returnURL: 'https://blytz.app/checkout/success',
          backgroundUrl: `${process.env.NEXT_PUBLIC_API_URL}/api/v1/webhooks/fiuu`,
          sandbox: isSandbox,
          scriptUrl: isSandbox 
            ? 'https://sandbox.merchant.razer.com/RMS2/IPGSeamless/IPGSeamless.js'
            : 'https://api.merchant.razer.com/RMS2/IPGSeamless/IPGSeamless.js'
        }
      }
    }
  }

  async createPayment(paymentRequest: PaymentRequest): Promise<ApiResponse<PaymentResponse>> {
    await this.delay(2000)
    
    const paymentId = `PAY_${Date.now()}`
    const orderNumber = `ORDER_${Date.now()}`
    
    return {
      success: true,
      data: {
        paymentId,
        orderNumber,
        amount: paymentRequest.amount,
        currency: 'MYR',
        status: 'pending',
        paymentMethod: paymentRequest.paymentMethod,
        redirectUrl: `https://blytz.app/checkout/success?payment_id=${paymentId}`,
        createdAt: new Date().toISOString(),
        expiresAt: new Date(Date.now() + 30 * 60 * 1000).toISOString()
      }
    }
  }

  async getPaymentStatus(paymentId: string): Promise<ApiResponse<PaymentResponse>> {
    await this.delay()
    
    // Mock payment status - in real implementation this would check the actual payment status
    return {
      success: true,
      data: {
        paymentId,
        orderNumber: `ORDER_${Date.now()}`,
        amount: 100.00,
        currency: 'MYR',
        status: Math.random() > 0.3 ? 'success' : 'pending',
        paymentMethod: 'fpx',
        redirectUrl: '',
        createdAt: new Date().toISOString(),
        expiresAt: new Date(Date.now() + 30 * 60 * 1000).toISOString()
      }
    }
  }
}

// Remote API Adapter (for production)
export class RemoteApiAdapter implements ApiAdapter {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  private async fetchApi<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers
        }
      })

      if (!response.ok) {
        let errorMessage = `HTTP ${response.status}: ${response.statusText}`
        
        try {
          const errorText = await response.text()
          if (errorText) {
            errorMessage = errorText
          }
        } catch {
          // Use default error message if response text is not readable
        }
        
        return { success: false, error: errorMessage }
      }

      const data = await response.json()
      return { success: true, data }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Network error occurred'
      return { success: false, error: errorMessage }
    }
  }

  async getProducts(filter?: ProductFilter): Promise<ApiResponse<PaginatedResponse<Product>>> {
    const params = new URLSearchParams()
    if (filter?.search) params.append('search', filter.search)
    if (filter?.category) params.append('category', filter.category)
    if (filter?.sortBy) params.append('sortBy', filter.sortBy)
    if (filter?.sortOrder) params.append('sortOrder', filter.sortOrder)

    return this.fetchApi(`/products?${params}`)
  }

  async getProduct(id: string): Promise<ApiResponse<Product>> {
    return this.fetchApi(`/products/${id}`)
  }

  async getFeaturedProducts(): Promise<ApiResponse<Product[]>> {
    return this.fetchApi('/products/featured')
  }

  async getAuctions(filter?: AuctionFilter): Promise<ApiResponse<PaginatedResponse<Auction>>> {
    const params = new URLSearchParams()
    if (filter?.status) params.append('status', filter.status)

    return this.fetchApi(`/auctions?${params}`)
  }

  async getAuction(id: string): Promise<ApiResponse<Auction>> {
    return this.fetchApi(`/auctions/${id}`)
  }

  async placeBid(auctionId: string, amount: number): Promise<ApiResponse<Bid>> {
    return this.fetchApi(`/auctions/${auctionId}/bid`, {
      method: 'POST',
      body: JSON.stringify({ amount })
    })
  }

  async getActiveAuctions(): Promise<ApiResponse<Auction[]>> {
    return this.fetchApi('/auctions/active')
  }

  async getLivestreams(filter?: LivestreamFilter): Promise<ApiResponse<PaginatedResponse<Livestream>>> {
    const params = new URLSearchParams()
    if (filter?.status) params.append('status', filter.status)

    return this.fetchApi(`/livestreams?${params}`)
  }

  async getLivestream(id: string): Promise<ApiResponse<Livestream>> {
    return this.fetchApi(`/livestreams/${id}`)
  }

  async getActiveLivestreams(): Promise<ApiResponse<Livestream[]>> {
    return this.fetchApi('/livestreams/active')
  }

  async getCart(): Promise<ApiResponse<Cart>> {
    return this.fetchApi('/cart')
  }

  async addToCart(productId: string, quantity: number, auctionId?: string): Promise<ApiResponse<Cart>> {
    return this.fetchApi('/cart/add', {
      method: 'POST',
      body: JSON.stringify({ productId, quantity, auctionId })
    })
  }

  async removeFromCart(itemId: string): Promise<ApiResponse<Cart>> {
    return this.fetchApi(`/cart/remove/${itemId}`, { method: 'DELETE' })
  }

  async updateCartItemQuantity(itemId: string, quantity: number): Promise<ApiResponse<Cart>> {
    return this.fetchApi(`/cart/update/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify({ quantity })
    })
  }

  async clearCart(): Promise<ApiResponse<Cart>> {
    return this.fetchApi('/cart/clear', { method: 'DELETE' })
  }

  async login(email: string, password: string): Promise<ApiResponse<User>> {
    return this.fetchApi('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    })
  }

  async register(userData: Partial<User>): Promise<ApiResponse<User>> {
    return this.fetchApi('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData)
    })
  }

  async logout(): Promise<ApiResponse<void>> {
    return this.fetchApi('/auth/logout', { method: 'POST' })
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    return this.fetchApi('/auth/me')
  }

  // Payment methods
  async getPaymentMethods(): Promise<ApiResponse<PaymentMethodInfo[]>> {
    return this.fetchApi('/payments/methods')
  }

  async getFiuuSeamlessConfig(): Promise<ApiResponse<FiuuSeamlessConfig>> {
    try {
      // Call payment service directly using container name in Docker network
      const response = await fetch(`http://blytz-payment-prod:8086/api/v1/public/seamless/config?order_id=TEST123&amount=10000&bill_name=Test%20User&bill_email=test@example.com&bill_mobile=01234567890&bill_desc=Test%20Payment&channel=FPX`)
      const data = await response.json()
      
      if (data.success) {
        return {
          success: true,
          data: data.data
        }
      } else {
        throw new Error('Failed to get Fiuu config')
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      }
    }
  }

  async createPayment(paymentRequest: PaymentRequest): Promise<ApiResponse<PaymentResponse>> {
    return this.fetchApi('/payments/create', {
      method: 'POST',
      body: JSON.stringify(paymentRequest)
    })
  }

  async getPaymentStatus(paymentId: string): Promise<ApiResponse<PaymentResponse>> {
    return this.fetchApi(`/payments/${paymentId}/status`)
  }
}

// Factory function to create appropriate adapter
export function createApiAdapter(): ApiAdapter {
  const mode = process.env.MODE || 'mock'

  if (mode === 'remote') {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL
    if (!baseUrl) {
      throw new Error('NEXT_PUBLIC_API_URL environment variable is required in remote mode')
    }
    return new RemoteApiAdapter(baseUrl)
  }

  return new MockApiAdapter()
}

export const api = createApiAdapter()
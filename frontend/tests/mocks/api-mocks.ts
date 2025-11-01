import { Page, Route } from '@playwright/test';
import { DEFAULT_USERS, DEFAULT_PRODUCTS, API_ENDPOINTS } from '../utils/test-data';

export class ApiMocks {
  constructor(private page: Page) {}

  /**
   * Mock authentication endpoints
   */
  mockAuthEndpoints(): void {
    // Mock login endpoint
    this.page.route(API_ENDPOINTS.auth.login, async (route: Route) => {
      const request = route.request();
      const postData = request.postDataJSON();

      // Check if credentials match test users
      const validUser = Object.values(DEFAULT_USERS).find(
        user => user.email === postData.email && user.password === postData.password
      );

      if (validUser) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              user: {
                id: `user-${validUser.email.split('@')[0]}`,
                email: validUser.email,
                displayName: validUser.displayName,
                role: validUser.role,
              },
              token: 'mock-jwt-token',
              refreshToken: 'mock-refresh-token',
            },
          }),
        });
      } else {
        await route.fulfill({
          status: 401,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            error: 'Invalid credentials',
          }),
        });
      }
    });

    // Mock register endpoint
    this.page.route(API_ENDPOINTS.auth.register, async (route: Route) => {
      const request = route.request();
      const postData = request.postDataJSON();

      // Simple validation
      if (postData.email && postData.password && postData.displayName) {
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              user: {
                id: `user-${Date.now()}`,
                email: postData.email,
                displayName: postData.displayName,
                role: 'user',
              },
              token: 'mock-jwt-token-new-user',
              refreshToken: 'mock-refresh-token-new-user',
            },
          }),
        });
      } else {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            error: 'Missing required fields',
          }),
        });
      }
    });

    // Mock logout endpoint
    this.page.route(API_ENDPOINTS.auth.logout, async (route: Route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          message: 'Logged out successfully',
        }),
      });
    });

    // Mock profile endpoint
    this.page.route(API_ENDPOINTS.auth.profile, async (route: Route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: {
            id: 'user-test',
            email: 'testuser@blytz.app',
            displayName: 'Test User',
            role: 'user',
            createdAt: new Date().toISOString(),
          },
        }),
      });
    });
  }

  /**
   * Mock auction endpoints
   */
  mockAuctionEndpoints(): void {
    // Mock auctions list
    this.page.route(API_ENDPOINTS.auctions.list, async (route: Route) => {
      const mockAuctions = [
        {
          id: 'auction-1',
          productId: DEFAULT_PRODUCTS.watch.id,
          title: DEFAULT_PRODUCTS.watch.title + ' - Live Auction',
          currentBid: 550.00,
          endTime: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(), // 2 hours from now
          status: 'active',
          bidCount: 3,
          product: DEFAULT_PRODUCTS.watch,
        },
        {
          id: 'auction-2',
          productId: DEFAULT_PRODUCTS.phone.id,
          title: DEFAULT_PRODUCTS.phone.title + ' - Live Auction',
          currentBid: 850.00,
          endTime: new Date(Date.now() + 1 * 60 * 60 * 1000).toISOString(), // 1 hour from now
          status: 'active',
          bidCount: 5,
          product: DEFAULT_PRODUCTS.phone,
        },
      ];

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: mockAuctions,
          pagination: {
            page: 1,
            limit: 10,
            total: mockAuctions.length,
          },
        }),
      });
    });

    // Mock auction detail
    this.page.route('**/api/auctions/*', async (route: Route) => {
      const auctionId = route.request().url().split('/').pop();

      const mockAuction = {
        id: auctionId,
        productId: DEFAULT_PRODUCTS.watch.id,
        title: DEFAULT_PRODUCTS.watch.title + ' - Live Auction',
        currentBid: 550.00,
        endTime: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(),
        status: 'active',
        bidCount: 3,
        minBidIncrement: 5.00,
        product: DEFAULT_PRODUCTS.watch,
        bids: [
          {
            id: 'bid-1',
            amount: 500.00,
            timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
            user: {
              displayName: 'User1',
            },
          },
          {
            id: 'bid-2',
            amount: 525.00,
            timestamp: new Date(Date.now() - 20 * 60 * 1000).toISOString(),
            user: {
              displayName: 'User2',
            },
          },
          {
            id: 'bid-3',
            amount: 550.00,
            timestamp: new Date(Date.now() - 10 * 60 * 1000).toISOString(),
            user: {
              displayName: 'User3',
            },
          },
        ],
      };

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: mockAuction,
        }),
      });
    });

    // Mock place bid endpoint
    this.page.route('**/api/auctions/*/bid', async (route: Route) => {
      const request = route.request();
      const postData = request.postDataJSON();
      const auctionId = request.url().split('/')[4];

      // Validate bid amount
      if (postData.amount && postData.amount >= 100) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: {
              id: `bid-${Date.now()}`,
              auctionId,
              amount: postData.amount,
              timestamp: new Date().toISOString(),
              user: {
                displayName: 'Test User',
              },
            },
          }),
        });
      } else {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            error: 'Bid amount too low',
          }),
        });
      }
    });
  }

  /**
   * Mock product endpoints
   */
  mockProductEndpoints(): void {
    // Mock products list
    this.page.route(API_ENDPOINTS.products.list, async (route: Route) => {
      const mockProducts = Object.values(DEFAULT_PRODUCTS);

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: mockProducts,
          pagination: {
            page: 1,
            limit: 10,
            total: mockProducts.length,
          },
        }),
      });
    });

    // Mock product detail
    this.page.route('**/api/products/*', async (route: Route) => {
      const productId = route.request().url().split('/').pop();
      const product = Object.values(DEFAULT_PRODUCTS).find(p => p.id === productId);

      if (product) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: product,
          }),
        });
      } else {
        await route.fulfill({
          status: 404,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            error: 'Product not found',
          }),
        });
      }
    });
  }

  /**
   * Mock cart endpoints
   */
  mockCartEndpoints(): void {
    // Mock get cart
    this.page.route(API_ENDPOINTS.cart.items, async (route: Route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: {
            items: [],
            total: 0,
            itemCount: 0,
          },
        }),
      });
    });

    // Mock add to cart
    this.page.route(API_ENDPOINTS.cart.add, async (route: Route) => {
      const request = route.request();
      const postData = request.postDataJSON();

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: {
            id: `cart-item-${Date.now()}`,
            productId: postData.productId,
            quantity: postData.quantity || 1,
            addedAt: new Date().toISOString(),
          },
        }),
      });
    });
  }

  /**
   * Mock WebSocket connections for real-time updates
   */
  mockWebSocket(): void {
    // Mock WebSocket connection for real-time bid updates
    this.page.route('**/ws/**', async (route: Route) => {
      // In a real scenario, you'd set up a WebSocket server
      // For now, we'll just prevent the connection from failing
      await route.abort();
    });
  }

  /**
   * Mock API errors
   */
  mockApiErrors(): void {
    // Mock network errors
    this.page.route('**/api/auctions/*/bid', async (route: Route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({
          success: false,
          error: 'Internal server error',
        }),
      });
    });
  }

  /**
   * Mock slow API responses
   */
  mockSlowResponses(delay: number = 2000): void {
    this.page.route('**/api/**', async (route: Route) => {
      await new Promise(resolve => setTimeout(resolve, delay));
      await route.continue();
    });
  }

  /**
   * Enable all mocks
   */
  enableAllMocks(): void {
    this.mockAuthEndpoints();
    this.mockAuctionEndpoints();
    this.mockProductEndpoints();
    this.mockCartEndpoints();
    this.mockWebSocket();
  }

  /**
   * Clear all mocks
   */
  clearMocks(): void {
    // This would be called when you want to stop mocking
    // Playwright automatically clears routes when the page is closed
  }
}
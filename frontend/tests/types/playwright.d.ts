// Global type definitions for Playwright tests
import { Page, BrowserContext, Browser } from '@playwright/test';

// Extend Playwright types with custom fixtures
declare global {
  namespace PlaywrightTest {
    interface Fixtures {
      authenticatedPage: any; // Will be properly typed by auth.fixture.ts
      testUser: any;
      adminUser: any;
      auctionPage: any;
      productPage: any;
      testAuction: any;
      testProduct: any;
    }
  }
}

// Custom test data types
export interface TestEnvironment {
  baseURL: string;
  apiURL: string;
  wsURL: string;
  environment: 'development' | 'staging' | 'production';
}

export interface TestConfig {
  timeout: number;
  retries: number;
  headless: boolean;
  slowMo: number;
  screenshotOnFailure: boolean;
  videoOnFailure: boolean;
}

// Test result types
export interface TestResult {
  status: 'passed' | 'failed' | 'skipped' | 'timedOut';
  duration: number;
  errors?: string[];
  screenshots?: string[];
  video?: string;
  trace?: string;
}

// Mock data types
export interface MockUser {
  id: string;
  email: string;
  password: string;
  displayName: string;
  role: 'user' | 'admin';
  token?: string;
}

export interface MockAuction {
  id: string;
  productId: string;
  title: string;
  description: string;
  startingPrice: number;
  currentBid: number;
  endTime: string;
  status: 'active' | 'ended' | 'scheduled';
  bids: MockBid[];
}

export interface MockBid {
  id: string;
  auctionId: string;
  userId: string;
  amount: number;
  timestamp: string;
}

export interface MockProduct {
  id: string;
  title: string;
  description: string;
  price: number;
  images: string[];
  category: string;
  condition: 'new' | 'used' | 'refurbished';
}

// API response types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface AuthResponse {
  user: MockUser;
  token: string;
  refreshToken: string;
}

export interface AuctionResponse {
  auction: MockAuction;
  product: MockProduct;
  bids: MockBid[];
}

// Custom matchers
export interface CustomMatchers {
  toBeAccessible(): Promise<void>;
  toHaveValidAriaLabels(labels: Record<string, string>): Promise<void>;
  toHaveCorrectColorContrast(): Promise<void>;
  toBeKeyboardNavigable(): Promise<void>;
  toHaveValidFormStructure(): Promise<void>;
}

// Page object types
export interface PageObject {
  page: Page;
  goto(): Promise<void>;
  waitForLoad(): Promise<void>;
  isLoaded(): Promise<boolean>;
  takeScreenshot(name?: string): Promise<void>;
}

// Global test utilities
export interface TestHelpers {
  waitForAPI(endpoint: string): Promise<any>;
  mockAPI(endpoint: string, response: any): Promise<void>;
  clearStorage(): Promise<void>;
  setStorage(key: string, value: any): Promise<void>;
  getStorage(key: string): Promise<any>;
  generateTestData<T>(type: string): T;
  compareScreenshots(baseline: string, current: string): Promise<boolean>;
}

// Test context extensions
export interface ExtendedTestContext {
  testHelpers: TestHelpers;
  mockData: {
    users: MockUser[];
    auctions: MockAuction[];
    products: MockProduct[];
  };
  apiResponses: Map<string, any>;
  customMatchers: CustomMatchers;
}
// Export all fixtures for easy importing
export * from './auth.fixture';
export * from './auction.fixture';

// Re-export Playwright test and expect
import { test as baseTest, expect as baseExpect } from '@playwright/test';
import { test as authTest } from './auth.fixture';
import { test as auctionTest } from './auction.fixture';

// Default test with auth fixtures
export const test = authTest;
export const expect = baseExpect;

// Specialized test configurations
export { authTest as testWithAuth };
export { auctionTest as testWithAuctions };
export { baseTest as plainTest };
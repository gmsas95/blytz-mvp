import { test as base, expect } from '@playwright/test';
import { AuthFixtures } from './auth.fixture';
import { AuctionPage } from '../pages/auction.page';
import { ProductPage } from '../pages/product.page';
import { TestAuction, TestProduct } from '../utils/test-data';

// Define auction fixture types
export interface AuctionFixtures extends AuthFixtures {
  auctionPage: AuctionPage;
  productPage: ProductPage;
  testAuction: TestAuction;
  testProduct: TestProduct;
}

// Extend auth test with auction fixtures
export const test = base.extend<AuctionFixtures>({
  auctionPage: async ({ page }, use) => {
    const auctionPage = new AuctionPage(page);
    await use(auctionPage);
  },

  productPage: async ({ page }, use) => {
    const productPage = new ProductPage(page);
    await use(productPage);
  },

  testProduct: async ({}, use) => {
    const product: TestProduct = {
      id: 'test-product-1',
      title: 'Test Luxury Watch',
      description: 'A beautiful luxury watch for testing',
      startingPrice: 100.00,
      currentBid: 100.00,
      images: ['https://picsum.photos/400/300?random=1'],
      category: 'watches',
      condition: 'new',
    };
    await use(product);
  },

  testAuction: async ({ testProduct }, use) => {
    const auction: TestAuction = {
      id: 'test-auction-1',
      productId: testProduct.id,
      title: `${testProduct.title} - Live Auction`,
      startingPrice: testProduct.startingPrice,
      currentBid: testProduct.currentBid,
      endTime: new Date(Date.now() + 60 * 60 * 1000).toISOString(), // 1 hour from now
      status: 'active',
      bidCount: 0,
      minBidIncrement: 5.00,
    };
    await use(auction);
  },
});

export { expect } from '@playwright/test';
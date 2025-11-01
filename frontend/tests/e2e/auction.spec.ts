import { test, expect } from '../fixtures';
import { AuctionPage } from '../pages/auction.page';
import { ProductPage } from '../pages/product.page';
import { ApiMocks } from '../mocks/api-mocks';

test.describe('Auction Functionality', () => {
  let auctionPage: AuctionPage;
  let apiMocks: ApiMocks;

  test.use({ storageState: { cookies: [], origins: [] } }); // Start with clean storage

  test.beforeEach(async ({ page, authenticatedPage }) => {
    auctionPage = new AuctionPage(page);
    apiMocks = new ApiMocks(page);

    // Enable API mocks
    apiMocks.enableAllMocks();

    // Start the app and login
    await authenticatedPage.goto();
    await authenticatedPage.verifyUserIsLoggedIn();
  });

  test('should display auction list', async ({ page }) => {
    await auctionPage.goto();
    await auctionPage.waitForAuctionsToLoad();

    const auctionCount = await auctionPage.getActiveAuctionCount();
    expect(auctionCount).toBeGreaterThan(0);
  });

  test('should navigate to auction details', async ({ page }) => {
    await auctionPage.goto();
    await auctionPage.waitForAuctionsToLoad();

    const auctions = await auctionPage.getAuctionList();
    if (auctions.length > 0) {
      await auctionPage.clickAuction(auctions[0]);

      // Verify auction details page
      await expect(page.locator('[data-testid="auction-title"]')).toBeVisible();
      await expect(page.locator('[data-testid="current-bid"]')).toBeVisible();
      await expect(page.locator('[data-testid="bid-input"]')).toBeVisible();
      await expect(page.locator('[data-testid="place-bid-button"]')).toBeVisible();
    }
  });

  test('should place a valid bid', async ({ page }) => {
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    const bidAmount = 600.00;
    await auctionPage.placeBid(bidAmount);
    await auctionPage.verifyBidSuccess(bidAmount);
  });

  test('should reject invalid bid amount', async ({ page }) => {
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    // Try to place a bid that's too low
    const invalidBidAmount = 50.00;
    await auctionPage.placeBid(invalidBidAmount);
    await auctionPage.verifyBidError('Bid amount too low');
  });

  test('should display bid history', async ({ page }) => {
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    await expect(page.locator('[data-testid="bid-history"]')).toBeVisible();

    const bidItems = page.locator('[data-testid="bid-item"]');
    const bidCount = await bidItems.count();
    expect(bidCount).toBeGreaterThan(0);
  });

  test('should watch and unwatch auctions', async ({ page }) => {
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    // Watch auction
    await auctionPage.watchAuction();
    await auctionPage.verifyAuctionIsWatched();

    // Unwatch auction
    await auctionPage.unwatchAuction();
    await auctionPage.verifyAuctionIsNotWatched();
  });

  test('should display countdown timer', async ({ page }) => {
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    await auctionPage.verifyCountdown(true);
  });

  test('should handle real-time bid updates', async ({ page }) => {
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    const initialBid = 550.00;
    const newBid = 600.00;

    // Simulate real-time bid update
    await auctionPage.simulateRealtimeBidUpdate(newBid);

    // Verify UI updates
    await auctionPage.waitForBidUpdate();
    const currentBidElement = page.locator('[data-testid="current-bid"]');
    await expect(currentBidElement).toContainText(newBid.toString());
  });

  test('should disable bidding when auction ends', async ({ page }) => {
    // Mock ended auction
    await apiMocks.mockApiErrors();
    await auctionPage.goto('auction-1');
    await auctionPage.waitForAuctionsToLoad();

    await auctionPage.verifyAuctionEnd();
    await auctionPage.verifyBidButtonState('disabled');
  });

  test('should search auctions', async ({ page }) => {
    await auctionPage.goto();
    await auctionPage.waitForAuctionsToLoad();

    await auctionPage.searchAuctions('watch');
    await auctionPage.waitForAuctionsToLoad();

    // Verify search results
    const auctions = await auctionPage.getAuctionList();
    const hasWatchAuctions = auctions.some(auction =>
      auction.toLowerCase().includes('watch')
    );
    expect(hasWatchAuctions).toBe(true);
  });

  test('should filter auctions by category', async ({ page }) => {
    await auctionPage.goto();
    await auctionPage.waitForAuctionsToLoad();

    await auctionPage.filterByCategory('electronics');
    await auctionPage.waitForAuctionsToLoad();

    // Verify filtered results
    const auctions = await auctionPage.getAuctionList();
    expect(auctions.length).toBeGreaterThan(0);
  });

  test('should sort auctions', async ({ page }) => {
    await auctionPage.goto();
    await auctionPage.waitForAuctionsToLoad();

    await auctionPage.sortBy('Price: Low to High');
    await auctionPage.waitForAuctionsToLoad();

    // Verify sorting worked
    const auctions = await auctionPage.getAuctionList();
    expect(auctions.length).toBeGreaterThan(0);
  });

  test('should navigate from product to auction', async ({ page }) => {
    const productPage = new ProductPage(page);

    await productPage.gotoProduct('product-watch-1');
    await productPage.waitForProductsToLoad();

    // Click view auction button if available
    await productPage.viewAuction();

    // Should be on auction page
    await expect(page).toHaveURL(/.*auctions\/.*/);
    await auctionPage.waitForAuctionsToLoad();
  });
});

test.describe('Auction Mobile Experience', () => {
  test.use({
    viewport: { width: 375, height: 667 }, // Mobile viewport
    storageState: { cookies: [], origins: [] }
  });

  test('should work on mobile devices', async ({ page, authenticatedPage }) => {
    const auctionPage = new AuctionPage(page);
    const apiMocks = new ApiMocks(page);
    apiMocks.enableAllMocks();

    await authenticatedPage.goto();
    await authenticatedPage.verifyUserIsLoggedIn();

    await auctionPage.goto();
    await auctionPage.waitForAuctionsToLoad();

    const auctionCount = await auctionPage.getActiveAuctionCount();
    expect(auctionCount).toBeGreaterThan(0);

    // Test mobile-specific interactions
    if (auctionCount > 0) {
      const auctions = await auctionPage.getAuctionList();
      await auctionPage.clickAuction(auctions[0]);

      // Verify mobile layout
      await expect(page.locator('[data-testid="auction-title"]')).toBeVisible();
      await expect(page.locator('[data-testid="mobile-bid-section"]')).toBeVisible();
    }
  });
});
import { Page, expect } from '@playwright/test';
import { TestHelpers } from '../utils/test-helpers';
import { TestAuction, TestBid } from '../utils/test-data';

export class AuctionPage {
  private helpers: TestHelpers;

  constructor(private page: Page) {
    this.helpers = new TestHelpers(page);
  }

  // Selectors
  private readonly selectors = {
    auctionList: '[data-testid="auction-list"], .auction-list',
    auctionCard: '[data-testid="auction-card"], .auction-card',
    auctionItem: '[data-testid="auction-item"], .auction-item',
    auctionTitle: '[data-testid="auction-title"], .auction-title',
    currentBid: '[data-testid="current-bid"], .current-bid',
    bidInput: '[data-testid="bid-input"], input[type="number"][placeholder*="bid"]',
    placeBidButton: '[data-testid="place-bid-button"], button:has-text("Place Bid")',
    bidHistory: '[data-testid="bid-history"], .bid-history',
    bidHistoryItem: '[data-testid="bid-item"], .bid-item',
    countdown: '[data-testid="countdown"], .countdown, .timer',
    watchButton: '[data-testid="watch-button"], button:has-text("Watch")',
    unwatchButton: '[data-testid="unwatch-button"], button:has-text("Unwatch")',
    shareButton: '[data-testid="share-button"], button:has-text("Share")',
    productImage: '[data-testid="product-image"], .product-image img',
    productGallery: '[data-testid="product-gallery"], .product-gallery',
    bidConfirmation: '[data-testid="bid-confirmation"], .bid-confirmation',
    errorMessage: '[data-testid="error-message"], .error-message',
    successMessage: '[data-testid="success-message"], .success-message',
    loading: '[data-testid="loading"], .loading, .spinner',
    noAuctions: '[data-testid="no-auctions"], .no-auctions',
    auctionStatus: '[data-testid="auction-status"], .auction-status',
    minBidInfo: '[data-testid="min-bid-info"], .min-bid-info',
    bidIncrement: '[data-testid="bid-increment"], .bid-increment',
  };

  async goto(auctionId?: string): Promise<void> {
    const url = auctionId ? `/auctions/${auctionId}` : '/auctions';
    await this.page.goto(url);
    await this.helpers.waitForPageLoad();
  }

  async gotoLive(): Promise<void> {
    await this.page.goto('/auctions/live');
    await this.helpers.waitForPageLoad();
  }

  async gotoWatchlist(): Promise<void> {
    await this.page.goto('/auctions/watchlist');
    await this.helpers.waitForPageLoad();
  }

  async waitForAuctionsToLoad(): Promise<void> {
    await this.helpers.waitForLoadingComplete();

    // Wait for auction list or no auctions message
    await this.page.waitForSelector(
      `${this.selectors.auctionList}, ${this.selectors.noAuctions}`,
      { state: 'visible' }
    );
  }

  async clickAuction(auctionTitle: string): Promise<void> {
    const auctionCard = this.page.locator(this.selectors.auctionCard).filter({
      hasText: auctionTitle,
    });

    await expect(auctionCard).toBeVisible();
    await auctionCard.click();
    await this.helpers.waitForLoadingComplete();
  }

  async clickAuctionById(auctionId: string): Promise<void> {
    const auctionItem = this.page.locator(`[data-auction-id="${auctionId}"]`);
    await expect(auctionItem).toBeVisible();
    await auctionItem.click();
    await this.helpers.waitForLoadingComplete();
  }

  async placeBid(amount: number): Promise<void> {
    const bidInput = this.page.locator(this.selectors.bidInput);
    await bidInput.fill(amount.toString());

    await this.page.click(this.selectors.placeBidButton);
    await this.helpers.waitForLoadingComplete();
  }

  async verifyBidSuccess(expectedAmount: number): Promise<void> {
    const successMessage = this.page.locator(this.selectors.successMessage);
    const bidConfirmation = this.page.locator(this.selectors.bidConfirmation);

    await expect(
      successMessage.or(bidConfirmation)
    ).toBeVisible({ timeout: 10000 });

    // Verify current bid is updated
    const currentBid = this.page.locator(this.selectors.currentBid);
    await expect(currentBid).toContainText(expectedAmount.toString());
  }

  async verifyBidError(expectedMessage?: string): Promise<void> {
    const errorMessage = this.page.locator(this.selectors.errorMessage);
    await expect(errorMessage).toBeVisible({ timeout: 5000 });

    if (expectedMessage) {
      await expect(errorMessage).toContainText(expectedMessage);
    }
  }

  async verifyAuctionDetails(auction: TestAuction): Promise<void> {
    const title = this.page.locator(this.selectors.auctionTitle);
    await expect(title).toContainText(auction.title);

    const currentBid = this.page.locator(this.selectors.currentBid);
    await expect(currentBid).toContainText(auction.currentBid.toString());

    const status = this.page.locator(this.selectors.auctionStatus);
    if (await status.isVisible()) {
      await expect(status).toContainText(auction.status);
    }
  }

  async verifyBidHistory(bids: TestBid[]): Promise<void> {
    const bidHistory = this.page.locator(this.selectors.bidHistory);
    await expect(bidHistory).toBeVisible();

    const bidItems = this.page.locator(this.selectors.bidHistoryItem);
    const bidCount = await bidItems.count();

    expect(bidCount).toBeGreaterThanOrEqual(bids.length);

    // Verify most recent bids are displayed
    for (let i = 0; i < Math.min(bids.length, bidCount); i++) {
      const bidItem = bidItems.nth(i);
      await expect(bidItem).toContainText(bids[bids.length - 1 - i].amount.toString());
    }
  }

  async watchAuction(): Promise<void> {
    const watchButton = this.page.locator(this.selectors.watchButton);
    if (await watchButton.isVisible()) {
      await watchButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async unwatchAuction(): Promise<void> {
    const unwatchButton = this.page.locator(this.selectors.unwatchButton);
    if (await unwatchButton.isVisible()) {
      await unwatchButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async verifyAuctionIsWatched(): Promise<void> {
    const unwatchButton = this.page.locator(this.selectors.unwatchButton);
    await expect(unwatchButton).toBeVisible();
  }

  async verifyAuctionIsNotWatched(): Promise<void> {
    const watchButton = this.page.locator(this.selectors.watchButton);
    await expect(watchButton).toBeVisible();
  }

  async verifyCountdown(isActive: boolean = true): Promise<void> {
    const countdown = this.page.locator(this.selectors.countdown);
    if (isActive) {
      await expect(countdown).toBeVisible();
      // Verify countdown is actually counting down
      const initialTime = await countdown.textContent();
      await this.page.waitForTimeout(2000);
      const laterTime = await countdown.textContent();
      expect(initialTime).not.toBe(laterTime);
    } else {
      // For ended auctions, show "Auction Ended" message
      await expect(countdown).toContainText(/ended|finished|closed/i);
    }
  }

  async verifyMinimumBidIncrement(minIncrement: number): Promise<void> {
    const incrementInfo = this.page.locator(this.selectors.bidIncrement);
    if (await incrementInfo.isVisible()) {
      await expect(incrementInfo).toContainText(minIncrement.toString());
    }
  }

  async verifyProductImages(): Promise<void> {
    const productImage = this.page.locator(this.selectors.productImage);
    await expect(productImage).toBeVisible();

    // Verify image has loaded
    const hasLoaded = await productImage.evaluate(img =>
      (img as HTMLImageElement).complete &&
      (img as HTMLImageElement).naturalHeight !== 0
    );
    expect(hasLoaded).toBe(true);
  }

  async shareAuction(): Promise<void> {
    const shareButton = this.page.locator(this.selectors.shareButton);
    if (await shareButton.isVisible()) {
      await shareButton.click();
      // Handle share dialog or clipboard interaction
      await this.page.waitForTimeout(1000);
    }
  }

  async verifyLiveAuctionFeatures(): Promise<void> {
    // Verify live auction specific features
    await expect(this.page.locator(this.selectors.countdown)).toBeVisible();
    await expect(this.page.locator(this.selectors.bidInput)).toBeVisible();
    await expect(this.page.locator(this.selectors.placeBidButton)).toBeVisible();
    await expect(this.page.locator(this.selectors.bidHistory)).toBeVisible();
  }

  async getAuctionList(): Promise<string[]> {
    const auctionCards = this.page.locator(this.selectors.auctionCard);
    const titles = await auctionCards.locator(this.selectors.auctionTitle).allTextContents();
    return titles;
  }

  async getActiveAuctionCount(): Promise<number> {
    const activeAuctions = this.page.locator(
      `${this.selectors.auctionCard}[data-status="active"]`
    );
    return await activeAuctions.count();
  }

  async searchAuctions(searchTerm: string): Promise<void> {
    const searchInput = this.page.locator('input[placeholder*="search"], [data-testid="search-input"]');
    if (await searchInput.isVisible()) {
      await searchInput.fill(searchTerm);
      await searchInput.press('Enter');
      await this.helpers.waitForLoadingComplete();
    }
  }

  async filterByCategory(category: string): Promise<void> {
    const categoryFilter = this.page.locator(
      `[data-testid="category-filter"], select[name="category"]`
    );
    if (await categoryFilter.isVisible()) {
      await categoryFilter.selectOption({ label: category });
      await this.helpers.waitForLoadingComplete();
    }
  }

  async sortBy(sortOption: string): Promise<void> {
    const sortSelect = this.page.locator(
      `[data-testid="sort-select"], select[name="sort"]`
    );
    if (await sortSelect.isVisible()) {
      await sortSelect.selectOption({ label: sortOption });
      await this.helpers.waitForLoadingComplete();
    }
  }

  async verifyBidButtonState(expectedState: 'enabled' | 'disabled'): Promise<void> {
    const bidButton = this.page.locator(this.selectors.placeBidButton);

    if (expectedState === 'enabled') {
      await expect(bidButton).toBeEnabled();
    } else {
      await expect(bidButton).toBeDisabled();
    }
  }

  async verifyAuctionEnd(): Promise<void> {
    const status = this.page.locator(this.selectors.auctionStatus);
    await expect(status).toContainText(/ended|finished|closed/i);

    const bidButton = this.page.locator(this.selectors.placeBidButton);
    await expect(bidButton).toBeDisabled();
  }

  async waitForBidUpdate(): Promise<void> {
    // Wait for bid to be reflected in UI
    await this.page.waitForFunction(() => {
      const currentBid = document.querySelector('[data-testid="current-bid"]');
      return currentBid && currentBid.textContent !== '';
    }, { timeout: 10000 });
  }

  async simulateRealtimeBidUpdate(newAmount: number): Promise<void> {
    // Mock real-time bid update via WebSocket or polling
    await this.page.evaluate((amount) => {
      // Simulate real-time bid update
      const currentBidElement = document.querySelector('[data-testid="current-bid"]');
      if (currentBidElement) {
        currentBidElement.textContent = `$${amount.toFixed(2)}`;
      }

      // Add to bid history
      const bidHistory = document.querySelector('[data-testid="bid-history"]');
      if (bidHistory) {
        const newBidItem = document.createElement('div');
        newBidItem.setAttribute('data-testid', 'bid-item');
        newBidItem.textContent = `New bid: $${amount.toFixed(2)}`;
        bidHistory.insertBefore(newBidItem, bidHistory.firstChild);
      }
    }, newAmount);
  }
}
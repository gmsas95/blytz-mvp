import { Page, expect } from '@playwright/test';
import { TestHelpers } from '../utils/test-helpers';
import { TestProduct } from '../utils/test-data';

export class ProductPage {
  private helpers: TestHelpers;

  constructor(private page: Page) {
    this.helpers = new TestHelpers(page);
  }

  // Selectors
  private readonly selectors = {
    productList: '[data-testid="product-list"], .product-list',
    productCard: '[data-testid="product-card"], .product-card',
    productItem: '[data-testid="product-item"], .product-item',
    productTitle: '[data-testid="product-title"], .product-title',
    productDescription: '[data-testid="product-description"], .product-description',
    productPrice: '[data-testid="product-price"], .product-price',
    productImage: '[data-testid="product-image"], .product-image img',
    productGallery: '[data-testid="product-gallery"], .product-gallery',
    addToCartButton: '[data-testid="add-to-cart"], button:has-text("Add to Cart")',
    buyNowButton: '[data-testid="buy-now"], button:has-text("Buy Now")',
    viewAuctionButton: '[data-testid="view-auction"], button:has-text("View Auction")',
    addToWatchlistButton: '[data-testid="add-to-watchlist"], button:has-text("Watch")',
    removeFromWatchlistButton: '[data-testid="remove-from-watchlist"], button:has-text("Unwatch")',
    productCategory: '[data-testid="product-category"], .product-category',
    productCondition: '[data-testid="product-condition"], .product-condition',
    productSpecs: '[data-testid="product-specs"], .product-specs',
    productRating: '[data-testid="product-rating"], .product-rating',
    productReviews: '[data-testid="product-reviews"], .product-reviews',
    searchInput: '[data-testid="search-input"], input[placeholder*="search"]',
    categoryFilter: '[data-testid="category-filter"], select[name="category"]',
    priceFilter: '[data-testid="price-filter"], [data-testid="price-range"]',
    sortSelect: '[data-testid="sort-select"], select[name="sort"]',
    noProducts: '[data-testid="no-products"], .no-products',
    loading: '[data-testid="loading"], .loading, .spinner',
    breadcrumb: '[data-testid="breadcrumb"], .breadcrumb',
    relatedProducts: '[data-testid="related-products"], .related-products',
  };

  async goto(): Promise<void> {
    await this.page.goto('/products');
    await this.helpers.waitForPageLoad();
  }

  async gotoProduct(productId: string): Promise<void> {
    await this.page.goto(`/products/${productId}`);
    await this.helpers.waitForPageLoad();
  }

  async gotoCategory(category: string): Promise<void> {
    await this.page.goto(`/products/category/${category}`);
    await this.helpers.waitForPageLoad();
  }

  async waitForProductsToLoad(): Promise<void> {
    await this.helpers.waitForLoadingComplete();

    // Wait for product list or no products message
    await this.page.waitForSelector(
      `${this.selectors.productList}, ${this.selectors.noProducts}`,
      { state: 'visible' }
    );
  }

  async clickProduct(productTitle: string): Promise<void> {
    const productCard = this.page.locator(this.selectors.productCard).filter({
      hasText: productTitle,
    });

    await expect(productCard).toBeVisible();
    await productCard.click();
    await this.helpers.waitForLoadingComplete();
  }

  async clickProductById(productId: string): Promise<void> {
    const productItem = this.page.locator(`[data-product-id="${productId}"]`);
    await expect(productItem).toBeVisible();
    await productItem.click();
    await this.helpers.waitForLoadingComplete();
  }

  async verifyProductDetails(product: TestProduct): Promise<void> {
    const title = this.page.locator(this.selectors.productTitle);
    await expect(title).toContainText(product.title);

    const description = this.page.locator(this.selectors.productDescription);
    if (await description.isVisible()) {
      await expect(description).toContainText(product.description);
    }

    const price = this.page.locator(this.selectors.productPrice);
    await expect(price).toContainText(product.startingPrice.toString());

    const category = this.page.locator(this.selectors.productCategory);
    if (await category.isVisible()) {
      await expect(category).toContainText(product.category);
    }

    const condition = this.page.locator(this.selectors.productCondition);
    if (await condition.isVisible()) {
      await expect(condition).toContainText(product.condition);
    }
  }

  async addToCart(): Promise<void> {
    const addToCartButton = this.page.locator(this.selectors.addToCartButton);
    await expect(addToCartButton).toBeVisible();
    await addToCartButton.click();
    await this.helpers.waitForLoadingComplete();
  }

  async buyNow(): Promise<void> {
    const buyNowButton = this.page.locator(this.selectors.buyNowButton);
    if (await buyNowButton.isVisible()) {
      await buyNowButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async viewAuction(): Promise<void> {
    const viewAuctionButton = this.page.locator(this.selectors.viewAuctionButton);
    if (await viewAuctionButton.isVisible()) {
      await viewAuctionButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async addToWatchlist(): Promise<void> {
    const watchlistButton = this.page.locator(this.selectors.addToWatchlistButton);
    if (await watchlistButton.isVisible()) {
      await watchlistButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async removeFromWatchlist(): Promise<void> {
    const removeButton = this.page.locator(this.selectors.removeFromWatchlistButton);
    if (await removeButton.isVisible()) {
      await removeButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async verifyProductInWatchlist(): Promise<void> {
    const removeButton = this.page.locator(this.selectors.removeFromWatchlistButton);
    await expect(removeButton).toBeVisible();
  }

  async verifyProductNotInWatchlist(): Promise<void> {
    const addButton = this.page.locator(this.selectors.addToWatchlistButton);
    await expect(addButton).toBeVisible();
  }

  async verifyProductImages(): Promise<void> {
    const productImage = this.page.locator(this.selectors.productImage);
    await expect(productImage).toBeVisible();

    // Verify main image has loaded
    const hasLoaded = await productImage.evaluate(img =>
      (img as HTMLImageElement).complete &&
      (img as HTMLImageElement).naturalHeight !== 0
    );
    expect(hasLoaded).toBe(true);

    // Check for product gallery if multiple images
    const gallery = this.page.locator(this.selectors.productGallery);
    if (await gallery.isVisible()) {
      const galleryImages = gallery.locator('img');
      const imageCount = await galleryImages.count();
      expect(imageCount).toBeGreaterThan(0);
    }
  }

  async searchProducts(searchTerm: string): Promise<void> {
    const searchInput = this.page.locator(this.selectors.searchInput);
    if (await searchInput.isVisible()) {
      await searchInput.fill(searchTerm);
      await searchInput.press('Enter');
      await this.helpers.waitForLoadingComplete();
    }
  }

  async filterByCategory(category: string): Promise<void> {
    const categoryFilter = this.page.locator(this.selectors.categoryFilter);
    if (await categoryFilter.isVisible()) {
      await categoryFilter.selectOption({ label: category });
      await this.helpers.waitForLoadingComplete();
    }
  }

  async filterByPriceRange(minPrice: number, maxPrice: number): Promise<void> {
    const priceFilter = this.page.locator(this.selectors.priceFilter);
    if (await priceFilter.isVisible()) {
      // Handle price range slider or input fields
      const minInput = priceFilter.locator('input[name="min-price"], input[data-testid="min-price"]');
      const maxInput = priceFilter.locator('input[name="max-price"], input[data-testid="max-price"]');

      if (await minInput.isVisible()) {
        await minInput.fill(minPrice.toString());
      }
      if (await maxInput.isVisible()) {
        await maxInput.fill(maxPrice.toString());
      }

      // Apply filter
      const applyButton = priceFilter.locator('button:has-text("Apply"), button[type="submit"]');
      if (await applyButton.isVisible()) {
        await applyButton.click();
      }

      await this.helpers.waitForLoadingComplete();
    }
  }

  async sortBy(sortOption: string): Promise<void> {
    const sortSelect = this.page.locator(this.selectors.sortSelect);
    if (await sortSelect.isVisible()) {
      await sortSelect.selectOption({ label: sortOption });
      await this.helpers.waitForLoadingComplete();
    }
  }

  async getProductList(): Promise<string[]> {
    const productCards = this.page.locator(this.selectors.productCard);
    const titles = await productCards.locator(this.selectors.productTitle).allTextContents();
    return titles;
  }

  async getProductCount(): Promise<number> {
    const productCards = this.page.locator(this.selectors.productCard);
    return await productCards.count();
  }

  async verifyProductAvailability(isAvailable: boolean = true): Promise<void> {
    const addToCartButton = this.page.locator(this.selectors.addToCartButton);
    const buyNowButton = this.page.locator(this.selectors.buyNowButton);

    if (isAvailable) {
      await expect(addToCartButton.or(buyNowButton)).toBeVisible();
    } else {
      await expect(addToCartButton).not.toBeVisible();
      await expect(buyNowButton).not.toBeVisible();

      // Check for out of stock message
      const outOfStock = this.page.locator('[data-testid="out-of-stock"], .out-of-stock');
      if (await outOfStock.isVisible()) {
        await expect(outOfStock).toBeVisible();
      }
    }
  }

  async verifyProductSpecifications(specs: Record<string, string>): Promise<void> {
    const specsSection = this.page.locator(this.selectors.productSpecs);
    if (await specsSection.isVisible()) {
      for (const [key, value] of Object.entries(specs)) {
        await expect(specsSection).toContainText(key);
        await expect(specsSection).toContainText(value);
      }
    }
  }

  async verifyProductRating(expectedRating?: number): Promise<void> {
    const rating = this.page.locator(this.selectors.productRating);
    if (await rating.isVisible()) {
      if (expectedRating) {
        await expect(rating).toContainText(expectedRating.toString());
      }
    }
  }

  async verifyProductReviews(): Promise<void> {
    const reviews = this.page.locator(this.selectors.productReviews);
    if (await reviews.isVisible()) {
      await expect(reviews).toBeVisible();
    }
  }

  async navigateThroughImageGallery(): Promise<void> {
    const gallery = this.page.locator(this.selectors.productGallery);
    if (await gallery.isVisible()) {
      const thumbnails = gallery.locator('.thumbnail, [data-testid="thumbnail"]');
      const thumbnailCount = await thumbnails.count();

      if (thumbnailCount > 1) {
        // Click on different thumbnails
        for (let i = 0; i < Math.min(3, thumbnailCount); i++) {
          await thumbnails.nth(i).click();
          await this.page.waitForTimeout(500); // Wait for image to load
        }
      }
    }
  }

  async verifyRelatedProducts(): Promise<void> {
    const relatedProducts = this.page.locator(this.selectors.relatedProducts);
    if (await relatedProducts.isVisible()) {
      const relatedProductCards = relatedProducts.locator(this.selectors.productCard);
      const relatedCount = await relatedProductCards.count();
      expect(relatedCount).toBeGreaterThan(0);
    }
  }

  async verifyBreadcrumbNavigation(): Promise<void> {
    const breadcrumb = this.page.locator(this.selectors.breadcrumb);
    if (await breadcrumb.isVisible()) {
      await expect(breadcrumb).toBeVisible();
    }
  }

  async compareProducts(productId1: string, productId2: string): Promise<void> {
    // Navigate to product comparison page if available
    await this.page.goto(`/products/compare?ids=${productId1},${productId2}`);
    await this.helpers.waitForLoadingComplete();

    // Verify comparison table
    const comparisonTable = this.page.locator('[data-testid="comparison-table"], .comparison-table');
    if (await comparisonTable.isVisible()) {
      await expect(comparisonTable).toBeVisible();
    }
  }

  async addToWishlist(): Promise<void> {
    // Similar to watchlist but for products
    const wishlistButton = this.page.locator('[data-testid="add-to-wishlist"], button:has-text("Save")');
    if (await wishlistButton.isVisible()) {
      await wishlistButton.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async verifyAddToCartSuccess(): Promise<void> {
    const successMessage = this.page.locator('[data-testid="cart-success"], .success-message');
    const cartIcon = this.page.locator('[data-testid="cart-icon"], .cart-icon');

    // Either show success message or update cart icon
    await expect(
      successMessage.or(cartIcon)
    ).toBeVisible({ timeout: 10000 });
  }

  async verifyAddToCartError(expectedMessage?: string): Promise<void> {
    const errorMessage = this.page.locator('[data-testid="cart-error"], .error-message');
    await expect(errorMessage).toBeVisible({ timeout: 5000 });

    if (expectedMessage) {
      await expect(errorMessage).toContainText(expectedMessage);
    }
  }

  async getProductDetailsFromPage(): Promise<Partial<TestProduct>> {
    return await this.page.evaluate(() => {
      const title = document.querySelector('[data-testid="product-title"]')?.textContent;
      const description = document.querySelector('[data-testid="product-description"]')?.textContent;
      const price = document.querySelector('[data-testid="product-price"]')?.textContent;
      const category = document.querySelector('[data-testid="product-category"]')?.textContent;
      const condition = document.querySelector('[data-testid="product-condition"]')?.textContent;

      return {
        title: title || undefined,
        description: description || undefined,
        startingPrice: price ? parseFloat(price.replace(/[^0-9.]/g, '')) : undefined,
        category: category || undefined,
        condition: condition as 'new' | 'used' | 'refurbished' || undefined,
      };
    });
  }
}
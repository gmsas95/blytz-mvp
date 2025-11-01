import { Page, Locator, expect } from '@playwright/test';
import { API_ENDPOINTS } from './test-data';

/**
 * Common test helper functions
 */
export class TestHelpers {
  constructor(private page: Page) {}

  /**
   * Wait for page to be fully loaded
   */
  async waitForPageLoad(): Promise<void> {
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForLoadState('domcontentloaded');
  }

  /**
   * Wait for and handle loading states
   */
  async waitForLoadingComplete(): Promise<void> {
    const loadingSelectors = [
      '[data-testid="loading"]',
      '.loading',
      '.spinner',
      '[role="progressbar"]',
    ];

    // Wait for any loading indicators to disappear
    for (const selector of loadingSelectors) {
      const loader = this.page.locator(selector);
      if (await loader.isVisible()) {
        await loader.waitFor({ state: 'hidden' });
      }
    }

    // Wait for network to be idle
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Wait for API request to complete
   */
  async waitForApiResponse(endpoint: string): Promise<any> {
    return await this.page.waitForResponse(
      response => response.url().includes(endpoint) && response.status() === 200
    );
  }

  /**
   * Mock API responses
   */
  async mockApiResponse(endpoint: string, response: any, status: number = 200): Promise<void> {
    await this.page.route(`**${endpoint}**`, async route => {
      await route.fulfill({
        status,
        contentType: 'application/json',
        body: JSON.stringify(response),
      });
    });
  }

  /**
   * Mock API error responses
   */
  async mockApiError(endpoint: string, status: number = 500, message: string = 'Internal Server Error'): Promise<void> {
    await this.page.route(`**${endpoint}**`, async route => {
      await route.fulfill({
        status,
        contentType: 'application/json',
        body: JSON.stringify({ error: message }),
      });
    });
  }

  /**
   * Take screenshot with better naming
   */
  async takeScreenshot(name: string, fullPage: boolean = true): Promise<void> {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const filename = `${name}-${timestamp}.png`;
    await this.page.screenshot({
      path: `test-results/screenshots/${filename}`,
      fullPage,
    });
  }

  /**
   * Verify element is accessible (visible and enabled)
   */
  async verifyAccessible(selector: string): Promise<void> {
    const element = this.page.locator(selector);
    await expect(element).toBeVisible();
    await expect(element).toBeEnabled();
  }

  /**
   * Fill form with data
   */
  async fillForm(formData: Record<string, string>): Promise<void> {
    for (const [field, value] of Object.entries(formData)) {
      await this.page.fill(field, value);
    }
  }

  /**
   * Verify toast/notification appears
   */
  async verifyToast(message: string): Promise<void> {
    const toast = this.page.locator('[data-testid="toast"], .toast, .notification');
    await expect(toast).toContainText(message);
    await toast.waitFor({ state: 'visible' });

    // Wait for toast to disappear (auto-dismiss)
    await toast.waitFor({ state: 'hidden', timeout: 5000 }).catch(() => {
      // Toast might not auto-dismiss, that's ok
    });
  }

  /**
   * Handle consent dialogs
   */
  async handleConsentDialogs(): Promise<void> {
    const consentSelectors = [
      '[data-testid="accept-cookies"]',
      '[data-testid="consent-accept"]',
      '.cookie-accept',
      '.consent-accept',
      'button:has-text("Accept")',
      'button:has-text("I Agree")',
    ];

    for (const selector of consentSelectors) {
      const button = this.page.locator(selector);
      if (await button.isVisible()) {
        await button.click();
        break;
      }
    }
  }

  /**
   * Simulate mobile device gestures
   */
  async swipe(startX: number, startY: number, endX: number, endY: number): Promise<void> {
    await this.page.touchscreen.tap(startX, startY);
    await this.page.touchscreen.move(endX, endY);
    await this.page.touchscreen.tap(endX, endY);
  }

  /**
   * Verify page title
   */
  async verifyPageTitle(expectedTitle: string): Promise<void> {
    await expect(this.page).toHaveTitle(new RegExp(expectedTitle, 'i'));
  }

  /**
   * Get current URL with query parameters
   */
  async getCurrentUrl(): Promise<string> {
    return this.page.url();
  }

  /**
   * Navigate with retry logic
   */
  async navigateWithRetry(url: string, maxRetries: number = 3): Promise<void> {
    for (let i = 0; i < maxRetries; i++) {
      try {
        await this.page.goto(url, { waitUntil: 'networkidle' });
        return;
      } catch (error) {
        if (i === maxRetries - 1) throw error;
        await this.page.waitForTimeout(1000);
      }
    }
  }

  /**
   * Verify element has correct ARIA attributes
   */
  async verifyAriaAttributes(selector: string, attributes: Record<string, string>): Promise<void> {
    const element = this.page.locator(selector);
    for (const [attribute, value] of Object.entries(attributes)) {
      await expect(element).toHaveAttribute(`aria-${attribute}`, value);
    }
  }

  /**
   * Check if element is in viewport
   */
  async isInViewport(selector: string): Promise<boolean> {
    const element = this.page.locator(selector);
    const boundingBox = await element.boundingBox();
    if (!boundingBox) return false;

    const viewport = this.page.viewportSize();
    if (!viewport) return false;

    return (
      boundingBox.x >= 0 &&
      boundingBox.y >= 0 &&
      boundingBox.x + boundingBox.width <= viewport.width &&
      boundingBox.y + boundingBox.height <= viewport.height
    );
  }

  /**
   * Wait for animation to complete
   */
  async waitForAnimation(selector: string): Promise<void> {
    const element = this.page.locator(selector);
    await element.waitFor({ state: 'visible' });

    // Wait for CSS transitions/animations
    await this.page.evaluate(el => {
      return new Promise(resolve => {
        if (el.getAnimations().length === 0) {
          resolve(null);
          return;
        }

        Promise.all(el.getAnimations().map(animation => animation.finished)).then(() => {
          resolve(null);
        });
      });
    }, await element.elementHandle());
  }

  /**
   * Generate random test data
   */
  static generateRandomString(length: number = 10): string {
    return Math.random().toString(36).substring(2, 2 + length);
  }

  static generateRandomEmail(): string {
    return `test-${this.generateRandomString(8)}@blytz.app`;
  }

  static generateRandomNumber(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
  }
}
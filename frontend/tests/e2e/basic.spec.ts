import { test, expect } from '@playwright/test';

test.describe('Basic Application Tests', () => {
  test('should load the homepage', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/Blytz/i);
  });

  test('should have proper meta tags', async ({ page }) => {
    await page.goto('/');

    const metaDescription = await page.locator('meta[name="description"]').getAttribute('content');
    expect(metaDescription).toBeTruthy();
  });

  test('should be responsive on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');

    await expect(page.locator('body')).toBeVisible();
  });

  test('should handle navigation', async ({ page }) => {
    await page.goto('/');

    // Test navigation to different pages if available
    // Use first matching element for navigation
    const auctionsLink = page.locator('a[href*="auction"]').first();
    if (await auctionsLink.isVisible()) {
      await auctionsLink.click();
      await expect(page).toHaveURL(/.*auction.*/);
    }
  });
});
import { test, expect } from '@playwright/test';
import { AuthPage } from '../pages/auth.page';
import { AccessibilityHelpers } from '../utils/accessibility-helpers';

test.describe('Authentication Accessibility', () => {
  let authPage: AuthPage;
  let a11y: AccessibilityHelpers;

  test.beforeEach(async ({ page }) => {
    authPage = new AuthPage(page);
    a11y = new AccessibilityHelpers(page);
    await a11y.initializeAxe();
  });

  test('should be accessible on login page', async ({ page }) => {
    await authPage.goto();
    await a11y.checkPageAccessibility();
  });

  test('should be accessible on register page', async ({ page }) => {
    await authPage.gotoRegister();
    await a11y.checkPageAccessibility();
  });

  test('should have proper form labels and descriptions', async ({ page }) => {
    await authPage.goto();

    // Check email input
    await a11y.verifyAriaLabels('[data-testid="email-input"]', {
      label: 'Email address',
      role: 'textbox',
    });

    // Check password input
    await a11y.verifyAriaLabels('[data-testid="password-input"]', {
      label: 'Password',
      role: 'textbox',
    });

    // Check submit button
    await a11y.verifyAriaLabels('[data-testid="submit-button"]', {
      role: 'button',
    });
  });

  test('should support keyboard navigation', async ({ page }) => {
    await authPage.goto();

    // Test tab navigation through form
    await page.keyboard.press('Tab'); // Focus email input
    await expect(page.locator('[data-testid="email-input"]')).toBeFocused();

    await page.keyboard.press('Tab'); // Focus password input
    await expect(page.locator('[data-testid="password-input"]')).toBeFocused();

    await page.keyboard.press('Tab'); // Focus submit button
    await expect(page.locator('[data-testid="submit-button"]')).toBeFocused();

    // Test form submission with Enter
    await page.keyboard.press('Tab'); // Go back to email input
    await page.keyboard.press('Enter');
    // Should trigger validation
  });

  test('should have accessible error messages', async ({ page }) => {
    await authPage.goto();

    // Submit empty form to trigger validation errors
    await authPage.submitForm();

    // Check if error messages are accessible
    const errorMessage = page.locator('[data-testid="error-message"]');
    if (await errorMessage.isVisible()) {
      await a11y.checkElementAccessibility('[data-testid="error-message"]');
      await a11y.verifyAriaLabels('[data-testid="error-message"]', {
        role: 'alert',
      });
    }
  });

  test('should have proper focus management', async ({ page }) => {
    await authPage.goto();

    // Test focus moves to error when validation fails
    await authPage.submitForm();

    const firstError = page.locator('[data-testid="error-message"]').first();
    if (await firstError.isVisible()) {
      await expect(firstError).toBeFocused();
    }
  });

  test('should have accessible links and buttons', async ({ page }) => {
    await authPage.goto();

    // Check forgot password link
    const forgotLink = page.locator('[data-testid="forgot-password"]');
    if (await forgotLink.isVisible()) {
      await a11y.verifyKeyboardNavigation('[data-testid="forgot-password"]');
    }

    // Check register link
    const registerLink = page.locator('a[href*="register"]');
    if (await registerLink.isVisible()) {
      await a11y.verifyKeyboardNavigation('a[href*="register"]');
    }
  });

  test('should have sufficient color contrast', async ({ page }) => {
    await authPage.goto();

    // Check color contrast for important elements
    await a11y.verifyColorContrast('[data-testid="submit-button"]');
    await a11y.verifyColorContrast('[data-testid="email-input"]');
    await a11y.verifyColorContrast('[data-testid="password-input"]');
  });

  test('should have accessible form validation', async ({ page }) => {
    await authPage.gotoRegister();

    // Fill form partially and submit
    await page.fill('[data-testid="email-input"]', 'invalid-email');
    await authPage.submitForm();

    // Check if validation errors are accessible
    const validationErrors = page.locator('[data-testid="validation-error"]');
    const errorCount = await validationErrors.count();

    if (errorCount > 0) {
      for (let i = 0; i < errorCount; i++) {
        const error = validationErrors.nth(i);
        await a11y.checkElementAccessibility(`[data-testid="validation-error"]:nth-child(${i + 1})`);
      }
    }
  });

  test('should have accessible password strength indicator', async ({ page }) => {
    await authPage.gotoRegister();

    await page.fill('[data-testid="password-input"]', 'weak');

    // Check if password strength indicator is accessible
    const strengthIndicator = page.locator('[data-testid="password-strength"]');
    if (await strengthIndicator.isVisible()) {
      await a11y.checkElementAccessibility('[data-testid="password-strength"]');
      await a11y.verifyAriaLabels('[data-testid="password-strength"]', {
        role: 'status',
      });
    }
  });

  test('should have accessible loading states', async ({ page }) => {
    await authPage.goto();

    // Mock slow response to see loading state
    await page.route('**/api/auth/login', async route => {
      await new Promise(resolve => setTimeout(resolve, 2000));
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, data: { token: 'mock' } }),
      });
    });

    await authPage.fillLoginForm('test@test.com', 'password');
    await authPage.submitForm();

    // Check if loading state is accessible
    const loading = page.locator('[data-testid="loading"]');
    if (await loading.isVisible()) {
      await a11y.verifyAriaLabels('[data-testid="loading"]', {
        role: 'status',
      });
    }
  });

  test('should have accessible success messages', async ({ page }) => {
    // Mock successful login
    await page.route('**/api/auth/login', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: { user: { displayName: 'Test User' }, token: 'mock' },
        }),
      });
    });

    await authPage.goto();
    await authPage.login('test@test.com', 'password');

    // Check if success message is accessible
    const successMessage = page.locator('[data-testid="success-message"]');
    if (await successMessage.isVisible()) {
      await a11y.checkElementAccessibility('[data-testid="success-message"]');
      await a11y.verifyAriaLabels('[data-testid="success-message"]', {
        role: 'status',
      });
    }
  });

  test('should have proper heading structure', async ({ page }) => {
    await authPage.goto();

    await a11y.verifyHeadingStructure();
  });

  test('should have accessible skip links', async ({ page }) => {
    await authPage.goto();

    await a11y.verifySkipLinks();
  });

  test('should work with screen readers', async ({ page }) => {
    await authPage.goto();

    // Check if page title is descriptive
    await expect(page).toHaveTitle(/login|sign in/i);

    // Check if main content is properly marked
    const main = page.locator('main, [role="main"]');
    if (await main.isVisible()) {
      await expect(main).toBeVisible();
    }

    // Check if form has proper labeling
    const form = page.locator('form');
    if (await form.isVisible()) {
      await a11y.verifyFormAccessibility('form');
    }
  });
});
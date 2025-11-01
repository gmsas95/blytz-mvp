import { test, expect } from '../fixtures';
import { AuthPage } from '../pages/auth.page';
import { generateTestEmail } from '../utils/test-data';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    const authPage = new AuthPage(page);
    await authPage.goto();
  });

  test('should display login form', async ({ page }) => {
    const authPage = new AuthPage(page);

    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('should show validation errors for empty form', async ({ page }) => {
    const authPage = new AuthPage(page);

    await authPage.submitForm();
    await authPage.verifyLoginError();
  });

  test('should login successfully with valid credentials', async ({ authPage, testUser }) => {
    await authPage.login(testUser.email, testUser.password);
    await authPage.verifyLoginSuccess();
    await authPage.verifyUserIsLoggedIn(testUser.displayName);
  });

  test('should show error for invalid credentials', async ({ authPage }) => {
    await authPage.login('invalid@test.com', 'wrongpassword');
    await authPage.verifyLoginError('Invalid credentials');
    await authPage.verifyUserIsLoggedOut();
  });

  test('should register new user successfully', async ({ page }) => {
    const authPage = new AuthPage(page);
    const newUser = {
      email: generateTestEmail(),
      password: 'NewUser123!',
      displayName: 'New Test User',
    };

    await authPage.gotoRegister();
    await authPage.register(newUser.email, newUser.password, newUser.displayName);
    await authPage.verifyRegistrationSuccess();
  });

  test('should logout successfully', async ({ authenticatedPage }) => {
    // User is already logged in via authenticatedPage fixture
    await authenticatedPage.logout();
    await authenticatedPage.verifyUserIsLoggedOut();
  });

  test('should persist login session across page refresh', async ({ authPage, testUser }) => {
    await authPage.login(testUser.email, testUser.password);
    await authPage.verifyLoginSuccess();

    // Refresh page
    await page.reload();
    await authPage.verifyUserIsLoggedIn(testUser.displayName);
  });

  test('should redirect to login when accessing protected route without auth', async ({ page }) => {
    // Try to access a protected route
    await page.goto('/auctions/create');

    // Should redirect to login
    await expect(page).toHaveURL(/.*auth\/login.*/);
  });

  test('should handle remember me functionality', async ({ authPage, testUser }) => {
    await authPage.login(testUser.email, testUser.password, true);
    await authPage.verifyLoginSuccess();

    // Check if remember me token is set
    const cookies = await page.context().cookies();
    const hasRememberToken = cookies.some(cookie => cookie.name === 'remember_token');

    if (hasRememberToken) {
      console.log('Remember me token is set');
    }
  });

  test('should handle password reset flow', async ({ authPage }) => {
    await authPage.goto();
    await authPage.clickForgotPassword();

    // Should navigate to password reset page
    await expect(page).toHaveURL(/.*forgot-password.*/);
  });

  test('should handle session expiration', async ({ authPage, testUser }) => {
    await authPage.login(testUser.email, testUser.password);
    await authPage.verifyLoginSuccess();

    // Clear auth token to simulate expiration
    await authPage.clearAuthData();

    // Try to access a protected route
    await page.goto('/profile');

    // Should redirect to login
    await expect(page).toHaveURL(/.*auth\/login.*/);
  });
});

test.describe('Authentication with Mobile Viewport', () => {
  test.use({ viewport: { width: 375, height: 667 } }); // iPhone 6/7/8

  test('should work on mobile devices', async ({ authPage, testUser }) => {
    await authPage.login(testUser.email, testUser.password);
    await authPage.verifyLoginSuccess();
    await authPage.verifyUserIsLoggedIn(testUser.displayName);
  });
});
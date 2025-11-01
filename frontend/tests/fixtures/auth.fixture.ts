import { test as base, expect } from '@playwright/test';
import { AuthPage } from '../pages/auth.page';
import { TestUser } from '../utils/test-data';

// Define auth fixture types
export interface AuthFixtures {
  authPage: AuthPage;
  authenticatedPage: AuthPage;
  testUser: TestUser;
  adminUser: TestUser;
}

// Extend base test with auth fixtures
export const test = base.extend<AuthFixtures>({
  authPage: async ({ page }, use) => {
    const authPage = new AuthPage(page);
    await use(authPage);
  },

  testUser: async ({}, use) => {
    const user: TestUser = {
      email: 'testuser@blytz.app',
      password: 'TestPassword123!',
      displayName: 'Test User',
      role: 'user',
    };
    await use(user);
  },

  adminUser: async ({}, use) => {
    const admin: TestUser = {
      email: 'admin@blytz.app',
      password: 'AdminPassword123!',
      displayName: 'Admin User',
      role: 'admin',
    };
    await use(admin);
  },

  authenticatedPage: async ({ page, testUser }, use) => {
    const authPage = new AuthPage(page);

    // Login the test user
    await authPage.goto();
    await authPage.login(testUser.email, testUser.password);

    // Verify successful login
    await expect(page.locator('[data-testid="user-menu"]')).toBeVisible();

    await use(authPage);

    // Cleanup: logout after test
    await authPage.logout();
  },
});

export { expect } from '@playwright/test';
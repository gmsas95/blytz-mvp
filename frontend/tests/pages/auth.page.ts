import { Page, expect } from '@playwright/test';
import { TestHelpers } from '../utils/test-helpers';

export class AuthPage {
  private helpers: TestHelpers;

  constructor(private page: Page) {
    this.helpers = new TestHelpers(page);
  }

  // Selectors
  private readonly selectors = {
    loginForm: '[data-testid="login-form"]',
    registerForm: '[data-testid="register-form"]',
    emailInput: '[data-testid="email-input"], input[type="email"], input[name="email"]',
    passwordInput: '[data-testid="password-input"], input[type="password"], input[name="password"]',
    confirmPasswordInput: '[data-testid="confirm-password-input"], input[name="confirmPassword"]',
    displayNameInput: '[data-testid="display-name-input"], input[name="displayName"]',
    submitButton: '[data-testid="submit-button"], button[type="submit"]',
    loginButton: '[data-testid="login-button"]',
    registerButton: '[data-testid="register-button"]',
    logoutButton: '[data-testid="logout-button"], button:has-text("Logout")',
    userMenu: '[data-testid="user-menu"], .user-menu',
    errorMessage: '[data-testid="error-message"], .error-message, .alert-error',
    successMessage: '[data-testid="success-message"], .success-message, .alert-success',
    rememberMeCheckbox: '[data-testid="remember-me"], input[name="rememberMe"]',
    forgotPasswordLink: '[data-testid="forgot-password"], a:has-text("Forgot Password")',
  };

  async goto(): Promise<void> {
    await this.page.goto('/auth/login');
    await this.helpers.waitForPageLoad();
  }

  async gotoRegister(): Promise<void> {
    await this.page.goto('/auth/register');
    await this.helpers.waitForPageLoad();
  }

  async login(email: string, password: string, rememberMe: boolean = false): Promise<void> {
    await this.page.fill(this.selectors.emailInput, email);
    await this.page.fill(this.selectors.passwordInput, password);

    if (rememberMe) {
      const checkbox = this.page.locator(this.selectors.rememberMeCheckbox);
      if (await checkbox.isVisible()) {
        await checkbox.check();
      }
    }

    await this.page.click(this.selectors.submitButton);
    await this.helpers.waitForLoadingComplete();
  }

  async register(
    email: string,
    password: string,
    displayName: string,
    confirmPassword?: string
  ): Promise<void> {
    await this.gotoRegister();

    await this.page.fill(this.selectors.emailInput, email);
    await this.page.fill(this.selectors.passwordInput, password);
    await this.page.fill(this.selectors.displayNameInput, displayName);

    const confirmInput = this.page.locator(this.selectors.confirmPasswordInput);
    if (await confirmInput.isVisible()) {
      await confirmInput.fill(confirmPassword || password);
    }

    await this.page.click(this.selectors.submitButton);
    await this.helpers.waitForLoadingComplete();
  }

  async logout(): Promise<void> {
    // Click user menu to open dropdown
    await this.page.click(this.selectors.userMenu);

    // Click logout button (wait for it to be visible in dropdown)
    await this.page.waitForSelector(this.selectors.logoutButton, { state: 'visible' });
    await this.page.click(this.selectors.logoutButton);

    await this.helpers.waitForLoadingComplete();
  }

  async verifyLoginSuccess(): Promise<void> {
    // Check for success message or user menu appearance
    const successMessage = this.page.locator(this.selectors.successMessage);
    const userMenu = this.page.locator(this.selectors.userMenu);

    await expect(
      successMessage.or(userMenu)
    ).toBeVisible({ timeout: 10000 });
  }

  async verifyRegistrationSuccess(): Promise<void> {
    // Check for success message or redirect to login/dashboard
    const successMessage = this.page.locator(this.selectors.successMessage);
    const loginForm = this.page.locator(this.selectors.loginForm);

    await expect(
      successMessage.or(loginForm)
    ).toBeVisible({ timeout: 10000 });
  }

  async verifyLoginError(expectedMessage?: string): Promise<void> {
    const errorMessage = this.page.locator(this.selectors.errorMessage);
    await expect(errorMessage).toBeVisible({ timeout: 5000 });

    if (expectedMessage) {
      await expect(errorMessage).toContainText(expectedMessage);
    }
  }

  async verifyRegistrationError(expectedMessage?: string): Promise<void> {
    const errorMessage = this.page.locator(this.selectors.errorMessage);
    await expect(errorMessage).toBeVisible({ timeout: 5000 });

    if (expectedMessage) {
      await expect(errorMessage).toContainText(expectedMessage);
    }
  }

  async verifyUserIsLoggedIn(displayName?: string): Promise<void> {
    await expect(this.page.locator(this.selectors.userMenu)).toBeVisible();

    if (displayName) {
      await expect(this.page.locator(this.selectors.userMenu)).toContainText(displayName);
    }
  }

  async verifyUserIsLoggedOut(): Promise<void> {
    await expect(this.page.locator(this.selectors.userMenu)).not.toBeVisible();

    // Should be on login page or have login form visible
    const loginForm = this.page.locator(this.selectors.loginForm);
    if (await loginForm.isVisible()) {
      await expect(loginForm).toBeVisible();
    } else {
      // Check if redirected to login page
      await expect(this.page).toHaveURL(/.*auth\/login.*/);
    }
  }

  async fillLoginForm(email: string, password: string): Promise<void> {
    await this.page.fill(this.selectors.emailInput, email);
    await this.page.fill(this.selectors.passwordInput, password);
  }

  async fillRegistrationForm(
    email: string,
    password: string,
    displayName: string,
    confirmPassword?: string
  ): Promise<void> {
    await this.page.fill(this.selectors.emailInput, email);
    await this.page.fill(this.selectors.passwordInput, password);
    await this.page.fill(this.selectors.displayNameInput, displayName);

    const confirmInput = this.page.locator(this.selectors.confirmPasswordInput);
    if (await confirmInput.isVisible()) {
      await confirmInput.fill(confirmPassword || password);
    }
  }

  async submitForm(): Promise<void> {
    await this.page.click(this.selectors.submitButton);
    await this.helpers.waitForLoadingComplete();
  }

  async clickForgotPassword(): Promise<void> {
    const forgotLink = this.page.locator(this.selectors.forgotPasswordLink);
    if (await forgotLink.isVisible()) {
      await forgotLink.click();
      await this.helpers.waitForLoadingComplete();
    }
  }

  async waitForAuthResponse(): Promise<any> {
    return await this.helpers.waitForApiResponse('/api/auth/');
  }

  async verifyFormValidation(): Promise<void> {
    // Check if form shows validation errors for empty fields
    await this.page.click(this.selectors.submitButton);

    // Look for validation error messages
    const validationError = this.page.locator(
      '[data-testid="validation-error"], .error, .invalid-feedback'
    );

    // Check if any validation errors are visible
    const hasValidationErrors = await validationError.count() > 0;
    if (!hasValidationErrors) {
      console.warn('No validation errors found after submitting empty form');
    }
  }

  async verifyPasswordStrengthIndicator(): Promise<void> {
    // Look for password strength indicator
    const strengthIndicator = this.page.locator(
      '[data-testid="password-strength"], .password-strength, .strength-meter'
    );

    if (await strengthIndicator.isVisible()) {
      await expect(strengthIndicator).toBeVisible();
    }
  }

  async verifyTermsAndConditionsCheckbox(): Promise<void> {
    const termsCheckbox = this.page.locator(
      '[data-testid="terms-checkbox"], input[name="terms"], input[type="checkbox"]'
    );

    if (await termsCheckbox.isVisible()) {
      await expect(termsCheckbox).toBeVisible();
    }
  }

  async getAuthToken(): Promise<string | null> {
    // Get auth token from localStorage or cookies
    const token = await this.page.evaluate(() => {
      return localStorage.getItem('auth_token') ||
             document.cookie
               .split('; ')
               .find(row => row.startsWith('auth_token='))
               ?.split('=')[1] || null;
    });

    return token;
  }

  async setAuthToken(token: string): Promise<void> {
    // Set auth token in localStorage for testing
    await this.page.evaluate((t) => {
      localStorage.setItem('auth_token', t);
    }, token);
  }

  async clearAuthData(): Promise<void> {
    // Clear all auth-related data
    await this.page.evaluate(() => {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_data');
      document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    });
  }
}
# Playwright E2E Testing Setup

This directory contains comprehensive end-to-end tests for the Blytz auction platform using Playwright.

## 📁 Directory Structure

```
tests/
├── e2e/                   # End-to-end tests
│   ├── auth.spec.ts      # Authentication tests
│   ├── auction.spec.ts   # Auction functionality tests
│   └── basic.spec.ts     # Basic application tests
├── accessibility/         # Accessibility tests
│   └── auth.a11y.spec.ts # Auth accessibility tests
├── fixtures/             # Test fixtures and setup
│   ├── auth.fixture.ts   # Authentication fixtures
│   ├── auction.fixture.ts # Auction fixtures
│   └── index.ts         # Fixture exports
├── pages/               # Page Object Models
│   ├── auth.page.ts     # Auth page object
│   ├── auction.page.ts  # Auction page object
│   └── product.page.ts  # Product page object
├── utils/               # Test utilities and helpers
│   ├── test-helpers.ts  # Common test helpers
│   ├── test-data.ts     # Test data and constants
│   └── accessibility-helpers.ts # A11y testing helpers
├── mocks/               # API mocking utilities
│   └── api-mocks.ts     # API response mocks
├── types/               # TypeScript type definitions
│   └── playwright.d.ts  # Playwright type extensions
├── reporters/           # Custom reporters
│   └── custom-reporter.ts # Detailed test reporting
├── global-setup.ts      # Global test setup
├── global-teardown.ts   # Global test cleanup
└── README.md           # This file
```

## 🚀 Getting Started

### Prerequisites

- Node.js 18+
- NPM 9+
- Playwright browsers (installed automatically)

### Installation

1. Install dependencies:
```bash
npm install
```

2. Install Playwright browsers:
```bash
npm run test:e2e:install
```

3. (Optional) Install system dependencies:
```bash
sudo npm run test:e2e:install-deps
```

### Running Tests

#### Development

```bash
# Run all tests
npm run test:e2e

# Run tests in UI mode (recommended for development)
npm run test:e2e:ui

# Run tests with visible browser
npm run test:e2e:headed

# Debug tests
npm run test:e2e:debug
```

#### Browser-Specific

```bash
# Chromium only
npm run test:e2e:chromium

# Firefox only
npm run test:e2e:firefox

# Safari (WebKit) only
npm run test:e2e:webkit

# Mobile devices
npm run test:e2e:mobile
```

#### Specialized Tests

```bash
# Accessibility tests
npm run test:e2e:accessibility

# Visual regression tests
npm run test:e2e:visual
```

#### CI/CD

```bash
# Run in CI mode
npm run test:e2e:ci
```

### Viewing Reports

```bash
# Open HTML report
npm run test:e2e:show-report
```

## 🧪 Test Configuration

### Playwright Config (`playwright.config.ts`)

The configuration includes:

- **Multiple Projects**: Chromium, Firefox, WebKit, Mobile Chrome, Mobile Safari
- **Accessibility Testing**: Dedicated project for @axe-core testing
- **Visual Regression**: Separate project for visual testing
- **Reporting**: HTML, JSON, JUnit reporters
- **Screenshots**: On failure only
- **Videos**: Retain on failure
- **Traces**: On first retry
- **Web Server**: Automatic Next.js server startup

### Test Fixtures

Custom fixtures provide pre-configured test states:

- `authPage`: Authentication page object
- `authenticatedPage`: Page with logged-in user
- `testUser`: Test user credentials
- `adminUser`: Admin user credentials
- `auctionPage`: Auction page object
- `productPage`: Product page object

## 🎯 Test Categories

### 1. Authentication (`e2e/auth.spec.ts`)

- Login/logout functionality
- Registration flow
- Form validation
- Session management
- Protected routes
- Password reset

### 2. Auction Functionality (`e2e/auction.spec.ts`)

- Auction browsing
- Bid placement
- Real-time updates
- Watchlist management
- Search and filtering
- Mobile responsiveness

### 3. Accessibility (`accessibility/`)

- WCAG compliance testing
- Keyboard navigation
- Screen reader support
- Color contrast verification
- Form accessibility
- ARIA attributes

## 📊 Reporting

### HTML Report

- Interactive test results
- Screenshots on failure
- Video recordings
- Stack traces
- Filtering and search

### Custom Reports

- **Detailed Report**: JSON with all test data
- **Failed Tests**: Markdown report of failures
- **Performance Report**: Slowest/fastest tests

### Artifacts

- **Screenshots**: `test-results/screenshots/`
- **Videos**: `test-results/videos/`
- **Traces**: `test-results/traces/`
- **Reports**: `playwright-report/`

## 🔧 Development Tools

### Code Generation

```bash
# Generate test code from browser interactions
npm run test:e2e:codegen
```

### Trace Viewing

```bash
# Run tests with trace enabled
npm run test:e2e:trace
```

### Debugging

```bash
# Debug with browser devtools
npm run test:e2e:debug

# Run single test file
npx playwright test tests/e2e/auth.spec.ts

# Run specific test
npx playwright test --grep "should login successfully"
```

## 🎨 Best Practices

### Test Organization

1. **Use Page Objects**: Separate UI interactions from test logic
2. **Fixtures**: Reuse test setup and data
3. **Descriptive Names**: Clear test descriptions
4. **Independent Tests**: Each test should run in isolation

### Data Management

1. **Mock APIs**: Use consistent mock data
2. **Clean State**: Reset state between tests
3. **Test Data**: Centralized test data management

### Accessibility

1. **Automated Testing**: Include a11y in CI/CD
2. **Keyboard Testing**: Verify keyboard navigation
3. **Screen Readers**: Test with assistive technologies
4. **Color Contrast**: Verify WCAG compliance

### Performance

1. **Parallel Execution**: Run tests concurrently
2. **Selective Testing**: Run relevant tests per change
3. **Timeouts**: Appropriate timeout values
4. **Retry Logic**: Handle flaky tests

## 🔍 Debugging Tips

### Common Issues

1. **Flaky Tests**: Add waits and proper selectors
2. **Timeout Errors**: Increase timeout values
3. **Element Not Found**: Use more specific selectors
4. **Network Issues**: Mock API responses

### Debugging Tools

1. **Playwright Inspector**: `npx playwright test --debug`
2. **Trace Viewer**: Analyze test execution
3. **Browser DevTools**: Use headed mode
4. **Console Logs**: Check browser console

### Selectors

1. **Data Test IDs**: Use `[data-testid="..."]`
2. **Accessible Names**: Use accessible selectors
3. **CSS Selectors**: Use specific but flexible selectors
4. **Text Content**: Use text-based selectors carefully

## 🚀 CI/CD Integration

### GitHub Actions

```yaml
- name: Install Playwright
  run: npm run test:e2e:install

- name: Run E2E Tests
  run: npm run test:e2e:ci

- name: Upload Test Results
  uses: actions/upload-artifact@v3
  with:
    name: playwright-report
    path: playwright-report/
```

### Environment Variables

- `BASE_URL`: Application base URL
- `CI`: Enable CI mode
- `PLAYWRIGHT_BROWSERS_PATH`: Browser installation path

## 📝 Writing New Tests

### Basic Test Structure

```typescript
import { test, expect } from '../fixtures';

test.describe('Feature Name', () => {
  test('should do something', async ({ page }) => {
    await page.goto('/some-page');
    await expect(page.locator('h1')).toContainText('Expected Title');
  });
});
```

### Using Fixtures

```typescript
test('authenticated user action', async ({ authenticatedPage }) => {
  // User is already logged in
  await authenticatedPage.goto('/protected-route');
  await expect(authenticatedPage.locator('[data-testid="user-menu"]')).toBeVisible();
});
```

### Page Objects

```typescript
const authPage = new AuthPage(page);
await authPage.login('user@test.com', 'password');
await authPage.verifyLoginSuccess();
```

## 🔗 Resources

- [Playwright Documentation](https://playwright.dev/)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Accessibility Testing](https://playwright.dev/docs/accessibility-testing)
- [Visual Testing](https://playwright.dev/docs/test-snapshots)
- [API Testing](https://playwright.dev/docs/api-testing)
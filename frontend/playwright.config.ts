import { defineConfig, devices } from '@playwright/test';
import { devices as mobileDevices } from '@playwright/test';
import path from 'path';

/**
 * @see https://playwright.dev/docs/test-configuration
 */
export default defineConfig({
  testDir: './tests',
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [
    ['html', {
      outputFolder: 'playwright-report',
      open: process.env.CI ? 'never' : 'on-failure'
    }],
    ['json', { outputFile: 'test-results/results.json' }],
    ['junit', { outputFile: 'test-results/results.xml' }],
    ['list'],
  ],
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.BASE_URL || 'http://localhost:3000',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',

    /* Take screenshot on failure */
    screenshot: 'only-on-failure',

    /* Record video on failure */
    video: 'retain-on-failure',

    /* Global timeout for each action */
    actionTimeout: 10 * 1000, // 10 seconds

    /* Global timeout for navigation */
    navigationTimeout: 30 * 1000, // 30 seconds

    /* Ignore HTTPS errors */
    ignoreHTTPSErrors: true,

    /* User agent */
    userAgent: 'Blytz-E2E-Tests',

    /* Locale */
    locale: 'en-US',

    /* Timezone */
    timezoneId: 'America/New_York',

    /* Color scheme */
    colorScheme: 'light',

    /* Add custom headers for API requests */
    extraHTTPHeaders: {
      'X-Test-Environment': 'playwright',
    },
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        viewport: { width: 1280, height: 720 },
        contextOptions: {
          permissions: ['clipboard-read', 'clipboard-write'],
        },
      },
      testIgnore: '**/*.mobile.spec.ts',
    },

    {
      name: 'firefox',
      use: {
        ...devices['Desktop Firefox'],
        viewport: { width: 1280, height: 720 },
      },
      testIgnore: '**/*.mobile.spec.ts',
    },

    {
      name: 'webkit',
      use: {
        ...devices['Desktop Safari'],
        viewport: { width: 1280, height: 720 },
      },
      testIgnore: '**/*.mobile.spec.ts',
    },

    /* Test against mobile viewports. */
    {
      name: 'Mobile Chrome',
      use: {
        ...mobileDevices['Pixel 5'],
        contextOptions: {
          permissions: ['clipboard-read', 'clipboard-write'],
        },
      },
      testMatch: '**/*.mobile.spec.ts',
    },

    {
      name: 'Mobile Safari',
      use: {
        ...mobileDevices['iPhone 12'],
      },
      testMatch: '**/*.mobile.spec.ts',
    },

    /* Test against branded browsers. */
    {
      name: 'Microsoft Edge',
      use: {
        ...devices['Desktop Edge'],
        channel: 'msedge',
        viewport: { width: 1280, height: 720 },
      },
      testIgnore: '**/*.mobile.spec.ts',
    },

    /* Accessibility testing project */
    {
      name: 'accessibility',
      use: {
        ...devices['Desktop Chrome'],
        viewport: { width: 1280, height: 720 },
      },
      testMatch: '**/*.a11y.spec.ts',
      dependencies: ['chromium'],
    },

    /* Visual regression testing */
    {
      name: 'visual-regression',
      use: {
        ...devices['Desktop Chrome'],
        viewport: { width: 1280, height: 720 },
        screenshot: 'only-on-failure',
        ignoreHTTPSErrors: true,
      },
      testMatch: '**/*.visual.spec.ts',
    },
  ],

  /* Global setup and teardown */
  globalSetup: require.resolve('./tests/global-setup.ts'),
  globalTeardown: require.resolve('./tests/global-teardown.ts'),

  /* Test timeout */
  timeout: 60 * 1000, // 60 seconds

  /* Expect timeout */
  expect: {
    timeout: 10 * 1000, // 10 seconds
  },

  /* Output directory for test artifacts */
  outputDir: 'test-results',

  /* Web server configuration for running the app during tests */
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000, // 2 minutes
  },

  /* Metadata */
  metadata: {
    'Test Environment': 'E2E Testing',
    'Application': 'Blytz Auction Platform',
    'Browser Support': ['Chromium', 'Firefox', 'WebKit', 'Mobile'],
  },
});
import { chromium, FullConfig } from '@playwright/test';
import path from 'path';

async function globalSetup(config: FullConfig) {
  console.log('🚀 Starting Playwright global setup...');

  // Ensure test directories exist
  const testDirs = [
    'test-results',
    'test-results/screenshots',
    'test-results/videos',
    'test-results/traces',
    'playwright-report',
  ];

  for (const dir of testDirs) {
    try {
      await fs.promises.mkdir(dir, { recursive: true });
    } catch (error) {
      console.log(`Directory ${dir} already exists or could not be created`);
    }
  }

  // Set up global test data or state if needed
  const browser = await chromium.launch();
  const context = await browser.newContext();

  // Example: Set up authentication tokens for testing
  try {
    const page = await context.newPage();

    // Navigate to the app and ensure it's ready
    await page.goto(config.webServer?.url || 'http://localhost:3000');
    await page.waitForLoadState('networkidle');

    // Check if health endpoint is available
    const healthResponse = await page.goto('/health').catch(() => null);
    if (healthResponse && healthResponse.status() === 200) {
      console.log('✅ Application is healthy and ready for testing');
    } else {
      console.warn('⚠️  Health check failed, but continuing with tests');
    }

    await page.close();
  } catch (error) {
    console.warn('⚠️  Could not verify application health during setup:', error);
  }

  await context.close();
  await browser.close();

  console.log('✅ Playwright global setup completed');
}

// Import fs for directory creation
import * as fs from 'fs';

export default globalSetup;
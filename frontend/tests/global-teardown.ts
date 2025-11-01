import { FullConfig } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';

async function globalTeardown(config: FullConfig) {
  console.log('🧹 Starting Playwright global teardown...');

  // Clean up temporary files or state if needed
  const tempDirs = [
    'test-results/temp',
    'test-results/.tmp',
  ];

  for (const dir of tempDirs) {
    try {
      if (fs.existsSync(dir)) {
        fs.rmSync(dir, { recursive: true, force: true });
        console.log(`🗑️  Cleaned up temporary directory: ${dir}`);
      }
    } catch (error) {
      console.warn(`⚠️  Could not clean up directory ${dir}:`, error);
    }
  }

  // Generate test summary report if needed
  try {
    const resultsDir = 'test-results';
    if (fs.existsSync(resultsDir)) {
      const files = fs.readdirSync(resultsDir);
      const testFiles = files.filter(file =>
        file.endsWith('.json') ||
        file.endsWith('.xml') ||
        file.endsWith('.txt')
      );

      console.log(`📊 Test artifacts generated: ${testFiles.length} files`);
    }
  } catch (error) {
    console.warn('⚠️  Could not generate test summary:', error);
  }

  console.log('✅ Playwright global teardown completed');
}

export default globalTeardown;
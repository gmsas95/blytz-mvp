import { FullConfig, FullResult, Reporter, Suite, TestCase, TestResult } from '@playwright/test/reporter';
import * as fs from 'fs';
import * as path from 'path';

class CustomReporter implements Reporter {
  private config: FullConfig;
  private testResults: any[] = [];
  private startTime: number = Date.now();

  onBegin(config: FullConfig, suite: Suite) {
    console.log('\n🚀 Starting Playwright Test Run');
    console.log(`📊 Total test files: ${suite.allTests().length}`);
    console.log(`🌐 Projects: ${config.projects.map(p => p.name).join(', ')}`);
    this.config = config;
    this.startTime = Date.now();
  }

  onTestEnd(test: TestCase, result: TestResult) {
    const testResult = {
      title: test.title,
      file: test.location?.file,
      line: test.location?.line,
      status: result.status,
      duration: result.duration,
      errors: result.errors.map(e => e.message),
      retries: result.retry,
      attachments: result.attachments.map(a => a.name),
      project: test.parent.project()?.name,
    };

    this.testResults.push(testResult);

    // Log test result with emoji
    const statusEmoji = {
      passed: '✅',
      failed: '❌',
      skipped: '⏭️',
      timedOut: '⏰',
    }[result.status] || '❓';

    const duration = `${result.duration}ms`;
    const retryText = result.retry > 0 ? ` (retry ${result.retry})` : '';
    console.log(`${statusEmoji} ${test.title} - ${duration}${retryText}`);

    // Log errors if any
    if (result.errors.length > 0) {
      result.errors.forEach(error => {
        console.log(`   ❌ Error: ${error.message}`);
      });
    }
  }

  onEnd(result: FullResult) {
    const duration = Date.now() - this.startTime;
    const totalTests = this.testResults.length;
    const passedTests = this.testResults.filter(t => t.status === 'passed').length;
    const failedTests = this.testResults.filter(t => t.status === 'failed').length;
    const skippedTests = this.testResults.filter(t => t.status === 'skipped').length;
    const timedOutTests = this.testResults.filter(t => t.status === 'timedOut').length;

    console.log('\n📋 Test Summary:');
    console.log(`⏱️  Total duration: ${duration}ms`);
    console.log(`📊 Total tests: ${totalTests}`);
    console.log(`✅ Passed: ${passedTests}`);
    console.log(`❌ Failed: ${failedTests}`);
    console.log(`⏭️  Skipped: ${skippedTests}`);
    console.log(`⏰ Timed out: ${timedOutTests}`);

    // Generate detailed report
    this.generateDetailedReport();

    // Generate failed tests report
    if (failedTests > 0) {
      this.generateFailedTestsReport();
    }

    // Generate performance report
    this.generatePerformanceReport();
  }

  private generateDetailedReport(): void {
    const reportDir = 'test-results';
    if (!fs.existsSync(reportDir)) {
      fs.mkdirSync(reportDir, { recursive: true });
    }

    const report = {
      timestamp: new Date().toISOString(),
      duration: Date.now() - this.startTime,
      summary: {
        total: this.testResults.length,
        passed: this.testResults.filter(t => t.status === 'passed').length,
        failed: this.testResults.filter(t => t.status === 'failed').length,
        skipped: this.testResults.filter(t => t.status === 'skipped').length,
        timedOut: this.testResults.filter(t => t.status === 'timedOut').length,
      },
      results: this.testResults,
      projects: this.config.projects.map(p => p.name),
    };

    const reportPath = path.join(reportDir, 'detailed-report.json');
    fs.writeFileSync(reportPath, JSON.stringify(report, null, 2));
    console.log(`📄 Detailed report generated: ${reportPath}`);
  }

  private generateFailedTestsReport(): void {
    const failedTests = this.testResults.filter(t => t.status === 'failed');
    const reportDir = 'test-results';

    const report = {
      timestamp: new Date().toISOString(),
      totalFailed: failedTests.length,
      failures: failedTests.map(test => ({
        title: test.title,
        file: test.file,
        line: test.line,
        errors: test.errors,
        project: test.project,
      })),
    };

    const reportPath = path.join(reportDir, 'failed-tests.json');
    fs.writeFileSync(reportPath, JSON.stringify(report, null, 2));
    console.log(`❌ Failed tests report generated: ${reportPath}`);

    // Generate markdown report for easy reading
    const markdownReport = this.generateMarkdownReport(failedTests);
    const markdownPath = path.join(reportDir, 'failed-tests.md');
    fs.writeFileSync(markdownPath, markdownReport);
    console.log(`📝 Failed tests markdown report generated: ${markdownPath}`);
  }

  private generateMarkdownReport(failedTests: any[]): string {
    let markdown = '# Failed Tests Report\n\n';
    markdown += `**Generated:** ${new Date().toISOString()}\n`;
    markdown += `**Total Failed:** ${failedTests.length}\n\n`;

    failedTests.forEach((test, index) => {
      markdown += `## ${index + 1}. ${test.title}\n\n`;
      markdown += `- **File:** ${test.file}:${test.line}\n`;
      markdown += `- **Project:** ${test.project}\n\n`;

      if (test.errors && test.errors.length > 0) {
        markdown += '### Errors:\n\n';
        test.errors.forEach((error: string) => {
          markdown += '```\n';
          markdown += error;
          markdown += '\n```\n\n';
        });
      }
      markdown += '---\n\n';
    });

    return markdown;
  }

  private generatePerformanceReport(): void {
    const reportDir = 'test-results';
    const slowTests = this.testResults
      .filter(t => t.status === 'passed')
      .sort((a, b) => b.duration - a.duration)
      .slice(0, 10);

    const fastTests = this.testResults
      .filter(t => t.status === 'passed')
      .sort((a, b) => a.duration - b.duration)
      .slice(0, 10);

    const performanceData = {
      timestamp: new Date().toISOString(),
      slowestTests: slowTests.map(test => ({
        title: test.title,
        duration: test.duration,
        file: test.file,
      })),
      fastestTests: fastTests.map(test => ({
        title: test.title,
        duration: test.duration,
        file: test.file,
      })),
      averageDuration: this.testResults.reduce((sum, t) => sum + t.duration, 0) / this.testResults.length,
    };

    const reportPath = path.join(reportDir, 'performance-report.json');
    fs.writeFileSync(reportPath, JSON.stringify(performanceData, null, 2));
    console.log(`⚡ Performance report generated: ${reportPath}`);
  }

  onError(error: Error) {
    console.error(`\n🚨 Reporter Error: ${error.message}`);
  }
}

export default CustomReporter;
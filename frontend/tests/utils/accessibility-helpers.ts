import { Page, expect } from '@playwright/test';
import { injectAxe, getViolations, checkA11y } from '@axe-core/playwright';

/**
 * Accessibility testing helpers
 */
export class AccessibilityHelpers {
  constructor(private page: Page) {}

  /**
   * Initialize axe-core for accessibility testing
   */
  async initializeAxe(): Promise<void> {
    await injectAxe(this.page);
  }

  /**
   * Run accessibility tests on the entire page
   */
  async checkPageAccessibility(
    options?: {
      include?: string[];
      exclude?: string[];
      rules?: Record<string, any>;
    }
  ): Promise<void> {
    await checkA11y(this.page, null, {
      detailedReport: true,
      detailedReportOptions: { html: true },
      ...options,
    });
  }

  /**
   * Run accessibility tests on a specific element
   */
  async checkElementAccessibility(
    selector: string,
    options?: {
      include?: string[];
      exclude?: string[];
      rules?: Record<string, any>;
    }
  ): Promise<void> {
    await checkA11y(this.page, selector, {
      detailedReport: true,
      detailedReportOptions: { html: true },
      ...options,
    });
  }

  /**
   * Get all accessibility violations
   */
  async getAccessibilityViolations(
    context?: string,
    options?: { include?: string[]; exclude?: string[] }
  ): Promise<any[]> {
    return await getViolations(this.page, context, options);
  }

  /**
   * Verify color contrast ratios
   */
  async verifyColorContrast(selector: string, expectedRatio: number = 4.5): Promise<void> {
    const violations = await getViolations(this.page, selector, {
      rules: {
        'color-contrast': { enabled: true },
      },
    });

    if (violations.length > 0) {
      const contrastViolations = violations.filter(v => v.id === 'color-contrast');
      expect(contrastViolations).toHaveLength(0,
        `Color contrast violations found for ${selector}: ${JSON.stringify(contrastViolations, null, 2)}`
      );
    }
  }

  /**
   * Verify keyboard navigation
   */
  async verifyKeyboardNavigation(selector: string): Promise<void> {
    const element = this.page.locator(selector);

    // Check if element is focusable
    const tabIndex = await element.getAttribute('tabindex');
    const isFocusable = tabIndex !== null && parseInt(tabIndex) >= 0;

    expect(isFocusable).toBe(true, `Element ${selector} is not keyboard focusable`);

    // Test keyboard navigation
    await element.focus();
    await expect(element).toBeFocused();

    // Test common keyboard interactions
    await this.page.keyboard.press('Enter');
    await this.page.keyboard.press('Space');
    await this.page.keyboard.press('Escape');
    await this.page.keyboard.press('Tab');
  }

  /**
   * Verify ARIA labels and descriptions
   */
  async verifyAriaLabels(selector: string, expectedLabels: {
    label?: string;
    description?: string;
    role?: string;
  }): Promise<void> {
    const element = this.page.locator(selector);

    if (expectedLabels.label) {
      await expect(element).toHaveAttribute('aria-label', expectedLabels.label);
    }

    if (expectedLabels.description) {
      await expect(element).toHaveAttribute('aria-describedby', expectedLabels.description);
    }

    if (expectedLabels.role) {
      await expect(element).toHaveAttribute('role', expectedLabels.role);
    }
  }

  /**
   * Verify focus management
   */
  async verifyFocusManagement(
    triggerSelector: string,
    expectedFocusSelector: string
  ): Promise<void> {
    const trigger = this.page.locator(triggerSelector);
    const expectedFocus = this.page.locator(expectedFocusSelector);

    await trigger.click();
    await expect(expectedFocus).toBeFocused();
  }

  /**
   * Verify screen reader announcements
   */
  async verifyScreenReaderAnnouncement(expectedMessage: string): Promise<void> {
    const liveRegion = this.page.locator('[aria-live="polite"], [aria-live="assertive"]');

    if (await liveRegion.isVisible()) {
      await expect(liveRegion).toContainText(expectedMessage);
    } else {
      console.warn('No live region found for screen reader announcements');
    }
  }

  /**
   * Verify form accessibility
   */
  async verifyFormAccessibility(formSelector: string): Promise<void> {
    const form = this.page.locator(formSelector);

    // Check form labels
    const inputs = form.locator('input, select, textarea');
    const inputCount = await inputs.count();

    for (let i = 0; i < inputCount; i++) {
      const input = inputs.nth(i);
      const hasLabel = await input.evaluate(el => {
        const labels = ['aria-label', 'aria-labelledby', 'title'];
        const hasAriaLabel = labels.some(attr => el.hasAttribute(attr));
        const hasHtmlLabel = document.querySelector(`label[for="${el.id}"]`);
        return hasAriaLabel || hasHtmlLabel;
      });

      expect(hasLabel).toBe(true, `Input at index ${i} in form ${formSelector} is missing label`);
    }

    // Check form validation
    await this.checkElementAccessibility(formSelector, {
      rules: {
        'label': { enabled: true },
        'form-field-multiple-labels': { enabled: true },
      },
    });
  }

  /**
   * Verify heading structure
   */
  async verifyHeadingStructure(): Promise<void> {
    const headings = this.page.locator('h1, h2, h3, h4, h5, h6');
    const headingCount = await headings.count();

    // Check for proper heading hierarchy
    let previousLevel = 0;
    for (let i = 0; i < headingCount; i++) {
      const heading = headings.nth(i);
      const tagName = await heading.evaluate(el => el.tagName.toLowerCase());
      const level = parseInt(tagName.substring(1));

      if (i === 0 && level !== 1) {
        console.warn(`Page should start with h1, but found ${tagName}`);
      }

      if (previousLevel > 0 && level > previousLevel + 1) {
        console.warn(`Heading skip detected: h${previousLevel} to ${tagName}`);
      }

      previousLevel = level;
    }
  }

  /**
   * Verify image accessibility
   */
  async verifyImageAccessibility(): Promise<void> {
    const images = this.page.locator('img');
    const imageCount = await images.count();

    for (let i = 0; i < imageCount; i++) {
      const image = images.nth(i);
      const alt = await image.getAttribute('alt');
      const role = await image.getAttribute('role');

      // Images should have alt text unless they're decorative
      if (role !== 'presentation') {
        expect(alt).toBeTruthy();
      }
    }
  }

  /**
   * Verify skip links
   */
  async verifySkipLinks(): Promise<void> {
    const skipLinks = this.page.locator('a[href^="#"], [data-testid="skip-link"]');
    const hasSkipLinks = await skipLinks.count() > 0;

    if (hasSkipLinks) {
      const skipLink = skipLinks.first();
      await expect(skipLink).toBeVisible();

      // Test skip link functionality
      const href = await skipLink.getAttribute('href');
      if (href) {
        const targetId = href.substring(1);
        const target = this.page.locator(`#${targetId}, [id="${targetId}"]`);
        await expect(target).toBeVisible();
      }
    }
  }

  /**
   * Generate accessibility report
   */
  async generateAccessibilityReport(): Promise<any> {
    const violations = await getViolations(this.page);

    return {
      timestamp: new Date().toISOString(),
      url: this.page.url(),
      violations: violations.map(v => ({
        id: v.id,
        impact: v.impact,
        description: v.description,
        help: v.help,
        helpUrl: v.helpUrl,
        nodes: v.nodes.map(n => ({
          html: n.html,
          target: n.target,
          failureSummary: n.failureSummary,
        })),
      })),
      summary: {
        total: violations.length,
        critical: violations.filter(v => v.impact === 'critical').length,
        serious: violations.filter(v => v.impact === 'serious').length,
        moderate: violations.filter(v => v.impact === 'moderate').length,
        minor: violations.filter(v => v.impact === 'minor').length,
      },
    };
  }
}
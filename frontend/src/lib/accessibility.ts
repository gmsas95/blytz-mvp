/**
 * Accessibility utilities for WCAG 2.1 AA compliance
 */

/**
 * Announce messages to screen readers
 */
export function announceMessage(message: string, priority: 'polite' | 'assertive' = 'polite') {
  const announcement = document.createElement('div');
  announcement.setAttribute('aria-live', priority);
  announcement.setAttribute('aria-atomic', 'true');
  announcement.className = 'sr-only';
  announcement.textContent = message;

  document.body.appendChild(announcement);

  // Remove after announcement
  setTimeout(() => {
    document.body.removeChild(announcement);
  }, 1000);
}

/**
 * Trap focus within a container element
 */
export function trapFocus(container: HTMLElement) {
  const focusableElements = container.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  ) as NodeListOf<HTMLElement>;

  if (focusableElements.length === 0) return;

  const firstFocusable = focusableElements[0];
  const lastFocusable = focusableElements[focusableElements.length - 1];

  const handleTabKey = (e: KeyboardEvent) => {
    if (e.key !== 'Tab') return;

    if (e.shiftKey) {
      if (document.activeElement === firstFocusable) {
        lastFocusable.focus();
        e.preventDefault();
      }
    } else {
      if (document.activeElement === lastFocusable) {
        firstFocusable.focus();
        e.preventDefault();
      }
    }
  };

  container.addEventListener('keydown', handleTabKey);

  // Return cleanup function
  return () => {
    container.removeEventListener('keydown', handleTabKey);
  };
}

/**
 * Set focus to first focusable element in container
 */
export function focusFirstElement(container: HTMLElement) {
  const focusableElements = container.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  ) as NodeListOf<HTMLElement>;

  if (focusableElements.length > 0) {
    focusableElements[0].focus();
    return true;
  }
  return false;
}

/**
 * Generate unique IDs for accessibility attributes
 */
let idCounter = 0;
export function generateId(prefix = 'a11y') {
  return `${prefix}-${++idCounter}`;
}

/**
 * Check if element is in viewport
 */
export function isInViewport(element: HTMLElement) {
  const rect = element.getBoundingClientRect();
  return (
    rect.top >= 0 &&
    rect.left >= 0 &&
    rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
    rect.right <= (window.innerWidth || document.documentElement.clientWidth)
  );
}

/**
 * Scroll element into view with focus management
 */
export function scrollIntoView(element: HTMLElement, options: ScrollIntoViewOptions = {}) {
  const defaultOptions: ScrollIntoViewOptions = {
    behavior: 'smooth',
    block: 'start',
    inline: 'nearest',
    ...options,
  };

  element.scrollIntoView(defaultOptions);

  // Focus after scroll completion
  setTimeout(() => {
    if (element.tabIndex === -1) {
      element.tabIndex = -1;
    }
    element.focus();
  }, 500);
}

/**
 * Manage keyboard navigation for dropdown menus
 */
export function setupDropdownNavigation(trigger: HTMLElement, menu: HTMLElement) {
  const menuItems = menu.querySelectorAll(
    'a, button, [role="menuitem"]'
  ) as NodeListOf<HTMLElement>;

  if (menuItems.length === 0) return;

  let currentIndex = -1;

  const handleKeyDown = (e: KeyboardEvent) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        currentIndex = (currentIndex + 1) % menuItems.length;
        menuItems[currentIndex].focus();
        break;

      case 'ArrowUp':
        e.preventDefault();
        currentIndex = currentIndex <= 0 ? menuItems.length - 1 : currentIndex - 1;
        menuItems[currentIndex].focus();
        break;

      case 'Home':
        e.preventDefault();
        currentIndex = 0;
        menuItems[currentIndex].focus();
        break;

      case 'End':
        e.preventDefault();
        currentIndex = menuItems.length - 1;
        menuItems[currentIndex].focus();
        break;

      case 'Escape':
        e.preventDefault();
        menu.style.display = 'none';
        trigger.focus();
        break;

      case 'Enter':
      case ' ':
        if (document.activeElement === trigger) {
          e.preventDefault();
          menu.style.display = menu.style.display === 'none' ? 'block' : 'none';
          if (menu.style.display !== 'none') {
            currentIndex = 0;
            menuItems[currentIndex].focus();
          }
        }
        break;
    }
  };

  const handleOutsideClick = (e: MouseEvent) => {
    if (!menu.contains(e.target as Node) && !trigger.contains(e.target as Node)) {
      menu.style.display = 'none';
    }
  };

  trigger.addEventListener('keydown', handleKeyDown);
  menu.addEventListener('keydown', handleKeyDown);
  document.addEventListener('click', handleOutsideClick);

  return () => {
    trigger.removeEventListener('keydown', handleKeyDown);
    menu.removeEventListener('keydown', handleKeyDown);
    document.removeEventListener('click', handleOutsideClick);
  };
}

/**
 * Setup skip links functionality
 */
export function setupSkipLinks() {
  const skipLinks = document.querySelectorAll('a[href^="#"]') as NodeListOf<HTMLAnchorElement>;

  skipLinks.forEach((link) => {
    link.addEventListener('click', (e) => {
      e.preventDefault();
      const targetId = link.getAttribute('href')?.slice(1);
      if (targetId) {
        const targetElement = document.getElementById(targetId);
        if (targetElement) {
          scrollIntoView(targetElement);
          announceMessage(`Navigated to ${targetElement.textContent || targetElement.id}`);
        }
      }
    });
  });
}

/**
 * Validate color contrast for WCAG compliance
 */
export function validateContrast(
  color1: string,
  color2: string
): { ratio: number; passes: boolean } {
  // Simple contrast ratio calculation
  // In production, you'd want to use a more sophisticated color parsing library
  const getLuminance = (color: string) => {
    // Simplified luminance calculation
    const rgb = parseInt(color.slice(1), 16);
    const r = (rgb >> 16) & 0xff;
    const g = (rgb >> 8) & 0xff;
    const b = (rgb >> 0) & 0xff;

    const rsRGB = r / 255;
    const gsRGB = g / 255;
    const bsRGB = b / 255;

    const rLinear = rsRGB <= 0.03928 ? rsRGB / 12.92 : Math.pow((rsRGB + 0.055) / 1.055, 2.4);
    const gLinear = gsRGB <= 0.03928 ? gsRGB / 12.92 : Math.pow((gsRGB + 0.055) / 1.055, 2.4);
    const bLinear = bsRGB <= 0.03928 ? bsRGB / 12.92 : Math.pow((bsRGB + 0.055) / 1.055, 2.4);

    return 0.2126 * rLinear + 0.7152 * gLinear + 0.0722 * bLinear;
  };

  const lum1 = getLuminance(color1);
  const lum2 = getLuminance(color2);

  const brightest = Math.max(lum1, lum2);
  const darkest = Math.min(lum1, lum2);

  const ratio = (brightest + 0.05) / (darkest + 0.05);

  return {
    ratio,
    passes: ratio >= 4.5, // WCAG AA standard
  };
}

/**
 * Initialize accessibility features
 */
export function initializeAccessibility() {
  setupSkipLinks();

  // Add global keyboard navigation
  document.addEventListener('keydown', (e) => {
    // Skip to main content on Alt+M
    if (e.altKey && e.key === 'm') {
      const mainContent = document.getElementById('main-content');
      if (mainContent) {
        scrollIntoView(mainContent);
        announceMessage('Skipped to main content');
      }
    }
  });

  // Announce page changes
  announceMessage(`Page loaded: ${  document.title}`, 'polite');
}

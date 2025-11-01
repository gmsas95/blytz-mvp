import { formatPrice, cn, formatTimeRemaining, debounce, isMobile } from '../utils';

describe('Utils', () => {
  describe('formatPrice', () => {
    it('formats positive numbers correctly', () => {
      expect(formatPrice(100)).toBe('$100.00');
      expect(formatPrice(100.5)).toBe('$100.50');
      expect(formatPrice(100.99)).toBe('$100.99');
    });

    it('formats zero correctly', () => {
      expect(formatPrice(0)).toBe('$0.00');
    });

    it('formats large numbers correctly', () => {
      expect(formatPrice(1000000)).toBe('$1,000,000.00');
    });

    it('handles negative numbers', () => {
      expect(formatPrice(-100)).toBe('-$100.00');
    });
  });

  describe('cn', () => {
    it('merges class names correctly', () => {
      expect(cn('foo', 'bar')).toBe('foo bar');
    });

    it('filters out falsy values', () => {
      expect(cn('foo', null, 'bar', undefined, false, 'baz')).toBe('foo bar baz');
    });

    it('handles empty input', () => {
      expect(cn()).toBe('');
      expect(cn('', null, undefined, false)).toBe('');
    });

    it('handles conditional classes', () => {
      expect(cn('base', true && 'active', false && 'inactive')).toBe('base active');
    });
  });

  describe('formatTimeRemaining', () => {
    it('formats seconds correctly', () => {
      expect(formatTimeRemaining(30)).toBe('30s');
      expect(formatTimeRemaining(59)).toBe('59s');
    });

    it('formats minutes correctly', () => {
      expect(formatTimeRemaining(60)).toBe('1h 0m');
      expect(formatTimeRemaining(90)).toBe('1h 1m');
      expect(formatTimeRemaining(3599)).toBe('59m 59s');
    });

    it('formats hours correctly', () => {
      expect(formatTimeRemaining(3600)).toBe('1h 0m');
      expect(formatTimeRemaining(3661)).toBe('1h 1m');
      expect(formatTimeRemaining(86399)).toBe('23h 59m');
    });

    it('handles zero time', () => {
      expect(formatTimeRemaining(0)).toBe('Ended');
    });

    it('handles negative time', () => {
      expect(formatTimeRemaining(-10)).toBe('Ended');
    });
  });

  describe('debounce', () => {
    jest.useFakeTimers();

    it('delays function execution', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn, 1000);

      debouncedFn();
      expect(mockFn).not.toHaveBeenCalled();

      jest.advanceTimersByTime(1000);
      expect(mockFn).toHaveBeenCalledTimes(1);
    });

    it('cancels previous calls', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn, 1000);

      debouncedFn();
      debouncedFn();
      debouncedFn();

      jest.advanceTimersByTime(1000);
      expect(mockFn).toHaveBeenCalledTimes(1);
    });

    it('passes arguments correctly', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn, 1000);

      debouncedFn('arg1', 'arg2');
      jest.advanceTimersByTime(1000);

      expect(mockFn).toHaveBeenCalledWith('arg1', 'arg2');
    });

    afterEach(() => {
      jest.clearAllTimers();
    });
  });

  describe('isMobile', () => {
    const mockWindow = Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
    });

    it('returns true for mobile screen width', () => {
      mockWindow.value = 500;
      expect(isMobile()).toBe(true);
    });

    it('returns false for desktop screen width', () => {
      mockWindow.value = 1024;
      expect(isMobile()).toBe(false);
    });

    it('returns false on server side', () => {
      const originalWindow = global.window;
      delete global.window;
      expect(isMobile()).toBe(false);
      global.window = originalWindow;
    });
  });
});

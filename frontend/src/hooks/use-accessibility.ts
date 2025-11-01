import { useEffect, useRef, useCallback } from 'react';

import { announceMessage, trapFocus, focusFirstElement, generateId } from '@/lib/accessibility';

/**
 * Hook for managing screen reader announcements
 */
export function useAnnouncement() {
  const announce = useCallback((message: string, priority?: 'polite' | 'assertive') => {
    announceMessage(message, priority);
  }, []);

  return { announce };
}

/**
 * Hook for managing focus within modals and dialogs
 */
export function useFocusTrap(isActive: boolean) {
  const containerRef = useRef<HTMLElement>(null);
  const cleanupRef = useRef<(() => void) | null>(null);

  useEffect(() => {
    if (isActive && containerRef.current) {
      // Focus first element when activated
      focusFirstElement(containerRef.current);

      // Setup focus trap
      cleanupRef.current = trapFocus(containerRef.current);
    }

    return () => {
      if (cleanupRef.current) {
        cleanupRef.current();
      }
    };
  }, [isActive]);

  return containerRef;
}

/**
 * Hook for managing keyboard navigation
 */
export function useKeyboardNavigation(
  items: HTMLElement[],
  options: {
    orientation?: 'horizontal' | 'vertical';
    loop?: boolean;
    onActivate?: (index: number, element: HTMLElement) => void;
  } = {}
) {
  const { orientation = 'vertical', loop = true, onActivate } = options;
  const activeIndexRef = useRef(-1);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      const isVertical = orientation === 'vertical';
      const nextKey = isVertical ? 'ArrowDown' : 'ArrowRight';
      const prevKey = isVertical ? 'ArrowUp' : 'ArrowLeft';

      switch (e.key) {
        case nextKey:
          e.preventDefault();
          activeIndexRef.current = (activeIndexRef.current + 1) % items.length;
          items[activeIndexRef.current]?.focus();
          onActivate?.(activeIndexRef.current, items[activeIndexRef.current]);
          break;

        case prevKey:
          e.preventDefault();
          activeIndexRef.current =
            activeIndexRef.current <= 0
              ? loop
                ? items.length - 1
                : 0
              : activeIndexRef.current - 1;
          items[activeIndexRef.current]?.focus();
          onActivate?.(activeIndexRef.current, items[activeIndexRef.current]);
          break;

        case 'Home':
          e.preventDefault();
          activeIndexRef.current = 0;
          items[0]?.focus();
          onActivate?.(0, items[0]);
          break;

        case 'End':
          e.preventDefault();
          activeIndexRef.current = items.length - 1;
          items[items.length - 1]?.focus();
          onActivate?.(items.length - 1, items[items.length - 1]);
          break;
      }
    },
    [items, orientation, loop, onActivate]
  );

  useEffect(() => {
    if (items.length > 0) {
      document.addEventListener('keydown', handleKeyDown);
      return () => document.removeEventListener('keydown', handleKeyDown);
    }
  }, [handleKeyDown, items.length]);

  return {
    activeIndex: activeIndexRef.current,
    setActiveIndex: (index: number) => {
      if (index >= 0 && index < items.length) {
        activeIndexRef.current = index;
        items[index]?.focus();
        onActivate?.(index, items[index]);
      }
    },
  };
}

/**
 * Hook for generating accessible IDs
 */
export function useA11yId(prefix = 'a11y') {
  const idRef = useRef<string>(generateId(prefix));

  return idRef.current;
}

/**
 * Hook for managing form accessibility
 */
export function useFormA11y(
  options: {
    onSubmit?: () => void;
    onError?: (errors: Record<string, string>) => void;
  } = {}
) {
  const { announce } = useAnnouncement();
  const formRef = useRef<HTMLFormElement>(null);

  const announceFormError = useCallback(
    (errors: Record<string, string>) => {
      const errorMessages = Object.values(errors);
      if (errorMessages.length > 0) {
        announce(
          `Form has ${errorMessages.length} error${errorMessages.length > 1 ? 's' : ''}. ${errorMessages.join('. ')}`,
          'assertive'
        );
        options.onError?.(errors);
      }
    },
    [announce, options.onError]
  );

  const announceFormSuccess = useCallback(() => {
    announce('Form submitted successfully', 'polite');
    options.onSubmit?.();
  }, [announce, options.onSubmit]);

  const focusFirstError = useCallback(() => {
    if (formRef.current) {
      const firstError = formRef.current.querySelector('[aria-invalid="true"]') as HTMLElement;
      if (firstError) {
        firstError.focus();
        firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    }
  }, []);

  return {
    formRef,
    announceFormError,
    announceFormSuccess,
    focusFirstError,
  };
}

/**
 * Hook for managing live regions
 */
export function useLiveRegion(priority: 'polite' | 'assertive' = 'polite') {
  const regionRef = useRef<HTMLDivElement>(null);
  const idRef = useRef<string>(generateId('live-region'));

  useEffect(() => {
    if (regionRef.current) {
      regionRef.current.setAttribute('aria-live', priority);
      regionRef.current.setAttribute('aria-atomic', 'true');
    }
  }, [priority]);

  const announce = useCallback((message: string) => {
    if (regionRef.current) {
      regionRef.current.textContent = message;
    }
  }, []);

  return {
    regionRef,
    id: idRef.current,
    announce,
  };
}

/**
 * Hook for managing carousel/slider accessibility
 */
export function useCarouselAccessibility(totalItems: number) {
  const { announce } = useAnnouncement();
  const currentIndexRef = useRef(0);

  const goToSlide = useCallback(
    (index: number) => {
      if (index >= 0 && index < totalItems) {
        currentIndexRef.current = index;
        announce(`Slide ${index + 1} of ${totalItems}`);
      }
    },
    [totalItems, announce]
  );

  const nextSlide = useCallback(() => {
    const nextIndex = (currentIndexRef.current + 1) % totalItems;
    goToSlide(nextIndex);
  }, [totalItems, goToSlide]);

  const prevSlide = useCallback(() => {
    const prevIndex = currentIndexRef.current === 0 ? totalItems - 1 : currentIndexRef.current - 1;
    goToSlide(prevIndex);
  }, [totalItems, goToSlide]);

  return {
    currentIndex: currentIndexRef.current,
    goToSlide,
    nextSlide,
    prevSlide,
  };
}

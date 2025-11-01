'use client';

import { useEffect } from 'react';

import { initializeAccessibility } from '@/lib/accessibility';

export function AccessibilityInit() {
  useEffect(() => {
    // Initialize accessibility features
    initializeAccessibility();
  }, []);

  return null;
}

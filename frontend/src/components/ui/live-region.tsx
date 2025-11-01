import * as React from 'react';

import { cn } from '@/lib/utils';

interface LiveRegionProps extends React.HTMLAttributes<HTMLDivElement> {
  priority?: 'polite' | 'assertive' | 'off';
  atomic?: boolean;
  busy?: boolean;
  clearOnUpdate?: boolean;
}

const LiveRegion = React.forwardRef<HTMLDivElement, LiveRegionProps>(
  (
    {
      className,
      priority = 'polite',
      atomic = true,
      busy = false,
      clearOnUpdate = false,
      children,
      ...props
    },
    ref
  ) => {
    const regionRef = React.useRef<HTMLDivElement>(null);

    React.useImperativeHandle(ref, () => regionRef.current!);

    React.useEffect(() => {
      if (regionRef.current) {
        // Clear content if requested on update
        if (clearOnUpdate && children) {
          regionRef.current.textContent = '';
        }
      }
    }, [children, clearOnUpdate]);

    return (
      <div
        ref={regionRef}
        className={cn('sr-only', className)}
        aria-live={priority}
        aria-atomic={atomic}
        aria-busy={busy}
        {...props}
      >
        {children}
      </div>
    );
  }
);
LiveRegion.displayName = 'LiveRegion';

export { LiveRegion };

import * as React from 'react';

import { cn } from '@/lib/utils';

interface LoadingSpinnerProps extends React.HTMLAttributes<HTMLDivElement> {
  size?: 'sm' | 'md' | 'lg';
  label?: string;
}

const LoadingSpinner = React.forwardRef<HTMLDivElement, LoadingSpinnerProps>(
  ({ className, size = 'md', label = 'Loading...', ...props }, ref) => {
    const sizeClasses = {
      sm: 'w-4 h-4',
      md: 'w-6 h-6',
      lg: 'w-8 h-8',
    };

    return (
      <div
        ref={ref}
        className={cn('flex items-center space-x-2', className)}
        role="status"
        aria-label={label}
        {...props}
      >
        <div
          className={cn(
            'animate-spin rounded-full border-2 border-primary border-t-transparent',
            sizeClasses[size]
          )}
          aria-hidden="true"
        />
        <span className="sr-only">{label}</span>
      </div>
    );
  }
);
LoadingSpinner.displayName = 'LoadingSpinner';

export { LoadingSpinner };

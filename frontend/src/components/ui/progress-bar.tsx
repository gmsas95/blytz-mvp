import * as ProgressPrimitive from '@radix-ui/react-progress';
import * as React from 'react';

import { cn } from '@/lib/utils';

const Progress = React.forwardRef<
  React.ElementRef<typeof ProgressPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root> & {
    label?: string;
    showValue?: boolean;
  }
>(({ className, value, label, showValue = true, ...props }, ref) => {
  const percentage = value ? Math.min(Math.max((value / (props.max || 100)) * 100, 0), 100) : 0;

  return (
    <div className={cn('space-y-2', className)}>
      {label && (
        <div className="flex justify-between items-center text-sm">
          <span className="font-medium">{label}</span>
          {showValue && (
            <span className="text-muted-foreground" aria-live="polite">
              {Math.round(percentage)}%
            </span>
          )}
        </div>
      )}
      <ProgressPrimitive.Root
        ref={ref}
        className={cn('relative h-4 w-full overflow-hidden rounded-full bg-secondary', className)}
        {...props}
      >
        <ProgressPrimitive.Indicator
          className="h-full w-full flex-1 bg-primary transition-all duration-300 ease-out"
          style={{ transform: `translateX(-${100 - percentage}%)` }}
          aria-label={
            label
              ? `${label}: ${Math.round(percentage)}% complete`
              : `${Math.round(percentage)}% complete`
          }
        />
      </ProgressPrimitive.Root>
    </div>
  );
});
Progress.displayName = ProgressPrimitive.Root.displayName;

export { Progress };

'use client';

import * as ToastPrimitives from '@radix-ui/react-toast';
import { cva, type VariantProps } from 'class-variance-authority';
import { X } from 'lucide-react';
import * as React from 'react';

import { cn } from '@/lib/utils';

const ToastProviderPrimitive = ToastPrimitives.Provider;

const ToastViewport = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Viewport>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Viewport>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Viewport
    ref={ref}
    className={cn(
      'fixed top-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:bottom-0 sm:right-0 sm:top-auto sm:flex-col md:max-w-[420px]',
      className
    )}
    {...props}
  />
));
ToastViewport.displayName = ToastPrimitives.Viewport.displayName;

const toastVariants = cva(
  'group pointer-events-auto relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-6 pr-8 shadow-lg transition-all data-[swipe=cancel]:translate-x-0 data-[swipe=end]:translate-x-[var(--radix-toast-swipe-end-x)] data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)] data-[swipe=move]:transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[swipe=end]:animate-out data-[state=closed]:fade-out-80 data-[state=closed]:slide-out-to-right-full data-[state=open]:slide-in-from-top-full data-[state=open]:sm:slide-in-from-bottom-full',
  {
    variants: {
      variant: {
        default: 'border bg-background text-foreground',
        destructive: 'destructive border-destructive bg-destructive text-destructive-foreground',
        success: 'border-green-200 bg-green-50 text-green-800',
        warning: 'border-yellow-200 bg-yellow-50 text-yellow-800',
        info: 'border-blue-200 bg-blue-50 text-blue-800',
        bid: 'border-purple-200 bg-purple-50 text-purple-800 animate-pulse',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
);

const Toast = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Root>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Root> & VariantProps<typeof toastVariants>
>(({ className, variant, ...props }, ref) => {
  return (
    <ToastPrimitives.Root
      ref={ref}
      className={cn(toastVariants({ variant }), className)}
      {...props}
    />
  );
});
Toast.displayName = ToastPrimitives.Root.displayName;

const ToastAction = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Action>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Action>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Action
    ref={ref}
    className={cn(
      'inline-flex h-8 shrink-0 items-center justify-center rounded-md border bg-transparent px-3 text-sm font-medium ring-offset-background transition-colors hover:bg-secondary focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 group-[.destructive]:border-muted/40 group-[.destructive]:hover:border-destructive/30 group-[.destructive]:hover:bg-destructive group-[.destructive]:hover:text-destructive-foreground group-[.destructive]:focus:ring-destructive',
      className
    )}
    {...props}
  />
));
ToastAction.displayName = ToastPrimitives.Action.displayName;

const ToastClose = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Close>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Close>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Close
    ref={ref}
    className={cn(
      'absolute right-2 top-2 rounded-md p-1 text-foreground/50 opacity-0 transition-opacity hover:text-foreground focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100 group-[.destructive]:text-red-300 group-[.destructive]:hover:text-red-50 group-[.destructive]:focus:ring-red-400 group-[.destructive]:focus:ring-offset-red-600',
      className
    )}
    toast-close=""
    {...props}
  >
    <X className="h-4 w-4" />
  </ToastPrimitives.Close>
));
ToastClose.displayName = ToastPrimitives.Close.displayName;

const ToastTitle = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Title>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Title>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Title ref={ref} className={cn('text-sm font-semibold', className)} {...props} />
));
ToastTitle.displayName = ToastPrimitives.Title.displayName;

const ToastDescription = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Description>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Description>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Description
    ref={ref}
    className={cn('text-sm opacity-90', className)}
    {...props}
  />
));
ToastDescription.displayName = ToastPrimitives.Description.displayName;

type ToastProps = React.ComponentPropsWithoutRef<typeof Toast>;

type ToastActionElement = React.ReactElement<typeof ToastAction>;

// Toast context for managing toasts
interface ToastContextType {
  toasts: Toast[];
  addToast: (toast: Omit<Toast, 'id'>) => void;
  removeToast: (id: string) => void;
  clearToasts: () => void;
}

interface Toast {
  id: string;
  title?: string;
  description?: string;
  variant?: 'default' | 'destructive' | 'success' | 'warning' | 'info' | 'bid';
  duration?: number;
  action?: ToastActionElement;
  onDismiss?: () => void;
}

const ToastContext = React.createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = React.useState<Toast[]>([]);

  const addToast = React.useCallback((toast: Omit<Toast, 'id'>) => {
    const id = Math.random().toString(36).substr(2, 9);
    const newToast = { ...toast, id };

    setToasts((prev) => [...prev, newToast]);

    // Auto-remove after duration
    const duration = toast.duration ?? 5000;
    if (duration > 0) {
      setTimeout(() => {
        removeToast(id);
      }, duration);
    }

    return id;
  }, []);

  const removeToast = React.useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  const clearToasts = React.useCallback(() => {
    setToasts([]);
  }, []);

  const contextValue = React.useMemo(() => ({
    toasts,
    addToast,
    removeToast,
    clearToasts
  }), [toasts, addToast, removeToast, clearToasts]);

  // Set global context for convenience functions
  React.useEffect(() => {
    setGlobalToastContext(contextValue);
    return () => setGlobalToastContext(null);
  }, [contextValue]);

  // Only render ToastViewport on client side
  const [isClient, setIsClient] = React.useState(false);

  React.useEffect(() => {
    setIsClient(true);
  }, []);

  return (
    <ToastContext.Provider value={contextValue}>
      {children}
      {isClient && (
        <ToastViewport>
          {toasts.map((toast) => (
            <Toast
              key={toast.id}
              variant={toast.variant}
              onOpenChange={(open) => !open && removeToast(toast.id)}
            >
              <div className="grid gap-1">
                {toast.title && <ToastTitle>{toast.title}</ToastTitle>}
                {toast.description && <ToastDescription>{toast.description}</ToastDescription>}
              </div>
              {toast.action}
              <ToastClose />
            </Toast>
          ))}
        </ToastViewport>
      )}
    </ToastContext.Provider>
  );
}

export function useToast() {
  const context = React.useContext(ToastContext);
  if (!context) {
    // Return a safe fallback for SSR cases
    if (typeof window === 'undefined') {
      return {
        toasts: [],
        addToast: () => '',
        removeToast: () => {},
        clearToasts: () => {},
        success: () => '',
        error: () => '',
        warning: () => '',
        info: () => '',
        bid: () => '',
      };
    }
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
}

// Global toast state for convenience functions
let globalToastContext: ToastContextType | null = null;

// Function to set the global toast context (called by ToastProvider)
export const setGlobalToastContext = (context: ToastContextType) => {
  globalToastContext = context;
};

// Convenience functions for common toast types
export const toast = {
  success: (
    title: string,
    description?: string,
    options?: Partial<Omit<Toast, 'id' | 'title' | 'description' | 'variant'>>
  ) => {
    if (!globalToastContext) {
      console.warn('Toast context not available. Make sure ToastProvider is mounted.');
      return null;
    }
    return globalToastContext.addToast({ title, description, variant: 'success', ...options });
  },
  error: (
    title: string,
    description?: string,
    options?: Partial<Omit<Toast, 'id' | 'title' | 'description' | 'variant'>>
  ) => {
    if (!globalToastContext) {
      console.warn('Toast context not available. Make sure ToastProvider is mounted.');
      return null;
    }
    return globalToastContext.addToast({ title, description, variant: 'destructive', ...options });
  },
  warning: (
    title: string,
    description?: string,
    options?: Partial<Omit<Toast, 'id' | 'title' | 'description' | 'variant'>>
  ) => {
    if (!globalToastContext) {
      console.warn('Toast context not available. Make sure ToastProvider is mounted.');
      return null;
    }
    return globalToastContext.addToast({ title, description, variant: 'warning', ...options });
  },
  info: (
    title: string,
    description?: string,
    options?: Partial<Omit<Toast, 'id' | 'title' | 'description' | 'variant'>>
  ) => {
    if (!globalToastContext) {
      console.warn('Toast context not available. Make sure ToastProvider is mounted.');
      return null;
    }
    return globalToastContext.addToast({ title, description, variant: 'info', ...options });
  },
  bid: (
    title: string,
    description?: string,
    options?: Partial<Omit<Toast, 'id' | 'title' | 'description' | 'variant'>>
  ) => {
    if (!globalToastContext) {
      console.warn('Toast context not available. Make sure ToastProvider is mounted.');
      return null;
    }
    return globalToastContext.addToast({ title, description, variant: 'bid', duration: 8000, ...options });
  },
};

export {
  type ToastProps,
  type ToastActionElement,
  Toast,
  ToastClose,
  ToastTitle,
  ToastDescription,
  ToastViewport,
  ToastProviderPrimitive,
};

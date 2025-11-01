import * as React from 'react';

import { cn } from '@/lib/utils';

interface AccessibleImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  fallback?: string;
  loadingStrategy?: 'lazy' | 'eager';
  onLoad?: () => void;
  onError?: () => void;
}

const AccessibleImage = React.forwardRef<HTMLImageElement, AccessibleImageProps>(
  ({ className, alt, src, fallback, loadingStrategy = 'lazy', onLoad, onError, ...props }, ref) => {
    const [imgSrc, setImgSrc] = React.useState(src);
    const [isLoading, setIsLoading] = React.useState(true);
    const [hasError, setHasError] = React.useState(false);

    React.useEffect(() => {
      setImgSrc(src);
      setIsLoading(true);
      setHasError(false);
    }, [src]);

    const handleLoad = React.useCallback(() => {
      setIsLoading(false);
      onLoad?.();
    }, [onLoad]);

    const handleError = React.useCallback(() => {
      setIsLoading(false);
      setHasError(true);
      if (fallback) {
        setImgSrc(fallback);
      }
      onError?.();
    }, [fallback, onError]);

    // If no alt text provided, it's a decorative image
    const isDecorative = !alt || alt.trim() === '';
    const altText = isDecorative ? '' : alt;

    return (
      <div className={cn('relative overflow-hidden', className)}>
        {isLoading && (
          <div
            className="absolute inset-0 bg-muted animate-pulse flex items-center justify-center"
            aria-hidden="true"
          >
            <div className="w-8 h-8 border-2 border-muted-foreground border-t-transparent rounded-full animate-spin" />
          </div>
        )}

        {hasError && !fallback && (
          <div
            className="absolute inset-0 bg-muted flex items-center justify-center text-muted-foreground"
            role="img"
            aria-label={altText || 'Image failed to load'}
          >
            <svg
              className="w-8 h-8"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        )}

        <img
          ref={ref}
          src={imgSrc}
          alt={altText}
          className={cn(
            'w-full h-full object-cover transition-opacity duration-300',
            isLoading ? 'opacity-0' : 'opacity-100'
          )}
          loading={loadingStrategy}
          onLoad={handleLoad}
          onError={handleError}
          role={isDecorative ? 'presentation' : undefined}
          {...props}
        />
      </div>
    );
  }
);

AccessibleImage.displayName = 'AccessibleImage';

export { AccessibleImage };

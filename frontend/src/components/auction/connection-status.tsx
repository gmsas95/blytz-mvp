'use client';

import {
  Wifi,
  WifiOff,
  Loader2,
  AlertTriangle,
  RefreshCw,
  SignalHigh,
  SignalMedium,
  SignalLow,
} from 'lucide-react';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface ConnectionStatusProps {
  connected: boolean;
  connecting: boolean;
  reconnectAttempts: number;
  error?: string;
  connectionQuality: 'excellent' | 'good' | 'poor' | 'disconnected';
  isUsingFallback: boolean;
  onReconnect?: () => void;
  className?: string;
  variant?: 'badge' | 'alert' | 'detailed';
}

export function ConnectionStatus({
  connected,
  connecting,
  reconnectAttempts,
  error,
  connectionQuality,
  isUsingFallback,
  onReconnect,
  className,
  variant = 'badge',
}: ConnectionStatusProps) {
  const getStatusIcon = () => {
    if (connecting) {
      return <Loader2 className="h-3 w-3 animate-spin" />;
    }

    if (!connected) {
      return <WifiOff className="h-3 w-3" />;
    }

    switch (connectionQuality) {
      case 'excellent':
        return <SignalHigh className="h-3 w-3" />;
      case 'good':
        return <SignalMedium className="h-3 w-3" />;
      case 'poor':
        return <SignalLow className="h-3 w-3" />;
      default:
        return <Wifi className="h-3 w-3" />;
    }
  };

  const getStatusText = () => {
    if (connecting) return 'Connecting...';
    if (!connected) return 'Disconnected';
    if (isUsingFallback) return 'Polling Mode';
    return 'Live';
  };

  const getStatusColor = () => {
    if (connecting) return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    if (!connected) return 'bg-red-100 text-red-800 border-red-200';
    if (isUsingFallback) return 'bg-orange-100 text-orange-800 border-orange-200';

    switch (connectionQuality) {
      case 'excellent':
        return 'bg-green-100 text-green-800 border-green-200';
      case 'good':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'poor':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const shouldShowWarning = () => {
    return !connected || isUsingFallback || connectionQuality === 'poor' || reconnectAttempts > 3;
  };

  if (variant === 'detailed') {
    return (
      <div className={cn('space-y-3', className)}>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Badge className={cn('flex items-center gap-1', getStatusColor())}>
              {getStatusIcon()}
              <span className="text-xs font-medium">{getStatusText()}</span>
            </Badge>
            {reconnectAttempts > 0 && (
              <Badge variant="outline" className="text-xs">
                Attempt {reconnectAttempts}
              </Badge>
            )}
          </div>
          {!connected && onReconnect && (
            <Button size="sm" variant="outline" onClick={onReconnect} disabled={connecting}>
              <RefreshCw className={cn('h-3 w-3 mr-1', connecting && 'animate-spin')} />
              Reconnect
            </Button>
          )}
        </div>

        {shouldShowWarning() && (
          <Alert
            className={cn(
              'border-l-4',
              !connected && 'border-red-500 bg-red-50',
              isUsingFallback && 'border-orange-500 bg-orange-50',
              connectionQuality === 'poor' && 'border-yellow-500 bg-yellow-50'
            )}
          >
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription className="text-sm">
              {!connected && 'Real-time connection lost. Some features may be delayed.'}
              {isUsingFallback && 'Using fallback polling mode. Updates may be delayed.'}
              {connected &&
                connectionQuality === 'poor' &&
                'Poor connection quality. Updates may be delayed.'}
              {reconnectAttempts > 3 &&
                `Connection unstable. ${reconnectAttempts} reconnection attempts made.`}
            </AlertDescription>
          </Alert>
        )}

        {error && (
          <Alert className="border-red-200 bg-red-50">
            <AlertTriangle className="h-4 w-4 text-red-600" />
            <AlertDescription className="text-red-800 text-sm">{error}</AlertDescription>
          </Alert>
        )}

        <div className="text-xs text-muted-foreground">
          {connected && (
            <div className="flex items-center gap-4">
              <span>Quality: {connectionQuality}</span>
              <span>Status: {getStatusText()}</span>
              {isUsingFallback && <span>Mode: Polling</span>}
            </div>
          )}
        </div>
      </div>
    );
  }

  if (variant === 'alert') {
    if (shouldShowWarning()) {
      return (
        <Alert
          className={cn(
            'border-l-4',
            !connected && 'border-red-500 bg-red-50',
            isUsingFallback && 'border-orange-500 bg-orange-50',
            connectionQuality === 'poor' && 'border-yellow-500 bg-yellow-50',
            className
          )}
        >
          <div className="flex items-center gap-2">
            {getStatusIcon()}
            <AlertDescription className="text-sm">
              {!connected && 'Real-time connection lost. '}
              {isUsingFallback && 'Using fallback polling mode. '}
              {connected && connectionQuality === 'poor' && 'Poor connection quality. '}
              Updates may be delayed.
              {!connected && onReconnect && (
                <Button
                  size="sm"
                  variant="link"
                  onClick={onReconnect}
                  disabled={connecting}
                  className="ml-2 h-auto p-0 text-sm"
                >
                  Reconnect
                </Button>
              )}
            </AlertDescription>
          </div>
        </Alert>
      );
    }
    return null;
  }

  // Default badge variant
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Badge className={cn('flex items-center gap-1', getStatusColor())}>
        {getStatusIcon()}
        <span className="text-xs font-medium">{getStatusText()}</span>
      </Badge>

      {(isUsingFallback || reconnectAttempts > 0) && (
        <Badge variant="outline" className="text-xs">
          {isUsingFallback ? 'Polling' : `Retry ${reconnectAttempts}`}
        </Badge>
      )}

      {shouldShowWarning() && <AlertTriangle className="h-3 w-3 text-yellow-600" />}
    </div>
  );
}

// Mini connection indicator for inline use
export function ConnectionIndicator({
  connected,
  connecting,
  connectionQuality,
  className,
}: {
  connected: boolean;
  connecting: boolean;
  connectionQuality: 'excellent' | 'good' | 'poor' | 'disconnected';
  className?: string;
}) {
  if (connecting) {
    return (
      <div className={cn('flex items-center gap-1', className)}>
        <Loader2 className="h-3 w-3 animate-spin text-yellow-600" />
        <span className="text-xs text-yellow-600">Connecting...</span>
      </div>
    );
  }

  if (!connected) {
    return (
      <div className={cn('flex items-center gap-1', className)}>
        <WifiOff className="h-3 w-3 text-red-600" />
        <span className="text-xs text-red-600">Offline</span>
      </div>
    );
  }

  const qualityColors = {
    excellent: 'text-green-600',
    good: 'text-blue-600',
    poor: 'text-yellow-600',
    disconnected: 'text-red-600',
  };

  return (
    <div className={cn('flex items-center gap-1', className)}>
      <Wifi className={cn('h-3 w-3', qualityColors[connectionQuality])} />
      <span className={cn('text-xs', qualityColors[connectionQuality])}>
        {connectionQuality === 'excellent' && 'Live'}
        {connectionQuality === 'good' && 'Connected'}
        {connectionQuality === 'poor' && 'Unstable'}
        {connectionQuality === 'disconnected' && 'Offline'}
      </span>
    </div>
  );
}

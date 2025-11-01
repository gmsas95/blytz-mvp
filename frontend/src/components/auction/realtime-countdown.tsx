'use client';

import { useEffect, useState, useCallback } from 'react';

import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface RealtimeCountdownProps {
  endTime: string | Date;
  onEnd?: () => void;
  onEndingSoon?: (remainingTime: number) => void;
  className?: string;
  showLabel?: boolean;
  variant?: 'default' | 'compact' | 'detailed';
  isActive?: boolean;
}

interface TimeLeft {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  total: number;
}

export function RealtimeCountdown({
  endTime,
  onEnd,
  onEndingSoon,
  className,
  showLabel = true,
  variant = 'default',
  isActive = true,
}: RealtimeCountdownProps) {
  const [timeLeft, setTimeLeft] = useState<TimeLeft>({
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0,
    total: 0,
  });
  const [isEndingSoon, setIsEndingSoon] = useState(false);
  const [hasEnded, setHasEnded] = useState(false);

  const calculateTimeLeft = useCallback((): TimeLeft => {
    const difference = new Date(endTime).getTime() - new Date().getTime();

    if (difference <= 0) {
      return {
        days: 0,
        hours: 0,
        minutes: 0,
        seconds: 0,
        total: 0,
      };
    }

    return {
      days: Math.floor(difference / (1000 * 60 * 60 * 24)),
      hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
      minutes: Math.floor((difference / 1000 / 60) % 60),
      seconds: Math.floor((difference / 1000) % 60),
      total: Math.floor(difference / 1000),
    };
  }, [endTime]);

  useEffect(() => {
    if (!isActive) return;

    const timer = setInterval(() => {
      const newTimeLeft = calculateTimeLeft();
      setTimeLeft(newTimeLeft);

      // Check if auction is ending soon (less than 5 minutes)
      const endingSoon = newTimeLeft.total > 0 && newTimeLeft.total < 300;
      if (endingSoon && !isEndingSoon) {
        setIsEndingSoon(true);
        onEndingSoon?.(newTimeLeft.total);
      }

      // Check if auction has ended
      if (newTimeLeft.total === 0 && !hasEnded) {
        setHasEnded(true);
        onEnd?.();
      }
    }, 1000);

    // Initial calculation
    setTimeLeft(calculateTimeLeft());

    return () => clearInterval(timer);
  }, [isActive, calculateTimeLeft, isEndingSoon, hasEnded, onEnd, onEndingSoon]);

  const formatTime = (time: TimeLeft): string => {
    if (time.total === 0) return 'Ended';

    if (time.days > 0) {
      return `${time.days}d ${time.hours}h`;
    } else if (time.hours > 0) {
      return `${time.hours}h ${time.minutes}m`;
    } else if (time.minutes > 0) {
      return `${time.minutes}m ${time.seconds}s`;
    } else {
      return `${time.seconds}s`;
    }
  };

  const getVariantStyles = () => {
    switch (variant) {
      case 'compact':
        return 'text-xs px-2 py-1';
      case 'detailed':
        return 'text-sm px-3 py-2';
      default:
        return 'text-sm px-2 py-1';
    }
  };

  const getStatusColor = () => {
    if (hasEnded) return 'bg-gray-100 text-gray-800';
    if (isEndingSoon) return 'bg-orange-100 text-orange-800 animate-pulse';
    if (timeLeft.total < 3600) return 'bg-yellow-100 text-yellow-800'; // Less than 1 hour
    return 'bg-green-100 text-green-800';
  };

  const getStatusText = () => {
    if (hasEnded) return 'ENDED';
    if (isEndingSoon) return 'ENDING SOON';
    if (timeLeft.total < 3600) return 'ENDING TODAY';
    return 'LIVE';
  };

  if (variant === 'detailed') {
    return (
      <div className={cn('border rounded-lg p-3 bg-card', className)}>
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm font-medium">Auction Ends</span>
          <Badge className={getVariantStyles()}>{getStatusText()}</Badge>
        </div>

        <div className="grid grid-cols-4 gap-2 text-center">
          {timeLeft.days > 0 && (
            <div>
              <div className="text-2xl font-bold">{timeLeft.days}</div>
              <div className="text-xs text-muted-foreground">Days</div>
            </div>
          )}
          <div>
            <div className="text-2xl font-bold">{timeLeft.hours}</div>
            <div className="text-xs text-muted-foreground">Hours</div>
          </div>
          <div>
            <div className="text-2xl font-bold">{timeLeft.minutes}</div>
            <div className="text-xs text-muted-foreground">Minutes</div>
          </div>
          <div>
            <div className={cn('text-2xl font-bold', isEndingSoon && 'text-orange-600')}>
              {timeLeft.seconds}
            </div>
            <div className="text-xs text-muted-foreground">Seconds</div>
          </div>
        </div>

        {showLabel && (
          <div className="mt-2 text-center">
            <span
              className={cn(
                'text-sm font-medium',
                hasEnded
                  ? 'text-muted-foreground'
                  : isEndingSoon
                    ? 'text-orange-600'
                    : 'text-foreground'
              )}
            >
              {formatTime(timeLeft)}
            </span>
          </div>
        )}
      </div>
    );
  }

  if (variant === 'compact') {
    return (
      <Badge className={cn(getVariantStyles(), getStatusColor(), className)}>
        {formatTime(timeLeft)}
      </Badge>
    );
  }

  // Default variant
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Badge className={getVariantStyles()}>
        {showLabel ? getStatusText() : formatTime(timeLeft)}
      </Badge>
      {showLabel && (
        <span
          className={cn(
            'text-sm font-medium',
            hasEnded
              ? 'text-muted-foreground'
              : isEndingSoon
                ? 'text-orange-600'
                : 'text-foreground'
          )}
        >
          {formatTime(timeLeft)}
        </span>
      )}
    </div>
  );
}

// Hook for countdown functionality
export function useCountdown(
  endTime: string | Date,
  options: {
    onEnd?: () => void;
    onEndingSoon?: (remainingTime: number) => void;
    endingSoonThreshold?: number;
  } = {}
) {
  const [timeLeft, setTimeLeft] = useState(0);
  const [isEndingSoon, setIsEndingSoon] = useState(false);
  const [hasEnded, setHasEnded] = useState(false);

  const { onEnd, onEndingSoon, endingSoonThreshold = 300 } = options;

  useEffect(() => {
    const timer = setInterval(() => {
      const remaining = Math.max(0, Math.floor((new Date(endTime).getTime() - Date.now()) / 1000));
      setTimeLeft(remaining);

      const endingSoon = remaining > 0 && remaining < endingSoonThreshold;
      if (endingSoon && !isEndingSoon) {
        setIsEndingSoon(true);
        onEndingSoon?.(remaining);
      }

      if (remaining === 0 && !hasEnded) {
        setHasEnded(true);
        onEnd?.();
      }
    }, 1000);

    // Initial calculation
    const initial = Math.max(0, Math.floor((new Date(endTime).getTime() - Date.now()) / 1000));
    setTimeLeft(initial);

    return () => clearInterval(timer);
  }, [endTime, onEnd, onEndingSoon, endingSoonThreshold, isEndingSoon, hasEnded]);

  return {
    timeLeft,
    isEndingSoon,
    hasEnded,
    formatTime: (seconds: number) => {
      if (seconds <= 0) return 'Ended';

      const days = Math.floor(seconds / 86400);
      const hours = Math.floor((seconds % 86400) / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      const secs = seconds % 60;

      if (days > 0) return `${days}d ${hours}h`;
      if (hours > 0) return `${hours}h ${minutes}m`;
      if (minutes > 0) return `${minutes}m ${secs}s`;
      return `${secs}s`;
    },
  };
}

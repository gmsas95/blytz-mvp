'use client';

import { Loader2, TrendingUp, AlertCircle, CheckCircle } from 'lucide-react';
import { useState, useCallback, useEffect } from 'react';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { formatPrice , cn } from '@/lib/utils';

interface BidButtonProps {
  currentBid: number;
  minBidIncrement: number;
  onPlaceBid: (amount: number) => Promise<{ success: boolean; error?: string; data?: any }>;
  auctionId: string;
  currentUserId?: string;
  disabled?: boolean;
  disabledReason?: string;
  className?: string;
  variant?: 'button' | 'card' | 'compact';
  showQuickBids?: boolean;
  quickBidAmounts?: number[];
}

interface BidState {
  isPlacing: boolean;
  lastBidAmount?: number;
  error?: string;
  success?: boolean;
  suggestedBid: number;
  customAmount: string;
}

export function BidButton({
  currentBid,
  minBidIncrement,
  onPlaceBid,
  auctionId,
  currentUserId,
  disabled = false,
  disabledReason,
  className,
  variant = 'button',
  showQuickBids = true,
  quickBidAmounts,
}: BidButtonProps) {
  const [bidState, setBidState] = useState<BidState>({
    isPlacing: false,
    suggestedBid: currentBid + minBidIncrement,
    customAmount: '',
  });

  const [showCustomBid, setShowCustomBid] = useState(false);

  // Update suggested bid when current bid changes
  useEffect(() => {
    setBidState((prev) => ({
      ...prev,
      suggestedBid: currentBid + minBidIncrement,
      error: undefined,
      success: false,
    }));
  }, [currentBid, minBidIncrement]);

  const generateQuickBids = useCallback(() => {
    if (quickBidAmounts) {
      return quickBidAmounts.filter((amount) => amount > currentBid);
    }

    // Generate default quick bid amounts
    const base = currentBid + minBidIncrement;
    return [
      base,
      base + minBidIncrement * 2,
      base + minBidIncrement * 5,
      base + minBidIncrement * 10,
    ].slice(0, 4);
  }, [currentBid, minBidIncrement, quickBidAmounts]);

  const handlePlaceBid = useCallback(
    async (amount: number) => {
      if (disabled || bidState.isPlacing) return;

      setBidState((prev) => ({
        ...prev,
        isPlacing: true,
        error: undefined,
        success: false,
      }));

      try {
        const result = await onPlaceBid(amount);

        if (result.success) {
          setBidState((prev) => ({
            ...prev,
            isPlacing: false,
            lastBidAmount: amount,
            success: true,
            customAmount: '',
            error: undefined,
          }));

          // Clear success message after 3 seconds
          setTimeout(() => {
            setBidState((prev) => ({ ...prev, success: false }));
          }, 3000);

          // Close custom bid input if open
          if (showCustomBid) {
            setShowCustomBid(false);
          }
        } else {
          setBidState((prev) => ({
            ...prev,
            isPlacing: false,
            error: result.error || 'Failed to place bid',
            success: false,
          }));
        }
      } catch (error) {
        setBidState((prev) => ({
          ...prev,
          isPlacing: false,
          error: error instanceof Error ? error.message : 'Unknown error occurred',
          success: false,
        }));
      }
    },
    [disabled, bidState.isPlacing, onPlaceBid, showCustomBid]
  );

  const handleCustomBidSubmit = useCallback(() => {
    const amount = parseFloat(bidState.customAmount);
    if (isNaN(amount) || amount <= currentBid) {
      setBidState((prev) => ({
        ...prev,
        error: `Bid must be greater than ${formatPrice(currentBid)}`,
      }));
      return;
    }

    if (amount < currentBid + minBidIncrement) {
      setBidState((prev) => ({
        ...prev,
        error: `Minimum bid is ${formatPrice(currentBid + minBidIncrement)}`,
      }));
      return;
    }

    handlePlaceBid(amount);
  }, [bidState.customAmount, currentBid, minBidIncrement, handlePlaceBid]);

  const quickBids = generateQuickBids();

  if (variant === 'compact') {
    return (
      <div className={cn('flex items-center gap-2', className)}>
        <Button
          size="sm"
          onClick={() => handlePlaceBid(bidState.suggestedBid)}
          disabled={disabled || bidState.isPlacing}
          className="flex-1"
        >
          {bidState.isPlacing ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            `Bid ${formatPrice(bidState.suggestedBid)}`
          )}
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={() => setShowCustomBid(!showCustomBid)}
          disabled={disabled}
        >
          Custom
        </Button>
      </div>
    );
  }

  if (variant === 'card') {
    return (
      <Card className={cn('w-full', className)}>
        <CardContent className="p-6">
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <div className="text-sm text-muted-foreground">Current Bid</div>
                <div className="text-2xl font-bold">{formatPrice(currentBid)}</div>
              </div>
              <div className="text-right">
                <div className="text-sm text-muted-foreground">Min Increment</div>
                <div className="text-lg font-semibold">+{formatPrice(minBidIncrement)}</div>
              </div>
            </div>

            {bidState.success && (
              <Alert className="bg-green-50 border-green-200">
                <CheckCircle className="h-4 w-4 text-green-600" />
                <AlertDescription className="text-green-800">
                  Your bid of {formatPrice(bidState.lastBidAmount!)} was placed successfully!
                </AlertDescription>
              </Alert>
            )}

            {bidState.error && (
              <Alert className="bg-red-50 border-red-200">
                <AlertCircle className="h-4 w-4 text-red-600" />
                <AlertDescription className="text-red-800">{bidState.error}</AlertDescription>
              </Alert>
            )}

            {disabled && disabledReason && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{disabledReason}</AlertDescription>
              </Alert>
            )}

            <div className="space-y-3">
              {showQuickBids && quickBids.length > 0 && !showCustomBid && (
                <div className="grid grid-cols-2 gap-2">
                  {quickBids.slice(0, 4).map((amount, index) => (
                    <Button
                      key={`${auctionId}-quick-${index}`}
                      variant="outline"
                      onClick={() => handlePlaceBid(amount)}
                      disabled={disabled || bidState.isPlacing}
                      className="h-auto py-3"
                    >
                      <div className="text-center">
                        <div className="font-semibold">{formatPrice(amount)}</div>
                        {index === 0 && (
                          <Badge variant="secondary" className="text-xs mt-1">
                            Min Bid
                          </Badge>
                        )}
                      </div>
                    </Button>
                  ))}
                </div>
              )}

              {showCustomBid && (
                <div className="space-y-2">
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      placeholder={`Enter amount > ${formatPrice(currentBid)}`}
                      value={bidState.customAmount}
                      onChange={(e) =>
                        setBidState((prev) => ({ ...prev, customAmount: e.target.value }))
                      }
                      disabled={bidState.isPlacing}
                      min={currentBid + minBidIncrement}
                      step={minBidIncrement}
                    />
                    <Button onClick={handleCustomBidSubmit} disabled={bidState.isPlacing}>
                      {bidState.isPlacing ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        'Place Bid'
                      )}
                    </Button>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setShowCustomBid(false)}
                    disabled={bidState.isPlacing}
                  >
                    Cancel
                  </Button>
                </div>
              )}

              {!showCustomBid && (
                <div className="flex gap-2">
                  <Button
                    onClick={() => handlePlaceBid(bidState.suggestedBid)}
                    disabled={disabled || bidState.isPlacing}
                    className="flex-1"
                    size="lg"
                  >
                    {bidState.isPlacing ? (
                      <Loader2 className="h-4 w-4 animate-spin mr-2" />
                    ) : (
                      <TrendingUp className="h-4 w-4 mr-2" />
                    )}
                    Place Bid {formatPrice(bidState.suggestedBid)}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => setShowCustomBid(true)}
                    disabled={disabled || bidState.isPlacing}
                  >
                    Custom
                  </Button>
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Default button variant
  return (
    <div className={cn('space-y-3', className)}>
      {bidState.success && (
        <Alert className="bg-green-50 border-green-200">
          <CheckCircle className="h-4 w-4 text-green-600" />
          <AlertDescription className="text-green-800">Bid placed successfully!</AlertDescription>
        </Alert>
      )}

      {bidState.error && (
        <Alert className="bg-red-50 border-red-200">
          <AlertCircle className="h-4 w-4 text-red-600" />
          <AlertDescription className="text-red-800">{bidState.error}</AlertDescription>
        </Alert>
      )}

      {disabled && disabledReason && (
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{disabledReason}</AlertDescription>
        </Alert>
      )}

      <Button
        onClick={() => handlePlaceBid(bidState.suggestedBid)}
        disabled={disabled || bidState.isPlacing}
        className="w-full"
        size="lg"
      >
        {bidState.isPlacing ? (
          <Loader2 className="h-4 w-4 animate-spin mr-2" />
        ) : (
          <TrendingUp className="h-4 w-4 mr-2" />
        )}
        Place Bid {formatPrice(bidState.suggestedBid)}
      </Button>

      {showQuickBids && quickBids.length > 0 && (
        <div className="flex gap-2">
          {quickBids.slice(0, 3).map((amount, index) => (
            <Button
              key={`${auctionId}-quick-${index}`}
              variant="outline"
              onClick={() => handlePlaceBid(amount)}
              disabled={disabled || bidState.isPlacing}
              className="flex-1"
            >
              {formatPrice(amount)}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}

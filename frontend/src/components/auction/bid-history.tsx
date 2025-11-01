'use client';

import { useState, useEffect } from 'react';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { formatPrice , cn } from '@/lib/utils';
import { Bid } from '@/types';

interface BidHistoryProps {
  bids: Bid[];
  currentUserId?: string;
  auctionId?: string;
  maxVisible?: number;
  showMoreButton?: boolean;
  className?: string;
  realtime?: boolean;
  onNewBid?: (bid: Bid) => void;
}

interface BidItemProps {
  bid: Bid;
  isCurrentUser?: boolean;
  isWinning?: boolean;
  isNew?: boolean;
  showAnimation?: boolean;
}

function BidItem({ bid, isCurrentUser, isWinning, isNew, showAnimation = true }: BidItemProps) {
  const timeAgo = new Date(bid.timestamp).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
  });

  return (
    <div
      className={cn(
        'flex items-center justify-between p-3 rounded-lg border transition-all duration-300',
        isWinning && 'bg-green-50 border-green-200',
        isCurrentUser && 'bg-blue-50 border-blue-200',
        isNew && showAnimation && 'animate-pulse bg-yellow-50 border-yellow-200',
        !isWinning && !isCurrentUser && 'bg-gray-50 border-gray-200'
      )}
    >
      <div className="flex items-center gap-3">
        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-sm font-semibold">
          {bid.user.name.charAt(0).toUpperCase()}
        </div>
        <div>
          <div className="font-medium text-sm">
            {bid.user.name}
            {isCurrentUser && (
              <Badge variant="secondary" className="ml-2 text-xs">
                You
              </Badge>
            )}
          </div>
          <div className="text-xs text-muted-foreground">{timeAgo}</div>
        </div>
      </div>
      <div className="text-right">
        <div
          className={cn(
            'font-semibold text-sm',
            isWinning && 'text-green-600',
            isCurrentUser && !isWinning && 'text-blue-600'
          )}
        >
          {formatPrice(bid.amount)}
        </div>
        {isWinning && (
          <Badge
            variant="outline"
            className="text-xs mt-1 bg-green-100 text-green-800 border-green-300"
          >
            Leading
          </Badge>
        )}
      </div>
    </div>
  );
}

export function BidHistory({
  bids,
  currentUserId,
  auctionId,
  maxVisible = 10,
  showMoreButton = true,
  className,
  realtime = false,
  onNewBid,
}: BidHistoryProps) {
  const [visibleBids, setVisibleBids] = useState<Bid[]>([]);
  const [showAll, setShowAll] = useState(false);
  const [newBidIds, setNewBidIds] = useState<Set<string>>(new Set());

  // Update visible bids when bids change
  useEffect(() => {
    if (realtime && bids.length > visibleBids.length) {
      // New bid detected
      const newBid = bids[0];
      if (newBid && !visibleBids.find((b) => b.id === newBid.id)) {
        onNewBid?.(newBid);
        setNewBidIds((prev) => new Set(prev).add(newBid.id));

        // Remove new bid highlight after 3 seconds
        setTimeout(() => {
          setNewBidIds((prev) => {
            const updated = new Set(prev);
            updated.delete(newBid.id);
            return updated;
          });
        }, 3000);
      }
    }

    setVisibleBids(bids);
  }, [bids, visibleBids.length, realtime, onNewBid]);

  const displayBids = showAll ? visibleBids : visibleBids.slice(0, maxVisible);
  const hasMoreBids = visibleBids.length > maxVisible;

  const getLeadingBidderId = () => {
    if (visibleBids.length === 0) return null;
    return visibleBids[0].userId;
  };

  const leadingBidderId = getLeadingBidderId();

  return (
    <Card className={cn('w-full', className)}>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center justify-between">
          <span>Bid History</span>
          {visibleBids.length > 0 && (
            <Badge variant="secondary" className="text-xs">
              {visibleBids.length} {visibleBids.length === 1 ? 'bid' : 'bids'}
            </Badge>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        {visibleBids.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            <div className="text-sm">No bids yet</div>
            <div className="text-xs mt-1">Be the first to place a bid!</div>
          </div>
        ) : (
          <ScrollArea className="h-[400px] pr-4">
            <div className="space-y-2">
              {displayBids.map((bid, index) => (
                <BidItem
                  key={`${auctionId}-${bid.id}`}
                  bid={bid}
                  isCurrentUser={bid.userId === currentUserId}
                  isWinning={bid.userId === leadingBidderId}
                  isNew={newBidIds.has(bid.id)}
                />
              ))}
            </div>

            {hasMoreBids && showMoreButton && !showAll && (
              <div className="mt-4 text-center">
                <button
                  onClick={() => setShowAll(true)}
                  className="text-sm text-blue-600 hover:text-blue-800 font-medium"
                >
                  Show {visibleBids.length - maxVisible} more bids
                </button>
              </div>
            )}

            {showAll && hasMoreBids && (
              <div className="mt-4 text-center">
                <button
                  onClick={() => setShowAll(false)}
                  className="text-sm text-blue-600 hover:text-blue-800 font-medium"
                >
                  Show less
                </button>
              </div>
            )}
          </ScrollArea>
        )}
      </CardContent>
    </Card>
  );
}

// Live bid activity feed component
export function LiveBidActivity({
  bids,
  maxVisible = 5,
  className,
}: {
  bids: Bid[];
  maxVisible?: number;
  className?: string;
}) {
  const [activity, setActivity] = useState<
    Array<{
      bid: Bid;
      timestamp: Date;
      type: 'new_bid' | 'outbid' | 'winning';
    }>
  >([]);

  useEffect(() => {
    const newActivity = bids.slice(0, maxVisible).map((bid, index) => ({
      bid,
      timestamp: new Date(bid.timestamp),
      type: index === 0 ? 'winning' : ('new_bid' as const),
    }));
    setActivity(newActivity);
  }, [bids, maxVisible]);

  return (
    <Card className={cn('w-full', className)}>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2">
          <span>Live Activity</span>
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {activity.length === 0 ? (
            <div className="text-center py-4 text-muted-foreground text-sm">No recent activity</div>
          ) : (
            activity.map((item, index) => (
              <div key={item.bid.id} className="flex items-center gap-3 text-sm">
                <div
                  className={cn(
                    'w-2 h-2 rounded-full',
                    item.type === 'winning' ? 'bg-green-500' : 'bg-blue-500'
                  )}
                ></div>
                <div className="flex-1">
                  <span className="font-medium">{item.bid.user.name}</span>
                  <span className="text-muted-foreground ml-1">
                    {item.type === 'winning' ? 'is leading with' : 'bid'}
                  </span>
                  <span className="font-semibold ml-1">{formatPrice(item.bid.amount)}</span>
                </div>
                <div className="text-xs text-muted-foreground">
                  {item.timestamp.toLocaleTimeString([], {
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </div>
              </div>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
}

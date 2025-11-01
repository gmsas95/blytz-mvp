'use client';

import { Loader2, Users, Gavel, Clock, AlertTriangle, CheckCircle } from 'lucide-react';
import { useState, useEffect } from 'react';

import { BidButton } from '@/components/auction/bid-button';
import { BidHistory, LiveBidActivity } from '@/components/auction/bid-history';
import { ConnectionStatus } from '@/components/auction/connection-status';
import { RealtimeCountdown } from '@/components/auction/realtime-countdown';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ToastProvider, useToast } from '@/components/ui/toast';
import { useAuctionRealtime } from '@/hooks/useAuctionRealtime';
import { formatPrice } from '@/lib/utils';


interface AuctionDetailPageProps {
  params: { id: string };
}

function AuctionDetailContent({ params }: AuctionDetailPageProps) {
  const { success } = useToast();
  const [currentUserId] = useState('current-user'); // This would come from auth context

  const {
    auction,
    isConnected,
    isUsingFallback,
    lastUpdate,
    connectionQuality,
    error,
    remainingTime,
    isEndingSoon,
    hasEnded,
    placeBid,
    connectionStatus,
  } = useAuctionRealtime({
    auctionId: params.id,
    enableRealtime: true,
    fallbackToPolling: true,
    onBidUpdate: (data) => {
      success('New Bid Placed!', `${data.bid.user.name} bid ${formatPrice(data.bid.amount)}`);
    },
    onCountdownUpdate: (data) => {
      if (data.isEndingSoon && !data.remainingTime) {
        success('Auction Ended!', 'The auction has concluded.');
      }
    },
    onStatusChange: (data) => {
      if (data.newStatus === 'ended') {
        success(
          'Auction Ended!',
          data.winner
            ? `${data.winner.user.name} won with ${formatPrice(data.winner.amount)}`
            : 'The auction has ended without a winner.'
        );
      }
    },
    onError: (error) => {
      console.error('Real-time error:', error);
    },
  });

  // Handle bid placement
  const handlePlaceBid = async (amount: number) => {
    try {
      const result = await placeBid(amount);
      return result;
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to place bid',
      };
    }
  };

  // Join auction when component mounts
  useEffect(() => {
    if (auction && isConnected) {
      // Auction will be joined automatically by the hook
    }
  }, [auction, isConnected]);

  if (error && !auction) {
    return (
      <section className="w-full py-16 md:py-24">
        <div className="container mx-auto px-4">
          <Alert className="max-w-2xl mx-auto">
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>Failed to load auction. Please try again later.</AlertDescription>
          </Alert>
        </div>
      </section>
    );
  }

  if (!auction) {
    return (
      <section className="w-full py-16 md:py-24">
        <div className="container mx-auto px-4 flex items-center justify-center">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      </section>
    );
  }

  const isAuctionActive = auction.status === 'active';
  const isAuctionScheduled = auction.status === 'scheduled';
  const isAuctionEnded = auction.status === 'ended';
  const isCurrentUserWinning = auction.bids[0]?.userId === currentUserId;

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4">
        {/* Connection Status */}
        <div className="mb-6">
          <ConnectionStatus
            connected={isConnected}
            connecting={connectionStatus.connecting}
            reconnectAttempts={connectionStatus.reconnectAttempts}
            error={connectionStatus.error}
            connectionQuality={connectionQuality}
            isUsingFallback={isUsingFallback}
            variant="detailed"
          />
        </div>

        <div className="grid gap-8 lg:grid-cols-2">
          {/* Left Column - Product Images */}
          <div className="space-y-6">
            <div className="aspect-square bg-muted rounded-2xl overflow-hidden">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={auction.product.images[0]}
                alt={auction.product.title}
                className="object-cover w-full h-full"
              />
            </div>

            {/* Product Details */}
            <Card>
              <CardHeader>
                <CardTitle>Product Details</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <h3 className="font-semibold text-lg">{auction.product.title}</h3>
                    <p className="text-muted-foreground mt-2">{auction.product.description}</p>
                  </div>

                  <Separator />

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    {Object.entries(auction.product.specifications).map(([key, value]) => (
                      <div key={key}>
                        <span className="text-muted-foreground">{key}:</span>
                        <div className="font-medium">{value}</div>
                      </div>
                    ))}
                  </div>

                  <Separator />

                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                      {auction.product.seller.name.charAt(0).toUpperCase()}
                    </div>
                    <div>
                      <div className="font-medium">{auction.product.seller.name}</div>
                      <div className="text-sm text-muted-foreground">
                        {auction.product.seller.storeName} • Rating: {auction.product.seller.rating}
                        /5
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Right Column - Auction Details */}
          <div className="space-y-6">
            {/* Auction Status */}
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <Badge
                  className={cn(
                    'text-xs px-3 py-1',
                    auction.isLive ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'
                  )}
                >
                  {auction.isLive ? '🔴 LIVE' : auction.status.toUpperCase()}
                </Badge>

                <RealtimeCountdown
                  endTime={auction.endTime}
                  variant="default"
                  showLabel={true}
                  isActive={isAuctionActive}
                />

                {isEndingSoon && isAuctionActive && (
                  <Badge className="bg-orange-100 text-orange-800 animate-pulse">
                    ⏰ Ending Soon
                  </Badge>
                )}
              </div>

              {/* Current Price */}
              <div className="space-y-2">
                <div className="text-3xl font-bold text-primary">
                  {formatPrice(auction.currentBid)}
                </div>
                <div className="text-sm text-muted-foreground">
                  Starting from {formatPrice(auction.startingPrice)} • Min increment +
                  {formatPrice(auction.minBidIncrement)}
                </div>
              </div>

              {/* Auction Stats */}
              <div className="grid grid-cols-3 gap-4">
                <div className="text-center p-3 bg-muted rounded-lg">
                  <Users className="h-4 w-4 mx-auto mb-1 text-muted-foreground" />
                  <div className="font-semibold">{auction.participants}</div>
                  <div className="text-xs text-muted-foreground">Bidders</div>
                </div>
                <div className="text-center p-3 bg-muted rounded-lg">
                  <Gavel className="h-4 w-4 mx-auto mb-1 text-muted-foreground" />
                  <div className="font-semibold">{auction.totalBids}</div>
                  <div className="text-xs text-muted-foreground">Total Bids</div>
                </div>
                <div className="text-center p-3 bg-muted rounded-lg">
                  <Clock className="h-4 w-4 mx-auto mb-1 text-muted-foreground" />
                  <div className="font-semibold">
                    {remainingTime > 0 ? formatTimeRemaining(remainingTime) : 'Ended'}
                  </div>
                  <div className="text-xs text-muted-foreground">Time Left</div>
                </div>
              </div>

              {/* Winner Information */}
              {isAuctionEnded && auction.winner && (
                <Alert className="bg-green-50 border-green-200">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <AlertDescription className="text-green-800">
                    <strong>Auction Winner:</strong> {auction.winner.name} won with{' '}
                    {formatPrice(auction.currentBid)}
                  </AlertDescription>
                </Alert>
              )}

              {/* Bid Button */}
              {isAuctionActive && (
                <BidButton
                  currentBid={auction.currentBid}
                  minBidIncrement={auction.minBidIncrement}
                  onPlaceBid={handlePlaceBid}
                  auctionId={auction.id}
                  currentUserId={currentUserId}
                  disabled={hasEnded}
                  disabledReason={hasEnded ? 'Auction has ended' : undefined}
                  variant="card"
                  showQuickBids={true}
                />
              )}

              {isAuctionScheduled && (
                <Card>
                  <CardContent className="p-6 text-center">
                    <Clock className="h-8 w-8 mx-auto mb-3 text-muted-foreground" />
                    <h3 className="font-semibold mb-2">Auction Starts Soon</h3>
                    <p className="text-sm text-muted-foreground mb-4">
                      This auction is scheduled to start on{' '}
                      {new Date(auction.startTime).toLocaleString()}
                    </p>
                    <Button variant="outline" disabled>
                      Not Started Yet
                    </Button>
                  </CardContent>
                </Card>
              )}

              {isAuctionEnded && (
                <Card>
                  <CardContent className="p-6 text-center">
                    <Gavel className="h-8 w-8 mx-auto mb-3 text-muted-foreground" />
                    <h3 className="font-semibold mb-2">Auction Ended</h3>
                    <p className="text-sm text-muted-foreground mb-4">
                      This auction has concluded. Final price: {formatPrice(auction.currentBid)}
                    </p>
                    <Button variant="outline" disabled>
                      Auction Closed
                    </Button>
                  </CardContent>
                </Card>
              )}
            </div>

            {/* Product Actions */}
            <div className="flex gap-3">
              <Button
                variant="outline"
                className="flex-1"
                onClick={() => window.open(`/product/${auction.productId}`, '_blank')}
              >
                View Product Details
              </Button>
              {auction.streamUrl && isAuctionActive && (
                <Button className="flex-1" onClick={() => window.open(auction.streamUrl, '_blank')}>
                  Watch Live Stream
                </Button>
              )}
            </div>

            {/* Last Update Info */}
            {lastUpdate && (
              <div className="text-xs text-muted-foreground text-center">
                Last updated: {lastUpdate.toLocaleTimeString()}
              </div>
            )}
          </div>
        </div>

        {/* Bottom Section - Tabs */}
        <div className="mt-12">
          <Tabs defaultValue="bids" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="bids">Bid History</TabsTrigger>
              <TabsTrigger value="activity">Live Activity</TabsTrigger>
            </TabsList>

            <TabsContent value="bids" className="mt-6">
              <BidHistory
                bids={auction.bids}
                currentUserId={currentUserId}
                auctionId={auction.id}
                maxVisible={20}
                realtime={true}
                onNewBid={(bid) => {
                  console.log('New bid received:', bid);
                }}
              />
            </TabsContent>

            <TabsContent value="activity" className="mt-6">
              <LiveBidActivity bids={auction.bids} maxVisible={10} />
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </section>
  );
}

// Helper function to format time remaining
function formatTimeRemaining(seconds: number): string {
  if (seconds <= 0) return 'Ended';

  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;

  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  } else if (minutes > 0) {
    return `${minutes}m ${remainingSeconds}s`;
  } else {
    return `${remainingSeconds}s`;
  }
}

export default function AuctionDetailPage(props: AuctionDetailPageProps) {
  return (
    <ToastProvider>
      <AuctionDetailContent {...props} />
    </ToastProvider>
  );
}

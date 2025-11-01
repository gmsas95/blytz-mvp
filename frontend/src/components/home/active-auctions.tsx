'use client';

import { Clock, Users, TrendingUp } from 'lucide-react';
import Link from 'next/link';
import { useEffect, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { api } from '@/lib/api-adapter';
import { formatPrice, formatTimeRemaining } from '@/lib/utils';
import { Auction } from '@/types';


export function ActiveAuctions() {
  const [auctions, setAuctions] = useState<Auction[]>([]);
  const [loading, setLoading] = useState(true);
  const [timeRemaining, setTimeRemaining] = useState<Record<string, number>>({});

  useEffect(() => {
    async function loadActiveAuctions() {
      try {
        const response = await api.getActiveAuctions();
        if (response.success && response.data) {
          setAuctions(response.data);

          // Initialize time remaining
          const initialTimes: Record<string, number> = {};
          response.data.forEach((auction) => {
            const endTime = new Date(auction.endTime).getTime();
            const now = new Date().getTime();
            initialTimes[auction.id] = Math.max(0, Math.floor((endTime - now) / 1000));
          });
          setTimeRemaining(initialTimes);
        }
      } catch (error) {
        console.error('Failed to load active auctions:', error);
      } finally {
        setLoading(false);
      }
    }

    loadActiveAuctions();
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      setTimeRemaining((prev) => {
        const updated = { ...prev };
        Object.keys(updated).forEach((auctionId) => {
          if (updated[auctionId] > 0) {
            updated[auctionId] -= 1;
          }
        });
        return updated;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <section className="w-full py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="text-center space-y-4 mb-12">
            <div className="h-8 bg-muted rounded-lg w-64 mx-auto animate-pulse" />
            <div className="h-4 bg-muted rounded w-96 mx-auto animate-pulse" />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(3)].map((_, i) => (
              <Card key={i} className="overflow-hidden border-red-200 bg-red-50/50">
                <CardHeader className="pb-3">
                  <div className="h-4 bg-muted rounded animate-pulse mb-2" />
                  <div className="h-3 bg-muted rounded w-1/3 animate-pulse" />
                </CardHeader>
                <CardContent className="pb-3">
                  <div className="h-6 bg-muted rounded w-1/3 animate-pulse mb-4" />
                  <div className="h-4 bg-muted rounded w-2/3 animate-pulse mb-2" />
                  <div className="h-4 bg-muted rounded w-1/2 animate-pulse" />
                </CardContent>
                <CardFooter>
                  <div className="h-10 bg-muted rounded-lg animate-pulse w-full" />
                </CardFooter>
              </Card>
            ))}
          </div>
        </div>
      </section>
    );
  }

  if (auctions.length === 0) {
    return null;
  }

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4">
        <div className="text-center space-y-4 mb-12">
          <div className="inline-flex items-center gap-2 bg-red-100 text-red-800 px-4 py-2 rounded-full text-sm font-medium">
            <div className="w-2 h-2 bg-red-500 rounded-full animate-pulse" />
            Live Auctions
          </div>

          <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
            🔥 Hot Auctions
          </h2>

          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Don't miss out on these exciting live auctions ending soon
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {auctions.map((auction, index) => {
            const remaining = timeRemaining[auction.id] || 0;
            const isEndingSoon = remaining < 300; // 5 minutes

            return (
              <div key={auction.id} className="flex">
                <Card
                  className={`overflow-hidden flex-1 border-2 ${
                    isEndingSoon ? 'border-orange-300 animate-pulse' : 'border-border'
                  }`}
                >
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between mb-3">
                      <div
                        className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${
                          auction.isLive
                            ? 'bg-red-100 text-red-800'
                            : 'bg-yellow-100 text-yellow-800'
                        }`}
                      >
                        {auction.isLive ? (
                          <>
                            <div className="w-2 h-2 bg-red-500 rounded-full animate-pulse" />
                            LIVE
                          </>
                        ) : (
                          'SCHEDULED'
                        )}
                      </div>

                      <div
                        className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${
                          isEndingSoon
                            ? 'bg-orange-100 text-orange-800 animate-pulse'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        <Clock className="w-3 h-3" />
                        {formatTimeRemaining(remaining)}
                      </div>
                    </div>

                    <h3 className="font-semibold text-lg leading-tight">{auction.product.title}</h3>
                    <p className="text-sm text-muted-foreground line-clamp-2">
                      {auction.product.description}
                    </p>
                  </CardHeader>

                  <CardContent className="pb-3 space-y-4">
                    <div className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Current Bid</span>
                        <span className="text-2xl font-bold text-primary">
                          {formatPrice(auction.currentBid)}
                        </span>
                      </div>

                      <div className="flex items-center justify-between text-sm">
                        <div className="flex items-center gap-1 text-muted-foreground">
                          <TrendingUp className="w-4 h-4" />
                          <span>{auction.totalBids} bids</span>
                        </div>

                        <div className="flex items-center gap-1 text-muted-foreground">
                          <Users className="w-4 h-4" />
                          <span>{auction.participants} participants</span>
                        </div>
                      </div>
                    </div>

                    <div className="pt-2 border-t border-border">
                      <div className="flex items-center justify-between">
                        <div className="text-sm">
                          <div className="font-medium">{auction.product.seller.storeName}</div>
                          <div className="text-muted-foreground">
                            ⭐ {auction.product.seller.rating} ({auction.product.seller.totalSales}{' '}
                            sales)
                          </div>
                        </div>

                        <div className="text-right text-sm">
                          <div className="text-muted-foreground">Min bid</div>
                          <div className="font-medium">+{formatPrice(auction.minBidIncrement)}</div>
                        </div>
                      </div>
                    </div>
                  </CardContent>

                  <CardFooter>
                    <Link href={`/auctions/${auction.id}`} className="w-full">
                      <Button className="w-full">Join Auction</Button>
                    </Link>
                  </CardFooter>
                </Card>
              </div>
            );
          })}
        </div>

        <div className="text-center mt-12">
          <Link href="/auctions">
            <Button size="lg" variant="outline">
              View All Auctions
            </Button>
          </Link>
        </div>
      </div>
    </section>
  );
}

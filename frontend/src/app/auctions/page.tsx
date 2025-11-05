'use client';

import { Loader2, Search, Filter, Users, Gavel, Clock } from 'lucide-react';
import Link from 'next/link';
import { useState, useEffect } from 'react';

import { ConnectionIndicator } from '@/components/auction/connection-status';
import { RealtimeCountdown } from '@/components/auction/realtime-countdown';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { ToastProvider } from '@/components/ui/toast';
import { mockAuctions } from '@/data/mock-data';
import { useWebSocket } from '@/hooks/useWebSocket';
import { formatPrice } from '@/lib/utils';
import { Auction } from '@/types';


// Mock data for demonstration

function AuctionsPageContent() {
  const [auctions, setAuctions] = useState<Auction[]>(mockAuctions);
  const [filteredAuctions, setFilteredAuctions] = useState<Auction[]>(mockAuctions);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [sortBy, setSortBy] = useState<string>('ending_soon');

  // WebSocket connection for real-time updates
  const { connectionStatus, isConnected } = useWebSocket({
    enableMockMode: true,
    autoReconnect: true,
  });

  // Apply filters
  useEffect(() => {
    let filtered = [...auctions];

    // Apply search filter
    if (searchTerm) {
      filtered = filtered.filter(
        (auction) =>
          auction.product.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
          auction.product.description.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Apply status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter((auction) => auction.status === statusFilter);
    }

    // Apply sorting
    filtered.sort((a, b) => {
      switch (sortBy) {
        case 'ending_soon':
          return new Date(a.endTime).getTime() - new Date(b.endTime).getTime();
        case 'price_low':
          return a.currentBid - b.currentBid;
        case 'price_high':
          return b.currentBid - a.currentBid;
        case 'most_bids':
          return b.totalBids - a.totalBids;
        case 'newest':
          return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
        default:
          return 0;
      }
    });

    setFilteredAuctions(filtered);
  }, [auctions, searchTerm, statusFilter, sortBy]);

  // Simulate real-time updates in mock mode
  useEffect(() => {
    if (!isConnected) return;

    const interval = setInterval(() => {
      setAuctions((prevAuctions) =>
        prevAuctions.map((auction) => {
          if (auction.status === 'active' && Math.random() > 0.8) {
            // Random bid update - use minBidIncrement with fallback
            const minIncrement = auction.minBidIncrement || 10.0;
            const newBid = auction.currentBid + (Math.random() * minIncrement * 2);
            return {
              ...auction,
              currentBid: newBid,
              totalBids: auction.totalBids + 1,
              participants: auction.participants + (Math.random() > 0.5 ? 1 : 0),
            };
          }
          return auction;
        })
      );
    }, 5000);

    return () => clearInterval(interval);
  }, [isConnected]);

  const getStatusColor = (status: string, isLive: boolean) => {
    if (isLive) return 'bg-red-100 text-red-800';
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'scheduled':
        return 'bg-blue-100 text-blue-800';
      case 'ended':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusText = (status: string, isLive: boolean) => {
    if (isLive) return '🔴 LIVE';
    return status.toUpperCase();
  };

  return (
    <section className="w-full py-16 md:py-24" aria-labelledby="auctions-heading">
      <div className="container mx-auto px-4">
        {/* Header */}
        <header className="text-center space-y-4 mb-12">
          <h1 id="auctions-heading" className="text-3xl md:text-4xl font-bold tracking-tight">
            Live & Upcoming Auctions
          </h1>
          <p className="text-muted-foreground">
            Join live auctions and place bids on amazing products
          </p>

          {/* Connection Status */}
          <div className="flex justify-center">
            <ConnectionIndicator
              connected={isConnected}
              connecting={connectionStatus.connecting}
              connectionQuality="excellent"
            />
          </div>
        </header>

        {/* Filters */}
        <div className="flex flex-col md:flex-row gap-4 mb-8">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              placeholder="Search auctions..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10"
            />
          </div>

          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-full md:w-48">
              <Filter className="h-4 w-4 mr-2" />
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="scheduled">Scheduled</SelectItem>
              <SelectItem value="ended">Ended</SelectItem>
            </SelectContent>
          </Select>

          <Select value={sortBy} onValueChange={setSortBy}>
            <SelectTrigger className="w-full md:w-48">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ending_soon">Ending Soon</SelectItem>
              <SelectItem value="price_low">Price: Low to High</SelectItem>
              <SelectItem value="price_high">Price: High to Low</SelectItem>
              <SelectItem value="most_bids">Most Bids</SelectItem>
              <SelectItem value="newest">Newest First</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Auction Grid */}
        {loading ? (
          <div className="flex justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin" />
          </div>
        ) : filteredAuctions.length === 0 ? (
          <div className="text-center text-muted-foreground py-12" role="status" aria-live="polite">
            <div className="text-lg font-medium mb-2">No auctions found</div>
            <div className="text-sm">Try adjusting your filters or search terms</div>
          </div>
        ) : (
          <div role="region" aria-label="Auction listings">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredAuctions.map((auction) => (
                <article key={auction.id} className="group">
                  <Link href={`/auctions/${auction.id}`} className="block">
                    <Card className="overflow-hidden h-full transition-all duration-300 hover:shadow-lg hover:scale-[1.02]">
                      <CardHeader className="pb-3">
                        <div className="flex items-center justify-between mb-2">
                          <Badge
                            className={getStatusColor(auction.status, auction.isLive)}
                            role="status"
                            aria-live={auction.isLive ? 'polite' : 'off'}
                          >
                            {getStatusText(auction.status, auction.isLive)}
                          </Badge>
                          <RealtimeCountdown
                            endTime={auction.endTime}
                            variant="compact"
                            isActive={auction.status === 'active'}
                          />
                        </div>

                        {/* Product Image */}
                        <div className="aspect-square bg-muted rounded-lg overflow-hidden mb-3">
                          {/* eslint-disable-next-line @next/next/no-img-element */}
                          <img
                            src={auction.product.images[0]}
                            alt={auction.product.title}
                            className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
                          />
                        </div>

                        <h2 className="font-semibold text-lg leading-tight overflow-hidden" style={{
                          display: '-webkit-box',
                          WebkitLineClamp: 2,
                          WebkitBoxOrient: 'vertical'
                        }}>
                          {auction.product.title}
                        </h2>
                        <p className="text-sm text-muted-foreground overflow-hidden mt-1" style={{
                          display: '-webkit-box',
                          WebkitLineClamp: 2,
                          WebkitBoxOrient: 'vertical'
                        }}>
                          {auction.product.description}
                        </p>
                      </CardHeader>

                      <CardContent className="pb-3">
                        <div className="space-y-3">
                          {/* Current Bid */}
                          <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">Current Bid</span>
                            <span
                              className="font-semibold text-lg"
                              aria-label={`Current bid: ${formatPrice(auction.currentBid)}`}
                            >
                              {formatPrice(auction.currentBid)}
                            </span>
                          </div>

                          {/* Auction Stats */}
                          <div className="flex items-center justify-between text-sm text-muted-foreground">
                            <div className="flex items-center gap-1">
                              <Users className="h-3 w-3" />
                              <span>{auction.participants}</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <Gavel className="h-3 w-3" />
                              <span>{auction.totalBids} bids</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              <span>+{formatPrice(auction.minBidIncrement)}</span>
                            </div>
                          </div>

                          {/* Seller Info */}
                          <div className="flex items-center gap-2 pt-2 border-t">
                            <div className="w-6 h-6 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-xs font-semibold">
                              {auction.product.seller.name.charAt(0).toUpperCase()}
                            </div>
                            <span className="text-xs text-muted-foreground truncate">
                              {auction.product.seller.storeName}
                            </span>
                            <span className="text-xs text-muted-foreground ml-auto">
                              {auction.product.seller.rating}★
                            </span>
                          </div>
                        </div>
                      </CardContent>

                      <CardFooter>
                        <Button className="w-full" aria-describedby={`view-auction-${auction.id}`}>
                          {auction.status === 'active'
                            ? 'Join Auction'
                            : auction.status === 'scheduled'
                              ? 'View Details'
                              : 'View Results'}
                        </Button>
                        <span id={`view-auction-${auction.id}`} className="sr-only">
                          View details for {auction.product.title} auction
                        </span>
                      </CardFooter>
                    </Card>
                  </Link>
                </article>
              ))}
            </div>
          </div>
        )}

        {/* Results Summary */}
        <div className="mt-8 text-center text-sm text-muted-foreground">
          Showing {filteredAuctions.length} of {auctions.length} auctions
          {isConnected && ' • Live updates enabled'}
        </div>
      </div>
    </section>
  );
}

export default function AuctionsPage() {
  return (
    <ToastProvider>
      <AuctionsPageContent />
    </ToastProvider>
  );
}

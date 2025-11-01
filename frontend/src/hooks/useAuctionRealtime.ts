'use client';

import { useEffect, useState, useRef, useCallback } from 'react';

import { Auction, Bid } from '@/types';

import {
  useWebSocket,
  isBidUpdate,
  isCountdownUpdate,
  isStatusChange,
  isAuctionUpdate,
} from './useWebSocket';

interface UseAuctionRealtimeOptions {
  auctionId: string;
  enableRealtime?: boolean;
  fallbackToPolling?: boolean;
  pollingInterval?: number;
  onBidUpdate?: (data: {
    bid: Bid;
    isWinning: boolean;
    newCurrentBid: number;
    totalBids: number;
  }) => void;
  onCountdownUpdate?: (data: { remainingTime: number; isEndingSoon: boolean }) => void;
  onStatusChange?: (data: { oldStatus: string; newStatus: string; winner?: Bid }) => void;
  onAuctionUpdate?: (auction: Auction) => void;
  onError?: (error: string) => void;
}

interface RealtimeState {
  auction: Auction | null;
  isConnected: boolean;
  isUsingFallback: boolean;
  lastUpdate: Date | null;
  connectionQuality: 'excellent' | 'good' | 'poor' | 'disconnected';
  error?: string;
}

export function useAuctionRealtime({
  auctionId,
  enableRealtime = true,
  fallbackToPolling = true,
  pollingInterval = 5000,
  onBidUpdate,
  onCountdownUpdate,
  onStatusChange,
  onAuctionUpdate,
  onError,
}: UseAuctionRealtimeOptions) {
  const [state, setState] = useState<RealtimeState>({
    auction: null,
    isConnected: false,
    isUsingFallback: false,
    lastUpdate: null,
    connectionQuality: 'disconnected',
  });

  const pollingTimerRef = useRef<NodeJS.Timeout | null>(null);
  const unsubscribeRef = useRef<(() => void) | null>(null);

  // WebSocket integration
  const { connectionStatus, subscribeToAuction, sendMessage, lastMessage, connectionQuality } =
    useWebSocket({
      enableMockMode: process.env.NODE_ENV === 'development',
      autoReconnect: true,
      maxReconnectAttempts: 10,
      reconnectInterval: 3000,
    });

  // Start polling as fallback
  const startPolling = useCallback(() => {
    if (pollingTimerRef.current) {
      clearInterval(pollingTimerRef.current);
    }

    setState((prev) => ({ ...prev, isUsingFallback: true }));

    pollingTimerRef.current = setInterval(async () => {
      try {
        const response = await fetch(`/api/auctions/${auctionId}/realtime`);
        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data) {
            setState((prev) => ({
              ...prev,
              auction: data.data.auction,
              lastUpdate: new Date(),
              error: undefined,
            }));
          }
        }
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Polling error';
        setState((prev) => ({ ...prev, error: errorMessage }));
        onError?.(errorMessage);
      }
    }, pollingInterval);
  }, [auctionId, pollingInterval, onError]);

  // Stop polling
  const stopPolling = useCallback(() => {
    if (pollingTimerRef.current) {
      clearInterval(pollingTimerRef.current);
      pollingTimerRef.current = null;
    }
    setState((prev) => ({ ...prev, isUsingFallback: false }));
  }, []);

  // Handle WebSocket messages
  useEffect(() => {
    if (!lastMessage || lastMessage.auctionId !== auctionId) return;

    setState((prev) => ({ ...prev, lastUpdate: new Date(), error: undefined }));

    // Handle different message types
    if (isBidUpdate(lastMessage)) {
      const { bid, isWinning, newCurrentBid, totalBids } = lastMessage.data;

      setState((prev) => ({
        ...prev,
        auction: prev.auction
          ? {
              ...prev.auction,
              currentBid: newCurrentBid,
              totalBids,
              bids: [bid, ...prev.auction.bids],
            }
          : null,
      }));

      onBidUpdate?.({ bid, isWinning, newCurrentBid, totalBids });
    }

    if (isCountdownUpdate(lastMessage)) {
      const { remainingTime, isEndingSoon } = lastMessage.data;

      setState((prev) => ({
        ...prev,
        auction: prev.auction
          ? {
              ...prev.auction,
              endTime: new Date(Date.now() + remainingTime * 1000).toISOString(),
            }
          : null,
      }));

      onCountdownUpdate?.({ remainingTime, isEndingSoon });
    }

    if (isStatusChange(lastMessage)) {
      const { oldStatus, newStatus, winner } = lastMessage.data;

      setState((prev) => ({
        ...prev,
        auction: prev.auction
          ? {
              ...prev.auction,
              status: newStatus as 'scheduled' | 'active' | 'ended',
              winner: winner?.user,
            }
          : null,
      }));

      onStatusChange?.({ oldStatus, newStatus, winner });
    }

    if (isAuctionUpdate(lastMessage)) {
      const { auction } = lastMessage.data;
      setState((prev) => ({ ...prev, auction }));
      onAuctionUpdate?.(auction);
    }
  }, [lastMessage, auctionId, onBidUpdate, onCountdownUpdate, onStatusChange, onAuctionUpdate]);

  // Handle connection status changes
  useEffect(() => {
    const isConnected = connectionStatus.connected;
    setState((prev) => ({
      ...prev,
      isConnected,
      connectionQuality,
    }));

    // If WebSocket disconnects and fallback is enabled, start polling
    if (!isConnected && fallbackToPolling && enableRealtime) {
      startPolling();
    } else if (isConnected) {
      stopPolling();
    }

    // Handle connection errors
    if (connectionStatus.error) {
      setState((prev) => ({ ...prev, error: connectionStatus.error }));
      onError?.(connectionStatus.error);
    }
  }, [
    connectionStatus,
    connectionQuality,
    fallbackToPolling,
    enableRealtime,
    startPolling,
    stopPolling,
    onError,
  ]);

  // Subscribe to auction updates
  useEffect(() => {
    if (enableRealtime && connectionStatus.connected) {
      unsubscribeRef.current = subscribeToAuction(auctionId);
    }

    return () => {
      if (unsubscribeRef.current) {
        unsubscribeRef.current();
      }
    };
  }, [enableRealtime, connectionStatus.connected, auctionId, subscribeToAuction]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      stopPolling();
      if (unsubscribeRef.current) {
        unsubscribeRef.current();
      }
    };
  }, [stopPolling]);

  // Initialize auction data
  useEffect(() => {
    const initializeAuction = async () => {
      try {
        const response = await fetch(`/api/auctions/${auctionId}`);
        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data) {
            setState((prev) => ({
              ...prev,
              auction: data.data,
              lastUpdate: new Date(),
            }));
          }
        }
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to load auction';
        setState((prev) => ({ ...prev, error: errorMessage }));
        onError?.(errorMessage);
      }
    };

    if (auctionId && !state.auction) {
      initializeAuction();
    }
  }, [auctionId, state.auction, onError]);

  // Place bid through WebSocket or fallback
  const placeBid = useCallback(
    async (amount: number) => {
      try {
        if (connectionStatus.connected) {
          // Send bid via WebSocket
          sendMessage({
            type: 'bid_update',
            auctionId,
            data: { action: 'place_bid', amount },
          });

          // Optimistic update (optional)
          // This will be overridden by the server response
          return { success: true, data: { message: 'Bid placed via real-time connection' } };
        } else {
          // Fallback to HTTP
          const response = await fetch(`/api/auctions/${auctionId}/bid`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ amount }),
          });

          const data = await response.json();
          return data;
        }
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to place bid';
        return { success: false, error: errorMessage };
      }
    },
    [connectionStatus.connected, auctionId, sendMessage]
  );

  // Join auction room for real-time updates
  const joinAuction = useCallback(() => {
    if (connectionStatus.connected) {
      sendMessage({
        type: 'auction_update',
        auctionId,
        data: { action: 'join_auction' },
      });
    }
  }, [connectionStatus.connected, auctionId, sendMessage]);

  // Leave auction room
  const leaveAuction = useCallback(() => {
    if (connectionStatus.connected) {
      sendMessage({
        type: 'auction_update',
        auctionId,
        data: { action: 'leave_auction' },
      });
    }
  }, [connectionStatus.connected, auctionId, sendMessage]);

  // Get remaining time in seconds
  const getRemainingTime = useCallback((): number => {
    if (!state.auction) return 0;

    const endTime = new Date(state.auction.endTime).getTime();
    const now = Date.now();
    return Math.max(0, Math.floor((endTime - now) / 1000));
  }, [state.auction]);

  // Check if auction is ending soon (less than 5 minutes)
  const isEndingSoon = useCallback((): boolean => {
    const remaining = getRemainingTime();
    return remaining > 0 && remaining < 300;
  }, [getRemainingTime]);

  // Check if auction has ended
  const hasEnded = useCallback((): boolean => {
    return getRemainingTime() === 0 || state.auction?.status === 'ended';
  }, [getRemainingTime, state.auction]);

  return {
    // State
    auction: state.auction,
    isConnected: state.isConnected,
    isUsingFallback: state.isUsingFallback,
    lastUpdate: state.lastUpdate,
    connectionQuality: state.connectionQuality,
    error: state.error,

    // Computed values
    remainingTime: getRemainingTime(),
    isEndingSoon: isEndingSoon(),
    hasEnded: hasEnded(),

    // Actions
    placeBid,
    joinAuction,
    leaveAuction,

    // Utilities
    refresh: () => {
      if (state.isUsingFallback) {
        startPolling();
      }
    },

    // Connection info
    connectionStatus: {
      connected: connectionStatus.connected,
      connecting: connectionStatus.connecting,
      reconnectAttempts: connectionStatus.reconnectAttempts,
      error: connectionStatus.error,
    },
  };
}

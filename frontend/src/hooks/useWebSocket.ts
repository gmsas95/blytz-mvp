'use client';

import { useEffect, useRef, useState, useCallback } from 'react';

import { Auction, Bid } from '@/types';

// WebSocket message types
export interface WebSocketMessage {
  type:
    | 'auction_update'
    | 'bid_update'
    | 'countdown_update'
    | 'status_change'
    | 'ping'
    | 'pong'
    | 'error';
  data: any;
  timestamp: string;
  auctionId?: string;
}

export interface AuctionUpdateData {
  auction: Auction;
  updateType: 'current_bid' | 'status' | 'participants' | 'ending_soon';
}

export interface BidUpdateData {
  auctionId: string;
  bid: Bid;
  isWinning: boolean;
  newCurrentBid: number;
  totalBids: number;
}

export interface CountdownUpdateData {
  auctionId: string;
  remainingTime: number;
  isEndingSoon: boolean;
}

export interface StatusChangeData {
  auctionId: string;
  oldStatus: string;
  newStatus: string;
  winner?: Bid;
}

export interface ConnectionStatus {
  connected: boolean;
  connecting: boolean;
  error?: string;
  reconnectAttempts: number;
  lastConnected?: Date;
}

interface UseWebSocketOptions {
  autoReconnect?: boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  heartbeatInterval?: number;
  enableMockMode?: boolean;
}

interface UseWebSocketReturn {
  connectionStatus: ConnectionStatus;
  subscribeToAuction: (auctionId: string) => () => void;
  unsubscribeFromAuction: (auctionId: string) => void;
  sendMessage: (message: Omit<WebSocketMessage, 'timestamp'>) => void;
  lastMessage: WebSocketMessage | null;
  connectionQuality: 'excellent' | 'good' | 'poor' | 'disconnected';
}

// Global WebSocket connection management
class WebSocketManager {
  private static instance: WebSocketManager | null = null;
  private ws: WebSocket | null = null;
  private subscribers: Map<string, Set<(message: WebSocketMessage) => void>> = new Map();
  private globalSubscribers: Set<(message: WebSocketMessage) => void> = new Set();
  private connectionStatus: ConnectionStatus = {
    connected: false,
    connecting: false,
    reconnectAttempts: 0,
  };
  private options: Required<UseWebSocketOptions>;
  private heartbeatTimer: NodeJS.Timeout | null = null;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private connectionStatusSubscribers: Set<(status: ConnectionStatus) => void> = new Set();

  private constructor(options: UseWebSocketOptions = {}) {
    this.options = {
      autoReconnect: options.autoReconnect ?? true,
      reconnectInterval: options.reconnectInterval ?? 3000,
      maxReconnectAttempts: options.maxReconnectAttempts ?? 10,
      heartbeatInterval: options.heartbeatInterval ?? 30000,
      enableMockMode: options.enableMockMode ?? false,
    };
  }

  static getInstance(options?: UseWebSocketOptions): WebSocketManager {
    if (!WebSocketManager.instance) {
      WebSocketManager.instance = new WebSocketManager(options);
    }
    return WebSocketManager.instance;
  }

  private updateConnectionStatus(updates: Partial<ConnectionStatus>) {
    this.connectionStatus = { ...this.connectionStatus, ...updates };
    this.connectionStatusSubscribers.forEach((callback) => callback(this.connectionStatus));
  }

  private startHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
    }

    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.sendMessage({ type: 'ping', data: null });
      }
    }, this.options.heartbeatInterval);
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private scheduleReconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }

    if (this.connectionStatus.reconnectAttempts >= this.options.maxReconnectAttempts) {
      this.updateConnectionStatus({
        connected: false,
        connecting: false,
        error: 'Max reconnection attempts reached',
      });
      return;
    }

    this.reconnectTimer = setTimeout(() => {
      this.updateConnectionStatus({
        reconnectAttempts: this.connectionStatus.reconnectAttempts + 1,
      });
      this.connect();
    }, this.options.reconnectInterval);
  }

  private connect() {
    if (this.options.enableMockMode) {
      this.updateConnectionStatus({
        connected: true,
        connecting: false,
        lastConnected: new Date(),
        reconnectAttempts: 0,
      });
      this.startMockMode();
      return;
    }

    if (this.ws?.readyState === WebSocket.OPEN || this.connectionStatus.connecting) {
      return;
    }

    this.updateConnectionStatus({ connecting: true });

    try {
      const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8085';
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.updateConnectionStatus({
          connected: true,
          connecting: false,
          error: undefined,
          lastConnected: new Date(),
          reconnectAttempts: 0,
        });
        this.startHeartbeat();
      };

      this.ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          this.handleMessage(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onclose = (event) => {
        this.updateConnectionStatus({
          connected: false,
          connecting: false,
          error: event.reason || 'Connection closed',
        });
        this.stopHeartbeat();

        if (this.options.autoReconnect && !event.wasClean) {
          this.scheduleReconnect();
        }
      };

      this.ws.onerror = (error) => {
        this.updateConnectionStatus({
          connected: false,
          connecting: false,
          error: 'WebSocket connection error',
        });
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      this.updateConnectionStatus({
        connected: false,
        connecting: false,
        error: error instanceof Error ? error.message : 'Unknown connection error',
      });

      if (this.options.autoReconnect) {
        this.scheduleReconnect();
      }
    }
  }

  private startMockMode() {
    // Simulate WebSocket messages for development
    const mockInterval = setInterval(() => {
      if (!this.connectionStatus.connected) {
        clearInterval(mockInterval);
        return;
      }

      // Simulate random bid updates
      if (Math.random() > 0.7) {
        const mockMessage: WebSocketMessage = {
          type: 'bid_update',
          auctionId: '1',
          data: {
            auctionId: '1',
            bid: {
              id: Date.now().toString(),
              auctionId: '1',
              userId: 'mock-user',
              amount: 275 + Math.random() * 50,
              timestamp: new Date().toISOString(),
              user: {
                id: 'mock-user',
                name: 'Mock Bidder',
                email: 'mock@example.com',
                isSeller: false,
                rating: 0,
                totalSales: 0,
                createdAt: new Date().toISOString(),
              },
            },
            isWinning: true,
            newCurrentBid: 275 + Math.random() * 50,
            totalBids: Math.floor(Math.random() * 20) + 1,
          },
          timestamp: new Date().toISOString(),
        };
        this.handleMessage(mockMessage);
      }
    }, 5000);
  }

  private handleMessage(message: WebSocketMessage) {
    // Handle ping/pong
    if (message.type === 'ping') {
      this.sendMessage({ type: 'pong', data: null });
      return;
    }

    // Notify global subscribers
    this.globalSubscribers.forEach((callback) => callback(message));

    // Notify auction-specific subscribers
    if (message.auctionId) {
      const auctionSubscribers = this.subscribers.get(message.auctionId);
      if (auctionSubscribers) {
        auctionSubscribers.forEach((callback) => callback(message));
      }
    }
  }

  public sendMessage(message: Omit<WebSocketMessage, 'timestamp'>) {
    const fullMessage: WebSocketMessage = {
      ...message,
      timestamp: new Date().toISOString(),
    };

    if (this.options.enableMockMode) {
      // In mock mode, just log the message
      console.log('Mock WebSocket send:', fullMessage);
      return;
    }

    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(fullMessage));
    } else {
      console.warn('WebSocket not connected, cannot send message:', fullMessage);
    }
  }

  public subscribe(auctionId: string, callback: (message: WebSocketMessage) => void) {
    if (!this.subscribers.has(auctionId)) {
      this.subscribers.set(auctionId, new Set());
    }

    this.subscribers.get(auctionId)!.add(callback);

    // Subscribe to the auction on the server
    this.sendMessage({
      type: 'auction_update',
      auctionId,
      data: { action: 'subscribe' },
    });

    // Return unsubscribe function
    return () => {
      const subscribers = this.subscribers.get(auctionId);
      if (subscribers) {
        subscribers.delete(callback);
        if (subscribers.size === 0) {
          this.subscribers.delete(auctionId);
          // Unsubscribe from the server
          this.sendMessage({
            type: 'auction_update',
            auctionId,
            data: { action: 'unsubscribe' },
          });
        }
      }
    };
  }

  public subscribeGlobal(callback: (message: WebSocketMessage) => void) {
    this.globalSubscribers.add(callback);
    return () => {
      this.globalSubscribers.delete(callback);
    };
  }

  public subscribeToConnectionStatus(callback: (status: ConnectionStatus) => void) {
    callback(this.connectionStatus);
    this.connectionStatusSubscribers.add(callback);
    return () => {
      this.connectionStatusSubscribers.delete(callback);
    };
  }

  public getConnectionStatus(): ConnectionStatus {
    return this.connectionStatus;
  }

  public connect() {
    this.connect();
  }

  public disconnect() {
    this.options.autoReconnect = false;
    this.stopHeartbeat();

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.updateConnectionStatus({
      connected: false,
      connecting: false,
      reconnectAttempts: 0,
    });
  }

  public destroy() {
    this.disconnect();
    this.subscribers.clear();
    this.globalSubscribers.clear();
    this.connectionStatusSubscribers.clear();
    WebSocketManager.instance = null;
  }
}

export function useWebSocket(options: UseWebSocketOptions = {}): UseWebSocketReturn {
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>({
    connected: false,
    connecting: false,
    reconnectAttempts: 0,
  });
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);

  const managerRef = useRef<WebSocketManager>();
  const unsubscribeGlobalRef = useRef<(() => void) | null>(null);
  const unsubscribeStatusRef = useRef<(() => void) | null>(null);

  // Initialize or get WebSocket manager
  useEffect(() => {
    if (!managerRef.current) {
      managerRef.current = WebSocketManager.getInstance(options);
      managerRef.current.connect();
    }
  }, [options]);

  // Subscribe to connection status updates
  useEffect(() => {
    if (managerRef.current) {
      unsubscribeStatusRef.current =
        managerRef.current.subscribeToConnectionStatus(setConnectionStatus);
    }

    return () => {
      if (unsubscribeStatusRef.current) {
        unsubscribeStatusRef.current();
      }
    };
  }, []);

  // Subscribe to global messages
  useEffect(() => {
    if (managerRef.current) {
      unsubscribeGlobalRef.current = managerRef.current.subscribeGlobal(setLastMessage);
    }

    return () => {
      if (unsubscribeGlobalRef.current) {
        unsubscribeGlobalRef.current();
      }
    };
  }, []);

  // Subscribe to specific auction
  const subscribeToAuction = useCallback((auctionId: string) => {
    if (!managerRef.current) return () => {};

    return managerRef.current.subscribe(auctionId, setLastMessage);
  }, []);

  // Unsubscribe from specific auction
  const unsubscribeFromAuction = useCallback((auctionId: string) => {
    // Note: The actual unsubscribing is handled by the unsubscribe function
    // returned from subscribeToAuction. This is just for convenience.
    if (managerRef.current) {
      managerRef.current.sendMessage({
        type: 'auction_update',
        auctionId,
        data: { action: 'unsubscribe' },
      });
    }
  }, []);

  // Send message
  const sendMessage = useCallback((message: Omit<WebSocketMessage, 'timestamp'>) => {
    if (managerRef.current) {
      managerRef.current.sendMessage(message);
    }
  }, []);

  // Calculate connection quality
  const connectionQuality = useCallback((): 'excellent' | 'good' | 'poor' | 'disconnected' => {
    if (!connectionStatus.connected) return 'disconnected';
    if (connectionStatus.reconnectAttempts > 5) return 'poor';
    if (connectionStatus.reconnectAttempts > 2) return 'good';
    return 'excellent';
  }, [connectionStatus]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      // Don't destroy the manager here as it might be used by other components
      // Just cleanup this component's subscriptions
      if (unsubscribeGlobalRef.current) {
        unsubscribeGlobalRef.current();
      }
      if (unsubscribeStatusRef.current) {
        unsubscribeStatusRef.current();
      }
    };
  }, []);

  return {
    connectionStatus,
    subscribeToAuction,
    unsubscribeFromAuction,
    sendMessage,
    lastMessage,
    connectionQuality: connectionQuality(),
  };
}

// Utility functions for working with WebSocket messages
export const isAuctionUpdate = (
  message: WebSocketMessage
): message is WebSocketMessage & { data: AuctionUpdateData } => {
  return message.type === 'auction_update';
};

export const isBidUpdate = (
  message: WebSocketMessage
): message is WebSocketMessage & { data: BidUpdateData } => {
  return message.type === 'bid_update';
};

export const isCountdownUpdate = (
  message: WebSocketMessage
): message is WebSocketMessage & { data: CountdownUpdateData } => {
  return message.type === 'countdown_update';
};

export const isStatusChange = (
  message: WebSocketMessage
): message is WebSocketMessage & { data: StatusChangeData } => {
  return message.type === 'status_change';
};

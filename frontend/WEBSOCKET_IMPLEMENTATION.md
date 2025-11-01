# WebSocket Real-Time Auction Implementation

## Overview

This implementation provides a comprehensive real-time auction system using
WebSockets in the Next.js frontend. It replaces inefficient polling with instant
bid updates, synchronized countdown timers, and live auction status changes.

## Architecture

### Core Components

1. **WebSocket Connection Manager** (`/src/hooks/useWebSocket.ts`)
   - Singleton pattern for managing WebSocket connections
   - Automatic reconnection with exponential backoff
   - Connection quality monitoring
   - Mock mode for development
   - Heartbeat/ping-pong mechanism

2. **Auction Real-time Hook** (`/src/hooks/useAuctionRealtime.ts`)
   - High-level hook for auction-specific real-time updates
   - Automatic fallback to polling when WebSocket fails
   - Bid updates, countdown synchronization, and status changes
   - Optimistic updates with server confirmation

3. **Real-time Components**
   - **Countdown Timer** (`/src/components/auction/realtime-countdown.tsx`)
   - **Bid History** (`/src/components/auction/bid-history.tsx`)
   - **Bid Button** (`/src/components/auction/bid-button.tsx`)
   - **Connection Status** (`/src/components/auction/connection-status.tsx`)

4. **Notification System** (`/src/components/ui/toast.tsx`)
   - Real-time toast notifications for bid updates
   - Custom toast variants for different auction events
   - Context-based notification management

## Key Features

### 1. Real-time Bid Updates

- Instant bid propagation across all connected clients
- Bid history with user identification
- Winning bid highlighting
- Animated new bid notifications

### 2. Synchronized Countdown Timers

- Server-synchronized countdown to prevent timing issues
- "Ending Soon" warnings (5 minutes)
- Cross-client time synchronization
- Graceful handling of clock drift

### 3. Connection Management

- Automatic reconnection with configurable attempts
- Connection quality indicators (excellent/good/poor/disconnected)
- Fallback to HTTP polling when WebSocket fails
- Connection status visibility to users

### 4. Auction Status Management

- Real-time status transitions (scheduled → active → ended)
- Winner announcements
- Live stream integration indicators
- Anti-snipe protection awareness

### 5. Performance Optimizations

- Singleton WebSocket connection to prevent multiple connections
- Efficient message routing to specific auction subscribers
- Component cleanup to prevent memory leaks
- Debounced updates to prevent excessive re-renders

## Usage Examples

### Basic Real-time Auction Page

```tsx
'use client';

import { useAuctionRealtime } from '@/hooks/useAuctionRealtime';
import { RealtimeCountdown } from '@/components/auction/realtime-countdown';
import { BidButton } from '@/components/auction/bid-button';
import { BidHistory } from '@/components/auction/bid-history';
import { ConnectionStatus } from '@/components/auction/connection-status';

export default function AuctionPage({ params }: { params: { id: string } }) {
  const {
    auction,
    isConnected,
    remainingTime,
    isEndingSoon,
    placeBid,
    connectionStatus,
  } = useAuctionRealtime({
    auctionId: params.id,
    onBidUpdate: (data) => {
      console.log('New bid:', data.bid);
    },
    onStatusChange: (data) => {
      console.log('Status changed:', data.newStatus);
    },
  });

  if (!auction) return <div>Loading...</div>;

  return (
    <div>
      <ConnectionStatus
        connected={isConnected}
        connecting={connectionStatus.connecting}
        reconnectAttempts={connectionStatus.reconnectAttempts}
        connectionQuality="excellent"
        variant="detailed"
      />

      <RealtimeCountdown
        endTime={auction.endTime}
        onEndingSoon={(remaining) => console.log('Ending soon!')}
        onEnd={() => console.log('Auction ended!')}
      />

      <BidButton
        currentBid={auction.currentBid}
        minBidIncrement={auction.minBidIncrement}
        onPlaceBid={placeBid}
        auctionId={auction.id}
      />

      <BidHistory
        bids={auction.bids}
        currentUserId="current-user-id"
        realtime={true}
      />
    </div>
  );
}
```

### WebSocket Hook Direct Usage

```tsx
import { useWebSocket } from '@/hooks/useWebSocket';

export function RealtimeComponent() {
  const {
    connectionStatus,
    subscribeToAuction,
    sendMessage,
    lastMessage,
    connectionQuality,
  } = useWebSocket({
    enableMockMode: process.env.NODE_ENV === 'development',
    autoReconnect: true,
    maxReconnectAttempts: 10,
  });

  useEffect(() => {
    const unsubscribe = subscribeToAuction('auction-id', (message) => {
      if (message.type === 'bid_update') {
        // Handle bid update
      }
    });

    return unsubscribe;
  }, [subscribeToAuction]);

  return (
    <div>
      <div>
        Connection: {connectionStatus.connected ? 'Connected' : 'Disconnected'}
      </div>
      <div>Quality: {connectionQuality}</div>
    </div>
  );
}
```

## WebSocket Message Types

### Bid Update

```typescript
{
  type: 'bid_update',
  auctionId: 'auction-123',
  data: {
    auctionId: 'auction-123',
    bid: {
      id: 'bid-456',
      userId: 'user-789',
      amount: 250.00,
      timestamp: '2024-01-15T10:30:00Z',
      user: { /* user data */ }
    },
    isWinning: true,
    newCurrentBid: 250.00,
    totalBids: 15
  },
  timestamp: '2024-01-15T10:30:00Z'
}
```

### Countdown Update

```typescript
{
  type: 'countdown_update',
  auctionId: 'auction-123',
  data: {
    auctionId: 'auction-123',
    remainingTime: 120,
    isEndingSoon: true
  },
  timestamp: '2024-01-15T10:30:00Z'
}
```

### Status Change

```typescript
{
  type: 'status_change',
  auctionId: 'auction-123',
  data: {
    auctionId: 'auction-123',
    oldStatus: 'active',
    newStatus: 'ended',
    winner: {
      id: 'bid-456',
      user: { /* user data */ },
      amount: 300.00
    }
  },
  timestamp: '2024-01-15T10:30:00Z'
}
```

## Configuration

### Environment Variables

```env
# WebSocket server URL
NEXT_PUBLIC_WS_URL=ws://localhost:8085

# API URL for fallback polling
NEXT_PUBLIC_API_URL=http://localhost:8080

# Development mode (enables mock WebSocket)
NODE_ENV=development
```

### WebSocket Options

```typescript
const wsOptions = {
  autoReconnect: true, // Auto-reconnect on disconnection
  reconnectInterval: 3000, // Reconnection delay (ms)
  maxReconnectAttempts: 10, // Maximum reconnection attempts
  heartbeatInterval: 30000, // Ping interval (ms)
  enableMockMode: false, // Enable mock mode for testing
};
```

## Development Features

### Mock Mode

- Simulates WebSocket messages for development
- Generates random bid updates every 5 seconds
- No backend WebSocket server required
- Enable with `enableMockMode: true`

### Connection Status Indicators

- **Excellent**: Stable connection, < 2 reconnection attempts
- **Good**: Minor issues, 2-5 reconnection attempts
- **Poor**: Unstable connection, > 5 reconnection attempts
- **Disconnected**: No connection

### Fallback Polling

- Automatically switches to HTTP polling when WebSocket fails
- Configurable polling interval (default: 5 seconds)
- Maintains user experience during connection issues

## Performance Considerations

1. **Connection Pooling**: Singleton WebSocket prevents multiple connections
2. **Message Routing**: Efficient subscription-based message delivery
3. **Memory Management**: Proper cleanup on component unmount
4. **Debouncing**: Prevents excessive re-renders during rapid updates
5. **Lazy Loading**: Components only subscribe when needed

## Security

1. **Authentication**: WebSocket connections require valid JWT tokens
2. **Authorization**: Users can only subscribe to auctions they have access to
3. **Input Validation**: All incoming messages are validated
4. **Rate Limiting**: Client-side rate limiting for bid submissions
5. **HTTPS**: WebSocket connections use secure WebSocket protocol (wss://) in
   production

## Testing

### Mock Mode Testing

```typescript
// Enable mock mode for testing
const { auction } = useAuctionRealtime({
  auctionId: 'test-auction',
  enableRealtime: true,
  fallbackToPolling: false,
});

// Mock WebSocket will simulate bid updates
```

### Connection Testing

```typescript
// Test connection quality
const { connectionQuality, connectionStatus } = useWebSocket();

// Monitor reconnection attempts
console.log('Reconnection attempts:', connectionStatus.reconnectAttempts);
```

## Troubleshooting

### Common Issues

1. **WebSocket Connection Fails**
   - Check `NEXT_PUBLIC_WS_URL` environment variable
   - Verify WebSocket server is running
   - Check network connectivity and firewall settings

2. **Frequent Disconnections**
   - Monitor connection quality indicator
   - Check network stability
   - Verify server-side WebSocket configuration

3. **Missing Real-time Updates**
   - Ensure auction subscription is active
   - Check WebSocket message routing
   - Verify component is properly mounted

4. **Memory Leaks**
   - Ensure proper cleanup in useEffect
   - Check for multiple WebSocket connections
   - Verify subscription cleanup

### Debug Mode

Enable debug logging:

```typescript
const ws = useWebSocket({
  enableMockMode: true,
  // Additional debug options
});
```

Monitor WebSocket events in browser DevTools:

```javascript
// In browser console
window.websocketManager = WebSocketManager.getInstance();
```

## Future Enhancements

1. **Offline Support**: Service worker integration for offline bidding
2. **Push Notifications**: Mobile push notifications for bid updates
3. **Analytics**: Real-time bidding analytics and insights
4. **Multi-language**: Internationalization support
5. **Accessibility**: Enhanced screen reader support for live updates
6. **Performance**: WebWorker for message processing
7. **Security**: Enhanced authentication and authorization

## Conclusion

This WebSocket implementation provides a robust, scalable real-time auction
system that significantly improves user experience over traditional
polling-based approaches. The architecture is designed for maintainability,
performance, and extensibility while providing fallback mechanisms to ensure
reliability.

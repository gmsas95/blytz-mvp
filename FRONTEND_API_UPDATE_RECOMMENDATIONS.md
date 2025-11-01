# Frontend API Update Recommendations

## 🎯 Executive Summary

This document provides specific code changes needed to fix the frontend-backend API integration issues. The changes are organized by priority and impact.

---

## 🚨 Phase 1: Critical API Path Fixes

### 1. Update RemoteApiAdapter in `/frontend/src/lib/api-adapter.ts`

**Current Issues:**
- Missing `/api` prefix on most endpoints
- Wrong auction bidding endpoint
- Missing authentication paths

**Required Changes:**

```typescript
// === PRODUCTS ===
// FROM:
return this.fetchApi(`/products?${params}`);
// TO:
return this.fetchApi(`/api/v1/products?${params}`);

// FROM:
return this.fetchApi(`/products/${id}`);
// TO:
return this.fetchApi(`/api/v1/products/${id}`);

// FROM:
return this.fetchApi('/products/featured');
// TO:
return this.fetchApi('/api/v1/products/featured');

// === AUCTIONS ===
// FROM:
return this.fetchApi(`/auctions?${params}`);
// TO:
return this.fetchApi(`/api/v1/auctions?${params}`);

// FROM:
return this.fetchApi(`/auctions/${id}`);
// TO:
return this.fetchApi(`/api/v1/auctions/${id}`);

// FROM:
return this.fetchApi(`/auctions/${auctionId}/bid`, {
  method: 'POST',
  body: JSON.stringify({ amount }),
});
// TO:
return this.fetchApi(`/api/v1/auctions/${auctionId}/bids`, {
  method: 'POST',
  body: JSON.stringify({ amount }),
});

// FROM:
return this.fetchApi('/auctions/active');
// TO:
return this.fetchApi('/api/v1/auctions/active');

// === AUTHENTICATION ===
// FROM:
return this.fetchApi('/auth/login', {
  method: 'POST',
  body: JSON.stringify({ email, password }),
});
// TO:
return this.fetchApi('/api/auth/login', {
  method: 'POST',
  body: JSON.stringify({ email, password }),
});

// FROM:
return this.fetchApi('/auth/register', {
  method: 'POST',
  body: JSON.stringify(userData),
});
// TO:
return this.fetchApi('/api/auth/register', {
  method: 'POST',
  body: JSON.stringify(userData),
});

// FROM:
return this.fetchApi('/auth/logout', { method: 'POST' });
// TO:
return this.fetchApi('/api/auth/logout', { method: 'POST' });

// FROM:
return this.fetchApi('/auth/me');
// TO:
return this.fetchApi('/api/auth/me');

// === PAYMENTS ===
// FROM:
return this.fetchApi('/payments/methods');
// TO:
return this.fetchApi('/api/v1/payments/methods');

// FROM:
return this.fetchApi('/payments/create', {
  method: 'POST',
  body: JSON.stringify(paymentRequest),
});
// TO:
return this.fetchApi('/api/v1/payments/create', {
  method: 'POST',
  body: JSON.stringify(paymentRequest),
});

// FROM:
return this.fetchApi(`/payments/${paymentId}/status`);
// TO:
return this.fetchApi(`/api/v1/payments/${paymentId}/status`);
```

---

## 🔧 Phase 2: Missing Backend Endpoints Implementation

### 2. Cart Endpoints (Add to Order Service)

**Backend Implementation Needed:**

```go
// In services/order-service/internal/handlers/cart.go
package handlers

type CartHandler struct {
    db *gorm.DB
}

func (h *CartHandler) GetCart(c *gin.Context) {
    userID := c.GetString("userID")
    
    var cart models.Cart
    if err := h.db.Where("user_id = ?", userID).Preload("Items.Product").First(&cart).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Create empty cart
            cart = models.Cart{
                UserID:    userID,
                Items:     []models.CartItem{},
                Total:     0,
                ItemCount: 0,
            }
            h.db.Create(&cart)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }
    
    c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
    userID := c.GetString("userID")
    
    var req struct {
        ProductID string `json:"productId"`
        Quantity  int    `json:"quantity"`
        AuctionID string `json:"auctionId,omitempty"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Implementation details...
    c.JSON(http.StatusOK, gin.H{"data": cart})
}

// Add routes: /api/v1/cart/, /api/v1/cart/add, etc.
```

**Frontend Updates:**

```typescript
// Add to RemoteApiAdapter class
async getCart(): Promise<ApiResponse<Cart>> {
  return this.fetchApi('/api/v1/cart/');
}

async addToCart(
  productId: string,
  quantity: number,
  auctionId?: string
): Promise<ApiResponse<Cart>> {
  return this.fetchApi('/api/v1/cart/add', {
    method: 'POST',
    body: JSON.stringify({ productId, quantity, auctionId }),
  });
}

async removeFromCart(itemId: string): Promise<ApiResponse<Cart>> {
  return this.fetchApi(`/api/v1/cart/remove/${itemId}`, { method: 'DELETE' });
}

async updateCartItemQuantity(itemId: string, quantity: number): Promise<ApiResponse<Cart>> {
  return this.fetchApi(`/api/v1/cart/update/${itemId}`, {
    method: 'PUT',
    body: JSON.stringify({ quantity }),
  });
}

async clearCart(): Promise<ApiResponse<Cart>> {
  return this.fetchApi('/api/v1/cart/clear', { method: 'DELETE' });
}
```

---

### 3. LiveKit Token Endpoint

**Backend Implementation (Add to Auction Service):**

```go
// In services/auction-service/internal/handlers/livekit.go
package handlers

import (
    "github.com/livekit/protocol/auth"
    lksdk "github.com/livekit/server-sdk-go"
)

type LiveKitHandler struct {
    roomService *lksdk.RoomServiceClient
    apiKey      string
    apiSecret   string
}

func (h *LiveKitHandler) GenerateToken(c *gin.Context) {
    room := c.Query("room")
    role := c.Query("role") // "viewer" or "broadcaster"
    
    if room == "" || role == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "room and role required"})
        return
    }
    
    at := auth.NewAccessToken(h.apiKey, h.apiSecret)
    
    if role == "broadcaster" {
        grant := &auth.VideoGrant{
            RoomJoin: true,
            Room:     room,
            CanPublish: true,
            CanSubscribe: true,
        }
        at.AddGrant(grant)
    } else {
        grant := &auth.VideoGrant{
            RoomJoin: true,
            Room:     room,
            CanPublish: false,
            CanSubscribe: true,
        }
        at.AddGrant(grant)
    }
    
    token, err := at.ToJWT()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "token": token,
            "room":  room,
            "role":  role,
        },
    })
}

// Add route: GET /api/v1/livekit/token
```

**Frontend Updates:**

```typescript
// Update LiveKit components to use correct endpoint
// In frontend-demo/src/components/LiveKitViewer.tsx
const response = await fetch(`/api/v1/livekit/token?room=${auctionId}&role=viewer`)

// In frontend-seller/src/components/LiveKitBroadcaster.tsx  
const response = await fetch(`/api/v1/livekit/token?room=${auctionId}&role=broadcaster`)
```

---

### 4. Livestream Endpoints

**Backend Implementation (New Livestream Service):**

```go
// Create new service: services/livestream-service/
// Endpoints:
// GET /api/v1/livestreams/
// GET /api/v1/livestreams/{id}
// GET /api/v1/livestreams/active
```

**Frontend Updates:**

```typescript
// Add to RemoteApiAdapter class
async getLivestreams(
  filter?: LivestreamFilter
): Promise<ApiResponse<PaginatedResponse<Livestream>>> {
  const params = new URLSearchParams();
  if (filter?.status) params.append('status', filter.status);

  return this.fetchApi(`/api/v1/livestreams?${params}`);
}

async getLivestream(id: string): Promise<ApiResponse<Livestream>> {
  return this.fetchApi(`/api/v1/livestreams/${id}`);
}

async getActiveLivestreams(): Promise<ApiResponse<Livestream[]>> {
  return this.fetchApi('/api/v1/livestreams/active');
}
```

---

## 🔄 Phase 3: Environment Configuration Updates

### 5. Update Environment Variables

**Frontend `.env.local`:**

```bash
# Update API base URL to point to gateway
NEXT_PUBLIC_API_URL=http://localhost:8080

# LiveKit configuration
NEXT_PUBLIC_LIVEKIT_URL=wss://livekit.blytz.app

# Mode configuration
MODE=remote  # Change from 'mock' to 'remote' for production
```

---

## 🧪 Phase 4: Testing & Validation

### 6. Integration Tests

**Create test file: `/frontend/src/lib/__tests__/api-adapter.test.ts`**

```typescript
import { RemoteApiAdapter } from '../api-adapter'

describe('RemoteApiAdapter', () => {
  let adapter: RemoteApiAdapter
  
  beforeEach(() => {
    adapter = new RemoteApiAdapter('http://localhost:8080')
  })
  
  test('should use correct API paths', () => {
    // Test that all endpoints use correct prefixes
    expect(adapter.baseUrl).toBe('http://localhost:8080')
    // Add more specific tests...
  })
})
```

---

## 📋 Implementation Checklist

### Frontend Changes:
- [ ] Update RemoteApiAdapter endpoints (Phase 1)
- [ ] Add cart API methods (Phase 2)
- [ ] Update LiveKit token endpoints (Phase 2)
- [ ] Add livestream API methods (Phase 2)
- [ ] Update environment variables (Phase 3)
- [ ] Run integration tests (Phase 4)

### Backend Changes:
- [ ] Implement cart endpoints in order-service
- [ ] Implement LiveKit token endpoint in auction-service
- [ ] Create livestream-service with required endpoints
- [ ] Update OpenAPI specs
- [ ] Add authentication middleware to new endpoints

---

## ⚡ Quick Start Fix (2-Hour Solution)

For immediate testing, apply only Phase 1 changes:

```bash
# In frontend/src/lib/api-adapter.ts
# Replace all endpoint paths with correct backend paths
# This will fix: auth, products, auctions, payments
# Remaining: cart, livekit, livestreams (need backend implementation)
```

This will restore basic functionality (login, browse, bid, pay) while cart and streaming features are implemented.

---

## 🚀 Deployment Strategy

1. **Stage 1:** Deploy frontend API path fixes
2. **Stage 2:** Deploy backend cart endpoints
3. **Stage 3:** Deploy LiveKit token endpoint
4. **Stage 4:** Deploy livestream service
5. **Stage 5:** Full integration testing

Each stage can be deployed independently to minimize risk.

---

**Estimated Timeline:**
- Phase 1: 2 hours
- Phase 2: 1-2 days  
- Phase 3: 30 minutes
- Phase 4: 4 hours
- **Total: 2-3 days for complete fix**
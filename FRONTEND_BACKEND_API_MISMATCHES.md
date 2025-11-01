# Frontend-Backend API Integration Analysis

## 🚨 Critical Mismatches Found

### 1. Authentication Endpoints

**Frontend expects:**
- `POST /auth/login`
- `POST /auth/register` 
- `POST /auth/logout`
- `GET /auth/me`

**Backend provides:**
- `POST /api/auth/login`
- `POST /api/auth/register`
- `POST /api/auth/logout`
- `GET /api/auth/me`

**Impact:** 🔥 **CRITICAL** - User authentication completely broken

---

### 2. Product Endpoints

**Frontend expects:**
- `GET /products`
- `GET /products/{id}`
- `GET /products/featured`

**Backend provides:**
- `GET /api/v1/products`
- `GET /api/v1/products/{id}`
- `GET /api/v1/products/featured`

**Impact:** 🔥 **CRITICAL** - Product browsing broken

---

### 3. Auction Endpoints

**Frontend expects:**
- `GET /auctions`
- `GET /auctions/{id}`
- `POST /auctions/{id}/bid` ⚠️
- `GET /auctions/active`

**Backend provides:**
- `GET /api/v1/auctions`
- `GET /api/v1/auctions/{id}`
- `POST /api/v1/auctions/{id}/bids` ✅
- `GET /api/v1/auctions/active`

**Impact:** 🔥 **CRITICAL** - Auction functionality broken (bidding endpoint mismatch)

---

### 4. Cart Endpoints

**Frontend expects:**
- `GET /cart`
- `POST /cart/add`
- `DELETE /cart/remove/{id}`
- `PUT /cart/update/{id}`
- `DELETE /cart/clear`

**Backend provides:**
- ❌ **NO CART ENDPOINTS** - Only `/api/v1/orders/`

**Impact:** 🔥 **CRITICAL** - Shopping cart completely missing

---

### 5. Payment Endpoints

**Frontend expects:**
- `GET /payments/methods`
- `POST /payments/create`
- `GET /payments/{id}/status`

**Backend provides:**
- `GET /api/v1/payments/methods`
- `POST /api/v1/payments/create`
- `GET /api/v1/payments/{id}/status`

**Impact:** 🔥 **CRITICAL** - Payment processing broken

---

### 6. LiveKit Token Endpoints

**Frontend expects:**
- `GET /api/livekit/token?room={id}&role=viewer|broadcaster`

**Backend provides:**
- ❌ **NO LIVEKIT TOKEN ENDPOINTS**

**Impact:** 🔥 **CRITICAL** - Live streaming completely broken

---

### 7. Livestream Endpoints

**Frontend expects:**
- `GET /livestreams`
- `GET /livestreams/{id}`
- `GET /livestreams/active`

**Backend provides:**
- ❌ **NO LIVESTREAM ENDPOINTS** - Only auction-related streaming

**Impact:** 🔥 **CRITICAL** - Livestream discovery broken

---

## 📊 Mismatch Summary

| Service | Frontend Path | Backend Path | Status |
|---------|---------------|--------------|---------|
| Auth | `/auth/*` | `/api/auth/*` | ❌ Prefix mismatch |
| Products | `/products/*` | `/api/v1/products/*` | ❌ Prefix mismatch |
| Auctions | `/auctions/{id}/bid` | `/api/v1/auctions/{id}/bids` | ❌ Endpoint name |
| Cart | `/cart/*` | **NONE** | ❌ Missing endpoints |
| Payments | `/payments/*` | `/api/v1/payments/*` | ❌ Prefix mismatch |
| LiveKit | `/api/livekit/token` | **NONE** | ❌ Missing endpoints |
| Livestreams | `/livestreams/*` | **NONE** | ❌ Missing endpoints |

---

## 🔧 Required Fixes

### Phase 1: Critical Path Fixes (Authentication & Core Features)

1. **Update Frontend API Adapter** (`/frontend/src/lib/api-adapter.ts`)
   ```typescript
   // Change all endpoints to match backend
   '/auth/login' → '/api/auth/login'
   '/products' → '/api/v1/products'
   '/auctions/{id}/bid' → '/api/v1/auctions/{id}/bids'
   '/payments' → '/api/v1/payments'
   ```

2. **Implement Missing Backend Endpoints**
   - Cart management endpoints in order-service
   - LiveKit token generation endpoint
   - Livestream discovery endpoints

### Phase 2: Feature Implementation

3. **Add Cart Service** (New microservice or extend order-service)
   - `GET /api/v1/cart/`
   - `POST /api/v1/cart/add`
   - `DELETE /api/v1/cart/remove/{id}`
   - `PUT /api/v1/cart/update/{id}`
   - `DELETE /api/v1/cart/clear`

4. **Add LiveKit Token Service** (Extend auction-service or new service)
   - `GET /api/v1/livekit/token?room={id}&role={viewer|broadcaster}`

5. **Add Livestream Service** (New microservice)
   - `GET /api/v1/livestreams/`
   - `GET /api/v1/livestreams/{id}`
   - `GET /api/v1/livestreams/active`

---

## 🚨 Immediate Action Required

The frontend applications **cannot communicate** with backend services due to these mismatches. Users will experience:

- ❌ Unable to login/register
- ❌ Unable to browse products
- ❌ Unable to place bids
- ❌ Unable to use shopping cart
- ❌ Unable to make payments
- ❌ Unable to access live streams

**Estimated Development Time:** 2-3 days for critical fixes, 1-2 weeks for complete implementation.

---

## 📝 Recommended Implementation Order

1. **Fix API path prefixes** (2 hours)
2. **Fix auction bidding endpoint** (1 hour)
3. **Implement cart endpoints** (1 day)
4. **Implement LiveKit token endpoint** (4 hours)
5. **Implement livestream endpoints** (1 day)
6. **Test all integrations** (4 hours)

---

## 🔍 Testing Strategy

After fixes:
1. Unit test each API adapter method
2. Integration test frontend-backend communication
3. End-to-end test user flows (login → browse → bid → pay)
4. Load test live streaming functionality

---

**Status:** 🚨 **BLOCKING** - Frontend completely non-functional with current backend
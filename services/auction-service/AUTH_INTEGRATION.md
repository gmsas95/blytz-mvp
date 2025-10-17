# Auction Service Authentication Integration

This document explains how authentication has been integrated into the auction service.

## Overview

The auction service now uses JWT-based authentication provided by the auth service. Authentication is required for operations that modify data (create, update, delete, bid), while read operations remain public.

## Authentication Flow

1. **User Registration/Login**: Users authenticate with the auth service
2. **Token Acquisition**: Auth service returns JWT token
3. **API Requests**: Include token in `Authorization: Bearer <token>` header
4. **Validation**: Auction service validates token with auth service
5. **User Context**: Authenticated user ID available in request context

## Protected vs Public Endpoints

### Public Endpoints (No Authentication Required)
- `GET /api/v1/auctions` - List auctions
- `GET /api/v1/auctions/:auction_id` - Get auction details
- `GET /api/v1/auctions/:auction_id/status` - Get auction status
- `GET /api/v1/auctions/:auction_id/bids` - Get bids for auction

### Protected Endpoints (Authentication Required)
- `POST /api/v1/auctions` - Create new auction
- `PUT /api/v1/auctions/:auction_id` - Update auction
- `DELETE /api/v1/auctions/:auction_id` - Delete auction
- `POST /api/v1/auctions/:auction_id/bids` - Place bid

## Implementation Details

### Router Configuration
```go
// Initialize auth client
authClient := auth.NewAuthClient("http://auth-service:8084")

// Public routes (no auth required)
auctions := api.Group("/auctions")
{
    auctions.GET("", auctionHandler.ListAuctions)
    auctions.GET("/:auction_id", auctionHandler.GetAuction)
    auctions.GET("/:auction_id/status", auctionHandler.GetAuctionStatus)
    auctions.GET("/:auction_id/bids", auctionHandler.GetBids)
}

// Protected routes (auth required)
protectedAuctions := api.Group("/auctions")
protectedAuctions.Use(auth.GinAuthMiddleware(authClient))
{
    protectedAuctions.POST("", auctionHandler.CreateAuction)
    protectedAuctions.PUT("/:auction_id", auctionHandler.UpdateAuction)
    protectedAuctions.DELETE("/:auction_id", auctionHandler.DeleteAuction)
    protectedAuctions.POST("/:auction_id/bids", auctionHandler.PlaceBid)
}
```

### Handler Implementation
```go
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
    // Get authenticated user ID from context
    userID, exists := c.Get("userID")
    if !exists {
        utils.ErrorResponse(c, errors.AuthenticationError("NO_USER_ID", "User ID not found"))
        return
    }

    // Use userID for creating auction
    auction, err := h.auctionService.CreateAuction(c.Request.Context(), userID.(string), &req)
    // ... rest of implementation
}
```

## API Usage Examples

### 1. Register User (with Auth Service)
```bash
curl -X POST http://localhost:8084/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "display_name": "John Doe",
    "phone_number": "+1234567890"
  }'
```

### 2. Login User (get token)
```bash
curl -X POST http://localhost:8084/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 3. Create Auction (with auth token)
```bash
curl -X POST http://localhost:8083/api/v1/auctions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "product_id": "product-123",
    "title": "Vintage Watch",
    "description": "A beautiful vintage watch",
    "starting_price": 100.00,
    "reserve_price": 500.00,
    "min_bid_increment": 10.00,
    "start_time": "2025-10-20T10:00:00Z",
    "end_time": "2025-10-20T11:00:00Z",
    "type": "scheduled",
    "images": ["https://example.com/watch.jpg"]
  }'
```

### 4. Place Bid (with auth token)
```bash
curl -X POST http://localhost:8083/api/v1/auctions/AUCTION_ID/bids \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 150.00
  }'
```

### 5. List Auctions (no auth required)
```bash
curl -X GET "http://localhost:8083/api/v1/auctions?status=active&page=1&limit=10"
```

## Security Features

- **JWT Token Validation**: All protected endpoints validate tokens with auth service
- **User Context**: Authenticated user ID available in request context
- **Error Handling**: Proper authentication error responses
- **Async Firebase**: Firebase persistence happens asynchronously for performance
- **Bid Authorization**: Only authenticated users can place bids
- **Auction Ownership**: Users can only modify their own auctions

## Error Responses

### Authentication Errors
```json
{
  "success": false,
  "error": "NO_TOKEN",
  "message": "No authorization token provided"
}
```

```json
{
  "success": false,
  "error": "INVALID_TOKEN",
  "message": "Invalid or expired token"
}
```

### Success Response
```json
{
  "success": true,
  "message": "Auction created successfully",
  "data": {
    "auction": {
      "auction_id": "auction-123",
      "seller_id": "user-456",
      "title": "Vintage Watch",
      "current_price": 100.00,
      "status": "scheduled"
    }
  }
}
```

## Testing

Run the authentication integration test:
```bash
./test-auction-auth.sh
```

## Integration with Frontend

React Native app should:
1. Store JWT token after login
2. Include token in API requests to protected endpoints
3. Handle authentication errors appropriately
4. Refresh tokens when expired

## Next Steps

1. **Test Integration**: Run the test script to verify auth works
2. **Update Frontend**: Modify React Native app to use new auth flow
3. **Other Services**: Add auth to product, order, payment services
4. **Production Deploy**: Deploy complete system with auth

## Troubleshooting

### Token Validation Fails
- Check auth service is running on port 8084
- Verify JWT token is valid and not expired
- Ensure correct Authorization header format

### User Context Missing
- Verify auth middleware is applied to route
- Check that auth client is properly initialized
- Ensure token is being passed correctly

### Performance Issues
- Auth validation adds ~5-10ms to requests
- Consider caching user context for repeated requests
- Monitor auth service health and performance

The auction service authentication integration is complete and ready for testing! 🚀
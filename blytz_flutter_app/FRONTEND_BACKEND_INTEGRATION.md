# Flutter Frontend - Backend Integration Guide

This document describes the complete integration between the Flutter frontend and the Go microservices backend for the Blytz Live Auction MVP.

## Architecture Overview

### Microservices Structure
- **Auth Service** (Port 8084): User authentication with JWT tokens
- **Auction Service** (Port 8081): Core auction functionality and bidding
- **Product Service** (Port 8082): Product management and catalog
- **Order Service** (Port 8083): Order processing and management
- **Payment Service** (Port 8086): Payment processing and transactions
- **Chat Service** (Port 8088): Real-time messaging
- **Logistics Service** (Port 8087): Shipping and delivery management
- **API Gateway** (Port 8080): Central routing and load balancing

### Flutter App Structure
```
lib/
├── core/
│   ├── constants/api_constants.dart     # API endpoints and URLs
│   ├── network/api_client.dart          # Retrofit API clients
│   ├── network/interceptors.dart        # Auth and logging interceptors
│   ├── services/websocket_service.dart  # Real-time WebSocket handling
│   ├── di/service_locator.dart          # Dependency injection
│   └── storage/                        # Secure storage and preferences
├── features/
│   ├── auth/
│   │   ├── data/models/auth_model.dart
│   │   └── data/repositories/auth_repository_impl.dart
│   └── auction/
│       ├── data/models/auction_model.dart
│       └── data/repositories/auction_repository_impl.dart
└── main.dart                           # App initialization
```

## API Integration

### Authentication Flow
1. **Login**: POST `/api/v1/auth/login` → JWT token
2. **Register**: POST `/api/v1/auth/register` → JWT token
3. **Token Refresh**: POST `/api/v1/auth/refresh` → New JWT token
4. **Auto-refresh**: Automatic token refresh on 401 responses

### Key Features Implemented

#### ✅ **Authentication Service Integration**
- JWT token management with automatic refresh
- Secure token storage using flutter_secure_storage
- Auth interceptor for automatic token injection
- Complete auth repository with error handling

#### ✅ **Auction Service Integration**
- Full CRUD operations for auctions
- Real-time bid placement
- Auction status management
- Bid history retrieval

#### ✅ **Product Service Integration**
- Product catalog with filtering
- Featured products
- User-specific product management
- Category and search support

#### ✅ **Order Service Integration**
- Order creation and management
- Status tracking
- Order history with pagination

#### ✅ **Real-time WebSocket Connections**
- Auction bid updates
- Live streaming WebSocket
- Chat functionality
- Automatic reconnection handling

#### ✅ **State Management**
- Riverpod providers for all services
- Dependency injection setup
- Error handling with Either pattern
- Network connectivity monitoring

## Configuration

### Environment Setup
Update `lib/core/constants/api_constants.dart` with your backend URLs:

```dart
class ApiConstants {
  // Update these URLs to match your deployment
  static const String authBaseUrl = 'http://localhost:8084';
  static const String auctionBaseUrl = 'http://localhost:8081';
  static const String productBaseUrl = 'http://localhost:8082';
  // ... other services
}
```

### Development Configuration
For local development:
1. Start backend services: `cd infra && docker-compose up -d`
2. Run Flutter app: `flutter run`
3. Ensure all ports are accessible (8081-8088)

## API Endpoints

### Authentication (Port 8084)
```dart
// Login
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

// Register
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe"
}
```

### Auctions (Port 8081)
```dart
// Get active auctions
GET /api/v1/auctions/active?page=1&limit=20

// Place bid
POST /api/v1/auctions/{id}/bids
{
  "amount": 150.00
}

// Get auction details
GET /api/v1/auctions/{id}
```

### Products (Port 8082)
```dart
// Get products with filtering
GET /api/v1/products?page=1&limit=20&category=electronics

// Get featured products
GET /api/v1/products/featured?limit=10
```

### Orders (Port 8083)
```dart
// Create order
POST /api/v1/orders
{
  "auctionId": "auction123",
  "shippingAddress": { ... }
}

// Get user orders
GET /api/v1/orders?page=1&limit=20
```

## WebSocket Events

### Auction Updates
```dart
// Connect to auction updates
ws://localhost:8081/ws/auctions/{auctionId}

// Event types
{
  "type": "bid_placed",
  "data": {
    "auctionId": "123",
    "bid": { "amount": 150.00, "bidderName": "John" }
  }
}

{
  "type": "auction_ended",
  "data": { "winnerId": "user123", "finalPrice": 200.00 }
}
```

### Live Streaming
```dart
// Connect to live stream
wss://api.blytz.app/ws/live/{streamId}

// Event types
{
  "type": "viewer_count",
  "data": { "count": 150 }
}
```

## Error Handling

### HTTP Error Handling
- **401 Unauthorized**: Automatic token refresh attempt
- **403 Forbidden**: Permission denied
- **404 Not Found**: Resource not found
- **500 Server Error**: Server-side error

### Error Types
```dart
// Authentication errors
class AuthFailure extends Failure { ... }

// Network errors
class NetworkFailure extends Failure { ... }

// Server errors
class ServerFailure extends Failure { ... }

// Validation errors
class ValidationFailure extends Failure { ... }
```

## Testing

### Integration Tests
Run the complete integration test suite:

```bash
flutter test test/integration/backend_integration_test.dart
```

### Test Coverage
- API client creation and configuration
- Authentication flow testing
- Auction model parsing
- WebSocket service initialization
- Dependency injection setup

## Performance Optimizations

### Network Optimization
- Connection pooling with Dio
- Request/response interceptors
- Automatic retry on failures
- Request caching for static data

### State Management
- Riverpod providers for efficient state updates
- Stream-based real-time updates
- Automatic disposal of resources

## Security

### Token Security
- Secure storage using flutter_secure_storage
- Automatic token refresh
- Token validation on app start
- Secure token transmission

### Network Security
- HTTPS in production
- Certificate pinning (optional)
- Request/response encryption

## Deployment

### Production Configuration
1. Update API URLs to production endpoints
2. Configure SSL/TLS certificates
3. Set up proper CORS configuration
4. Enable request signing (optional)

### Environment Variables
```dart
// Use environment-specific configuration
class ApiConstants {
  static const String baseUrl =
    String.fromEnvironment('API_BASE_URL',
    defaultValue: 'http://localhost:8080');
}
```

## Troubleshooting

### Common Issues

#### Connection Refused
- Ensure backend services are running
- Check if ports 8081-8088 are accessible
- Verify firewall settings

#### Authentication Failures
- Check JWT token format
- Verify auth service is running on port 8084
- Ensure proper token storage

#### WebSocket Connection Issues
- Check WebSocket server status
- Verify WebSocket URL format
- Ensure proper event handling

#### Build Issues
- Run `flutter pub get` to install dependencies
- Clean build: `flutter clean && flutter pub get`
- Check for conflicting dependencies

### Debug Tools
- Enable logging interceptor for API debugging
- Use Flutter Inspector for state debugging
- Monitor network requests with Charles Proxy or similar

## Next Steps

### Features to Implement
1. **Complete Payment Integration**: Connect to payment service
2. **Enhanced Chat Features**: Message history, typing indicators
3. **Push Notifications**: Firebase Cloud Messaging integration
4. **Offline Support**: Local caching and synchronization
5. **Advanced Search**: Elasticsearch integration
6. **File Upload**: Product images and documents

### Performance Improvements
1. **Image Optimization**: Caching and compression
2. **Pagination**: Infinite scroll for large datasets
3. **Background Sync**: Data synchronization in background
4. **Analytics**: User behavior tracking

This integration provides a solid foundation for the Blytz Live Auction platform with complete frontend-backend connectivity, real-time features, and robust error handling.
# Blytz Live Auction Platform

A comprehensive microservices-based live auction platform built with Go, featuring real-time bidding, chat, and payment processing.

## Architecture

The platform consists of 8 microservices:

- **Auth Service** (Port 8084) - Self-hosted authentication with Better Auth and JWT tokens
- **Product Service** (Port 8082) - Product catalog management
- **Auction Service** (Port 8083) - Auction management and bidding with Redis backend
- **Chat Service** (Port 8088) - Real-time chat functionality
- **Order Service** (Port 8085) - Order processing and management
- **Payment Service** (Port 8086) - Payment processing
- **Logistics Service** (Port 8087) - Shipping and delivery management
- **API Gateway** (Port 8080) - Central API gateway and routing

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- PostgreSQL 15+
- Redis 7+

### Deployment

1. **Start the entire stack:**
   ```bash
   docker-compose up -d
   ```

2. **Check service health:**
   ```bash
   # Check all services
   docker-compose ps

   # Check individual service health
   curl http://localhost:8081/health  # Auth service
   curl http://localhost:8082/health  # Product service
   curl http://localhost:8083/health  # Auction service
   # ... and so on for all services

   # Check gateway
   curl http://localhost:8080/health
   ```

3. **Access the services:**
   - API Gateway: http://localhost:8080
   - Auth Service: http://localhost:8081
   - Product Service: http://localhost:8082
   - Auction Service: http://localhost:8083
   - Chat Service: http://localhost:8084
   - Order Service: http://localhost:8085
   - Payment Service: http://localhost:8086
   - Logistics Service: http://localhost:8087

4. **Monitoring:**
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000 (admin/admin)

### API Documentation

The API documentation is available through the OpenAPI specification in the `/openapi` directory. You can view it using Swagger UI or any OpenAPI viewer.

### Database

- **PostgreSQL**: Main database running on port 5432
- **Redis**: Cache and session storage on port 6379

Each service has its own database:
- `auth_db` - Authentication data
- `products_db` - Product catalog
- `auction_db` - Auction and bidding data
- `chat_db` - Chat messages
- `orders_db` - Order information
- `payments_db` - Payment transactions
- `logistics_db` - Shipping and delivery data

### Development

For local development without Docker:

1. **Install dependencies:**
   ```bash
   # Install Go dependencies for each service
   cd backend/auth-service && go mod download
   cd backend/product-service && go mod download
   # ... repeat for all services
   ```

2. **Set up local databases:**
   ```bash
   # Start PostgreSQL and Redis (using Docker for convenience)
   docker-compose up -d postgres redis
   ```

3. **Run individual services:**
   ```bash
   cd backend/auth-service
   go run cmd/main.go
   ```

### Environment Variables

Each service supports the following environment variables:

- `PORT` - Service port (default: service-specific)
- `ENVIRONMENT` - Environment (development/production)
- `LOG_LEVEL` - Logging level (debug/info/warn/error)
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - JWT signing secret (auth service)

### Monitoring

The platform includes comprehensive monitoring:

- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization and dashboards
- **Health checks**: Built-in health endpoints for all services

### Scaling

To scale individual services:

```bash
# Scale auction service to 3 instances
docker-compose up -d --scale auction-service=3
```

### Troubleshooting

1. **Service won't start:**
   - Check Docker logs: `docker-compose logs <service-name>`
   - Verify database connectivity
   - Check environment variables

2. **Database connection issues:**
   - Ensure PostgreSQL is running and accessible
   - Check database initialization script
   - Verify connection strings

3. **Build failures:**
   - Ensure Go 1.21+ is installed
   - Check for missing dependencies: `go mod tidy`
   - Verify shared module is accessible
   - **CRITICAL**: Ensure all shared imports use `/pkg/` prefix (e.g., `github.com/blytz/shared/pkg/auth`)

4. **Authentication issues:**
   - Verify auth service is running on port 8084
   - Check JWT_SECRET and BETTER_AUTH_SECRET environment variables
   - Test auth endpoints manually before service integration
   - Use provided test scripts for validation

## Authentication System (NEW - Phase 2 Complete)

### Overview
The platform now features a **self-hosted authentication system** using Better Auth, providing significant cost savings and performance benefits over Firebase Auth.

### Key Features
- **JWT-based authentication** with HS256 algorithm
- **Complete user management**: registration, login, token refresh, profile updates
- **Microservices integration** with shared auth client
- **Cost-effective**: 97% savings vs Firebase Auth ($480-585/month)
- **Performance**: Local queries (~5ms) vs external API calls (~100ms)
- **Security**: Complete data ownership, no vendor lock-in

### Authentication Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │  API Gateway    │    │  Auth Service   │
│                 │────│                 │────│                 │
│ React Native    │    │   Nginx/Go      │    │   Go + Better   │
│                 │    │                 │    │   Auth + JWT    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │   PostgreSQL    │
                                               │  (Auth Database)│
                                               └─────────────────┘
```

### Testing Authentication
```bash
# Test complete auth flow
cd services/auth-service && ./test-auth-service.sh

# Test auction service with auth
cd services/auction-service && ./test-auction-auth.sh

# Verify shared package structure
./verify-shared-migration.sh

# Manual testing
curl -X POST http://localhost:8084/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "display_name": "Test User"}'
```

### Protected Endpoints
The following endpoints require authentication:
- `POST /api/v1/auctions` - Create auction
- `PUT /api/v1/auctions/:id` - Update auction
- `DELETE /api/v1/auctions/:id` - Delete auction
- `POST /api/v1/auctions/:id/bids` - Place bid

### Integration Pattern
For any new service, use the shared auth client:
```go
import "github.com/blytz/shared/pkg/auth"

authClient := auth.NewAuthClient("http://auth-service:8084")
router.Use(auth.GinAuthMiddleware(authClient))
```

### Production Deployment

For production deployment:

1. Update environment variables for production
2. Configure SSL/TLS certificates
3. Set up proper database backups
4. Configure monitoring and alerting
5. Set up log aggregation
6. Configure rate limiting and security policies

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh
- `GET /api/v1/auth/verify` - Token verification
- `POST /api/v1/auth/logout` - User logout

### Products
- `GET /api/v1/products` - List products
- `POST /api/v1/products` - Create product
- `GET /api/v1/products/:id` - Get product details
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

### Auctions
- `GET /api/v1/auctions` - List auctions
- `POST /api/v1/auctions` - Create auction
- `GET /api/v1/auctions/:id` - Get auction details
- `PUT /api/v1/auctions/:id` - Update auction
- `DELETE /api/v1/auctions/:id` - Delete auction
- `POST /api/v1/auctions/:id/bids` - Place bid
- `GET /api/v1/auctions/:id/bids` - Get bids

### Chat
- `GET /api/v1/chat/rooms` - List chat rooms
- `POST /api/v1/chat/rooms` - Create chat room
- `GET /api/v1/chat/rooms/:id/messages` - Get messages
- `POST /api/v1/chat/rooms/:id/messages` - Send message

### Orders
- `GET /api/v1/orders` - List orders
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id` - Update order

### Payments
- `POST /api/v1/payments` - Process payment
- `GET /api/v1/payments/:id` - Get payment details
- `POST /api/v1/payments/:id/refund` - Process refund

### Logistics
- `GET /api/v1/logistics/shipments` - List shipments
- `POST /api/v1/logistics/shipments` - Create shipment
- `GET /api/v1/logistics/shipments/:id` - Get shipment details
- `PUT /api/v1/logistics/shipments/:id` - Update shipment

## Contributing

Please see CONTRIBUTING.md for development guidelines and contribution process.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
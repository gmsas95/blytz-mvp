# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Blytz Live Auction MVP - A real-time livestream commerce platform built with Go, Redis, and Nginx. This is a Docker-based microservices deployment setup designed for Hostinger KVM (1vCore/4GB) environments.

## Architecture

### Core Services
- **Auction Engine**: Go service with Redis for core bidding logic, anti-snipe features, and atomic Lua scripts
- **Authentication Service**: Self-hosted Better Auth microservice with JWT tokens (NEW - Phase 2 Complete)
- **Redis**: In-memory cache for real-time bid state, session cache, and product data
- **Nginx API Gateway**: Reverse proxy and load balancer for HTTP routing
- **PostgreSQL**: Primary database for user data and business logic

### External Integrations
- **Frontend**: React Native with Expo (external to this repo)
- **Firebase**: Cloud functions for payments, notifications, and data persistence
- **LiveKit**: Real-time streaming service

### Authentication Architecture (Phase 2 Complete)
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

**Key Benefits:**
- ✅ **97% cost savings** vs Firebase Auth ($480-585/month savings)
- ✅ **Performance**: Local queries (~5ms) vs external API calls (~100ms)
- ✅ **Control**: Complete data ownership, no vendor lock-in
- ✅ **Scalability**: Fixed costs regardless of user count

## Common Commands

### Docker Operations
```bash
# Build and start all services
sudo docker-compose up -d --build

# Stop all services
sudo docker-compose down

# View logs
sudo docker-compose logs -f

# Restart specific service
sudo docker restart blytz-auction

# Access Redis CLI
sudo docker exec -it blytz-redis redis-cli

# Load test bidding endpoints
k6 run infra/k6-bid-test.js
```

### Authentication Operations (NEW - Phase 2 Complete)
```bash
# Test auth service
cd services/auth-service && ./test-auth-service.sh

# Test auction service with auth
cd services/auction-service && ./test-auction-auth.sh

# Verify shared package migration
./verify-shared-migration.sh

# Manual auth testing
# Register user
curl -X POST http://localhost:8084/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "display_name": "Test User"}'

# Login and get token
curl -X POST http://localhost:8084/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'

# Test protected auction endpoint
curl -X POST http://localhost:8083/api/v1/auctions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"product_id": "test", "title": "Test Auction", "starting_price": 10.00}'
```

### Health Check
```bash
curl http://localhost:8080/health
```

### Development Commands
```bash
# Run Go tests
cd services/auction-service && go test ./...

# Build auction service locally
cd services/auction-service && go build -o auction ./cmd/main.go

# Deploy Firebase functions
cd functions && firebase deploy --only functions

# Check Docker stats
docker stats
```

## Service Configuration

### Docker Compose Services
- **auth-service**: Port 8084, authentication microservice with Better Auth
- **auction-service**: Port 8083, auction engine with Redis backend
- **product-service**: Port 8082, product management service
- **order-service**: Port 8085, order processing service
- **payment-service**: Port 8086, payment processing service
- **chat-service**: Port 8088, real-time messaging service
- **logistics-service**: Port 8087, shipping and logistics service
- **redis**: Port 6379, in-memory cache for auctions and sessions
- **postgres**: Port 5432, primary database for user data and business logic
- **nginx**: Ports 80/443, API gateway and load balancer

### API Routing
All auction API requests go through:
```
http://api.blytz.app/auction/ -> auction:8080
```

## Directory Structure
```
/blytz/
├── services/
│   ├── auth-service/           # Authentication microservice (Better Auth + JWT)
│   ├── auction-service/        # Auction engine with Redis backend
│   ├── product-service/        # Product management service
│   ├── order-service/          # Order processing service
│   ├── payment-service/        # Payment processing service
│   ├── chat-service/           # Real-time messaging service
│   └── logistics-service/      # Shipping and logistics service
├── shared/                     # Centralized shared packages (NEW)
│   └── pkg/
│       ├── auth/               # Shared auth client for microservices
│       ├── errors/             # Common error handling
│       ├── utils/              # Utility functions
│       ├── constants/          # Shared constants
│       └── proto/              # Protocol buffer definitions
├── functions/                  # Firebase Cloud Functions
├── frontend/                   # React Native Expo app
├── infra/
│   ├── docker-compose.yml      # Service orchestration
│   ├── nginx.conf              # API gateway configuration
│   ├── prometheus.yml          # Monitoring config
│   └── k6-bid-test.js          # Load test script
├── specs/                      # OpenAPI specifications
└── .github/workflows/          # CI/CD pipelines
```

## Key Technical Details

- **Deployment Target**: Hostinger KVM Ubuntu 22.04
- **Redis Scripts**: Uses atomic Lua scripts for bid operations
- **Anti-Snipe**: Implemented in auction engine logic
- **Real-time**: Redis pub/sub for bid updates
- **State Management**: All auction state in Redis with persistence
- **Load Testing**: k6 scripts for performance validation
- **Monitoring**: Prometheus metrics endpoint at `/metrics`

## MVP Development Plan (Priority Order)

### ✅ PHASE 1: Core Infrastructure (COMPLETED)
1. **Day 1-2**: ✅ Complete auction service with Go + Gin + Redis Lua scripts
2. **Day 3**: ✅ Set up VPS with Docker, deploy Redis + auction service
3. **Day 4**: ✅ Configure Nginx gateway with SSL and API routing
4. **Day 5**: ✅ Implement Firebase Cloud Functions for persistence
5. **Day 6**: ✅ Add monitoring with Prometheus metrics
6. **Day 7**: ✅ Load test with k6 (target: 50 VUs, <300ms bid latency)

### ✅ PHASE 2: Authentication System (COMPLETED - MAJOR MILESTONE)
1. **Day 8-9**: ✅ Complete self-hosted authentication system with Better Auth
2. **Day 10**: ✅ Create shared auth client for microservices
3. **Day 11**: ✅ Integrate auth middleware with auction service
4. **Day 12**: ✅ Implement comprehensive test suite
5. **Day 13**: ✅ Update Docker configuration and deployment
6. **Day 14**: ✅ Complete shared package migration and fix structural issues

**🎉 ACHIEVEMENTS:**
- **97% cost savings** vs Firebase Auth ($480-585/month savings)
- **Performance**: Local queries (~5ms) vs external API calls (~100ms)
- **Security**: JWT-based authentication with complete user control
- **Architecture**: Consistent microservices authentication pattern

### 🔄 PHASE 3: Complete Service Integration (IN PROGRESS)
1. **Day 15-16**: Integrate auth with product, order, payment services
2. **Day 17**: Integrate auth with chat and logistics services
3. **Day 18**: Implement service-to-service authentication
4. **Day 19**: Update React Native frontend auth flow
5. **Day 20**: Complete end-to-end authentication testing
6. **Day 21**: Production deployment with complete auth system

## Production Readiness Checklist

### ✅ PHASE 1 & 2 COMPLETED
- [x] **Core Auction Service**: Go + Redis with atomic Lua scripts
- [x] **Docker Infrastructure**: Complete microservices orchestration
- [x] **Nginx Gateway**: SSL configuration and API routing
- [x] **Firebase Integration**: Cloud functions for persistence
- [x] **Monitoring**: Prometheus metrics and health checks
- [x] **Load Testing**: k6 scripts with <300ms bid latency target
- [x] **Authentication System**: Self-hosted Better Auth with JWT
- [x] **Shared Packages**: Centralized auth client for microservices
- [x] **Auction Auth Integration**: Protected endpoints with user context
- [x] **Comprehensive Testing**: Unit, integration, and benchmark tests

### 🔄 PHASE 3: IN PROGRESS
- [ ] **Complete Service Auth**: All microservices with authentication
- [ ] **Service-to-Service Auth**: Internal communication security
- [ ] **Frontend Auth Integration**: React Native authentication flow
- [ ] **SSL Certificates**: Let's Encrypt configuration
- [ ] **Redis Persistence**: Backup and recovery strategy
- [ ] **Rate Limiting**: Bid endpoint protection
- [ ] **Secrets Management**: Environment variables and secure storage
- [ ] **End-to-End Testing**: Complete auth flow validation
- [ ] **Production Deployment**: VPS deployment with Dokploy
- [ ] **Log Aggregation**: Centralized logging and alerting

## Development Notes

### Shared Package Architecture (CRITICAL - Phase 2 Update)
**All shared packages now live in `/shared/pkg/` - this is MANDATORY for consistency.**

```go
// ✅ CORRECT - Use pkg/ prefix
import "github.com/blytz/shared/pkg/auth"
import "github.com/blytz/shared/pkg/errors"
import "github.com/blytz/shared/pkg/utils"

// ❌ INCORRECT - Will cause build failures
import "github.com/blytz/shared/auth"
import "github.com/blytz/shared/errors"
```

**Package Structure:**
- `shared/pkg/auth/` - Authentication client and middleware
- `shared/pkg/errors/` - Common error types and handling
- `shared/pkg/utils/` - Utility functions (logger, validation, responses)
- `shared/pkg/constants/` - Shared constants and configurations
- `shared/pkg/proto/` - Protocol buffer definitions

### Authentication Integration Pattern
**For any new microservice, follow this exact pattern:**

```go
import "github.com/blytz/shared/pkg/auth"

// In router setup
authClient := auth.NewAuthClient("http://auth-service:8084")

// Public routes (no auth required)
public := router.Group("/api/v1/public")
{
    public.GET("/health", healthHandler)
}

// Protected routes (auth required)
protected := router.Group("/api/v1/protected")
protected.Use(auth.GinAuthMiddleware(authClient))
{
    protected.POST("/create", createHandler)  // Requires auth
}

// In handlers
userID := c.GetString("userID")  // Get authenticated user
```

### Service Communication
- **Auth Service**: Port 8084 - JWT token validation
- **Service-to-Service**: Use shared auth client for internal calls
- **External APIs**: All protected endpoints require Bearer token

### This is a deployment-focused repository with Docker configurations
- Source code lives in `services/` directories
- Use environment variables for all configuration
- Follow Go module structure for all services
- Implement atomic Redis operations for bid processing
- Use WebSocket for real-time bid updates to frontend
- **CRITICAL**: All shared imports must use `/pkg/` prefix
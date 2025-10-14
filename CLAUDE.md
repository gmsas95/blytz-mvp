# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Blytz Live Auction MVP - A real-time livestream commerce platform built with Go, Redis, and Nginx. This is a Docker-based microservices deployment setup designed for Hostinger KVM (1vCore/4GB) environments.

## Architecture

- **Auction Engine**: Go service with Redis for core bidding logic, anti-snipe features, and atomic Lua scripts
- **Redis**: In-memory cache for real-time bid state, session cache, and product data
- **Nginx API Gateway**: Reverse proxy and load balancer for HTTP routing
- **Frontend**: React Native with Expo (external to this repo)
- **Firebase**: External cloud functions for payments, authentication, and notifications
- **LiveKit**: Real-time streaming service

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
- **redis**: Port 6379, persistent volume for data
- **auction**: Port 8080, depends on Redis
- **nginx**: Ports 80/443, routes `/auction/` to auction service

### API Routing
All auction API requests go through:
```
http://api.blytz.app/auction/ -> auction:8080
```

## Directory Structure
```
/blytz/
├── services/
│   ├── auction-service/        # Go backend service
│   ├── payment-service/        # Payment processing
│   └── order-service/          # Order management
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

### Week 1: Core Auction Service
1. **Day 1-2**: Complete auction service with Go + Gin + Redis Lua scripts
2. **Day 3**: Set up VPS with Docker, deploy Redis + auction service
3. **Day 4**: Configure Nginx gateway with SSL and API routing
4. **Day 5**: Implement Firebase Cloud Functions for persistence
5. **Day 6**: Add monitoring with Prometheus metrics
6. **Day 7**: Load test with k6 (target: 50 VUs, <300ms bid latency)

### Week 2: Integration & Frontend
1. **Day 1-2**: Integrate LiveKit streaming tokens
2. **Day 3-4**: Connect React Native frontend to auction API
3. **Day 5**: End-to-end testing: create auction → bid → end → persist
4. **Day 6**: Set up CI/CD pipeline with GitHub Actions
5. **Day 7**: Production readiness checklist completion

## Production Readiness Checklist

- [ ] SSL certificates via Let's Encrypt
- [ ] Firebase authentication token validation
- [ ] Redis persistence and backup strategy
- [ ] Rate limiting for bid endpoints
- [ ] Health checks and monitoring
- [ ] Secrets management (environment variables)
- [ ] Unit and integration tests
- [ ] Load testing validation
- [ ] Rollback deployment strategy
- [ ] Log aggregation and alerting

## Development Notes

- This is a deployment-focused repository with Docker configurations
- Source code lives in `services/` directories
- Use environment variables for all configuration
- Follow Go module structure for auction service
- Implement atomic Redis operations for bid processing
- Use WebSocket for real-time bid updates to frontend
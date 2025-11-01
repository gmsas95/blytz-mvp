#!/bin/bash
# Deploy Missing Backend Services to Fix 502 Errors

echo "🔧 Deploying missing backend services to fix 502 errors..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 Step 1: Checking current service status...${NC}"
cd /home/sas/blytzmvp-clean
docker compose ps

echo -e "${YELLOW}🛑 Step 2: Stopping any existing backend services...${NC}"
docker compose stop auth-service auction-service product-service payment-service order-service logistics-service chat-service gateway 2>/dev/null || echo "No backend services to stop"

echo -e "${YELLOW}🏗️ Step 3: Starting database and cache services...${NC}"
docker compose up -d postgres redis

echo -e "${YELLOW}⏳ Step 4: Waiting for database to be ready...${NC}"
sleep 30

echo -e "${YELLOW}🔐 Step 5: Starting authentication service...${NC}"
docker compose up -d auth-service

echo -e "${YELLOW}⏳ Step 6: Waiting for auth service...${NC}"
sleep 20

echo -e "${YELLOW}🚀 Step 7: Starting all backend microservices...${NC}"
docker compose up -d auction-service product-service payment-service order-service logistics-service chat-service gateway

echo -e "${YELLOW}⏳ Step 8: Waiting for services to initialize...${NC}"
sleep 45

echo -e "${YELLOW}🧪 Step 9: Testing service connectivity...${NC}"

# Test each service
test_service() {
    local service_name=$1
    local port=$2
    local path=${3:-/health}
    
    echo -n "Testing $service_name on port $port... "
    if curl -f -s --max-time 10 "http://localhost:$port$path" > /dev/null; then
        echo -e "${GREEN}✅ OK${NC}"
        return 0
    else
        echo -e "${RED}❌ FAILED${NC}"
        return 1
    fi
}

# Test all services
test_service "Gateway" 8080
test_service "Auth Service" 8084
test_service "Auction Service" 8083
test_service "Product Service" 8082
test_service "Payment Service" 8086
test_service "Order Service" 8085
test_service "Logistics Service" 8087
test_service "Chat Service" 8088

echo -e "${YELLOW}📊 Step 10: Showing final service status...${NC}"
docker compose ps

echo -e "${YELLOW}🌐 Step 11: Testing external API access...${NC}"
echo "Testing API gateway from external access..."
curl -I https://api.blytz.app/health || echo -e "${RED}❌ External API test failed${NC}"

echo -e "${GREEN}✅ Backend services deployment complete!${NC}"
echo ""
echo -e "${BLUE}🔍 Troubleshooting:${NC}"
echo "- If services are failing, check logs: docker compose logs [service-name]"
echo "- If 502 errors persist, check Traefik logs: sudo docker logs dokploy-traefik"
echo "- Verify environment variables in .env file"
echo ""
echo -e "${BLUE}📊 Service URLs:${NC}"
echo "- API Gateway: https://api.blytz.app"
echo "- Auth Service: https://api.blytz.app/auth/health"
echo "- Auction Service: https://api.blytz.app/auctions/health"
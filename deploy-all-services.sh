#!/bin/bash
# Deploy All Blytz Services to VPS

echo "🚀 Deploying complete Blytz platform to VPS..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENV_FILE="/home/sas/blytzmvp-clean/.env"
COMPOSE_FILE="/home/sas/blytzmvp-clean/docker-compose.yml"

echo -e "${BLUE}📋 Step 1: Checking prerequisites...${NC}"

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}❌ .env file not found at $ENV_FILE${NC}"
    echo "Please create .env file with production credentials"
    exit 1
fi

echo -e "${GREEN}✅ .env file found${NC}"

# Check if docker-compose.yml exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo -e "${RED}❌ docker-compose.yml not found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ docker-compose.yml found${NC}"

echo -e "${YELLOW}🛑 Step 2: Stopping existing services...${NC}"
cd /home/sas/blytzmvp-clean
docker compose down 2>/dev/null || echo "No existing services to stop"

echo -e "${YELLOW}🏗️ Step 3: Building and starting all services...${NC}"

# Start core services first (database, redis)
echo "Starting core services (PostgreSQL, Redis)..."
docker compose up -d postgres redis

# Wait for database to be ready
echo "Waiting for database to be ready..."
sleep 30

# Start backend services
echo "Starting backend services..."
docker compose up -d auth-service

# Wait for auth service
echo "Waiting for auth service..."
sleep 20

# Start remaining backend services
docker compose up -d auction-service product-service payment-service order-service logistics-service chat-service gateway

# Wait for all backend services
echo "Waiting for backend services to be ready..."
sleep 30

# Start frontend services
echo "Starting frontend services..."
docker compose up -d frontend

echo -e "${YELLOW}⏳ Step 4: Waiting for all services to be healthy...${NC}"
sleep 60

echo -e "${YELLOW}🧪 Step 5: Checking service health...${NC}"

# Check each service
services=("postgres:5432" "redis:6379" "auth-service:8084" "auction-service:8083" "product-service:8082" "payment-service:8086" "order-service:8085" "logistics-service:8087" "chat-service:8088" "gateway:8080" "frontend:3000")

for service in "${services[@]}"; do
    service_name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    
    echo "Checking $service_name on port $port..."
    if docker compose ps | grep -q "$service_name.*Up"; then
        echo -e "${GREEN}✅ $service_name is running${NC}"
    else
        echo -e "${RED}❌ $service_name is not running${NC}"
    fi
done

echo -e "${YELLOW}📊 Step 6: Testing API endpoints...${NC}"

# Test API gateway
echo "Testing API gateway..."
curl -f http://localhost:8080/health && echo -e "${GREEN} ✅ Gateway healthy${NC}" || echo -e "${RED} ❌ Gateway unhealthy${NC}"

# Test auth service
echo "Testing auth service..."
curl -f http://localhost:8084/health && echo -e "${GREEN} ✅ Auth service healthy${NC}" || echo -e "${RED} ❌ Auth service unhealthy${NC}"

echo -e "${YELLOW}🌐 Step 7: Deploying Traefik configuration...${NC}"

# Deploy Traefik config
./deploy-to-vps.sh

echo -e "${YELLOW}📋 Step 8: Showing service status...${NC}"
docker compose ps

echo -e "${GREEN}✅ Complete deployment finished!${NC}"
echo ""
echo -e "${BLUE}🎯 Next steps:${NC}"
echo "1. Check Traefik dashboard at https://sudo.blytz.app"
echo "2. Test main site at https://blytz.app"
echo "3. Test API at https://api.blytz.app/health"
echo "4. Deploy demo and seller frontends through Dokploy"
echo ""
echo -e "${BLUE}📚 Useful commands:${NC}"
echo "- View logs: docker compose logs -f [service-name]"
echo "- Restart service: docker compose restart [service-name]"
echo "- Check health: docker compose ps"
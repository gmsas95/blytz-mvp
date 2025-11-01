#!/bin/bash
# Fix 502 Errors by Deploying Missing Backend Services
# This script assumes environment variables are stored in Dokploy secrets manager

echo "🔧 Fixing 502 errors by deploying missing backend services..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 Step 1: Checking current service status...${NC}"
cd /home/sas/blytzmvp-clean

# Check what's currently running
echo "Current Docker containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo -e "${YELLOW}🏗️ Step 2: Deploying backend services via Dokploy...${NC}"

# Services that need to be deployed through Dokploy
services=(
    "blytz-postgres-prod:Database:5432"
    "blytz-redis-prod:Cache:6379" 
    "blytz-auth-prod:Authentication:8084"
    "blytz-gateway-prod:API Gateway:8080"
    "blytz-auction-prod:Auction Service:8083"
    "blytz-product-prod:Product Service:8082"
    "blytz-payment-prod:Payment Service:8086"
    "blytz-order-prod:Order Service:8085"
    "blytz-logistics-prod:Logistics Service:8087"
    "blytz-chat-prod:Chat Service:8088"
)

echo -e "${BLUE}📝 Services that need to be deployed in Dokploy:${NC}"
for service in "${services[@]}"; do
    name=$(echo $service | cut -d: -f1)
    desc=$(echo $service | cut -d: -f2)
    port=$(echo $service | cut -d: -f3)
    echo "  - $name ($desc) - Port $port"
done

echo ""
echo -e "${YELLOW}🚀 Step 3: Instructions for deploying via Dokploy:${NC}"
echo ""
echo "1. Go to https://sudo.blytz.app (Dokploy dashboard)"
echo "2. For each service below, create a new application:"
echo ""
echo -e "${BLUE}📋 Required Applications:${NC}"

for service in "${services[@]}"; do
    name=$(echo $service | cut -d: -f1)
    desc=$(echo $service | cut -d: -f2)
    port=$(echo $service | cut -d: -f3)
    
    echo ""
    echo -e "${GREEN}🏷️  Application Name: $name${NC}"
    echo -e "${GREEN}📝 Description: $desc${NC}"
    echo -e "${GREEN}🔌 Port: $port${NC}"
    
    case $name in
        *postgres*)
            echo "📦 Image: postgres:15-alpine"
            echo "🔧 Environment Variables:"
            echo "   - POSTGRES_USER=blytz"
            echo "   - POSTGRES_PASSWORD=[from secrets]"
            echo "   - POSTGRES_DB=blytz_prod"
            echo "💾 Volume: /var/lib/postgresql/data"
            ;;
        *redis*)
            echo "📦 Image: redis:7-alpine"
            echo "🔧 Command: redis-server --appendonly yes"
            echo "💾 Volume: /data"
            ;;
        *auth*)
            echo "📦 Image: gmsas95/blytz-mvp:auth-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8084"
            echo "   - NODE_ENV=production"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - JWT_SECRET=[from secrets]"
            echo "   - BETTER_AUTH_SECRET=[from secrets]"
            echo "   - BETTER_AUTH_URL=https://api.blytz.app"
            ;;
        *gateway*)
            echo "📦 Image: gmsas95/blytz-mvp:gateway-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8080"
            echo "   - NODE_ENV=production"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            ;;
        *auction*)
            echo "📦 Image: gmsas95/blytz-mvp:auction-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8083"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - REDIS_URL=redis://blytz-redis-prod:6379"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            ;;
        *product*)
            echo "📦 Image: gmsas95/blytz-mvp:product-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8082"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            ;;
        *payment*)
            echo "📦 Image: gmsas95/blytz-mvp:payment-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8086"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            echo "   - FIUU_MERCHANT_ID=[from secrets]"
            echo "   - FIUU_SECRET_KEY=[from secrets]"
            ;;
        *order*)
            echo "📦 Image: gmsas95/blytz-mvp:order-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8085"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            ;;
        *logistics*)
            echo "📦 Image: gmsas95/blytz-mvp:logistics-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8087"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            ;;
        *chat*)
            echo "📦 Image: gmsas95/blytz-mvp:chat-service-latest"
            echo "🔧 Environment Variables:"
            echo "   - PORT=8088"
            echo "   - DATABASE_URL=[from secrets]"
            echo "   - AUTH_SERVICE_URL=http://blytz-auth-prod:8084"
            ;;
    esac
    echo "🌐 Network: dokploy-network"
    echo "🔄 Restart Policy: unless-stopped"
    echo "---"
done

echo ""
echo -e "${YELLOW}🔧 Step 4: After deployment, test services:${NC}"
echo ""
echo "Test API Gateway:"
echo "curl -I https://api.blytz.app/health"
echo ""
echo "Test Auth Service:"
echo "curl -I https://api.blytz.app/auth/health"
echo ""
echo "Test Auction Service:"
echo "curl -I https://api.blytz.app/auctions/health"

echo ""
echo -e "${YELLOW}📊 Step 5: Check Traefik logs if 502 errors persist:${NC}"
echo "sudo docker logs dokploy-traefik -f"

echo ""
echo -e "${GREEN}✅ Instructions complete!${NC}"
echo ""
echo -e "${BLUE}🎯 Quick Summary:${NC}"
echo "1. Deploy all backend services via Dokploy dashboard"
echo "2. Ensure secrets are properly configured"
echo "3. Test API endpoints"
echo "4. Check Traefik logs if issues persist"
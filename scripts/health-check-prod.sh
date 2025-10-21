#!/bin/bash

# Blytz Production Health Check Script
echo "🏥 Running production health checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test endpoints
endpoints=(
    "https://blytz.app:Frontend"
    "https://api.blytz.app/health:API Gateway"
    "https://api.blytz.app/api/v1/auth/health:Auth Service"
    "https://api.blytz.app/api/v1/auctions/health:Auction Service"
    "https://api.blytz.app/api/v1/products/health:Product Service"
    "https://api.blytz.app/api/v1/orders/health:Order Service"
    "https://api.blytz.app/api/v1/payments/health:Payment Service"
    "https://api.blytz.app/api/v1/chat/health:Chat Service"
    "https://api.blytz.app/api/v1/logistics/health:Logistics Service"
)

echo "Testing production endpoints..."
echo "================================"

all_healthy=true

for endpoint in "${endpoints[@]}"; do
    IFS=':' read -r url service_name <<< "$endpoint"

    echo -n "Testing $service_name ($url)... "

    if response=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "$url" 2>/dev/null); then
        if [ "$response" = "200" ]; then
            echo -e "${GREEN}✅ HEALTHY (HTTP $response)${NC}"
        elif [ "$response" = "000" ]; then
            echo -e "${RED}❌ TIMEOUT/CONNECTION FAILED${NC}"
            all_healthy=false
        else
            echo -e "${YELLOW}⚠️  WARNING (HTTP $response)${NC}"
        fi
    else
        echo -e "${RED}❌ FAILED${NC}"
        all_healthy=false
    fi
done

echo "================================"

if [ "$all_healthy" = true ]; then
    echo -e "${GREEN}🎉 All services are healthy!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some services are unhealthy!${NC}"
    echo ""
    echo "🔧 Troubleshooting:"
    echo "1. Check Dokploy logs: dokploy logs blytz-live-auction"
    echo "2. Check service status: dokploy status blytz-live-auction"
    echo "3. Verify domain DNS: nslookup blytz.app"
    echo "4. Check SSL certificates: curl -v https://blytz.app"
    exit 1
fi
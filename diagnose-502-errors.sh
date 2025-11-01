#!/bin/bash
# Diagnose 502 Errors on VPS

echo "🔍 Diagnosing 502 errors on Blytz VPS..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 Step 1: Checking Traefik status...${NC}"
if sudo docker ps | grep -q "dokploy-traefik.*Up"; then
    echo -e "${GREEN}✅ Traefik is running${NC}"
else
    echo -e "${RED}❌ Traefik is not running${NC}"
    echo "Run: sudo docker start dokploy-traefik"
fi

echo -e "${BLUE}📋 Step 2: Checking Traefik configuration...${NC}"
if [ -f "/etc/dokploy/traefik/dynamic/dokploy.yml" ]; then
    echo -e "${GREEN}✅ Traefik config exists${NC}"
    echo "Checking for syntax errors..."
    sudo docker exec dokploy-traefik traefik config check --configFile=/etc/dokploy/traefik/traefik.yml 2>/dev/null && echo -e "${GREEN}✅ Config syntax OK${NC}" || echo -e "${RED}❌ Config syntax error${NC}"
else
    echo -e "${RED}❌ Traefik config missing${NC}"
    echo "Run: ./deploy-to-vps.sh"
fi

echo -e "${BLUE}📋 Step 3: Checking backend services...${NC}"

# Services that should be running
services=(
    "blytz-postgres-prod:5432"
    "blytz-redis-prod:6379"
    "blytz-auth-prod:8084"
    "blytz-gateway-prod:8080"
    "blytz-auction-prod:8083"
    "blytz-product-prod:8082"
    "blytz-payment-prod:8086"
    "blytz-order-prod:8085"
    "blytz-logistics-prod:8087"
    "blytz-chat-prod:8088"
)

missing_services=()
running_services=()

for service in "${services[@]}"; do
    name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    
    if sudo docker ps --format "{{.Names}}" | grep -q "^${name}$"; then
        echo -e "${GREEN}✅ $name is running${NC}"
        running_services+=("$name")
    else
        echo -e "${RED}❌ $name is NOT running${NC}"
        missing_services+=("$name")
    fi
done

echo ""
echo -e "${BLUE}📊 Summary:${NC}"
echo -e "${GREEN}Running services: ${#running_services[@]}${NC}"
echo -e "${RED}Missing services: ${#missing_services[@]}${NC}"

if [ ${#missing_services[@]} -gt 0 ]; then
    echo ""
    echo -e "${YELLOW}🚨 Missing services causing 502 errors:${NC}"
    for service in "${missing_services[@]}"; do
        echo "  - $service"
    done
    echo ""
    echo -e "${BLUE}💡 Solution:${NC}"
    echo "Deploy missing services via Dokploy dashboard at https://sudo.blytz.app"
    echo "Or run: ./fix-502-deployment-guide.sh"
fi

echo -e "${BLUE}📋 Step 4: Testing external connectivity...${NC}"

# Test main endpoints
echo "Testing main domain..."
if curl -s -I https://blytz.app | head -1 | grep -q "200\|301\|302"; then
    echo -e "${GREEN}✅ blytz.app is accessible${NC}"
else
    echo -e "${RED}❌ blytz.app is NOT accessible${NC}"
fi

echo "Testing API gateway..."
if curl -s -I https://api.blytz.app/health | head -1 | grep -q "200"; then
    echo -e "${GREEN}✅ api.blytz.app is accessible${NC}"
else
    echo -e "${RED}❌ api.blytz.app is NOT accessible (502 error)${NC}"
fi

echo -e "${BLUE}📋 Step 5: Recent Traefik logs...${NC}"
echo "Last 10 Traefik log entries:"
sudo docker logs dokploy-traefik --tail 10 2>/dev/null | grep -E "(502|error|warning)" || echo "No recent errors in logs"

echo ""
echo -e "${BLUE}📋 Step 6: Network connectivity test...${NC}"
echo "Testing if services can reach each other..."

# Test if gateway can reach auth (if both are running)
if [[ " ${running_services[@]} " =~ " blytz-gateway-prod " ]] && [[ " ${running_services[@]} " =~ " blytz-auth-prod " ]]; then
    echo "Testing gateway -> auth connectivity..."
    sudo docker exec blytz-gateway-prod curl -s --connect-timeout 5 http://blytz-auth-prod:8084/health >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Gateway can reach Auth service${NC}"
    else
        echo -e "${RED}❌ Gateway cannot reach Auth service${NC}"
    fi
fi

echo ""
echo -e "${GREEN}✅ Diagnosis complete!${NC}"
echo ""
echo -e "${BLUE}🎯 Next Steps:${NC}"
if [ ${#missing_services[@]} -gt 0 ]; then
    echo "1. Deploy missing backend services via Dokploy"
    echo "2. Ensure secrets are configured in Dokploy"
    echo "3. Test API endpoints after deployment"
else
    echo "1. Check Traefik logs for routing issues"
    echo "2. Verify service health endpoints"
    echo "3. Check network connectivity between services"
fi

echo ""
echo -e "${BLUE}🔧 Useful Commands:${NC}"
echo "- View Traefik logs: sudo docker logs dokploy-traefik -f"
echo "- Restart Traefik: sudo docker restart dokploy-traefik"
echo "- Check all containers: sudo docker ps -a"
echo "- Test API locally: curl -I http://localhost:8080/health"
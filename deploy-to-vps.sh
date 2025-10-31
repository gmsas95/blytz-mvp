#!/bin/bash
# Deploy Traefik Configuration to Hostinger VPS

echo "🚀 Deploying Traefik configuration to VPS..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}📋 Step 1: Copying configuration files to VPS...${NC}"

# Copy dokploy.yml
if [ -f "dokploy.yml" ]; then
    echo "Copying dokploy.yml..."
    sudo cp dokploy.yml /etc/dokploy/traefik/dynamic/dokploy.yml
    echo -e "${GREEN}✅ dokploy.yml copied${NC}"
else
    echo -e "${RED}❌ dokploy.yml not found${NC}"
    exit 1
fi

# Copy traefik.yml (if exists)
if [ -f "traefik.yml" ]; then
    echo "Copying traefik.yml..."
    sudo cp traefik.yml /etc/dokploy/traefik/traefik.yml 2>/dev/null || echo "No traefik.yml to update"
    echo -e "${GREEN}✅ traefik.yml copied${NC}"
fi

# Copy middlewares.yml (if exists)
if [ -f "middlewares.yml" ]; then
    echo "Copying middlewares.yml..."
    sudo cp middlewares.yml /etc/dokploy/traefik/dynamic/middlewares.yml 2>/dev/null || echo "No middlewares.yml to update"
    echo -e "${GREEN}✅ middlewares.yml copied${NC}"
fi

echo -e "${YELLOW}🔄 Step 2: Restarting Traefik...${NC}"
sudo docker restart dokploy-traefik

echo -e "${YELLOW}⏳ Step 3: Waiting for Traefik to restart...${NC}"
sleep 10

echo -e "${YELLOW}🧪 Step 4: Testing domains...${NC}"
echo "Testing main domain..."
curl -I https://blytz.app
echo ""
echo "Testing demo subdomain..."
curl -I https://demo.blytz.app
echo ""
echo "Testing seller subdomain..."
curl -I https://seller.blytz.app
echo ""
echo "Testing Dokploy dashboard..."
curl -I https://sudo.blytz.app

echo -e "${YELLOW}📊 Step 5: Checking Traefik logs...${NC}"
sudo docker logs dokploy-traefik | tail -20

echo -e "${GREEN}✅ Deployment complete!${NC}"
echo ""
echo "🎯 Next steps:"
echo "1. Create frontend-demo and frontend-seller directories with Next.js apps"
echo "2. Add LiveKit Cloud credentials to your .env file"
echo "3. Deploy the new frontend containers through Dokploy"
echo ""
echo "📚 LiveKit Cloud: https://cloud.livekit.io/"
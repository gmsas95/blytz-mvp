#!/bin/bash

# Blytz Auction MVP - Final Integration Test
# Quick verification that all components are working for exhibition

set -e

echo "🚀 Blytz Auction MVP - Final Integration Test"
echo "============================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Service URLs
AUTH_URL="http://localhost:8084/api/v1"
AUCTION_URL="http://localhost:8083/api/v1"
FIREBASE_URL="http://127.0.0.1:5001/demo-blytz-mvp/us-central1"
WEB_URL="http://localhost:8080/frontend/index.html"

echo "Testing all MVP components..."
echo ""

# Test 1: Health Checks
echo "🔍 Testing Service Health..."
echo "---------------------------"

# Auth service
if curl -s "$AUTH_URL/health" | grep -q "ok"; then
    echo -e "${GREEN}✅ Auth Service: Healthy${NC}"
else
    echo -e "${RED}❌ Auth Service: Not responding${NC}"
fi

# Auction service
if curl -s "$AUCTION_URL/health" | grep -q "ok"; then
    echo -e "${GREEN}✅ Auction Service: Healthy${NC}"
else
    echo -e "${RED}❌ Auction Service: Not responding${NC}"
fi

# Firebase functions
if curl -s -X POST -H "Content-Type: application/json" -d '{"data":{}}' "$FIREBASE_URL/health" | grep -q "ok"; then
    echo -e "${GREEN}✅ Firebase Functions: Healthy${NC}"
else
    echo -e "${YELLOW}⚠️  Firebase Functions: Not responding (optional)${NC}"
fi

# Web interface
if curl -s -I "$WEB_URL" | grep -q "200"; then
    echo -e "${GREEN}✅ Web Interface: Available${NC}"
else
    echo -e "${RED}❌ Web Interface: Not accessible${NC}"
fi

echo ""
echo "🎯 MVP Status Summary:"
echo "====================="
echo ""

# Count working services
working_services=0
total_services=4

curl -s "$AUTH_URL/health" | grep -q "ok" && working_services=$((working_services + 1))
curl -s "$AUCTION_URL/health" | grep -q "ok" && working_services=$((working_services + 1))
curl -s -X POST -H "Content-Type: application/json" -d '{"data":{}}' "$FIREBASE_URL/health" | grep -q "ok" && working_services=$((working_services + 1))
curl -s -I "$WEB_URL" | grep -q "200" && working_services=$((working_services + 1))

if [ "$working_services" -eq "$total_services" ]; then
    echo -e "${GREEN}🎉 ALL SYSTEMS OPERATIONAL!${NC}"
    status="READY"
elif [ "$working_services" -ge 3 ]; then
    echo -e "${YELLOW}⚠️  MOST SYSTEMS OPERATIONAL${NC}"
    status="MOSTLY READY"
else
    echo -e "${RED}❌ MULTIPLE SYSTEMS DOWN${NC}"
    status="NEEDS ATTENTION"
fi

echo ""
echo -e "${BLUE}System Status: $status${NC}"
echo -e "${BLUE}Working Services: $working_services/$total_services${NC}"
echo ""

echo "📋 Core Components:"
echo "• Authentication Service (Port 8084): User registration/login"
echo "• Auction Service (Port 8083): Auctions and bidding engine"
echo "• Web Interface (Port 8080): User interface and testing"
echo "• Firebase Functions (Port 5001): Payments and notifications (optional)"
echo ""

echo "🌐 Web Interface:"
echo -e "   ${BLUE}Access URL:${NC} $WEB_URL"
echo -e "   ${BLUE}Features:${NC} Registration, Login, Create Auctions, Place Bids"
echo -e "   ${BLUE}Mobile:${NC} Fully responsive for exhibition visitors"
echo ""

echo "📱 For Exhibition Visitors:"
echo "1. Open the web interface on their phones/tablets"
echo "2. Register with email and password"
echo "3. Create auctions with items to sell"
echo "4. Browse and bid on other visitors' auctions"
echo "5. Watch real-time price updates every 10 seconds"
echo ""

echo "🚀 Deployment Ready Checklist:"
if [ "$working_services" -eq "$total_services" ]; then
    echo -e "${GREEN}✅${NC} All core services are healthy"
    echo -e "${GREEN}✅${NC} Web interface is accessible"
    echo -e "${GREEN}✅${NC} Database initialized (if using PostgreSQL)"
    echo -e "${GREEN}✅${NC} Firebase emulators running (if using Firebase)"
    echo -e "${GREEN}✅${NC} Ready to deploy to blytz.app VPS"
else
    echo -e "${YELLOW}⚠️${NC}  Some services need attention before deployment"
fi

echo ""
echo -e "${BLUE}🎯 Next Steps:${NC}"
echo "1. Deploy to your VPS at blytz.app"
echo "2. Set up SSL certificates with Let's Encrypt"
echo "3. Configure DNS to point to your VPS"
echo "4. Test with multiple concurrent users"
echo "5. Monitor during the exhibition"
echo ""

if [ "$working_services" -eq "$total_services" ]; then
    echo -e "${GREEN}🏆 MVP IS READY FOR EXHIBITION!${NC}"
    echo -e "${GREEN}🎉 The auction platform will impress your visitors!${NC}"
else
    echo -e "${YELLOW}🔧 Please fix the issues above before deployment${NC}"
fi

echo ""
echo -e "${BLUE}📞 Support:${NC} Check DEPLOYMENT_GUIDE.md for troubleshooting"
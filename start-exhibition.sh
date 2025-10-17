#!/bin/bash

# Blytz Auction MVP - Exhibition Startup Script
# This script starts the core services for exhibition demonstration

echo "🚀 Starting Blytz Auction MVP for Exhibition"
echo "============================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Docker is running
if ! command_exists docker || ! docker info >/dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running or not installed${NC}"
    echo "Please install and start Docker first"
    exit 1
fi

# Check if docker compose is available
if ! docker compose version >/dev/null 2>&1; then
    echo -e "${RED}❌ docker compose is not available${NC}"
    echo "Please ensure Docker Compose is installed"
    exit 1
fi

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
docker compose -f docker-compose.simple.yml down 2>/dev/null || true

# Build and start the services
echo "🏗️ Building and starting services..."
echo -e "${BLUE}This may take a few minutes on first run...${NC}"

# Use the simple docker-compose configuration
docker compose -f docker-compose.simple.yml up -d --build

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 10

# Function to check service health
check_service() {
    local service=$1
    local url=$2
    local max_attempts=30
    local attempt=1

    echo -n "🔍 Checking $service... "

    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "$url" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Ready${NC}"
            return 0
        fi
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done

    echo -e "${RED}❌ Failed${NC}"
    return 1
}

# Check each service
echo ""
echo "🔍 Checking service health..."
echo "-----------------------------"

# Check PostgreSQL
if docker compose -f docker-compose.simple.yml exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
    echo -e "${GREEN}✅ PostgreSQL: Ready${NC}"
else
    echo -e "${RED}❌ PostgreSQL: Not ready${NC}"
fi

# Check Redis
if docker compose -f docker-compose.simple.yml exec -T redis redis-cli ping >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Redis: Ready${NC}"
else
    echo -e "${RED}❌ Redis: Not ready${NC}"
fi

# Check Auth Service
check_service "Auth Service" "http://localhost:8084/health"

# Check Auction Service
check_service "Auction Service" "http://localhost:8083/health"

# Check Nginx
check_service "Nginx Gateway" "http://localhost/health"

# Test the frontend
echo ""
echo "🌐 Testing frontend access..."
if curl -s -f "http://localhost/frontend/index.html" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend: Accessible${NC}"
else
    echo -e "${RED}❌ Frontend: Not accessible${NC}"
fi

echo ""
echo "🎯 Final System Status:"
echo "======================"

# Run the final test
echo "Running comprehensive system test..."
./tests/integration/mvp-final-test.sh

echo ""
echo "🌐 Access Points:"
echo "================"
echo -e "${BLUE}Web Interface:${NC} http://localhost/frontend/index.html"
echo -e "${BLUE}Auth Service:${NC} http://localhost:8084/api/v1"
echo -e "${BLUE}Auction Service:${NC} http://localhost:8083/api/v1"
echo -e "${BLUE}Health Check:${NC} http://localhost/health"
echo ""

echo "📋 Quick Test:"
echo "=============="
echo "1. Open http://localhost/frontend/index.html in your browser"
echo "2. Register a new user account"
echo "3. Create a test auction"
echo "4. Place a bid on the auction"
echo "5. Watch the real-time price updates!"
echo ""

echo -e "${GREEN}🎉 Exhibition setup complete!${NC}"
echo -e "${GREEN}The auction platform is ready for visitors!${NC}"
echo ""
echo -e "${YELLOW}💡 Pro Tip:${NC} Use demo@blytz.app / password123 for quick testing"
echo -e "${YELLOW}📱 Mobile:${NC} Perfect for phone/tablet interaction at exhibitions"

# Keep containers running
echo ""
echo "🔄 Services are running in the background..."
echo "To stop: docker compose -f docker-compose.simple.yml down"
echo "To view logs: docker compose -f docker-compose.simple.yml logs -f"
echo "To restart: ./start-exhibition.sh"
echo "To check status: ./tests/integration/mvp-final-test.sh"
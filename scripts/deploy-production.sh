#!/bin/bash

# Blytz Production Deployment with Secrets Manager
# This script deploys to production using secrets manager integration

set -e

echo "🚀 Starting Blytz production deployment with secrets manager..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if .env.production exists
if [ ! -f .env.production ]; then
    print_error ".env.production not found. Please run:"
    print_error "./scripts/setup-production-secrets.sh"
    exit 1
fi

# Verify critical production settings
echo "🔍 Verifying production configuration..."

if ! grep -q "FIUU_SANDBOX=false" .env.production; then
    print_error "FIUU_SANDBOX must be set to false in production"
    exit 1
fi

if ! grep -q "MODE=remote" .env.production; then
    print_error "MODE must be set to remote in production"
    exit 1
fi

if ! grep -q "NODE_ENV=production" .env.production; then
    print_error "NODE_ENV must be set to production"
    exit 1
fi

print_status "Production configuration verified"

# Stop existing services
echo "🛑 Stopping existing services..."
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production down || true

# Pull latest images
echo "📦 Pulling latest images..."
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production pull

# Start production services
echo "🚀 Starting production services..."
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check service health
echo "🔍 Checking service health..."

services=("postgres" "redis" "auth-service" "payment-service" "gateway" "frontend")
for service in "${services[@]}"; do
    if docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production ps $service | grep -q "Up"; then
        print_status "$service is running"
    else
        print_error "$service failed to start"
        docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production logs $service
        exit 1
    fi
done

# Test critical endpoints
echo "🧪 Testing production endpoints..."

# Test API Gateway
if curl -f -s http://localhost:8080/health > /dev/null; then
    print_status "API Gateway is healthy"
else
    print_error "API Gateway health check failed"
fi

# Test Payment Service
if curl -f -s http://localhost:8086/health > /dev/null; then
    print_status "Payment Service is healthy"
else
    print_error "Payment Service health check failed"
fi

# Test Fiuu Configuration
echo "🔧 Testing Fiuu configuration..."
fiuu_config=$(curl -s "http://localhost:8080/api/v1/payments/seamless/config?order_id=TEST123&amount=10050&bill_name=Test%20User&bill_email=test@example.com&bill_mobile=0123456789&bill_desc=Test%20Payment&channel=FPX" || "")

if echo "$fiuu_config" | grep -q '"sandbox":false'; then
    print_status "Fiuu is configured for PRODUCTION mode"
elif echo "$fiuu_config" | grep -q '"sandbox":true'; then
    print_error "Fiuu is still configured for SANDBOX mode"
    echo "Response: $fiuu_config"
else
    print_warning "Could not verify Fiuu configuration"
fi

# Test script URL
if echo "$fiuu_config" | grep -q "api.merchant.fiuu.com.my"; then
    print_status "Fiuu script URL is set to production"
elif echo "$fiuu_config" | grep -q "sandbox.merchant.fiuu.com.my"; then
    print_error "Fiuu script URL is still set to sandbox"
else
    print_warning "Could not verify Fiuu script URL"
fi

echo ""
echo "🎉 Production deployment completed!"
echo ""
echo "🌐 Production URLs:"
echo "   Frontend: https://blytz.app"
echo "   API: https://api.blytz.app"
echo "   Demo: https://demo.blytz.app"
echo "   Seller: https://seller.blytz.app"
echo ""
echo "📊 Monitoring:"
echo "   Check logs: docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production logs -f [service]"
echo "   Check status: docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.production ps"
echo ""
echo "🔧 Environment variables are loaded from .env.production"
echo "   Ensure your secrets manager provides these values in production"
#!/bin/bash

# Blytz Live Auction MVP - Dokploy Deployment Script
set -e

echo "🚀 Starting Blytz production deployment via Dokploy..."

# Check if Dokploy CLI is available
if ! command -v dokploy > /dev/null 2>&1; then
    echo "❌ Dokploy CLI not found. Please install it first:"
    echo "curl -sSL https://dokploy.com/install.sh | sh"
    exit 1
fi

# Generate production secrets if not exists
if [ ! -f .env.production ]; then
    echo "🔐 Generating production secrets..."
    ./scripts/generate-secrets.sh
fi

# Create Dokploy application
echo "📦 Creating Dokploy application..."
dokploy create blytz-live-auction \
    --compose dokploy-compose.yml \
    --env .env.production \
    --domains "blytz.app:3000,api.blytz.app:8080"

# Deploy the application
echo "🚀 Deploying application..."
dokploy deploy blytz-live-auction

# Wait for deployment to complete
echo "⏳ Waiting for deployment to complete..."
sleep 30

# Check deployment status
echo "🔍 Checking deployment status..."
dokploy status blytz-live-auction

# Test endpoints
echo "🧪 Testing production endpoints..."
echo "Testing frontend: https://blytz.app"
curl -f -s -o /dev/null -w "%{http_code}" https://blytz.app || echo "Frontend test failed"

echo "Testing API: https://api.blytz.app/health"
curl -f -s -o /dev/null -w "%{http_code}" https://api.blytz.app/health || echo "API test failed"

echo "✅ Deployment completed!"
echo ""
echo "🌐 Production URLs:"
echo "   Frontend: https://blytz.app"
echo "   API: https://api.blytz.app"
echo ""
echo "📊 Monitor your deployment:"
echo "   Dokploy Dashboard: https://your-dokploy-domain.com"
echo ""
echo "🔧 To view logs: dokploy logs blytz-live-auction"
echo "🔧 To restart: dokploy restart blytz-live-auction"
echo "🔧 To stop: dokploy stop blytz-live-auction"
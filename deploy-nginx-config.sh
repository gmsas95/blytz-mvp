#!/bin/bash

# Blytz Nginx Subdomain Configuration Deployment Script
# This script configures nginx with subdomain routing for all services

set -e

echo "🚀 Starting nginx subdomain configuration deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root or with sudo
if [[ $EUID -eq 0 ]]; then
   print_warning "Running as root - be careful with permissions"
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose could not be found. Please install it first."
    exit 1
fi

# Backup current nginx configuration
print_status "Backing up current nginx configuration..."
if [ -f "nginx.conf" ]; then
    cp nginx.conf "nginx.conf.backup.$(date +%Y%m%d_%H%M%S)"
    print_status "Backup created: nginx.conf.backup.$(date +%Y%m%d_%H%M%S)"
else
    print_warning "No existing nginx.conf found, skipping backup"
fi

# Copy the subdomain configuration
print_status "Installing subdomain nginx configuration..."
cp nginx-subdomains.conf nginx.conf

# Validate nginx configuration syntax
print_status "Validating nginx configuration..."
if docker-compose run --rm nginx nginx -t; then
    print_status "✅ Nginx configuration syntax is valid"
else
    print_error "❌ Nginx configuration has syntax errors"
    print_status "Restoring backup..."
    [ -f "nginx.conf.backup.$(date +%Y%m%d_%H%M%S)" ] && cp "nginx.conf.backup.$(date +%Y%m%d_%H%M%S)" nginx.conf
    exit 1
fi

# Restart nginx container
print_status "Restarting nginx container..."
docker-compose restart nginx

# Wait for nginx to be ready
print_status "Waiting for nginx to be ready..."
sleep 5

# Test nginx health
print_status "Testing nginx health..."
if curl -f http://localhost/health &>/dev/null; then
    print_status "✅ Nginx is responding correctly"
else
    print_warning "⚠️  Nginx health check failed, but configuration was applied"
fi

# Test main domain
print_status "Testing main domain routing..."
if curl -f http://localhost &>/dev/null; then
    print_status "✅ Main domain routing works"
else
    print_warning "⚠️  Main domain routing test failed"
fi

# Show container status
print_status "Checking container status..."
docker-compose ps nginx

echo ""
echo "🎉 Nginx subdomain configuration deployment completed!"
echo ""
echo "📋 Next steps:"
echo "1. Configure DNS records in Cloudflare for all subdomains"
echo "2. Test each subdomain once DNS propagates"
echo "3. Set up SSL certificates if not already done"
echo ""
echo "🔧 Test commands:"
echo "   curl -I http://localhost                    # Main frontend"
echo "   curl -I -H 'Host: api.blytz.app' localhost  # API gateway"
echo "   curl -I -H 'Host: auth.blytz.app' localhost # Auth service"
echo ""
echo "📚 Subdomains configured:"
echo "   blytz.app, api.blytz.app, auth.blytz.app, auction.blytz.app,"
echo "   products.blytz.app, orders.blytz.app, payments.blytz.app,"
echo "   chat.blytz.app, logistics.blytz.app, monitor.blytz.app,"
echo "   analytics.blytz.app, status.blytz.app"
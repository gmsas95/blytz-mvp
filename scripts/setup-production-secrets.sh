#!/bin/bash

# Blytz Production Environment Setup for Secrets Manager Integration
# This script configures production environment without committing secrets to Git

echo "🔧 Setting up Blytz production environment for secrets manager..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
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
    print_warning ".env.production not found. Creating from template..."
    
    # Create .env.production with production settings
    cat > .env.production << 'EOF'
# ===========================================
# Blytz Live Auction Platform - Production Environment
# ===========================================
# IMPORTANT: This file should be managed by your secrets manager
# Do not commit secrets to version control

# ===========================================
# DATABASE CONFIGURATION
# ===========================================
POSTGRES_USER=blytz
POSTGRES_PASSWORD=YOUR_SECURE_POSTGRES_PASSWORD_HERE
POSTGRES_DB=blytz_prod
DATABASE_URL=postgres://blytz:YOUR_SECURE_POSTGRES_PASSWORD_HERE@postgres:5432/blytz_prod?sslmode=disable

# ===========================================
# AUTHENTICATION SECRETS
# ===========================================
JWT_SECRET=YOUR_JWT_SECRET_HERE_MINIMUM_32_CHARACTERS
BETTER_AUTH_SECRET=YOUR_BETTER_AUTH_SECRET_HERE_MINIMUM_32_CHARACTERS
BETTER_AUTH_URL=https://api.blytz.app

# ===========================================
# FIUU PAYMENT GATEWAY CONFIGURATION
# ===========================================
FIUU_MERCHANT_ID=YOUR_FIUU_MERCHANT_ID_HERE
FIUU_VERIFY_KEY=YOUR_FIUU_VERIFY_KEY_HERE
FIUU_SANDBOX=false

# Fiuu API Endpoints
FIUU_API_URL=https://api.merchant.fiuu.com.my
FIUU_SANDBOX_URL=https://sandbox.merchant.fiuu.com.my

# Fiuu Callback URLs
FIUU_RETURN_URL=https://blytz.app/checkout/success
FIUU_NOTIFY_URL=https://api.blytz.app/api/v1/webhooks/fiuu
FIUU_CALLBACK_URL=https://api.blytz.app/api/v1/webhooks/fiuu
FIUU_CANCEL_URL=https://blytz.app/checkout/cancel

# ===========================================
# FRONTEND CONFIGURATION
# ===========================================
NEXT_PUBLIC_API_URL=https://api.blytz.app
NEXT_PUBLIC_APP_URL=https://blytz.app
NEXT_PUBLIC_WS_URL=wss://api.blytz.app
NEXT_PUBLIC_FIUU_SANDBOX=false

# ===========================================
# PRODUCTION SETTINGS
# ===========================================
NODE_ENV=production
ENVIRONMENT=production
MODE=remote
PORT=8080

# ===========================================
# INTERNAL SERVICE URLS (Docker Network)
# ===========================================
AUTH_SERVICE_URL=http://auth-service:8084
PRODUCT_SERVICE_URL=http://product-service:8082
AUCTION_SERVICE_URL=http://auction-service:8083
ORDER_SERVICE_URL=http://order-service:8085
PAYMENT_SERVICE_URL=http://payment-service:8086
LOGISTICS_SERVICE_URL=http://logistics-service:8087
CHAT_SERVICE_URL=http://chat-service:8088

# ===========================================
# REDIS CONFIGURATION
# ===========================================
REDIS_URL=redis://redis:6379

# ===========================================
# SECURITY SETTINGS
# ===========================================
CORS_ORIGINS=https://blytz.app,https://www.blytz.app,https://demo.blytz.app,https://seller.blytz.app
RATE_LIMIT_WINDOW_MS=900000
RATE_LIMIT_MAX_REQUESTS=100

# ===========================================
# MONITORING & LOGGING
# ===========================================
LOG_LEVEL=info
DEBUG=false
ENABLE_SWAGGER=false
ENABLE_PROFILER=false
EOF

    print_status ".env.production template created"
    print_warning "Please update the following values in your secrets manager:"
    print_warning "- POSTGRES_PASSWORD"
    print_warning "- JWT_SECRET" 
    print_warning "- BETTER_AUTH_SECRET"
    print_warning "- FIUU_MERCHANT_ID"
    print_warning "- FIUU_VERIFY_KEY"
else
    print_status ".env.production already exists"
fi

# Ensure critical production settings are correct
echo ""
echo "🔍 Verifying critical production settings..."

# Check FIUU_SANDBOX is set to false
if grep -q "FIUU_SANDBOX=false" .env.production; then
    print_status "FIUU_SANDBOX is correctly set to false"
else
    print_warning "Setting FIUU_SANDBOX=false for production..."
    sed -i 's/FIUU_SANDBOX=.*/FIUU_SANDBOX=false/g' .env.production
fi

# Check NEXT_PUBLIC_FIUU_SANDBOX is set to false
if grep -q "NEXT_PUBLIC_FIUU_SANDBOX=false" .env.production; then
    print_status "NEXT_PUBLIC_FIUU_SANDBOX is correctly set to false"
else
    print_warning "Setting NEXT_PUBLIC_FIUU_SANDBOX=false for production..."
    if grep -q "NEXT_PUBLIC_FIUU_SANDBOX" .env.production; then
        sed -i 's/NEXT_PUBLIC_FIUU_SANDBOX=.*/NEXT_PUBLIC_FIUU_SANDBOX=false/g' .env.production
    else
        echo "NEXT_PUBLIC_FIUU_SANDBOX=false" >> .env.production
    fi
fi

# Check MODE is set to remote
if grep -q "MODE=remote" .env.production; then
    print_status "MODE is correctly set to remote"
else
    print_warning "Setting MODE=remote for production..."
    sed -i 's/MODE=.*/MODE=remote/g' .env.production
fi

# Check NODE_ENV is set to production
if grep -q "NODE_ENV=production" .env.production; then
    print_status "NODE_ENV is correctly set to production"
else
    print_warning "Setting NODE_ENV=production for production..."
    sed -i 's/NODE_ENV=.*/NODE_ENV=production/g' .env.production
fi

echo ""
echo "🚀 Production environment configuration complete!"
echo ""
echo "📋 Next Steps:"
echo "1. Update your secrets manager with the required values:"
echo "   - AWS Secrets Manager, HashiCorp Vault, or similar"
echo "2. Configure your deployment pipeline to use secrets manager"
echo "3. Ensure .env.production is never committed to Git"
echo ""
echo "🔧 For local testing with production settings:"
echo "   cp .env.production .env.local"
echo "   docker compose down && docker compose up -d"
echo ""
echo "🌐 Production URLs will use:"
echo "   - Fiuu Production API: https://api.merchant.fiuu.com.my"
echo "   - Frontend: https://blytz.app"
echo "   - API: https://api.blytz.app"
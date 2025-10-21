#!/bin/bash

# Generate secure secrets for production deployment
echo "🔐 Generating secure secrets for Blytz production deployment..."

# Generate PostgreSQL password
POSTGRES_PASSWORD=$(openssl rand -base64 32)

# Generate JWT secret (minimum 32 characters)
JWT_SECRET=$(openssl rand -base64 64)

# Generate Better Auth secret (minimum 32 characters)
BETTER_AUTH_SECRET=$(openssl rand -base64 64)

# Create production environment file
cat > .env.production << EOF
# Database Configuration
POSTGRES_USER=blytz
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Authentication Secrets
JWT_SECRET=${JWT_SECRET}
BETTER_AUTH_SECRET=${BETTER_AUTH_SECRET}

# API Configuration
NEXT_PUBLIC_API_URL=https://api.blytz.app
NEXT_PUBLIC_WS_URL=wss://api.blytz.app

# Redis Configuration (internal)
REDIS_URL=redis://redis:6379

# Database URLs (internal)
DATABASE_URL=postgres://blytz:${POSTGRES_PASSWORD}@postgres:5432/blytz_prod?sslmode=disable

# Service URLs (internal)
AUTH_SERVICE_URL=http://auth-service:8084
AUCTION_SERVICE_URL=http://auction-service:8083
PRODUCT_SERVICE_URL=http://product-service:8082
ORDER_SERVICE_URL=http://order-service:8085
PAYMENT_SERVICE_URL=http://payment-service:8086
CHAT_SERVICE_URL=http://chat-service:8088
LOGISTICS_SERVICE_URL=http://logistics-service:8087

# Production Settings
NODE_ENV=production
PORT=8080
MODE=remote
EOF

echo "✅ Production environment file created: .env.production"
echo "📝 IMPORTANT: Save these secrets securely!"
echo ""
echo "PostgreSQL Password: ${POSTGRES_PASSWORD}"
echo "JWT Secret: ${JWT_SECRET}"
echo "Better Auth Secret: ${BETTER_AUTH_SECRET}"
echo ""
echo "🔒 Keep this file secure and never commit it to version control!"

# Make the script executable
chmod +x scripts/generate-secrets.sh
#!/bin/bash

# Traefik Configuration Fix Script for Blytz
# This script helps fix the routing issues by applying unified configuration

echo "🔧 Blytz Traefik Configuration Fix"
echo "=================================="

# Check if we're on the server
if [ ! -d "/etc/dokploy" ]; then
    echo "❌ This script must be run on the server with Dokploy installed"
    echo "   Run: ssh root@blytz.app 'bash -s' < traefik-fix.sh"
    exit 1
fi

# Backup existing configurations
echo "📦 Backing up existing configurations..."
mkdir -p /etc/dokploy/traefik/dynamic/backup
cp /etc/dokploy/traefik/dynamic/*.yml /etc/dokploy/traefik/dynamic/backup/ 2>/dev/null || true

# Remove conflicting configuration files
echo "🗑️  Removing conflicting configuration files..."
rm -f /etc/dokploy/traefik/dynamic/dokploy.yml
rm -f /etc/dokploy/traefik/dynamic/middlewares.yml
rm -f /etc/dokploy/traefik/dynamic/traefik.yml

# Apply unified configuration
echo "📝 Applying unified Traefik configuration..."
cat > /etc/dokploy/traefik/dynamic/configuration.yml << 'EOF'
# Unified Traefik Dynamic Configuration for Blytz

http:
  routers:
    # Main Frontend HTTP (redirects to HTTPS)
    frontend-router-http:
      rule: Host(`blytz.app`) || Host(`www.blytz.app`)
      service: frontend-service
      entryPoints:
        - web
      middlewares:
        - redirect-to-https

    # Main Frontend HTTPS
    frontend-router-secure:
      rule: Host(`blytz.app`) || Host(`www.blytz.app`)
      service: frontend-service
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

    # API Gateway Router (for api.blytz.app)
    api-gateway-router-http:
      rule: Host(`api.blytz.app`)
      service: gateway-service
      entryPoints:
        - web
      middlewares:
        - redirect-to-https

    api-gateway-router-secure:
      rule: Host(`api.blytz.app`)
      service: gateway-service
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

    # Demo Interface - Viewer Platform
    demo-router-http:
      rule: Host(`demo.blytz.app`)
      service: demo-service
      entryPoints:
        - web
      middlewares:
        - redirect-to-https

    demo-router-secure:
      rule: Host(`demo.blytz.app`)
      service: demo-service
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

    # Seller Interface - Broadcaster Platform
    seller-router-http:
      rule: Host(`seller.blytz.app`)
      service: seller-service
      entryPoints:
        - web
      middlewares:
        - redirect-to-https

    seller-router-secure:
      rule: Host(`seller.blytz.app`)
      service: seller-service
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

    # Dokploy Management
    dokploy-router-app:
      rule: Host(`sudo.blytz.app`)
      service: dokploy-service-app
      entryPoints:
        - web
      middlewares:
        - redirect-to-https
    dokploy-router-app-secure:
      rule: Host(`sudo.blytz.app`)
      service: dokploy-service-app
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

  services:
    # Main Frontend Service
    frontend-service:
      loadBalancer:
        servers:
          - url: http://frontend:3000
        passHostHeader: true

    # API Gateway Service
    gateway-service:
      loadBalancer:
        servers:
          - url: http://gateway:8080
        passHostHeader: true

    # Demo Frontend Service - Viewer Platform
    demo-service:
      loadBalancer:
        servers:
          - url: http://demo-frontend:3001
        passHostHeader: true

    # Seller Frontend Service - Broadcaster Platform
    seller-service:
      loadBalancer:
        servers:
          - url: http://seller-frontend:3002
        passHostHeader: true

    dokploy-service-app:
      loadBalancer:
        servers:
          - url: http://dokploy:3000
        passHostHeader: true

  middlewares:
    redirect-to-https:
      redirectScheme:
        scheme: https
        permanent: true

    # CORS middleware for API endpoints
    cors-headers:
      headers:
        customResponseHeaders:
          Access-Control-Allow-Origin: "*"
          Access-Control-Allow-Methods: "GET, POST, PUT, DELETE, OPTIONS"
          Access-Control-Allow-Headers: "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization,LiveKit-Token"
          Access-Control-Allow-Credentials: "true"
          Access-Control-Max-Age: "86400"

    # Security headers
    security-headers:
      headers:
        customResponseHeaders:
          X-Frame-Options: "SAMEORIGIN"
          X-Content-Type-Options: "nosniff"
          X-XSS-Protection: "1; mode=block"
          Strict-Transport-Security: "max-age=31536000; includeSubDomains"
EOF

# Restart Traefik
echo "🔄 Restarting Traefik to apply configuration..."
docker restart traefik 2>/dev/null || docker-compose restart traefik 2>/dev/null || echo "⚠️  Could not restart Traefik automatically"

# Wait for Traefik to restart
echo "⏳ Waiting for Traefik to restart..."
sleep 10

# Test the configuration
echo "🧪 Testing configuration..."
echo "Testing frontend..."
curl -s -o /dev/null -w "%{http_code}" http://localhost/health || echo "❌ Frontend not responding"

echo "Testing API gateway..."
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health || echo "❌ API gateway not responding"

echo ""
echo "✅ Configuration fix completed!"
echo ""
echo "📋 Next steps:"
echo "1. Check if services are responding: curl -s https://api.blytz.app/health"
echo "2. Check frontend: curl -s https://blytz.app"
echo "3. If still not working, check Docker logs: docker logs traefik"
echo "4. Check service status: docker ps"
echo ""
echo "📁 Backups saved to: /etc/dokploy/traefik/dynamic/backup/"
#!/bin/bash
# Setup Demo and Seller Interfaces with LiveKit for Blytz Live Auction

echo "🚀 Setting up Demo and Seller interfaces with LiveKit..."

# 1. Update Traefik configuration with new subdomains
echo "📋 Updating Traefik configuration..."
sudo cp /etc/dokploy/traefik/dynamic/dokploy.yml /etc/dokploy/traefik/dynamic/dokploy.yml.backup

# Add the new subdomain routes to existing config
cat > /tmp/dokploy-with-demo-seller.yml << 'EOF'
http:
  routers:
    # Frontend HTTP (redirects to HTTPS)
    frontend-router-http:
      rule: Host(`blytz.app`) || Host(`www.blytz.app`)
      service: frontend-service
      entryPoints:
        - web
      middlewares:
        - redirect-to-https

    # Frontend HTTPS (main domain)
    frontend-router-secure:
      rule: Host(`blytz.app`) || Host(`www.blytz.app`)
      service: frontend-service
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

    # Demo Interface - Viewer Platform with LiveKit Client
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

    # Seller Interface - Broadcaster Platform with LiveKit SDK
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
      rule: Host(`dokploy.docker.localhost`) && PathPrefix(`/`) || Host(`sudo.blytz.app`)
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
          - url: http://blytz-frontend-prod:3000
        passHostHeader: true

    # Demo Interface Service - will be created for LiveKit client
    demo-service:
      loadBalancer:
        servers:
          - url: http://blytz-demo-prod:3001
        passHostHeader: true

    # Seller Interface Service - will be created for LiveKit broadcaster
    seller-service:
      loadBalancer:
        servers:
          - url: http://blytz-seller-prod:3002
        passHostHeader: true

    dokploy-service-app:
      loadBalancer:
        servers:
          - url: http://dokploy:3000
        passHostHeader: true
EOF

sudo mv /tmp/dokploy-with-demo-seller.yml /etc/dokploy/traefik/dynamic/dokploy.yml

# 2. Restart Traefik to pick up new routes
echo "🔄 Restarting Traefik..."
docker restart dokploy-traefik

# 3. Wait for Traefik to restart
echo "⏳ Waiting for Traefik to restart..."
sleep 10

# 4. Test new subdomains
echo "🧪 Testing new subdomains..."
echo "Testing demo.blytz.app..."
curl -I https://demo.blytz.app
echo "Testing seller.blytz.app..."
curl -I https://seller.blytz.app

echo "✅ Demo and Seller subdomains configured!"
echo ""
echo "🎯 Next steps:"
echo "1. Create frontend-demo and frontend-seller directories with Next.js apps"
echo "2. Set up LiveKit server configuration"
echo "3. Build and deploy the new frontend containers"
echo "4. Configure LiveKit client SDK for demo interface"
echo "5. Configure LiveKit broadcaster SDK for seller interface"
echo ""
echo "📚 LiveKit Documentation: https://docs.livekit.io/"
#!/bin/bash
# Setup Demo and Seller Interfaces on Hostinger VPS

echo "🚀 Setting up Demo and Seller interfaces on Hostinger VPS..."

# 1. Backup current configuration
echo "📋 Backing up current configuration..."
cp docker-compose.yml docker-compose.yml.backup
cp /etc/dokploy/traefik/dynamic/dokploy.yml /etc/dokploy/traefik/dynamic/dokploy.yml.backup 2>/dev/null || echo "No existing dokploy.yml to backup"

# 2. Copy new configuration files
echo "📁 Copying configuration files..."
cp docker-compose-hostinger.yml docker-compose.yml
cp livekit-config-hostinger.yaml livekit-config.yaml
cp traefik-additional-routes.yml /tmp/traefik-additional-routes.yml

# 3. Update Traefik configuration
echo "🔧 Updating Traefik configuration..."
# Merge the additional routes with existing config
cat /tmp/traefik-additional-routes.yml >> /etc/dokploy/traefik/dynamic/dokploy.yml

# 4. Restart Traefik to pick up new routes
echo "🔄 Restarting Traefik..."
docker restart dokploy-traefik

# 5. Wait for Traefik to restart
echo "⏳ Waiting for Traefik to restart..."
sleep 10

# 6. Deploy new services
echo "🚀 Deploying new services..."
docker compose up -d --build demo-frontend seller-frontend livekit-server livekit-redis

# 7. Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 30

# 8. Test new subdomains
echo "🧪 Testing new subdomains..."
echo "Testing demo.blytz.app..."
curl -I https://demo.blytz.app
echo "Testing seller.blytz.app..."
curl -I https://seller.blytz.app
echo "Testing livekit.blytz.app..."
curl -I https://livekit.blytz.app

echo "✅ Setup complete!"
echo ""
echo "🎯 Next steps:"
echo "1. Create frontend-demo and frontend-seller directories with Next.js apps"
echo "2. Build and deploy the frontend applications"
echo "3. Configure LiveKit client SDK in demo interface"
echo "4. Configure LiveKit broadcaster SDK in seller interface"
echo ""
echo "📚 LiveKit Documentation: https://docs.livekit.io/"
echo "🔧 LiveKit React Components: https://docs.livekit.io/references/components/react/"
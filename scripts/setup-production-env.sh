#!/bin/bash

# Blytz Production Environment Setup Script
# This script ensures production environment has correct FIUU_SANDBOX setting

echo "🔧 Setting up production environment for Blytz..."

# Check if .env.production exists, if not create it from example
if [ ! -f .env.production ]; then
    echo "📝 Creating .env.production from example..."
    cp .env.production.example .env.production
    
    # Update FIUU_SANDBOX to false for production
    sed -i 's/FIUU_SANDBOX=true/FIUU_SANDBOX=false/g' .env.production
    
    # Add NEXT_PUBLIC_FIUU_SANDBOX if not present
    if ! grep -q "NEXT_PUBLIC_FIUU_SANDBOX" .env.production; then
        echo "NEXT_PUBLIC_FIUU_SANDBOX=false" >> .env.production
    fi
    
    echo "✅ .env.production created with production settings"
else
    echo "📝 Updating existing .env.production..."
    # Ensure FIUU_SANDBOX is set to false
    sed -i 's/FIUU_SANDBOX=.*/FIUU_SANDBOX=false/g' .env.production
    
    # Ensure NEXT_PUBLIC_FIUU_SANDBOX is set to false
    if grep -q "NEXT_PUBLIC_FIUU_SANDBOX" .env.production; then
        sed -i 's/NEXT_PUBLIC_FIUU_SANDBOX=.*/NEXT_PUBLIC_FIUU_SANDBOX=false/g' .env.production
    else
        echo "NEXT_PUBLIC_FIUU_SANDBOX=false" >> .env.production
    fi
    
    echo "✅ .env.production updated with production settings"
fi

# Also update current .env for local testing
if [ -f .env ]; then
    echo "📝 Updating current .env for production testing..."
    sed -i 's/FIUU_SANDBOX=.*/FIUU_SANDBOX=false/g' .env
    
    if grep -q "NEXT_PUBLIC_FIUU_SANDBOX" .env; then
        sed -i 's/NEXT_PUBLIC_FIUU_SANDBOX=.*/NEXT_PUBLIC_FIUU_SANDBOX=false/g' .env
    else
        echo "NEXT_PUBLIC_FIUU_SANDBOX=false" >> .env
    fi
    
    echo "✅ Current .env updated for production testing"
fi

echo ""
echo "🚀 Production environment is now configured!"
echo "📋 Summary of changes:"
echo "   - FIUU_SANDBOX=false (backend)"
echo "   - NEXT_PUBLIC_FIUU_SANDBOX=false (frontend)"
echo "   - Production Fiuu URLs will be used"
echo ""
echo "🔄 To apply changes:"
echo "   1. Restart your services: docker-compose down && docker-compose up -d"
echo "   2. Or redeploy: ./scripts/deploy-dokploy.sh"
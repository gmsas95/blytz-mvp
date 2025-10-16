#!/bin/bash
# Firebase Local Development Setup for Blytz MVP
# This sets up Firebase for local development without requiring a real Firebase project

set -e

echo "🔥 Setting up Firebase for local development..."
echo "=============================================="

# Navigate to functions directory
cd functions/

echo "📦 Installing Firebase dependencies..."
npm install

echo "🏗️ Building TypeScript functions..."
npm run build

echo "🔧 Setting up local Firebase configuration..."

# Create local Firebase config for development
cat > .runtimeconfig.json << 'EOF'
{
  "stripe": {
    "secret_key": "sk_test_demo",
    "webhook_secret": "whsec_demo"
  },
  "jwt": {
    "secret": "your-jwt-secret-for-development"
  }
}
EOF

echo "🚀 Starting Firebase emulators for local development..."
echo ""
echo "Available functions:"
echo "• createUser - Create new users"
echo "• validateToken - Validate JWT tokens"
echo "• createPaymentIntent - Create payment intents"
echo "• placeBid - Place bids on auctions"
echo "• createAuction - Create new auctions"
echo ""
echo "Local endpoints:"
echo "• http://localhost:5001/demo-blytz-mvp/us-central1/health"
echo "• http://localhost:5001/demo-blytz-mvp/us-central1/createUser"
echo "• http://localhost:5001/demo-blytz-mvp/us-central1/placeBid"
echo ""
echo "Starting Firebase emulators..."
echo "Press Ctrl+C to stop"
echo "=============================================="

# Start Firebase emulators
firebase emulators:start --only functions --project demo-blytz-mvp
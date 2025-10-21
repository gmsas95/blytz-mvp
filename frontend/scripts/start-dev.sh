#!/bin/bash

# Blytz Frontend Development Startup Script

echo "🚀 Starting Blytz Frontend Development Server..."
echo "Environment: MODE=${MODE:-mock}"
echo "API Base: ${NEXT_PUBLIC_API_BASE:-http://localhost:8080}"
echo ""

# Check if Node.js is available
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ to continue."
    exit 1
fi

# Check if npm is available
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm to continue."
    exit 1
fi

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
    if [ $? -ne 0 ]; then
        echo "❌ Failed to install dependencies."
        exit 1
    fi
fi

# Set default environment variables
export MODE=${MODE:-mock}
export NEXT_PUBLIC_API_BASE=${NEXT_PUBLIC_API_BASE:-http://localhost:8080}

# Start the development server
echo "🌟 Starting Next.js development server..."
echo "📱 Frontend will be available at: http://localhost:3000"
echo "🔍 API Mode: $MODE"
echo ""

npm run dev
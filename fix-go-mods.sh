#!/bin/bash
set -e

echo "🔧 Fixing and tidying Go modules across all services..."

# Define services
SERVICES=(
  "backend/shared"
  "backend/auth-service"
  "backend/product-service"
  "backend/auction-service"
  "backend/chat-service"
)

# Ensure Go is installed
if ! command -v go &> /dev/null; then
  echo "❌ Go not found. Please install Go 1.23+ before running this script."
  exit 1
fi

# Loop through services
for SERVICE in "${SERVICES[@]}"; do
  if [ -d "$SERVICE" ]; then
    echo "--------------------------------------------"
    echo "📦 Processing: $SERVICE"
    echo "--------------------------------------------"
    cd "$SERVICE"

    # Ensure go.mod exists
    if [ ! -f "go.mod" ]; then
      echo "⚙️  Initializing go.mod..."
      go mod init "blytz/$(basename "$SERVICE")"
    fi

    echo "🧹 Running go mod tidy..."
    go mod tidy

    echo "🔍 Verifying dependencies..."
    go mod verify

    cd - > /dev/null
  else
    echo "⚠️  Directory not found: $SERVICE"
  fi
done

echo "✅ All Go modules have been tidied and verified successfully!"

#!/bin/bash
set -e

echo "🔧 Fixing and tidying Go modules across all microservices..."

# Define services
SERVICES=(
  "services/shared"
  "services/auth-service"
  "services/product-service"
  "services/auction-service"
  "services/chat-service"
  "services/gateway"
  "services/order-service"
  "services/payment-service"
  "services/logistics-service"
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

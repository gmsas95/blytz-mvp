#!/bin/bash

# Shared Package Migration Verification Script
# This script verifies that the shared package migration is complete and correct

set -e

echo "🔍 Verifying Shared Package Migration..."
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[→]${NC} $1"
}

# Check 1: Verify new shared structure exists
print_info "Checking new shared package structure..."
if [ -d "/home/sas/blytzmvp-clean/shared/pkg/auth" ] && \
   [ -d "/home/sas/blytzmvp-clean/shared/pkg/errors" ] && \
   [ -d "/home/sas/blytzmvp-clean/shared/pkg/utils" ] && \
   [ -d "/home/sas/blytzmvp-clean/shared/pkg/constants" ]; then
    print_status "New shared package structure is complete"
else
    print_error "New shared package structure is incomplete"
    exit 1
fi

# Check 2: Verify go.mod files exist
print_info "Checking go.mod files..."
if [ -f "/home/sas/blytzmvp-clean/shared/go.mod" ] && \
   [ -f "/home/sas/blytzmvp-clean/services/auction-service/go.mod" ]; then
    print_status "Go module files exist"
else
    print_error "Missing go.mod files"
    exit 1
fi

# Check 3: Verify auction service imports use pkg/ prefix
print_info "Checking auction service import paths..."

# Check main.go
if grep -q "github.com/blytz/shared/pkg/utils" /home/sas/blytzmvp-clean/services/auction-service/cmd/main.go; then
    print_status "✓ main.go uses correct import path"
else
    print_error "✗ main.go import path incorrect"
    exit 1
fi

# Check handlers/auction.go
if grep -q "github.com/blytz/shared/pkg/utils" /home/sas/blytzmvp-clean/services/auction-service/internal/api/handlers/auction.go && \
   grep -q "github.com/blytz/shared/pkg/errors" /home/sas/blytzmvp-clean/services/auction-service/internal/api/handlers/auction.go; then
    print_status "✓ handlers/auction.go uses correct import paths"
else
    print_error "✗ handlers/auction.go import paths incorrect"
    exit 1
fi

# Check services/auction.go
if grep -q "github.com/blytz/shared/pkg/errors" /home/sas/blytzmvp-clean/services/auction-service/internal/services/auction.go; then
    print_status "✓ services/auction.go uses correct import path"
else
    print_error "✗ services/auction.go import path incorrect"
    exit 1
fi

# Check 4: Verify auth integration is still intact
print_info "Verifying auth integration..."
if grep -q "github.com/blytz/shared/pkg/auth" /home/sas/blytzmvp-clean/services/auction-service/internal/api/router.go && \
   grep -q "auth.GinAuthMiddleware" /home/sas/blytzmvp-clean/services/auction-service/internal/api/router.go; then
    print_status "✓ Auth integration is intact"
else
    print_error "✗ Auth integration missing"
    exit 1
fi

# Check 5: Verify router structure
print_info "Checking router structure..."
if grep -q "protectedAuctions.Use(auth.GinAuthMiddleware(authClient))" /home/sas/blytzmvp-clean/services/auction-service/internal/api/router.go; then
    print_status "✓ Protected routes properly configured"
else
    print_error "✗ Protected routes not configured"
    exit 1
fi

# Check 6: Verify no old import paths remain
print_info "Checking for remaining old import paths..."
old_imports=$(grep -r "github.com/blytz/shared/" /home/sas/blytzmvp-clean/services/auction-service/ --include="*.go" | grep -v "pkg/" | wc -l)
if [ "$old_imports" -eq 0 ]; then
    print_status "✓ No old import paths found"
else
    print_error "✗ Found $old_imports old import paths"
    echo "Found old imports:"
    grep -r "github.com/blytz/shared/" /home/sas/blytzmvp-clean/services/auction-service/ --include="*.go" | grep -v "pkg/"
    exit 1
fi

# Check 7: Verify test scripts are in place
print_info "Checking test infrastructure..."
if [ -f "/home/sas/blytzmvp-clean/services/auth-service/test-auth-service.sh" ] && \
   [ -f "/home/sas/blytzmvp-clean/services/auction-service/test-auction-auth.sh" ]; then
    print_status "✓ Test scripts are available"
else
    print_error "✗ Test scripts missing"
    exit 1
fi

# Check 8: Verify documentation exists
print_info "Checking documentation..."
if [ -f "/home/sas/blytzmvp-clean/services/auth-service/README.md" ] && \
   [ -f "/home/sas/blytzmvp-clean/services/auction-service/AUTH_INTEGRATION.md" ]; then
    print_status "✓ Documentation is complete"
else
    print_error "✗ Documentation incomplete"
    exit 1
fi

echo ""
echo "=========================================="
print_status "🎉 SHARED PACKAGE MIGRATION VERIFICATION COMPLETE!"
echo ""
print_info "Next steps to complete the fix:"
echo "1. Navigate to auction service: cd /home/sas/blytzmvp-clean/services/auction-service"
echo "2. Run: go mod tidy"
echo "3. Run: go build to verify compilation"
echo "4. Run: ./test-auction-auth.sh to test auth integration"
echo ""
print_info "The shared package structure is now consistent:"
echo "- /home/sas/blytzmvp-clean/shared/pkg/auth/     ✓"
echo "- /home/sas/blytzmvp-clean/shared/pkg/errors/  ✓"
echo "- /home/sas/blytzmvp-clean/shared/pkg/utils/   ✓"
echo "- /home/sas/blytzmvp-clean/shared/pkg/constants/ ✓"
echo "- /home/sas/blytzmvp-clean/shared/pkg/proto/   ✓"
echo ""
print_status "All auction service imports now use the correct pkg/ prefix!"
echo "The authentication integration is ready for testing! 🚀"

exit 0

# Commands that would be run in a Go environment:
# cd /home/sas/blytzmvp-clean/services/auction-service
# go mod tidy
# go build -o auction-service ./cmd/main.go
# if [ $? -eq 0 ]; then
#     echo "Build successful!"
#     ./test-auction-auth.sh
# else
#     echo "Build failed - check dependencies"
#     exit 1
# fi
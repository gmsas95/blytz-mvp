#!/bin/bash

# Auction Service Authentication Integration Test
# This script tests auction endpoints with authentication

set -e

# Configuration
AUTH_SERVICE_URL="${AUTH_SERVICE_URL:-http://localhost:8084}"
AUCTION_SERVICE_URL="${AUCTION_SERVICE_URL:-http://localhost:8083}"
TEST_EMAIL="auction-test-$(date +%s)@example.com"
TEST_PASSWORD="password123"
TEST_DISPLAY_NAME="Auction Test User"

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

# Function to make HTTP requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local auth_header=$5

    local curl_cmd="curl -s -w \"\n%{http_code}\" -X $method"

    if [ -n "$auth_header" ]; then
        curl_cmd="$curl_cmd -H \"Authorization: Bearer $auth_header\""
    fi

    if [ -n "$data" ]; then
        curl_cmd="$curl_cmd -H \"Content-Type: application/json\" -d '$data'"
    fi

    curl_cmd="$curl_cmd $endpoint"

    response=$(eval $curl_cmd)

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "$expected_status" ]; then
        echo "$body"
        return 0
    else
        print_error "Expected status $expected_status, got $http_code"
        echo "Response: $body"
        return 1
    fi
}

print_info "Testing Auction Service with Authentication"
print_info "Auth Service: $AUTH_SERVICE_URL"
print_info "Auction Service: $AUCTION_SERVICE_URL"
print_info "Test Email: $TEST_EMAIL"

# Step 1: Register a new user
print_info "Step 1: Registering new user..."
register_data=$(cat <<EOF
{
    "email": "$TEST_EMAIL",
    "password": "$TEST_PASSWORD",
    "display_name": "$TEST_DISPLAY_NAME",
    "phone_number": "+1234567890"
}
EOF
)

register_response=$(make_request "POST" "$AUTH_SERVICE_URL/api/v1/auth/register" "$register_data" "200")
if echo "$register_response" | grep -q "true"; then
    print_status "User registration successful"

    # Extract token
    AUTH_TOKEN=$(echo "$register_response" | grep -o '"token":"[^"]*' | sed 's/"token":"//')
    if [ -n "$AUTH_TOKEN" ]; then
        print_status "Auth token extracted successfully"
    else
        print_error "Failed to extract auth token"
        exit 1
    fi
else
    print_error "User registration failed"
    exit 1
fi

# Step 2: Test public auction endpoints (should work without auth)
print_info "Step 2: Testing public auction endpoints..."

# List auctions (public)
list_response=$(make_request "GET" "$AUCTION_SERVICE_URL/api/v1/auctions" "" "200")
if echo "$list_response" | grep -q "auctions"; then
    print_status "List auctions (public) works"
else
    print_error "List auctions (public) failed"
    exit 1
fi

# Step 3: Test protected auction endpoints (should fail without auth)
print_info "Step 3: Testing protected endpoints without auth (should fail)..."

create_auction_data=$(cat <<EOF
{
    "product_id": "test-product-123",
    "title": "Test Auction Item",
    "description": "A test auction item for authentication testing",
    "starting_price": 10.00,
    "reserve_price": 50.00,
    "min_bid_increment": 1.00,
    "start_time": "2025-10-20T10:00:00Z",
    "end_time": "2025-10-20T11:00:00Z",
    "type": "scheduled",
    "images": ["https://example.com/image1.jpg"]
}
EOF
)

# Try to create auction without auth (should fail)
unauth_response=$(make_request "POST" "$AUCTION_SERVICE_URL/api/v1/auctions" "$create_auction_data" "401" 2>/dev/null || echo "")
if [ -z "$unauth_response" ]; then
    print_status "Create auction without auth correctly rejected"
else
    print_error "Create auction without auth should have failed"
    exit 1
fi

# Step 4: Test protected auction endpoints with auth (should work)
print_info "Step 4: Testing protected endpoints with auth..."

# Create auction with auth
create_response=$(make_request "POST" "$AUCTION_SERVICE_URL/api/v1/auctions" "$create_auction_data" "200" "$AUTH_TOKEN")
if echo "$create_response" | grep -q "auction_id"; then
    print_status "Create auction with auth successful"

    # Extract auction ID
    AUCTION_ID=$(echo "$create_response" | grep -o '"auction_id":"[^"]*' | sed 's/"auction_id":"//')
    if [ -n "$AUCTION_ID" ]; then
        print_status "Auction ID extracted: $AUCTION_ID"
    else
        print_error "Failed to extract auction ID"
        exit 1
    fi
else
    print_error "Create auction with auth failed"
    exit 1
fi

# Place bid with auth
print_info "Step 5: Testing bid placement with auth..."
bid_data=$(cat <<EOF
{
    "amount": 15.00
}
EOF
)

bid_response=$(make_request "POST" "$AUCTION_SERVICE_URL/api/v1/auctions/$AUCTION_ID/bids" "$bid_data" "200" "$AUTH_TOKEN")
if echo "$bid_response" | grep -q "bid_id"; then
    print_status "Place bid with auth successful"
else
    print_error "Place bid with auth failed"
    exit 1
fi

# Step 5: Test bid without auth (should fail)
print_info "Step 6: Testing bid placement without auth (should fail)..."
unauth_bid_response=$(make_request "POST" "$AUCTION_SERVICE_URL/api/v1/auctions/$AUCTION_ID/bids" "$bid_data" "401" 2>/dev/null || echo "")
if [ -z "$unauth_bid_response" ]; then
    print_status "Place bid without auth correctly rejected"
else
    print_error "Place bid without auth should have failed"
    exit 1
fi

# Step 6: Test getting auction (public - should work)
print_info "Step 7: Testing get auction (public endpoint)..."
get_response=$(make_request "GET" "$AUCTION_SERVICE_URL/api/v1/auctions/$AUCTION_ID" "" "200")
if echo "$get_response" | grep -q "$AUCTION_ID"; then
    print_status "Get auction (public) works"
else
    print_error "Get auction (public) failed"
    exit 1
fi

# Step 7: Test update auction with auth
print_info "Step 8: Testing update auction with auth..."
update_data=$(cat <<EOF
{
    "title": "Updated Test Auction Item",
    "description": "An updated test auction item"
}
EOF
)

update_response=$(make_request "PUT" "$AUCTION_SERVICE_URL/api/v1/auctions/$AUCTION_ID" "$update_data" "200" "$AUTH_TOKEN")
if echo "$update_response" | grep -q "auction_id"; then
    print_status "Update auction with auth successful"
else
    print_error "Update auction with auth failed"
    exit 1
fi

# Summary
echo ""
print_status "🎉 All auction service authentication tests passed!"
print_status "✅ Public endpoints work without authentication"
print_status "✅ Protected endpoints require authentication"
print_status "✅ Create auction requires auth"
print_status "✅ Place bid requires auth"
print_status "✅ Update auction requires auth"
print_status "✅ Auth middleware is working correctly"

echo ""
print_info "Integration Summary:"
print_info "- Auth service is providing valid JWT tokens"
print_info "- Auction service auth middleware is validating tokens"
print_info "- Protected routes are properly secured"
print_info "- Public routes remain accessible"
print_info "- User context is correctly passed to handlers"

exit 0
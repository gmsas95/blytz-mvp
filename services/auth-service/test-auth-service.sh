#!/bin/bash

# Blytz Auth Service Integration Test Script
# This script tests the complete authentication flow

set -e

# Configuration
AUTH_SERVICE_URL="${AUTH_SERVICE_URL:-http://localhost:8084}"
TEST_EMAIL="test-$(date +%s)@example.com"
TEST_PASSWORD="password123"
TEST_DISPLAY_NAME="Test User"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
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

    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$AUTH_SERVICE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            "$AUTH_SERVICE_URL$endpoint")
    fi

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

print_info "Testing Blytz Authentication Service"
print_info "Service URL: $AUTH_SERVICE_URL"
print_info "Test Email: $TEST_EMAIL"

# Test 1: Health Check
print_info "Testing health check endpoint..."
health_response=$(make_request "GET" "/health" "" "200")
if echo "$health_response" | grep -q "healthy"; then
    print_status "Health check passed"
else
    print_error "Health check failed"
    exit 1
fi

# Test 2: User Registration
print_info "Testing user registration..."
register_data=$(cat <<EOF
{
    "email": "$TEST_EMAIL",
    "password": "$TEST_PASSWORD",
    "display_name": "$TEST_DISPLAY_NAME",
    "phone_number": "+1234567890"
}
EOF
)

register_response=$(make_request "POST" "/api/v1/auth/register" "$register_data" "200")
if echo "$register_response" | grep -q "true"; then
    print_status "User registration successful"

    # Extract token from response
    TOKEN=$(echo "$register_response" | grep -o '"token":"[^"]*' | sed 's/"token":"//')
    if [ -n "$TOKEN" ]; then
        print_status "Token extracted successfully"
    else
        print_error "Failed to extract token from registration response"
        exit 1
    fi
else
    print_error "User registration failed"
    exit 1
fi

# Test 3: User Login
print_info "Testing user login..."
login_data=$(cat <<EOF
{
    "email": "$TEST_EMAIL",
    "password": "$TEST_PASSWORD"
}
EOF
)

login_response=$(make_request "POST" "/api/v1/auth/login" "$login_data" "200")
if echo "$login_response" | grep -q "true"; then
    print_status "User login successful"

    # Extract login token
    LOGIN_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | sed 's/"token":"//')
    if [ -n "$LOGIN_TOKEN" ]; then
        print_status "Login token extracted successfully"
    else
        print_error "Failed to extract token from login response"
        exit 1
    fi
else
    print_error "User login failed"
    exit 1
fi

# Test 4: Token Validation
print_info "Testing token validation..."
validate_data=$(cat <<EOF
{
    "token": "$LOGIN_TOKEN"
}
EOF
)

validate_response=$(make_request "POST" "/api/v1/auth/validate" "$validate_data" "200")
if echo "$validate_response" | grep -q "true" && echo "$validate_response" | grep -q "$TEST_EMAIL"; then
    print_status "Token validation successful"
else
    print_error "Token validation failed"
    exit 1
fi

# Test 5: Token Refresh
print_info "Testing token refresh..."
refresh_data=$(cat <<EOF
{
    "refresh_token": "$LOGIN_TOKEN"
}
EOF
)

refresh_response=$(make_request "POST" "/api/v1/auth/refresh" "$refresh_data" "200")
if echo "$refresh_response" | grep -q "true"; then
    print_status "Token refresh successful"
else
    print_error "Token refresh failed"
    exit 1
fi

# Test 6: Invalid Login (should fail)
print_info "Testing invalid login (should fail)..."
invalid_login_data=$(cat <<EOF
{
    "email": "$TEST_EMAIL",
    "password": "wrongpassword"
}
EOF
)

invalid_login_response=$(make_request "POST" "/api/v1/auth/login" "$invalid_login_data" "401")
if [ $? -eq 0 ]; then
    print_status "Invalid login correctly rejected"
else
    print_error "Invalid login test failed"
    exit 1
fi

# Test 7: Invalid Token Validation (should fail)
print_info "Testing invalid token validation (should fail)..."
invalid_validate_data=$(cat <<EOF
{
    "token": "invalid.token.here"
}
EOF
)

invalid_validate_response=$(make_request "POST" "/api/v1/auth/validate" "$invalid_validate_data" "200")
if echo "$invalid_validate_response" | grep -q "false"; then
    print_status "Invalid token correctly rejected"
else
    print_error "Invalid token test failed"
    exit 1
fi

# Test 8: Duplicate Registration (should fail)
print_info "Testing duplicate registration (should fail)..."
duplicate_response=$(make_request "POST" "/api/v1/auth/register" "$register_data" "400")
if [ $? -eq 0 ]; then
    print_status "Duplicate registration correctly rejected"
else
    print_error "Duplicate registration test failed"
    exit 1
fi

# Summary
print_info "Running performance test..."
start_time=$(date +%s%N)

# Make 10 rapid requests
for i in {1..10}; do
    make_request "POST" "/api/v1/auth/validate" "$validate_data" "200" > /dev/null
done

end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 ))  # Convert to milliseconds
avg_response_time=$((duration / 10))

print_status "Performance test completed"
print_status "Average response time: ${avg_response_time}ms"

# Final summary
echo ""
print_status "🎉 All authentication tests passed successfully!"
print_status "Service is working correctly and ready for integration"
print_status "Average response time: ${avg_response_time}ms (target: <50ms)"

# Cleanup note
echo ""
print_info "Test completed. User created with email: $TEST_EMAIL"
print_info "You can now integrate this auth service with other microservices"

exit 0
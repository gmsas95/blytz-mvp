#!/bin/bash

# Complete Firebase Integration Test Suite
# Fixed version with proper v2 API format

set -e

echo "🚀 COMPLETE Firebase Integration Test Suite"
echo "==========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Firebase emulators are running
echo "🔍 Checking Firebase emulators..."
if ! curl -s http://localhost:5001/demo-blytz-mvp/us-central1/health > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Firebase emulators not running, starting them...${NC}"
    cd functions && npm run serve > /dev/null 2>&1 &
    sleep 5
    cd ..
fi

# Wait for emulators to be ready
echo "⏳ Waiting for Firebase emulators..."
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:5001/demo-blytz-mvp/us-central1/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Firebase emulators are running${NC}"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    sleep 1
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}❌ Failed to start Firebase emulators${NC}"
    exit 1
fi

# Test Results Summary
TESTS_PASSED=0
TESTS_FAILED=0

# Function to test Firebase function
test_firebase_function() {
    local function_name=$1
    local test_name=$2
    local request_data=$3
    local success_pattern=$4

    echo ""
    echo -e "${BLUE}Testing $test_name...${NC}"

    RESPONSE=$(curl -s -X POST "http://localhost:5001/demo-blytz-mvp/us-central1/$function_name" \
        -H "Content-Type: application/json" \
        -d "$request_data")

    if echo "$RESPONSE" | grep -q "$success_pattern"; then
        echo -e "${GREEN}✅ $test_name PASSED${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}❌ $test_name FAILED${NC}"
        echo "Response: $RESPONSE"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Test 1: Health Check
test_firebase_function "health" "Health Check" \
    '{"data": {}}' \
    '"status":"healthy"'

# Test 2: User Creation
test_firebase_function "createUser" "User Creation" \
    '{
        "data": {
            "email": "testuser@example.com",
            "password": "testpassword123",
            "displayName": "Test User",
            "phoneNumber": "+1234567890"
        }
    }' \
    '"success":true'

# Test 3: Auction Creation
test_firebase_function "createAuction" "Auction Creation" \
    '{
        "data": {
            "title": "Vintage Rolex Submariner",
            "description": "Authentic vintage Rolex from 1965",
            "startingPrice": 5000.00,
            "duration": 24,
            "category": "watches",
            "images": ["watch1.jpg", "watch2.jpg"]
        }
    }' \
    '"success":true'

# Test 4: Bid Placement
test_firebase_function "placeBid" "Bid Placement" \
    '{
        "data": {
            "auctionId": "test-auction-123",
            "amount": 5500.00
        }
    }' \
    '"success":true'

# Test 5: Payment Intent Creation
test_firebase_function "createPaymentIntent" "Payment Intent Creation" \
    '{
        "data": {
            "amount": 5500.00,
            "currency": "usd",
            "auctionId": "test-auction-123",
            "bidId": "test-bid-456"
        }
    }' \
    '"success":true'

# Test 6: Send Notification
test_firebase_function "sendNotification" "Notification Sending" \
    '{
        "data": {
            "userId": "test-user-id",
            "title": "New Bid Alert",
            "body": "Someone placed a bid on your auction!",
            "data": {
                "type": "bid_alert",
                "auctionId": "test-auction-123"
            }
        }
    }' \
    '"success":true'

# Test 7: Get Auction Details
test_firebase_function "getAuctionDetails" "Get Auction Details" \
    '{
        "data": {
            "auctionId": "test-auction-123"
        }
    }' \
    '"success":true'

# Test 8: Validate Token
test_firebase_function "validateToken" "Token Validation" \
    '{
        "data": {}
    }' \
    '"success"'

# Test 9: Error Handling - Invalid Data
echo ""
echo -e "${BLUE}Testing Error Handling...${NC}"

# Test invalid auction data
ERROR_RESPONSE=$(curl -s -X POST "http://localhost:5001/demo-blytz-mvp/us-central1/createAuction" \
    -H "Content-Type: application/json" \
    -d '{
        "data": {
            "title": "",
            "description": "Invalid auction",
            "startingPrice": -100,
            "duration": 0,
            "category": "test"
        }
    }')

if echo "$ERROR_RESPONSE" | grep -q "error"; then
    echo -e "${GREEN}✅ Error Handling PASSED - Correctly rejected invalid data${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}❌ Error Handling FAILED - Should have rejected invalid data${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Performance Test
echo ""
echo -e "${BLUE}Running Performance Test...${NC}"
START_TIME=$(date +%s%3N)

for i in {1..5}; do
    curl -s -X POST "http://localhost:5001/demo-blytz-mvp/us-central1/health" \
        -H "Content-Type: application/json" \
        -d '{"data": {}}' > /dev/null
done

END_TIME=$(date +%s%3N)
TOTAL_TIME=$((END_TIME - START_TIME))
AVG_TIME=$((TOTAL_TIME / 5))

echo -e "${GREEN}✅ Performance Test Completed${NC}"
echo "Average response time: ${AVG_TIME}ms"

if [ $AVG_TIME -lt 500 ]; then
    echo -e "${GREEN}🚀 Excellent performance (< 500ms)${NC}"
elif [ $AVG_TIME -lt 1000 ]; then
    echo -e "${YELLOW}⚡ Good performance (< 1000ms)${NC}"
else
    echo -e "${RED}⚠️  Performance could be improved (> 1000ms)${NC}"
fi

# Final Summary
echo ""
echo "=========================================="
echo "🎯 FINAL TEST RESULTS"
echo "=========================================="
echo -e "${GREEN}✅ Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}❌ Tests Failed: $TESTS_FAILED${NC}"
echo -e "${BLUE}📊 Total Tests: $((TESTS_PASSED + TESTS_FAILED))${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! Firebase integration is working perfectly.${NC}"
    echo ""
    echo "🚀 READY FOR GO SERVICE INTEGRATION:"
    echo "1. Import: firebase.NewClient(logger)"
    echo "2. Use: client.CreateAuction(), client.PlaceBid(), etc."
    echo "3. Test with your Go microservices"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed. Check the responses above.${NC}"
    exit 1
fi
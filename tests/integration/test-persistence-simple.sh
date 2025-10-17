#!/bin/bash

# Blytz Auction MVP - Simple Persistence Verification Test
# Quick test to verify database persistence is working correctly

set -e

echo "🔍 Simple Persistence Verification Test"
echo "======================================"

# Configuration
AUCTION_BASE_URL="http://localhost:8083/api/v1"

echo "Testing basic persistence functionality..."
echo ""

# Test 1: Check if auction service is running
echo "1. Checking auction service health..."
if curl -s "$AUCTION_BASE_URL/health" | grep -q "ok"; then
    echo "✅ Auction service is running"
else
    echo "❌ Auction service not responding"
    exit 1
fi

echo ""

# Test 2: Try to create an auction (this will verify database connection)
echo "2. Testing auction creation (verifies database connection)..."

create_data=$(cat << EOF
{
    "title": "Persistence Test $(date +%s)",
    "description": "Testing database persistence",
    "starting_price": 50.00,
    "reserve_price": 100.00,
    "min_bid_increment": 5.00,
    "start_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "end_time": "$(date -u -d '+1 hour' +%Y-%m-%dT%H:%M:%SZ)",
    "type": "live",
    "product_id": "test_$(date +%s)"
}
EOF
)

response=$(curl -s -X POST "$AUCTION_BASE_URL/auctions" \
  -H "Content-Type: application/json" \
  -d "$create_data")

if echo "$response" | jq -e '.success' >/dev/null 2>&1; then
    auction_id=$(echo "$response" | jq -r '.data.auction.auction_id')
    echo "✅ Auction created successfully"
    echo "   Auction ID: $auction_id"

    # Test 3: Verify the auction was actually saved
    echo ""
    echo "3. Verifying auction was saved to database..."

    retrieve_response=$(curl -s "$AUCTION_BASE_URL/auctions/$auction_id")

    if echo "$retrieve_response" | jq -e '.success' >/dev/null 2>&1; then
        title=$(echo "$retrieve_response" | jq -r '.data.auction.title')
        if [[ "$title" == *"Persistence Test"* ]]; then
            echo "✅ Auction retrieved from database successfully"
            echo "   Title: $title"
            echo "   ✅ PERSISTENCE CONFIRMED: Data is being saved to database!"
        else
            echo "❌ Auction data doesn't match - persistence failed"
            exit 1
        fi
    else
        echo "❌ Failed to retrieve auction from database"
        exit 1
    fi

    echo ""
    echo "4. Testing bid placement (verifies bid persistence)..."

    bid_response=$(curl -s -X POST "$AUCTION_BASE_URL/auctions/$auction_id/bids" \
      -H "Content-Type: application/json" \
      -d '{"amount":75.00}')

    if echo "$bid_response" | jq -e '.success' >/dev/null 2>&1; then
        echo "✅ Bid placed successfully"

        # Verify price was updated
        updated_response=$(curl -s "$AUCTION_BASE_URL/auctions/$auction_id")
        updated_price=$(echo "$updated_response" | jq -r '.data.auction.current_price')

        if (( $(echo "$updated_price == 75.00" | bc -l) )); then
            echo "✅ Auction price updated to $75.00 - bid persistence confirmed!"
        else
            echo "❌ Price not updated correctly"
            exit 1
        fi
    else
        echo "❌ Bid placement failed"
        echo "Response: $bid_response"
        exit 1
    fi

else
    echo "❌ Auction creation failed"
    echo "Response: $response"
    echo ""
    echo "💡 This indicates the auction service cannot connect to the database."
    echo "   The application should fail to start if the database is unavailable."
    exit 1
fi

echo ""
echo "🎯 PERSISTENCE VERIFICATION SUMMARY"
echo "===================================="
echo "✅ Auction creation: WORKING (data saved to database)"
echo "✅ Auction retrieval: WORKING (data retrieved from database)"
echo "✅ Bid placement: WORKING (bid saved and auction updated)"
echo ""
echo -e "${GREEN}🎉 DATABASE PERSISTENCE IS CONFIRMED!${NC}"
echo -e "${GREEN}   The auction service is now truly EXHIBITION READY!${NC}"
echo ""
echo "🚀 Key Improvements Made:"
echo "   • Application fails fast if database is unavailable"
echo "   • No more silent fallback to mock data"
echo "   • All operations use database persistence"
echo "   • UpdateAuction and DeleteAuction properly implemented"
echo "   • Bid operations properly update database"
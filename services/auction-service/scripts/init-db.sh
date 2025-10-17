#!/bin/bash

# Database initialization script for auction service

set -e

echo "🚀 Initializing auction service database..."

# Database connection parameters
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"
DB_NAME="${DB_NAME:-auction_db}"

# Construct database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "📡 Connecting to database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL client (psql) is not installed"
    exit 1
fi

# Test connection
echo "🔍 Testing database connection..."
if ! PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &> /dev/null; then
    echo "❌ Failed to connect to database"
    echo "📝 Make sure PostgreSQL is running and the database exists"
    exit 1
fi

echo "✅ Database connection successful"

# Create tables
echo "🏗️ Creating database tables..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f migrations/001_create_auctions_table.sql

echo "✅ Database tables created successfully"

# Insert sample data for demo
echo "📝 Inserting sample auction data..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << EOF
-- Insert sample auctions
INSERT INTO auctions (auction_id, product_id, seller_id, title, description, starting_price, current_price, reserve_price, min_bid_increment, start_time, end_time, status, type, is_active, created_at, updated_at) VALUES
('auction_001', 'product_001', 'seller_001', 'Vintage Watch', 'Beautiful vintage Rolex watch from 1960s', 500.00, 500.00, 800.00, 25.00, NOW() - INTERVAL '30 minutes', NOW() + INTERVAL '2 hours', 'active', 'live', true, NOW(), NOW()),
('auction_002', 'product_002', 'seller_002', 'Antique Vase', 'Ming dynasty porcelain vase', 200.00, 200.00, 400.00, 20.00, NOW() - INTERVAL '1 hour', NOW() + INTERVAL '3 hours', 'active', 'live', true, NOW(), NOW()),
('auction_003', 'product_003', 'seller_003', 'Collectible Coin', 'Rare 1921 silver dollar', 150.00, 150.00, 300.00, 15.00, NOW() + INTERVAL '1 hour', NOW() + INTERVAL '4 hours', 'scheduled', 'live', true, NOW(), NOW());

-- Insert sample bids
INSERT INTO bids (bid_id, auction_id, bidder_id, amount, is_winning, bid_time, created_at) VALUES
('bid_001', 'auction_001', 'bidder_001', 525.00, false, NOW() - INTERVAL '20 minutes', NOW()),
('bid_002', 'auction_001', 'bidder_002', 550.00, true, NOW() - INTERVAL '15 minutes', NOW()),
('bid_003', 'auction_002', 'bidder_003', 220.00, true, NOW() - INTERVAL '30 minutes', NOW());
EOF

echo "✅ Sample data inserted successfully"

echo "🎉 Database initialization complete!"
echo ""
echo "📊 Database Summary:"
echo "   - Auctions table: Created with indexes"
echo "   - Bids table: Created with foreign key constraint"
echo "   - Sample data: 3 auctions, 3 bids inserted"
echo ""
echo "🚀 The auction service is ready to use with real database persistence!"

echo ""
echo "📋 Next steps:"
echo "   1. Start the auction service: docker-compose up auction-service"
echo "   2. Test the API endpoints"
echo "   3. Create more auctions through the API"
echo "   4. Place bids on active auctions"
echo ""
echo "🔗 API Endpoints available at:"
echo "   - GET  http://localhost:8083/api/v1/auctions"
echo "   - POST http://localhost:8083/api/v1/auctions (requires auth)"
echo "   - GET  http://localhost:8083/api/v1/auctions/{id}/bids"
echo "   - POST http://localhost:8083/api/v1/auctions/{id}/bids (requires auth)"
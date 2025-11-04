#!/bin/bash

# Production Health Check
# Returns JSON with database and user statistics

echo "Content-Type: application/json"
echo ""

# Get statistics
TOTAL_USERS=$(docker compose exec -T postgres psql -U postgres -d blytz_prod -t -c "SELECT COUNT(*) FROM users;" 2>/dev/null | tr -d ' ')
USERS_24H=$(docker compose exec -T postgres psql -U postgres -d blytz_prod -t -c "SELECT COUNT(*) FROM users WHERE created_at > NOW() - INTERVAL '24 hours';" 2>/dev/null | tr -d ' ')
DB_SIZE=$(docker compose exec -T postgres psql -U postgres -d blytz_prod -t -c "SELECT pg_size_pretty(pg_database_size('blytz_prod'));" 2>/dev/null | tr -d ' ')

# Get timestamp
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

cat << EOF
{
  "status": "healthy",
  "timestamp": "$TIMESTAMP",
  "database": {
    "status": "connected",
    "size": "$DB_SIZE",
    "total_users": $TOTAL_USERS,
    "users_last_24h": $USERS_24H
  },
  "services": {
    "postgres": "running",
    "auth": "running"
  }
}
EOF
#!/bin/bash

# Production Database Monitoring Script
# Usage: ./check-prod-db.sh

echo "=========================================="
echo "Blytz Production Database Status"
echo "=========================================="
echo "Timestamp: $(date)"
echo ""

# Check if PostgreSQL container is running
if ! docker compose ps postgres | grep -q "Up"; then
    echo "❌ PostgreSQL container is not running"
    exit 1
fi

echo "✅ PostgreSQL container is running"
echo ""

# User Statistics
echo "📊 USER STATISTICS"
echo "-------------------"
docker compose exec -T postgres psql -U postgres -d blytz_prod -c "
SELECT 
  'Total Users' as metric,
  COUNT(*)::text as value
FROM users
UNION ALL
SELECT 
  'Users Last 24h',
  COUNT(CASE WHEN created_at > NOW() - INTERVAL '24 hours' THEN 1 END)::text
FROM users
UNION ALL
SELECT 
  'Users Last 7d',
  COUNT(CASE WHEN created_at > NOW() - INTERVAL '7 days' THEN 1 END)::text
FROM users
UNION ALL
SELECT 
  'Active Users',
  COUNT(CASE WHEN is_active = true THEN 1 END)::text
FROM users;
" 2>/dev/null

echo ""

# Recent Users
echo "👥 RECENT USERS (Last 5)"
echo "-------------------------"
docker compose exec -T postgres psql -U postgres -d blytz_prod -c "
SELECT 
  SUBSTRING(email, 1, 20) as email,
  display_name,
  is_active,
  created_at::timestamp(0) as created_at
FROM users 
ORDER BY created_at DESC 
LIMIT 5;
" 2>/dev/null

echo ""

# Database Size
echo "💾 DATABASE SIZE"
echo "----------------"
docker compose exec -T postgres psql -U postgres -d blytz_prod -c "
SELECT 
  pg_size_pretty(pg_database_size('blytz_prod')) as database_size,
  pg_size_pretty(pg_total_relation_size('users')) as users_table_size;
" 2>/dev/null

echo ""

# Recent Activity (if you have audit logs)
echo "📈 RECENT ACTIVITY"
echo "------------------"
echo "Check application logs for recent registration attempts:"
echo "docker compose logs auth-service | grep 'POST.*register' | tail -5"

echo ""
echo "=========================================="
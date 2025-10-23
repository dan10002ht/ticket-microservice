#!/bin/bash
# Check PgPool-II Configuration
# This script checks if all services are configured to use PgPool-II

set -e

echo "🔍 Checking PgPool-II Configuration..."

# Check Auth Service Config
echo "📊 Checking Auth Service..."
if grep -q "pgpool-auth" auth-service/src/config/databaseConfig.js; then
    echo "  ✅ Auth Service configured for PgPool-II"
else
    echo "  ❌ Auth Service NOT configured for PgPool-II"
fi

# Check Event Service Config
echo "📊 Checking Event Service..."
if grep -q "pgpool-event" event-service/config/config.go; then
    echo "  ✅ Event Service configured for PgPool-II"
else
    echo "  ❌ Event Service NOT configured for PgPool-II"
fi

# Check Ticket Service Config
echo "📊 Checking Ticket Service..."
if grep -q "pgpool-ticket" ticket-service/config/config.go; then
    echo "  ✅ Ticket Service configured for PgPool-II"
else
    echo "  ❌ Ticket Service NOT configured for PgPool-II"
fi

# Check Docker Compose
echo "📊 Checking Docker Compose..."
if grep -q "PGPOOL_AUTH_HOST" deploy/docker-compose.dev.yml; then
    echo "  ✅ Docker Compose configured for PgPool-II"
else
    echo "  ❌ Docker Compose NOT configured for PgPool-II"
fi

# Check PgPool-II Infrastructure
echo "📊 Checking PgPool-II Infrastructure..."
if [ -f "deploy/pgpool/docker-compose.pgpool.yml" ]; then
    echo "  ✅ PgPool-II Infrastructure files exist"
else
    echo "  ❌ PgPool-II Infrastructure files missing"
fi

echo ""
echo "📋 Configuration Summary:"
echo "  - Auth Service:     pgpool-auth:5432"
echo "  - Event Service:    pgpool-event:5432"  
echo "  - Ticket Service:   pgpool-ticket:5432"
echo ""
echo "🚀 Next Steps:"
echo "  1. Start PgPool-II:  ./scripts/start-dev-with-pgpool.sh"
echo "  2. Start Services:   ./scripts/dev-all.sh"
echo "  3. Check Health:      ./deploy/pgpool/scripts/check-health.sh"
echo "  4. Stop All:          ./scripts/stop-dev-with-pgpool.sh"
echo ""
echo "✅ Configuration check completed!"

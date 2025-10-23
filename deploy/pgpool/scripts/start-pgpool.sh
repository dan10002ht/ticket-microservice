#!/bin/bash
# Start PgPool-II Infrastructure
# This script starts PgPool-II and PostgreSQL clusters

set -e

echo "🚀 Starting PgPool-II Infrastructure..."

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p logs/pgpool-auth
mkdir -p logs/pgpool-event
mkdir -p logs/pgpool-ticket

# Start PgPool-II infrastructure
echo "🐳 Starting Docker containers..."
docker-compose -f docker-compose.pgpool.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check PgPool-II health
echo "🔍 Checking PgPool-II health..."

# Check Auth PgPool-II
echo "  - Checking pgpool-auth..."
if docker exec pgpool-auth pgpool -c "show pool_status" > /dev/null 2>&1; then
    echo "    ✅ pgpool-auth is healthy"
else
    echo "    ❌ pgpool-auth is not responding"
fi

# Check Event PgPool-II
echo "  - Checking pgpool-event..."
if docker exec pgpool-event pgpool -c "show pool_status" > /dev/null 2>&1; then
    echo "    ✅ pgpool-event is healthy"
else
    echo "    ❌ pgpool-event is not responding"
fi

# Check Ticket PgPool-II
echo "  - Checking pgpool-ticket..."
if docker exec pgpool-ticket pgpool -c "show pool_status" > /dev/null 2>&1; then
    echo "    ✅ pgpool-ticket is healthy"
else
    echo "    ❌ pgpool-ticket is not responding"
fi

# Show connection info
echo ""
echo "📊 PgPool-II Connection Information:"
echo "  - Auth Database:    localhost:5432"
echo "  - Event Database:  localhost:5433"
echo "  - Ticket Database: localhost:5434"
echo ""
echo "🔧 PgPool-II Management:"
echo "  - Auth PCP:        localhost:9898"
echo "  - Event PCP:       localhost:9899"
echo "  - Ticket PCP:      localhost:9900"
echo ""
echo "✅ PgPool-II Infrastructure started successfully!"
echo ""
echo "💡 Next steps:"
echo "   1. Update your service configs to use PgPool-II endpoints"
echo "   2. Start your application services"
echo "   3. Test the load balancing and failover"


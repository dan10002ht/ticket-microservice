#!/bin/bash
# Check PgPool-II Health
# This script checks the health of all PgPool-II instances

set -e

echo "🔍 Checking PgPool-II Health..."

# Function to check PgPool-II status
check_pgpool() {
    local name=$1
    local port=$2
    local pcp_port=$3
    
    echo "  - Checking $name (port $port)..."
    
    # Check if container is running
    if ! docker ps | grep -q "$name"; then
        echo "    ❌ Container $name is not running"
        return 1
    fi
    
    # Check PgPool-II status
    if docker exec "$name" pgpool -c "show pool_status" > /dev/null 2>&1; then
        echo "    ✅ $name is healthy"
        
        # Show pool status
        echo "    📊 Pool Status:"
        docker exec "$name" pgpool -c "show pool_status" | head -10
        
        # Show backend status
        echo "    🗄️  Backend Status:"
        docker exec "$name" pgpool -c "show pool_nodes" | head -10
        
    else
        echo "    ❌ $name is not responding"
        return 1
    fi
}

# Check all PgPool-II instances
check_pgpool "pgpool-auth" "5432" "9898"
check_pgpool "pgpool-event" "5433" "9899"
check_pgpool "pgpool-ticket" "5434" "9900"

echo ""
echo "📊 Connection Test:"
echo "  - Auth Database:    psql -h localhost -p 5432 -U postgres -d booking_system_auth"
echo "  - Event Database:   psql -h localhost -p 5433 -U postgres -d booking_system_event"
echo "  - Ticket Database:  psql -h localhost -p 5434 -U postgres -d booking_system_ticket"
echo ""
echo "✅ Health check completed!"


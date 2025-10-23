#!/bin/bash
# Start Development Environment with PgPool-II
# This script starts PgPool-II infrastructure and all services

echo "🚀 Starting Development Environment with PgPool-II"

# Function to check if port is in use
kill_port() {
    local port=$1
    local service_name=$2
    
    echo "🔍 Checking if port $port is in use by $service_name..."
    
    local pids=$(ss -tlnp | grep ":$port " | awk '{print $7}' | sed 's/.*pid=\([0-9]*\).*/\1/' | sort -u)
    
    if [ ! -z "$pids" ]; then
        echo "⚠️  Found processes using port $port: $pids, killing them..."
        echo $pids | xargs kill -9 2>/dev/null
        sleep 3
    else
        echo "✅ Port $port is available"
    fi
}

# Check prerequisites
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker"
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose v2 is not installed. Please install Docker Desktop or docker compose plugin."
    exit 1
fi

echo "✅ Prerequisites check passed"

# Kill existing processes
echo "🧹 Cleaning up existing processes..."
kill_port 5432 "pgpool-auth"
kill_port 5433 "pgpool-event"
kill_port 5434 "pgpool-ticket"

# Start PgPool-II infrastructure
echo "🐳 Starting PgPool-II infrastructure..."
cd deploy/pgpool

if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    echo "✅ Using Docker Compose v2 for PgPool-II"
    docker compose -f docker-compose.pgpool.yml up -d
else
    echo "❌ Docker Compose v2 not available, trying docker-compose..."
    docker-compose -f docker-compose.pgpool.yml up -d
fi

# Wait for PgPool-II to be ready
echo "⏳ Waiting for PgPool-II to be ready..."
sleep 20

# Check PgPool-II health
echo "🔍 Checking PgPool-II health..."
./scripts/check-health.sh

# Go back to root directory
cd ../..

echo ""
echo "🎉 PgPool-II infrastructure started successfully!"
echo ""
echo "📊 Available PgPool-II endpoints:"
echo "   - Auth Service:    localhost:5432 (pgpool-auth)"
echo "   - Event Service:   localhost:5433 (pgpool-event)"
echo "   - Ticket Service:  localhost:5434 (pgpool-ticket)"
echo ""
echo "🚀 Next steps:"
echo "   1. Start services: ./scripts/dev-all.sh"
echo "   2. Or start individual services:"
echo "      - Auth Service:  cd auth-service && yarn dev:local"
echo "      - Event Service: cd event-service && go run main.go"
echo "      - Ticket Service: cd ticket-service && go run main.go"
echo ""
echo "💡 Tips:"
echo "   - All services will now connect through PgPool-II"
echo "   - PgPool-II handles master/slave routing automatically"
echo "   - Use Ctrl+C to stop PgPool-II infrastructure"
echo ""

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping PgPool-II infrastructure..."
    cd deploy/pgpool
    docker compose -f docker-compose.pgpool.yml down
    echo "✅ PgPool-II infrastructure stopped"
    exit 0
}

# Set trap to cleanup on exit
trap cleanup SIGINT SIGTERM

# Wait for user to stop
echo "Press Ctrl+C to stop PgPool-II infrastructure..."
wait
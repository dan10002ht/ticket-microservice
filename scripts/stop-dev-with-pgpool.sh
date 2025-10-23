#!/bin/bash
# Stop Development Environment with PgPool-II
# This script stops PgPool-II infrastructure and all services

echo "🛑 Stopping Development Environment with PgPool-II"

# Stop PgPool-II infrastructure
echo "🐳 Stopping PgPool-II infrastructure..."
cd deploy/pgpool

if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    echo "✅ Using Docker Compose v2 for PgPool-II"
    docker compose -f docker-compose.pgpool.yml down
else
    echo "❌ Docker Compose v2 not available, trying docker-compose..."
    docker-compose -f docker-compose.pgpool.yml down
fi

# Go back to root directory
cd ../..

# Stop other infrastructure services
echo "🐳 Stopping other infrastructure services..."
cd deploy

if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    echo "✅ Using Docker Compose v2"
    docker compose -f docker-compose.dev.yml down
else
    echo "❌ Docker Compose v2 not available, trying docker-compose..."
    docker-compose -f docker-compose.dev.yml down
fi

# Go back to root directory
cd ..

echo ""
echo "✅ All development services stopped!"
echo ""
echo "💡 Tips:"
echo "   - All PgPool-II infrastructure has been stopped"
echo "   - All application services have been stopped"
echo "   - To start again: ./scripts/start-dev-with-pgpool.sh"
echo ""
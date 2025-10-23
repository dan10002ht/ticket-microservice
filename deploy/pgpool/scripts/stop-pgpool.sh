#!/bin/bash
# Stop PgPool-II Infrastructure
# This script stops PgPool-II and PostgreSQL clusters

set -e

echo "🛑 Stopping PgPool-II Infrastructure..."

# Stop PgPool-II infrastructure
echo "🐳 Stopping Docker containers..."
docker-compose -f docker-compose.pgpool.yml down

# Clean up volumes (optional)
read -p "🗑️  Do you want to remove volumes? This will delete all data! (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️  Removing volumes..."
    docker-compose -f docker-compose.pgpool.yml down -v
    echo "✅ Volumes removed"
else
    echo "💾 Volumes preserved"
fi

echo "✅ PgPool-II Infrastructure stopped!"


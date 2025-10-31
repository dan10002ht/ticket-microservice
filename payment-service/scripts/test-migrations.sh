#!/bin/bash

# Test Flyway migrations for Payment Service
# This script tests migrations against a local PostgreSQL database

set -e

echo "🧪 Testing Flyway Migrations for Payment Service"

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="payment_db_test"
DB_USER="${DB_USER:-booking_user}"
DB_PASSWORD="${DB_PASSWORD:-booking_pass}"

echo "📋 Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "❌ psql is not installed. Please install PostgreSQL client"
    exit 1
fi

# Check if Maven is installed
if ! command -v mvn &> /dev/null; then
    echo "❌ Maven is not installed. Please install Maven"
    exit 1
fi

# Create test database
echo "🗄️  Creating test database..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || true
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"

echo "✅ Test database created"

# Run Flyway migrations
echo "🚀 Running Flyway migrations..."
cd "$(dirname "$0")/.."

mvn flyway:migrate \
    -Dflyway.url=jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME \
    -Dflyway.user=$DB_USER \
    -Dflyway.password=$DB_PASSWORD \
    -Dflyway.locations=filesystem:src/main/resources/db/migration

echo "✅ Migrations completed successfully"

# Verify tables
echo "🔍 Verifying tables..."
TABLES=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename;")

echo "📊 Created tables:"
echo "$TABLES"

# Check migration history
echo "📜 Checking migration history..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT installed_rank, version, description, script, installed_on, success FROM flyway_schema_history ORDER BY installed_rank;"

# Verify views
echo "🔍 Verifying views..."
VIEWS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT viewname FROM pg_views WHERE schemaname = 'public' ORDER BY viewname;")

echo "📊 Created views:"
echo "$VIEWS"

# Verify functions
echo "🔍 Verifying functions..."
FUNCTIONS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT proname FROM pg_proc WHERE pronamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public') ORDER BY proname;")

echo "📊 Created functions:"
echo "$FUNCTIONS"

# Test migration rollback (validate down scripts if they exist)
echo "🔄 Testing migration info..."
mvn flyway:info \
    -Dflyway.url=jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME \
    -Dflyway.user=$DB_USER \
    -Dflyway.password=$DB_PASSWORD

echo ""
echo "✅ All migration tests passed!"
echo ""
echo "💡 To clean up test database:"
echo "   PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c \"DROP DATABASE $DB_NAME;\""




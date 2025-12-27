#!/bin/bash

# Database initialization script for Ticket Service

set -e

echo "🚀 Initializing Ticket Service Database..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5433}
DB_NAME=${DB_NAME:-booking_system}
DB_USER=${DB_USER:-booking_user}
DB_PASSWORD=${DB_PASSWORD:-booking_pass}

echo "📊 Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"

# Export password for psql
export PGPASSWORD=$DB_PASSWORD

# Check if database exists
echo "🔍 Checking if database exists..."
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo "✅ Database '$DB_NAME' already exists"
else
    echo "📝 Creating database '$DB_NAME'..."
    createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    echo "✅ Database '$DB_NAME' created successfully"
fi

# Create database connection string
DSN="host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable"

# Run migrations
echo "🔄 Running database migrations..."
if [ -d "./migrations" ]; then
    # If migrations directory exists, run migrations
    go run main.go migrate
else
    echo "⚠️  Migrations directory not found, skipping migrations"
fi

# Create utility functions
echo "🔧 Creating utility functions..."
psql "$DSN" -c "$(cat << 'EOF'
-- Create update_updated_at_column function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
EOF
)"

echo "✅ Database initialization completed successfully!"
echo ""
echo "🎯 Next steps:"
echo "  1. Run 'go mod tidy' to install dependencies"
echo "  2. Run 'go run main.go' to start the service"
echo "  3. Check logs for any errors"
echo ""
echo "📊 Database connection details:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo ""
echo "🔗 Connection string: $DSN"

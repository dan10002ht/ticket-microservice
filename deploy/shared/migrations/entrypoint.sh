#!/bin/sh
set -e

# Migration entrypoint for Node.js services (Knex)
# Usage: SERVICE_NAME=auth-service DB_HOST=postgres DB_NAME=mydb ...

echo "============================================"
echo "  Migration Runner - $SERVICE_NAME"
echo "============================================"

# Required environment variables
: "${SERVICE_NAME:?SERVICE_NAME is required}"
: "${DB_HOST:?DB_HOST is required}"
: "${DB_PORT:=5432}"
: "${DB_NAME:?DB_NAME is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"

echo "Waiting for PostgreSQL at $DB_HOST:$DB_PORT..."

# Wait for database to be ready
max_attempts=30
attempt=0
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo "ERROR: Database not ready after $max_attempts attempts"
        exit 1
    fi
    echo "  Attempt $attempt/$max_attempts - waiting..."
    sleep 2
done

echo "PostgreSQL is ready!"

# Check if database exists, create if not
echo "Checking database $DB_NAME..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || {
    echo "Creating database $DB_NAME..."
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME"
}

# Change to service directory
cd /app/$SERVICE_NAME

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install --production=false
fi

# Set database connection for Knex
export DB_MASTER_HOST=$DB_HOST
export DB_MASTER_PORT=$DB_PORT
export DB_MASTER_NAME=$DB_NAME
export DB_MASTER_USER=$DB_USER
export DB_MASTER_PASSWORD=$DB_PASSWORD

# Run migrations
echo "Running migrations..."
npx knex migrate:latest --env development

# Run seeds if SEED=true
if [ "$SEED" = "true" ]; then
    echo "Running seeds..."
    npx knex seed:run --env development
fi

echo "============================================"
echo "  Migration completed for $SERVICE_NAME"
echo "============================================"

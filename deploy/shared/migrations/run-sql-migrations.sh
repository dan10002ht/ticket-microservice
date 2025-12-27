#!/bin/sh
set -e

# Migration runner for raw SQL files
# Usage: SERVICE_NAME=user-service DB_HOST=postgres DB_NAME=mydb ...

echo "============================================"
echo "  SQL Migration Runner - $SERVICE_NAME"
echo "============================================"

# Required environment variables
: "${SERVICE_NAME:?SERVICE_NAME is required}"
: "${DB_HOST:?DB_HOST is required}"
: "${DB_PORT:=5432}"
: "${DB_NAME:?DB_NAME is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"
: "${MIGRATIONS_PATH:=/migrations}"

echo "Waiting for PostgreSQL at $DB_HOST:$DB_PORT..."

# Wait for database to be ready
max_attempts=30
attempt=0
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; do
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

# Create migrations tracking table
echo "Creating migrations tracking table..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) UNIQUE NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF

# Run SQL migrations in order
echo "Running SQL migrations from $MIGRATIONS_PATH..."

# Find and sort migration files
for migration_file in $(find "$MIGRATIONS_PATH" -name "*.up.sql" -o -name "*.sql" ! -name "*.down.sql" | sort); do
    filename=$(basename "$migration_file")

    # Check if already applied
    applied=$(PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT 1 FROM schema_migrations WHERE filename='$filename'" 2>/dev/null | tr -d ' ')

    if [ "$applied" = "1" ]; then
        echo "  SKIP: $filename (already applied)"
    else
        echo "  APPLY: $filename"
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"

        # Record migration
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "INSERT INTO schema_migrations (filename) VALUES ('$filename')"
    fi
done

echo "============================================"
echo "  SQL Migration completed for $SERVICE_NAME"
echo "============================================"

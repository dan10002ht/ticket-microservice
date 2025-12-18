#!/bin/bash

# Setup databases for all services
# Usage: ./scripts/setup-databases.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${GREEN}[INFO]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_header() { echo -e "${BLUE}========================================${NC}"; echo -e "${BLUE}$1${NC}"; echo -e "${BLUE}========================================${NC}"; }

# Database configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-55435}"  # docker-compose port
DB_USER="${DB_USER:-booking_user}"
DB_PASSWORD="${DB_PASSWORD:-booking_pass}"

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

print_header "Database Setup for Ticket Microservices"
print_status "Host: $DB_HOST:$DB_PORT"
print_status "User: $DB_USER"

# Export password for psql
export PGPASSWORD="$DB_PASSWORD"

# Function to create database if not exists
create_db_if_not_exists() {
    local db_name=$1
    print_status "Checking database: $db_name"

    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$db_name"; then
        print_warning "Database '$db_name' already exists, skipping creation"
    else
        print_status "Creating database: $db_name"
        createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$db_name"
        print_status "Database '$db_name' created successfully"
    fi
}

# Function to run migrations
run_migration() {
    local db_name=$1
    local migration_file=$2

    if [ -f "$migration_file" ]; then
        print_status "Running migration for $db_name: $(basename "$migration_file")"
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$db_name" -f "$migration_file"
        print_status "Migration completed for $db_name"
    else
        print_warning "Migration file not found: $migration_file"
    fi
}

# ============================================
# Setup User Service Database
# ============================================
print_header "Setting up User Service Database"
create_db_if_not_exists "user_service"
run_migration "user_service" "$PROJECT_ROOT/user-service/migrations/001_init.sql"

# ============================================
# Setup Auth Service Database (if migration exists)
# ============================================
if [ -d "$PROJECT_ROOT/auth-service/migrations" ]; then
    print_header "Setting up Auth Service Database"
    create_db_if_not_exists "auth_service"
    for migration in "$PROJECT_ROOT/auth-service/migrations"/*.sql; do
        [ -f "$migration" ] && run_migration "auth_service" "$migration"
    done
fi

# ============================================
# Setup Booking System Database
# ============================================
print_header "Setting up Main Booking System Database"
create_db_if_not_exists "booking_system"

# ============================================
# Summary
# ============================================
print_header "Setup Complete"
print_status "Databases created/verified:"
print_status "  - user_service"
print_status "  - auth_service (if migrations exist)"
print_status "  - booking_system"

print_status ""
print_status "To verify, run:"
print_status "  psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d user_service -c '\\dt'"

unset PGPASSWORD

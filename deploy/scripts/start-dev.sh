#!/bin/bash
# Start Development Environment
# Usage: ./start-dev.sh [up|down|logs|ps|migrate]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_DIR="$SCRIPT_DIR/../environments/development"
COMPOSE_FILE="$ENV_DIR/docker-compose.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════╗"
    echo "║     Development Environment                      ║"
    echo "║     Single PostgreSQL + Auto Migration           ║"
    echo "╚══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_connection_info() {
    echo -e "${YELLOW}Connection Information:${NC}"
    echo "  Auth DB:    localhost:5432 (booking_system_auth)"
    echo "  Main DB:    localhost:5433 (booking_system)"
    echo "  Redis:      localhost:6379"
    echo "  Kafka:      localhost:9092"
    echo ""
    echo -e "${YELLOW}Database Credentials:${NC}"
    echo "  User:       booking_user"
    echo "  Password:   booking_pass"
}

wait_for_migrations() {
    echo -e "${CYAN}Waiting for migrations to complete...${NC}"

    # Wait for migration containers to finish
    docker compose -f "$COMPOSE_FILE" wait migrate-auth 2>/dev/null || true
    docker compose -f "$COMPOSE_FILE" wait migrate-user 2>/dev/null || true
    docker compose -f "$COMPOSE_FILE" wait migrate-event 2>/dev/null || true
    docker compose -f "$COMPOSE_FILE" wait migrate-ticket 2>/dev/null || true

    echo -e "${GREEN}All migrations completed!${NC}"
}

show_migration_logs() {
    echo -e "${CYAN}Migration Logs:${NC}"
    echo ""
    echo -e "${YELLOW}=== Auth Service ===${NC}"
    docker compose -f "$COMPOSE_FILE" logs migrate-auth 2>/dev/null | tail -10
    echo ""
    echo -e "${YELLOW}=== User Service ===${NC}"
    docker compose -f "$COMPOSE_FILE" logs migrate-user 2>/dev/null | tail -10
    echo ""
    echo -e "${YELLOW}=== Event Service ===${NC}"
    docker compose -f "$COMPOSE_FILE" logs migrate-event 2>/dev/null | tail -10
    echo ""
    echo -e "${YELLOW}=== Ticket Service ===${NC}"
    docker compose -f "$COMPOSE_FILE" logs migrate-ticket 2>/dev/null | tail -10
}

case "${1:-up}" in
    up)
        print_header
        echo "Starting development environment..."

        # Start databases first
        echo -e "${CYAN}Starting databases...${NC}"
        docker compose -f "$COMPOSE_FILE" up -d postgres-auth postgres-main redis

        # Wait for databases to be healthy
        echo -e "${CYAN}Waiting for databases to be ready...${NC}"
        sleep 5

        # Run migrations
        echo -e "${CYAN}Running migrations...${NC}"
        docker compose -f "$COMPOSE_FILE" up migrate-auth migrate-user migrate-event migrate-ticket

        # Start remaining services
        echo -e "${CYAN}Starting remaining services...${NC}"
        docker compose -f "$COMPOSE_FILE" up -d zookeeper kafka

        echo ""
        echo -e "${GREEN}Development environment started successfully!${NC}"
        echo ""
        print_connection_info
        ;;
    up-quick)
        # Start without waiting for migrations
        print_header
        echo "Starting development environment (quick mode)..."
        docker compose -f "$COMPOSE_FILE" up -d
        echo ""
        print_connection_info
        ;;
    down)
        echo "Stopping development environment..."
        docker compose -f "$COMPOSE_FILE" down
        ;;
    logs)
        docker compose -f "$COMPOSE_FILE" logs -f ${@:2}
        ;;
    ps)
        docker compose -f "$COMPOSE_FILE" ps
        ;;
    restart)
        echo "Restarting development environment..."
        docker compose -f "$COMPOSE_FILE" restart ${@:2}
        ;;
    migrate)
        echo -e "${CYAN}Running migrations...${NC}"
        docker compose -f "$COMPOSE_FILE" up migrate-auth migrate-user migrate-event migrate-ticket
        show_migration_logs
        ;;
    migrate-logs)
        show_migration_logs
        ;;
    clean)
        echo -e "${RED}WARNING: This will delete all data!${NC}"
        read -p "Are you sure? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker compose -f "$COMPOSE_FILE" down -v
            echo "All containers and volumes removed."
        fi
        ;;
    *)
        echo "Usage: $0 {up|up-quick|down|logs|ps|restart|migrate|migrate-logs|clean}"
        echo ""
        echo "Commands:"
        echo "  up           Start all services with migrations"
        echo "  up-quick     Start all services without waiting for migrations"
        echo "  down         Stop all services"
        echo "  logs         View logs (optional: service name)"
        echo "  ps           List running services"
        echo "  restart      Restart services (optional: service name)"
        echo "  migrate      Run all migrations"
        echo "  migrate-logs Show migration logs"
        echo "  clean        Stop and remove all data (WARNING: destructive)"
        exit 1
        ;;
esac

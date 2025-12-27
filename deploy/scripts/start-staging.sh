#!/bin/bash
# Start Staging Environment
# Usage: ./start-staging.sh [up|down|logs|ps|status]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_DIR="$SCRIPT_DIR/../environments/staging"
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
    echo "║     Staging Environment                          ║"
    echo "║     HAProxy + PostgreSQL Primary/Replica         ║"
    echo "╚══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_connection_info() {
    echo -e "${YELLOW}Connection Information:${NC}"
    echo ""
    echo -e "${CYAN}Auth Database:${NC}"
    echo "  Write:      localhost:5432 (via HAProxy -> Primary)"
    echo "  Read:       localhost:5433 (via HAProxy -> Replicas)"
    echo ""
    echo -e "${CYAN}Main Database:${NC}"
    echo "  Write:      localhost:5434 (via HAProxy -> Primary)"
    echo "  Read:       localhost:5435 (via HAProxy -> Replicas)"
    echo ""
    echo -e "${CYAN}Direct Access (for debugging):${NC}"
    echo "  Auth Primary:   localhost:55432"
    echo "  Auth Replica:   localhost:55433"
    echo "  Main Primary:   localhost:55434"
    echo "  Main Replica:   localhost:55435"
    echo ""
    echo -e "${CYAN}Other Services:${NC}"
    echo "  HAProxy Stats:  http://localhost:8404/stats (admin/admin)"
    echo "  Redis:          localhost:6379"
    echo "  Kafka:          localhost:9092"
    echo ""
    echo -e "${YELLOW}Database Credentials:${NC}"
    echo "  User:       booking_user"
    echo "  Password:   booking_pass"
}

check_replication_status() {
    echo -e "${CYAN}Checking replication status...${NC}"
    echo ""

    # Check Auth cluster
    echo "Auth Cluster:"
    docker exec staging-auth-postgres-primary psql -U booking_user -d booking_system_auth -c \
        "SELECT client_addr, state, sent_lsn, replay_lsn, sync_state FROM pg_stat_replication;" 2>/dev/null || echo "  Primary not ready"

    echo ""

    # Check Main cluster
    echo "Main Cluster:"
    docker exec staging-main-postgres-primary psql -U booking_user -d booking_system -c \
        "SELECT client_addr, state, sent_lsn, replay_lsn, sync_state FROM pg_stat_replication;" 2>/dev/null || echo "  Primary not ready"
}

case "${1:-up}" in
    up)
        print_header
        echo "Starting staging environment..."
        docker compose -f "$COMPOSE_FILE" up -d
        echo ""
        echo "Waiting for services to be healthy..."
        sleep 10
        print_connection_info
        ;;
    down)
        echo "Stopping staging environment..."
        docker compose -f "$COMPOSE_FILE" down
        ;;
    logs)
        docker compose -f "$COMPOSE_FILE" logs -f ${@:2}
        ;;
    ps)
        docker compose -f "$COMPOSE_FILE" ps
        ;;
    restart)
        echo "Restarting staging environment..."
        docker compose -f "$COMPOSE_FILE" restart ${@:2}
        ;;
    status)
        print_header
        echo -e "${YELLOW}Service Status:${NC}"
        docker compose -f "$COMPOSE_FILE" ps
        echo ""
        check_replication_status
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
        echo "Usage: $0 {up|down|logs|ps|restart|status|clean}"
        echo ""
        echo "Commands:"
        echo "  up       Start all services"
        echo "  down     Stop all services"
        echo "  logs     View logs (optional: service name)"
        echo "  ps       List running services"
        echo "  restart  Restart services (optional: service name)"
        echo "  status   Check service and replication status"
        echo "  clean    Stop and remove all data (WARNING: destructive)"
        exit 1
        ;;
esac

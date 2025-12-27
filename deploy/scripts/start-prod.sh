#!/bin/bash
# Start Production Environment
# Usage: ./start-prod.sh [up|down|logs|ps|status]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_DIR="$SCRIPT_DIR/../environments/production"
COMPOSE_FILE="$ENV_DIR/docker-compose.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${MAGENTA}"
    echo "╔══════════════════════════════════════════════════╗"
    echo "║     Production Environment                       ║"
    echo "║     PgBouncer + HAProxy + PostgreSQL Cluster     ║"
    echo "╚══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

check_env_file() {
    if [ ! -f "$ENV_DIR/.env" ]; then
        echo -e "${RED}ERROR: .env file not found!${NC}"
        echo "Please copy .env.example to .env and configure:"
        echo "  cp $ENV_DIR/.env.example $ENV_DIR/.env"
        echo ""
        echo -e "${YELLOW}IMPORTANT: Change all default passwords before starting!${NC}"
        exit 1
    fi
}

print_connection_info() {
    echo -e "${YELLOW}Connection Information:${NC}"
    echo ""
    echo -e "${CYAN}Application Connection (via PgBouncer):${NC}"
    echo "  Auth DB (RW):   localhost:6432/booking_system_auth"
    echo "  Auth DB (RO):   localhost:6432/booking_system_auth_ro"
    echo "  Main DB (RW):   localhost:6432/booking_system"
    echo "  Main DB (RO):   localhost:6432/booking_system_ro"
    echo ""
    echo -e "${CYAN}Monitoring:${NC}"
    echo "  HAProxy Stats:  http://localhost:8404/stats"
    echo "  Prometheus:     http://localhost:9090"
    echo "  Grafana:        http://localhost:3001"
    echo ""
    echo -e "${CYAN}Other Services:${NC}"
    echo "  Redis:          localhost:6379"
    echo "  Kafka:          localhost:9092"
}

check_replication_status() {
    echo -e "${CYAN}Checking replication status...${NC}"
    echo ""

    # Check Auth cluster
    echo -e "${YELLOW}Auth Cluster:${NC}"
    docker exec prod-auth-postgres-primary psql -U booking_user -d booking_system_auth -c \
        "SELECT client_addr, application_name, state, sent_lsn, replay_lsn, sync_state FROM pg_stat_replication;" 2>/dev/null || echo "  Primary not ready"

    echo ""

    # Check Main cluster
    echo -e "${YELLOW}Main Cluster:${NC}"
    docker exec prod-main-postgres-primary psql -U booking_user -d booking_system -c \
        "SELECT client_addr, application_name, state, sent_lsn, replay_lsn, sync_state FROM pg_stat_replication;" 2>/dev/null || echo "  Primary not ready"
}

check_pgbouncer_status() {
    echo -e "${CYAN}PgBouncer Connection Pools:${NC}"
    docker exec prod-pgbouncer psql -p 6432 -U pgbouncer_admin pgbouncer -c "SHOW POOLS;" 2>/dev/null || echo "  PgBouncer not ready"
}

case "${1:-up}" in
    up)
        print_header
        check_env_file
        echo "Starting production environment..."
        docker compose -f "$COMPOSE_FILE" up -d
        echo ""
        echo "Waiting for services to be healthy..."
        sleep 15
        print_connection_info
        ;;
    down)
        echo "Stopping production environment..."
        docker compose -f "$COMPOSE_FILE" down
        ;;
    logs)
        docker compose -f "$COMPOSE_FILE" logs -f ${@:2}
        ;;
    ps)
        docker compose -f "$COMPOSE_FILE" ps
        ;;
    restart)
        echo "Restarting production environment..."
        docker compose -f "$COMPOSE_FILE" restart ${@:2}
        ;;
    status)
        print_header
        echo -e "${YELLOW}Service Status:${NC}"
        docker compose -f "$COMPOSE_FILE" ps
        echo ""
        check_replication_status
        echo ""
        check_pgbouncer_status
        ;;
    health)
        echo -e "${CYAN}Health Check:${NC}"
        echo ""
        echo "PostgreSQL Primaries:"
        docker exec prod-auth-postgres-primary pg_isready -U booking_user -d booking_system_auth && echo "  Auth Primary: OK" || echo "  Auth Primary: FAIL"
        docker exec prod-main-postgres-primary pg_isready -U booking_user -d booking_system && echo "  Main Primary: OK" || echo "  Main Primary: FAIL"
        echo ""
        echo "PostgreSQL Replicas:"
        docker exec prod-auth-postgres-replica1 pg_isready -U booking_user -d booking_system_auth && echo "  Auth Replica 1: OK" || echo "  Auth Replica 1: FAIL"
        docker exec prod-auth-postgres-replica2 pg_isready -U booking_user -d booking_system_auth && echo "  Auth Replica 2: OK" || echo "  Auth Replica 2: FAIL"
        docker exec prod-main-postgres-replica1 pg_isready -U booking_user -d booking_system && echo "  Main Replica 1: OK" || echo "  Main Replica 1: FAIL"
        docker exec prod-main-postgres-replica2 pg_isready -U booking_user -d booking_system && echo "  Main Replica 2: OK" || echo "  Main Replica 2: FAIL"
        echo ""
        echo "Redis:"
        docker exec prod-redis redis-cli ping && echo "  Redis: OK" || echo "  Redis: FAIL"
        ;;
    clean)
        echo -e "${RED}╔══════════════════════════════════════════════════╗${NC}"
        echo -e "${RED}║  WARNING: PRODUCTION DATA WILL BE DELETED!       ║${NC}"
        echo -e "${RED}╚══════════════════════════════════════════════════╝${NC}"
        echo ""
        read -p "Type 'DELETE PRODUCTION DATA' to confirm: " confirm
        if [ "$confirm" = "DELETE PRODUCTION DATA" ]; then
            docker compose -f "$COMPOSE_FILE" down -v
            echo "All containers and volumes removed."
        else
            echo "Aborted."
        fi
        ;;
    *)
        echo "Usage: $0 {up|down|logs|ps|restart|status|health|clean}"
        echo ""
        echo "Commands:"
        echo "  up       Start all services"
        echo "  down     Stop all services"
        echo "  logs     View logs (optional: service name)"
        echo "  ps       List running services"
        echo "  restart  Restart services (optional: service name)"
        echo "  status   Check service, replication and PgBouncer status"
        echo "  health   Quick health check of all databases"
        echo "  clean    Stop and remove all data (WARNING: destructive)"
        exit 1
        ;;
esac

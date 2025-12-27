#!/bin/bash

# Development script using native PostgreSQL and Redis (no Docker)
# Usage:
#   ./dev-native.sh                                    # Start all services
#   ./dev-native.sh --services auth,gateway,user       # Start specific services
#   ./dev-native.sh --help                             # Show help

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${GREEN}"
echo "╔══════════════════════════════════════════════════╗"
echo "║  Booking System - Native Development (No Docker) ║"
echo "╚══════════════════════════════════════════════════╝"
echo -e "${NC}"

# Parse command line arguments
SELECTED_SERVICES=""
START_ALL=true
SKIP_MIGRATION=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --services)
            SELECTED_SERVICES="$2"
            START_ALL=false
            shift 2
            ;;
        --skip-migration)
            SKIP_MIGRATION=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --services <svc1,svc2,...>  Start only specific services"
            echo "  --skip-migration            Skip database migrations"
            echo "  --help, -h                  Show this help message"
            echo ""
            echo "Available services:"
            echo "  - auth           Auth Service (Node.js, gRPC 50051)"
            echo "  - user           User Service (Go, gRPC 50052)"
            echo "  - event          Event Service (Go, gRPC 50053)"
            echo "  - booking        Booking Service (Go, gRPC 50054)"
            echo "  - payment        Payment Service (Java, gRPC 50056)"
            echo "  - realtime       Realtime Service (Go, HTTP 3003, gRPC 50057)"
            echo "  - ticket         Ticket Service (Go, gRPC 50058)"
            echo "  - booking-worker Booking Worker (Go, gRPC 50059)"
            echo "  - email-worker   Email Worker (Go, HTTP 8080, gRPC 50060)"
            echo "  - gateway        API Gateway (Node.js, HTTP 3000)"
            echo ""
            echo "Examples:"
            echo "  $0                                  # Start all"
            echo "  $0 --services auth,gateway          # Specific services"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Function to check if service should start
should_start_service() {
    local service=$1
    if [ "$START_ALL" = true ]; then
        return 0
    fi
    IFS=',' read -ra SERVICES <<< "$SELECTED_SERVICES"
    for s in "${SERVICES[@]}"; do
        if [ "$s" = "$service" ]; then
            return 0
        fi
    done
    return 1
}

# Function to kill process using a specific port
kill_port() {
    local port=$1
    local service_name=$2
    local pids=$(lsof -ti:$port 2>/dev/null)
    if [ ! -z "$pids" ]; then
        echo -e "${YELLOW}  Killing process on port $port ($service_name)${NC}"
        echo $pids | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# Function to kill all development processes
kill_all_dev_processes() {
    echo -e "${CYAN}Cleaning up existing processes...${NC}"
    pkill -f nodemon 2>/dev/null || true
    pkill -f "node.*src/index.js" 2>/dev/null || true
    pkill -f "go run" 2>/dev/null || true
    pkill -f "air" 2>/dev/null || true
    pkill -f "mvn spring-boot:run" 2>/dev/null || true
    sleep 2
}

# Check prerequisites
check_prerequisites() {
    echo -e "${CYAN}Checking prerequisites...${NC}"

    local missing=()

    if ! command -v node &> /dev/null; then
        missing+=("node")
    fi

    if ! command -v yarn &> /dev/null; then
        missing+=("yarn")
    fi

    if ! command -v psql &> /dev/null; then
        missing+=("postgresql")
    fi

    if ! command -v redis-cli &> /dev/null; then
        missing+=("redis")
    fi

    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Go not installed - Go services will be skipped${NC}"
    fi

    if ! command -v mvn &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Maven not installed - Java services will be skipped${NC}"
    fi

    if [ ${#missing[@]} -ne 0 ]; then
        echo -e "${RED}Missing required tools: ${missing[*]}${NC}"
        exit 1
    fi

    echo -e "${GREEN}  ✓ Prerequisites OK${NC}"
}

# Setup native infrastructure
setup_native_infrastructure() {
    echo ""
    echo -e "${CYAN}════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  Setting up Native Infrastructure${NC}"
    echo -e "${CYAN}════════════════════════════════════════════════${NC}"

    # Check PostgreSQL
    echo -e "${CYAN}Checking PostgreSQL...${NC}"
    if ! pg_isready -q 2>/dev/null; then
        echo -e "${YELLOW}  Starting PostgreSQL...${NC}"
        brew services start postgresql@14 2>/dev/null || brew services start postgresql 2>/dev/null || true
        sleep 3
    fi

    if pg_isready -q 2>/dev/null; then
        echo -e "${GREEN}  ✓ PostgreSQL is running${NC}"
    else
        echo -e "${RED}  ✗ PostgreSQL is not running. Please start it manually.${NC}"
        exit 1
    fi

    # Check Redis
    echo -e "${CYAN}Checking Redis...${NC}"
    if ! redis-cli ping > /dev/null 2>&1; then
        echo -e "${YELLOW}  Starting Redis...${NC}"
        brew services start redis 2>/dev/null || true
        sleep 2
    fi

    if redis-cli ping > /dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Redis is running${NC}"
    else
        echo -e "${RED}  ✗ Redis is not running. Please start it manually.${NC}"
        exit 1
    fi

    # Setup databases
    echo -e "${CYAN}Setting up databases...${NC}"

    # Get current macOS user for PostgreSQL connection
    CURRENT_USER=$(whoami)

    psql -U "$CURRENT_USER" -d postgres << 'EOF' 2>/dev/null || true
-- Create user for booking system (if not exists)
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'booking_user') THEN
      CREATE ROLE booking_user WITH LOGIN PASSWORD 'booking_pass';
   END IF;
END
$$;

-- Create databases (if not exist)
SELECT 'CREATE DATABASE booking_system_auth OWNER booking_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'booking_system_auth')\gexec

SELECT 'CREATE DATABASE booking_system OWNER booking_user'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'booking_system')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE booking_system_auth TO booking_user;
GRANT ALL PRIVILEGES ON DATABASE booking_system TO booking_user;

-- Grant schema permissions
\c booking_system_auth
GRANT ALL ON SCHEMA public TO booking_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO booking_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO booking_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO booking_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO booking_user;

\c booking_system
GRANT ALL ON SCHEMA public TO booking_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO booking_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO booking_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO booking_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO booking_user;
EOF

    echo -e "${GREEN}  ✓ Databases configured${NC}"

    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}  Native Infrastructure Ready!${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "  ${CYAN}Databases:${NC}"
    echo "    Auth DB:  localhost:5432/booking_system_auth"
    echo "    Main DB:  localhost:5432/booking_system"
    echo ""
    echo -e "  ${CYAN}Services:${NC}"
    echo "    Redis:    localhost:6379"
    echo ""
    echo -e "  ${CYAN}Credentials:${NC}"
    echo "    User:     booking_user"
    echo "    Password: booking_pass"
}

# Service PIDs tracking
AUTH_PID=""
USER_PID=""
EVENT_PID=""
BOOKING_PID=""
PAYMENT_PID=""
REALTIME_PID=""
TICKET_PID=""
BOOKING_WORKER_PID=""
EMAIL_WORKER_PID=""
GATEWAY_PID=""

# Start individual services
start_auth_service() {
    echo ""
    echo -e "${CYAN}Starting Auth Service...${NC}"
    kill_port 50051 "auth-service"
    cd "$PROJECT_DIR/auth-service"
    if [ ! -d "node_modules" ]; then
        echo "  Installing dependencies..."
        yarn install --silent
    fi

    # Create .env file for auth-service (native - all on port 5432)
    cat > .env << 'ENVEOF'
# Server Configuration
PORT=50051
HOST=0.0.0.0
NODE_ENV=development

# Database Configuration (Native PostgreSQL - port 5432)
DB_MASTER_HOST=localhost
DB_MASTER_PORT=5432
DB_MASTER_NAME=booking_system_auth
DB_MASTER_USER=booking_user
DB_MASTER_PASSWORD=booking_pass

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# JWT Configuration
JWT_SECRET=dev_jwt_secret_key_change_in_production
JWT_REFRESH_SECRET=dev_refresh_secret_key_change_in_production
JWT_EXPIRES_IN=15m
JWT_REFRESH_EXPIRES_IN=7d

# Development Configuration
DEBUG=true
SKIP_EMAIL_VERIFICATION=true
ENVEOF

    # Run migrations
    if [ "$SKIP_MIGRATION" = false ]; then
        echo "  Running migrations..."
        yarn migrate:latest 2>/dev/null || true
        yarn seed:run 2>/dev/null || true
    fi

    yarn dev &
    AUTH_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Auth Service started (gRPC: 50051)${NC}"
}

start_user_service() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping User Service (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting User Service...${NC}"
    kill_port 50052 "user-service"
    cd "$PROJECT_DIR/user-service"

    # Create .env for native
    cat > .env << 'ENVEOF'
GRPC_PORT=50052
DB_HOST=localhost
DB_PORT=5432
DB_NAME=booking_system
DB_USER=booking_user
DB_PASSWORD=booking_pass
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    USER_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ User Service started (gRPC: 50052)${NC}"
}

start_event_service() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Event Service (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Event Service...${NC}"
    kill_port 50053 "event-service"
    cd "$PROJECT_DIR/event-service"

    # Create .env for native
    cat > .env << 'ENVEOF'
GRPC_PORT=50053
DB_HOST=localhost
DB_PORT=5432
DB_NAME=booking_system
DB_USER=booking_user
DB_PASSWORD=booking_pass
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    EVENT_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Event Service started (gRPC: 50053)${NC}"
}

start_booking_service() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Booking Service (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Booking Service...${NC}"
    kill_port 50054 "booking-service"
    cd "$PROJECT_DIR/booking-service"

    # Create .env for native
    cat > .env << 'ENVEOF'
GRPC_PORT=50054
DB_HOST=localhost
DB_PORT=5432
DB_NAME=booking_system
DB_USER=booking_user
DB_PASSWORD=booking_pass
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    BOOKING_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Booking Service started (gRPC: 50054)${NC}"
}

start_payment_service() {
    if ! command -v mvn &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Payment Service (Maven not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Payment Service...${NC}"
    kill_port 8081 "payment-service-http"
    kill_port 50056 "payment-service-grpc"
    cd "$PROJECT_DIR/payment-service"
    mvn spring-boot:run -Dspring-boot.run.profiles=dev &
    PAYMENT_PID=$!
    cd "$PROJECT_DIR"
    sleep 10
    echo -e "${GREEN}  ✓ Payment Service started (HTTP: 8081, gRPC: 50056)${NC}"
}

start_realtime_service() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Realtime Service (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Realtime Service...${NC}"
    kill_port 3003 "realtime-service-http"
    kill_port 50057 "realtime-service-grpc"
    cd "$PROJECT_DIR/realtime-service"

    # Create .env for native
    cat > .env << 'ENVEOF'
HTTP_PORT=3003
GRPC_PORT=50057
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    REALTIME_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Realtime Service started (HTTP: 3003, gRPC: 50057)${NC}"
}

start_ticket_service() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Ticket Service (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Ticket Service...${NC}"
    kill_port 50058 "ticket-service"
    cd "$PROJECT_DIR/ticket-service"

    # Create .env for native (port 5432 for native PostgreSQL)
    cat > .env << 'ENVEOF'
GRPC_PORT=50058
DB_HOST=localhost
DB_PORT=5432
DB_NAME=booking_system
DB_USER=booking_user
DB_PASSWORD=booking_pass
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    TICKET_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Ticket Service started (gRPC: 50058)${NC}"
}

start_booking_worker() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Booking Worker (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Booking Worker...${NC}"
    kill_port 50059 "booking-worker"
    cd "$PROJECT_DIR/booking-worker"

    # Create .env for native
    cat > .env << 'ENVEOF'
GRPC_PORT=50059
DB_HOST=localhost
DB_PORT=5432
DB_NAME=booking_system
DB_USER=booking_user
DB_PASSWORD=booking_pass
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    BOOKING_WORKER_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Booking Worker started (gRPC: 50059)${NC}"
}

start_email_worker() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Email Worker (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Email Worker...${NC}"
    kill_port 8080 "email-worker-http"
    kill_port 50060 "email-worker-grpc"
    cd "$PROJECT_DIR/email-worker"

    # Create .env for native
    cat > .env << 'ENVEOF'
HTTP_PORT=8080
GRPC_PORT=50060
REDIS_HOST=localhost
REDIS_PORT=6379
ENV=development
ENVEOF

    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    EMAIL_WORKER_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Email Worker started (HTTP: 8080, gRPC: 50060)${NC}"
}

start_gateway() {
    echo ""
    echo -e "${CYAN}Starting API Gateway...${NC}"
    kill_port 3000 "gateway"
    cd "$PROJECT_DIR/gateway"
    if [ ! -d "node_modules" ]; then
        echo "  Installing dependencies..."
        yarn install --silent
    fi

    # Create .env for gateway
    cat > .env << 'ENVEOF'
PORT=3000
NODE_ENV=development

# gRPC Services
AUTH_SERVICE_URL=localhost:50051
USER_SERVICE_URL=localhost:50052
EVENT_SERVICE_URL=localhost:50053
BOOKING_SERVICE_URL=localhost:50054
PAYMENT_SERVICE_URL=localhost:50056
REALTIME_SERVICE_URL=localhost:50057
TICKET_SERVICE_URL=localhost:50058

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
ENVEOF

    yarn dev &
    GATEWAY_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Gateway started (HTTP: 3000)${NC}"
}

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping all services...${NC}"

    for pid_var in AUTH_PID USER_PID EVENT_PID BOOKING_PID PAYMENT_PID REALTIME_PID TICKET_PID BOOKING_WORKER_PID EMAIL_WORKER_PID GATEWAY_PID; do
        eval "pid=\$$pid_var"
        if [ ! -z "$pid" ]; then
            echo "  Stopping $pid_var (PID: $pid)..."
            pkill -P $pid 2>/dev/null || true
            kill -9 $pid 2>/dev/null || true
        fi
    done

    echo -e "${GREEN}All services stopped${NC}"
    echo -e "${YELLOW}Note: PostgreSQL and Redis are still running (managed by Homebrew)${NC}"
    exit 0
}

# Main execution
check_prerequisites
kill_all_dev_processes
setup_native_infrastructure

# Run migrations for other services if not skipped
if [ "$SKIP_MIGRATION" = false ]; then
    echo ""
    echo -e "${CYAN}Running database migrations...${NC}"

    # User service migrations
    if [ -d "$PROJECT_DIR/user-service/migrations" ]; then
        echo "  Running user-service migrations..."
        for f in "$PROJECT_DIR/user-service/migrations"/*.sql; do
            if [ -f "$f" ]; then
                PGPASSWORD=booking_pass psql -h localhost -U booking_user -d booking_system -f "$f" 2>/dev/null || true
            fi
        done
    fi

    # Event service migrations
    if [ -d "$PROJECT_DIR/event-service/database/migrations" ]; then
        echo "  Running event-service migrations..."
        for f in "$PROJECT_DIR/event-service/database/migrations"/*.up.sql; do
            if [ -f "$f" ]; then
                PGPASSWORD=booking_pass psql -h localhost -U booking_user -d booking_system -f "$f" 2>/dev/null || true
            fi
        done
    fi

    # Ticket service migrations
    if [ -d "$PROJECT_DIR/ticket-service/migrations" ]; then
        echo "  Running ticket-service migrations..."
        for f in "$PROJECT_DIR/ticket-service/migrations"/*.sql; do
            if [ -f "$f" ]; then
                PGPASSWORD=booking_pass psql -h localhost -U booking_user -d booking_system -f "$f" 2>/dev/null || true
            fi
        done
    fi

    echo -e "${GREEN}  ✓ Migrations completed${NC}"
fi

# Start services
echo ""
echo -e "${CYAN}════════════════════════════════════════════════${NC}"
echo -e "${CYAN}  Starting Application Services${NC}"
echo -e "${CYAN}════════════════════════════════════════════════${NC}"

should_start_service "auth" && start_auth_service
should_start_service "user" && start_user_service
should_start_service "event" && start_event_service
should_start_service "booking" && start_booking_service
should_start_service "payment" && start_payment_service
should_start_service "realtime" && start_realtime_service
should_start_service "ticket" && start_ticket_service
should_start_service "booking-worker" && start_booking_worker
should_start_service "email-worker" && start_email_worker
should_start_service "gateway" && start_gateway

# Print summary
echo ""
echo -e "${GREEN}════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  All Services Started! (Native Mode)${NC}"
echo -e "${GREEN}════════════════════════════════════════════════${NC}"
echo ""
echo -e "${CYAN}Endpoints:${NC}"
echo "  Gateway:     http://localhost:3000"
echo "  Swagger:     http://localhost:3000/api-docs"
echo ""
echo -e "${CYAN}gRPC Services:${NC}"
echo "  Auth:        localhost:50051"
echo "  User:        localhost:50052"
echo "  Event:       localhost:50053"
echo "  Booking:     localhost:50054"
echo "  Payment:     localhost:50056"
echo "  Realtime:    localhost:50057"
echo "  Ticket:      localhost:50058"
echo "  BookingWorker: localhost:50059"
echo "  EmailWorker: localhost:50060"
echo ""
echo -e "${CYAN}Infrastructure (Native):${NC}"
echo "  PostgreSQL:  localhost:5432"
echo "  Redis:       localhost:6379"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Set trap and wait
trap cleanup SIGINT SIGTERM
wait

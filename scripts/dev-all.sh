#!/bin/bash

# Development script for all microservices with hot reload
# Usage:
#   ./dev-all.sh                                    # Start all services
#   ./dev-all.sh --services auth,gateway,user       # Start specific services
#   ./dev-all.sh --infra-only                       # Start only infrastructure
#   ./dev-all.sh --help                             # Show help

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
DEPLOY_DIR="$PROJECT_DIR/deploy"

echo -e "${GREEN}"
echo "╔══════════════════════════════════════════════════╗"
echo "║  Booking System - Development Environment        ║"
echo "╚══════════════════════════════════════════════════╝"
echo -e "${NC}"

# Parse command line arguments
SELECTED_SERVICES=""
START_ALL=true
INFRA_ONLY=false
SKIP_INFRA=false
SKIP_MIGRATION=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --services)
            SELECTED_SERVICES="$2"
            START_ALL=false
            shift 2
            ;;
        --infra-only)
            INFRA_ONLY=true
            shift
            ;;
        --skip-infra)
            SKIP_INFRA=true
            shift
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
            echo "  --infra-only                Start only infrastructure (DB, Redis, Kafka)"
            echo "  --skip-infra                Skip infrastructure, start services only"
            echo "  --skip-migration            Skip database migrations"
            echo "  --help, -h                  Show this help message"
            echo ""
            echo "Available services:"
            echo "  - infra          Infrastructure (Postgres, Redis, Kafka)"
            echo "  - auth           Auth Service (Node.js, gRPC 50051)"
            echo "  - user           User Service (Go, gRPC 50052)"
            echo "  - event          Event Service (Go, gRPC 50053)"
            echo "  - booking        Booking Service (Java, HTTP 8084, gRPC 50058)"
            echo "  - payment        Payment Service (Java, gRPC 50056)"
            echo "  - realtime       Realtime Service (Go, HTTP 3003, gRPC 50057)"
            echo "  - ticket         Ticket Service (Go, gRPC 50054)"
            echo "  - booking-worker Booking Worker (Go, gRPC 50059)"
            echo "  - email-worker   Email Worker (Go, HTTP 8080, gRPC 50060)"
            echo "  - gateway        API Gateway (Node.js, HTTP 3000)"
            echo ""
            echo "Examples:"
            echo "  $0                                  # Start all"
            echo "  $0 --infra-only                     # Infrastructure only"
            echo "  $0 --services auth,gateway          # Specific services"
            echo "  $0 --skip-infra --services gateway  # Gateway only (infra already running)"
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

    if ! command -v docker &> /dev/null; then
        missing+=("docker")
    fi

    if ! command -v node &> /dev/null; then
        missing+=("node")
    fi

    if ! command -v yarn &> /dev/null; then
        missing+=("yarn")
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

    # Check Docker daemon is running
    echo -e "${CYAN}Checking Docker daemon...${NC}"
    if ! docker info > /dev/null 2>&1; then
        echo -e "${YELLOW}  Docker daemon not running. Attempting to start Colima...${NC}"
        if command -v colima &> /dev/null; then
            colima start --cpu 4 --memory 8 2>&1 | tail -5
            sleep 3
            if ! docker info > /dev/null 2>&1; then
                echo -e "${RED}  ✗ Failed to start Docker daemon${NC}"
                exit 1
            fi
        else
            echo -e "${RED}  ✗ Docker daemon not running and Colima not found${NC}"
            echo -e "${RED}  Please start Docker Desktop or run 'colima start'${NC}"
            exit 1
        fi
    fi
    echo -e "${GREEN}  ✓ Docker daemon running${NC}"

    echo -e "${GREEN}  ✓ Prerequisites OK${NC}"
}

# Start infrastructure
start_infrastructure() {
    echo ""
    echo -e "${CYAN}════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}  Starting Infrastructure (Full Stack)${NC}"
    echo -e "${CYAN}════════════════════════════════════════════════${NC}"

    cd "$DEPLOY_DIR/environments/development"

    # Use full docker-compose with Kafka
    echo -e "${CYAN}Starting all infrastructure services...${NC}"
    docker compose -f docker-compose.yml up -d

    # Wait for databases
    echo -e "${CYAN}Waiting for databases to be healthy...${NC}"
    local max_attempts=30
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if docker compose -f docker-compose.yml exec -T postgres-auth pg_isready -U booking_user -d booking_system_auth > /dev/null 2>&1 && \
           docker compose -f docker-compose.yml exec -T postgres-main pg_isready -U booking_user -d booking_system > /dev/null 2>&1; then
            echo -e "${GREEN}  ✓ Databases ready${NC}"
            break
        fi
        attempt=$((attempt + 1))
        echo -e "  Waiting... ($attempt/$max_attempts)"
        sleep 2
    done

    if [ $attempt -eq $max_attempts ]; then
        echo -e "${RED}  ✗ Databases not ready after $max_attempts attempts${NC}"
        exit 1
    fi

    # Wait for Redis
    echo -e "${CYAN}Waiting for Redis...${NC}"
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker compose -f docker-compose.yml exec -T redis redis-cli ping > /dev/null 2>&1; then
            echo -e "${GREEN}  ✓ Redis ready${NC}"
            break
        fi
        attempt=$((attempt + 1))
        echo -e "  Waiting for Redis... ($attempt/$max_attempts)"
        sleep 2
    done

    # Wait for Kafka
    echo -e "${CYAN}Waiting for Kafka to be ready...${NC}"
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        # Check if Kafka container is running and can list topics
        if docker compose -f docker-compose.yml exec -T kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
            echo -e "${GREEN}  ✓ Kafka ready${NC}"
            break
        fi
        attempt=$((attempt + 1))
        echo -e "  Waiting for Kafka... ($attempt/$max_attempts)"
        sleep 3
    done

    if [ $attempt -eq $max_attempts ]; then
        echo -e "${YELLOW}  ⚠ Kafka may not be fully ready, but continuing...${NC}"
    fi

    cd "$PROJECT_DIR"

    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}  Infrastructure Ready!${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "  ${CYAN}Databases:${NC}"
    echo "    Auth DB:  localhost:5432 (booking_system_auth)"
    echo "    Main DB:  localhost:5433 (booking_system)"
    echo ""
    echo -e "  ${CYAN}Services:${NC}"
    echo "    Redis:    localhost:6379"
    echo "    Kafka:    localhost:9092"
    echo "    Zookeeper: localhost:2181"
    echo ""
    echo -e "  ${CYAN}Credentials:${NC}"
    echo "    User:     booking_user"
    echo "    Password: booking_pass"
}

# Service PIDs tracking (Bash 3 compatible - use regular variables)
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

    # Create .env file for auth-service
    cat > .env << 'ENVEOF'
# Server Configuration
PORT=50051
HOST=0.0.0.0
NODE_ENV=development

# Database Configuration (for knexfile)
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

    # Run migrations (infrastructure already started by start_infrastructure)
    echo "  Running migrations..."
    yarn migrate:latest 2>/dev/null || true
    yarn seed:run 2>/dev/null || true

    # Start with yarn dev (not dev:local, infrastructure already started)
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
    kill_port 9192 "user-service-metrics"
    cd "$PROJECT_DIR/user-service"
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
    if ! command -v mvn &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Booking Service (Maven not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Booking Service...${NC}"
    kill_port 8084 "booking-service-http"
    kill_port 50058 "booking-service-grpc"
    cd "$PROJECT_DIR/booking-service"
    mvn spring-boot:run -Dspring-boot.run.profiles=dev &
    BOOKING_PID=$!
    cd "$PROJECT_DIR"
    sleep 10
    echo -e "${GREEN}  ✓ Booking Service started (HTTP: 8084, gRPC: 50058)${NC}"
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
    kill_port 50054 "ticket-service"
    cd "$PROJECT_DIR/ticket-service"
    if [ -f "scripts/dev.sh" ]; then
        chmod +x scripts/dev.sh
        ./scripts/dev.sh &
    else
        go run main.go &
    fi
    TICKET_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Ticket Service started (gRPC: 50054)${NC}"
}

start_booking_worker() {
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}  ⚠ Skipping Booking Worker (Go not installed)${NC}"
        return
    fi
    echo ""
    echo -e "${CYAN}Starting Booking Worker...${NC}"
    kill_port 50059 "booking-worker"
    kill_port 9091 "booking-worker-metrics"
    cd "$PROJECT_DIR/booking-worker"
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
    # Use scripts/dev.sh or go run (infrastructure already started by dev-all.sh)
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
    # Use yarn dev instead of dev:local (infrastructure already started)
    yarn dev &
    GATEWAY_PID=$!
    cd "$PROJECT_DIR"
    sleep 5
    echo -e "${GREEN}  ✓ Gateway started (HTTP: 3000)${NC}"
}

# Cleanup function (Bash 3 compatible)
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping all services...${NC}"

    # Stop each service by individual PID variable
    if [ ! -z "$AUTH_PID" ]; then
        echo "  Stopping auth (PID: $AUTH_PID)..."
        pkill -P $AUTH_PID 2>/dev/null || true
        kill -9 $AUTH_PID 2>/dev/null || true
    fi

    if [ ! -z "$USER_PID" ]; then
        echo "  Stopping user (PID: $USER_PID)..."
        pkill -P $USER_PID 2>/dev/null || true
        kill -9 $USER_PID 2>/dev/null || true
    fi

    if [ ! -z "$EVENT_PID" ]; then
        echo "  Stopping event (PID: $EVENT_PID)..."
        pkill -P $EVENT_PID 2>/dev/null || true
        kill -9 $EVENT_PID 2>/dev/null || true
    fi

    if [ ! -z "$BOOKING_PID" ]; then
        echo "  Stopping booking (PID: $BOOKING_PID)..."
        pkill -P $BOOKING_PID 2>/dev/null || true
        kill -9 $BOOKING_PID 2>/dev/null || true
    fi

    if [ ! -z "$PAYMENT_PID" ]; then
        echo "  Stopping payment (PID: $PAYMENT_PID)..."
        pkill -P $PAYMENT_PID 2>/dev/null || true
        kill -9 $PAYMENT_PID 2>/dev/null || true
    fi

    if [ ! -z "$REALTIME_PID" ]; then
        echo "  Stopping realtime (PID: $REALTIME_PID)..."
        pkill -P $REALTIME_PID 2>/dev/null || true
        kill -9 $REALTIME_PID 2>/dev/null || true
    fi

    if [ ! -z "$TICKET_PID" ]; then
        echo "  Stopping ticket (PID: $TICKET_PID)..."
        pkill -P $TICKET_PID 2>/dev/null || true
        kill -9 $TICKET_PID 2>/dev/null || true
    fi

    if [ ! -z "$BOOKING_WORKER_PID" ]; then
        echo "  Stopping booking-worker (PID: $BOOKING_WORKER_PID)..."
        pkill -P $BOOKING_WORKER_PID 2>/dev/null || true
        kill -9 $BOOKING_WORKER_PID 2>/dev/null || true
    fi

    if [ ! -z "$EMAIL_WORKER_PID" ]; then
        echo "  Stopping email-worker (PID: $EMAIL_WORKER_PID)..."
        pkill -P $EMAIL_WORKER_PID 2>/dev/null || true
        kill -9 $EMAIL_WORKER_PID 2>/dev/null || true
    fi

    if [ ! -z "$GATEWAY_PID" ]; then
        echo "  Stopping gateway (PID: $GATEWAY_PID)..."
        pkill -P $GATEWAY_PID 2>/dev/null || true
        kill -9 $GATEWAY_PID 2>/dev/null || true
    fi

    echo -e "${GREEN}All services stopped${NC}"
    exit 0
}

# Main execution
check_prerequisites
kill_all_dev_processes

# Start infrastructure if not skipped
if [ "$SKIP_INFRA" = false ]; then
    if should_start_service "infra" || [ "$INFRA_ONLY" = true ]; then
        start_infrastructure
    fi
fi

# Exit if infra only
if [ "$INFRA_ONLY" = true ]; then
    echo ""
    echo -e "${GREEN}Infrastructure started. Use --services to start application services.${NC}"
    exit 0
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
echo -e "${GREEN}  All Services Started!${NC}"
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
echo "  Ticket:      localhost:50054"
echo "  Payment:     localhost:50056"
echo "  Realtime:    localhost:50057"
echo "  Booking:     localhost:50058"
echo "  BookingWorker: localhost:50059"
echo "  EmailWorker: localhost:50060"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Set trap and wait
trap cleanup SIGINT SIGTERM
wait

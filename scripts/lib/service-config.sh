#!/bin/bash

# ============================================================================
# Shared service configuration and health check utilities
# Source this file from other scripts: source "$SCRIPT_DIR/lib/service-config.sh"
# ============================================================================

# Colors
HC_RED='\033[0;31m'
HC_GREEN='\033[0;32m'
HC_YELLOW='\033[1;33m'
HC_CYAN='\033[0;36m'
HC_NC='\033[0m'

# ============================================================================
# Service Registry
# Parallel arrays for Bash 3 compatibility
# ============================================================================

# --- Infrastructure services ---
INFRA_NAMES=("postgres-auth" "postgres-main" "redis" "kafka" "zookeeper")
INFRA_PORTS=("50432" "50433" "50379" "50092" "50181")
INFRA_CONTAINERS=("dev-postgres-auth" "dev-postgres-main" "dev-redis-cache" "dev-kafka" "dev-zookeeper")

# --- Application services ---
APP_NAMES=("auth" "user" "event" "ticket" "payment" "realtime" "booking" "booking-worker" "checkin" "invoice" "email-worker" "gateway")

# Primary port for each app service (used for basic port check)
APP_PORTS=("50051" "50052" "50053" "50054" "50062" "3003" "8084" "50056" "50059" "8083" "50061" "53000")

# gRPC ports (empty if no gRPC)
APP_GRPC_PORTS=("50051" "50052" "50053" "50054" "50062" "50057" "50058" "50056" "50059" "50060" "50061" "")

# HTTP health endpoints (empty if none)
APP_HEALTH_URLS=("" "" "" "" "" "http://localhost:3003/health" "http://localhost:8084/actuator/health" "" "" "http://localhost:8083/actuator/health" "" "http://localhost:53000/health")

# Check type: port, grpc, http, spring, http+grpc
APP_CHECK_TYPES=("port" "grpc" "port" "port" "port" "http" "spring" "port" "port" "spring" "port" "http")

# Display labels (keep under 19 chars for table alignment)
APP_PORT_LABELS=("gRPC:50051" "gRPC:50052" "gRPC:50053" "gRPC:50054" "gRPC:50062" ":3003/:50057" ":8084/:50058" "gRPC:50056" "gRPC:50059" ":8083/:50060" "gRPC:50061" "HTTP:53000")

# ============================================================================
# Log color codes for each service (for prefix_log)
# ============================================================================
get_log_color() {
    case "$1" in
        auth)          echo "36" ;;   # cyan
        user)          echo "33" ;;   # yellow
        event)         echo "35" ;;   # magenta
        ticket)        echo "32" ;;   # green
        payment)       echo "91" ;;   # light red
        realtime)      echo "94" ;;   # light blue
        booking)       echo "95" ;;   # light magenta
        booking-worker) echo "93" ;;  # light yellow
        checkin)       echo "34" ;;   # blue
        invoice)       echo "92" ;;   # light green
        email-worker)  echo "96" ;;   # light cyan
        gateway)       echo "97" ;;   # white/bright
        *)             echo "0"  ;;
    esac
}

# ============================================================================
# Health Check Functions
# ============================================================================

# Check if a TCP port is listening
# Usage: check_port <port>
check_port() {
    local port=$1
    if command -v nc &> /dev/null; then
        nc -z localhost "$port" 2>/dev/null
        return $?
    fi
    # Fallback to /dev/tcp
    (echo > /dev/tcp/localhost/"$port") 2>/dev/null
    return $?
}

# Check HTTP health endpoint
# Usage: check_http_health <url> [timeout_seconds]
check_http_health() {
    local url=$1
    local timeout=${2:-3}
    local response
    response=$(curl -s -o /dev/null -w "%{http_code}" --max-time "$timeout" "$url" 2>/dev/null)
    [ "$response" = "200" ]
}

# Check gRPC health (grpcurl if available, fallback to port check)
# Usage: check_grpc_health <port>
check_grpc_health() {
    local port=$1
    if command -v grpcurl &> /dev/null; then
        grpcurl -plaintext "localhost:$port" grpc.health.v1.Health/Check &>/dev/null
        return $?
    fi
    check_port "$port"
}

# Check Docker container health
# Usage: check_docker_health <container_name>
check_docker_health() {
    local container=$1
    local status
    # Try health status first
    status=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}no-healthcheck{{end}}' "$container" 2>/dev/null)
    if [ "$status" = "healthy" ]; then
        echo "healthy"
        return 0
    fi
    # Fallback: check if running
    local running
    running=$(docker inspect --format='{{.State.Running}}' "$container" 2>/dev/null)
    if [ "$running" = "true" ]; then
        echo "running"
        return 0
    fi
    echo "down"
    return 1
}

# Get PID listening on a port
# Usage: get_pid_for_port <port>
get_pid_for_port() {
    local port=$1
    lsof -ti:"$port" 2>/dev/null | head -1
}

# Perform health check for an app service by index
# Usage: check_app_service <index>
# Returns: 0=healthy, 1=down, 2=degraded
# Sets: HC_STATUS_TEXT (e.g., "UP", "DOWN", "DEGRADED")
# Sets: HC_CHECK_METHOD (e.g., "port", "http", "grpc")
check_app_service() {
    local idx=$1
    local port="${APP_PORTS[$idx]}"
    local check_type="${APP_CHECK_TYPES[$idx]}"
    local health_url="${APP_HEALTH_URLS[$idx]}"
    local grpc_port="${APP_GRPC_PORTS[$idx]}"

    HC_CHECK_METHOD="$check_type"

    case "$check_type" in
        http)
            if [ -n "$health_url" ] && check_http_health "$health_url"; then
                HC_STATUS_TEXT="UP"
                HC_CHECK_METHOD="http"
                return 0
            elif check_port "$port"; then
                HC_STATUS_TEXT="DEGRADED"
                HC_CHECK_METHOD="port"
                return 2
            fi
            ;;
        spring)
            if [ -n "$health_url" ] && check_http_health "$health_url"; then
                HC_STATUS_TEXT="UP"
                HC_CHECK_METHOD="actuator"
                return 0
            elif check_port "$port"; then
                HC_STATUS_TEXT="DEGRADED"
                HC_CHECK_METHOD="port"
                return 2
            fi
            ;;
        grpc)
            if check_grpc_health "$port"; then
                HC_STATUS_TEXT="UP"
                HC_CHECK_METHOD="grpc"
                return 0
            elif check_port "$port"; then
                HC_STATUS_TEXT="UP"
                HC_CHECK_METHOD="port"
                return 0
            fi
            ;;
        "http+grpc")
            local http_ok=false
            local grpc_ok=false
            if [ -n "$health_url" ] && check_http_health "$health_url"; then
                http_ok=true
            fi
            if [ -n "$grpc_port" ] && check_port "$grpc_port"; then
                grpc_ok=true
            fi
            if $http_ok && $grpc_ok; then
                HC_STATUS_TEXT="UP"
                HC_CHECK_METHOD="http+grpc"
                return 0
            elif $http_ok || $grpc_ok; then
                HC_STATUS_TEXT="DEGRADED"
                HC_CHECK_METHOD="partial"
                return 2
            fi
            ;;
        port|*)
            if check_port "$port"; then
                HC_STATUS_TEXT="UP"
                HC_CHECK_METHOD="port"
                return 0
            fi
            ;;
    esac

    HC_STATUS_TEXT="DOWN"
    return 1
}

# Wait for a service to become healthy with timeout
# Usage: wait_for_healthy <name> <port> <check_type> [timeout_seconds] [http_url]
wait_for_healthy() {
    local name=$1
    local port=$2
    local check_type=$3
    local timeout=${4:-30}
    local http_url=$5
    local elapsed=0

    while [ $elapsed -lt $timeout ]; do
        case "$check_type" in
            port)
                if check_port "$port"; then
                    echo -e "${HC_GREEN}  [OK] $name is listening on port $port (${elapsed}s)${HC_NC}"
                    return 0
                fi
                ;;
            http|spring)
                if [ -n "$http_url" ] && check_http_health "$http_url"; then
                    echo -e "${HC_GREEN}  [OK] $name health check passed (${elapsed}s)${HC_NC}"
                    return 0
                elif check_port "$port"; then
                    : # Port up, health not ready yet
                fi
                ;;
            grpc)
                if check_grpc_health "$port"; then
                    echo -e "${HC_GREEN}  [OK] $name gRPC health passed (${elapsed}s)${HC_NC}"
                    return 0
                fi
                ;;
        esac
        sleep 2
        elapsed=$((elapsed + 2))
        printf "  Waiting for %s... (%d/%ds)\r" "$name" "$elapsed" "$timeout"
    done

    # Timeout - check if at least port is up
    if check_port "$port"; then
        echo -e "${HC_YELLOW}  [WARN] $name port $port is open but health check didn't pass within ${timeout}s${HC_NC}"
        return 0
    fi

    echo -e "${HC_RED}  [FAIL] $name did not become healthy within ${timeout}s${HC_NC}"
    return 1
}

# Log prefixing function
# Usage: some_command 2>&1 | prefix_log "service-name" "color-code"
prefix_log() {
    local name=$1
    local color=$2
    local padded
    padded=$(printf "%-16s" "$name")
    # Use sed -u for unbuffered output on Linux
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -l "s/^/$(printf '\033[0;%sm[%s]\033[0m ' "$color" "$padded")/"
    else
        sed -u "s/^/$(printf '\033[0;%sm[%s]\033[0m ' "$color" "$padded")/"
    fi
}

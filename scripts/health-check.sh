#!/bin/bash

# ============================================================================
# Service Health Dashboard
# Check the status of all infrastructure and application services
#
# Usage:
#   ./scripts/health-check.sh              # One-shot check
#   ./scripts/health-check.sh --watch      # Auto-refresh every 5s
#   ./scripts/health-check.sh --watch 3    # Auto-refresh every 3s
#   ./scripts/health-check.sh --services   # App services only
#   ./scripts/health-check.sh --infra      # Infrastructure only
#   ./scripts/health-check.sh --help       # Show help
# ============================================================================

set -euo pipefail

# Disable exit-on-error for health checks (services may be down - that's expected)
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/service-config.sh"

# Defaults
WATCH_MODE=false
WATCH_INTERVAL=5
SHOW_INFRA=true
SHOW_SERVICES=true

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --watch|-w)
            WATCH_MODE=true
            if [[ "${2:-}" =~ ^[0-9]+$ ]]; then
                WATCH_INTERVAL="$2"
                shift
            fi
            shift
            ;;
        --services|-s)
            SHOW_INFRA=false
            shift
            ;;
        --infra|-i)
            SHOW_SERVICES=false
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --watch, -w [N]    Auto-refresh every N seconds (default: 5)"
            echo "  --services, -s     Show only application services"
            echo "  --infra, -i        Show only infrastructure"
            echo "  --help, -h         Show this help"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# ============================================================================
# Display Functions
# ============================================================================

# Print a colored status indicator
print_status() {
    local status=$1
    case "$status" in
        UP)       printf "${HC_GREEN}● UP      ${HC_NC}" ;;
        DOWN)     printf "${HC_RED}● DOWN    ${HC_NC}" ;;
        DEGRADED) printf "${HC_YELLOW}● DEGRADED${HC_NC}" ;;
        *)        printf "  %-8s" "$status" ;;
    esac
}

display_header() {
    local now
    now=$(date '+%Y-%m-%d %H:%M:%S')
    echo ""
    echo -e "${HC_GREEN}╔═══════════════════════════════════════════════════════════════╗${HC_NC}"
    echo -e "${HC_GREEN}║  Booking System - Service Health Dashboard                    ║${HC_NC}"
    echo -e "${HC_GREEN}║  ${HC_NC}${now}${HC_GREEN}                                              ║${HC_NC}"
    echo -e "${HC_GREEN}╚═══════════════════════════════════════════════════════════════╝${HC_NC}"
}

display_infra() {
    echo ""
    echo -e "  ${HC_CYAN}INFRASTRUCTURE${HC_NC}"
    echo -e "  ${HC_CYAN}┌──────────────────┬────────────┬────────┬───────────────────┐${HC_NC}"
    printf "  ${HC_CYAN}│${HC_NC} %-16s ${HC_CYAN}│${HC_NC} %-10s ${HC_CYAN}│${HC_NC} %-6s ${HC_CYAN}│${HC_NC} %-17s ${HC_CYAN}│${HC_NC}\n" \
        "Service" "Status" "Port" "Details"
    echo -e "  ${HC_CYAN}├──────────────────┼────────────┼────────┼───────────────────┤${HC_NC}"

    local total=0
    local healthy=0

    for i in "${!INFRA_NAMES[@]}"; do
        local name="${INFRA_NAMES[$i]}"
        local port="${INFRA_PORTS[$i]}"
        local container="${INFRA_CONTAINERS[$i]}"
        local detail=""
        local status="DOWN"

        total=$((total + 1))

        # Check Docker container
        if command -v docker &> /dev/null; then
            if detail=$(check_docker_health "$container" 2>/dev/null); then
                status="UP"
                healthy=$((healthy + 1))
            fi
            detail="Container: ${detail:-down}"
        else
            # Fallback: port check
            if check_port "$port"; then
                status="UP"
                detail="Port: open"
                healthy=$((healthy + 1))
            else
                detail="Port: closed"
            fi
        fi

        printf "  ${HC_CYAN}│${HC_NC} %-16s ${HC_CYAN}│${HC_NC} " "$name"
        print_status "$status"
        printf " ${HC_CYAN}│${HC_NC} %-6s ${HC_CYAN}│${HC_NC} %-17s ${HC_CYAN}│${HC_NC}\n" "$port" "$detail"
    done

    echo -e "  ${HC_CYAN}└──────────────────┴────────────┴────────┴───────────────────┘${HC_NC}"

    INFRA_TOTAL=$total
    INFRA_HEALTHY=$healthy
}

display_services() {
    echo ""
    echo -e "  ${HC_CYAN}APPLICATION SERVICES${HC_NC}"
    echo -e "  ${HC_CYAN}┌──────────────────┬────────────┬─────────────────────┬────────┬──────────┐${HC_NC}"
    printf "  ${HC_CYAN}│${HC_NC} %-16s ${HC_CYAN}│${HC_NC} %-10s ${HC_CYAN}│${HC_NC} %-19s ${HC_CYAN}│${HC_NC} %-6s ${HC_CYAN}│${HC_NC} %-8s ${HC_CYAN}│${HC_NC}\n" \
        "Service" "Status" "Ports" "PID" "Method"
    echo -e "  ${HC_CYAN}├──────────────────┼────────────┼─────────────────────┼────────┼──────────┤${HC_NC}"

    local total=0
    local healthy=0
    local degraded=0

    for i in "${!APP_NAMES[@]}"; do
        local name="${APP_NAMES[$i]}"
        local port="${APP_PORTS[$i]}"
        local port_label="${APP_PORT_LABELS[$i]}"
        local pid=""

        total=$((total + 1))

        # Get PID
        pid=$(get_pid_for_port "$port" 2>/dev/null)
        if [ -z "$pid" ]; then
            pid="-"
        fi

        # Health check
        check_app_service "$i"
        local result=$?
        local status="$HC_STATUS_TEXT"
        local method="$HC_CHECK_METHOD"

        if [ $result -eq 0 ]; then
            healthy=$((healthy + 1))
        elif [ $result -eq 2 ]; then
            degraded=$((degraded + 1))
        fi

        # Truncate port_label if too long
        if [ ${#port_label} -gt 19 ]; then
            port_label="${port_label:0:19}"
        fi

        printf "  ${HC_CYAN}│${HC_NC} %-16s ${HC_CYAN}│${HC_NC} " "$name"
        print_status "$status"
        printf " ${HC_CYAN}│${HC_NC} %-19s ${HC_CYAN}│${HC_NC} %-6s ${HC_CYAN}│${HC_NC} %-8s ${HC_CYAN}│${HC_NC}\n" \
            "$port_label" "$pid" "$method"
    done

    echo -e "  ${HC_CYAN}└──────────────────┴────────────┴─────────────────────┴────────┴──────────┘${HC_NC}"

    APP_TOTAL=$total
    APP_HEALTHY=$healthy
    APP_DEGRADED=$degraded
}

display_summary() {
    local total=$((INFRA_TOTAL + APP_TOTAL))
    local healthy=$((INFRA_HEALTHY + APP_HEALTHY))
    local degraded=${APP_DEGRADED:-0}
    local down=$((total - healthy - degraded))

    echo ""
    printf "  Summary: "
    echo -e "${HC_GREEN}$healthy healthy${HC_NC} | ${HC_RED}$down down${HC_NC} | ${HC_YELLOW}$degraded degraded${HC_NC} | $total total"
    echo ""
}

display_dashboard() {
    display_header

    if [ "$SHOW_INFRA" = true ]; then
        display_infra
    else
        INFRA_TOTAL=0
        INFRA_HEALTHY=0
    fi

    if [ "$SHOW_SERVICES" = true ]; then
        display_services
    else
        APP_TOTAL=0
        APP_HEALTHY=0
        APP_DEGRADED=0
    fi

    display_summary
}

# ============================================================================
# Main
# ============================================================================

if [ "$WATCH_MODE" = true ]; then
    trap 'echo ""; echo "Stopped."; exit 0' SIGINT SIGTERM
    while true; do
        clear
        display_dashboard
        echo -e "  ${HC_CYAN}Refreshing every ${WATCH_INTERVAL}s... Press Ctrl+C to exit${HC_NC}"
        sleep "$WATCH_INTERVAL"
    done
else
    display_dashboard
fi

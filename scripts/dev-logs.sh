#!/bin/bash

# ============================================================================
# Dev Log Viewer
# View and filter logs from individual or all microservices
#
# Usage:
#   ./scripts/dev-logs.sh auth              # View last 100 lines of auth log
#   ./scripts/dev-logs.sh --tail auth       # Tail auth log in realtime
#   ./scripts/dev-logs.sh --tail -n 50 auth # Tail last 50 lines first
#   ./scripts/dev-logs.sh --all             # Tail all service logs
#   ./scripts/dev-logs.sh --all --filter ERROR  # Filter for ERROR
#   ./scripts/dev-logs.sh --list            # List available log files
#   ./scripts/dev-logs.sh --help            # Show help
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
LOG_DIR="$PROJECT_DIR/logs"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# Defaults
MODE="view"
FILTER=""
LINES=100
SERVICE=""

show_help() {
    echo "Usage: $0 [OPTIONS] [SERVICE]"
    echo ""
    echo "Options:"
    echo "  --tail, -f           Tail log in realtime"
    echo "  --all, -a            Show all service logs"
    echo "  --list, -l           List available log files"
    echo "  --filter <pattern>   Filter log lines by pattern"
    echo "  -n <lines>           Number of lines to show (default: 100)"
    echo "  --help, -h           Show this help"
    echo ""
    echo "Service names (fuzzy match):"
    echo "  auth, user, event, ticket, payment"
    echo "  realtime, booking, booking-worker, email-worker, gateway"
    echo ""
    echo "Examples:"
    echo "  $0 auth                          # View last 100 lines"
    echo "  $0 --tail auth                   # Tail in realtime"
    echo "  $0 --tail -n 20 gateway          # Tail with 20 lines of context"
    echo "  $0 --all --filter 'ERROR|WARN'   # Filter errors from all logs"
    echo "  $0 --list                        # List available logs"
}

# Resolve fuzzy service name to log file name
resolve_service() {
    local input=$1
    case "$input" in
        auth*)          echo "auth-service" ;;
        user*)          echo "user-service" ;;
        event*)         echo "event-service" ;;
        ticket*)        echo "ticket-service" ;;
        payment*)       echo "payment-service" ;;
        realtime*)      echo "realtime-service" ;;
        booking-w*)     echo "booking-worker" ;;
        booking*)       echo "booking-service" ;;
        email*)         echo "email-worker" ;;
        gate*|gw)       echo "gateway" ;;
        *)              echo "$input" ;;
    esac
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --tail|-f)    MODE="tail"; shift ;;
        --all|-a)     MODE="all"; shift ;;
        --list|-l)    MODE="list"; shift ;;
        --filter)     FILTER="$2"; shift 2 ;;
        -n)           LINES="$2"; shift 2 ;;
        --help|-h)    show_help; exit 0 ;;
        -*)           echo "Unknown option: $1"; show_help; exit 1 ;;
        *)            SERVICE="$1"; shift ;;
    esac
done

# Check log directory exists
if [ ! -d "$LOG_DIR" ]; then
    echo -e "${RED}Log directory not found: $LOG_DIR${NC}"
    echo "Start services with ./scripts/dev-all.sh first."
    exit 1
fi

case "$MODE" in
    list)
        echo -e "${CYAN}Available log files:${NC}"
        echo ""
        if ls "$LOG_DIR"/*.log &>/dev/null; then
            for f in "$LOG_DIR"/*.log; do
                local_name=$(basename "$f")
                local_size=$(du -h "$f" 2>/dev/null | cut -f1)
                local_lines=$(wc -l < "$f" 2>/dev/null)
                printf "  %-24s %6s  %6d lines\n" "$local_name" "$local_size" "$local_lines"
            done
        else
            echo "  No log files found."
        fi
        echo ""
        if ls "$LOG_DIR"/*.log.prev &>/dev/null; then
            echo -e "${CYAN}Previous run logs:${NC}"
            for f in "$LOG_DIR"/*.log.prev; do
                echo "  $(basename "$f")"
            done
        fi
        ;;

    view)
        if [ -z "$SERVICE" ]; then
            echo -e "${RED}Please specify a service name.${NC}"
            echo "Usage: $0 [--tail] <service>"
            echo "Run '$0 --list' to see available logs."
            exit 1
        fi

        svc=$(resolve_service "$SERVICE")
        logfile="$LOG_DIR/${svc}.log"

        if [ ! -f "$logfile" ]; then
            echo -e "${RED}Log file not found: $logfile${NC}"
            echo "Run '$0 --list' to see available logs."
            exit 1
        fi

        echo -e "${CYAN}=== $svc (last $LINES lines) ===${NC}"
        if [ -n "$FILTER" ]; then
            grep --color=always -iE "$FILTER" "$logfile" | tail -n "$LINES"
        else
            tail -n "$LINES" "$logfile"
        fi
        ;;

    tail)
        if [ -z "$SERVICE" ]; then
            echo -e "${RED}Please specify a service name.${NC}"
            echo "Usage: $0 --tail <service>"
            exit 1
        fi

        svc=$(resolve_service "$SERVICE")
        logfile="$LOG_DIR/${svc}.log"

        if [ ! -f "$logfile" ]; then
            echo -e "${RED}Log file not found: $logfile${NC}"
            exit 1
        fi

        echo -e "${CYAN}=== Tailing $svc (Ctrl+C to stop) ===${NC}"
        if [ -n "$FILTER" ]; then
            tail -n "$LINES" -f "$logfile" | grep --color=always --line-buffered -iE "$FILTER"
        else
            tail -n "$LINES" -f "$logfile"
        fi
        ;;

    all)
        echo -e "${CYAN}=== Tailing all service logs (Ctrl+C to stop) ===${NC}"
        if [ -n "$FILTER" ]; then
            tail -f "$LOG_DIR"/*.log 2>/dev/null | grep --color=always --line-buffered -iE "$FILTER"
        else
            tail -f "$LOG_DIR"/*.log 2>/dev/null
        fi
        ;;
esac

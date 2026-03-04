#!/bin/bash
#
# Generate Go protobuf stubs for all Go services.
#
# Each proto file is generated into its own subdirectory:
#   <service>/internal/protos/<proto_name>/<proto_name>.pb.go
#
# The shared proto files already contain the correct `option go_package`
# directives, so no manipulation of the proto files is required.
#
# Usage:
#   ./scripts/generate-go-protos.sh                  # all services
#   ./scripts/generate-go-protos.sh checkin-service  # one service

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status()  { echo -e "${GREEN}[INFO]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_error()   { echo -e "${RED}[ERROR]${NC} $1"; }
print_header()  { echo -e "${BLUE}=====${NC} $1"; }

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED_PROTO_DIR="$PROJECT_ROOT/shared-lib/protos"

# ---------------------------------------------------------------------------
# Pre-flight checks
# ---------------------------------------------------------------------------
if ! command -v protoc &>/dev/null; then
    print_error "protoc not found. Install protobuf-compiler and retry."
    exit 1
fi

if ! command -v protoc-gen-go &>/dev/null; then
    print_warning "protoc-gen-go not found — installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &>/dev/null; then
    print_warning "protoc-gen-go-grpc not found — installing..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# ---------------------------------------------------------------------------
# generate_for_service <service_name>
#
# For each .proto in shared-lib/protos/, creates
#   <service>/internal/protos/<proto_name>/
# and generates .pb.go + _grpc.pb.go there.
#
# email-worker is special: only email.proto is generated.
# ---------------------------------------------------------------------------
generate_for_service() {
    local service_name="$1"
    local service_dir="$PROJECT_ROOT/$service_name"

    if [ ! -d "$service_dir" ]; then
        print_warning "Directory $service_dir not found — skipping."
        return
    fi

    if [ ! -f "$service_dir/go.mod" ]; then
        print_warning "$service_name has no go.mod — skipping (not a Go service)."
        return
    fi

    print_status "Generating protos for $service_name..."

    local proto_glob
    if [ "$service_name" = "email-worker" ]; then
        proto_glob="$SHARED_PROTO_DIR/email.proto"
    else
        proto_glob="$SHARED_PROTO_DIR/*.proto"
    fi

    for proto_file in $proto_glob; do
        [ -f "$proto_file" ] || continue

        local proto_name
        proto_name="$(basename "$proto_file" .proto)"
        local output_dir="$service_dir/internal/protos/$proto_name"
        mkdir -p "$output_dir"

        print_status "  $proto_name → $output_dir"
        protoc \
            --go_out="$output_dir" \
            --go_opt=paths=source_relative \
            --go-grpc_out="$output_dir" \
            --go-grpc_opt=paths=source_relative \
            --proto_path="$SHARED_PROTO_DIR" \
            "$proto_file"
    done

    print_status "Done: $service_name"
}

# ---------------------------------------------------------------------------
# generate_all_services — iterate over known Go services
# ---------------------------------------------------------------------------
generate_all_services() {
    print_header "Generating Go protos for all services"

    local services=(
        "email-worker"
        "auth-service"
        "booking-service"
        "booking-worker"
        "ticket-service"
        "payment-service"
        "realtime-service"
        "checkin-service"
    )

    for svc in "${services[@]}"; do
        generate_for_service "$svc"
    done
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
print_header "Go Proto Generation"
print_status "Project root : $PROJECT_ROOT"
print_status "Proto dir    : $SHARED_PROTO_DIR"

if [ $# -eq 0 ]; then
    generate_all_services
elif [ $# -eq 1 ]; then
    generate_for_service "$1"
else
    print_error "Usage: $0 [service_name]"
    exit 1
fi

print_status "Proto generation complete."

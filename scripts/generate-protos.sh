#!/bin/bash

# Generate Protocol Buffers for all Go services
# This script generates Go code from .proto files in shared-lib/protos

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PROTO_DIR="$PROJECT_ROOT/shared-lib/protos"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║  Generating Protocol Buffers for Go Services     ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════╝${NC}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc is not installed${NC}"
    echo "Install with: sudo apt-get install -y protobuf-compiler"
    exit 1
fi

# Ensure PATH includes Go bin
export PATH="$PATH:$(go env GOPATH)/bin"

# Install Go protobuf plugins if not present
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW}Installing protoc-gen-go...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${YELLOW}Installing protoc-gen-go-grpc...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# user-service: imports "user-service/internal/protos" (no subdir)
echo -e "\n${CYAN}Generating protos for user-service...${NC}"
rm -rf "$PROJECT_ROOT/user-service/internal/protos"
mkdir -p "$PROJECT_ROOT/user-service/internal/protos"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/user-service/internal/protos" \
    --go_opt=paths=source_relative \
    --go_opt=Muser.proto=user-service/internal/protos \
    --go-grpc_out="$PROJECT_ROOT/user-service/internal/protos" \
    --go-grpc_opt=paths=source_relative \
    --go-grpc_opt=Muser.proto=user-service/internal/protos \
    "$PROTO_DIR/user.proto" 2>&1 || echo -e "${YELLOW}    Warning: user.proto had issues${NC}"
echo -e "${GREEN}  ✓ user-service protos generated${NC}"

# event-service: imports "event-service/internal/protos/event" (with subdir)
echo -e "\n${CYAN}Generating protos for event-service...${NC}"
rm -rf "$PROJECT_ROOT/event-service/internal/protos"
mkdir -p "$PROJECT_ROOT/event-service/internal/protos/event"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/event-service/internal/protos/event" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/event-service/internal/protos/event" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/event.proto" 2>&1 || echo -e "${YELLOW}    Warning: event.proto had issues${NC}"
echo -e "${GREEN}  ✓ event-service protos generated${NC}"

# realtime-service: check import pattern
echo -e "\n${CYAN}Generating protos for realtime-service...${NC}"
rm -rf "$PROJECT_ROOT/realtime-service/internal/protos"
mkdir -p "$PROJECT_ROOT/realtime-service/internal/protos"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/realtime-service/internal/protos" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/realtime-service/internal/protos" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/realtime.proto" 2>&1 || echo -e "${YELLOW}    Warning: realtime.proto had issues${NC}"
echo -e "${GREEN}  ✓ realtime-service protos generated${NC}"

# ticket-service: needs ticket, payment, and event protos
echo -e "\n${CYAN}Generating protos for ticket-service...${NC}"
rm -rf "$PROJECT_ROOT/ticket-service/internal/protos"
mkdir -p "$PROJECT_ROOT/ticket-service/internal/protos/ticket"
mkdir -p "$PROJECT_ROOT/ticket-service/internal/protos/payment"
mkdir -p "$PROJECT_ROOT/ticket-service/internal/protos/event"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/ticket-service/internal/protos/ticket" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/ticket-service/internal/protos/ticket" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/ticket.proto" 2>&1 || echo -e "${YELLOW}    Warning: ticket.proto had issues${NC}"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/ticket-service/internal/protos/payment" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/ticket-service/internal/protos/payment" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/payment.proto" 2>&1 || echo -e "${YELLOW}    Warning: payment.proto had issues${NC}"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/ticket-service/internal/protos/event" \
    --go_opt=paths=source_relative \
    --go_opt=Mevent.proto=ticket-service/internal/protos/event \
    --go-grpc_out="$PROJECT_ROOT/ticket-service/internal/protos/event" \
    --go-grpc_opt=paths=source_relative \
    --go-grpc_opt=Mevent.proto=ticket-service/internal/protos/event \
    "$PROTO_DIR/event.proto" 2>&1 || echo -e "${YELLOW}    Warning: event.proto had issues${NC}"
echo -e "${GREEN}  ✓ ticket-service protos generated${NC}"

# booking-worker: multiple protos with subdirs
echo -e "\n${CYAN}Generating protos for booking-worker...${NC}"
for subdir in booking booking_worker realtime; do
    mkdir -p "$PROJECT_ROOT/booking-worker/internal/protos/$subdir"
done
protoc --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/booking-worker/internal/protos/booking" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/booking-worker/internal/protos/booking" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/booking.proto" 2>&1 || true
protoc --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/booking-worker/internal/protos/booking_worker" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/booking-worker/internal/protos/booking_worker" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/booking_worker.proto" 2>&1 || true
protoc --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/booking-worker/internal/protos/realtime" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/booking-worker/internal/protos/realtime" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/realtime.proto" 2>&1 || true
echo -e "${GREEN}  ✓ booking-worker protos generated${NC}"

# email-worker
echo -e "\n${CYAN}Generating protos for email-worker...${NC}"
rm -rf "$PROJECT_ROOT/email-worker/internal/protos"
mkdir -p "$PROJECT_ROOT/email-worker/internal/protos"
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT/email-worker/internal/protos" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT/email-worker/internal/protos" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/email.proto" 2>&1 || echo -e "${YELLOW}    Warning: email.proto had issues${NC}"
echo -e "${GREEN}  ✓ email-worker protos generated${NC}"

echo -e "\n${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Proto generation complete!                      ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"

#!/bin/bash

# Generate protobuf code for event-service
# This script generates Go code from shared-lib/protos/event.proto

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[HEADER]${NC} $1"
}

# Get the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$SERVICE_DIR/.." && pwd)"
SHARED_PROTO_DIR="$PROJECT_ROOT/shared-lib/protos"
OUTPUT_DIR="$SERVICE_DIR/internal/protos"

print_header "Event Service Proto Generation"
print_status "Service directory: $SERVICE_DIR"
print_status "Project root: $PROJECT_ROOT"
print_status "Shared proto directory: $SHARED_PROTO_DIR"
print_status "Output directory: $OUTPUT_DIR"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    print_error "protoc is not installed. Please install protobuf-compiler first."
    print_status "Installation:"
    print_status "  macOS: brew install protobuf"
    print_status "  Ubuntu: sudo apt-get install protobuf-compiler"
    exit 1
fi

# Check if Go protobuf plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    print_warning "protoc-gen-go not found. Installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    print_warning "protoc-gen-go-grpc not found. Installing..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate Go code from event.proto
PROTO_FILE="$SHARED_PROTO_DIR/event.proto"
EVENT_OUTPUT_DIR="$OUTPUT_DIR/event"

if [ ! -f "$PROTO_FILE" ]; then
    print_error "Proto file not found: $PROTO_FILE"
    exit 1
fi

print_status "Generating Go code from event.proto..."
mkdir -p "$EVENT_OUTPUT_DIR"

protoc \
    --go_out="$EVENT_OUTPUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$EVENT_OUTPUT_DIR" \
    --go-grpc_opt=paths=source_relative \
    --proto_path="$SHARED_PROTO_DIR" \
    "$PROTO_FILE"

print_status "Proto files generated successfully!"
print_status "Generated files in: $OUTPUT_DIR"
find "$OUTPUT_DIR" -name "*.pb.go" -type f | while read file; do
    print_status "  - $(basename "$file")"
done

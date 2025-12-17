#!/bin/bash

# Generate Go protobuf files for realtime-service

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SHARED_PROTO_DIR="$(cd "$PROJECT_ROOT/../shared-lib/protos" && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/internal/protos"

print_status "Generating Go protobuf files for realtime-service"
print_status "Proto source: $SHARED_PROTO_DIR"
print_status "Output directory: $OUTPUT_DIR"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    print_error "protoc is not installed. Please install protobuf-compiler first."
    exit 1
fi

# Check if Go protobuf plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    print_status "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    print_status "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate Go code from realtime.proto
print_status "Generating realtime.proto..."
protoc \
    --go_out="$OUTPUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUTPUT_DIR" \
    --go-grpc_opt=paths=source_relative \
    --proto_path="$SHARED_PROTO_DIR" \
    "$SHARED_PROTO_DIR/realtime.proto"

print_status "Proto generation completed!"
print_status "Generated files in: $OUTPUT_DIR"
ls -la "$OUTPUT_DIR"

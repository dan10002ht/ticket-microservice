#!/bin/bash

# Build script for booking-worker

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔨 Building booking-worker..."

cd "$SERVICE_DIR"

# Generate protobuf code first
if [ -f "scripts/generate-proto.sh" ]; then
    echo "📝 Generating protobuf code..."
    chmod +x scripts/generate-proto.sh
    ./scripts/generate-proto.sh
fi

# Build the service
echo "🔨 Building Go binary..."
go mod tidy
go build -o booking-worker main.go

echo "✅ Build completed!"
echo "📦 Binary: $SERVICE_DIR/booking-worker"


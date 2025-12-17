#!/bin/bash

# Development script for booking-worker
# Runs the service with hot reload (requires air or similar tool)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🚀 Starting booking-worker in development mode..."

cd "$SERVICE_DIR"

# Check if .env file exists
if [ ! -f ".env" ]; then
    if [ -f "env.example" ]; then
        echo "📋 Creating .env from env.example..."
        cp env.example .env
        echo "⚠️  Please update .env with your configuration"
    else
        echo "⚠️  No .env file found. Using default configuration."
    fi
fi

# Generate protobuf code if needed
if [ ! -d "internal/protos" ] || [ -z "$(ls -A internal/protos 2>/dev/null)" ]; then
    echo "📝 Generating protobuf code..."
    if [ -f "scripts/generate-proto.sh" ]; then
        chmod +x scripts/generate-proto.sh
        ./scripts/generate-proto.sh
    fi
fi

# Run the service
echo "🚀 Starting service..."
go run main.go


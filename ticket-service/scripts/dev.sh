#!/bin/bash

# Development script for Ticket Service

set -e

echo "🚀 Starting Ticket Service Development Environment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. Please update the configuration as needed."
fi

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Initialize database
echo "🗄️  Initializing database..."
chmod +x scripts/init-db.sh
./scripts/init-db.sh

# Generate protobuf code (if needed)
if [ -d "../shared-lib/protos" ]; then
    echo "🔧 Generating protobuf code..."
    # TODO: Add protobuf generation commands
    echo "⚠️  Protobuf generation not implemented yet"
fi

# Run tests (if any)
echo "🧪 Running tests..."
if [ -d "tests" ]; then
    go test ./...
else
    echo "⚠️  No tests found, skipping test execution"
fi

# Start the service
echo "🚀 Starting Ticket Service..."
echo "  Port: ${PORT:-3003}"
echo "  gRPC Port: ${GRPC_PORT:-50053}"
echo "  Database: ${DB_NAME:-booking_system_ticket}"
echo ""

# Run the service
go run main.go

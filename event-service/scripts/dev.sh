#!/bin/bash

# Development script for event-service
echo "🚀 Starting Event Service in development mode..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  .env file not found, using env.example..."
    if [ -f env.example ]; then
        cp env.example .env
        export $(cat .env | grep -v '^#' | xargs)
    else
        echo "❌ No environment configuration found!"
        exit 1
    fi
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Run the service
echo "✅ Starting event-service on port ${GRPC_PORT:-50053}..."
go run main.go

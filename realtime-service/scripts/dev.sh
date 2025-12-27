#!/bin/bash

# Development script for realtime-service
echo "🚀 Starting Realtime Service in development mode..."

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
echo "✅ Starting realtime-service on HTTP:${HTTP_PORT:-3003} gRPC:${GRPC_PORT:-50057}..."
go run main.go

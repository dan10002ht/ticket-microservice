#!/bin/bash

# Development script for Checkin Service

set -e

echo "Starting Checkin Service Development Environment..."

cd "$(dirname "$0")/.."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp env.example .env
    echo ".env file created. Please update the configuration as needed."
fi

# Load environment variables
set -a
source .env
set +a

# Install dependencies
echo "Installing Go dependencies..."
go mod tidy

# Start the service
echo "Starting Checkin Service on gRPC port ${GRPC_PORT:-50059}..."
go run main.go

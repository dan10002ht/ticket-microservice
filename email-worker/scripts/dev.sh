#!/bin/bash

# Development script for email-worker
echo "Starting Email Worker in development mode..."

cd "$(dirname "$0")/.."

# Load environment variables safely
if [ -f .env ]; then
    echo "Loading .env file..."
    set -a
    source .env
    set +a
else
    echo "Creating default .env file..."
    cat > .env << 'EOF'
# Email Worker Environment Variables

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=email-notifications
KAFKA_GROUP_ID=email-worker-group

# SMTP Configuration
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=noreply@ticketing.local

# gRPC Configuration
GRPC_PORT=50060

# HTTP Configuration
HTTP_PORT=8080

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Database Configuration
DB_HOST=localhost
DB_PORT=5433
DB_NAME=booking_system
DB_USER=booking_user
DB_PASSWORD=booking_pass
DB_SSL_MODE=disable

# Metrics Configuration
METRICS_ENABLED=true
METRICS_PORT=2112

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=json

# Environment
ENV=development
EOF
    set -a
    source .env
    set +a
fi

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Run the service
echo "Starting email-worker on HTTP port ${HTTP_PORT:-8080}, gRPC port ${GRPC_PORT:-50060}..."
go run main.go

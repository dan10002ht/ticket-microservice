#!/bin/bash

# Development script for booking-service (Spring Boot)
echo "🚀 Starting Booking Service in development mode..."

# Check if Maven is installed
if ! command -v mvn &> /dev/null; then
    echo "❌ Maven is not installed. Please install Maven"
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run Spring Boot application
echo "✅ Starting booking-service with Spring Boot..."
mvn spring-boot:run -Dspring-boot.run.profiles=dev

#!/bin/bash
set -e

# Example: migrate -path ./migrations -database "$DATABASE_URL" up
# Hoặc goose postgres "$DATABASE_URL" up

echo "Chạy migration cho ticket-service (bạn cần cài migrate hoặc goose)" 
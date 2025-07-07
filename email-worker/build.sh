#!/bin/bash

# Email Worker Service Build Script

set -e

echo "🚀 Building Email Worker Service..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.23+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.23"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    print_error "Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION+"
    exit 1
fi

print_status "Go version: $GO_VERSION"

# Clean previous builds
print_status "Cleaning previous builds..."
rm -rf dist/
mkdir -p dist/

# Download dependencies
print_status "Downloading dependencies..."
go mod download

# Run tests
print_status "Running tests..."
if ! go test ./...; then
    print_error "Tests failed!"
    exit 1
fi

# Run linting
print_status "Running linting..."
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run; then
        print_warning "Linting issues found, but continuing..."
    fi
else
    print_warning "golangci-lint not found, skipping linting"
fi

# Build for different platforms
print_status "Building for multiple platforms..."

# Linux AMD64
print_status "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o dist/email-worker-linux-amd64 .

# Linux ARM64
print_status "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -a -installsuffix cgo -o dist/email-worker-linux-arm64 .

# macOS AMD64
print_status "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o dist/email-worker-darwin-amd64 .

# macOS ARM64
print_status "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -a -installsuffix cgo -o dist/email-worker-darwin-arm64 .

# Build Docker image
if command -v docker &> /dev/null; then
    print_status "Building Docker image..."
    docker build -t email-worker:latest .
    print_status "Docker image built successfully"
else
    print_warning "Docker not found, skipping Docker build"
fi

# Show build results
print_status "Build completed successfully!"
echo ""
echo "📦 Build artifacts:"
ls -la dist/
echo ""
echo "🐳 Docker image: email-worker:latest"
echo ""
print_status "You can now run the worker with:"
echo "  ./dist/email-worker-linux-amd64"
echo "  docker run email-worker:latest" 
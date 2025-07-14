#!/bin/bash

# Generate Go code from proto files for venue-service
# This script generates both the protobuf and gRPC Go code

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Generating Go code from proto files...${NC}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}❌ protoc is not installed. Please install Protocol Buffers compiler.${NC}"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW}⚠️  protoc-gen-go not found. Installing...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${YELLOW}⚠️  protoc-gen-go-grpc not found. Installing...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Set up paths
PROTO_DIR="protos"
SHARED_PROTO_DIR="../shared-lib/protos"

# Create protos directory if it doesn't exist
mkdir -p $PROTO_DIR

# Copy proto files from shared-lib if they don't exist
if [ ! -f "$PROTO_DIR/venue.proto" ]; then
    echo -e "${YELLOW}📋 Copying venue.proto from shared-lib...${NC}"
    cp $SHARED_PROTO_DIR/venue.proto $PROTO_DIR/
fi

# Generate Go code
echo -e "${GREEN}🔧 Generating Go code from venue.proto...${NC}"
protoc \
    --go_out=$PROTO_DIR \
    --go-grpc_out=$PROTO_DIR \
    --proto_path=$PROTO_DIR \
    --proto_path=$SHARED_PROTO_DIR \
    $PROTO_DIR/venue.proto

# Move generated files to correct location
if [ -d "$PROTO_DIR/booking-system" ]; then
    echo -e "${YELLOW}📁 Moving generated files...${NC}"
    mv $PROTO_DIR/booking-system/venue-service/protos/* $PROTO_DIR/
    rm -rf $PROTO_DIR/booking-system
fi

# Update go.mod
echo -e "${GREEN}📦 Updating dependencies...${NC}"
go mod tidy

# Build to check for errors
echo -e "${GREEN}🔨 Building venue-service...${NC}"
go build

echo -e "${GREEN}✅ Proto generation completed successfully!${NC}"
echo -e "${GREEN}📁 Generated files:${NC}"
echo -e "   - $PROTO_DIR/venue.pb.go"
echo -e "   - $PROTO_DIR/venue_grpc.pb.go" 
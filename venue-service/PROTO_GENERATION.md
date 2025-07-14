# Proto Generation Guide - Venue Service

## Overview

This guide explains how to generate Go code from Protocol Buffer (proto) files for the venue-service.

## Prerequisites

Make sure you have the following tools installed:

1. **Protocol Buffers Compiler (protoc)**

   ```bash
   # Ubuntu/Debian
   sudo apt-get install protobuf-compiler

   # macOS
   brew install protobuf

   # Windows
   # Download from https://github.com/protocolbuffers/protobuf/releases
   ```

2. **Go Protobuf Plugin**

   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```

3. **Go gRPC Plugin**
   ```bash
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

## Quick Start

### Using the Script (Recommended)

```bash
# From venue-service directory
./scripts/generate-proto.sh
```

This script will:

- Check if required tools are installed
- Copy proto files from shared-lib if needed
- Generate Go code from proto files
- Update dependencies
- Build the service to check for errors

### Manual Generation

If you prefer to run the commands manually:

```bash
# 1. Copy proto file from shared-lib
cp ../shared-lib/protos/venue.proto protos/

# 2. Generate Go code
protoc \
    --go_out=protos \
    --go-grpc_out=protos \
    --proto_path=protos \
    --proto_path=../shared-lib/protos \
    protos/venue.proto

# 3. Move generated files to correct location
mv protos/booking-system/venue-service/protos/* protos/
rm -rf protos/booking-system

# 4. Update dependencies
go mod tidy

# 5. Build to check for errors
go build
```

## Generated Files

After running the generation script, you'll have:

- `protos/venue.proto` - Original proto definition
- `protos/venue.pb.go` - Generated protobuf Go code
- `protos/venue_grpc.pb.go` - Generated gRPC Go code

## Proto File Structure

The venue service uses the following proto services:

### VenueService

- `CreateVenue` - Create a new venue
- `GetVenue` - Get venue by ID
- `UpdateVenue` - Update venue information
- `DeleteVenue` - Delete a venue
- `ListVenues` - List venues with pagination
- `SearchVenues` - Search venues by query
- `GetVenuesByLocation` - Get venues by location

### LayoutService

- `CreateLayout` - Create a new layout
- `GetLayout` - Get layout by ID
- `UpdateLayout` - Update layout information
- `DeleteLayout` - Delete a layout
- `ListLayouts` - List layouts with pagination
- `SetDefaultLayout` - Set default layout for venue
- `GetDefaultLayout` - Get default layout for venue
- `ValidateLayout` - Validate layout configuration

### ZoneService

- `CreateZone` - Create a new seating zone
- `GetZone` - Get zone by ID
- `UpdateZone` - Update zone information
- `DeleteZone` - Delete a zone
- `ListZones` - List zones with pagination
- `GetZonesByLayout` - Get zones by layout ID
- `UpdateZoneCoordinates` - Update zone coordinates

### SeatService

- `CreateSeat` - Create a new seat
- `GetSeat` - Get seat by ID
- `UpdateSeat` - Update seat information
- `DeleteSeat` - Delete a seat
- `CreateSeatsBulk` - Create multiple seats
- `UpdateSeatsBulk` - Update multiple seats
- `DeleteSeatsBulk` - Delete multiple seats
- `GetSeatsByZone` - Get seats by zone ID
- `GetSeatsByLayout` - Get seats by layout ID
- `GetSeatCoordinates` - Get seat coordinates

## Usage in Code

### Import the generated packages

```go
import (
    "booking-system/venue-service/protos"
    "google.golang.org/grpc"
)
```

### Create a gRPC client

```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

client := protos.NewVenueServiceClient(conn)
```

### Make gRPC calls

```go
// Create venue
resp, err := client.CreateVenue(ctx, &protos.CreateVenueRequest{
    Name:        "Concert Hall",
    Description: "A beautiful concert hall",
    Address:     "123 Main St",
    City:        "New York",
    Country:     "USA",
    Capacity:    1000,
    VenueType:   "concert_hall",
})
```

## Troubleshooting

### Common Issues

1. **"protoc-gen-go: program not found"**

   - Install the Go protobuf plugin: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

2. **"protoc-gen-go-grpc: program not found"**

   - Install the Go gRPC plugin: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

3. **"unable to determine Go import path"**

   - Make sure the proto file has `option go_package = "booking-system/venue-service/protos";`

4. **Import errors in generated code**
   - Run `go mod tidy` to update dependencies
   - Make sure you're using the correct import paths

### Updating Proto Files

When you update the proto files in `shared-lib/protos/venue.proto`:

1. Copy the updated file: `cp ../shared-lib/protos/venue.proto protos/`
2. Run the generation script: `./scripts/generate-proto.sh`
3. Update your code to use the new message types or service methods

## Best Practices

1. **Always use the script** - It handles all the setup and cleanup automatically
2. **Check generated code** - Review the generated files to ensure they're correct
3. **Update tests** - When proto changes, update your tests accordingly
4. **Version control** - Commit both the proto files and generated Go files
5. **Documentation** - Update this guide when adding new services or methods

## Related Files

- `shared-lib/protos/venue.proto` - Source proto definition
- `protos/venue.proto` - Local copy of proto file
- `protos/venue.pb.go` - Generated protobuf code
- `protos/venue_grpc.pb.go` - Generated gRPC code
- `scripts/generate-proto.sh` - Generation script

# Regenerate Protobuf Code

You are a backend developer. Regenerate Go protobuf code after a `.proto` file change.

## Instructions

1. Identify which proto file and service are affected from user input: `$ARGUMENTS`
2. Read the proto file in `shared-lib/protos/` to verify changes
3. Run the protoc generation command
4. Verify the generated code compiles

## Proto Files → Services Mapping

| Proto file | Go package | Target service |
|-----------|------------|----------------|
| `event.proto` | `event` | `event-service/internal/protos/event/` |
| `auth.proto` | `auth` | `auth-service/internal/protos/auth/` |
| `ticket.proto` | `ticket` | `ticket-service/internal/protos/ticket/` |
| `booking.proto` | `booking` | `booking-service/internal/protos/booking/` |
| `user.proto` | `user` | `user-service/internal/protos/user/` |
| `checkin.proto` | `checkin` | `checkin-service/internal/protos/checkin/` |
| `invoice.proto` | `invoice` | `invoice-service/internal/protos/invoice/` |
| `payment.proto` | `payment` | `payment-service/internal/protos/payment/` |

## Generation Command

```bash
# Pattern: generate to temp dir, then copy (handles go_package nested paths)
PROTO_FILE=<proto>.proto
PKG=<package>
SERVICE=<service>

TMPDIR=$(mktemp -d) && \
protoc --go_out=$TMPDIR --go-grpc_out=$TMPDIR \
  --proto_path=shared-lib/protos $PROTO_FILE && \
cp $TMPDIR/github.com/booking-system/shared-lib/protos/$PKG/${PKG}.pb.go \
  $SERVICE/internal/protos/$PKG/${PKG}.pb.go && \
cp $TMPDIR/github.com/booking-system/shared-lib/protos/$PKG/${PKG}_grpc.pb.go \
  $SERVICE/internal/protos/$PKG/${PKG}_grpc.pb.go && \
rm -rf $TMPDIR
```

## After Generation

1. Run `cd $SERVICE && go build ./...` to verify compilation
2. Check if the gRPC controller needs updating to handle new/changed fields
3. Check if the repository queries need updating

## Common Errors

- **"protoc-gen-go: unable to determine Go import path"** → proto file missing `option go_package`
- **Files generated to wrong path** → Use temp dir pattern above, don't use `paths=source_relative`
- **Missing fields in generated code** → Verify proto file saved, re-run protoc

## User Input

$ARGUMENTS

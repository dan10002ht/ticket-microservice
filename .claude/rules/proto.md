---
paths:
  - "shared-lib/protos/**/*.proto"
---

# Proto Rules

## Field Numbering
- NEVER reuse a field number, even after deleting a field
- Always increment to the next available number
- Check existing field numbers before adding new ones

## After Modifying a Proto File

1. **Regenerate** Go code using temp dir pattern:
   ```bash
   TMPDIR=$(mktemp -d) && \
   protoc --go_out=$TMPDIR --go-grpc_out=$TMPDIR \
     --proto_path=shared-lib/protos <file>.proto && \
   cp $TMPDIR/github.com/booking-system/shared-lib/protos/<pkg>/*.pb.go \
     <service>/internal/protos/<pkg>/ && \
   rm -rf $TMPDIR
   ```

2. **Update Go controller** to map new fields (proto ↔ model)
3. **Update repository** SQL queries if new DB columns are involved
4. **Update Go model** struct with `db:"col" json:"col"` tags
5. **Run** `go build ./...` to verify

## DO NOT
- Use `paths=source_relative` (generates to wrong location)
- Edit generated `.pb.go` files directly
- Remove or renumber existing fields in use

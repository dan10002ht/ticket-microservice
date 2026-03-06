---
paths:
  - "event-service/**/*.go"
  - "auth-service/**/*.go"
  - "ticket-service/**/*.go"
  - "booking-service/**/*.go"
  - "checkin-service/**/*.go"
  - "invoice-service/**/*.go"
  - "payment-service/**/*.go"
  - "user-service/**/*.go"
  - "booking-worker/**/*.go"
  - "email-worker/**/*.go"
  - "realtime-service/**/*.go"
---

# Go Service Rules

## Project Structure (per service)
```
<service>/
├── main.go                        # Entry point, gRPC server setup
├── models/                        # Go structs with db/json tags
├── repositories/                  # SQL queries (parameterized only!)
├── grpc/                          # gRPC controller implementations
├── internal/protos/<pkg>/         # Generated .pb.go files (DO NOT edit)
├── database/
│   ├── migrations/                # SQL migration files
│   └── database.go                # DB connection setup
└── go.mod
```

## Error Handling
- Always check and handle errors — never `_ = err`
- Use `status.Errorf(codes.X, "message: %v", err)` for gRPC errors
- Wrap errors with context: `fmt.Errorf("operation: %w", err)`

## SQL Safety
- ALWAYS use parameterized queries (`$1, $2, ...`) — NEVER string concatenation
- Use `QueryRowContext` / `QueryContext` (context-aware versions)
- Use `sql.NullString`, `sql.NullInt64` for nullable columns

## gRPC Controller Pattern
```go
func (c *Controller) Method(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // 1. Validate input
    // 2. Call repository
    // 3. Map model → proto response
}
```

## After Changes
- Run `cd <service> && go build ./...` to verify compilation
- If proto fields changed, regenerate .pb.go files (see proto rules)

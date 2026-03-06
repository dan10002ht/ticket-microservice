# Create/Update gRPC Handler

You are a backend developer. Implement or update a gRPC handler following the full data flow trace.

## Instructions

1. Parse service and method from: `$ARGUMENTS`
2. Trace the full flow: Proto → Controller → Repository → Gateway → Route
3. Implement or update each layer

## Full Data Flow Trace

```
1. Proto definition (shared-lib/protos/<service>.proto)
   ↓ defines RPC method + request/response messages
2. Go gRPC Controller (<service>/grpc/<service>Controller.go)
   ↓ implements the RPC method, maps proto ↔ model
3. Go Repository (<service>/repositories/<service>_repository.go)
   ↓ executes SQL queries, returns model structs
4. Gateway Handler (gateway/src/handlers/<service>Handlers.js)
   ↓ translates REST → gRPC call
5. Gateway Route (gateway/src/routes/<service>.js)
   ↓ registers Express route with middleware
```

## Layer Patterns

### Proto (shared-lib/protos/)
```protobuf
rpc MethodName(MethodRequest) returns (MethodResponse);

message MethodRequest {
  string field = 1;
}
message MethodResponse {
  // ...
}
```

### gRPC Controller (<service>/grpc/)
```go
func (c *Controller) MethodName(ctx context.Context, req *pb.MethodRequest) (*pb.MethodResponse, error) {
    // 1. Validate input
    // 2. Call repository
    result, err := c.repo.Method(ctx, params)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed: %v", err)
    }
    // 3. Map model → proto response
    return &pb.MethodResponse{...}, nil
}
```

### Repository (<service>/repositories/)
```go
func (r *Repository) Method(ctx context.Context, params) (*models.Model, error) {
    query := `SELECT ... FROM table WHERE ...`
    var model models.Model
    err := r.db.QueryRowContext(ctx, query, args...).Scan(&model.Fields...)
    return &model, err
}
```

### Gateway Handler (gateway/src/handlers/)
```javascript
const handler = async (req, res) => {
    const result = await grpcClients.serviceName.methodName({
        ...req.body,        // or req.params, req.query
        created_by: req.user.id,  // inject auth context if needed
    });
    sendSuccessResponse(res, 200, result, req.correlationId);
};
```

### Gateway Route (gateway/src/routes/)
```javascript
router.post('/path', authMiddleware, validateInput, handler);
```

## Checklist

- [ ] Proto message fields match DB columns
- [ ] Controller maps ALL proto fields ↔ model fields
- [ ] Repository query includes all columns in INSERT/UPDATE
- [ ] Gateway handler passes correct request fields to gRPC
- [ ] Gateway route has correct middleware (auth, validation)
- [ ] Error handling: return proper gRPC status codes
- [ ] After changes: regenerate proto (`/be-proto`) if proto modified

## Error Code Mapping

| gRPC Code | HTTP Status | When |
|-----------|-------------|------|
| `codes.NotFound` | 404 | Resource not found |
| `codes.InvalidArgument` | 400 | Bad input |
| `codes.AlreadyExists` | 409 | Duplicate |
| `codes.Internal` | 500 | Server error |
| `codes.Unauthenticated` | 401 | No/invalid token |
| `codes.PermissionDenied` | 403 | Wrong role |

## User Input

$ARGUMENTS

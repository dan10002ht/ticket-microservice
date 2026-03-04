# Shared Library

This directory contains shared code and definitions used across all microservices in the booking system.

## 📁 Structure

```
shared-lib/
├── protos/           # Protocol Buffer definitions
│   ├── auth.proto    # Authentication service
│   ├── user.proto    # User management service
│   ├── device.proto  # Device management service
│   ├── security.proto # Security service
│   ├── event.proto   # Event management service
│   ├── booking.proto # Booking service
│   ├── payment.proto # Payment service
│   └── ticket.proto  # Ticket service
└── docs/             # Documentation
    └── auth-service.md
```

## 🔧 Protocol Buffers

### Centralized Proto Approach

**IMPORTANT**: All services should use proto files from `shared-lib/protos/` instead of maintaining their own copies. This ensures:

- **Consistency**: All services use the same interface definitions
- **Maintainability**: Single source of truth for API contracts
- **Version Control**: Centralized versioning of service interfaces
- **Type Safety**: Consistent data types across all services

### Usage Guidelines

1. **Import from shared-lib**: Always reference proto files from `shared-lib/protos/`
2. **No Local Copies**: Do not create duplicate proto files in individual services
3. **Version Synchronization**: Keep all services using the same proto versions
4. **Breaking Changes**: Coordinate proto changes across all affected services

### Service Mappings

| Service          | Proto File         | Package Name | gRPC Port |
| ---------------- | ------------------ | ------------ | --------- |
| Auth Service     | `auth.proto`       | `auth`       | 50051     |
| User Service     | `user.proto`       | `user`       | 50052     |
| Ticket Service   | `ticket.proto`     | `ticket`     | 50053     |
| Event Service    | `event.proto`      | `event`      | 50055     |
| Booking Worker   | (queue consumer)   | -            | 50056     |
| Realtime Service | `realtime.proto`   | `realtime`   | 50057     |
| Booking Service  | `booking.proto`    | `booking`    | 50058     |
| Checkin Service  | `checkin.proto`    | `checkin`    | 50059     |
| Invoice Service  | `invoice.proto`    | `invoice`    | 50060     |
| Email Worker     | (queue consumer)   | -            | 50061     |
| Payment Service  | `payment.proto`    | `payment`    | 50062     |

### Loading Proto Files

Services should implement proto loading logic that prioritizes shared-lib:

```javascript
const loadProto = (protoFile) => {
  // 1. Try Docker shared proto (production)
  const dockerSharedProtoPath = path.join("/shared-lib", "protos", protoFile);

  // 2. Try local shared proto (development)
  const localSharedProtoPath = path.join(
    __dirname,
    "..",
    "shared-lib",
    "protos",
    protoFile
  );

  // 3. Fallback to local proto (legacy)
  const localProtoPath = path.join(__dirname, "protos", protoFile);

  let protoPath;
  if (fs.existsSync(dockerSharedProtoPath)) {
    protoPath = dockerSharedProtoPath;
  } else if (fs.existsSync(localSharedProtoPath)) {
    protoPath = localSharedProtoPath;
  } else {
    protoPath = localProtoPath;
  }

  return protoLoader.loadSync(protoPath, {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  });
};
```

## 📚 Documentation

- `docs/auth-service.md` - Authentication service integration guide

## 🔄 Version Management

### Proto Versioning

- Use semantic versioning for proto files
- Document breaking changes in proto definitions
- Maintain backward compatibility when possible
- Coordinate updates across all services

### Migration Guide

When updating proto files:

1. **Update shared-lib**: Modify proto files in `shared-lib/protos/`
2. **Regenerate Code**: All services should regenerate their gRPC code
3. **Test Integration**: Verify inter-service communication
4. **Deploy Together**: Deploy all affected services simultaneously

## 🚀 Best Practices

1. **Single Source of Truth**: Always use `shared-lib/protos/` as the authoritative source
2. **Consistent Naming**: Follow established naming conventions across all proto files
3. **Documentation**: Keep proto definitions well-documented with comments
4. **Validation**: Use proper validation rules in proto definitions
5. **Error Handling**: Define consistent error response structures
6. **Health Checks**: Include health check endpoints in all services

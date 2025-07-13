# Authorization System Guide

## Overview

The authorization system provides comprehensive permission and role-based access control (RBAC) for the booking system. It integrates with the auth-service via gRPC and includes caching for optimal performance.

## Architecture

```
Gateway → Authorization Client → gRPC → Auth Service → Database
                ↓
        Cache Layer (Redis)
```

## Components

### 1. Authorization Client (`src/clients/authorizationClient.js`)

Functional client that communicates with auth-service via gRPC:

- **checkPermission(userId, permission, correlationId)**: Check if user has specific permission
- **checkResourcePermission(userId, permission, resourceType, resourceId, correlationId)**: Check resource-specific permissions
- **getUserRoles(userId, correlationId)**: Get all user roles
- **batchCheckPermissions(userId, permissions, correlationId)**: Check multiple permissions at once
- **clearUserCache(userId)**: Clear cache for specific user
- **clearAllCache()**: Clear all permission cache
- **getCacheStats()**: Get cache statistics

### 2. Authorization Middleware (`src/middlewares/authorizationMiddleware.js`)

Express middleware functions for different authorization scenarios:

#### Basic Permission Check

```javascript
const requirePermission = (permissions) => {
  // permissions can be string or array of strings
  // Returns true if user has ANY of the specified permissions
};
```

#### Resource Permission Check

```javascript
const requireResourcePermission = (permission, resourceExtractor) => {
  // Checks permission for specific resource (e.g., event, booking)
  // resourceExtractor is a function that extracts resource info from request
};
```

#### Role-Based Access

```javascript
const requireRole = (roles) => {
  // roles is array of role names
  // Returns true if user has ANY of the specified roles
};
```

#### All Permissions Required

```javascript
const requireAllPermissions = (permissions) => {
  // permissions is array of permission names
  // Returns true only if user has ALL specified permissions
};
```

#### Ownership Check

```javascript
const requireOwnership = (ownerExtractor) => {
  // Checks if user owns the resource
  // ownerExtractor is a function that extracts owner ID from request
};
```

## Usage Examples

### 1. Basic Permission Check

```javascript
import { requirePermission } from '../middlewares/authorizationMiddleware.js';

router.get('/admin/dashboard', requirePermission('admin.dashboard.view'), (req, res) => {
  res.json({ message: 'Admin dashboard accessed' });
});
```

### 2. Multiple Permissions (Any of them)

```javascript
router.post('/events', requirePermission(['event.create', 'event.manage']), (req, res) => {
  res.json({ message: 'Event created' });
});
```

### 3. All Permissions Required

```javascript
router.put(
  '/events/:id/publish',
  requireAllPermissions(['event.edit', 'event.publish']),
  (req, res) => {
    res.json({ message: 'Event published' });
  }
);
```

### 4. Role-Based Access

```javascript
router.get('/admin/users', requireRole(['admin', 'super_admin']), (req, res) => {
  res.json({ message: 'User list accessed' });
});
```

### 5. Resource-Specific Permission

```javascript
import {
  requireResourcePermission,
  extractResourceFromParams,
} from '../middlewares/authorizationMiddleware.js';

router.get(
  '/events/:id',
  requireResourcePermission('event.view', extractResourceFromParams('event')),
  (req, res) => {
    res.json({ message: 'Event details accessed' });
  }
);
```

### 6. Ownership Check

```javascript
import {
  requireOwnership,
  extractOwnerFromParams,
} from '../middlewares/authorizationMiddleware.js';

router.delete(
  '/bookings/:bookingId',
  requireOwnership(extractOwnerFromParams('userId')),
  (req, res) => {
    res.json({ message: 'Booking deleted' });
  }
);
```

### 7. Complex Authorization Chain

```javascript
router.post(
  '/events/:eventId/bookings',
  requireResourcePermission('event.book', extractResourceFromParams('event')),
  requirePermission('booking.create'),
  (req, res) => {
    res.json({ message: 'Booking created' });
  }
);
```

## Helper Functions

### Resource Extractors

```javascript
// Extract from URL parameters
const extractResourceFromParams = (resourceType, idParam = 'id') => {
  return (req) => ({
    type: resourceType,
    id: req.params[idParam],
  });
};

// Extract from request body
const extractResourceFromBody = (resourceType, idField = 'id') => {
  return (req) => ({
    type: resourceType,
    id: req.body[idField],
  });
};
```

### Owner Extractors

```javascript
// Extract from URL parameters
const extractOwnerFromParams = (ownerParam = 'userId') => {
  return (req) => req.params[ownerParam];
};

// Extract from request body
const extractOwnerFromBody = (ownerField = 'userId') => {
  return (req) => req.body[ownerField];
};
```

## Caching

The authorization system includes intelligent caching:

- **Permission Cache**: 5-minute TTL for permission checks
- **Role Cache**: 5-minute TTL for user roles
- **Resource Cache**: 5-minute TTL for resource-specific permissions
- **Cache Keys**: Structured as `permission:userId:permission` or `resource:userId:permission:type:id`

### Cache Management

```javascript
import { clearUserCache, clearAllCache, getCacheStats } from '../clients/authorizationClient.js';

// Clear cache for specific user
clearUserCache('user123');

// Clear all cache
clearAllCache();

// Get cache statistics
const stats = getCacheStats();
console.log(stats);
```

## Error Handling

The system handles various error scenarios gracefully:

### Authentication Errors (401)

- No token provided
- Invalid token
- Token expired

### Authorization Errors (403)

- Insufficient permissions
- Missing required roles
- Resource access denied
- Ownership verification failed

### Service Errors (500)

- Auth service unavailable
- Database connection issues
- gRPC communication failures

## Performance Considerations

### Optimization Strategies

1. **Caching**: All permission checks are cached for 5 minutes
2. **Batch Operations**: Use `batchCheckPermissions` for multiple checks
3. **Circuit Breaker**: Built-in circuit breaker for gRPC calls
4. **Connection Pooling**: gRPC client with connection pooling
5. **Retry Logic**: Automatic retry with exponential backoff

### Best Practices

1. **Use Appropriate Middleware**: Choose the right middleware for your use case
2. **Cache Wisely**: Clear cache when permissions change
3. **Batch When Possible**: Group permission checks together
4. **Handle Errors**: Always handle authorization errors gracefully
5. **Log Appropriately**: Use debug logging for performance monitoring

## Migration from Legacy Auth

The system maintains backward compatibility:

```javascript
// Legacy middleware (still available)
import { requireRole, requirePermission } from '../middlewares/auth.js';

// New middleware (recommended)
import {
  requireRole as requireRoleNew,
  requirePermission as requirePermissionNew,
} from '../middlewares/authorizationMiddleware.js';
```

## Testing

### Unit Testing

```javascript
// Mock the authorization client
jest.mock('../clients/authorizationClient.js', () => ({
  checkPermission: jest.fn(),
  checkResourcePermission: jest.fn(),
  getUserRoles: jest.fn(),
  batchCheckPermissions: jest.fn(),
}));

// Test different scenarios
describe('Authorization', () => {
  it('should allow access with valid permission', async () => {
    checkPermission.mockResolvedValue(true);
    // Test implementation
  });
});
```

### Integration Testing

1. Test with real auth-service
2. Verify caching behavior
3. Test error scenarios
4. Performance testing

## Monitoring

### Key Metrics

- Authorization request latency
- Cache hit/miss ratios
- Permission check success/failure rates
- gRPC call performance
- Error rates by type

### Logging

The system provides comprehensive logging:

```javascript
logger.debug('Permission check completed', {
  userId: req.user.id,
  permission: 'event.create',
  hasPermission: true,
  correlationId: req.correlationId,
  cached: false,
});
```

## Security Considerations

1. **Fail-Safe**: On errors, deny access by default
2. **Input Validation**: All inputs are validated
3. **Rate Limiting**: Built-in rate limiting for auth endpoints
4. **Audit Logging**: All authorization decisions are logged
5. **Token Validation**: JWT tokens are properly validated

## Troubleshooting

### Common Issues

1. **Permission Denied**: Check user roles and permissions in database
2. **Cache Issues**: Clear cache if permissions were recently updated
3. **gRPC Errors**: Check auth-service connectivity
4. **Performance**: Monitor cache hit rates and gRPC latency

### Debug Mode

Enable debug logging to troubleshoot issues:

```javascript
// Set log level to debug
logger.setLevel('debug');
```

## Future Enhancements

1. **Dynamic Permissions**: Real-time permission updates
2. **Permission Inheritance**: Hierarchical permission system
3. **Advanced Caching**: Redis-based distributed caching
4. **Audit Dashboard**: Web interface for permission management
5. **API Rate Limiting**: Per-user permission check limits

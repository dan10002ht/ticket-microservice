# Auth Service Integration Flows

## Overview

This document describes the integration flows between Auth Service and other microservices in the ticketing system.

## Architecture

### Service Connections

```
                         ┌─────────────────┐
                         │     Gateway     │
                         │   (Port 3000)   │
                         └────────┬────────┘
                                  │ gRPC
                                  ▼
                         ┌─────────────────┐
                         │  Auth Service   │
                         │  (Port 50051)   │
                         └────────┬────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
              ▼                   ▼                   ▼
     ┌────────────────┐  ┌────────────────┐  ┌────────────────┐
     │  User Service  │  │ Booking Service│  │ Realtime Svc   │
     │  (Port 50052)  │  │  (Port 50056)  │  │  (Port 50057)  │
     └────────────────┘  └────────────────┘  └────────────────┘
```

### Communication Protocol

- **Protocol**: gRPC (Remote Procedure Call)
- **Serialization**: Protocol Buffers (protobuf)
- **Transport**: HTTP/2
- **Security**: JWT tokens for authentication

## Business Flows

### 1. User Authentication Flow

**Flow**: `Client → Gateway → Auth Service`

```
1. User submits login credentials (email/password)
2. Gateway forwards request to Auth Service via gRPC
3. Auth Service validates credentials against database
4. Auth Service generates JWT access token + refresh token
5. Gateway returns tokens to client
```

**API Endpoint**: `POST /api/auth/login`

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user-uuid",
      "email": "user@example.com",
      "role": "user"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "refresh-token-uuid",
    "expires_in": 3600
  }
}
```

### 2. User Registration Flow

**Flow**: `Client → Gateway → Auth Service`

```
1. User submits registration data
2. Gateway forwards to Auth Service
3. Auth Service validates email uniqueness
4. Auth Service hashes password and creates user
5. Auth Service generates verification token (optional)
6. Gateway returns success response
```

**API Endpoint**: `POST /api/auth/register`

**Request Body**:

```json
{
  "email": "newuser@example.com",
  "password": "securePassword123",
  "name": "John Doe"
}
```

### 3. Token Refresh Flow

**Flow**: `Client → Gateway → Auth Service`

```
1. Client sends refresh token
2. Gateway forwards to Auth Service
3. Auth Service validates refresh token
4. Auth Service generates new access token
5. Gateway returns new token to client
```

**API Endpoint**: `POST /api/auth/refresh`

### 4. Protected Resource Access Flow

**Flow**: `Client → Gateway → Target Service`

```
1. Client sends request with JWT in Authorization header
2. Gateway validates JWT using Auth Service's public key
3. Gateway extracts user_id from JWT
4. Gateway forwards request to target service with user context
5. Target service processes request and returns response
```

**Authentication Header**:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 5. User Logout Flow

**Flow**: `Client → Gateway → Auth Service`

```
1. User initiates logout
2. Gateway forwards to Auth Service
3. Auth Service invalidates refresh token
4. Auth Service removes session (if applicable)
5. Gateway returns success
```

**API Endpoint**: `POST /api/auth/logout`

## Integration with Other Services

### Auth Service → User Service

When a user registers, Auth Service creates the auth record, and User Service can be called to create the profile:

```
Auth Service                    User Service
     │                               │
     │ ──── CreateProfile ────────▶  │
     │                               │
     │ ◀──── Profile Created ─────   │
```

### Gateway → Auth Service (JWT Validation)

Gateway uses Auth Service for token validation on protected routes:

```javascript
// Gateway middleware
const validateToken = async (req, res, next) => {
  const token = req.headers.authorization?.split(" ")[1];

  try {
    // Validate with Auth Service or verify locally with public key
    const decoded = await authService.validateToken(token);
    req.userId = decoded.user_id;
    next();
  } catch (error) {
    res.status(401).json({ error: "Unauthorized" });
  }
};
```

### Booking Service → Auth Context

Booking Service receives user context from Gateway headers:

```java
// Booking Service receives user_id from Gateway
String userId = metadata.get("x-user-id");
```

## gRPC Methods

### Auth Service Proto

Location: `shared-lib/protos/auth.proto`

| Method           | Description                           |
| ---------------- | ------------------------------------- |
| `Login`          | Authenticate user with email/password |
| `Register`       | Create new user account               |
| `ValidateToken`  | Validate JWT access token             |
| `RefreshToken`   | Generate new access token             |
| `Logout`         | Invalidate refresh token              |
| `GetUser`        | Get user by ID                        |
| `UpdatePassword` | Change user password                  |

## Environment Variables

```bash
# Auth Service
GRPC_PORT=50051
JWT_SECRET=your-jwt-secret
JWT_EXPIRY=3600
REFRESH_TOKEN_EXPIRY=604800

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=auth_service
DB_USER=postgres
DB_PASSWORD=postgres

# Redis (session storage)
REDIS_HOST=localhost
REDIS_PORT=6379
```

## Security Features

### JWT Configuration

- **Algorithm**: HS256 (or RS256 for asymmetric)
- **Access Token Expiry**: 1 hour (configurable)
- **Refresh Token Expiry**: 7 days (configurable)
- **Token Storage**: Refresh tokens stored in database

### Password Security

- **Hashing**: bcrypt with salt rounds = 10
- **Validation**: Minimum 8 characters, complexity requirements

### Rate Limiting

- **Login attempts**: 5 per minute per IP
- **Registration**: 3 per hour per IP
- **Token refresh**: 10 per minute per user

## Testing

### Test gRPC Connection

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Test login
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "password123"
}' localhost:50051 auth.AuthService/Login
```

### Test via Gateway

```bash
# Login
curl -X POST http://localhost:53000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'

# Access protected route
curl http://localhost:53000/api/users/profile \
  -H "Authorization: Bearer <access_token>"
```

## Health Checks

```bash
# gRPC health check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# HTTP health check (if enabled)
curl http://localhost:50051/health
```

## Error Handling

### Common Error Codes

| Code                    | Description                        |
| ----------------------- | ---------------------------------- |
| `INVALID_CREDENTIALS`   | Email or password incorrect        |
| `USER_NOT_FOUND`        | User does not exist                |
| `EMAIL_EXISTS`          | Email already registered           |
| `TOKEN_EXPIRED`         | JWT has expired                    |
| `TOKEN_INVALID`         | JWT signature invalid              |
| `REFRESH_TOKEN_INVALID` | Refresh token not found or expired |

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

## Related Documentation

- [Auth Service README](../../auth-service/README.md)
- [Gateway Documentation](../../gateway/README.md)
- [Service Connections](../../architecture/SERVICE_CONNECTIONS.md)
- [User Service](../../user-service/README.md)

---

**Last Updated**: December 2024

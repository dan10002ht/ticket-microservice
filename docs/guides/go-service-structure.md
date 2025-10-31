# Boilerplate Service - Folder Structure

```
boilerplate-service/
├── main.go                    # Application entry point
├── go.mod                     # Go modules
├── go.sum                     # Dependencies checksums
├── Dockerfile                 # Multi-stage Docker build
├── env.example                # Environment variables example
├── README.md                  # Documentation
├── folder-structure.md        # This file
├── build.sh                   # Build script
├── internal/                  # Private application code
│   ├── app/                   # Application initialization and lifecycle
│   │   └── app.go             # Main application logic
│   ├── config/                # Configuration management
│   │   └── config.go          # Configuration structs and loading
│   ├── logger/                # Logging configuration
│   │   └── logger.go          # Logger setup with zap
│   └── server/                # HTTP server
│       └── server.go          # HTTP server with Gin
├── database/                  # Database connection and models
│   ├── connection.go          # PostgreSQL connection
│   ├── models/                # Database models
│   │   ├── user.go            # User model
│   │   ├── event.go           # Event model
│   │   └── booking.go         # Booking model
│   └── repositories/          # Data access layer
│       ├── user_repo.go       # User repository
│       ├── event_repo.go      # Event repository
│       └── booking_repo.go    # Booking repository
├── queue/                     # Message queue clients
│   ├── redis.go               # Redis client
│   ├── kafka.go               # Kafka consumer/producer
│   └── job_queue.go           # Job queue interface
├── grpcclient/                # gRPC client connections
│   ├── clients.go             # gRPC client manager
│   ├── auth_client.go         # Auth service client
│   ├── user_client.go         # User service client
│   └── booking_client.go      # Booking service client
├── services/                  # Business logic services
│   ├── services.go            # Service container
│   ├── user_service.go        # User business logic
│   ├── event_service.go       # Event business logic
│   ├── booking_service.go     # Booking business logic
│   └── admin_service.go       # Admin operations
├── providers/                 # External service providers
│   ├── provider.go            # Provider interface
│   ├── email_provider.go      # Email provider
│   └── payment_provider.go    # Payment provider
├── processor/                 # Job processing logic
│   ├── processor.go           # Main processor
│   ├── job_consumer.go        # Job consumer
│   ├── job_validator.go       # Job validation
│   └── batch_processor.go     # Batch processing
├── templates/                 # Template management
│   ├── templates.go           # Template manager
│   ├── renderer.go            # Template renderer
│   └── cache.go               # Template caching
├── metrics/                   # Prometheus metrics
│   ├── metrics.go             # Metrics definitions
│   ├── collector.go           # Metrics collector
│   └── exporter.go            # Metrics exporter
├── logging/                   # Logging configuration
│   ├── logger.go              # Logger setup
│   ├── middleware.go          # Logging middleware
│   └── formatter.go           # Log formatters
├── utils/                     # Utility functions
│   ├── crypto.go              # Cryptographic utilities
│   ├── validation.go          # Validation helpers
│   ├── time.go                # Time utilities
│   └── errors.go              # Error handling
├── models/                    # Domain models
│   ├── user.go                # User domain model
│   ├── event.go               # Event domain model
│   ├── booking.go             # Booking domain model
│   └── common.go              # Common models
├── types/                     # Type definitions
│   ├── user_types.go          # User-related types
│   ├── event_types.go         # Event-related types
│   └── common_types.go        # Common types
├── constants/                 # Constants
│   ├── user_constants.go      # User-related constants
│   ├── event_constants.go     # Event-related constants
│   └── system_constants.go    # System constants
├── tests/                     # Test files
│   ├── unit/                  # Unit tests
│   ├── integration/           # Integration tests
│   ├── benchmarks/            # Benchmark tests
│   └── fixtures/              # Test data
├── scripts/                   # Build and deployment scripts
│   ├── build.sh               # Build script
│   ├── deploy.sh              # Deployment script
│   └── migrate.sh             # Database migration script
├── docs/                      # Documentation
│   ├── api.md                 # API documentation
│   ├── deployment.md          # Deployment guide
│   └── troubleshooting.md     # Troubleshooting guide
├── protos/                    # Protocol buffer definitions
│   ├── user.proto             # User service protos
│   ├── event.proto            # Event service protos
│   └── booking.proto          # Booking service protos
└── grpc/                      # gRPC server
    └── server.go              # gRPC server implementation
```

## 📁 Key Directories Explained

### **internal/**

- **app/**: Application initialization and lifecycle management
- **config/**: Configuration management and environment variables
- **logger/**: Logging setup with structured logging
- **server/**: HTTP server with Gin framework

### **database/**

- Database connections and connection pooling
- Data models and repositories
- Database migrations

### **queue/**

- Message queue clients (Redis, Kafka)
- Job queue interfaces
- Queue management

### **grpcclient/**

- gRPC client connections
- Service-to-service communication
- Connection pooling

### **services/**

- Business logic implementation
- Service orchestration
- Core application logic

### **providers/**

- External service provider implementations
- Provider abstraction
- Fallback mechanisms

### **processor/**

- Job processing logic
- Consumer implementations
- Batch processing

### **templates/**

- Template management
- Template rendering
- Template caching

### **metrics/**

- Prometheus metrics
- Performance monitoring
- Health checks

### **logging/**

- Structured logging
- Log formatting
- Log levels

### **utils/**

- Utility functions
- Helper methods
- Common functionality

### **models/**

- Domain models
- Data structures
- Business entities

### **types/**

- Type definitions
- Custom types
- Type safety

### **constants/**

- Application constants
- Configuration constants
- System constants

## 🎯 Best Practices Followed

1. **Separation of Concerns**: Each directory has a specific responsibility
2. **Dependency Injection**: Services are injected rather than created directly
3. **Interface Segregation**: Clear interfaces for each component
4. **Error Handling**: Centralized error handling and logging
5. **Configuration Management**: Environment-based configuration
6. **Testing Structure**: Organized test files with fixtures
7. **Documentation**: Comprehensive documentation in each directory
8. **Monitoring**: Built-in metrics and health checks
9. **Logging**: Structured logging throughout the application
10. **Security**: Secure handling of sensitive data

## 🚀 Benefits of This Structure

- **Scalability**: Easy to add new features and services
- **Maintainability**: Clear organization and separation of concerns
- **Testability**: Easy to write unit and integration tests
- **Deployment**: Simple deployment and configuration
- **Monitoring**: Built-in observability and monitoring
- **Security**: Secure by design with proper abstractions

## 📝 Usage

This boilerplate follows the same pattern as `email-worker` service:

1. **Clean main.go**: Only initialization and orchestration
2. **Internal modules**: All business logic in `internal/`
3. **Structured logging**: Using zap logger
4. **Configuration**: Environment-based with validation
5. **Graceful shutdown**: Proper cleanup on exit
6. **Health checks**: Built-in health and readiness endpoints
7. **Metrics**: Prometheus metrics integration
8. **gRPC support**: Both client and server capabilities

# 🤖 AI Agent Project Overview

This file is intended for AI agents and automated tools to understand the structure, conventions, and integration points of the booking-system project.

## 🏗️ Project Structure

- **Microservices Architecture**: Each service is in its own folder at the root level.
- **Service Naming**: All services use `-service` or `-worker` suffixes for clarity.
- **Shared Libraries**: Common code and protobuf definitions are in `shared-lib/`.
- **Deployment**: Orchestration and deployment files are in `deploy/`.

## 📦 Core Services

- **gateway/**: API Gateway (Node.js, HTTP :53000), entrypoint for all client requests, Swagger docs
- **auth-service/**: Authentication and authorization (Node.js, gRPC :50051)
- **user-service/**: User profiles and address management (Go, gRPC :50052)
- **event-service/**: Event and venue management (Go, gRPC :50053)
- **ticket-service/**: Ticket inventory and types (Go, gRPC :50054)
- **booking-worker/**: Queue-based booking worker (Go, gRPC :50056), handles high concurrency
- **realtime-service/**: WebSocket real-time updates (Go, HTTP :3003, gRPC :50057)
- **booking-service/**: Booking orchestration with Saga pattern (Java, HTTP :8084, gRPC :50058)
- **checkin-service/**: Event check-in and QR code scanning (Go, gRPC :50059)
- **invoice-service/**: Invoice generation and PDF export (Java, HTTP :8083, gRPC :50060)
- **email-worker/**: Asynchronous email sending (Go, gRPC :50061)
- **payment-service/**: Payment processing with Stripe (Java, HTTP :8080, gRPC :50062)

## 🔗 Communication

- **gRPC**: All inter-service communication uses gRPC with protobuf definitions in `shared-lib/protos/`.
- **Kafka/Redis**: Used for distributed queues and event-driven flows (e.g., booking-worker, analytics).
- **WebSocket**: Real-time updates via `realtime-service`.

## 🧩 Integration Points

- **Booking Flow**: gateway → booking-worker (Go) → booking-service/ticket-service → payment-service → realtime-service
- **Check-in Flow**: gateway → checkin-service (validates ticket QR code, records check-in)
- **Invoice Flow**: payment confirmed → invoice-service generates invoice → PDF available via gateway
- **Queue Handling**: booking-worker manages distributed queue, notifies clients via realtime-service
- **Timeouts**: booking-worker enforces timeouts, releases tickets if payment not completed
- **Email Flow**: auth-service → email-worker (verification emails, notifications)

## 📝 Conventions

- **Service Naming**: Use lowercase, hyphen-separated, and `-service`/`-worker` suffix
- **gRPC Naming**: Service and method names use PascalCase in protobuf
- **Environment Variables**: All configuration via env vars, documented in each service's README
- **Monitoring**: All services expose Prometheus metrics
- **Security**: JWT/mTLS for gRPC, API keys for external APIs

## 💻 Coding Conventions

### JavaScript/Node.js Services

- **Use Functions, Not Classes**: Prefer functional programming approach with exported functions instead of classes
- **File Naming**: Use camelCase.js instead of service.js pattern
  - ✅ `userService.js` → `userService.js`
  - ✅ `authHandler.js` → `authHandler.js`
  - ✅ `databaseConnection.js` → `databaseConnection.js`
  - ❌ `user.service.js` → `userService.js`
  - ❌ `auth.handler.js` → `authHandler.js`
- **Module Exports**: Export functions directly, not class instances
- **Dependency Injection**: Pass dependencies as function parameters rather than class constructor injection

## 🗂️ File/Folder Conventions

- `README.md` in each service: Human-readable documentation
- `AI_README.md` (this file): Machine/AI agent documentation
- `shared-lib/protos/`: All protobuf definitions
- `deploy/`: Docker Compose, Kubernetes, and CI/CD files

## 🧠 AI Agent Guidance

- **Service Discovery**: List all folders at root, ignore `shared-lib`, `deploy`, and markdown files for service enumeration
- **API Discovery**: Parse gRPC proto files in `shared-lib/protos/` for service/method definitions
- **Queue/Stream Discovery**: Look for Kafka/Redis usage in service configs and code
- **Integration**: Use gRPC endpoints and service names for cross-service calls
- **Security**: Always authenticate via JWT or mTLS for gRPC

## 🚦 Example Booking Flow (AI Perspective)

1. **Client** sends booking request to `gateway`
2. `gateway` forwards to `booking-worker` (Go)
3. `booking-worker` enqueues request, notifies client via `realtime-service`
4. When turn comes, `booking-worker` reserves ticket via `booking-service`/`ticket-service`
5. Client proceeds to payment via `payment-service`
6. On success, notifications sent via `realtime-service`/`email-worker`

---

This file is maintained for AI agents and automated tools. For human documentation, see the main `README.md` and service-specific READMEs.

# ✅ Phase 1 - Project Setup (Planned)

This document will track booking-worker setup milestones once implementation starts.

## 🧱 Planned Deliverables
- Go module initialization (go.mod, go.sum)
- Standardized folder structure (cmd, internal, pkg, queue, grpc)
- Base dependencies (go-redis, grpc-go, zap, prometheus)
- Dockerfile + docker-compose instructions
- Development environment setup (local Redis, gRPC stubs)

## 📋 Checklist (to be updated)
- [ ] Initialize Go module (`go mod init github.com/ticketing/booking-worker`)
- [ ] Set up project structure:
  ```
  booking-worker/
  ├── cmd/
  │   └── server/
  │       └── main.go
  ├── internal/
  │   ├── config/
  │   ├── queue/
  │   ├── processor/
  │   ├── grpc/
  │   └── metrics/
  ├── pkg/
  │   └── logger/
  ├── go.mod
  ├── go.sum
  ├── Dockerfile
  └── README.md
  ```
- [ ] Add core dependencies:
  - `github.com/redis/go-redis/v9` - Redis client
  - `google.golang.org/grpc` - gRPC client/server
  - `go.uber.org/zap` - Structured logging
  - `github.com/prometheus/client_golang` - Metrics
  - `github.com/go-redsync/redsync/v4` - Distributed locking
- [ ] Create Dockerfile (multi-stage build)
- [ ] Add docker-compose integration
- [ ] Set up local development (env.example, scripts)

## 🧰 Tools & Dependencies
- Go 1.21+
- Redis (for queue backend)
- Kafka (optional, for queue backend)
- gRPC stubs (from shared-lib/protos)

## 📝 Notes
- Follow Go best practices (package naming, error handling)
- Use dependency injection pattern for testability
- Add structured logging from the start
- Wire Prometheus metrics early

_Last updated: Planning stage (2024)_


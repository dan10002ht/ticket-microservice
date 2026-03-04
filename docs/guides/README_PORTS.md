# Service Ports Reference

## Application Services

### gRPC Ports

| Service           | Port  | Language | Description                         |
|-------------------|-------|----------|-------------------------------------|
| Auth Service      | 50051 | Node.js  | Authentication, JWT tokens          |
| User Service      | 50052 | Go       | User profiles, addresses            |
| Ticket Service    | 50054 | Go       | Ticket inventory & lifecycle        |
| Event Service     | 50053 | Go       | Event management, zones, seats      |
| Booking Worker    | 50056 | Go       | Queue processing (gRPC server)      |
| Realtime Service  | 50057 | Go       | Real-time notifications             |
| Booking Service   | 50058 | Java     | Booking orchestration (Saga)        |
| Checkin Service   | 50059 | Go       | Event check-in (QR/barcode)         |
| Invoice Service   | 50060 | Java     | Invoice generation & PDF            |
| Email Worker      | 50061 | Go       | Email delivery                      |
| Payment Service   | 50062 | Java     | Payment processing (multi-gateway)  |

### HTTP Ports

| Service           | Port | Description                    |
|-------------------|------|--------------------------------|
| Gateway           | 53000 | REST API Gateway               |
| Realtime Service  | 3003 | WebSocket server               |
| Invoice Service   | 8083 | REST API + Actuator            |
| Booking Service   | 8084 | REST API + Actuator            |
| Payment Service   | 8080 | REST API + Actuator            |

### Metrics Ports (Prometheus `/metrics` or `/actuator/prometheus`)

| Service           | Port  | Endpoint              |
|-------------------|-------|-----------------------|
| Gateway           | 53000 | /metrics              |
| Auth Service      | 53001 | /metrics              |
| Event Service     | 9095  | /metrics              |
| Ticket Service    | 9096  | /metrics              |
| Checkin Service   | 2112  | /metrics              |
| User Service      | 9092  | /metrics              |
| Booking Worker    | 9091  | /metrics              |
| Email Worker      | 9090  | /metrics              |
| Realtime Service  | 9057  | /metrics              |
| Booking Service   | 8084  | /actuator/prometheus  |
| Payment Service   | 8085  | /actuator/prometheus  |
| Invoice Service   | 8083  | /actuator/prometheus  |

---

## Infrastructure Services (Docker Compose — Development)

| Service      | Host Port | Container Port | Notes                         |
|--------------|-----------|----------------|-------------------------------|
| postgres-auth | 5432     | 5432           | Auth DB only                  |
| postgres-main | 5433     | 5432           | All other services (shared)   |
| redis-cache   | 6379     | 6379           | Session, rate limiting, cache |
| redis-queue   | 6380     | 6379           | Booking queue — noeviction!   |
| redis-pubsub  | 6381     | 6379           | WebSocket pub/sub             |
| kafka         | 9092     | 9092           | Event streaming               |
| zookeeper     | 2181     | 2181           | Kafka coordination            |
| prometheus    | 9090     | 9090           | Metrics collection            |
| grafana       | 3001     | 3000           | Dashboards (admin/admin)      |

---

## Port Allocation Strategy

```
50051-50062  : gRPC services (increment for each new service)
3000-3010    : HTTP/WebSocket services
8080-8089    : Java Spring Boot HTTP
2112, 9057,
9091-9097    : Prometheus /metrics (Go services)
9090         : Prometheus server, email-worker metrics
```

### Next available ports (for new services)

| Type    | Next Available |
|---------|----------------|
| gRPC    | 50063          |
| HTTP    | 3004           |
| Metrics | 9097           |

---

## Quick Reference

```bash
# Gateway
curl http://localhost:53000/health
curl http://localhost:53000/api/v1/auth/health   # v1 routes
curl http://localhost:53000/api/auth/health      # legacy (deprecated)

# User Service (gRPC)
grpcurl -plaintext localhost:50052 list

# Event Service (gRPC)
grpcurl -plaintext localhost:50053 list

# Ticket Service (gRPC)
grpcurl -plaintext localhost:50054 list

# Checkin Service (gRPC)
grpcurl -plaintext localhost:50059 list

# Booking Service
curl http://localhost:8084/actuator/health

# Payment Service
curl http://localhost:8080/actuator/health

# Invoice Service
curl http://localhost:8083/actuator/health

# Realtime Service
curl http://localhost:3003/health
wscat -c ws://localhost:3003/ws
```

---

## Notes

- **Port 9090**: Used exclusively by the Prometheus server. Email-worker metrics also expose on 9090 (different host — no runtime conflict, but note when configuring).
- **Shared database**: `postgres-main` is shared by all services except auth.
  Each service uses its own PostgreSQL schema (e.g., `tickets`, `booking`, `invoice`).
- **Redis separation**: Three separate Redis instances for different eviction policies.
  Never point booking-worker at redis-cache — it uses noeviction policy.

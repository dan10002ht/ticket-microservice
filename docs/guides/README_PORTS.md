# Service Ports Reference

## Application Services

### gRPC Ports

| Service | Port | Language | Description |
|---------|------|----------|-------------|
| Auth Service | 50051 | Node.js | Authentication, JWT tokens |
| User Service | 50052 | Go | User profiles, addresses |
| Event Service | 50054 | Node.js | Event management |
| Ticket Service | 50055 | Node.js | Ticket/seat management |
| Booking Service | 50056 | Java | Booking orchestration (Saga) |
| Realtime Service | 50057 | Go | Real-time notifications (gRPC) |
| Payment Service | 50058 | Java | Payment processing |

### HTTP Ports

| Service | Port | Description |
|---------|------|-------------|
| Gateway | 3000 | API Gateway (REST) |
| Realtime Service | 3003 | WebSocket server |
| Booking Service | 8080 | REST API + Actuator |
| Payment Service | 8081 | REST API + Actuator |

### Metrics Ports (Prometheus)

| Service | Port | Endpoint |
|---------|------|----------|
| Gateway | 9090 | /metrics |
| Booking Service | 9091 | /actuator/prometheus |
| User Service | 9092 | /metrics |
| Payment Service | 9092 | /actuator/prometheus |
| Booking Worker | 9093 | /metrics |
| Email Worker | 9094 | /metrics |
| Realtime Service | 9057 | /metrics |

---

## Infrastructure Services (Docker Compose - WSL2)

Để tránh xung đột với các cổng mặc định trên WSL2/Ubuntu, các service trong docker-compose.dev.yml đã được đổi port như sau:

| Service | Host Port | Container Port | Ghi chú |
|---------|-----------|----------------|---------|
| PostgreSQL | 55432 | 5432 | Database |
| Redis | 56379 | 6379 | Cache/Queue |
| Zookeeper | 52181 | 2181 | Kafka dependency |
| Kafka | 59092 | 9092 | Message broker |
| Prometheus | 59090 | 9090 | Metrics collection |
| Grafana | 53001 | 3000 | Dashboard |
| Elasticsearch | 59200 | 9200 | Search engine |
| Kibana | 55601 | 5601 | Elasticsearch UI |

---

## Port Allocation Strategy

### Convention

```
50051-50059  : gRPC services
3000-3010    : HTTP/WebSocket services
8080-8089    : Java Spring Boot services
9090-9099    : Prometheus metrics
55xxx        : Docker infrastructure (WSL2)
```

### Available Ports (for new services)

| Range | Next Available | Purpose |
|-------|----------------|---------|
| gRPC | 50059 | New gRPC service |
| HTTP | 3004 | New HTTP service |
| Metrics | 9095 | New metrics endpoint |

---

## Quick Reference

### Local Development

```bash
# Gateway
curl http://localhost:3000/health

# User Service (gRPC)
grpcurl -plaintext localhost:50052 list

# Realtime Service
curl http://localhost:3003/health
wscat -c ws://localhost:3003/ws

# Booking Service
curl http://localhost:8080/actuator/health

# Payment Service
curl http://localhost:8081/actuator/health
```

### Docker Development (WSL2)

```bash
# PostgreSQL
psql -h localhost -p 55432 -U postgres

# Redis
redis-cli -p 56379 ping

# Kafka
kafka-topics.sh --list --bootstrap-server localhost:59092

# Prometheus
curl http://localhost:59090/api/v1/targets

# Grafana
# Open: http://localhost:53001
```

---

## Notes

- Khi kết nối tới các infrastructure services từ máy host (WSL2), hãy sử dụng Host Port
- Application services chạy trực tiếp (không qua Docker) sử dụng port mặc định
- Nếu cần thêm service mới, nên chọn port theo convention ở trên

# Development Scripts

## Quick Start

### Start ALL services
```bash
cd /Users/dantt1002/projects/ticket-mcrsv
./scripts/dev-all.sh
```

### Start specific services only
```bash
# Start only infrastructure and gateway
./scripts/dev-all.sh --services infra,auth,gateway

# Start only Go services
./scripts/dev-all.sh --services user,event,ticket,realtime

# Start full stack (recommended for development)
./scripts/dev-all.sh --services infra,auth,user,event,gateway
```

### Show help
```bash
./scripts/dev-all.sh --help
```

---

## Available Services

| Service | Type | Ports | Description |
|---------|------|-------|-------------|
| `infra` | Docker | Multiple | PostgreSQL, Redis, Kafka, Grafana, Prometheus, Elasticsearch |
| `auth` | Node.js | 50051 | Authentication & Authorization service |
| `user` | Go | 50052, 9092 | User profile management service |
| `event` | Go | 50053 | Event management service |
| `booking` | Java/Spring | 50054 | Booking service |
| `payment` | Java/Spring | 8081, 50056 | Payment processing service |
| `realtime` | Go | 3003, 50057, 9057 | Real-time notifications service |
| `ticket` | Go | 50058 | Ticket management service |
| `booking-worker` | Go | 50059, 9091 | Background worker for bookings |
| `email-worker` | Go | 8080, 50060, 2112 | Email sending worker |
| `gateway` | Node.js | 53000 | API Gateway |

---

## Prerequisites

### Required
- **Node.js** 18+
- **Yarn** package manager
- **Docker** & Docker Compose
- **Git**

### Optional (for specific services)
- **Go** 1.19+ (for Go services: user, event, ticket, realtime, booking-worker, email-worker)
- **Java** 17+ & **Maven** (for Java services: booking, payment)

### Check your installation
```bash
node --version    # Should be 18+
yarn --version
docker --version
go version        # Should be 1.19+
mvn --version     # Should be 3.x
```

---

## Usage Examples

### Development Scenarios

#### 1. **Full Stack Development**
Start everything for full system testing:
```bash
./scripts/dev-all.sh
```

#### 2. **Gateway Development Only**
```bash
./scripts/dev-all.sh --services infra,auth,gateway
```

#### 3. **Backend Services Development**
```bash
./scripts/dev-all.sh --services infra,auth,user,event,booking
```

#### 4. **Frontend + Gateway Development**
```bash
# Start backend infrastructure
./scripts/dev-all.sh --services infra,auth,user,event,booking,payment

# Then run your frontend separately
cd frontend
npm run dev
```

---

## Service Dependencies

Services have the following dependencies:

```
infra (must start first)
  ├── auth (depends on postgres-auth, redis)
  ├── user (depends on postgres-main, redis)
  ├── event (depends on postgres-main, redis)
  ├── booking (depends on postgres-main, redis, kafka)
  ├── payment (depends on postgres-main, redis, kafka)
  ├── ticket (depends on postgres-main, redis)
  ├── realtime (depends on redis)
  ├── booking-worker (depends on redis, realtime, booking)
  ├── email-worker (depends on postgres-main, redis, kafka)
  └── gateway (depends on all services above)
```

**Recommended startup order:**
1. `infra` - Always start first
2. `auth`, `user`, `event` - Core services
3. `booking`, `payment`, `ticket` - Business services
4. `realtime`, `booking-worker`, `email-worker` - Workers
5. `gateway` - Start last

---

## Accessing Services

### API Endpoints

| Service | Endpoint | Description |
|---------|----------|-------------|
| Gateway | http://localhost:53000 | Main API entry point |
| Email Worker | http://localhost:8080 | Email service HTTP API |
| Realtime | http://localhost:3003 | WebSocket/SSE endpoint |
| Payment Service | http://localhost:8081 | Payment HTTP API |

### gRPC Endpoints

All services expose gRPC on their respective ports. Use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) or [BloomRPC](https://github.com/bloomrpc/bloomrpc) to test.

### Development Tools

| Tool | URL | Credentials |
|------|-----|-------------|
| Grafana | http://localhost:53001 | admin/admin |
| Prometheus | http://localhost:59090 | - |
| Kibana | http://localhost:55601 | - |

### Databases

| Database | Host | Port | Database | User | Password |
|----------|------|------|----------|------|----------|
| Auth DB (Master) | localhost | 55432 | booking_system_auth | booking_user | booking_pass |
| Auth DB (Slave 1) | localhost | 55433 | booking_system_auth | booking_user | booking_pass |
| Main DB (Master) | localhost | 55435 | booking_system | booking_user | booking_pass |
| Main DB (Slave 1) | localhost | 55436 | booking_system | booking_user | booking_pass |
| Redis | localhost | 56379 | - | - | - |

---

## Troubleshooting

### Port Already in Use
```bash
# Find and kill process using a specific port
lsof -ti:53000 | xargs kill -9

# Or use the cleanup script
./scripts/stop-dev-with-pgpool.sh
```

### Service Won't Start
1. Check logs in the terminal window
2. Verify dependencies are installed:
   ```bash
   cd <service-directory>
   yarn install  # For Node.js services
   go mod tidy   # For Go services
   mvn clean install  # For Java services
   ```

### Database Connection Issues
1. Ensure infrastructure is running:
   ```bash
   docker ps | grep postgres
   ```
2. Wait for databases to be ready (can take 30-60 seconds)
3. Check database logs:
   ```bash
   cd deploy
   docker compose -f docker-compose.dev.yml logs postgres-master
   ```

### Hot Reload Not Working
- **Node.js services**: Using `nodemon`, should auto-reload
- **Go services**: Need to restart manually or use `air` for hot reload
- **Java services**: Changes require Maven recompile

---

## Stopping Services

### Stop all services
Press `Ctrl+C` in the terminal where `dev-all.sh` is running.

Or use the stop script:
```bash
./scripts/stop-dev-with-pgpool.sh
```

### Stop specific services
Kill the process manually or use:
```bash
# Find PID
ps aux | grep "go run main.go"
ps aux | grep "node"
ps aux | grep "mvn"

# Kill the process
kill -9 <PID>
```

---

## Advanced Configuration

### Environment Variables

Each service uses its own `.env` file. Copy from `env.example`:
```bash
cd auth-service
cp env.example .env
# Edit .env with your configuration
```

### Custom Ports

Edit the service's `.env` file or `docker-compose.dev.yml` for infrastructure.

---

## Performance Tips

1. **Use SSD** for Docker volumes (much faster database I/O)
2. **Allocate enough RAM** to Docker (minimum 8GB recommended)
3. **Start only needed services** using `--services` flag
4. **Use Docker Dashboard** to monitor resource usage

---

## Common Development Workflows

### Adding a new service
1. Create the service directory
2. Add dev script to `scripts/dev.sh`
3. Update `dev-all.sh` to include the service
4. Test with: `./scripts/dev-all.sh --services infra,<your-service>`

### Testing changes
```bash
# Start minimal stack
./scripts/dev-all.sh --services infra,auth,<service-to-test>

# Make changes to your service
# Hot reload will pick up changes automatically

# Test via Gateway
curl http://localhost:53000/api/...
```

---

## Need Help?

- Check service-specific README in each service directory
- View logs: Each service outputs to its terminal
- Infrastructure logs: `docker compose -f deploy/docker-compose.dev.yml logs -f <service>`

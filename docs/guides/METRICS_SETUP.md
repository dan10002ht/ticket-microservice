# 📊 Metrics Setup Guide

## Overview

This document describes the comprehensive metrics setup for the booking system using Prometheus, Grafana, and AlertManager.

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Gateway       │    │   Auth Service  │    │  Email Worker   │
│   (Node.js)     │    │   (Node.js)     │    │     (Go)        │
│   /metrics      │    │   /metrics      │    │   /metrics      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Prometheus Server                            │
│              (Scrapes metrics every 15s)                       │
└─────────────────┬───────────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                    AlertManager                                 │
│              (Handles alerts & notifications)                   │
└─────────────────┬───────────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Grafana                                      │
│              (Visualization & Dashboards)                      │
└─────────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### 1. Start Monitoring Stack

```bash
# Start all monitoring services
docker-compose -f deploy/docker-compose.dev.yml up -d prometheus grafana alertmanager node-exporter redis-exporter postgres-exporter

# Start application services
docker-compose -f deploy/docker-compose.dev.yml up -d gateway auth-service email-worker
```

### 2. Access Monitoring Tools

- **Prometheus**: http://localhost:59090
- **Grafana**: http://localhost:53001 (admin/admin)
- **AlertManager**: http://localhost:59093

### 3. Setup Grafana Data Source

1. Login to Grafana (admin/admin)
2. Go to Configuration → Data Sources
3. Add Prometheus data source:
   - URL: `http://prometheus:9090`
   - Access: Server (default)

## 📊 Available Metrics

### Gateway Metrics

- `http_request_duration_seconds` - Request latency
- `http_requests_total` - Request count by status
- `gateway_active_connections` - Active connections
- `gateway_upstream_latency_seconds` - Upstream service latency
- `gateway_rate_limit_exceeded_total` - Rate limit violations

### Auth Service Metrics

- `auth_attempts_total` - Authentication attempts
- `auth_active_sessions` - Active sessions
- `jwt_validations_total` - JWT validation count
- `password_reset_requests_total` - Password reset requests

### Email Worker Metrics

- `email_jobs_processed_total` - Email processing count
- `email_job_processing_duration_seconds` - Processing time
- `email_queue_size` - Queue size
- `email_delivery_success_total` - Successful deliveries
- `email_delivery_failure_total` - Failed deliveries

### Infrastructure Metrics

- `node_cpu_seconds_total` - CPU usage
- `node_memory_MemTotal_bytes` - Memory usage
- `node_filesystem_size_bytes` - Disk usage
- `redis_connected_clients` - Redis connections
- `pg_stat_database_numbackends` - PostgreSQL connections

## 🚨 Alert Rules

### Critical Alerts

- **Service Down**: Any service stops responding
- **High Error Rate**: >10% 5xx errors on gateway
- **High Auth Failure**: >20% authentication failures
- **Low Email Delivery**: <95% email delivery success

### Warning Alerts

- **High Latency**: 95th percentile >1s for gateway
- **Queue Backlog**: Email queue >1000 items
- **High Resource Usage**: CPU >80%, Memory >85%, Disk >85%

## 📈 Grafana Dashboards

### 1. System Overview Dashboard

**Panels:**

- Request rate by service
- Response time percentiles
- Error rate by service
- Active connections
- Resource usage (CPU, Memory, Disk)

**Query Example:**

```promql
# Request rate
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, http_request_duration_seconds)

# Error rate
rate(http_requests_total{status_code=~"5.."}[5m]) / rate(http_requests_total[5m])
```

### 2. Service-Specific Dashboards

#### Gateway Dashboard

- Request rate by endpoint
- Upstream service latency
- Rate limiting metrics
- Active connections

#### Auth Service Dashboard

- Authentication attempts by method
- JWT validation success rate
- Active sessions count
- Password reset requests

#### Email Worker Dashboard

- Email processing rate
- Queue size over time
- Delivery success rate
- Processing duration

## 🔧 Configuration Files

### Prometheus Configuration

- `deploy/prometheus.yml` - Main Prometheus config
- `deploy/alert_rules.yml` - Alert rules

### AlertManager Configuration

- `deploy/alertmanager.yml` - Alert routing and notifications

### Docker Compose

- `deploy/docker-compose.dev.yml` - Service definitions

## 🛠️ Implementation Details

### Node.js Services (Gateway & Auth)

```javascript
// Import shared metrics
import {
  metricsMiddleware,
  metricsEndpoint,
  serviceMetrics,
} from "../shared-lib/metrics.js";

// Apply middleware
app.use(metricsMiddleware("gateway"));

// Add metrics endpoint
app.get("/metrics", metricsEndpoint);

// Record business metrics
serviceMetrics.gateway.activeConnections.set(getActiveConnections());
```

### Go Service (Email Worker)

```go
// Initialize metrics
metrics.Init()

// Record metrics
metrics.EmailJobsProcessed.WithLabelValues("success", "welcome", "welcome_email").Inc()
metrics.EmailQueueSize.Set(float64(queueSize))
```

## 📊 Monitoring Best Practices

### 1. Metric Naming

- Use descriptive names
- Include units in help text
- Use consistent naming conventions

### 2. Labeling

- Add service labels for filtering
- Include relevant business context
- Avoid high cardinality labels

### 3. Alerting

- Set appropriate thresholds
- Use different severity levels
- Include meaningful descriptions

### 4. Dashboard Design

- Group related metrics
- Use appropriate visualization types
- Include time range selectors

## 🔍 Troubleshooting

### Common Issues

1. **Metrics not appearing**

   - Check service is running
   - Verify metrics endpoint is accessible
   - Check Prometheus targets page

2. **High cardinality**

   - Review label usage
   - Consider metric aggregation
   - Use recording rules

3. **Alert noise**
   - Adjust thresholds
   - Use alert grouping
   - Implement alert inhibition

### Debug Commands

```bash
# Check Prometheus targets
curl http://localhost:59090/api/v1/targets

# Check metrics endpoint
curl http://localhost:53000/metrics

# Check alert rules
curl http://localhost:59090/api/v1/rules

# Check alert manager
curl http://localhost:59093/api/v1/alerts
```

## 📚 Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [AlertManager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [PromQL Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)

## 🔄 Maintenance

### Regular Tasks

1. **Review alert thresholds** - Monthly
2. **Update dashboards** - As needed
3. **Monitor metric cardinality** - Weekly
4. **Backup configurations** - Before changes
5. **Update dependencies** - Quarterly

### Performance Optimization

1. **Use recording rules** for complex queries
2. **Optimize scrape intervals** based on needs
3. **Monitor Prometheus resource usage**
4. **Implement metric aggregation** where appropriate

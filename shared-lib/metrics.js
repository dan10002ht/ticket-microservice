import prometheus from "prom-client";

// Common HTTP metrics
export const httpRequestDuration = new prometheus.Histogram({
  name: "http_request_duration_seconds",
  help: "Duration of HTTP requests in seconds",
  labelNames: ["method", "route", "status_code", "service"],
  buckets: [0.1, 0.3, 0.5, 0.7, 1, 3, 5, 7, 10],
});

export const httpRequestsTotal = new prometheus.Counter({
  name: "http_requests_total",
  help: "Total number of HTTP requests",
  labelNames: ["method", "route", "status_code", "service"],
});

export const httpRequestSize = new prometheus.Histogram({
  name: "http_request_size_bytes",
  help: "Size of HTTP requests in bytes",
  labelNames: ["method", "route", "service"],
  buckets: [100, 1000, 5000, 10000, 50000, 100000],
});

export const httpResponseSize = new prometheus.Histogram({
  name: "http_response_size_bytes",
  help: "Size of HTTP responses in bytes",
  labelNames: ["method", "route", "status_code", "service"],
  buckets: [100, 1000, 5000, 10000, 50000, 100000],
});

// gRPC metrics
export const grpcRequestDuration = new prometheus.Histogram({
  name: "grpc_request_duration_seconds",
  help: "Duration of gRPC requests in seconds",
  labelNames: ["method", "service", "status"],
  buckets: [0.01, 0.05, 0.1, 0.5, 1, 2, 5],
});

export const grpcRequestsTotal = new prometheus.Counter({
  name: "grpc_requests_total",
  help: "Total number of gRPC requests",
  labelNames: ["method", "service", "status"],
});

// Database metrics
export const dbQueryDuration = new prometheus.Histogram({
  name: "db_query_duration_seconds",
  help: "Duration of database queries in seconds",
  labelNames: ["operation", "table", "service"],
  buckets: [0.001, 0.01, 0.1, 0.5, 1, 2, 5],
});

export const dbConnectionsActive = new prometheus.Gauge({
  name: "db_connections_active",
  help: "Number of active database connections",
  labelNames: ["service"],
});

// Cache metrics
export const cacheHits = new prometheus.Counter({
  name: "cache_hits_total",
  help: "Total number of cache hits",
  labelNames: ["cache_type", "service"],
});

export const cacheMisses = new prometheus.Counter({
  name: "cache_misses_total",
  help: "Total number of cache misses",
  labelNames: ["cache_type", "service"],
});

// Business metrics
export const businessEventsTotal = new prometheus.Counter({
  name: "business_events_total",
  help: "Total number of business events",
  labelNames: ["event_type", "service"],
});

// Service-specific metrics
export const serviceMetrics = {
  gateway: {
    activeConnections: new prometheus.Gauge({
      name: "gateway_active_connections",
      help: "Number of active connections to gateway",
    }),
    upstreamLatency: new prometheus.Histogram({
      name: "gateway_upstream_latency_seconds",
      help: "Upstream service latency",
      labelNames: ["upstream_service"],
      buckets: [0.01, 0.05, 0.1, 0.5, 1, 2, 5],
    }),
    rateLimitExceeded: new prometheus.Counter({
      name: "gateway_rate_limit_exceeded_total",
      help: "Total number of rate limit violations",
      labelNames: ["client_ip", "endpoint"],
    }),
  },
  auth: {
    authAttempts: new prometheus.Counter({
      name: "auth_attempts_total",
      help: "Total authentication attempts",
      labelNames: ["method", "status", "user_type"],
    }),
    activeSessions: new prometheus.Gauge({
      name: "auth_active_sessions",
      help: "Number of active sessions",
    }),
    jwtValidations: new prometheus.Counter({
      name: "jwt_validations_total",
      help: "Total JWT token validations",
      labelNames: ["status"],
    }),
    passwordResetRequests: new prometheus.Counter({
      name: "password_reset_requests_total",
      help: "Total password reset requests",
      labelNames: ["status"],
    }),
  },
  email: {
    emailsSent: new prometheus.Counter({
      name: "emails_sent_total",
      help: "Total emails sent",
      labelNames: ["template", "status", "provider"],
    }),
    emailQueueSize: new prometheus.Gauge({
      name: "email_queue_size",
      help: "Current email queue size",
    }),
    emailProcessingDuration: new prometheus.Histogram({
      name: "email_processing_duration_seconds",
      help: "Email processing duration in seconds",
      labelNames: ["template"],
      buckets: [0.1, 0.5, 1, 2, 5, 10],
    }),
  },
};

// Metrics middleware for Express
export const metricsMiddleware = (serviceName) => (req, res, next) => {
  const start = Date.now();
  const originalSend = res.send;

  // Track request size
  const requestSize = req.headers["content-length"] || 0;
  httpRequestSize
    .labels(req.method, req.route?.path || req.path, serviceName)
    .observe(parseInt(requestSize));

  // Override send to track response size
  res.send = function (data) {
    const responseSize = Buffer.byteLength(data);
    httpResponseSize
      .labels(
        req.method,
        req.route?.path || req.path,
        res.statusCode,
        serviceName
      )
      .observe(responseSize);

    originalSend.call(this, data);
  };

  res.on("finish", () => {
    const duration = (Date.now() - start) / 1000;

    httpRequestDuration
      .labels(
        req.method,
        req.route?.path || req.path,
        res.statusCode,
        serviceName
      )
      .observe(duration);

    httpRequestsTotal
      .labels(
        req.method,
        req.route?.path || req.path,
        res.statusCode,
        serviceName
      )
      .inc();
  });

  next();
};

// gRPC metrics interceptor
export const grpcMetricsInterceptor = (serviceName) => {
  return {
    intercept: (call, nextCall) => {
      const start = Date.now();

      return new prometheus.Gauge({
        name: "grpc_active_requests",
        help: "Number of active gRPC requests",
        labelNames: ["method", "service"],
      }).inc();

      return nextCall(call).then(
        (response) => {
          const duration = (Date.now() - start) / 1000;
          grpcRequestDuration
            .labels(call.method, serviceName, "OK")
            .observe(duration);
          grpcRequestsTotal.labels(call.method, serviceName, "OK").inc();
          return response;
        },
        (error) => {
          const duration = (Date.now() - start) / 1000;
          grpcRequestDuration
            .labels(call.method, serviceName, "ERROR")
            .observe(duration);
          grpcRequestsTotal.labels(call.method, serviceName, "ERROR").inc();
          throw error;
        }
      );
    },
  };
};

// Metrics endpoint
export const metricsEndpoint = async (req, res) => {
  res.set("Content-Type", prometheus.register.contentType);
  res.end(await prometheus.register.metrics());
};

// Health check endpoint
export const healthCheckEndpoint = (req, res) => {
  res.json({
    status: "healthy",
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
  });
};

// Initialize default metrics
prometheus.collectDefaultMetrics({
  prefix: "booking_system_",
  labels: { service: "unknown" },
});

export default {
  metricsMiddleware,
  grpcMetricsInterceptor,
  metricsEndpoint,
  healthCheckEndpoint,
  serviceMetrics,
  prometheus,
};

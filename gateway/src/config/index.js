import dotenv from 'dotenv';

dotenv.config();

const _DEFAULT_JWT_SECRET = 'your_jwt_secret_key';
const _DEFAULT_JWT_REFRESH_SECRET = 'your_refresh_secret_key';

if (process.env.NODE_ENV === 'production') {
  if (!process.env.JWT_SECRET || process.env.JWT_SECRET === _DEFAULT_JWT_SECRET) {
    throw new Error('[FATAL] JWT_SECRET must be explicitly set in production. Refusing to start.');
  }
  if (!process.env.JWT_REFRESH_SECRET || process.env.JWT_REFRESH_SECRET === _DEFAULT_JWT_REFRESH_SECRET) {
    throw new Error('[FATAL] JWT_REFRESH_SECRET must be explicitly set in production. Refusing to start.');
  }
} else {
  if (!process.env.JWT_SECRET) {
    console.warn('[WARN] JWT_SECRET not set — using insecure default. Set this before deploying to production.');
  }
  if (!process.env.JWT_REFRESH_SECRET) {
    console.warn('[WARN] JWT_REFRESH_SECRET not set — using insecure default. Set this before deploying to production.');
  }
}

const config = {
  // Server Configuration
  server: {
    port: process.env.PORT || 3000,
    nodeEnv: process.env.NODE_ENV || 'development',
    allowedOrigins: process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:53000'],
  },

  // Redis Configuration
  redis: {
    host: process.env.REDIS_HOST || 'localhost',
    port: process.env.REDIS_PORT || 6379,
    password: process.env.REDIS_PASSWORD,
    database: process.env.REDIS_DATABASE || 0,
    url:
      process.env.REDIS_URL ||
      `redis://${process.env.REDIS_HOST || 'localhost'}:${process.env.REDIS_PORT || 6379}`,
  },

  // JWT Configuration
  jwt: {
    secret: process.env.JWT_SECRET || 'your_jwt_secret_key',
    refreshSecret: process.env.JWT_REFRESH_SECRET || 'your_refresh_secret_key',
    expiresIn: process.env.JWT_EXPIRES_IN || '15m',
    refreshExpiresIn: process.env.JWT_REFRESH_EXPIRES_IN || '7d',
  },

  // Rate Limiting Configuration
  rateLimit: {
    windowMs: parseInt(process.env.RATE_LIMIT_WINDOW_MS) || 15 * 60 * 1000, // 15 minutes
    maxRequests: parseInt(process.env.RATE_LIMIT_MAX_REQUESTS) || 100,
    delayAfter: parseInt(process.env.RATE_LIMIT_DELAY_AFTER) || 50,
    delayMs: parseInt(process.env.RATE_LIMIT_DELAY_MS) || 500,
  },

  // gRPC Service Configuration
  grpc: {
    authService: {
      url: process.env.GRPC_AUTH_SERVICE_URL || '127.0.0.1:50051',
      maxReceiveMessageLength: parseInt(process.env.GRPC_MAX_RECEIVE_MESSAGE_LENGTH) || 4194304,
      maxSendMessageLength: parseInt(process.env.GRPC_MAX_SEND_MESSAGE_LENGTH) || 4194304,
      keepaliveTimeMs: parseInt(process.env.GRPC_KEEPALIVE_TIME_MS) || 30000,
      keepaliveTimeoutMs: parseInt(process.env.GRPC_KEEPALIVE_TIMEOUT_MS) || 5000,
    },
    userService: {
      url: process.env.GRPC_USER_SERVICE_URL || 'localhost:50052',
    },
    eventService: {
      url: process.env.GRPC_EVENT_SERVICE_URL || 'event-service:50055',
    },
    bookingService: {
      url: process.env.GRPC_BOOKING_SERVICE_URL || 'booking-service:50058',
    },
    paymentService: {
      url: process.env.GRPC_PAYMENT_SERVICE_URL || 'payment-service:50062',
    },
    ticketService: {
      url: process.env.GRPC_TICKET_SERVICE_URL || 'ticket-service:50053',
    },
  },

  // Circuit Breaker Configuration
  // Always enabled. Dev uses relaxed thresholds so failures trip the breaker
  // less aggressively, but the behaviour is still tested end-to-end.
  circuitBreaker: {
    enabled: process.env.CIRCUIT_BREAKER_ENABLED !== 'false',
    timeout: parseInt(process.env.CIRCUIT_BREAKER_TIMEOUT) || 30000,
    resetTimeout: parseInt(process.env.CIRCUIT_BREAKER_RESET_TIMEOUT) || 30000,
    errorThresholdPercentage:
      parseInt(process.env.CIRCUIT_BREAKER_ERROR_PERCENTAGE) ||
      (process.env.NODE_ENV === 'production' ? 50 : 75),
    volumeThreshold:
      parseInt(process.env.CIRCUIT_BREAKER_VOLUME_THRESHOLD) ||
      (process.env.NODE_ENV === 'production' ? 5 : 20),
  },

  // Logging Configuration
  logging: {
    level: process.env.LOG_LEVEL || 'info',
    format: process.env.LOG_FORMAT || 'json',
  },

  // Monitoring Configuration
  monitoring: {
    prometheus: {
      port: process.env.PROMETHEUS_PORT || 9090,
    },
  },

  // Security Configuration
  security: {
    cors: {
      origin: process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:53000'],
      credentials: true,
    },
    helmet: {
      contentSecurityPolicy: {
        directives: {
          defaultSrc: ["'self'"],
          styleSrc: ["'self'", "'unsafe-inline'"],
          scriptSrc: ["'self'"],
          imgSrc: ["'self'", 'data:', 'https:'],
        },
      },
    },
  },
};

export default config;

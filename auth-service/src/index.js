import * as grpc from '@grpc/grpc-js';
import dotenv from 'dotenv';
import logger from './utils/logger.js';
import { server } from './server.js';
import { initializeBackgroundService } from './background/registerJobHandlers.js';
import { healthCheck as grpcHealthCheck } from './grpc/clients.js';
import { markServiceAsStarted } from './controllers/healthController.js';
import express from 'express';
import promBundle from 'express-prom-bundle';

dotenv.config();

// Start server
const PORT = process.env.PORT || 50051;
const HOST = process.env.HOST || '0.0.0.0';

server.bindAsync(`${HOST}:${PORT}`, grpc.ServerCredentials.createInsecure(), async (err, port) => {
  if (err) {
    logger.error('Failed to bind server:', err);
    process.exit(1);
  }

  try {
    // Initialize background service with job handlers
    await initializeBackgroundService();
    logger.info('🔄 Background service initialized with job handlers');
  } catch (error) {
    logger.error('Failed to initialize background service:', error);
    // Don't exit, continue with server startup
  }

  logger.info(`🚀 Auth Service started on ${HOST}:${port}`);
  logger.info('📊 Master-Slave Database Pattern enabled');
  logger.info('🔐 JWT Authentication ready');
  logger.info('⚡ Functional Programming approach');
  logger.info('🔄 Background job processing enabled');

  // Mark service as started for startup probe
  markServiceAsStarted();

  // Check gRPC clients health in background (non-blocking)
  // This is for monitoring/observability only, not critical for startup
  grpcHealthCheck()
    .then((grpcHealth) => {
      logger.info('🔗 gRPC clients health check:', grpcHealth);
    })
    .catch((error) => {
      logger.warn('⚠️  gRPC clients health check failed (non-critical):', error.message);
    });
});

// Start metrics HTTP server using express-prom-bundle
const metricsApp = express();
const metricsMiddleware = promBundle({
  includeMethod: true,
  includePath: true,
  promClient: {
    collectDefaultMetrics: {},
  },
});
metricsApp.use(metricsMiddleware);
const METRICS_PORT = process.env.PROMETHEUS_PORT || 9190;
metricsApp.listen(METRICS_PORT, () => {
  logger.info(`📊 Metrics endpoint available at http://localhost:${METRICS_PORT}/metrics`);
});

export default server;

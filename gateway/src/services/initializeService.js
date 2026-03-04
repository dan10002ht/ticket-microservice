import express from 'express';
import logger from '../utils/logger.js';
import { requestContext } from '../utils/requestContext.js';
import {
  securityMiddleware,
  compressionMiddleware,
  bodyParsingMiddleware,
  loggingMiddleware,
  rateLimitMiddleware,
} from '../middlewares/index.js';
import { initializeSwagger } from './swaggerService.js';
import { initializeMetrics } from './metricsService.js';
import { initializeRoutes } from './routeService.js';
import { initializeErrorHandling } from './errorHandlingService.js';

/**
 * Initialize all gateway services and middleware
 * @param {express.Application} app - Express app instance
 */
export const initializeGateway = (app) => {
  const swaggerSpec = initializeSwagger();

  const metricsMiddleware = initializeMetrics();

  securityMiddleware(app);
  compressionMiddleware(app);
  bodyParsingMiddleware(app);
  loggingMiddleware(app);

  // Bind request-scoped context (correlationId flows into gRPC metadata automatically)
  app.use((req, _res, next) => {
    requestContext.run({ correlationId: req.correlationId || req.headers['x-correlation-id'] }, next);
  });

  const { limiter, speedLimiter } = rateLimitMiddleware();
  app.use(limiter);
  app.use(speedLimiter);

  app.use(metricsMiddleware);

  initializeRoutes(app, swaggerSpec);

  initializeErrorHandling(app);

  logger.info('Gateway initialization completed');
};

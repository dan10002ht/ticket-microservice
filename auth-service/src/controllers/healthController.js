/**
 * Health Check Endpoints following Kubernetes best practices
 *
 * - /health/startup:   Check if service has started (one-time check)
 * - /health/liveness:  Check if service is alive (critical resources only)
 * - /health/readiness: Check if service is ready to accept traffic (includes dependencies)
 */

import * as authService from '../services/internal/authService.js';
import { healthCheck as grpcHealthCheck } from '../grpc/clients.js';
import logger from '../utils/logger.js';

// Track if service has completed startup
let serviceStarted = false;

export function markServiceAsStarted() {
  serviceStarted = true;
  logger.info('✅ Service marked as started');
}

/**
 * Startup Probe - Has the service finished starting up?
 * This runs once during pod initialization
 * Returns 200 only after service initialization is complete
 */
export async function startup(call, callback) {
  try {
    if (!serviceStarted) {
      callback({
        code: 14, // UNAVAILABLE
        message: 'Service is still starting up',
      });
      return;
    }

    callback(null, {
      status: 'started',
      message: 'Service has started successfully',
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    logger.error('Startup probe error:', error);
    callback({
      code: 13, // INTERNAL
      message: error.message,
    });
  }
}

/**
 * Liveness Probe - Is the service alive?
 * Checks only critical resources (DB, cache)
 * If this fails, container should be restarted
 */
export async function liveness(call, callback) {
  try {
    // Check only critical resources - failures here mean service needs restart
    const result = await authService.healthCheck();

    // If database is down, service is not alive
    if (result.database?.status === 'unhealthy') {
      callback({
        code: 14, // UNAVAILABLE
        message: 'Critical resource unavailable: database',
      });
      return;
    }

    callback(null, {
      status: 'alive',
      message: 'Service is alive',
      timestamp: result.timestamp,
      critical: {
        database: result.database?.status || 'unknown',
        cache: result.cache?.status || 'unknown',
      },
    });
  } catch (error) {
    logger.error('Liveness probe error:', error);
    callback({
      code: 13, // INTERNAL
      message: error.message,
    });
  }
}

/**
 * Readiness Probe - Is the service ready to accept traffic?
 * Checks critical resources + optional dependencies
 * If this fails, traffic should be routed elsewhere (but container stays up)
 */
export async function readiness(call, callback) {
  try {
    // Check critical resources
    const serviceHealth = await authService.healthCheck();

    // Check if critical resources are healthy
    if (serviceHealth.status === 'unhealthy') {
      callback({
        code: 14, // UNAVAILABLE
        message: 'Service not ready: critical resources unhealthy',
        details: serviceHealth,
      });
      return;
    }

    // Check optional dependencies (non-blocking)
    let dependenciesHealth = { status: 'checking' };
    try {
      dependenciesHealth = await Promise.race([
        grpcHealthCheck(),
        new Promise((resolve) =>
          setTimeout(() => resolve({ status: 'timeout', message: 'Dependencies check timeout' }), 2000)
        ),
      ]);
    } catch (error) {
      logger.warn('Dependencies health check failed (non-critical):', error.message);
      dependenciesHealth = { status: 'unavailable', error: error.message };
    }

    callback(null, {
      status: 'ready',
      message: 'Service is ready to accept traffic',
      timestamp: serviceHealth.timestamp,
      critical: {
        database: serviceHealth.database?.status || 'unknown',
        cache: serviceHealth.cache?.status || 'unknown',
        redis: serviceHealth.redisPasswordReset?.status || 'unknown',
      },
      dependencies: dependenciesHealth,
    });
  } catch (error) {
    logger.error('Readiness probe error:', error);
    callback({
      code: 13, // INTERNAL
      message: error.message,
    });
  }
}

/**
 * Legacy health endpoint - combines all checks
 * Kept for backward compatibility
 */
export async function health(call, callback) {
  try {
    const result = await authService.healthCheck();

    callback(null, {
      status: result.status,
      message: result.status === 'healthy' ? 'Service is healthy' : 'Service is unhealthy',
      details: {
        timestamp: result.timestamp,
        database: result.database,
        service: result.service,
      },
    });
  } catch (error) {
    logger.error('Health check error:', error);
    callback({
      code: 13,
      message: error.message,
    });
  }
}

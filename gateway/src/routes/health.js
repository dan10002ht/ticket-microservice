import express from 'express';
import circuitBreakerService from '../services/circuitBreakerService.js';
import logger from '../utils/logger.js';

const router = express.Router();

const buildMeta = (correlationId) => ({
  correlationId,
  timestamp: new Date().toISOString(),
});

router.get('/', (req, res) => {
  res.status(200).json({
    data: {
      status: 'healthy',
      uptime: process.uptime(),
      version: process.env.npm_package_version || '1.0.0',
      service: 'gateway',
    },
    meta: buildMeta(req.correlationId),
  });
});

router.get('/ready', (req, res) => {
  res.status(200).json({
    data: {
      status: 'ready',
      service: 'gateway',
    },
    meta: buildMeta(req.correlationId),
  });
});

router.get('/live', (req, res) => {
  res.status(200).json({
    data: {
      status: 'alive',
      service: 'gateway',
    },
    meta: buildMeta(req.correlationId),
  });
});

// Circuit breaker stats endpoint
router.get('/circuit-breaker-stats', (req, res) => {
  try {
    const stats = circuitBreakerService.getStats();
    const health = circuitBreakerService.getHealth();

    res.json({
      data: { stats, health },
      meta: buildMeta(req.correlationId),
    });
  } catch (error) {
    logger.error('Failed to get circuit breaker stats', { error: error.message });
    res.status(500).json({
      error: { code: 'INTERNAL_ERROR', message: 'Failed to get circuit breaker stats' },
      meta: buildMeta(req.correlationId),
    });
  }
});

// Circuit breaker reset endpoint
router.post('/reset-circuit-breakers', (req, res) => {
  try {
    const { service } = req.body;

    if (service) {
      // Reset specific service breakers
      const breakers = circuitBreakerService.breakers;
      let resetCount = 0;

      for (const [name, breaker] of breakers) {
        if (name.startsWith(service + '.')) {
          breaker.close();
          resetCount++;
        }
      }

      logger.info(`Reset ${resetCount} circuit breakers for service: ${service}`);
      res.json({
        data: { message: `Reset ${resetCount} circuit breakers for service: ${service}` },
        meta: buildMeta(req.correlationId),
      });
    } else {
      // Reset all breakers
      circuitBreakerService.resetAllBreakers();
      res.json({
        data: { message: 'All circuit breakers reset' },
        meta: buildMeta(req.correlationId),
      });
    }
  } catch (error) {
    logger.error('Failed to reset circuit breakers', { error: error.message });
    res.status(500).json({
      error: { code: 'INTERNAL_ERROR', message: 'Failed to reset circuit breakers' },
      meta: buildMeta(req.correlationId),
    });
  }
});

export default router;

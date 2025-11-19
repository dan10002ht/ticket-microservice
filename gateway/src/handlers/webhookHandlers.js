import grpcClients from '../grpc/clients.js';
import logger from '../utils/logger.js';
import { createSimpleHandler } from '../utils/responseHandler.js';

const processPaymentWebhook = async (req, res) => {
  const gateway = req.params.gateway;

  const payload =
    req.rawBody ??
    (typeof req.body === 'string'
      ? req.body
      : req.body && Object.keys(req.body).length > 0
        ? JSON.stringify(req.body)
        : '');

  const headers = {};
  Object.entries(req.headers).forEach(([key, value]) => {
    if (Array.isArray(value)) {
      headers[key] = value.join(',');
    } else if (value !== undefined) {
      headers[key] = String(value);
    }
  });

  logger.info('Forwarding payment webhook to gRPC service', {
    gateway,
    correlationId: req.correlationId,
  });

  const response = await grpcClients.paymentService.processWebhook({
    gateway,
    payload,
    headers,
    correlationId: req.correlationId,
  });

  res.status(200).json({
    success: response.success,
    message: response.message,
    correlationId: req.correlationId,
  });
};

export const processPaymentWebhookHandler = createSimpleHandler(
  processPaymentWebhook,
  'payment',
  'processWebhook'
);

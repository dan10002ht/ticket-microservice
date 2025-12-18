import express from 'express';
import { body } from 'express-validator';

// Import handlers
import {
  processPaymentHandler,
  getUserPaymentsHandler,
  getPaymentHandler,
  refundPaymentHandler,
  getPaymentMethodsHandler,
  addPaymentMethodHandler,
  capturePaymentHandler,
  cancelPaymentHandler,
  listPaymentsHandler,
  listRefundsHandler,
  updateRefundStatusHandler,
} from '../handlers/paymentHandlers.js';
import { requireRole, requireAuth } from '../middlewares/index.js';

const router = express.Router();

// ============================================
// Admin routes (must be before /:paymentId)
// ============================================
router.get('/admin/list', requireRole(['admin']), listPaymentsHandler);

// ============================================
// Refund routes (must be before /:paymentId)
// ============================================
router.put('/refunds/:refundId', requireRole(['admin']), updateRefundStatusHandler);

// ============================================
// Payment Methods
// ============================================
router.get('/methods', requireAuth, getPaymentMethodsHandler);
router.post('/methods', requireAuth, addPaymentMethodHandler);

// ============================================
// Payment CRUD
// ============================================
router.post(
  '/',
  requireAuth,
  [
    body('booking_id').notEmpty().trim(),
    body('amount').isFloat({ min: 0.01 }),
    body('payment_method').notEmpty().trim(),
  ],
  processPaymentHandler
);

router.get('/', requireAuth, getUserPaymentsHandler);

router.get('/:paymentId', requireAuth, getPaymentHandler);

// ============================================
// Payment Operations
// ============================================
router.post('/:paymentId/capture', requireAuth, capturePaymentHandler);

router.post(
  '/:paymentId/cancel',
  requireAuth,
  [body('reason').optional().trim()],
  cancelPaymentHandler
);

// ============================================
// Refund Management
// ============================================
router.post(
  '/:paymentId/refund',
  requireAuth,
  [
    body('amount').isFloat({ min: 0.01 }),
    body('reason').optional().trim(),
  ],
  refundPaymentHandler
);

router.get('/:paymentId/refunds', requireAuth, listRefundsHandler);

export default router;

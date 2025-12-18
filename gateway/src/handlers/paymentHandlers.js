import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createHandler, createSimpleHandler } from '../utils/responseHandler.js';

// ============================================
// Payment CRUD
// ============================================

/**
 * Create/Process payment
 */
const processUserPayment = async (req, res) => {
  const result = await grpcClients.paymentService.CreatePayment({
    booking_id: req.body.booking_id,
    ticket_id: req.body.ticket_id || '',
    user_id: req.user.id,
    amount: req.body.amount,
    currency: req.body.currency || 'USD',
    payment_method: req.body.payment_method,
    gateway_provider: req.body.gateway_provider || 'stripe',
    idempotency_key: req.body.idempotency_key || '',
    metadata: req.body.metadata || {},
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get payment by ID
 */
const getPaymentById = async (req, res) => {
  const result = await grpcClients.paymentService.GetPayment({
    payment_id: req.params.paymentId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get user payments
 */
const getUserPayments = async (req, res) => {
  const result = await grpcClients.paymentService.ListPayments({
    user_id: req.user.id,
    status: req.query.status || '',
    page: parseInt(req.query.page) || 1,
    size: parseInt(req.query.size) || 20,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get payment methods
 */
const getPaymentMethods = async (req, res) => {
  // Note: This may need a separate gRPC method if available
  // For now returning static payment methods
  const methods = [
    { id: 'card', name: 'Credit/Debit Card', enabled: true },
    { id: 'bank_transfer', name: 'Bank Transfer', enabled: true },
    { id: 'wallet', name: 'Digital Wallet', enabled: true },
  ];
  sendSuccessResponse(res, 200, { methods }, req.correlationId);
};

/**
 * Add payment method
 */
const addPaymentMethod = async (req, res) => {
  // Note: This may need a separate gRPC method if available
  sendSuccessResponse(res, 201, { message: 'Payment method added' }, req.correlationId);
};

// ============================================
// Payment Operations
// ============================================

/**
 * Capture pre-authorized payment
 */
const capturePayment = async (req, res) => {
  const result = await grpcClients.paymentService.CapturePayment({
    payment_id: req.params.paymentId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Cancel payment
 */
const cancelPayment = async (req, res) => {
  const result = await grpcClients.paymentService.CancelPayment({
    payment_id: req.params.paymentId,
    reason: req.body.reason || '',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * List all payments (Admin)
 */
const listPayments = async (req, res) => {
  const result = await grpcClients.paymentService.ListPayments({
    user_id: req.query.user_id || '',
    status: req.query.status || '',
    page: parseInt(req.query.page) || 1,
    size: parseInt(req.query.size) || 20,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Refund Management
// ============================================

/**
 * Create refund for payment
 */
const createRefund = async (req, res) => {
  const result = await grpcClients.paymentService.CreateRefund({
    payment_id: req.params.paymentId,
    amount: req.body.amount,
    reason: req.body.reason || '',
    description: req.body.description || '',
    refund_type: req.body.refund_type || 'full',
    idempotency_key: req.body.idempotency_key || '',
    metadata: req.body.metadata || {},
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * List refunds for a payment
 */
const listRefunds = async (req, res) => {
  const result = await grpcClients.paymentService.ListRefunds({
    payment_id: req.params.paymentId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update refund status (Admin)
 */
const updateRefundStatus = async (req, res) => {
  const result = await grpcClients.paymentService.UpdateRefundStatus({
    refund_id: req.params.refundId,
    status: req.body.status,
    external_reference: req.body.external_reference || '',
    failure_reason: req.body.failure_reason || '',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Export handlers
// ============================================

// Payment CRUD
export const processPaymentHandler = createHandler(processUserPayment, 'payment', 'processPayment');
export const getPaymentHandler = createSimpleHandler(getPaymentById, 'payment', 'getPayment');
export const getUserPaymentsHandler = createSimpleHandler(getUserPayments, 'payment', 'getUserPayments');
export const getPaymentMethodsHandler = createSimpleHandler(getPaymentMethods, 'payment', 'getPaymentMethods');
export const addPaymentMethodHandler = createHandler(addPaymentMethod, 'payment', 'addPaymentMethod');

// Payment Operations
export const capturePaymentHandler = createHandler(capturePayment, 'payment', 'capturePayment');
export const cancelPaymentHandler = createHandler(cancelPayment, 'payment', 'cancelPayment');
export const listPaymentsHandler = createSimpleHandler(listPayments, 'payment', 'listPayments');

// Refund Management
export const refundPaymentHandler = createHandler(createRefund, 'payment', 'createRefund');
export const listRefundsHandler = createSimpleHandler(listRefunds, 'payment', 'listRefunds');
export const updateRefundStatusHandler = createHandler(updateRefundStatus, 'payment', 'updateRefundStatus'); 
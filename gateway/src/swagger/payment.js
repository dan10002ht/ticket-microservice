/**
 * @swagger
 * components:
 *   schemas:
 *     Payment:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: Payment ID
 *         booking_id:
 *           type: string
 *           description: Associated booking ID
 *         ticket_id:
 *           type: string
 *           description: Associated ticket ID
 *         user_id:
 *           type: string
 *           description: User ID
 *         amount:
 *           type: number
 *           description: Payment amount
 *         currency:
 *           type: string
 *           default: USD
 *           description: Payment currency
 *         payment_method:
 *           type: string
 *           enum: [credit_card, debit_card, bank_transfer, digital_wallet]
 *           description: Payment method used
 *         gateway_provider:
 *           type: string
 *           enum: [stripe, paypal, square]
 *           description: Payment gateway provider
 *         status:
 *           type: string
 *           enum: [pending, authorized, captured, completed, failed, refunded, cancelled]
 *           description: Payment status
 *         transaction_id:
 *           type: string
 *           description: External transaction ID
 *         idempotency_key:
 *           type: string
 *           description: Idempotency key
 *         metadata:
 *           type: object
 *           description: Additional payment metadata
 *         created_at:
 *           type: string
 *           format: date-time
 *           description: Payment creation date
 *         updated_at:
 *           type: string
 *           format: date-time
 *           description: Last update date
 *     PaymentCreate:
 *       type: object
 *       required:
 *         - booking_id
 *         - amount
 *         - payment_method
 *       properties:
 *         booking_id:
 *           type: string
 *           description: Booking ID to pay for
 *         ticket_id:
 *           type: string
 *           description: Ticket ID (optional)
 *         amount:
 *           type: number
 *           minimum: 0.01
 *           description: Payment amount
 *         currency:
 *           type: string
 *           default: USD
 *           description: Currency code
 *         payment_method:
 *           type: string
 *           enum: [credit_card, debit_card, bank_transfer, digital_wallet]
 *           description: Payment method to use
 *         gateway_provider:
 *           type: string
 *           enum: [stripe, paypal, square]
 *           default: stripe
 *           description: Payment gateway provider
 *         idempotency_key:
 *           type: string
 *           description: Idempotency key to prevent duplicate payments
 *         metadata:
 *           type: object
 *           description: Additional metadata
 *     PaymentMethod:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: Payment method ID
 *         name:
 *           type: string
 *           description: Payment method name
 *         enabled:
 *           type: boolean
 *           description: Whether this method is enabled
 *     Refund:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: Refund ID
 *         payment_id:
 *           type: string
 *           description: Original payment ID
 *         amount:
 *           type: number
 *           description: Refund amount
 *         currency:
 *           type: string
 *           description: Currency code
 *         reason:
 *           type: string
 *           description: Reason for refund
 *         description:
 *           type: string
 *           description: Refund description
 *         refund_type:
 *           type: string
 *           enum: [full, partial]
 *           description: Type of refund
 *         status:
 *           type: string
 *           enum: [pending, processing, completed, failed]
 *           description: Refund status
 *         external_reference:
 *           type: string
 *           description: External refund reference
 *         failure_reason:
 *           type: string
 *           description: Reason for failure (if failed)
 *         metadata:
 *           type: object
 *           description: Additional refund metadata
 *         created_at:
 *           type: string
 *           format: date-time
 *         updated_at:
 *           type: string
 *           format: date-time
 *     RefundCreate:
 *       type: object
 *       properties:
 *         amount:
 *           type: number
 *           description: Refund amount (optional for full refund)
 *         reason:
 *           type: string
 *           description: Reason for refund
 *         description:
 *           type: string
 *           description: Refund description
 *         refund_type:
 *           type: string
 *           enum: [full, partial]
 *           default: full
 *           description: Type of refund
 *         idempotency_key:
 *           type: string
 *           description: Idempotency key
 *         metadata:
 *           type: object
 *           description: Additional metadata
 */

/**
 * @swagger
 * /api/payments:
 *   post:
 *     summary: Process payment
 *     description: Process a payment for a booking
 *     tags: [Payments]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/PaymentCreate'
 *     responses:
 *       200:
 *         description: Payment processed successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 payment:
 *                   $ref: '#/components/schemas/Payment'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       402:
 *         description: Payment failed
 *       404:
 *         description: Booking not found
 *   get:
 *     summary: Get user payment history
 *     description: Retrieve payment history for the current user
 *     tags: [Payments]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [pending, authorized, captured, completed, failed, refunded, cancelled]
 *         description: Filter by payment status
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number for pagination
 *       - in: query
 *         name: size
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 20
 *         description: Number of items per page
 *     responses:
 *       200:
 *         description: Payment history retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 payments:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Payment'
 *                 total:
 *                   type: integer
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/payments/admin/list:
 *   get:
 *     summary: List all payments (Admin)
 *     description: Retrieve all payments with optional filters (Admin only)
 *     tags: [Payments Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: user_id
 *         schema:
 *           type: string
 *         description: Filter by user ID
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [pending, authorized, captured, completed, failed, refunded, cancelled]
 *         description: Filter by payment status
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number for pagination
 *       - in: query
 *         name: size
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 20
 *         description: Number of items per page
 *     responses:
 *       200:
 *         description: Payments list retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 payments:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Payment'
 *                 total:
 *                   type: integer
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 */

/**
 * @swagger
 * /api/payments/methods:
 *   get:
 *     summary: Get available payment methods
 *     description: Retrieve available payment methods
 *     tags: [Payments]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Payment methods retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 methods:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/PaymentMethod'
 *       401:
 *         description: Unauthorized
 *   post:
 *     summary: Add payment method
 *     description: Add a new payment method for the user
 *     tags: [Payments]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - type
 *               - name
 *             properties:
 *               type:
 *                 type: string
 *                 enum: [credit_card, debit_card, bank_account, digital_wallet]
 *               name:
 *                 type: string
 *               is_default:
 *                 type: boolean
 *     responses:
 *       201:
 *         description: Payment method added successfully
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/payments/{paymentId}:
 *   get:
 *     summary: Get payment by ID
 *     description: Retrieve a specific payment by ID
 *     tags: [Payments]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: paymentId
 *         required: true
 *         schema:
 *           type: string
 *         description: Payment ID
 *     responses:
 *       200:
 *         description: Payment retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 payment:
 *                   $ref: '#/components/schemas/Payment'
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Payment not found
 */

/**
 * @swagger
 * /api/payments/{paymentId}/capture:
 *   post:
 *     summary: Capture pre-authorized payment
 *     description: Capture a previously authorized payment
 *     tags: [Payment Operations]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: paymentId
 *         required: true
 *         schema:
 *           type: string
 *         description: Payment ID
 *     responses:
 *       200:
 *         description: Payment captured successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 payment:
 *                   $ref: '#/components/schemas/Payment'
 *                 message:
 *                   type: string
 *       400:
 *         description: Payment cannot be captured
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Payment not found
 */

/**
 * @swagger
 * /api/payments/{paymentId}/cancel:
 *   post:
 *     summary: Cancel payment
 *     description: Cancel a pending or authorized payment
 *     tags: [Payment Operations]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: paymentId
 *         required: true
 *         schema:
 *           type: string
 *         description: Payment ID
 *     requestBody:
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               reason:
 *                 type: string
 *                 description: Reason for cancellation
 *     responses:
 *       200:
 *         description: Payment cancelled successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 payment:
 *                   $ref: '#/components/schemas/Payment'
 *                 message:
 *                   type: string
 *       400:
 *         description: Payment cannot be cancelled
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Payment not found
 */

/**
 * @swagger
 * /api/payments/{paymentId}/refund:
 *   post:
 *     summary: Create refund
 *     description: Create a refund for a completed payment
 *     tags: [Refunds]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: paymentId
 *         required: true
 *         schema:
 *           type: string
 *         description: Payment ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/RefundCreate'
 *     responses:
 *       200:
 *         description: Refund created successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 refund:
 *                   $ref: '#/components/schemas/Refund'
 *       400:
 *         description: Payment cannot be refunded
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Payment not found
 */

/**
 * @swagger
 * /api/payments/{paymentId}/refunds:
 *   get:
 *     summary: List refunds for payment
 *     description: Get all refunds associated with a payment
 *     tags: [Refunds]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: paymentId
 *         required: true
 *         schema:
 *           type: string
 *         description: Payment ID
 *     responses:
 *       200:
 *         description: Refunds retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 refunds:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Refund'
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Payment not found
 */

/**
 * @swagger
 * /api/payments/refunds/{refundId}:
 *   put:
 *     summary: Update refund status (Admin)
 *     description: Update the status of a refund (Admin only)
 *     tags: [Refunds Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: refundId
 *         required: true
 *         schema:
 *           type: string
 *         description: Refund ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - status
 *             properties:
 *               status:
 *                 type: string
 *                 enum: [pending, processing, completed, failed]
 *                 description: New refund status
 *               external_reference:
 *                 type: string
 *                 description: External refund reference
 *               failure_reason:
 *                 type: string
 *                 description: Reason for failure (if status is failed)
 *     responses:
 *       200:
 *         description: Refund status updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 refund:
 *                   $ref: '#/components/schemas/Refund'
 *       400:
 *         description: Invalid status transition
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 *       404:
 *         description: Refund not found
 */

export default {};

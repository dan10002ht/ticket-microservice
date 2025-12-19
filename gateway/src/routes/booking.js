import express from 'express';
import { body } from 'express-validator';

// Import handlers
import {
  createBookingHandler,
  getUserBookingsHandler,
  getBookingHandler,
  cancelBookingHandler,
  updateBookingHandler,
  confirmBookingHandler,
  listBookingsHandler,
  reserveSeatsHandler,
  releaseSeatsHandler,
} from '../handlers/bookingHandlers.js';
import { requireRole, requireAuth, validateUUIDParam } from '../middlewares/index.js';

const router = express.Router();

// ============================================
// Admin routes (must be before /:bookingId)
// ============================================
router.get('/admin/list', requireRole(['admin']), listBookingsHandler);

// ============================================
// Seat Reservation
// ============================================
router.post(
  '/seats/reserve',
  requireAuth,
  [
    body('event_id').notEmpty().trim(),
    body('seat_numbers').isArray({ min: 1 }),
  ],
  reserveSeatsHandler
);

router.post(
  '/seats/release',
  requireAuth,
  [
    body('reservation_id').notEmpty().trim(),
  ],
  releaseSeatsHandler
);

// ============================================
// Booking CRUD
// ============================================
router.post(
  '/',
  requireAuth,
  [
    body('event_id').notEmpty().trim(),
    body('ticket_quantity').isInt({ min: 1 }),
    body('special_requests').optional().trim(),
  ],
  createBookingHandler
);

router.get('/', requireAuth, getUserBookingsHandler);

router.get('/:bookingId', requireAuth, validateUUIDParam('bookingId'), getBookingHandler);

router.put('/:bookingId', requireAuth, validateUUIDParam('bookingId'), updateBookingHandler);

router.post('/:bookingId/cancel', requireAuth, validateUUIDParam('bookingId'), cancelBookingHandler);

router.post(
  '/:bookingId/confirm',
  requireAuth,
  validateUUIDParam('bookingId'),
  [
    body('payment_reference').notEmpty().trim(),
  ],
  confirmBookingHandler
);

export default router;

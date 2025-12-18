import grpcClients from '../grpc/clients.js';
import {
  sendSuccessResponse,
  createHandler,
  createSimpleHandler,
} from '../utils/responseHandler.js';

// ============================================
// Booking CRUD
// ============================================

/**
 * Create a new booking
 */
const createNewBooking = async (req, res) => {
  const result = await grpcClients.bookingService.CreateBooking({
    user_id: req.user.id,
    event_id: req.body.event_id,
    ticket_quantity: req.body.ticket_quantity,
    seat_numbers: req.body.seat_numbers || [],
    special_requests: req.body.special_requests || '',
    metadata: req.body.metadata || {},
    idempotency_key: req.body.idempotency_key || '',
  });
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Get booking by ID
 */
const getBookingById = async (req, res) => {
  const result = await grpcClients.bookingService.GetBooking({
    booking_id: req.params.bookingId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get user bookings
 */
const getUserBookings = async (req, res) => {
  const result = await grpcClients.bookingService.GetUserBookings({
    user_id: req.user.id,
    page: parseInt(req.query.page) || 1,
    limit: parseInt(req.query.limit) || 20,
    status: req.query.status || '',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Cancel booking
 */
const cancelUserBooking = async (req, res) => {
  const result = await grpcClients.bookingService.CancelBooking({
    booking_id: req.params.bookingId,
    reason: req.body.reason || '',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update booking
 */
const updateUserBooking = async (req, res) => {
  const result = await grpcClients.bookingService.UpdateBooking({
    booking_id: req.params.bookingId,
    ticket_quantity: req.body.ticket_quantity,
    seat_numbers: req.body.seat_numbers,
    special_requests: req.body.special_requests,
    metadata: req.body.metadata,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Booking Operations
// ============================================

/**
 * Confirm booking (after payment)
 */
const confirmBooking = async (req, res) => {
  const result = await grpcClients.bookingService.ConfirmBooking({
    booking_id: req.params.bookingId,
    payment_reference: req.body.payment_reference,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * List all bookings (Admin)
 */
const listBookings = async (req, res) => {
  const result = await grpcClients.bookingService.ListBookings({
    page: parseInt(req.query.page) || 1,
    limit: parseInt(req.query.limit) || 20,
    status: req.query.status || '',
    event_id: req.query.event_id || '',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Seat Reservation
// ============================================

/**
 * Reserve seats for booking
 */
const reserveSeats = async (req, res) => {
  const result = await grpcClients.bookingService.ReserveSeats({
    event_id: req.body.event_id,
    seat_numbers: req.body.seat_numbers,
    user_id: req.user.id,
    timeout_seconds: req.body.timeout_seconds || 600,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Release reserved seats
 */
const releaseSeats = async (req, res) => {
  const result = await grpcClients.bookingService.ReleaseSeats({
    reservation_id: req.body.reservation_id,
    seat_numbers: req.body.seat_numbers || [],
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Export handlers
// ============================================

// Booking CRUD
export const createBookingHandler = createHandler(createNewBooking, 'booking', 'createBooking');
export const getBookingHandler = createSimpleHandler(getBookingById, 'booking', 'getBooking');
export const getUserBookingsHandler = createSimpleHandler(getUserBookings, 'booking', 'getUserBookings');
export const cancelBookingHandler = createHandler(cancelUserBooking, 'booking', 'cancelBooking');
export const updateBookingHandler = createHandler(updateUserBooking, 'booking', 'updateBooking');

// Booking Operations
export const confirmBookingHandler = createHandler(confirmBooking, 'booking', 'confirmBooking');
export const listBookingsHandler = createSimpleHandler(listBookings, 'booking', 'listBookings');

// Seat Reservation
export const reserveSeatsHandler = createHandler(reserveSeats, 'booking', 'reserveSeats');
export const releaseSeatsHandler = createHandler(releaseSeats, 'booking', 'releaseSeats');

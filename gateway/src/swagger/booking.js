/**
 * @swagger
 * components:
 *   schemas:
 *     Booking:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: Booking ID
 *         user_id:
 *           type: string
 *           description: User ID
 *         event_id:
 *           type: string
 *           description: Event ID
 *         ticket_quantity:
 *           type: integer
 *           minimum: 1
 *           description: Number of tickets
 *         seat_numbers:
 *           type: array
 *           items:
 *             type: string
 *           description: Reserved seat numbers
 *         total_amount:
 *           type: number
 *           description: Total booking amount
 *         currency:
 *           type: string
 *           description: Currency code
 *         status:
 *           type: string
 *           enum: [pending, confirmed, cancelled, completed, expired]
 *           description: Booking status
 *         special_requests:
 *           type: string
 *           description: Special requests for the booking
 *         payment_reference:
 *           type: string
 *           description: Payment reference ID
 *         metadata:
 *           type: object
 *           description: Additional booking metadata
 *         created_at:
 *           type: string
 *           format: date-time
 *           description: Booking creation date
 *         updated_at:
 *           type: string
 *           format: date-time
 *           description: Last update date
 *         expires_at:
 *           type: string
 *           format: date-time
 *           description: Booking expiration time
 *     BookingCreate:
 *       type: object
 *       required:
 *         - event_id
 *         - ticket_quantity
 *       properties:
 *         event_id:
 *           type: string
 *           description: Event ID
 *         ticket_quantity:
 *           type: integer
 *           minimum: 1
 *           description: Number of tickets
 *         seat_numbers:
 *           type: array
 *           items:
 *             type: string
 *           description: Specific seats to book
 *         special_requests:
 *           type: string
 *           description: Special requests for the booking
 *         idempotency_key:
 *           type: string
 *           description: Idempotency key to prevent duplicate bookings
 *         metadata:
 *           type: object
 *           description: Additional metadata
 *     BookingUpdate:
 *       type: object
 *       properties:
 *         ticket_quantity:
 *           type: integer
 *           minimum: 1
 *           description: Updated number of tickets
 *         seat_numbers:
 *           type: array
 *           items:
 *             type: string
 *           description: Updated seat numbers
 *         special_requests:
 *           type: string
 *           description: Updated special requests
 *         metadata:
 *           type: object
 *           description: Updated metadata
 *     SeatReservation:
 *       type: object
 *       properties:
 *         reservation_id:
 *           type: string
 *           description: Reservation ID
 *         event_id:
 *           type: string
 *           description: Event ID
 *         seat_numbers:
 *           type: array
 *           items:
 *             type: string
 *           description: Reserved seat numbers
 *         user_id:
 *           type: string
 *           description: User ID
 *         expires_at:
 *           type: string
 *           format: date-time
 *           description: Reservation expiration time
 *     SeatReservationRequest:
 *       type: object
 *       required:
 *         - event_id
 *         - seat_numbers
 *       properties:
 *         event_id:
 *           type: string
 *           description: Event ID
 *         seat_numbers:
 *           type: array
 *           items:
 *             type: string
 *           description: Seat numbers to reserve
 *         timeout_seconds:
 *           type: integer
 *           default: 600
 *           description: Reservation timeout in seconds
 */

/**
 * @swagger
 * /api/bookings:
 *   post:
 *     summary: Create a new booking
 *     description: Create a new booking for an event
 *     tags: [Bookings]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/BookingCreate'
 *     responses:
 *       201:
 *         description: Booking created successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 booking:
 *                   $ref: '#/components/schemas/Booking'
 *       400:
 *         description: Validation error or insufficient tickets
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Event not found
 *   get:
 *     summary: Get user bookings
 *     description: Retrieve all bookings for the current user
 *     tags: [Bookings]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [pending, confirmed, cancelled, completed, expired]
 *         description: Filter by booking status
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number for pagination
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 20
 *         description: Number of items per page
 *     responses:
 *       200:
 *         description: Bookings retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 bookings:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Booking'
 *                 total:
 *                   type: integer
 *                 page:
 *                   type: integer
 *                 limit:
 *                   type: integer
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/bookings/admin/list:
 *   get:
 *     summary: List all bookings (Admin)
 *     description: Retrieve all bookings with optional filters (Admin only)
 *     tags: [Bookings Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [pending, confirmed, cancelled, completed, expired]
 *         description: Filter by booking status
 *       - in: query
 *         name: event_id
 *         schema:
 *           type: string
 *         description: Filter by event ID
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number for pagination
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 20
 *         description: Number of items per page
 *     responses:
 *       200:
 *         description: Bookings list retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 bookings:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Booking'
 *                 total:
 *                   type: integer
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 */

/**
 * @swagger
 * /api/bookings/seats/reserve:
 *   post:
 *     summary: Reserve seats
 *     description: Temporarily reserve seats for an event before booking
 *     tags: [Seat Reservation]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/SeatReservationRequest'
 *     responses:
 *       200:
 *         description: Seats reserved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 reservation:
 *                   $ref: '#/components/schemas/SeatReservation'
 *                 message:
 *                   type: string
 *       400:
 *         description: Seats not available
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/bookings/seats/release:
 *   post:
 *     summary: Release reserved seats
 *     description: Release previously reserved seats back to available pool
 *     tags: [Seat Reservation]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - reservation_id
 *             properties:
 *               reservation_id:
 *                 type: string
 *                 description: Reservation ID to release
 *               seat_numbers:
 *                 type: array
 *                 items:
 *                   type: string
 *                 description: Specific seats to release (optional, releases all if empty)
 *     responses:
 *       200:
 *         description: Seats released successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 message:
 *                   type: string
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Reservation not found
 */

/**
 * @swagger
 * /api/bookings/{bookingId}:
 *   get:
 *     summary: Get booking by ID
 *     description: Retrieve a specific booking by ID
 *     tags: [Bookings]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: bookingId
 *         required: true
 *         schema:
 *           type: string
 *         description: Booking ID
 *     responses:
 *       200:
 *         description: Booking retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 booking:
 *                   $ref: '#/components/schemas/Booking'
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Booking not found
 *   put:
 *     summary: Update booking
 *     description: Update a booking (before confirmation)
 *     tags: [Bookings]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: bookingId
 *         required: true
 *         schema:
 *           type: string
 *         description: Booking ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/BookingUpdate'
 *     responses:
 *       200:
 *         description: Booking updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 booking:
 *                   $ref: '#/components/schemas/Booking'
 *       400:
 *         description: Booking cannot be updated
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Booking not found
 */

/**
 * @swagger
 * /api/bookings/{bookingId}/cancel:
 *   post:
 *     summary: Cancel booking
 *     description: Cancel a booking
 *     tags: [Bookings]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: bookingId
 *         required: true
 *         schema:
 *           type: string
 *         description: Booking ID
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
 *         description: Booking cancelled successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 booking:
 *                   $ref: '#/components/schemas/Booking'
 *                 message:
 *                   type: string
 *       400:
 *         description: Booking cannot be cancelled
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Booking not found
 */

/**
 * @swagger
 * /api/bookings/{bookingId}/confirm:
 *   post:
 *     summary: Confirm booking
 *     description: Confirm a booking after successful payment
 *     tags: [Bookings]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: bookingId
 *         required: true
 *         schema:
 *           type: string
 *         description: Booking ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - payment_reference
 *             properties:
 *               payment_reference:
 *                 type: string
 *                 description: Payment reference/transaction ID
 *     responses:
 *       200:
 *         description: Booking confirmed successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 booking:
 *                   $ref: '#/components/schemas/Booking'
 *                 message:
 *                   type: string
 *       400:
 *         description: Booking cannot be confirmed
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Booking not found
 */

export default {};

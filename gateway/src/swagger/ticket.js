/**
 * @swagger
 * components:
 *   schemas:
 *     Ticket:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: Ticket ID
 *         event_id:
 *           type: string
 *           description: Event ID
 *         ticket_type_id:
 *           type: string
 *           description: Ticket type ID
 *         user_id:
 *           type: string
 *           description: User ID
 *         booking_id:
 *           type: string
 *           description: Booking ID
 *         ticket_number:
 *           type: string
 *           description: Unique ticket number
 *         seat_number:
 *           type: string
 *           description: Seat number
 *         price:
 *           type: number
 *           description: Ticket price
 *         currency:
 *           type: string
 *           description: Currency code
 *         status:
 *           type: string
 *           enum: [available, reserved, sold, used, cancelled]
 *           description: Ticket status
 *         qr_code:
 *           type: string
 *           description: QR code for ticket validation
 *         valid_from:
 *           type: string
 *           format: date-time
 *         valid_until:
 *           type: string
 *           format: date-time
 *         created_at:
 *           type: string
 *           format: date-time
 *         updated_at:
 *           type: string
 *           format: date-time
 *     TicketType:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: Ticket type ID
 *         event_id:
 *           type: string
 *           description: Event ID
 *         name:
 *           type: string
 *           description: Ticket type name (e.g., VIP, Standard)
 *         description:
 *           type: string
 *           description: Description
 *         price:
 *           type: number
 *           description: Base price
 *         currency:
 *           type: string
 *           description: Currency code
 *         quantity:
 *           type: integer
 *           description: Total quantity available
 *         available_quantity:
 *           type: integer
 *           description: Remaining quantity
 *         max_per_purchase:
 *           type: integer
 *           description: Maximum tickets per purchase
 *         min_per_purchase:
 *           type: integer
 *           description: Minimum tickets per purchase
 *         valid_from:
 *           type: string
 *           format: date-time
 *         valid_until:
 *           type: string
 *           format: date-time
 *         status:
 *           type: string
 *           enum: [active, inactive, sold_out]
 *     TicketTypeCreate:
 *       type: object
 *       required:
 *         - event_id
 *         - name
 *         - price
 *         - quantity
 *       properties:
 *         event_id:
 *           type: string
 *         name:
 *           type: string
 *         description:
 *           type: string
 *         price:
 *           type: number
 *         currency:
 *           type: string
 *           default: USD
 *         quantity:
 *           type: integer
 *         max_per_purchase:
 *           type: integer
 *           default: 10
 *         min_per_purchase:
 *           type: integer
 *           default: 1
 *     ReservationResponse:
 *       type: object
 *       properties:
 *         success:
 *           type: boolean
 *         tickets:
 *           type: array
 *           items:
 *             $ref: '#/components/schemas/Ticket'
 *         reservation_id:
 *           type: string
 *         expires_at:
 *           type: string
 *           format: date-time
 *         message:
 *           type: string
 *     AvailabilityResponse:
 *       type: object
 *       properties:
 *         success:
 *           type: boolean
 *         available:
 *           type: boolean
 *         available_quantity:
 *           type: integer
 *         available_seats:
 *           type: array
 *           items:
 *             type: string
 *         message:
 *           type: string
 */

/**
 * @swagger
 * /api/tickets:
 *   get:
 *     summary: Get all tickets
 *     tags: [Tickets]
 *     parameters:
 *       - in: query
 *         name: event_id
 *         schema:
 *           type: string
 *         description: Filter by event ID
 *       - in: query
 *         name: ticket_type_id
 *         schema:
 *           type: string
 *         description: Filter by ticket type
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *         description: Filter by status
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 50
 *     responses:
 *       200:
 *         description: List of tickets
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 tickets:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Ticket'
 *                 total:
 *                   type: integer
 *   post:
 *     summary: Create a new ticket
 *     tags: [Tickets]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/Ticket'
 *     responses:
 *       201:
 *         description: Ticket created
 */

/**
 * @swagger
 * /api/tickets/{ticketId}:
 *   get:
 *     summary: Get ticket by ID
 *     tags: [Tickets]
 *     parameters:
 *       - in: path
 *         name: ticketId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Ticket detail
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Ticket'
 *   put:
 *     summary: Update ticket
 *     tags: [Tickets]
 *     parameters:
 *       - in: path
 *         name: ticketId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/Ticket'
 *     responses:
 *       200:
 *         description: Ticket updated
 *   delete:
 *     summary: Delete ticket
 *     tags: [Tickets]
 *     parameters:
 *       - in: path
 *         name: ticketId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Ticket deleted
 */

/**
 * @swagger
 * /api/tickets/types/{eventId}:
 *   get:
 *     summary: Get ticket types for an event
 *     tags: [Ticket Types]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *         description: Event ID
 *       - in: query
 *         name: include_availability
 *         schema:
 *           type: boolean
 *         description: Include availability info
 *     responses:
 *       200:
 *         description: List of ticket types
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 ticket_types:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/TicketType'
 */

/**
 * @swagger
 * /api/tickets/types:
 *   post:
 *     summary: Create a new ticket type
 *     tags: [Ticket Types]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/TicketTypeCreate'
 *     responses:
 *       201:
 *         description: Ticket type created
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/TicketType'
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Organization role required
 */

/**
 * @swagger
 * /api/tickets/types/{typeId}:
 *   put:
 *     summary: Update ticket type
 *     tags: [Ticket Types]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: typeId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/TicketTypeCreate'
 *     responses:
 *       200:
 *         description: Ticket type updated
 *   delete:
 *     summary: Delete ticket type
 *     tags: [Ticket Types]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: typeId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Ticket type deleted
 */

/**
 * @swagger
 * /api/tickets/availability/{eventId}:
 *   get:
 *     summary: Check ticket availability for an event
 *     tags: [Ticket Availability]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *         description: Event ID
 *       - in: query
 *         name: ticket_type_id
 *         schema:
 *           type: string
 *         description: Specific ticket type ID
 *       - in: query
 *         name: quantity
 *         schema:
 *           type: integer
 *           default: 1
 *         description: Quantity to check
 *       - in: query
 *         name: seat_numbers
 *         schema:
 *           type: string
 *         description: Comma-separated seat numbers to check
 *     responses:
 *       200:
 *         description: Availability status
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/AvailabilityResponse'
 */

/**
 * @swagger
 * /api/tickets/available/{eventId}:
 *   get:
 *     summary: Get available tickets for an event
 *     tags: [Ticket Availability]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: query
 *         name: ticket_type_id
 *         schema:
 *           type: string
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 50
 *     responses:
 *       200:
 *         description: List of available tickets
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 tickets:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Ticket'
 *                 total_available:
 *                   type: integer
 */

/**
 * @swagger
 * /api/tickets/reserve/{eventId}:
 *   post:
 *     summary: Reserve tickets for an event
 *     description: Temporarily reserve tickets before purchase. Reservation expires after timeout.
 *     tags: [Ticket Reservation]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - ticket_type_id
 *               - quantity
 *             properties:
 *               ticket_type_id:
 *                 type: string
 *               quantity:
 *                 type: integer
 *                 minimum: 1
 *               seat_numbers:
 *                 type: array
 *                 items:
 *                   type: string
 *                 description: Specific seats to reserve
 *               timeout_seconds:
 *                 type: integer
 *                 default: 600
 *                 description: Reservation timeout in seconds
 *     responses:
 *       200:
 *         description: Tickets reserved successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/ReservationResponse'
 *       400:
 *         description: Tickets not available
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /api/tickets/release:
 *   post:
 *     summary: Release reserved tickets
 *     description: Release previously reserved tickets back to available pool
 *     tags: [Ticket Reservation]
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
 *                 description: Reservation ID from reserve response
 *               ticket_ids:
 *                 type: array
 *                 items:
 *                   type: string
 *                 description: Specific ticket IDs to release (optional)
 *     responses:
 *       200:
 *         description: Tickets released successfully
 *       401:
 *         description: Unauthorized
 */

export default {};

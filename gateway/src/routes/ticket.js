import express from 'express';
import {
  getTicketsHandler,
  getTicketHandler,
  createTicketHandler,
  updateTicketHandler,
  deleteTicketHandler,
  // TODO: Thêm handlers cho booking, seat reservation
  // createBookingSessionHandler,
  // getBookingSessionHandler,
  // reserveSeatsHandler,
  // releaseSeatsHandler,
} from '../handlers/ticketHandlers.js';

const router = express.Router();

/**
 * @swagger
 * /api/tickets:
 *   get:
 *     summary: Get all tickets
 *     tags: [Ticket]
 *     responses:
 *       200:
 *         description: List of tickets
 *   post:
 *     summary: Create a new ticket
 *     tags: [Ticket]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/Ticket'
 *     responses:
 *       201:
 *         description: Ticket created
 * /api/tickets/{ticketId}:
 *   get:
 *     summary: Get ticket by ID
 *     tags: [Ticket]
 *     parameters:
 *       - in: path
 *         name: ticketId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Ticket detail
 *   put:
 *     summary: Update ticket
 *     tags: [Ticket]
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
 *     tags: [Ticket]
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

// Ticket CRUD
router.get('/', getTicketsHandler);
router.post('/', createTicketHandler);
router.get('/:ticketId', getTicketHandler);
router.put('/:ticketId', updateTicketHandler);
router.delete('/:ticketId', deleteTicketHandler);

// Booking Session (event-centric)
// TODO: Implement khi có handlers
// router.post('/booking-sessions', createBookingSessionHandler);
// router.get('/booking-sessions/:sessionId', getBookingSessionHandler);

// Seat Reservation (event-centric)
// TODO: Implement khi có handlers
// router.post('/seats/reserve', reserveSeatsHandler);
// router.post('/seats/release', releaseSeatsHandler);

export default router;

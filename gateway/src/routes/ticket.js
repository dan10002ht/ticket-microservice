import express from 'express';
import {
  getTicketsHandler,
  getTicketHandler,
  createTicketHandler,
  updateTicketHandler,
  deleteTicketHandler,
  getTicketTypesHandler,
  createTicketTypeHandler,
  updateTicketTypeHandler,
  deleteTicketTypeHandler,
  checkAvailabilityHandler,
  getAvailableTicketsHandler,
  reserveTicketsHandler,
  releaseTicketsHandler,
} from '../handlers/ticketHandlers.js';
import { requireRole, requireAuth } from '../middlewares/index.js';

const router = express.Router();

// ============================================
// Ticket Types (must be before /:ticketId routes)
// ============================================
router.get('/types/:eventId', getTicketTypesHandler);
router.post('/types', requireRole(['organization']), createTicketTypeHandler);
router.put('/types/:typeId', requireRole(['organization']), updateTicketTypeHandler);
router.delete('/types/:typeId', requireRole(['organization']), deleteTicketTypeHandler);

// ============================================
// Availability
// ============================================
router.get('/availability/:eventId', checkAvailabilityHandler);
router.get('/available/:eventId', getAvailableTicketsHandler);

// ============================================
// Reservation
// ============================================
router.post('/reserve/:eventId', requireAuth, reserveTicketsHandler);
router.post('/release', requireAuth, releaseTicketsHandler);

// ============================================
// Ticket CRUD
// ============================================
router.get('/', getTicketsHandler);
router.post('/', createTicketHandler);
router.get('/:ticketId', getTicketHandler);
router.put('/:ticketId', updateTicketHandler);
router.delete('/:ticketId', deleteTicketHandler);

export default router;

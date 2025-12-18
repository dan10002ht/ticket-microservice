import express from 'express';

// Import handlers
import {
  getEventsHandler,
  getEventHandler,
  createEventHandler,
  updateEventHandler,
  deleteEventHandler,
  saveEventDraftHandler,
  publishEventHandler,
  getEventTemplatesHandler,
  duplicateEventHandler,
} from '../handlers/eventHandlers.js';

// Zone handlers
import {
  createZoneHandler,
  getZoneHandler,
  updateZoneHandler,
  deleteZoneHandler,
  listZonesByEventHandler,
} from '../handlers/zoneHandlers.js';

// Seat handlers
import {
  createSeatHandler,
  bulkCreateSeatsHandler,
  getSeatHandler,
  updateSeatHandler,
  deleteSeatHandler,
  listSeatsByEventHandler,
} from '../handlers/seatHandlers.js';

// Pricing handlers
import {
  createPricingHandler,
  getPricingHandler,
  updatePricingHandler,
  deletePricingHandler,
  listPricingByEventHandler,
  getPricingByZoneHandler,
  calculatePriceHandler,
  applyDiscountHandler,
} from '../handlers/pricingHandlers.js';

// Availability handlers
import {
  getEventAvailabilityHandler,
  getZoneAvailabilityHandler,
  getSeatAvailabilityHandler,
  updateSeatAvailabilityHandler,
  blockSeatsHandler,
  releaseSeatsHandler,
} from '../handlers/availabilityHandlers.js';

import { validateEvent } from '../middlewares/index.js';
import { requireRole } from '../middlewares/index.js';

const router = express.Router();

// Event CRUD
router.get('/', getEventsHandler);
router.post('/', requireRole(['organization']), validateEvent, createEventHandler);
router.get('/:eventId', getEventHandler);
router.put('/:eventId', updateEventHandler);
router.delete('/:eventId', deleteEventHandler);

// Draft event
router.put('/:eventId/draft', validateEvent, saveEventDraftHandler);

// Publish event
router.post('/:eventId/publish', publishEventHandler);

// Event templates
router.get('/templates', getEventTemplatesHandler);

// Duplicate event
router.post('/:eventId/duplicate', duplicateEventHandler);

// ============================================
// Zone Management
// ============================================
router.get('/:eventId/zones', listZonesByEventHandler);
router.post('/:eventId/zones', requireRole(['organization']), createZoneHandler);
router.get('/:eventId/zones/:zoneId', getZoneHandler);
router.put('/:eventId/zones/:zoneId', requireRole(['organization']), updateZoneHandler);
router.delete('/:eventId/zones/:zoneId', requireRole(['organization']), deleteZoneHandler);

// ============================================
// Seat Management
// ============================================
router.get('/:eventId/seats', listSeatsByEventHandler);
router.post('/:eventId/seats', requireRole(['organization']), createSeatHandler);
router.post('/:eventId/seats/bulk', requireRole(['organization']), bulkCreateSeatsHandler);
router.get('/:eventId/seats/:seatId', getSeatHandler);
router.put('/:eventId/seats/:seatId', requireRole(['organization']), updateSeatHandler);
router.delete('/:eventId/seats/:seatId', requireRole(['organization']), deleteSeatHandler);

// ============================================
// Pricing Management
// ============================================
router.get('/:eventId/pricing', listPricingByEventHandler);
router.post('/:eventId/pricing', requireRole(['organization']), createPricingHandler);
router.post('/:eventId/pricing/calculate', calculatePriceHandler);
router.post('/:eventId/pricing/discount', requireRole(['organization']), applyDiscountHandler);
router.get('/:eventId/pricing/zone/:zoneId', getPricingByZoneHandler);
router.get('/:eventId/pricing/:pricingId', getPricingHandler);
router.put('/:eventId/pricing/:pricingId', requireRole(['organization']), updatePricingHandler);
router.delete('/:eventId/pricing/:pricingId', requireRole(['organization']), deletePricingHandler);

// ============================================
// Availability Management
// ============================================
router.get('/:eventId/availability', getEventAvailabilityHandler);
router.get('/:eventId/availability/zones/:zoneId', getZoneAvailabilityHandler);
router.get('/:eventId/availability/seats/:seatId', getSeatAvailabilityHandler);
router.put('/:eventId/availability/seats/:seatId', requireRole(['organization']), updateSeatAvailabilityHandler);
router.post('/:eventId/availability/block', requireRole(['organization']), blockSeatsHandler);
router.post('/:eventId/availability/release', requireRole(['organization']), releaseSeatsHandler);

export default router;

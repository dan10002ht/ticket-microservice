import express from 'express';

// Import handlers
import {
  getEventsHandler,
  getEventHandler,
  createEventHandler,
  updateEventHandler,
  deleteEventHandler,
  saveEventDraftHandler,
  saveEventLayoutHandler,
  saveEventPricingHandler,
  publishEventHandler,
  getEventTemplatesHandler,
  duplicateEventHandler,
} from '../handlers/eventHandlers.js';
import { validateEvent } from '../middlewares/index.js';
import { requireRole } from '../middlewares/index.js';
// TODO: Thêm handlers cho zone, seat, layout
// getEventZonesHandler,
// getEventSeatsHandler,
// getEventLayoutHandler,
// } from '../handlers/eventHandlers.js';

const router = express.Router();

// Event CRUD
router.get('/', getEventsHandler);
router.post('/', requireRole(['organization']), validateEvent, createEventHandler);
router.get('/:eventId', getEventHandler);
router.put('/:eventId', updateEventHandler);
router.delete('/:eventId', deleteEventHandler);

// Draft event
router.put('/:eventId/draft', validateEvent, saveEventDraftHandler);

// Layout management
router.post('/:eventId/layout', saveEventLayoutHandler);

// Pricing management
router.post('/:eventId/pricing', saveEventPricingHandler);

// Publish event
router.post('/:eventId/publish', publishEventHandler);

// Event templates
router.get('/templates', getEventTemplatesHandler);

// Duplicate event
router.post('/:eventId/duplicate', duplicateEventHandler);

// Event Zones (venue/layout/zone/seat là property của event)
// TODO: Implement khi có handlers
// router.get('/:eventId/zones', getEventZonesHandler);
// router.get('/:eventId/seats', getEventSeatsHandler);
// router.get('/:eventId/layout', getEventLayoutHandler);

export default router;

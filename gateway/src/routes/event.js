import express from 'express';

// Import handlers
import {
  getEventsHandler,
  getEventHandler,
  createEventHandler,
  updateEventHandler,
  deleteEventHandler,
  // TODO: Thêm handlers cho zone, seat, layout
  // getEventZonesHandler,
  // getEventSeatsHandler,
  // getEventLayoutHandler,
} from '../handlers/index.js';

const router = express.Router();

// Event CRUD
router.get('/', getEventsHandler);
router.post('/', createEventHandler);
router.get('/:eventId', getEventHandler);
router.put('/:eventId', updateEventHandler);
router.delete('/:eventId', deleteEventHandler);

// Event Zones (venue/layout/zone/seat là property của event)
// TODO: Implement khi có handlers
// router.get('/:eventId/zones', getEventZonesHandler);
// router.get('/:eventId/seats', getEventSeatsHandler);
// router.get('/:eventId/layout', getEventLayoutHandler);

export default router;

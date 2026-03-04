import express from 'express';
import {
  checkInHandler,
  getCheckInHandler,
  listCheckInsHandler,
  getEventStatsHandler,
} from '../handlers/checkinHandlers.js';
import { requireRole } from '../middlewares/index.js';

const router = express.Router();

// ============================================
// Check-in Operations (staff/admin only)
// ============================================
router.post('/:eventId/checkin', requireRole(['staff', 'admin']), checkInHandler);

// ============================================
// Check-in Records
// ============================================
router.get('/event/:eventId', listCheckInsHandler);
router.get('/event/:eventId/stats', requireRole(['staff', 'admin']), getEventStatsHandler);
router.get('/:checkinId', getCheckInHandler);

export default router;

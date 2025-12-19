import express from 'express';

// Import handlers
import {
  getProfileHandler,
  updateProfileHandler,
  getAddressesHandler,
  addAddressHandler,
  updateAddressHandler,
  deleteAddressHandler,
  createProfileHandler,
  getLegacyUserHandler,
  createLegacyUserHandler,
  listLegacyUsersHandler,
} from '../handlers/userHandlers.js';

// Import validation middleware
import {
  validateUserProfileCreate,
  validateUserProfileUpdate,
  validateUserAddress,
  validateUserAddressUpdate,
} from '../middlewares/index.js';

import { requireRole } from '../middlewares/authorizationMiddleware.js';

const router = express.Router();

// ============================================
// Profile Operations (requires auth - applied at app level)
// ============================================
router.get('/profile', getProfileHandler);
router.post('/profile', validateUserProfileCreate, createProfileHandler);
router.put('/profile', validateUserProfileUpdate, updateProfileHandler);

// ============================================
// Address Operations (requires auth - applied at app level)
// ============================================
router.get('/addresses', getAddressesHandler);
router.post('/addresses', validateUserAddress, addAddressHandler);
router.put('/addresses/:addressId', validateUserAddressUpdate, updateAddressHandler);
router.delete('/addresses/:addressId', deleteAddressHandler);

// ============================================
// Legacy Operations (Admin only)
// ============================================
router.get('/admin/list', requireRole(['admin', 'super_admin']), listLegacyUsersHandler);
router.get('/admin/:userId', requireRole(['admin', 'super_admin']), getLegacyUserHandler);
router.post('/admin', requireRole(['admin', 'super_admin']), createLegacyUserHandler);

export default router;

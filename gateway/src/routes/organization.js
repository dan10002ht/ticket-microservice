import express from 'express';
import { getOrganizationDashboardHandler } from '../handlers/organizationHandlers.js';

const router = express.Router();

router.get('/:orgId/dashboard', getOrganizationDashboardHandler);

export default router;

import express from 'express';
import {
  getInvoiceHandler,
  listInvoicesHandler,
  getInvoicePdfHandler,
} from '../handlers/invoiceHandlers.js';

const router = express.Router();

// ============================================
// Invoice Queries
// ============================================
router.get('/', listInvoicesHandler);
router.get('/:invoiceId', getInvoiceHandler);
router.get('/:invoiceId/pdf', getInvoicePdfHandler);

export default router;

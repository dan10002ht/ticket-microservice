import express from 'express';
import { processPaymentWebhookHandler } from '../handlers/webhookHandlers.js';

const router = express.Router();

router.post('/payment/:gateway', processPaymentWebhookHandler);

export default router;

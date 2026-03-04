import swaggerUi from 'swagger-ui-express';
import authMiddleware from '../middlewares/auth.js';

// Import routes
import authRoutes from '../routes/auth.js';
import userRoutes from '../routes/user.js';
import eventRoutes from '../routes/event.js';
import bookingRoutes from '../routes/booking.js';
import paymentRoutes from '../routes/payment.js';
import ticketRoutes from '../routes/ticket.js';
import checkinRoutes from '../routes/checkin.js';
import invoiceRoutes from '../routes/invoice.js';
import healthRoutes from '../routes/health.js';
import organizationRoutes from '../routes/organization.js';
import webhookRoutes from '../routes/webhook.js';

/**
 * Initialize all API routes
 * @param {express.Application} app - Express app instance
 * @param {Object} swaggerSpec - Swagger specification
 */
export const initializeRoutes = (app, swaggerSpec) => {
  // Health check route (no authentication required)
  app.use('/health', healthRoutes);

  // Webhooks (public endpoint, signature verification handled downstream)
  app.use('/webhooks', webhookRoutes);

  // API documentation
  app.use('/api/docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec));

  // ============================================
  // API v1 — primary versioned routes
  // ============================================
  app.use('/api/v1/auth', authRoutes);
  app.use('/api/v1/users', authMiddleware, userRoutes);
  app.use('/api/v1/events', authMiddleware, eventRoutes);
  app.use('/api/v1/bookings', authMiddleware, bookingRoutes);
  app.use('/api/v1/payments', authMiddleware, paymentRoutes);
  app.use('/api/v1/tickets', authMiddleware, ticketRoutes);
  app.use('/api/v1/checkins', authMiddleware, checkinRoutes);
  app.use('/api/v1/invoices', authMiddleware, invoiceRoutes);
  app.use('/api/v1/organizations', authMiddleware, organizationRoutes);

  // ============================================
  // Legacy /api/* — backward-compat aliases (deprecated)
  // Clients should migrate to /api/v1/
  // ============================================
  const deprecationWarning = (_req, res, next) => {
    res.setHeader('Deprecation', 'true');
    res.setHeader('Link', '</api/v1/>; rel="successor-version"');
    next();
  };
  app.use('/api/auth', deprecationWarning, authRoutes);
  app.use('/api/users', deprecationWarning, authMiddleware, userRoutes);
  app.use('/api/events', deprecationWarning, authMiddleware, eventRoutes);
  app.use('/api/bookings', deprecationWarning, authMiddleware, bookingRoutes);
  app.use('/api/payments', deprecationWarning, authMiddleware, paymentRoutes);
  app.use('/api/tickets', deprecationWarning, authMiddleware, ticketRoutes);
  app.use('/api/checkins', deprecationWarning, authMiddleware, checkinRoutes);
  app.use('/api/invoices', deprecationWarning, authMiddleware, invoiceRoutes);
  app.use('/api/organizations', deprecationWarning, authMiddleware, organizationRoutes);
};

import swaggerUi from 'swagger-ui-express';
import authMiddleware from '../middlewares/auth.js';

// Import routes
import authRoutes from '../routes/auth.js';
import userRoutes from '../routes/user.js';
import eventRoutes from '../routes/event.js';
import bookingRoutes from '../routes/booking.js';
import paymentRoutes from '../routes/payment.js';
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

  // API routes with authentication middleware
  app.use('/api/auth', authRoutes);
  app.use('/api/users', authMiddleware, userRoutes);
  app.use('/api/events', authMiddleware, eventRoutes);
  app.use('/api/bookings', authMiddleware, bookingRoutes);
  app.use('/api/payments', authMiddleware, paymentRoutes);
  app.use('/api/organizations', authMiddleware, organizationRoutes);
};

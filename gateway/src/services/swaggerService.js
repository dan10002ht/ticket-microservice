import swaggerJsdoc from 'swagger-jsdoc';
import config from '../config/index.js';

// Import all swagger documentation files
import '../swagger/auth.js';
import '../swagger/user.js';
import '../swagger/booking.js';
import '../swagger/event.js';
import '../swagger/payment.js';
import '../swagger/health.js';
import '../swagger/ticket.js';

/**
 * Initialize Swagger documentation
 * @returns {Object} Swagger specification
 */
export const initializeSwagger = () => {
  const swaggerOptions = {
    definition: {
      openapi: '3.0.0',
      info: {
        title: 'Booking System API Gateway',
        version: '1.0.0',
        description: 'API Gateway for Booking System Microservices',
        contact: {
          name: 'API Support',
          email: 'support@bookingsystem.com',
        },
        license: {
          name: 'MIT',
          url: 'https://opensource.org/licenses/MIT',
        },
      },
      servers: [
        {
          url: `http://localhost:${config.server.port}/api`,
          description: 'Development server',
        },
        {
          url: 'https://api.bookingsystem.com/api',
          description: 'Production server',
        },
      ],
      components: {
        securitySchemes: {
          bearerAuth: {
            type: 'http',
            scheme: 'bearer',
            bearerFormat: 'JWT',
            description: 'JWT token for authentication',
          },
        },
      },
      tags: [
        {
          name: 'Auth',
          description: 'Authentication and authorization endpoints',
        },
        {
          name: 'Users',
          description: 'User profile and address management',
        },
        {
          name: 'Users Admin',
          description: 'Admin user management endpoints',
        },
        {
          name: 'Events',
          description: 'Event management and discovery',
        },
        {
          name: 'Event Zones',
          description: 'Event zone/section management',
        },
        {
          name: 'Event Seats',
          description: 'Event seat management',
        },
        {
          name: 'Event Pricing',
          description: 'Event pricing and discounts',
        },
        {
          name: 'Event Availability',
          description: 'Event seat availability management',
        },
        {
          name: 'Bookings',
          description: 'Booking creation and management',
        },
        {
          name: 'Bookings Admin',
          description: 'Admin booking management',
        },
        {
          name: 'Seat Reservation',
          description: 'Temporary seat reservation',
        },
        {
          name: 'Tickets',
          description: 'Ticket management',
        },
        {
          name: 'Ticket Types',
          description: 'Ticket type management',
        },
        {
          name: 'Ticket Availability',
          description: 'Ticket availability checking',
        },
        {
          name: 'Ticket Reservation',
          description: 'Ticket reservation',
        },
        {
          name: 'Payments',
          description: 'Payment processing and management',
        },
        {
          name: 'Payments Admin',
          description: 'Admin payment management',
        },
        {
          name: 'Payment Operations',
          description: 'Payment capture and cancellation',
        },
        {
          name: 'Refunds',
          description: 'Payment refund management',
        },
        {
          name: 'Refunds Admin',
          description: 'Admin refund management',
        },
        {
          name: 'Health',
          description: 'Health check and monitoring endpoints',
        },
      ],
    },
    apis: ['./src/routes/*.js', './src/swagger/*.js'],
  };

  const spec = swaggerJsdoc(swaggerOptions);

  return spec;
};

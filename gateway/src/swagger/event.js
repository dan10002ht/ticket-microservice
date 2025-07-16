/**
 * @swagger
 * components:
 *   schemas:
 *     Event:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         organization_id:
 *           type: string
 *         name:
 *           type: string
 *         description:
 *           type: string
 *         start_date:
 *           type: string
 *         end_date:
 *           type: string
 *         venue_name:
 *           type: string
 *         venue_address:
 *           type: string
 *         venue_city:
 *           type: string
 *         venue_country:
 *           type: string
 *         venue_capacity:
 *           type: integer
 *         canvas_config:
 *           type: string
 *         zones:
 *           type: array
 *           items:
 *             $ref: '#/components/schemas/EventSeatingZone'
 *         seats:
 *           type: array
 *           items:
 *             $ref: '#/components/schemas/EventSeat'
 *         created_at:
 *           type: string
 *         updated_at:
 *           type: string
 *     EventSeatingZone:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         event_id:
 *           type: string
 *         name:
 *           type: string
 *         zone_type:
 *           type: string
 *         coordinates:
 *           type: string
 *         seat_count:
 *           type: integer
 *         color:
 *           type: string
 *         created_at:
 *           type: string
 *         updated_at:
 *           type: string
 *     EventSeat:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         event_id:
 *           type: string
 *         zone_id:
 *           type: string
 *         seat_number:
 *           type: string
 *         row_number:
 *           type: string
 *         coordinates:
 *           type: string
 *         created_at:
 *           type: string
 *         updated_at:
 *           type: string
 */

/**
 * @swagger
 * /events:
 *   get:
 *     summary: Get all events
 *     description: Retrieve all available events with optional filtering
 *     tags: [Events]
 *     parameters:
 *       - in: query
 *         name: category
 *         schema:
 *           type: string
 *         description: Filter by event category
 *       - in: query
 *         name: location
 *         schema:
 *           type: string
 *         description: Filter by event location
 *       - in: query
 *         name: date
 *         schema:
 *           type: string
 *           format: date
 *         description: Filter by event date
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [draft, published, cancelled, completed]
 *         description: Filter by event status
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number for pagination
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 10
 *         description: Number of items per page
 *     responses:
 *       200:
 *         description: Events retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 events:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Event'
 *                 pagination:
 *                   type: object
 *                   properties:
 *                     page:
 *                       type: integer
 *                     limit:
 *                       type: integer
 *                     total:
 *                       type: integer
 *                     totalPages:
 *                       type: integer
 */

/**
 * @swagger
 * /events:
 *   post:
 *     summary: Create new event
 *     description: Create a new event
 *     tags: [Events]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/EventCreate'
 *     responses:
 *       201:
 *         description: Event created successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Event'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}:
 *   get:
 *     summary: Get event by ID
 *     description: Retrieve a specific event by ID
 *     tags: [Events]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *         description: Event ID
 *     responses:
 *       200:
 *         description: Event retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Event'
 *       404:
 *         description: Event not found
 */

/**
 * @swagger
 * /events/{eventId}:
 *   put:
 *     summary: Update event
 *     description: Update an existing event
 *     tags: [Events]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *         description: Event ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/EventUpdate'
 *     responses:
 *       200:
 *         description: Event updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Event'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Event not found
 *       403:
 *         description: Forbidden - not the event organizer
 */

/**
 * @swagger
 * /events/{eventId}:
 *   delete:
 *     summary: Delete event
 *     description: Delete an event
 *     tags: [Events]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *         description: Event ID
 *     responses:
 *       200:
 *         description: Event deleted successfully
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Event not found
 *       403:
 *         description: Forbidden - not the event organizer
 */

export default {};

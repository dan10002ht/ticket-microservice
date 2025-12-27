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

/**
 * @swagger
 * components:
 *   schemas:
 *     EventCreate:
 *       type: object
 *       required:
 *         - name
 *         - start_date
 *         - venue_name
 *       properties:
 *         name:
 *           type: string
 *         description:
 *           type: string
 *         start_date:
 *           type: string
 *           format: date-time
 *         end_date:
 *           type: string
 *           format: date-time
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
 *     EventUpdate:
 *       type: object
 *       properties:
 *         name:
 *           type: string
 *         description:
 *           type: string
 *         start_date:
 *           type: string
 *           format: date-time
 *         end_date:
 *           type: string
 *           format: date-time
 *         venue_name:
 *           type: string
 *         venue_address:
 *           type: string
 *         status:
 *           type: string
 *           enum: [draft, published, cancelled, completed]
 *     ZoneCreate:
 *       type: object
 *       required:
 *         - name
 *         - zone_type
 *       properties:
 *         name:
 *           type: string
 *         zone_type:
 *           type: string
 *           enum: [seated, standing, vip]
 *         coordinates:
 *           type: string
 *         seat_count:
 *           type: integer
 *         color:
 *           type: string
 *     SeatCreate:
 *       type: object
 *       required:
 *         - zone_id
 *         - seat_number
 *       properties:
 *         zone_id:
 *           type: string
 *         seat_number:
 *           type: string
 *         row_number:
 *           type: string
 *         coordinates:
 *           type: string
 *     BulkSeatCreate:
 *       type: object
 *       required:
 *         - zone_id
 *         - seats
 *       properties:
 *         zone_id:
 *           type: string
 *         seats:
 *           type: array
 *           items:
 *             type: object
 *             properties:
 *               seat_number:
 *                 type: string
 *               row_number:
 *                 type: string
 *               coordinates:
 *                 type: string
 *     PricingCreate:
 *       type: object
 *       required:
 *         - zone_id
 *         - name
 *         - price
 *       properties:
 *         zone_id:
 *           type: string
 *         name:
 *           type: string
 *         price:
 *           type: number
 *         currency:
 *           type: string
 *           default: USD
 *     Pricing:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         event_id:
 *           type: string
 *         zone_id:
 *           type: string
 *         name:
 *           type: string
 *         price:
 *           type: number
 *         currency:
 *           type: string
 *     Availability:
 *       type: object
 *       properties:
 *         total_seats:
 *           type: integer
 *         available_seats:
 *           type: integer
 *         reserved_seats:
 *           type: integer
 *         sold_seats:
 *           type: integer
 *         blocked_seats:
 *           type: integer
 */

/**
 * @swagger
 * /events/{eventId}/draft:
 *   put:
 *     summary: Save event as draft
 *     tags: [Events]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/EventUpdate'
 *     responses:
 *       200:
 *         description: Draft saved
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden
 */

/**
 * @swagger
 * /events/{eventId}/publish:
 *   post:
 *     summary: Publish event
 *     tags: [Events]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Event published
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden
 */

/**
 * @swagger
 * /events/templates:
 *   get:
 *     summary: Get event templates
 *     tags: [Events]
 *     responses:
 *       200:
 *         description: Templates retrieved
 */

/**
 * @swagger
 * /events/{eventId}/duplicate:
 *   post:
 *     summary: Duplicate event
 *     tags: [Events]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       201:
 *         description: Event duplicated
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden
 */

/**
 * @swagger
 * /events/{eventId}/zones:
 *   get:
 *     summary: List zones for event
 *     tags: [Event Zones]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Zones retrieved
 *   post:
 *     summary: Create zone
 *     tags: [Event Zones]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/ZoneCreate'
 *     responses:
 *       201:
 *         description: Zone created
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden
 */

/**
 * @swagger
 * /events/{eventId}/zones/{zoneId}:
 *   get:
 *     summary: Get zone by ID
 *     tags: [Event Zones]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: zoneId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Zone retrieved
 *       404:
 *         description: Not found
 *   put:
 *     summary: Update zone
 *     tags: [Event Zones]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: zoneId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/ZoneCreate'
 *     responses:
 *       200:
 *         description: Zone updated
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Not found
 *   delete:
 *     summary: Delete zone
 *     tags: [Event Zones]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: zoneId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Zone deleted
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/seats:
 *   get:
 *     summary: List seats for event
 *     tags: [Event Seats]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: query
 *         name: zone_id
 *         schema:
 *           type: string
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [available, reserved, sold, blocked]
 *     responses:
 *       200:
 *         description: Seats retrieved
 *   post:
 *     summary: Create seat
 *     tags: [Event Seats]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/SeatCreate'
 *     responses:
 *       201:
 *         description: Seat created
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/seats/bulk:
 *   post:
 *     summary: Bulk create seats
 *     tags: [Event Seats]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/BulkSeatCreate'
 *     responses:
 *       201:
 *         description: Seats created
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/seats/{seatId}:
 *   get:
 *     summary: Get seat by ID
 *     tags: [Event Seats]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: seatId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Seat retrieved
 *       404:
 *         description: Not found
 *   put:
 *     summary: Update seat
 *     tags: [Event Seats]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: seatId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/SeatCreate'
 *     responses:
 *       200:
 *         description: Seat updated
 *       401:
 *         description: Unauthorized
 *   delete:
 *     summary: Delete seat
 *     tags: [Event Seats]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: seatId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Seat deleted
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/pricing:
 *   get:
 *     summary: List pricing for event
 *     tags: [Event Pricing]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Pricing list retrieved
 *   post:
 *     summary: Create pricing
 *     tags: [Event Pricing]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/PricingCreate'
 *     responses:
 *       201:
 *         description: Pricing created
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/pricing/calculate:
 *   post:
 *     summary: Calculate price
 *     tags: [Event Pricing]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               seat_ids:
 *                 type: array
 *                 items:
 *                   type: string
 *               quantity:
 *                 type: integer
 *               zone_id:
 *                 type: string
 *               discount_code:
 *                 type: string
 *     responses:
 *       200:
 *         description: Price calculated
 */

/**
 * @swagger
 * /events/{eventId}/pricing/discount:
 *   post:
 *     summary: Apply discount
 *     tags: [Event Pricing]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - discount_code
 *             properties:
 *               discount_code:
 *                 type: string
 *               pricing_id:
 *                 type: string
 *     responses:
 *       200:
 *         description: Discount applied
 *       400:
 *         description: Invalid discount code
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/pricing/zone/{zoneId}:
 *   get:
 *     summary: Get pricing by zone
 *     tags: [Event Pricing]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: zoneId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Pricing retrieved
 */

/**
 * @swagger
 * /events/{eventId}/pricing/{pricingId}:
 *   get:
 *     summary: Get pricing by ID
 *     tags: [Event Pricing]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: pricingId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Pricing retrieved
 *   put:
 *     summary: Update pricing
 *     tags: [Event Pricing]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: pricingId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/PricingCreate'
 *     responses:
 *       200:
 *         description: Pricing updated
 *       401:
 *         description: Unauthorized
 *   delete:
 *     summary: Delete pricing
 *     tags: [Event Pricing]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: pricingId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Pricing deleted
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/availability:
 *   get:
 *     summary: Get event availability
 *     tags: [Event Availability]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Availability retrieved
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Availability'
 */

/**
 * @swagger
 * /events/{eventId}/availability/zones/{zoneId}:
 *   get:
 *     summary: Get zone availability
 *     tags: [Event Availability]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: zoneId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Zone availability retrieved
 */

/**
 * @swagger
 * /events/{eventId}/availability/seats/{seatId}:
 *   get:
 *     summary: Get seat availability
 *     tags: [Event Availability]
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: seatId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Seat availability retrieved
 *   put:
 *     summary: Update seat availability
 *     tags: [Event Availability]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: path
 *         name: seatId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               status:
 *                 type: string
 *                 enum: [available, reserved, sold, blocked]
 *     responses:
 *       200:
 *         description: Seat availability updated
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/availability/block:
 *   post:
 *     summary: Block seats
 *     tags: [Event Availability]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - seat_ids
 *             properties:
 *               seat_ids:
 *                 type: array
 *                 items:
 *                   type: string
 *               reason:
 *                 type: string
 *     responses:
 *       200:
 *         description: Seats blocked
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /events/{eventId}/availability/release:
 *   post:
 *     summary: Release blocked seats
 *     tags: [Event Availability]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - seat_ids
 *             properties:
 *               seat_ids:
 *                 type: array
 *                 items:
 *                   type: string
 *     responses:
 *       200:
 *         description: Seats released
 *       401:
 *         description: Unauthorized
 */

export default {};

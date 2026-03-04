/**
 * @swagger
 * components:
 *   schemas:
 *     CheckInRecord:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *         ticket_id:
 *           type: string
 *         event_id:
 *           type: string
 *         user_id:
 *           type: string
 *         staff_id:
 *           type: string
 *         qr_code:
 *           type: string
 *         status:
 *           type: string
 *           enum: [success, invalid, already_used, cancelled]
 *         check_in_time:
 *           type: integer
 *           format: int64
 *         device_id:
 *           type: string
 *         gate:
 *           type: string
 *         notes:
 *           type: string
 *         created_at:
 *           type: integer
 *           format: int64
 *     EventStats:
 *       type: object
 *       properties:
 *         success:
 *           type: boolean
 *         total_checkins:
 *           type: integer
 *         unique_tickets:
 *           type: integer
 *         by_gate:
 *           type: object
 *           additionalProperties:
 *             type: integer
 *         last_checkin_at:
 *           type: integer
 *           format: int64
 */

/**
 * @swagger
 * /checkins/{eventId}/checkin:
 *   post:
 *     summary: Process ticket check-in
 *     tags: [Check-in]
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
 *             type: object
 *             required:
 *               - ticket_id
 *               - qr_code
 *             properties:
 *               ticket_id:
 *                 type: string
 *               qr_code:
 *                 type: string
 *               device_id:
 *                 type: string
 *               gate:
 *                 type: string
 *     responses:
 *       201:
 *         description: Check-in processed
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 record:
 *                   $ref: '#/components/schemas/CheckInRecord'
 *                 message:
 *                   type: string
 *       400:
 *         description: Invalid ticket or already checked in
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Staff/admin role required
 */

/**
 * @swagger
 * /checkins/{checkinId}:
 *   get:
 *     summary: Get check-in record by ID
 *     tags: [Check-in]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: checkinId
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Check-in record
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 record:
 *                   $ref: '#/components/schemas/CheckInRecord'
 *       404:
 *         description: Check-in not found
 */

/**
 * @swagger
 * /checkins/event/{eventId}:
 *   get:
 *     summary: List check-ins for an event
 *     tags: [Check-in]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: eventId
 *         required: true
 *         schema:
 *           type: string
 *       - in: query
 *         name: user_id
 *         schema:
 *           type: string
 *       - in: query
 *         name: gate
 *         schema:
 *           type: string
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 20
 *     responses:
 *       200:
 *         description: List of check-in records
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 records:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/CheckInRecord'
 *                 total:
 *                   type: integer
 *                 page:
 *                   type: integer
 *                 limit:
 *                   type: integer
 *                 has_more:
 *                   type: boolean
 */

/**
 * @swagger
 * /checkins/event/{eventId}/stats:
 *   get:
 *     summary: Get check-in statistics for an event
 *     tags: [Check-in]
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
 *         description: Event check-in statistics
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/EventStats'
 *       403:
 *         description: Staff/admin role required
 */

export default {};

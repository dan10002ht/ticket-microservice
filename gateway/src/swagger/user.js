/**
 * @swagger
 * components:
 *   schemas:
 *     UserProfile:
 *       type: object
 *       properties:
 *         id:
 *           type: string
 *           description: User ID
 *         email:
 *           type: string
 *           format: email
 *           description: User email
 *         firstName:
 *           type: string
 *           description: User first name
 *         lastName:
 *           type: string
 *           description: User last name
 *         phone:
 *           type: string
 *           description: User phone number
 *         dateOfBirth:
 *           type: string
 *           format: date
 *           description: User date of birth
 *         createdAt:
 *           type: string
 *           format: date-time
 *           description: Account creation date
 *         updatedAt:
 *           type: string
 *           format: date-time
 *           description: Last update date
 *     UserProfileUpdate:
 *       type: object
 *       properties:
 *         firstName:
 *           type: string
 *           minLength: 2
 *           maxLength: 50
 *           description: User first name
 *         lastName:
 *           type: string
 *           minLength: 2
 *           maxLength: 50
 *           description: User last name
 *         phone:
 *           type: string
 *           description: User phone number
 *         dateOfBirth:
 *           type: string
 *           format: date
 *           description: User date of birth
 *     UserAddress:
 *       type: object
 *       required:
 *         - street
 *         - city
 *         - state
 *         - zipCode
 *         - country
 *       properties:
 *         id:
 *           type: string
 *           description: Address ID
 *         street:
 *           type: string
 *           maxLength: 200
 *           description: Street address
 *         city:
 *           type: string
 *           maxLength: 100
 *           description: City
 *         state:
 *           type: string
 *           maxLength: 100
 *           description: State/Province
 *         zipCode:
 *           type: string
 *           maxLength: 20
 *           description: ZIP/Postal code
 *         country:
 *           type: string
 *           maxLength: 100
 *           description: Country
 *         isDefault:
 *           type: boolean
 *           description: Whether this is the default address
 *     UserAddressCreate:
 *       type: object
 *       required:
 *         - street
 *         - city
 *         - state
 *         - zipCode
 *         - country
 *       properties:
 *         street:
 *           type: string
 *           maxLength: 200
 *           description: Street address
 *         city:
 *           type: string
 *           maxLength: 100
 *           description: City
 *         state:
 *           type: string
 *           maxLength: 100
 *           description: State/Province
 *         zipCode:
 *           type: string
 *           maxLength: 20
 *           description: ZIP/Postal code
 *         country:
 *           type: string
 *           maxLength: 100
 *           description: Country
 *         isDefault:
 *           type: boolean
 *           description: Whether this is the default address
 */

/**
 * @swagger
 * /users/profile:
 *   get:
 *     summary: Get user profile
 *     description: Retrieve the current user's profile information
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: User profile retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/UserProfile'
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /users/profile:
 *   put:
 *     summary: Update user profile
 *     description: Update the current user's profile information
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UserProfileUpdate'
 *     responses:
 *       200:
 *         description: Profile updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/UserProfile'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /users/addresses:
 *   get:
 *     summary: Get user addresses
 *     description: Retrieve all addresses for the current user
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Addresses retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/UserAddress'
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /users/addresses:
 *   post:
 *     summary: Add user address
 *     description: Add a new address for the current user
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UserAddressCreate'
 *     responses:
 *       201:
 *         description: Address added successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/UserAddress'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 */

/**
 * @swagger
 * /users/addresses/{addressId}:
 *   put:
 *     summary: Update user address
 *     description: Update an existing address for the current user
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: addressId
 *         required: true
 *         schema:
 *           type: string
 *         description: Address ID
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UserAddressCreate'
 *     responses:
 *       200:
 *         description: Address updated successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/UserAddress'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Address not found
 */

/**
 * @swagger
 * /users/addresses/{addressId}:
 *   delete:
 *     summary: Delete user address
 *     description: Delete an address for the current user
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: addressId
 *         required: true
 *         schema:
 *           type: string
 *         description: Address ID
 *     responses:
 *       200:
 *         description: Address deleted successfully
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Address not found
 */

/**
 * @swagger
 * /users/profile:
 *   post:
 *     summary: Create user profile
 *     description: Create a new profile for the current user (called after registration)
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               first_name:
 *                 type: string
 *                 description: First name
 *               last_name:
 *                 type: string
 *                 description: Last name
 *               email:
 *                 type: string
 *                 format: email
 *                 description: Email address
 *               phone:
 *                 type: string
 *                 description: Phone number
 *               avatar_url:
 *                 type: string
 *                 description: Avatar URL
 *               date_of_birth:
 *                 type: string
 *                 format: date
 *                 description: Date of birth
 *               preferences:
 *                 type: object
 *                 additionalProperties:
 *                   type: string
 *                 description: User preferences
 *     responses:
 *       201:
 *         description: Profile created successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 profile:
 *                   $ref: '#/components/schemas/UserProfile'
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       409:
 *         description: Profile already exists
 */

/**
 * @swagger
 * /users/admin/list:
 *   get:
 *     summary: List all users (Admin)
 *     description: Retrieve a paginated list of all users (Admin only)
 *     tags: [Users Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number
 *       - in: query
 *         name: size
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 20
 *         description: Number of items per page
 *     responses:
 *       200:
 *         description: Users list retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 users:
 *                   type: array
 *                   items:
 *                     type: object
 *                     properties:
 *                       id:
 *                         type: string
 *                       name:
 *                         type: string
 *                       email:
 *                         type: string
 *                 total:
 *                   type: integer
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 */

/**
 * @swagger
 * /users/admin/{userId}:
 *   get:
 *     summary: Get user by ID (Admin)
 *     description: Retrieve a specific user by ID (Admin only)
 *     tags: [Users Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: userId
 *         required: true
 *         schema:
 *           type: string
 *         description: User ID
 *     responses:
 *       200:
 *         description: User retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 id:
 *                   type: string
 *                 name:
 *                   type: string
 *                 email:
 *                   type: string
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 *       404:
 *         description: User not found
 */

/**
 * @swagger
 * /users/admin:
 *   post:
 *     summary: Create user (Admin)
 *     description: Create a new user (Admin only)
 *     tags: [Users Admin]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - name
 *               - email
 *             properties:
 *               name:
 *                 type: string
 *                 description: User name
 *               email:
 *                 type: string
 *                 format: email
 *                 description: User email
 *     responses:
 *       201:
 *         description: User created successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 id:
 *                   type: string
 *                 name:
 *                   type: string
 *                 email:
 *                   type: string
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin role required
 *       409:
 *         description: User already exists
 */

export default {};

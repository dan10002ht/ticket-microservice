import grpcClients from '../grpc/clients.js';
import {
  sendSuccessResponse,
  createHandler,
  createSimpleHandler,
} from '../utils/responseHandler.js';

/**
 * Get user profile
 */
const getUserProfile = async (req, res) => {
  const result = await grpcClients.userService.getProfile({
    user_id: req.user.id,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update user profile
 */
const updateUserProfile = async (req, res) => {
  const result = await grpcClients.userService.updateProfile({
    user_id: req.user.id,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get user addresses
 */
const getUserAddresses = async (req, res) => {
  const result = await grpcClients.userService.getAddresses({
    user_id: req.user.id,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Add user address
 */
const addUserAddress = async (req, res) => {
  const result = await grpcClients.userService.addAddress({
    user_id: req.user.id,
    ...req.body,
  });
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Update user address
 */
const updateUserAddress = async (req, res) => {
  const result = await grpcClients.userService.updateAddress({
    user_id: req.user.id,
    address_id: req.params.addressId,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Delete user address
 */
const deleteUserAddress = async (req, res) => {
  await grpcClients.userService.deleteAddress({
    user_id: req.user.id,
    address_id: req.params.addressId,
  });
  sendSuccessResponse(res, 200, { message: 'Address deleted successfully' }, req.correlationId);
};

/**
 * Create user profile (called after registration)
 */
const createUserProfile = async (req, res) => {
  const result = await grpcClients.userService.CreateProfile({
    user_id: req.user.id,
    first_name: req.body.first_name,
    last_name: req.body.last_name,
    email: req.body.email,
    phone: req.body.phone,
    avatar_url: req.body.avatar_url,
    date_of_birth: req.body.date_of_birth,
    preferences: req.body.preferences || {},
  });
  sendSuccessResponse(res, 201, result, req.correlationId);
};

// ============================================
// Legacy Operations (Admin)
// ============================================

/**
 * Get user by ID (Legacy - Admin)
 */
const getLegacyUser = async (req, res) => {
  const result = await grpcClients.userService.GetUser({
    id: req.params.userId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Create user (Legacy - Admin)
 */
const createLegacyUser = async (req, res) => {
  const result = await grpcClients.userService.CreateUser({
    name: req.body.name,
    email: req.body.email,
  });
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * List users (Legacy - Admin)
 */
const listLegacyUsers = async (req, res) => {
  const result = await grpcClients.userService.ListUsers({
    page: parseInt(req.query.page) || 1,
    size: parseInt(req.query.size) || 20,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// Export wrapped handlers
export const getProfileHandler = createSimpleHandler(getUserProfile, 'user', 'getProfile');
export const updateProfileHandler = createHandler(updateUserProfile, 'user', 'updateProfile');

export const getAddressesHandler = createSimpleHandler(getUserAddresses, 'user', 'getAddresses');
export const addAddressHandler = createHandler(addUserAddress, 'user', 'addAddress');
export const updateAddressHandler = createHandler(updateUserAddress, 'user', 'updateAddress');
export const deleteAddressHandler = createSimpleHandler(deleteUserAddress, 'user', 'deleteAddress');

// Create profile
export const createProfileHandler = createHandler(createUserProfile, 'user', 'createProfile');

// Legacy operations (Admin)
export const getLegacyUserHandler = createSimpleHandler(getLegacyUser, 'user', 'getLegacyUser');
export const createLegacyUserHandler = createHandler(createLegacyUser, 'user', 'createLegacyUser');
export const listLegacyUsersHandler = createSimpleHandler(listLegacyUsers, 'user', 'listLegacyUsers');

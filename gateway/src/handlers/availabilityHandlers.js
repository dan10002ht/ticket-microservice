import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

/**
 * Get availability for entire event
 */
const getEventAvailability = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.availabilityService.GetEventAvailability({
    event_id: eventId,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get availability for specific zone
 */
const getZoneAvailability = async (req, res) => {
  const { eventId, zoneId } = req.params;
  const result = await grpcClients.availabilityService.GetZoneAvailability({
    event_id: eventId,
    zone_id: zoneId,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get availability for specific seat
 */
const getSeatAvailability = async (req, res) => {
  const { eventId, seatId } = req.params;
  const result = await grpcClients.availabilityService.GetSeatAvailability({
    event_id: eventId,
    seat_id: seatId,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 404;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update seat availability status
 */
const updateSeatAvailability = async (req, res) => {
  const { eventId, seatId } = req.params;
  const result = await grpcClients.availabilityService.UpdateSeatAvailability({
    event_id: eventId,
    seat_id: seatId,
    ...req.body,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Block multiple seats for booking
 */
const blockSeats = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.availabilityService.BlockSeats({
    event_id: eventId,
    ...req.body,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Release blocked seats
 */
const releaseSeats = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.availabilityService.ReleaseSeats({
    event_id: eventId,
    ...req.body,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

export const getEventAvailabilityHandler = createSimpleHandler(getEventAvailability, 'availability', 'getEventAvailability');
export const getZoneAvailabilityHandler = createSimpleHandler(getZoneAvailability, 'availability', 'getZoneAvailability');
export const getSeatAvailabilityHandler = createSimpleHandler(getSeatAvailability, 'availability', 'getSeatAvailability');
export const updateSeatAvailabilityHandler = createSimpleHandler(updateSeatAvailability, 'availability', 'updateSeatAvailability');
export const blockSeatsHandler = createSimpleHandler(blockSeats, 'availability', 'blockSeats');
export const releaseSeatsHandler = createSimpleHandler(releaseSeats, 'availability', 'releaseSeats');

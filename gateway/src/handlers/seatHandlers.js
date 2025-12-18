import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

/**
 * Create seat for an event
 */
const createSeat = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.seatService.CreateSeat({
    event_id: eventId,
    ...req.body,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Bulk create seats for an event
 */
const bulkCreateSeats = async (req, res) => {
  const { eventId } = req.params;
  const { seats } = req.body;

  const result = await grpcClients.seatService.BulkCreateSeats({
    event_id: eventId,
    seats: seats,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Get seat by ID
 */
const getSeat = async (req, res) => {
  const { seatId } = req.params;
  const result = await grpcClients.seatService.GetSeat({ seat_id: seatId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 404;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update seat
 */
const updateSeat = async (req, res) => {
  const { seatId } = req.params;
  const result = await grpcClients.seatService.UpdateSeat({
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
 * Delete seat
 */
const deleteSeat = async (req, res) => {
  const { seatId } = req.params;
  const result = await grpcClients.seatService.DeleteSeat({ seat_id: seatId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, { message: 'Seat deleted successfully' }, req.correlationId);
};

/**
 * List seats by event
 */
const listSeatsByEvent = async (req, res) => {
  const { eventId } = req.params;
  const { zone_id, status, page, limit } = req.query;

  const result = await grpcClients.seatService.ListSeatsByEvent({
    event_id: eventId,
    zone_id: zone_id || '',
    status: status || '',
    page: parseInt(page) || 1,
    limit: parseInt(limit) || 50,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

export const createSeatHandler = createSimpleHandler(createSeat, 'seat', 'createSeat');
export const bulkCreateSeatsHandler = createSimpleHandler(bulkCreateSeats, 'seat', 'bulkCreateSeats');
export const getSeatHandler = createSimpleHandler(getSeat, 'seat', 'getSeat');
export const updateSeatHandler = createSimpleHandler(updateSeat, 'seat', 'updateSeat');
export const deleteSeatHandler = createSimpleHandler(deleteSeat, 'seat', 'deleteSeat');
export const listSeatsByEventHandler = createSimpleHandler(listSeatsByEvent, 'seat', 'listSeatsByEvent');

import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

/**
 * Create zone for an event
 */
const createZone = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.zoneService.CreateZone({
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
 * Get zone by ID
 */
const getZone = async (req, res) => {
  const { zoneId } = req.params;
  const result = await grpcClients.zoneService.GetZone({ zone_id: zoneId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 404;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update zone
 */
const updateZone = async (req, res) => {
  const { zoneId } = req.params;
  const result = await grpcClients.zoneService.UpdateZone({
    zone_id: zoneId,
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
 * Delete zone
 */
const deleteZone = async (req, res) => {
  const { zoneId } = req.params;
  const result = await grpcClients.zoneService.DeleteZone({ zone_id: zoneId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, { message: 'Zone deleted successfully' }, req.correlationId);
};

/**
 * List zones by event
 */
const listZonesByEvent = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.zoneService.ListZonesByEvent({ event_id: eventId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

export const createZoneHandler = createSimpleHandler(createZone, 'zone', 'createZone');
export const getZoneHandler = createSimpleHandler(getZone, 'zone', 'getZone');
export const updateZoneHandler = createSimpleHandler(updateZone, 'zone', 'updateZone');
export const deleteZoneHandler = createSimpleHandler(deleteZone, 'zone', 'deleteZone');
export const listZonesByEventHandler = createSimpleHandler(listZonesByEvent, 'zone', 'listZonesByEvent');

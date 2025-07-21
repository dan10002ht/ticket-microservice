import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

/**
 * Get all events
 */
const getAllEvents = async (req, res) => {
  const result = await grpcClients.eventService.getEvents({});
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get event by ID
 */
const getEventById = async (req, res) => {
  const result = await grpcClients.eventService.getEvent({ eventId: req.params.eventId });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Create new event
 */
const createNewEvent = async (req, res) => {
  const result = await grpcClients.eventService.createEvent(req.body);
  sendSuccessResponse(res, 201, result, req.correlationId);
};

/**
 * Update event
 */
const updateEvent = async (req, res) => {
  const result = await grpcClients.eventService.updateEvent({
    eventId: req.params.eventId,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Delete event
 */
const deleteEvent = async (req, res) => {
  await grpcClients.eventService.deleteEvent({ eventId: req.params.eventId });
  sendSuccessResponse(res, 200, { message: 'Event deleted successfully' }, req.correlationId);
};

/**
 * Save event draft (partial update)
 */
const saveEventDraft = async (req, res) => {
  // Validate eventId param là UUID
  const eventId = req.params.eventId;
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  if (!uuidRegex.test(eventId)) {
    const error = new Error('Invalid eventId format (must be UUID)');
    error.status = 400;
    throw error;
  }
  // Cho phép partial body
  const result = await grpcClients.eventService.updateEvent({
    eventId,
    ...req.body,
  });
  return result;
};

export const saveEventDraftHandler = createSimpleHandler(saveEventDraft, 'event', 'saveEventDraft');

export const saveEventLayoutHandler = async (req, res) => {
  res.status(501).json({ message: 'Not implemented: saveEventLayoutHandler' });
};

export const saveEventPricingHandler = async (req, res) => {
  res.status(501).json({ message: 'Not implemented: saveEventPricingHandler' });
};

export const publishEventHandler = async (req, res) => {
  res.status(501).json({ message: 'Not implemented: publishEventHandler' });
};

export const getEventTemplatesHandler = async (req, res) => {
  res.status(501).json({ message: 'Not implemented: getEventTemplatesHandler' });
};

export const duplicateEventHandler = async (req, res) => {
  res.status(501).json({ message: 'Not implemented: duplicateEventHandler' });
};

export const getEventsHandler = createSimpleHandler(getAllEvents, 'event', 'getEvents');
export const getEventHandler = createSimpleHandler(getEventById, 'event', 'getEvent');
export const createEventHandler = createSimpleHandler(createNewEvent, 'event', 'createEvent');
export const updateEventHandler = createSimpleHandler(updateEvent, 'event', 'updateEvent');
export const deleteEventHandler = createSimpleHandler(deleteEvent, 'event', 'deleteEvent');

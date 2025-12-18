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

/**
 * Publish event (change status to published)
 */
const publishEvent = async (req, res) => {
  const { eventId } = req.params;

  // First get the event to validate it exists
  const eventResult = await grpcClients.eventService.getEvent({ id: eventId });

  if (!eventResult.event) {
    const error = new Error('Event not found');
    error.status = 404;
    throw error;
  }

  // Update event status to published
  const result = await grpcClients.eventService.updateEvent({
    id: eventId,
    status: 'published',
  });

  sendSuccessResponse(res, 200, {
    message: 'Event published successfully',
    event: result.event,
  }, req.correlationId);
};

export const publishEventHandler = createSimpleHandler(publishEvent, 'event', 'publishEvent');

/**
 * Get event templates
 * Returns a list of predefined event templates
 */
const getEventTemplates = async (req, res) => {
  // Static templates - in production, these could come from database or config
  const templates = [
    {
      id: 'concert-general',
      name: 'Concert - General Admission',
      description: 'Standard concert layout with standing area',
      eventType: 'concert',
      category: 'music',
      defaultCapacity: 5000,
      layoutPreset: 'standing',
    },
    {
      id: 'concert-seated',
      name: 'Concert - Seated',
      description: 'Concert with assigned seating',
      eventType: 'concert',
      category: 'music',
      defaultCapacity: 2000,
      layoutPreset: 'theater',
    },
    {
      id: 'theater-standard',
      name: 'Theater - Standard',
      description: 'Traditional theater seating arrangement',
      eventType: 'theater',
      category: 'arts',
      defaultCapacity: 500,
      layoutPreset: 'theater',
    },
    {
      id: 'sports-stadium',
      name: 'Sports - Stadium',
      description: 'Stadium layout for sports events',
      eventType: 'sports',
      category: 'sports',
      defaultCapacity: 50000,
      layoutPreset: 'stadium',
    },
    {
      id: 'conference-hall',
      name: 'Conference - Hall',
      description: 'Conference or seminar layout',
      eventType: 'conference',
      category: 'business',
      defaultCapacity: 300,
      layoutPreset: 'conference',
    },
  ];

  sendSuccessResponse(res, 200, { templates }, req.correlationId);
};

export const getEventTemplatesHandler = createSimpleHandler(getEventTemplates, 'event', 'getEventTemplates');

/**
 * Duplicate event
 * Creates a copy of an existing event
 */
const duplicateEvent = async (req, res) => {
  const { eventId } = req.params;
  const { newName, newStartDate, newEndDate } = req.body;

  // Get the original event
  const originalResult = await grpcClients.eventService.getEvent({ id: eventId });

  if (!originalResult.event) {
    const error = new Error('Event not found');
    error.status = 404;
    throw error;
  }

  const original = originalResult.event;

  // Create a new event based on the original
  const newEvent = {
    name: newName || `${original.name} (Copy)`,
    description: original.description,
    startDate: newStartDate || original.startDate,
    endDate: newEndDate || original.endDate,
    venueName: original.venueName,
    venueAddress: original.venueAddress,
    venueCity: original.venueCity,
    venueCountry: original.venueCountry,
    venueCapacity: original.venueCapacity,
    canvasConfig: original.canvasConfig,
    organizationId: original.organizationId,
    // Reset status to draft for the copy
    status: 'draft',
  };

  const result = await grpcClients.eventService.createEvent({ event: newEvent });

  sendSuccessResponse(res, 201, {
    message: 'Event duplicated successfully',
    event: result.event,
    originalEventId: eventId,
  }, req.correlationId);
};

export const duplicateEventHandler = createSimpleHandler(duplicateEvent, 'event', 'duplicateEvent');

export const getEventsHandler = createSimpleHandler(getAllEvents, 'event', 'getEvents');
export const getEventHandler = createSimpleHandler(getEventById, 'event', 'getEvent');
export const createEventHandler = createSimpleHandler(createNewEvent, 'event', 'createEvent');
export const updateEventHandler = createSimpleHandler(updateEvent, 'event', 'updateEvent');
export const deleteEventHandler = createSimpleHandler(deleteEvent, 'event', 'deleteEvent');

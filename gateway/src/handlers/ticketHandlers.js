import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler, createHandler } from '../utils/responseHandler.js';

// ============================================
// Ticket CRUD
// ============================================

const getTickets = async (req, res) => {
  const result = await grpcClients.ticketService.GetTickets({
    event_id: req.query.event_id || '',
    ticket_type_id: req.query.ticket_type_id || '',
    status: req.query.status || '',
    page: parseInt(req.query.page) || 1,
    limit: parseInt(req.query.limit) || 50,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const getTicketById = async (req, res) => {
  const result = await grpcClients.ticketService.GetTicket({
    ticket_id: req.params.ticketId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const createTicket = async (req, res) => {
  const result = await grpcClients.ticketService.CreateTicket(req.body);
  sendSuccessResponse(res, 201, result, req.correlationId);
};

const updateTicket = async (req, res) => {
  const result = await grpcClients.ticketService.UpdateTicket({
    ticket_id: req.params.ticketId,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const deleteTicket = async (req, res) => {
  const result = await grpcClients.ticketService.DeleteTicket({
    ticket_id: req.params.ticketId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Ticket Types
// ============================================

const getTicketTypes = async (req, res) => {
  const result = await grpcClients.ticketService.GetTicketTypes({
    event_id: req.params.eventId,
    include_availability: req.query.include_availability === 'true',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const createTicketType = async (req, res) => {
  const result = await grpcClients.ticketService.CreateTicketType({
    event_id: req.body.event_id,
    ...req.body,
  });
  sendSuccessResponse(res, 201, result, req.correlationId);
};

const updateTicketType = async (req, res) => {
  const result = await grpcClients.ticketService.UpdateTicketType({
    ticket_type_id: req.params.typeId,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const deleteTicketType = async (req, res) => {
  const result = await grpcClients.ticketService.DeleteTicketType({
    ticket_type_id: req.params.typeId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Availability
// ============================================

const checkAvailability = async (req, res) => {
  const result = await grpcClients.ticketService.CheckAvailability({
    event_id: req.params.eventId,
    ticket_type_id: req.query.ticket_type_id || '',
    quantity: parseInt(req.query.quantity) || 1,
    seat_numbers: req.query.seat_numbers ? req.query.seat_numbers.split(',') : [],
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const getAvailableTickets = async (req, res) => {
  const result = await grpcClients.ticketService.GetAvailableTickets({
    event_id: req.params.eventId,
    ticket_type_id: req.query.ticket_type_id || '',
    limit: parseInt(req.query.limit) || 50,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Reservation
// ============================================

const reserveTickets = async (req, res) => {
  const result = await grpcClients.ticketService.ReserveTickets({
    event_id: req.params.eventId,
    user_id: req.user.id,
    ticket_type_id: req.body.ticket_type_id,
    quantity: req.body.quantity || 1,
    seat_numbers: req.body.seat_numbers || [],
    timeout_seconds: req.body.timeout_seconds || 600,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const releaseTickets = async (req, res) => {
  const result = await grpcClients.ticketService.ReleaseTickets({
    reservation_id: req.body.reservation_id,
    ticket_ids: req.body.ticket_ids || [],
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Export handlers
// ============================================

// Ticket CRUD
export const getTicketsHandler = createSimpleHandler(getTickets, 'ticket', 'getTickets');
export const getTicketHandler = createSimpleHandler(getTicketById, 'ticket', 'getTicket');
export const createTicketHandler = createHandler(createTicket, 'ticket', 'createTicket');
export const updateTicketHandler = createHandler(updateTicket, 'ticket', 'updateTicket');
export const deleteTicketHandler = createSimpleHandler(deleteTicket, 'ticket', 'deleteTicket');

// Ticket Types
export const getTicketTypesHandler = createSimpleHandler(getTicketTypes, 'ticket', 'getTicketTypes');
export const createTicketTypeHandler = createHandler(createTicketType, 'ticket', 'createTicketType');
export const updateTicketTypeHandler = createHandler(updateTicketType, 'ticket', 'updateTicketType');
export const deleteTicketTypeHandler = createSimpleHandler(deleteTicketType, 'ticket', 'deleteTicketType');

// Availability
export const checkAvailabilityHandler = createSimpleHandler(checkAvailability, 'ticket', 'checkAvailability');
export const getAvailableTicketsHandler = createSimpleHandler(getAvailableTickets, 'ticket', 'getAvailableTickets');

// Reservation
export const reserveTicketsHandler = createHandler(reserveTickets, 'ticket', 'reserveTickets');
export const releaseTicketsHandler = createHandler(releaseTickets, 'ticket', 'releaseTickets');

import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

const getTickets = async (req, res) => {
  const result = await grpcClients.ticketService.getTickets({});
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const getTicketById = async (req, res) => {
  const result = await grpcClients.ticketService.getTicket({ ticketId: req.params.ticketId });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const createTicket = async (req, res) => {
  const result = await grpcClients.ticketService.createTicket(req.body);
  sendSuccessResponse(res, 201, result, req.correlationId);
};

const updateTicket = async (req, res) => {
  const result = await grpcClients.ticketService.updateTicket({
    ticketId: req.params.ticketId,
    ...req.body,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const deleteTicket = async (req, res) => {
  await grpcClients.ticketService.deleteTicket({ ticketId: req.params.ticketId });
  sendSuccessResponse(res, 200, { message: 'Ticket deleted successfully' }, req.correlationId);
};

export const getTicketsHandler = createSimpleHandler(getTickets, 'ticket', 'getTickets');
export const getTicketHandler = createSimpleHandler(getTicketById, 'ticket', 'getTicket');
export const createTicketHandler = createSimpleHandler(createTicket, 'ticket', 'createTicket');
export const updateTicketHandler = createSimpleHandler(updateTicket, 'ticket', 'updateTicket');
export const deleteTicketHandler = createSimpleHandler(deleteTicket, 'ticket', 'deleteTicket');

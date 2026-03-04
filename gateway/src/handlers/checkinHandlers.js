import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler, createHandler } from '../utils/responseHandler.js';

// ============================================
// Check-in Operations
// ============================================

const checkIn = async (req, res) => {
  const result = await grpcClients.checkinService.CheckIn({
    ticket_id: req.body.ticket_id,
    qr_code: req.body.qr_code,
    event_id: req.params.eventId,
    staff_id: req.user.id,
    device_id: req.body.device_id || '',
    gate: req.body.gate || '',
  });
  sendSuccessResponse(res, 201, result, req.correlationId);
};

const getCheckIn = async (req, res) => {
  const result = await grpcClients.checkinService.GetCheckIn({
    checkin_id: req.params.checkinId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const listCheckIns = async (req, res) => {
  const result = await grpcClients.checkinService.ListCheckIns({
    event_id: req.params.eventId,
    user_id: req.query.user_id || '',
    gate: req.query.gate || '',
    page: parseInt(req.query.page) || 1,
    limit: parseInt(req.query.limit) || 20,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const getEventStats = async (req, res) => {
  const result = await grpcClients.checkinService.GetEventStats({
    event_id: req.params.eventId,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

// ============================================
// Exports
// ============================================

export const checkInHandler = createHandler(checkIn, 'checkin', 'checkIn');
export const getCheckInHandler = createSimpleHandler(getCheckIn, 'checkin', 'getCheckIn');
export const listCheckInsHandler = createSimpleHandler(listCheckIns, 'checkin', 'listCheckIns');
export const getEventStatsHandler = createSimpleHandler(getEventStats, 'checkin', 'getEventStats');

import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

/**
 * Create pricing for an event
 */
const createPricing = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.pricingService.CreatePricing({
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
 * Get pricing by ID
 */
const getPricing = async (req, res) => {
  const { pricingId } = req.params;
  const result = await grpcClients.pricingService.GetPricing({ pricing_id: pricingId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 404;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Update pricing
 */
const updatePricing = async (req, res) => {
  const { pricingId } = req.params;
  const result = await grpcClients.pricingService.UpdatePricing({
    pricing_id: pricingId,
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
 * Delete pricing
 */
const deletePricing = async (req, res) => {
  const { pricingId } = req.params;
  const result = await grpcClients.pricingService.DeletePricing({ pricing_id: pricingId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, { message: 'Pricing deleted successfully' }, req.correlationId);
};

/**
 * List pricing by event
 */
const listPricingByEvent = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.pricingService.GetPricingByEvent({ event_id: eventId });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 400;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Get pricing by zone
 */
const getPricingByZone = async (req, res) => {
  const { eventId, zoneId } = req.params;
  const result = await grpcClients.pricingService.GetPricingByZone({
    event_id: eventId,
    zone_id: zoneId,
  });

  if (result.error) {
    const error = new Error(result.error);
    error.status = 404;
    throw error;
  }

  sendSuccessResponse(res, 200, result, req.correlationId);
};

/**
 * Calculate price for seats
 */
const calculatePrice = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.pricingService.CalculatePrice({
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
 * Apply discount to pricing
 */
const applyDiscount = async (req, res) => {
  const { eventId } = req.params;
  const result = await grpcClients.pricingService.ApplyDiscount({
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

export const createPricingHandler = createSimpleHandler(createPricing, 'pricing', 'createPricing');
export const getPricingHandler = createSimpleHandler(getPricing, 'pricing', 'getPricing');
export const updatePricingHandler = createSimpleHandler(updatePricing, 'pricing', 'updatePricing');
export const deletePricingHandler = createSimpleHandler(deletePricing, 'pricing', 'deletePricing');
export const listPricingByEventHandler = createSimpleHandler(listPricingByEvent, 'pricing', 'listPricingByEvent');
export const getPricingByZoneHandler = createSimpleHandler(getPricingByZone, 'pricing', 'getPricingByZone');
export const calculatePriceHandler = createSimpleHandler(calculatePrice, 'pricing', 'calculatePrice');
export const applyDiscountHandler = createSimpleHandler(applyDiscount, 'pricing', 'applyDiscount');

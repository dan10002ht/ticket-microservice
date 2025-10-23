package grpc

import (
	"context"
	"event-service/services"
	eventpb "shared-lib/protos/event"
)

type AdvancedPricingController struct {
	service *services.PricingService
	eventpb.UnimplementedPricingServiceServer
}

func NewAdvancedPricingController(service *services.PricingService) *AdvancedPricingController {
	return &AdvancedPricingController{service: service}
}

// CalculatePrice - Calculate price for seats
func (c *AdvancedPricingController) CalculatePrice(ctx context.Context, req *eventpb.CalculatePriceRequest) (*eventpb.CalculatePriceResponse, error) {
	basePrice, finalPrice, discountAmount, reason, err := c.service.CalculatePrice(ctx, req.EventId, req.ZoneId, int(req.Quantity), req.DiscountCode, req.UserId)
	if err != nil {
		return &eventpb.CalculatePriceResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.CalculatePriceResponse{
		Success:        true,
		BasePrice:      basePrice,
		FinalPrice:     finalPrice,
		DiscountAmount: discountAmount,
		DiscountReason: reason,
		Currency:       req.Currency,
		Quantity:       req.Quantity,
	}, nil
}

// GetPricingByEvent - Get all pricing for an event
func (c *AdvancedPricingController) GetPricingByEvent(ctx context.Context, req *eventpb.GetPricingByEventRequest) (*eventpb.GetPricingByEventResponse, error) {
	pricing, err := c.service.GetPricingByEvent(ctx, req.EventId)
	if err != nil {
		return &eventpb.GetPricingByEventResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	var pbPricing []*eventpb.EventPricing
	for _, p := range pricing {
		pbPricing = append(pbPricing, &eventpb.EventPricing{
			Id:          p.PublicID,
			EventId:     p.EventID,
			ZoneId:      p.ZoneID,
			Price:       p.Price,
			Currency:    p.Currency,
			PricingType: p.PricingType,
			ValidFrom:   p.ValidFrom,
			ValidUntil:  p.ValidUntil,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return &eventpb.GetPricingByEventResponse{
		Success: true,
		Pricing: pbPricing,
		Total:   int32(len(pbPricing)),
	}, nil
}

// GetPricingByZone - Get pricing for specific zone
func (c *AdvancedPricingController) GetPricingByZone(ctx context.Context, req *eventpb.GetPricingByZoneRequest) (*eventpb.GetPricingByZoneResponse, error) {
	pricing, err := c.service.GetPricingByZone(ctx, req.EventId, req.ZoneId)
	if err != nil {
		return &eventpb.GetPricingByZoneResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	var pbPricing []*eventpb.EventPricing
	for _, p := range pricing {
		pbPricing = append(pbPricing, &eventpb.EventPricing{
			Id:          p.PublicID,
			EventId:     p.EventID,
			ZoneId:      p.ZoneID,
			Price:       p.Price,
			Currency:    p.Currency,
			PricingType: p.PricingType,
			ValidFrom:   p.ValidFrom,
			ValidUntil:  p.ValidUntil,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return &eventpb.GetPricingByZoneResponse{
		Success: true,
		Pricing: pbPricing,
		Total:   int32(len(pbPricing)),
	}, nil
}

// ApplyDiscount - Apply discount to price
func (c *AdvancedPricingController) ApplyDiscount(ctx context.Context, req *eventpb.ApplyDiscountRequest) (*eventpb.ApplyDiscountResponse, error) {
	finalPrice, discountAmount, reason, err := c.service.ApplyDiscount(ctx, req.EventId, req.BasePrice, req.DiscountCode, req.UserId)
	if err != nil {
		return &eventpb.ApplyDiscountResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.ApplyDiscountResponse{
		Success:        true,
		BasePrice:      req.BasePrice,
		FinalPrice:     finalPrice,
		DiscountAmount: discountAmount,
		DiscountReason: reason,
		Currency:       req.Currency,
	}, nil
}

package grpc

import (
	"context"
	"event-service/services"
	eventpb "event-service/internal/protos/event"
)

type AdvancedPricingController struct {
	service *services.PricingService
	eventpb.UnimplementedPricingServiceServer
}

func NewAdvancedPricingController(service *services.PricingService) *AdvancedPricingController {
	return &AdvancedPricingController{service: service}
}

// CreatePricing - Create new pricing
func (c *AdvancedPricingController) CreatePricing(ctx context.Context, req *eventpb.CreatePricingRequest) (*eventpb.CreatePricingResponse, error) {
	pricing, err := c.service.CreatePricing(ctx, req.EventId, req.ZoneId, req.PricingCategory, req.BasePrice, req.Currency, req.PricingRules, req.DiscountRules, req.ValidFrom, req.ValidUntil, req.CreatedBy)
	if err != nil {
		return &eventpb.CreatePricingResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.CreatePricingResponse{
		Pricing: &eventpb.EventPricing{
			Id:              pricing.PublicID,
			EventId:         pricing.EventID,
			ZoneId:          pricing.ZoneID,
			PricingCategory: pricing.PricingCategory,
			BasePrice:       pricing.BasePrice,
			Currency:        pricing.Currency,
			PricingRules:    pricing.PricingRules,
			DiscountRules:   pricing.DiscountRules,
			IsActive:        pricing.IsActive,
			ValidFrom:       pricing.ValidFrom,
			ValidUntil:      pricing.ValidUntil,
			CreatedAt:       pricing.CreatedAt,
			UpdatedAt:       pricing.UpdatedAt,
			CreatedBy:       pricing.CreatedBy,
			UpdatedBy:       pricing.UpdatedBy,
		},
	}, nil
}

// GetPricing - Get pricing by ID
func (c *AdvancedPricingController) GetPricing(ctx context.Context, req *eventpb.GetPricingRequest) (*eventpb.GetPricingResponse, error) {
	pricing, err := c.service.GetPricing(ctx, req.Id)
	if err != nil {
		return &eventpb.GetPricingResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.GetPricingResponse{
		Pricing: &eventpb.EventPricing{
			Id:              pricing.PublicID,
			EventId:         pricing.EventID,
			ZoneId:          pricing.ZoneID,
			PricingCategory: pricing.PricingCategory,
			BasePrice:       pricing.BasePrice,
			Currency:        pricing.Currency,
			PricingRules:    pricing.PricingRules,
			DiscountRules:   pricing.DiscountRules,
			IsActive:        pricing.IsActive,
			ValidFrom:       pricing.ValidFrom,
			ValidUntil:      pricing.ValidUntil,
			CreatedAt:       pricing.CreatedAt,
			UpdatedAt:       pricing.UpdatedAt,
			CreatedBy:       pricing.CreatedBy,
			UpdatedBy:       pricing.UpdatedBy,
		},
	}, nil
}

// UpdatePricing - Update pricing
func (c *AdvancedPricingController) UpdatePricing(ctx context.Context, req *eventpb.UpdatePricingRequest) (*eventpb.UpdatePricingResponse, error) {
	pricing, err := c.service.UpdatePricing(ctx, req.Id, req.BasePrice, req.Currency, req.PricingRules, req.DiscountRules, req.IsActive, req.ValidFrom, req.ValidUntil, req.UpdatedBy)
	if err != nil {
		return &eventpb.UpdatePricingResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.UpdatePricingResponse{
		Pricing: &eventpb.EventPricing{
			Id:              pricing.PublicID,
			EventId:         pricing.EventID,
			ZoneId:          pricing.ZoneID,
			PricingCategory: pricing.PricingCategory,
			BasePrice:       pricing.BasePrice,
			Currency:        pricing.Currency,
			PricingRules:    pricing.PricingRules,
			DiscountRules:   pricing.DiscountRules,
			IsActive:        pricing.IsActive,
			ValidFrom:       pricing.ValidFrom,
			ValidUntil:      pricing.ValidUntil,
			CreatedAt:       pricing.CreatedAt,
			UpdatedAt:       pricing.UpdatedAt,
			CreatedBy:       pricing.CreatedBy,
			UpdatedBy:       pricing.UpdatedBy,
		},
	}, nil
}

// DeletePricing - Delete pricing
func (c *AdvancedPricingController) DeletePricing(ctx context.Context, req *eventpb.DeletePricingRequest) (*eventpb.DeletePricingResponse, error) {
	err := c.service.DeletePricing(ctx, req.Id)
	if err != nil {
		return &eventpb.DeletePricingResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.DeletePricingResponse{
		Success: true,
	}, nil
}

// ListPricing - List pricing with pagination
func (c *AdvancedPricingController) ListPricing(ctx context.Context, req *eventpb.ListPricingRequest) (*eventpb.ListPricingResponse, error) {
	pricings, total, err := c.service.ListPricing(ctx, req.EventId, req.IsActive, req.Page, req.Limit)
	if err != nil {
		return &eventpb.ListPricingResponse{
			Error: err.Error(),
		}, nil
	}

	var pbPricing []*eventpb.EventPricing
	for _, p := range pricings {
		pbPricing = append(pbPricing, &eventpb.EventPricing{
			Id:              p.PublicID,
			EventId:         p.EventID,
			ZoneId:          p.ZoneID,
			PricingCategory: p.PricingCategory,
			BasePrice:       p.BasePrice,
			Currency:        p.Currency,
			PricingRules:    p.PricingRules,
			DiscountRules:   p.DiscountRules,
			IsActive:        p.IsActive,
			ValidFrom:       p.ValidFrom,
			ValidUntil:      p.ValidUntil,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
			CreatedBy:       p.CreatedBy,
			UpdatedBy:       p.UpdatedBy,
		})
	}

	return &eventpb.ListPricingResponse{
		Pricing: pbPricing,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// CalculatePrice - Calculate price for seats
func (c *AdvancedPricingController) CalculatePrice(ctx context.Context, req *eventpb.CalculatePriceRequest) (*eventpb.CalculatePriceResponse, error) {
	basePrice, finalPrice, discountAmount, discountReason, currency, pricingDetails, err := c.service.CalculatePrice(ctx, req.EventId, req.ZoneId, req.PricingCategory, req.Quantity, req.DiscountCode, req.UserId)
	if err != nil {
		return &eventpb.CalculatePriceResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.CalculatePriceResponse{
		BasePrice:      basePrice,
		FinalPrice:     finalPrice,
		DiscountAmount: discountAmount,
		DiscountReason: discountReason,
		Currency:       currency,
		PricingDetails: pricingDetails,
	}, nil
}

// GetPricingByEvent - Get all pricing for an event
func (c *AdvancedPricingController) GetPricingByEvent(ctx context.Context, req *eventpb.GetPricingByEventRequest) (*eventpb.GetPricingByEventResponse, error) {
	pricings, err := c.service.GetPricingByEvent(ctx, req.EventId, req.IsActive)
	if err != nil {
		return &eventpb.GetPricingByEventResponse{
			Error: err.Error(),
		}, nil
	}

	var pbPricing []*eventpb.EventPricing
	for _, p := range pricings {
		pbPricing = append(pbPricing, &eventpb.EventPricing{
			Id:              p.PublicID,
			EventId:         p.EventID,
			ZoneId:          p.ZoneID,
			PricingCategory: p.PricingCategory,
			BasePrice:       p.BasePrice,
			Currency:        p.Currency,
			PricingRules:    p.PricingRules,
			DiscountRules:   p.DiscountRules,
			IsActive:        p.IsActive,
			ValidFrom:       p.ValidFrom,
			ValidUntil:      p.ValidUntil,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
			CreatedBy:       p.CreatedBy,
			UpdatedBy:       p.UpdatedBy,
		})
	}

	return &eventpb.GetPricingByEventResponse{
		Pricing: pbPricing,
	}, nil
}

// GetPricingByZone - Get pricing for specific zone
func (c *AdvancedPricingController) GetPricingByZone(ctx context.Context, req *eventpb.GetPricingByZoneRequest) (*eventpb.GetPricingByZoneResponse, error) {
	pricings, err := c.service.GetPricingByZone(ctx, req.EventId, req.ZoneId, req.IsActive)
	if err != nil {
		return &eventpb.GetPricingByZoneResponse{
			Error: err.Error(),
		}, nil
	}

	var pbPricing []*eventpb.EventPricing
	for _, p := range pricings {
		pbPricing = append(pbPricing, &eventpb.EventPricing{
			Id:              p.PublicID,
			EventId:         p.EventID,
			ZoneId:          p.ZoneID,
			PricingCategory: p.PricingCategory,
			BasePrice:       p.BasePrice,
			Currency:        p.Currency,
			PricingRules:    p.PricingRules,
			DiscountRules:   p.DiscountRules,
			IsActive:        p.IsActive,
			ValidFrom:       p.ValidFrom,
			ValidUntil:      p.ValidUntil,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
			CreatedBy:       p.CreatedBy,
			UpdatedBy:       p.UpdatedBy,
		})
	}

	return &eventpb.GetPricingByZoneResponse{
		Pricing: pbPricing,
	}, nil
}

// ApplyDiscount - Apply discount to price
func (c *AdvancedPricingController) ApplyDiscount(ctx context.Context, req *eventpb.ApplyDiscountRequest) (*eventpb.ApplyDiscountResponse, error) {
	discountAmount, finalPrice, reason, isValid, err := c.service.ApplyDiscount(ctx, req.EventId, req.OriginalPrice, req.DiscountCode, req.UserId)
	if err != nil {
		return &eventpb.ApplyDiscountResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.ApplyDiscountResponse{
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		DiscountReason: reason,
		IsValid:        isValid,
	}, nil
}

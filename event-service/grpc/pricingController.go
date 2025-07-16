package grpc

import (
	"context"
	"event-service/models"
	"event-service/services"
)

type PricingController struct {
	service *services.PricingService
}

func NewPricingController(service *services.PricingService) *PricingController {
	return &PricingController{service: service}
}

func (c *PricingController) CreatePricing(ctx context.Context, req *models.EventPricing) (*models.EventPricing, error) {
	if err := c.service.CreatePricing(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *PricingController) GetPricing(ctx context.Context, publicID string) (*models.EventPricing, error) {
	return c.service.GetPricing(ctx, publicID)
}

func (c *PricingController) UpdatePricing(ctx context.Context, req *models.EventPricing) (*models.EventPricing, error) {
	if err := c.service.UpdatePricing(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *PricingController) DeletePricing(ctx context.Context, publicID string) error {
	return c.service.DeletePricing(ctx, publicID)
}

func (c *PricingController) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventPricing, error) {
	return c.service.ListByEventID(ctx, eventID)
}

func (c *PricingController) ApplyDiscount(ctx context.Context, eventID int64, basePrice float64, discountCode string, userID string) (finalPrice float64, discountAmount float64, reason string, err error) {
	return c.service.ApplyDiscount(ctx, eventID, basePrice, discountCode, userID)
}

func (c *PricingController) CalculatePrice(ctx context.Context, eventID int64, zoneID int64, quantity int, discountCode string, userID string) (basePrice, finalPrice, discountAmount float64, reason string, err error) {
	return c.service.CalculatePrice(ctx, eventID, zoneID, quantity, discountCode, userID)
} 
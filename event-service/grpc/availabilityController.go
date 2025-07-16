package grpc

import (
	"context"
	"event-service/models"
	"event-service/services"
)

type AvailabilityController struct {
	service *services.AvailabilityService
}

func NewAvailabilityController(service *services.AvailabilityService) *AvailabilityController {
	return &AvailabilityController{service: service}
}

func (c *AvailabilityController) CreateAvailability(ctx context.Context, req *models.EventSeatAvailability) (*models.EventSeatAvailability, error) {
	if err := c.service.CreateAvailability(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *AvailabilityController) GetAvailability(ctx context.Context, publicID string) (*models.EventSeatAvailability, error) {
	return c.service.GetAvailability(ctx, publicID)
}

func (c *AvailabilityController) UpdateAvailability(ctx context.Context, req *models.EventSeatAvailability) (*models.EventSeatAvailability, error) {
	if err := c.service.UpdateAvailability(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *AvailabilityController) DeleteAvailability(ctx context.Context, publicID string) error {
	return c.service.DeleteAvailability(ctx, publicID)
}

func (c *AvailabilityController) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSeatAvailability, error) {
	return c.service.ListByEventID(ctx, eventID)
} 
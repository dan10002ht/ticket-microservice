package grpc

import (
	"context"
	"event-service/models"
	"event-service/services"
)

type ScheduleController struct {
	service *services.ScheduleService
}

func NewScheduleController(service *services.ScheduleService) *ScheduleController {
	return &ScheduleController{service: service}
}

func (c *ScheduleController) CreateSchedule(ctx context.Context, req *models.EventSchedule) (*models.EventSchedule, error) {
	if err := c.service.CreateSchedule(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *ScheduleController) GetSchedule(ctx context.Context, publicID string) (*models.EventSchedule, error) {
	return c.service.GetSchedule(ctx, publicID)
}

func (c *ScheduleController) UpdateSchedule(ctx context.Context, req *models.EventSchedule) (*models.EventSchedule, error) {
	if err := c.service.UpdateSchedule(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *ScheduleController) DeleteSchedule(ctx context.Context, publicID string) error {
	return c.service.DeleteSchedule(ctx, publicID)
}

func (c *ScheduleController) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSchedule, error) {
	return c.service.ListByEventID(ctx, eventID)
} 
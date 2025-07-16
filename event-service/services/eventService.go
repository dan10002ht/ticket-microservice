package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
)

type EventService struct {
	repo *repositories.EventRepository
}

func NewEventService(repo *repositories.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, event *models.Event) error {
	return s.repo.Create(ctx, event)
}

func (s *EventService) GetEvent(ctx context.Context, publicID string) (*models.Event, error) {
	return s.repo.GetByPublicID(ctx, publicID)
}

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return s.repo.Update(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, publicID)
}

func (s *EventService) ListByOrganizationID(ctx context.Context, organizationID int64) ([]*models.Event, error) {
	return s.repo.List(ctx, organizationID)
} 
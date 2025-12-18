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

func (s *EventService) ListByOrganizationID(ctx context.Context, organizationID string) ([]*models.Event, error) {
	return s.repo.ListByOrganizationID(ctx, organizationID)
}

func (s *EventService) List(ctx context.Context, page, limit int32) ([]*models.Event, error) {
	return s.repo.ListAll(ctx, page, limit)
}

// Advanced search and filtering methods
func (s *EventService) SearchEvents(ctx context.Context, query, eventType, category string, page, limit int32) ([]*models.Event, int, error) {
	return s.repo.SearchEvents(ctx, query, eventType, category, page, limit)
}

func (s *EventService) GetEventsByVenue(ctx context.Context, venueID, status string, page, limit int32) ([]*models.Event, int, error) {
	return s.repo.GetEventsByVenue(ctx, venueID, status, page, limit)
}

func (s *EventService) GetUpcomingEvents(ctx context.Context, days int32, eventType, category string, page, limit int32) ([]*models.Event, int, error) {
	return s.repo.GetUpcomingEvents(ctx, days, eventType, category, page, limit)
}

func (s *EventService) GetFeaturedEvents(ctx context.Context, eventType, category string, page, limit int32) ([]*models.Event, int, error) {
	return s.repo.GetFeaturedEvents(ctx, eventType, category, page, limit)
}

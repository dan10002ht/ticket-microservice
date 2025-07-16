package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
)

type EventSeatingZoneService struct {
	repo *repositories.EventSeatingZoneRepository
}

func NewEventSeatingZoneService(repo *repositories.EventSeatingZoneRepository) *EventSeatingZoneService {
	return &EventSeatingZoneService{repo: repo}
}

func (s *EventSeatingZoneService) CreateZone(ctx context.Context, zone *models.EventSeatingZone) error {
	return s.repo.Create(ctx, zone)
}

func (s *EventSeatingZoneService) GetZone(ctx context.Context, publicID string) (*models.EventSeatingZone, error) {
	return s.repo.GetByPublicID(ctx, publicID)
}

func (s *EventSeatingZoneService) UpdateZone(ctx context.Context, zone *models.EventSeatingZone) error {
	return s.repo.Update(ctx, zone)
}

func (s *EventSeatingZoneService) DeleteZone(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, publicID)
}

func (s *EventSeatingZoneService) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSeatingZone, error) {
	return s.repo.ListByEventID(ctx, eventID)
} 
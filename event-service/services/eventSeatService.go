package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
)

type EventSeatService struct {
	repo *repositories.EventSeatRepository
}

func NewEventSeatService(repo *repositories.EventSeatRepository) *EventSeatService {
	return &EventSeatService{repo: repo}
}

func (s *EventSeatService) CreateSeat(ctx context.Context, seat *models.EventSeat) error {
	return s.repo.Create(ctx, seat)
}

func (s *EventSeatService) GetSeat(ctx context.Context, publicID string) (*models.EventSeat, error) {
	return s.repo.GetByPublicID(ctx, publicID)
}

func (s *EventSeatService) UpdateSeat(ctx context.Context, seat *models.EventSeat) error {
	return s.repo.Update(ctx, seat)
}

func (s *EventSeatService) DeleteSeat(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, publicID)
}

func (s *EventSeatService) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSeat, error) {
	return s.repo.ListByEventID(ctx, eventID)
} 
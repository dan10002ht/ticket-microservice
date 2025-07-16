package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
	"fmt"

	"github.com/google/uuid"
)

type AvailabilityService struct {
	repo *repositories.EventSeatAvailabilityRepository
}

func NewAvailabilityService(repo *repositories.EventSeatAvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{repo: repo}
}

func (s *AvailabilityService) CreateAvailability(ctx context.Context, avail *models.EventSeatAvailability) error {
	if err := s.ValidateAvailability(avail); err != nil {
		return err
	}
	return s.repo.Create(ctx, avail)
}

func (s *AvailabilityService) GetAvailability(ctx context.Context, publicID string) (*models.EventSeatAvailability, error) {
	return s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
}

func (s *AvailabilityService) UpdateAvailability(ctx context.Context, avail *models.EventSeatAvailability) error {
	if err := s.ValidateAvailability(avail); err != nil {
		return err
	}
	return s.repo.Update(ctx, avail)
}

func (s *AvailabilityService) DeleteAvailability(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, uuid.MustParse(publicID))
}

func (s *AvailabilityService) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSeatAvailability, error) {
	return s.repo.ListByEventID(ctx, eventID)
}

func (s *AvailabilityService) ValidateAvailability(avail *models.EventSeatAvailability) error {
	if avail.EventID == 0 || avail.SeatID == 0 {
		return fmt.Errorf("invalid seat availability data")
	}
	return nil
} 
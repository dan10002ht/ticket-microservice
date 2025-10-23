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

// Advanced availability methods
func (s *AvailabilityService) GetEventAvailability(ctx context.Context, eventID string) (*models.EventAvailability, error) {
	return s.repo.GetEventAvailability(ctx, eventID)
}

func (s *AvailabilityService) GetZoneAvailability(ctx context.Context, eventID, zoneID string) (*models.ZoneAvailability, error) {
	return s.repo.GetZoneAvailability(ctx, eventID, zoneID)
}

func (s *AvailabilityService) GetSeatAvailability(ctx context.Context, eventID, seatID string) (*models.SeatAvailability, error) {
	return s.repo.GetSeatAvailability(ctx, eventID, seatID)
}

func (s *AvailabilityService) UpdateSeatAvailability(ctx context.Context, eventID, seatID, status, userID, bookingID string) error {
	return s.repo.UpdateSeatAvailability(ctx, eventID, seatID, status, userID, bookingID)
}

func (s *AvailabilityService) BlockSeats(ctx context.Context, eventID string, seatIDs []string, userID, bookingID, expiresAt string) (*models.BlockSeatsResult, error) {
	return s.repo.BlockSeats(ctx, eventID, seatIDs, userID, bookingID, expiresAt)
}

func (s *AvailabilityService) ReleaseSeats(ctx context.Context, eventID string, seatIDs []string, userID, reason string) (*models.ReleaseSeatsResult, error) {
	return s.repo.ReleaseSeats(ctx, eventID, seatIDs, userID, reason)
}

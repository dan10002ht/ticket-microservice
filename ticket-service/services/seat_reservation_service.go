package services

import (
	"context"
	"ticket-service/models"
	"ticket-service/repositories"
)

type SeatReservationService struct {
	repo *repositories.SeatReservationRepository
}

func NewSeatReservationService(repo *repositories.SeatReservationRepository) *SeatReservationService {
	return &SeatReservationService{repo: repo}
}

func (s *SeatReservationService) Create(ctx context.Context, reservation *models.SeatReservation) error {
	return s.repo.Create(ctx, reservation)
}

func (s *SeatReservationService) GetByID(ctx context.Context, id string) (*models.SeatReservation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SeatReservationService) Update(ctx context.Context, reservation *models.SeatReservation) error {
	return s.repo.Update(ctx, reservation)
}

func (s *SeatReservationService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *SeatReservationService) ListByEventID(ctx context.Context, eventID string) ([]*models.SeatReservation, error) {
	return s.repo.ListByEventID(ctx, eventID)
} 
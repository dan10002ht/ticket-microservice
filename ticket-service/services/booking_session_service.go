package services

import (
	"context"
	"ticket-service/models"
	"ticket-service/repositories"
)

type BookingSessionService struct {
	repo *repositories.BookingSessionRepository
}

func NewBookingSessionService(repo *repositories.BookingSessionRepository) *BookingSessionService {
	return &BookingSessionService{repo: repo}
}

func (s *BookingSessionService) Create(ctx context.Context, session *models.BookingSession) error {
	return s.repo.Create(ctx, session)
}

func (s *BookingSessionService) GetByID(ctx context.Context, id string) (*models.BookingSession, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookingSessionService) Update(ctx context.Context, session *models.BookingSession) error {
	return s.repo.Update(ctx, session)
}

func (s *BookingSessionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *BookingSessionService) ListByUserID(ctx context.Context, userID string) ([]*models.BookingSession, error) {
	return s.repo.ListByUserID(ctx, userID)
} 
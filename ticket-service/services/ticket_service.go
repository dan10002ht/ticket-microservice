package services

import (
	"context"
	"ticket-service/models"
	"ticket-service/repositories"
)

type TicketService struct {
	repo *repositories.TicketRepository
}

func NewTicketService(repo *repositories.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) Create(ctx context.Context, ticket *models.Ticket) error {
	return s.repo.Create(ctx, ticket)
}

func (s *TicketService) GetByID(ctx context.Context, id string) (*models.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TicketService) Update(ctx context.Context, ticket *models.Ticket) error {
	return s.repo.Update(ctx, ticket)
}

func (s *TicketService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *TicketService) ListByEventID(ctx context.Context, eventID string) ([]*models.Ticket, error) {
	return s.repo.ListByEventID(ctx, eventID)
} 
package services

import (
	"context"
	"fmt"
	"venue-service/database"
	"venue-service/grpcclient"
	"venue-service/models"
	"venue-service/queue"
	"venue-service/repositories"

	"github.com/google/uuid"
)

type SeatService struct {
	db          *database.Database
	redis       *queue.RedisClient
	grpcClients *grpcclient.Clients
	repo        *repositories.SeatRepository
}

func NewSeatService(db *database.Database, redis *queue.RedisClient, grpcClients *grpcclient.Clients) *SeatService {
	return &SeatService{
		db:          db,
		redis:       redis,
		grpcClients: grpcClients,
		repo:        repositories.NewSeatRepository(db.DB, nil), // TODO: pass logger
	}
}

func (s *SeatService) CreateSeat(ctx context.Context, seat *models.Seat) error {
	if err := s.ValidateSeat(seat); err != nil {
		return err
	}
	return s.repo.Create(ctx, seat)
}

func (s *SeatService) GetSeat(ctx context.Context, publicID string) (*models.Seat, error) {
	return s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
}

func (s *SeatService) UpdateSeat(ctx context.Context, seat *models.Seat) error {
	if err := s.ValidateSeat(seat); err != nil {
		return err
	}
	return s.repo.Update(ctx, seat)
}

func (s *SeatService) DeleteSeat(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, uuid.MustParse(publicID))
}

func (s *SeatService) BulkCreateSeats(ctx context.Context, seats []*models.Seat) error {
	for _, seat := range seats {
		if err := s.ValidateSeat(seat); err != nil {
			return err
		}
	}
	return s.repo.BulkCreate(ctx, seats)
}

func (s *SeatService) ValidateSeat(seat *models.Seat) error {
	if seat.SeatNumber == "" || seat.ZoneID == 0 {
		return fmt.Errorf("invalid seat data")
	}
	// TODO: Validate coordinates JSON, uniqueness, etc.
	return nil
} 
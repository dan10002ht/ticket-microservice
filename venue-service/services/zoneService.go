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

type ZoneService struct {
	db          *database.Database
	redis       *queue.RedisClient
	grpcClients *grpcclient.Clients
	repo        *repositories.SeatingZoneRepository
}

func NewZoneService(db *database.Database, redis *queue.RedisClient, grpcClients *grpcclient.Clients) *ZoneService {
	return &ZoneService{
		db:          db,
		redis:       redis,
		grpcClients: grpcClients,
		repo:        repositories.NewSeatingZoneRepository(db.DB, nil), // TODO: pass logger
	}
}

func (s *ZoneService) CreateZone(ctx context.Context, zone *models.SeatingZone) error {
	if err := s.ValidateZone(zone); err != nil {
		return err
	}
	return s.repo.Create(ctx, zone)
}

func (s *ZoneService) GetZone(ctx context.Context, publicID string) (*models.SeatingZone, error) {
	return s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
}

func (s *ZoneService) UpdateZone(ctx context.Context, zone *models.SeatingZone) error {
	if err := s.ValidateZone(zone); err != nil {
		return err
	}
	return s.repo.Update(ctx, zone)
}

func (s *ZoneService) DeleteZone(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, uuid.MustParse(publicID))
}

func (s *ZoneService) ValidateZone(zone *models.SeatingZone) error {
	if zone.Name == "" || zone.LayoutID == 0 {
		return fmt.Errorf("invalid zone data")
	}
	// TODO: Validate coordinates JSON, seat_count, etc.
	return nil
} 
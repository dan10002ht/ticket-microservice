package services

import (
	"context"
	"encoding/json"
	"fmt"
	"venue-service/database"
	"venue-service/grpcclient"
	"venue-service/models"
	"venue-service/queue"
	"venue-service/repositories"

	"github.com/google/uuid"
)

type VenueService struct {
	db          *database.Database
	redis       *queue.RedisClient
	grpcClients *grpcclient.Clients
	repo        *repositories.VenueRepository
}

func NewVenueService(db *database.Database, redis *queue.RedisClient, grpcClients *grpcclient.Clients) *VenueService {
	return &VenueService{
		db:          db,
		redis:       redis,
		grpcClients: grpcClients,
		repo:        repositories.NewVenueRepository(db.DB, nil), // TODO: pass logger
	}
}

func (s *VenueService) CreateVenue(ctx context.Context, venue *models.Venue) error {
	if err := s.ValidateVenue(venue); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, venue); err != nil {
		return err
	}
	// Cache venue
	b, _ := json.Marshal(venue)
	s.redis.Set(ctx, s.venueCacheKey(venue.PublicID.String()), b, 0)
	return nil
}

func (s *VenueService) GetVenue(ctx context.Context, publicID string) (*models.Venue, error) {
	// Try cache
	val, err := s.redis.Get(ctx, s.venueCacheKey(publicID)).Result()
	if err == nil {
		var v models.Venue
		if json.Unmarshal([]byte(val), &v) == nil {
			return &v, nil
		}
	}
	// Fallback DB
	venue, err := s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(venue)
	s.redis.Set(ctx, s.venueCacheKey(publicID), b, 0)
	return venue, nil
}

func (s *VenueService) UpdateVenue(ctx context.Context, venue *models.Venue) error {
	if err := s.ValidateVenue(venue); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, venue); err != nil {
		return err
	}
	b, _ := json.Marshal(venue)
	s.redis.Set(ctx, s.venueCacheKey(venue.PublicID.String()), b, 0)
	return nil
}

func (s *VenueService) DeleteVenue(ctx context.Context, publicID string) error {
	if err := s.repo.Delete(ctx, uuid.MustParse(publicID)); err != nil {
		return err
	}
	s.redis.Del(ctx, s.venueCacheKey(publicID))
	return nil
}

func (s *VenueService) ValidateVenue(venue *models.Venue) error {
	if venue.Name == "" || venue.Address == "" || venue.Capacity <= 0 {
		return fmt.Errorf("invalid venue data")
	}
	return nil
}

func (s *VenueService) venueCacheKey(publicID string) string {
	return "venue:" + publicID
} 
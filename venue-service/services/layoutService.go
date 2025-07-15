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

type LayoutService struct {
	db          *database.Database
	redis       *queue.RedisClient
	grpcClients *grpcclient.Clients
	repo        *repositories.VenueLayoutRepository
}

func NewLayoutService(db *database.Database, redis *queue.RedisClient, grpcClients *grpcclient.Clients) *LayoutService {
	return &LayoutService{
		db:          db,
		redis:       redis,
		grpcClients: grpcClients,
		repo:        repositories.NewVenueLayoutRepository(db.DB, nil), // TODO: pass logger
	}
}

func (s *LayoutService) CreateLayout(ctx context.Context, layout *models.VenueLayout) error {
	if err := s.ValidateLayout(layout); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, layout); err != nil {
		return err
	}
	b, _ := json.Marshal(layout)
	s.redis.Set(ctx, s.layoutCacheKey(layout.PublicID.String()), b, 0)
	return nil
}

func (s *LayoutService) GetLayout(ctx context.Context, publicID string) (*models.VenueLayout, error) {
	val, err := s.redis.Get(ctx, s.layoutCacheKey(publicID)).Result()
	if err == nil {
		var l models.VenueLayout
		if json.Unmarshal([]byte(val), &l) == nil {
			return &l, nil
		}
	}
	layout, err := s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(layout)
	s.redis.Set(ctx, s.layoutCacheKey(publicID), b, 0)
	return layout, nil
}

func (s *LayoutService) UpdateLayout(ctx context.Context, layout *models.VenueLayout) error {
	if err := s.ValidateLayout(layout); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, layout); err != nil {
		return err
	}
	b, _ := json.Marshal(layout)
	s.redis.Set(ctx, s.layoutCacheKey(layout.PublicID.String()), b, 0)
	return nil
}

func (s *LayoutService) DeleteLayout(ctx context.Context, publicID string) error {
	if err := s.repo.Delete(ctx, uuid.MustParse(publicID)); err != nil {
		return err
	}
	s.redis.Del(ctx, s.layoutCacheKey(publicID))
	return nil
}

func (s *LayoutService) ValidateLayout(layout *models.VenueLayout) error {
	if layout.Name == "" || layout.VenueID == 0 {
		return fmt.Errorf("invalid layout data")
	}
	// TODO: Validate canvas_config JSON, seat_count, etc.
	return nil
}

func (s *LayoutService) layoutCacheKey(publicID string) string {
	return "layout:" + publicID
} 
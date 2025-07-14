package services

import (
	"venue-service/database"
	"venue-service/grpcclient"
	"venue-service/queue"
)

type Services struct {
	Database *database.Database
	Redis    *queue.RedisClient
	User     *UserService
	Admin    *AdminService
}

func NewServices(db *database.Database, redis *queue.RedisClient, grpcClients *grpcclient.Clients) *Services {
	return &Services{
		Database: db,
		Redis:    redis,
		User:     NewUserService(db, redis, grpcClients),
		Admin:    NewAdminService(db, redis, grpcClients),
	}
}

func (s *Services) Close() error {
	// Close database connection
	if s.Database != nil {
		if err := s.Database.Close(); err != nil {
			return err
		}
	}

	// Close Redis connection
	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil {
			return err
		}
	}

	return nil
} 
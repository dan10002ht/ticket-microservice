package queue

import (
	"context"
	"fmt"

	"venue-service/internal/config"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		PoolSize:     10,
		MinIdleConns: 5,
	})
	
	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	
	return &RedisClient{client}, nil
}

func (r *RedisClient) Ping() error {
	ctx := context.Background()
	return r.Client.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
} 
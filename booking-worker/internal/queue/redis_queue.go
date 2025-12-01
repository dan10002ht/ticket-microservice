package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"booking-worker/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisQueueManager implements QueueManager using Redis
type RedisQueueManager struct {
	client *redis.Client
	config *config.Config
	logger *zap.Logger
}

// NewRedisQueueManager creates a new Redis queue manager
func NewRedisQueueManager(cfg *config.Config, logger *zap.Logger) (QueueManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis",
		zap.String("host", cfg.Redis.Host),
		zap.Int("port", cfg.Redis.Port),
	)

	return &RedisQueueManager{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// Enqueue adds an item to the queue
func (r *RedisQueueManager) Enqueue(ctx context.Context, item *QueueItem) error {
	key := fmt.Sprintf("booking-queue:%s", item.EventID)
	
	// Serialize item
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal queue item: %w", err)
	}

	// Add to queue (LPUSH)
	if err := r.client.LPush(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("failed to enqueue item: %w", err)
	}

	// Add to sorted set for position tracking
	positionKey := fmt.Sprintf("booking-queue-positions:%s", item.EventID)
	score := float64(time.Now().UnixNano())
	if err := r.client.ZAdd(ctx, positionKey, redis.Z{
		Score:  score,
		Member: item.ID,
	}).Err(); err != nil {
		r.logger.Warn("Failed to update position tracking", zap.Error(err))
	}

	// Set expiry for timeout tracking
	timeoutKey := fmt.Sprintf("booking-timeouts:%s", item.EventID)
	if err := r.client.ZAdd(ctx, timeoutKey, redis.Z{
		Score:  float64(item.ExpiresAt.Unix()),
		Member: item.ID,
	}).Err(); err != nil {
		r.logger.Warn("Failed to update timeout tracking", zap.Error(err))
	}

	r.logger.Debug("Item enqueued",
		zap.String("item_id", item.ID),
		zap.String("event_id", item.EventID),
	)

	return nil
}

// Dequeue removes and returns an item from the queue (blocking)
func (r *RedisQueueManager) Dequeue(ctx context.Context, timeout time.Duration) (*QueueItem, error) {
	// TODO: Implement dequeue logic with BRPOP
	// For now, return error as this needs proper implementation
	return nil, fmt.Errorf("dequeue not yet implemented")
}

// GetPosition returns the position of an item in the queue
func (r *RedisQueueManager) GetPosition(ctx context.Context, itemID string) (int, error) {
	// TODO: Implement position tracking using sorted sets
	return 0, fmt.Errorf("get position not yet implemented")
}

// GetQueueLength returns the current queue length
func (r *RedisQueueManager) GetQueueLength(ctx context.Context, eventID string) (int, error) {
	key := fmt.Sprintf("booking-queue:%s", eventID)
	length, err := r.client.LLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return int(length), nil
}

// Remove removes an item from the queue
func (r *RedisQueueManager) Remove(ctx context.Context, itemID string) error {
	// TODO: Implement remove logic
	return fmt.Errorf("remove not yet implemented")
}

// Close closes the queue connection
func (r *RedisQueueManager) Close() error {
	return r.client.Close()
}


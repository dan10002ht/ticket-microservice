package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"booking-worker/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// Redis key for tracking active event queues (SET)
	activeQueuesKey = "booking:active-queues"
	// Redis key prefix for item ID to event ID mapping (HASH)
	itemEventMapKey = "booking:item-event-map"
)

// RedisQueueManager implements QueueManager using Redis
type RedisQueueManager struct {
	client         *redis.Client
	config         *config.Config
	logger         *zap.Logger
	roundRobinIdx  int
	roundRobinLock sync.Mutex
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

	// Use pipeline for atomic operations
	pipe := r.client.Pipeline()

	// Add to queue (LPUSH)
	pipe.LPush(ctx, key, data)

	// Register this event queue as active (SADD - O(1))
	pipe.SAdd(ctx, activeQueuesKey, item.EventID)

	// Map item ID to event ID for O(1) lookup (HSET)
	pipe.HSet(ctx, itemEventMapKey, item.ID, item.EventID)

	// Add to sorted set for position tracking
	positionKey := fmt.Sprintf("booking-queue-positions:%s", item.EventID)
	score := float64(time.Now().UnixNano())
	pipe.ZAdd(ctx, positionKey, redis.Z{
		Score:  score,
		Member: item.ID,
	})

	// Set expiry for timeout tracking
	timeoutKey := fmt.Sprintf("booking-timeouts:%s", item.EventID)
	pipe.ZAdd(ctx, timeoutKey, redis.Z{
		Score:  float64(item.ExpiresAt.Unix()),
		Member: item.ID,
	})

	// Execute pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to enqueue item: %w", err)
	}

	r.logger.Debug("Item enqueued",
		zap.String("item_id", item.ID),
		zap.String("event_id", item.EventID),
	)

	return nil
}

// Dequeue removes and returns an item from the queue (blocking)
// Uses round-robin to fairly process items from different event queues
func (r *RedisQueueManager) Dequeue(ctx context.Context, timeout time.Duration) (*QueueItem, error) {
	// Get active event queues from SET (O(N) but N is small - number of events, not items)
	eventIDs, err := r.client.SMembers(ctx, activeQueuesKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get active queues: %w", err)
	}

	if len(eventIDs) == 0 {
		// No queues available, wait a bit
		time.Sleep(500 * time.Millisecond)
		return nil, nil
	}

	// Round-robin: select next queue in rotation
	r.roundRobinLock.Lock()
	startIdx := r.roundRobinIdx
	r.roundRobinIdx = (r.roundRobinIdx + 1) % len(eventIDs)
	r.roundRobinLock.Unlock()

	// Try each queue starting from round-robin index
	for i := 0; i < len(eventIDs); i++ {
		idx := (startIdx + i) % len(eventIDs)
		eventID := eventIDs[idx]
		queueKey := fmt.Sprintf("booking-queue:%s", eventID)

		// Use RPOP (non-blocking) for fair distribution
		// BRPOP would block on first queue, starving others
		result, err := r.client.RPop(ctx, queueKey).Result()
		if err != nil {
			if err == redis.Nil {
				continue // Empty queue, try next
			}
			r.logger.Warn("Failed to pop from queue", zap.String("event_id", eventID), zap.Error(err))
			continue
		}

		// Parse the item
		var item QueueItem
		if err := json.Unmarshal([]byte(result), &item); err != nil {
			r.logger.Error("Failed to unmarshal queue item", zap.Error(err))
			continue
		}

		// Cleanup tracking data using pipeline
		r.cleanupItemTracking(ctx, &item)

		r.logger.Debug("Item dequeued",
			zap.String("item_id", item.ID),
			zap.String("event_id", item.EventID),
		)

		return &item, nil
	}

	// No items found in any queue, wait before next poll
	time.Sleep(100 * time.Millisecond)
	return nil, nil
}

// cleanupItemTracking removes item from all tracking structures
func (r *RedisQueueManager) cleanupItemTracking(ctx context.Context, item *QueueItem) {
	pipe := r.client.Pipeline()

	// Remove from position tracking
	positionKey := fmt.Sprintf("booking-queue-positions:%s", item.EventID)
	pipe.ZRem(ctx, positionKey, item.ID)

	// Remove from timeout tracking
	timeoutKey := fmt.Sprintf("booking-timeouts:%s", item.EventID)
	pipe.ZRem(ctx, timeoutKey, item.ID)

	// Remove from item-event mapping
	pipe.HDel(ctx, itemEventMapKey, item.ID)

	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.Warn("Failed to cleanup item tracking", zap.String("item_id", item.ID), zap.Error(err))
	}

	// Check if queue is now empty and remove from active queues
	go r.cleanupEmptyQueue(context.Background(), item.EventID)
}

// cleanupEmptyQueue removes event from active queues if its queue is empty
func (r *RedisQueueManager) cleanupEmptyQueue(ctx context.Context, eventID string) {
	queueKey := fmt.Sprintf("booking-queue:%s", eventID)
	length, err := r.client.LLen(ctx, queueKey).Result()
	if err != nil {
		return
	}

	if length == 0 {
		// Queue is empty, remove from active queues
		r.client.SRem(ctx, activeQueuesKey, eventID)
		r.logger.Debug("Removed empty queue from active queues", zap.String("event_id", eventID))
	}
}

// GetPosition returns the position of an item in the queue
func (r *RedisQueueManager) GetPosition(ctx context.Context, itemID string) (int, error) {
	// O(1) lookup: Get event ID from item-event mapping
	eventID, err := r.client.HGet(ctx, itemEventMapKey, itemID).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("item not found in any queue")
		}
		return 0, fmt.Errorf("failed to get item event mapping: %w", err)
	}

	// O(log N) lookup: Get position from sorted set
	positionKey := fmt.Sprintf("booking-queue-positions:%s", eventID)
	rank, err := r.client.ZRank(ctx, positionKey, itemID).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("item not found in queue")
		}
		return 0, fmt.Errorf("failed to get rank: %w", err)
	}

	// Return rank (0-indexed, so add 1 for human-readable position)
	return int(rank) + 1, nil
}

// GetPositionWithEventID returns position when event ID is known (faster)
func (r *RedisQueueManager) GetPositionWithEventID(ctx context.Context, itemID, eventID string) (int, error) {
	positionKey := fmt.Sprintf("booking-queue-positions:%s", eventID)
	rank, err := r.client.ZRank(ctx, positionKey, itemID).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("item not found in queue")
		}
		return 0, fmt.Errorf("failed to get rank: %w", err)
	}
	return int(rank) + 1, nil
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

// Remove removes an item from the queue by item ID
// Uses Lua script for atomic removal - O(N) but only scans specific event queue
func (r *RedisQueueManager) Remove(ctx context.Context, itemID string) error {
	// O(1) lookup: Get event ID from item-event mapping
	eventID, err := r.client.HGet(ctx, itemEventMapKey, itemID).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("item not found in any queue")
		}
		return fmt.Errorf("failed to get item event mapping: %w", err)
	}

	return r.RemoveWithEventID(ctx, itemID, eventID)
}

// RemoveWithEventID removes an item when event ID is known (faster, avoids hash lookup)
func (r *RedisQueueManager) RemoveWithEventID(ctx context.Context, itemID, eventID string) error {
	queueKey := fmt.Sprintf("booking-queue:%s", eventID)

	// Lua script for atomic find-and-remove operation
	// This is more efficient than LRANGE + iterate + LSET + LREM in Go
	luaScript := redis.NewScript(`
		local queue_key = KEYS[1]
		local item_id = ARGV[1]
		local items = redis.call('LRANGE', queue_key, 0, -1)

		for i, item_data in ipairs(items) do
			local item = cjson.decode(item_data)
			if item.ID == item_id then
				-- Found the item, remove it using LSET + LREM pattern
				redis.call('LSET', queue_key, i - 1, '___REMOVED___')
				redis.call('LREM', queue_key, 1, '___REMOVED___')
				return 1
			end
		end
		return 0
	`)

	result, err := luaScript.Run(ctx, r.client, []string{queueKey}, itemID).Int()
	if err != nil {
		return fmt.Errorf("failed to execute remove script: %w", err)
	}

	if result == 0 {
		return fmt.Errorf("item not found in queue")
	}

	// Cleanup tracking data using pipeline
	pipe := r.client.Pipeline()

	positionKey := fmt.Sprintf("booking-queue-positions:%s", eventID)
	pipe.ZRem(ctx, positionKey, itemID)

	timeoutKey := fmt.Sprintf("booking-timeouts:%s", eventID)
	pipe.ZRem(ctx, timeoutKey, itemID)

	pipe.HDel(ctx, itemEventMapKey, itemID)

	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.Warn("Failed to cleanup item tracking after remove", zap.String("item_id", itemID), zap.Error(err))
	}

	// Cleanup empty queue
	go r.cleanupEmptyQueue(context.Background(), eventID)

	r.logger.Debug("Item removed from queue", zap.String("item_id", itemID), zap.String("event_id", eventID))
	return nil
}

// Close closes the queue connection
func (r *RedisQueueManager) Close() error {
	return r.client.Close()
}

// GetClient returns the Redis client (for internal use)
func (r *RedisQueueManager) GetClient() *redis.Client {
	return r.client
}

// GetActiveQueues returns all active event queue IDs
func (r *RedisQueueManager) GetActiveQueues(ctx context.Context) ([]string, error) {
	return r.client.SMembers(ctx, activeQueuesKey).Result()
}

// GetItemEventID returns the event ID for a given item ID
func (r *RedisQueueManager) GetItemEventID(ctx context.Context, itemID string) (string, error) {
	eventID, err := r.client.HGet(ctx, itemEventMapKey, itemID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("item not found")
		}
		return "", err
	}
	return eventID, nil
}

// GetQueueStats returns statistics for all active queues
func (r *RedisQueueManager) GetQueueStats(ctx context.Context) (map[string]int64, error) {
	eventIDs, err := r.client.SMembers(ctx, activeQueuesKey).Result()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64, len(eventIDs))
	for _, eventID := range eventIDs {
		queueKey := fmt.Sprintf("booking-queue:%s", eventID)
		length, err := r.client.LLen(ctx, queueKey).Result()
		if err != nil {
			continue
		}
		stats[eventID] = length
	}

	return stats, nil
}

// GetItem retrieves a queue item by ID without removing it (for authorization checks)
func (r *RedisQueueManager) GetItem(ctx context.Context, itemID string) (*QueueItem, error) {
	// Get event ID from mapping
	eventID, err := r.client.HGet(ctx, itemEventMapKey, itemID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item event mapping: %w", err)
	}

	// Search for item in the queue
	queueKey := fmt.Sprintf("booking-queue:%s", eventID)
	items, err := r.client.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search queue: %w", err)
	}

	for _, itemData := range items {
		var item QueueItem
		if err := json.Unmarshal([]byte(itemData), &item); err != nil {
			continue
		}
		if item.ID == itemID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("item not found in queue")
}

// RemoveWithAuthorization removes an item only if the requesting user owns it
func (r *RedisQueueManager) RemoveWithAuthorization(ctx context.Context, itemID, userID string) error {
	// First, get the item to verify ownership
	item, err := r.GetItem(ctx, itemID)
	if err != nil {
		return err
	}

	// Verify user ownership
	if item.UserID != userID {
		r.logger.Warn("Unauthorized cancel attempt",
			zap.String("item_id", itemID),
			zap.String("requesting_user", userID),
			zap.String("item_owner", item.UserID),
		)
		return fmt.Errorf("unauthorized: user does not own this queue item")
	}

	// User is authorized, proceed with removal
	return r.RemoveWithEventID(ctx, itemID, item.EventID)
}



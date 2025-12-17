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

// TimeoutHandler handles cleanup of expired queue items
type TimeoutHandler struct {
	client *redis.Client
	config *config.Config
	logger *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTimeoutHandler creates a new timeout handler
func NewTimeoutHandler(cfg *config.Config, client *redis.Client, logger *zap.Logger) (*TimeoutHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &TimeoutHandler{
		client: client,
		config: cfg,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start starts the timeout handler cleanup loop
func (t *TimeoutHandler) Start() {
	t.logger.Info("Starting timeout handler",
		zap.Duration("cleanup_interval", t.config.Queue.CleanupInterval),
	)

	ticker := time.NewTicker(t.config.Queue.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.ctx.Done():
			t.logger.Info("Timeout handler stopping")
			return
		case <-ticker.C:
			if err := t.cleanupExpiredItems(); err != nil {
				t.logger.Error("Error cleaning up expired items", zap.Error(err))
			}
		}
	}
}

// Stop stops the timeout handler
func (t *TimeoutHandler) Stop() {
	t.cancel()
}

// cleanupExpiredItems removes expired items from queues
func (t *TimeoutHandler) cleanupExpiredItems() error {
	now := time.Now().Unix()

	// Find all timeout tracking keys
	pattern := "booking-timeouts:*"
	keys, err := t.client.Keys(t.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get timeout keys: %w", err)
	}

	totalCleaned := 0

	for _, timeoutKey := range keys {
		// Get expired items (score < now)
		expired, err := t.client.ZRangeByScoreWithScores(t.ctx, timeoutKey, &redis.ZRangeBy{
			Min:   "0",
			Max:   fmt.Sprintf("%d", now),
			Count: 100, // Process in batches
		}).Result()

		if err != nil {
			t.logger.Warn("Failed to get expired items",
				zap.String("key", timeoutKey),
				zap.Error(err),
			)
			continue
		}

		if len(expired) == 0 {
			continue
		}

		// Extract event ID from key (booking-timeouts:{eventId})
		eventID := extractEventIDFromKey(timeoutKey)

		for _, item := range expired {
			itemID := item.Member.(string)

			// Remove from queue
			if err := t.removeItemFromQueue(t.ctx, eventID, itemID); err != nil {
				t.logger.Warn("Failed to remove expired item from queue",
					zap.String("item_id", itemID),
					zap.String("event_id", eventID),
					zap.Error(err),
				)
				continue
			}

			// Remove from position tracking
			positionKey := fmt.Sprintf("booking-queue-positions:%s", eventID)
			if err := t.client.ZRem(t.ctx, positionKey, itemID).Err(); err != nil {
				t.logger.Warn("Failed to remove from position tracking",
					zap.String("item_id", itemID),
					zap.Error(err),
				)
			}

			// Remove from timeout tracking
			if err := t.client.ZRem(t.ctx, timeoutKey, itemID).Err(); err != nil {
				t.logger.Warn("Failed to remove from timeout tracking",
					zap.String("item_id", itemID),
					zap.Error(err),
				)
			}

			totalCleaned++
			t.logger.Debug("Cleaned up expired item",
				zap.String("item_id", itemID),
				zap.String("event_id", eventID),
			)
		}
	}

	if totalCleaned > 0 {
		t.logger.Info("Cleaned up expired items",
			zap.Int("count", totalCleaned),
		)
	}

	return nil
}

// removeItemFromQueue removes an item from a specific queue
func (t *TimeoutHandler) removeItemFromQueue(ctx context.Context, eventID, itemID string) error {
	queueKey := fmt.Sprintf("booking-queue:%s", eventID)

	// Get all items
	items, err := t.client.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	// Find and remove the item
	for i, itemData := range items {
		// Try to parse as JSON to check ID
		var item QueueItem
		if err := json.Unmarshal([]byte(itemData), &item); err != nil {
			continue
		}

		if item.ID == itemID {
			// Remove by index using LSET + LREM trick
			if err := t.client.LSet(ctx, queueKey, int64(i), "___REMOVED___").Err(); err != nil {
				return err
			}
			if err := t.client.LRem(ctx, queueKey, 1, "___REMOVED___").Err(); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("item not found in queue")
}

// extractEventIDFromKey extracts event ID from a Redis key
// e.g., "booking-timeouts:event-123" -> "event-123"
func extractEventIDFromKey(key string) string {
	prefix := "booking-timeouts:"
	if len(key) > len(prefix) {
		return key[len(prefix):]
	}
	return ""
}


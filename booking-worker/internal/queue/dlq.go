package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// DLQ key patterns
	dlqKeyPrefix = "booking:dlq:"
	dlqListKey   = "booking:dlq:items"
	dlqStatsKey  = "booking:dlq:stats"
)

// DLQReason represents the reason for DLQ placement
type DLQReason string

const (
	DLQReasonMaxRetries     DLQReason = "max_retries_exceeded"
	DLQReasonCircuitOpen    DLQReason = "circuit_breaker_open"
	DLQReasonInvalidData    DLQReason = "invalid_data"
	DLQReasonServiceError   DLQReason = "service_error"
	DLQReasonTimeout        DLQReason = "timeout"
	DLQReasonUnknown        DLQReason = "unknown"
)

// DLQItem represents an item in the dead letter queue
type DLQItem struct {
	ID            string            `json:"id"`
	OriginalItem  *QueueItem        `json:"original_item"`
	Reason        DLQReason         `json:"reason"`
	ErrorMessage  string            `json:"error_message"`
	RetryCount    int               `json:"retry_count"`
	FailedAt      time.Time         `json:"failed_at"`
	LastAttemptAt time.Time         `json:"last_attempt_at"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// DLQStats holds DLQ statistics
type DLQStats struct {
	TotalItems      int64            `json:"total_items"`
	ItemsByReason   map[string]int64 `json:"items_by_reason"`
	OldestItemAge   time.Duration    `json:"oldest_item_age"`
	ProcessedToday  int64            `json:"processed_today"`
}

// DLQManager manages the dead letter queue
type DLQManager struct {
	client *redis.Client
	logger *zap.Logger
}

// NewDLQManager creates a new DLQ manager
func NewDLQManager(client *redis.Client, logger *zap.Logger) *DLQManager {
	return &DLQManager{
		client: client,
		logger: logger,
	}
}

// AddToDLQ adds a failed item to the dead letter queue
func (d *DLQManager) AddToDLQ(ctx context.Context, item *QueueItem, reason DLQReason, errMsg string, retryCount int) error {
	dlqItem := &DLQItem{
		ID:            fmt.Sprintf("dlq:%s:%d", item.ID, time.Now().UnixNano()),
		OriginalItem:  item,
		Reason:        reason,
		ErrorMessage:  errMsg,
		RetryCount:    retryCount,
		FailedAt:      time.Now(),
		LastAttemptAt: time.Now(),
		Metadata: map[string]string{
			"event_id": item.EventID,
			"user_id":  item.UserID,
		},
	}

	data, err := json.Marshal(dlqItem)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ item: %w", err)
	}

	// Use pipeline for atomic operation
	pipe := d.client.Pipeline()

	// Add to DLQ list (sorted by time)
	pipe.ZAdd(ctx, dlqListKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: dlqItem.ID,
	})

	// Store item data
	pipe.Set(ctx, dlqKeyPrefix+dlqItem.ID, data, 30*24*time.Hour) // 30 days retention

	// Update stats
	pipe.HIncrBy(ctx, dlqStatsKey, "total_items", 1)
	pipe.HIncrBy(ctx, dlqStatsKey, string(reason), 1)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add item to DLQ: %w", err)
	}

	d.logger.Info("Added item to DLQ",
		zap.String("dlq_item_id", dlqItem.ID),
		zap.String("original_item_id", item.ID),
		zap.String("reason", string(reason)),
		zap.String("error", errMsg),
		zap.Int("retry_count", retryCount),
	)

	return nil
}

// GetDLQItem retrieves a DLQ item by ID
func (d *DLQManager) GetDLQItem(ctx context.Context, dlqItemID string) (*DLQItem, error) {
	data, err := d.client.Get(ctx, dlqKeyPrefix+dlqItemID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("DLQ item not found: %s", dlqItemID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get DLQ item: %w", err)
	}

	var item DLQItem
	if err := json.Unmarshal([]byte(data), &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DLQ item: %w", err)
	}

	return &item, nil
}

// GetDLQItems retrieves items from the DLQ with pagination
func (d *DLQManager) GetDLQItems(ctx context.Context, offset, limit int64) ([]*DLQItem, error) {
	// Get IDs from sorted set
	ids, err := d.client.ZRange(ctx, dlqListKey, offset, offset+limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get DLQ item IDs: %w", err)
	}

	items := make([]*DLQItem, 0, len(ids))
	for _, id := range ids {
		item, err := d.GetDLQItem(ctx, id)
		if err != nil {
			d.logger.Warn("Failed to get DLQ item", zap.String("id", id), zap.Error(err))
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// GetDLQItemsByReason retrieves DLQ items filtered by reason
func (d *DLQManager) GetDLQItemsByReason(ctx context.Context, reason DLQReason, limit int64) ([]*DLQItem, error) {
	allItems, err := d.GetDLQItems(ctx, 0, 1000) // Get more and filter
	if err != nil {
		return nil, err
	}

	filtered := make([]*DLQItem, 0)
	for _, item := range allItems {
		if item.Reason == reason {
			filtered = append(filtered, item)
			if int64(len(filtered)) >= limit {
				break
			}
		}
	}

	return filtered, nil
}

// RemoveFromDLQ removes an item from the DLQ (after successful reprocessing)
func (d *DLQManager) RemoveFromDLQ(ctx context.Context, dlqItemID string) error {
	pipe := d.client.Pipeline()

	// Remove from sorted set
	pipe.ZRem(ctx, dlqListKey, dlqItemID)

	// Delete item data
	pipe.Del(ctx, dlqKeyPrefix+dlqItemID)

	// Update stats
	pipe.HIncrBy(ctx, dlqStatsKey, "total_items", -1)
	pipe.HIncrBy(ctx, dlqStatsKey, "processed_today", 1)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove item from DLQ: %w", err)
	}

	d.logger.Info("Removed item from DLQ", zap.String("dlq_item_id", dlqItemID))
	return nil
}

// RequeueFromDLQ moves an item from DLQ back to the main queue for reprocessing
func (d *DLQManager) RequeueFromDLQ(ctx context.Context, dlqItemID string, mainQueue QueueManager) error {
	// Get DLQ item
	dlqItem, err := d.GetDLQItem(ctx, dlqItemID)
	if err != nil {
		return err
	}

	// Enqueue to main queue
	if err := mainQueue.Enqueue(ctx, dlqItem.OriginalItem); err != nil {
		return fmt.Errorf("failed to requeue item: %w", err)
	}

	// Remove from DLQ
	if err := d.RemoveFromDLQ(ctx, dlqItemID); err != nil {
		d.logger.Warn("Item requeued but failed to remove from DLQ",
			zap.String("dlq_item_id", dlqItemID),
			zap.Error(err),
		)
	}

	d.logger.Info("Requeued item from DLQ",
		zap.String("dlq_item_id", dlqItemID),
		zap.String("original_item_id", dlqItem.OriginalItem.ID),
	)

	return nil
}

// GetDLQStats returns DLQ statistics
func (d *DLQManager) GetDLQStats(ctx context.Context) (*DLQStats, error) {
	// Get total count
	total, err := d.client.ZCard(ctx, dlqListKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get DLQ count: %w", err)
	}

	// Get stats by reason
	statsData, err := d.client.HGetAll(ctx, dlqStatsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get DLQ stats: %w", err)
	}

	itemsByReason := make(map[string]int64)
	var processedToday int64
	for key, value := range statsData {
		var count int64
		fmt.Sscanf(value, "%d", &count)
		if key == "processed_today" {
			processedToday = count
		} else if key != "total_items" {
			itemsByReason[key] = count
		}
	}

	// Get oldest item age
	var oldestAge time.Duration
	oldest, err := d.client.ZRange(ctx, dlqListKey, 0, 0).Result()
	if err == nil && len(oldest) > 0 {
		score, _ := d.client.ZScore(ctx, dlqListKey, oldest[0]).Result()
		oldestTime := time.Unix(int64(score), 0)
		oldestAge = time.Since(oldestTime)
	}

	return &DLQStats{
		TotalItems:     total,
		ItemsByReason:  itemsByReason,
		OldestItemAge:  oldestAge,
		ProcessedToday: processedToday,
	}, nil
}

// GetDLQDepth returns the current depth of the DLQ
func (d *DLQManager) GetDLQDepth(ctx context.Context) (int64, error) {
	return d.client.ZCard(ctx, dlqListKey).Result()
}

// CleanupOldItems removes items older than retention period
func (d *DLQManager) CleanupOldItems(ctx context.Context, retentionDays int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays).Unix()

	// Get old item IDs
	oldIDs, err := d.client.ZRangeByScore(ctx, dlqListKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%d", cutoff),
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get old DLQ items: %w", err)
	}

	if len(oldIDs) == 0 {
		return 0, nil
	}

	// Delete old items
	pipe := d.client.Pipeline()
	for _, id := range oldIDs {
		pipe.ZRem(ctx, dlqListKey, id)
		pipe.Del(ctx, dlqKeyPrefix+id)
	}
	pipe.HIncrBy(ctx, dlqStatsKey, "total_items", -int64(len(oldIDs)))

	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old DLQ items: %w", err)
	}

	d.logger.Info("Cleaned up old DLQ items",
		zap.Int("count", len(oldIDs)),
		zap.Int("retention_days", retentionDays),
	)

	return int64(len(oldIDs)), nil
}

package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// QueueManager manages user queues for high-demand events
type QueueManager struct {
	redis  *redis.Client
	logger *zap.Logger
}

// QueueConfig holds queue configuration
type QueueConfig struct {
	EventID                string
	MaxConcurrentUsers     int
	QueueTimeout           time.Duration
	PositionUpdateInterval time.Duration
}

// QueuePosition represents user position in queue
type QueuePosition struct {
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Position  int       `json:"position"`
	JoinedAt  time.Time `json:"joined_at"`
	ExpiresAt time.Time `json:"expires_at"`
	SessionID string    `json:"session_id"`
}

// QueueStatus represents queue status
type QueueStatus struct {
	EventID           string `json:"event_id"`
	TotalUsers        int    `json:"total_users"`
	ActiveUsers       int    `json:"active_users"`
	EstimatedWaitTime int    `json:"estimated_wait_time"` // seconds
	IsOpen            bool   `json:"is_open"`
}

// NewQueueManager creates a new queue manager
func NewQueueManager(redis *redis.Client, logger *zap.Logger) *QueueManager {
	return &QueueManager{
		redis:  redis,
		logger: logger,
	}
}

// JoinQueue adds user to event queue
func (qm *QueueManager) JoinQueue(ctx context.Context, userID, eventID, sessionID string) (*QueuePosition, error) {
	queueKey := fmt.Sprintf("queue:%s", eventID)
	userKey := fmt.Sprintf("queue_user:%s:%s", eventID, userID)

	// Check if user already in queue
	exists, err := qm.redis.Exists(ctx, userKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check user queue status: %w", err)
	}

	if exists > 0 {
		// User already in queue, return current position
		return qm.GetUserPosition(ctx, userID, eventID)
	}

	// Add user to queue
	position := &QueuePosition{
		UserID:    userID,
		EventID:   eventID,
		JoinedAt:  time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute), // 30 min queue timeout
		SessionID: sessionID,
	}

	// Add to Redis sorted set (position based on join time)
	score := float64(time.Now().UnixNano())
	err = qm.redis.ZAdd(ctx, queueKey, redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to add user to queue: %w", err)
	}

	// Store user details
	err = qm.redis.HSet(ctx, userKey, map[string]interface{}{
		"user_id":    userID,
		"event_id":   eventID,
		"joined_at":  position.JoinedAt.Format(time.RFC3339),
		"expires_at": position.ExpiresAt.Format(time.RFC3339),
		"session_id": sessionID,
	}).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store user details: %w", err)
	}

	// Set expiration
	qm.redis.Expire(ctx, userKey, 30*time.Minute)

	// Get position
	rank, err := qm.redis.ZRank(ctx, queueKey, userID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	position.Position = int(rank) + 1

	qm.logger.Info("User joined queue",
		zap.String("user_id", userID),
		zap.String("event_id", eventID),
		zap.Int("position", position.Position),
	)

	return position, nil
}

// GetUserPosition gets user position in queue
func (qm *QueueManager) GetUserPosition(ctx context.Context, userID, eventID string) (*QueuePosition, error) {
	queueKey := fmt.Sprintf("queue:%s", eventID)
	userKey := fmt.Sprintf("queue_user:%s:%s", eventID, userID)

	// Get position
	rank, err := qm.redis.ZRank(ctx, queueKey, userID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not in queue")
		}
		return nil, fmt.Errorf("failed to get user position: %w", err)
	}

	// Get user details
	userData, err := qm.redis.HGetAll(ctx, userKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	joinedAt, _ := time.Parse(time.RFC3339, userData["joined_at"])
	expiresAt, _ := time.Parse(time.RFC3339, userData["expires_at"])

	return &QueuePosition{
		UserID:    userID,
		EventID:   eventID,
		Position:  int(rank) + 1,
		JoinedAt:  joinedAt,
		ExpiresAt: expiresAt,
		SessionID: userData["session_id"],
	}, nil
}

// GetQueueStatus gets queue status
func (qm *QueueManager) GetQueueStatus(ctx context.Context, eventID string) (*QueueStatus, error) {
	queueKey := fmt.Sprintf("queue:%s", eventID)

	// Get total users in queue
	totalUsers, err := qm.redis.ZCard(ctx, queueKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue size: %w", err)
	}

	// Get active users (users currently booking)
	activeUsers, err := qm.redis.SCard(ctx, fmt.Sprintf("active_users:%s", eventID)).Result()
	if err != nil {
		activeUsers = 0 // If no active users set exists
	}

	// Calculate estimated wait time (rough estimate)
	estimatedWaitTime := int(totalUsers * 2) // 2 seconds per user

	return &QueueStatus{
		EventID:           eventID,
		TotalUsers:        int(totalUsers),
		ActiveUsers:       int(activeUsers),
		EstimatedWaitTime: estimatedWaitTime,
		IsOpen:            true, // Queue is always open
	}, nil
}

// ProcessNextUsers processes next batch of users from queue
func (qm *QueueManager) ProcessNextUsers(ctx context.Context, eventID string, batchSize int) ([]string, error) {
	queueKey := fmt.Sprintf("queue:%s", eventID)
	activeUsersKey := fmt.Sprintf("active_users:%s", eventID)

	// Get next batch of users
	users, err := qm.redis.ZRange(ctx, queueKey, 0, int64(batchSize-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get next users: %w", err)
	}

	if len(users) == 0 {
		return []string{}, nil
	}

	// Move users to active set
	for _, userID := range users {
		err = qm.redis.SAdd(ctx, activeUsersKey, userID).Err()
		if err != nil {
			qm.logger.Error("Failed to add user to active set", zap.Error(err))
			continue
		}

		// Remove from queue
		err = qm.redis.ZRem(ctx, queueKey, userID).Err()
		if err != nil {
			qm.logger.Error("Failed to remove user from queue", zap.Error(err))
		}

		// Remove user details
		userKey := fmt.Sprintf("queue_user:%s:%s", eventID, userID)
		qm.redis.Del(ctx, userKey)
	}

	qm.logger.Info("Processed users from queue",
		zap.String("event_id", eventID),
		zap.Int("batch_size", len(users)),
	)

	return users, nil
}

// CompleteUserBooking marks user booking as completed
func (qm *QueueManager) CompleteUserBooking(ctx context.Context, userID, eventID string) error {
	activeUsersKey := fmt.Sprintf("active_users:%s", eventID)

	err := qm.redis.SRem(ctx, activeUsersKey, userID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove user from active set: %w", err)
	}

	qm.logger.Info("User booking completed",
		zap.String("user_id", userID),
		zap.String("event_id", eventID),
	)

	return nil
}

// CleanupExpiredUsers removes expired users from queue
func (qm *QueueManager) CleanupExpiredUsers(ctx context.Context, eventID string) error {
	queueKey := fmt.Sprintf("queue:%s", eventID)
	now := time.Now()

	// Get all users in queue
	users, err := qm.redis.ZRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get queue users: %w", err)
	}

	expiredUsers := []string{}
	for _, userID := range users {
		userKey := fmt.Sprintf("queue_user:%s:%s", eventID, userID)
		expiresAtStr, err := qm.redis.HGet(ctx, userKey, "expires_at").Result()
		if err != nil {
			continue
		}

		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			continue
		}

		if now.After(expiresAt) {
			expiredUsers = append(expiredUsers, userID)
		}
	}

	// Remove expired users
	for _, userID := range expiredUsers {
		qm.redis.ZRem(ctx, queueKey, userID)
		userKey := fmt.Sprintf("queue_user:%s:%s", eventID, userID)
		qm.redis.Del(ctx, userKey)
	}

	if len(expiredUsers) > 0 {
		qm.logger.Info("Cleaned up expired users",
			zap.String("event_id", eventID),
			zap.Int("count", len(expiredUsers)),
		)
	}

	return nil
}

package queue

import (
	"context"
	"time"
)

// QueueItem represents an item in the queue
type QueueItem struct {
	ID        string
	EventID   string
	UserID    string
	SeatNumbers []string
	SeatCount int
	TotalAmount float64
	Currency  string
	Metadata  map[string]string
	EnqueuedAt time.Time
	ExpiresAt  time.Time
}

// QueueManager defines the interface for queue operations
type QueueManager interface {
	// Enqueue adds an item to the queue
	Enqueue(ctx context.Context, item *QueueItem) error

	// Dequeue removes and returns an item from the queue (blocking)
	Dequeue(ctx context.Context, timeout time.Duration) (*QueueItem, error)

	// GetPosition returns the position of an item in the queue
	GetPosition(ctx context.Context, itemID string) (int, error)

	// GetQueueLength returns the current queue length
	GetQueueLength(ctx context.Context, eventID string) (int, error)

	// Remove removes an item from the queue
	Remove(ctx context.Context, itemID string) error

	// Close closes the queue connection
	Close() error
}


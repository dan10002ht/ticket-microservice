package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"booking-worker/config"
	"booking-worker/grpcclient"
	"booking-worker/internal/circuitbreaker"
	"booking-worker/internal/queue"
	"booking-worker/metrics"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// ErrCircuitOpen is returned when the circuit breaker is open
var ErrCircuitOpen = errors.New("circuit breaker is open, booking service unavailable")

// Processor processes queue items
type Processor struct {
	config              *config.Config
	queue               queue.QueueManager
	dlq                 *queue.DLQManager
	logger              *zap.Logger
	ctx                 context.Context
	cancel              context.CancelFunc
	bookingClient       *grpcclient.BookingServiceClient
	realtimeClient      *grpcclient.RealtimeServiceClient
	metrics             *metrics.Exporter
	workerWg            sync.WaitGroup
	workers             []*worker
	bookingBreaker      *circuitbreaker.CircuitBreaker
}

type worker struct {
	id       int
	processor *Processor
	ctx      context.Context
}

// NewProcessor creates a new processor
func NewProcessor(cfg *config.Config, q queue.QueueManager, dlq *queue.DLQManager, logger *zap.Logger) (*Processor, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize gRPC clients
	bookingClient, err := grpcclient.NewBookingServiceClient(cfg, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create booking service client: %w", err)
	}

	realtimeClient, err := grpcclient.NewRealtimeServiceClient(cfg, logger)
	if err != nil {
		bookingClient.Close()
		cancel()
		return nil, fmt.Errorf("failed to create realtime service client: %w", err)
	}

	// Initialize circuit breaker for booking service
	breakerCfg := circuitbreaker.DefaultConfig("booking-service", logger)
	bookingBreaker := circuitbreaker.New(breakerCfg, logger)

	return &Processor{
		config:         cfg,
		queue:          q,
		dlq:            dlq,
		logger:         logger,
		ctx:            ctx,
		cancel:         cancel,
		bookingClient:  bookingClient,
		realtimeClient: realtimeClient,
		workers:        make([]*worker, 0, cfg.Worker.PoolSize),
		bookingBreaker: bookingBreaker,
	}, nil
}

// Start starts the processor with worker pool
func (p *Processor) Start() error {
	p.logger.Info("Starting worker processor",
		zap.Int("pool_size", p.config.Worker.PoolSize),
	)

	// Start worker pool
	for i := 0; i < p.config.Worker.PoolSize; i++ {
		w := &worker{
			id:        i + 1,
			processor: p,
			ctx:       p.ctx,
		}
		p.workers = append(p.workers, w)

		p.workerWg.Add(1)
		go w.run()
	}

	p.logger.Info("Worker processor started", zap.Int("workers", len(p.workers)))
	return nil
}

// Stop stops the processor
func (p *Processor) Stop() error {
	p.logger.Info("Stopping worker processor")
	p.cancel()

	// Wait for all workers to finish
	p.workerWg.Wait()

	// Close gRPC clients
	if p.bookingClient != nil {
		if err := p.bookingClient.Close(); err != nil {
			p.logger.Error("Error closing booking client", zap.Error(err))
		}
	}

	if p.realtimeClient != nil {
		if err := p.realtimeClient.Close(); err != nil {
			p.logger.Error("Error closing realtime client", zap.Error(err))
		}
	}

	p.logger.Info("Worker processor stopped")
	return nil
}

// run is the main loop for a worker
func (w *worker) run() {
	defer w.processor.workerWg.Done()

	w.processor.logger.Info("Worker started", zap.Int("worker_id", w.id))

	for {
		select {
		case <-w.ctx.Done():
			w.processor.logger.Info("Worker stopping", zap.Int("worker_id", w.id))
			return
		default:
			// Dequeue with timeout
			item, err := w.processor.queue.Dequeue(w.ctx, 5*time.Second)
			if err != nil {
				w.processor.logger.Error("Error dequeuing item",
					zap.Int("worker_id", w.id),
					zap.Error(err),
				)
				time.Sleep(1 * time.Second)
				continue
			}

			if item == nil {
				// Timeout, no item available, continue
				continue
			}

			// Process the item
			if err := w.processor.ProcessItem(w.ctx, item); err != nil {
				w.processor.logger.Error("Error processing item",
					zap.Int("worker_id", w.id),
					zap.String("item_id", item.ID),
					zap.Error(err),
				)
				// TODO: Handle retry logic or dead letter queue
			}
		}
	}
}

// GetWorkerCount returns the number of active workers
func (p *Processor) GetWorkerCount() int {
	return len(p.workers)
}

// ProcessItem processes a single queue item
func (p *Processor) ProcessItem(ctx context.Context, item *queue.QueueItem) error {
	startTime := time.Now()
	p.logger.Info("Processing queue item",
		zap.String("item_id", item.ID),
		zap.String("event_id", item.EventID),
		zap.String("user_id", item.UserID),
	)

	// Update metrics
	if p.metrics != nil {
		defer func() {
			duration := time.Since(startTime)
			p.metrics.ProcessingDuration.Observe(duration.Seconds())
		}()
	}

	// Check circuit breaker state before processing
	if p.bookingBreaker.IsOpen() {
		p.logger.Warn("Circuit breaker is open, skipping item processing",
			zap.String("item_id", item.ID),
			zap.String("circuit_state", p.bookingBreaker.State().String()),
		)
		// Re-queue the item for later processing (optional: could use DLQ instead)
		return ErrCircuitOpen
	}

	// Build request with idempotency key (item.ID ensures no duplicate bookings on retry)
	bookingReq := &grpcclient.CreateBookingRequest{
		UserID:         item.UserID,
		EventID:        item.EventID,
		SeatNumbers:    item.SeatNumbers,
		SeatCount:      item.SeatCount,
		TotalAmount:    item.TotalAmount,
		Currency:       item.Currency,
		IdempotencyKey: item.ID, // Queue item ID as idempotency key
	}

	// Retry logic with exponential backoff and circuit breaker
	var lastErr error
	for attempt := 1; attempt <= p.config.Worker.MaxRetries; attempt++ {
		if attempt > 1 {
			// Check circuit breaker before retry
			if p.bookingBreaker.IsOpen() {
				p.logger.Warn("Circuit breaker opened during retries",
					zap.String("item_id", item.ID),
					zap.Int("attempt", attempt),
				)
				lastErr = ErrCircuitOpen
				break
			}

			// Exponential backoff: baseDelay * 2^(attempt-1) with jitter
			backoffDelay := p.config.Worker.RetryInterval * time.Duration(1<<(attempt-1))
			if backoffDelay > 30*time.Second {
				backoffDelay = 30 * time.Second // Cap at 30 seconds
			}
			p.logger.Info("Retrying processing",
				zap.String("item_id", item.ID),
				zap.Int("attempt", attempt),
				zap.Duration("backoff", backoffDelay),
			)
			time.Sleep(backoffDelay)
		}

		// Call booking-service through circuit breaker
		result, cbErr := p.bookingBreaker.Execute(func() (interface{}, error) {
			return p.bookingClient.CreateBooking(ctx, bookingReq)
		})

		if cbErr != nil {
			lastErr = cbErr
			// Check if circuit breaker tripped
			if errors.Is(cbErr, gobreaker.ErrOpenState) || errors.Is(cbErr, gobreaker.ErrTooManyRequests) {
				p.logger.Warn("Circuit breaker prevented request",
					zap.String("item_id", item.ID),
					zap.String("circuit_state", p.bookingBreaker.State().String()),
					zap.Error(cbErr),
				)
				break // Don't retry if circuit is open
			}
			p.logger.Warn("Failed to create booking",
				zap.String("item_id", item.ID),
				zap.String("idempotency_key", item.ID),
				zap.Int("attempt", attempt),
				zap.Error(cbErr),
			)
			continue
		}

		bookingID := result.(string)

		// Notify realtime-service of success
		if err := p.realtimeClient.NotifyBookingResult(ctx, item.UserID, bookingID, true, "Booking created successfully"); err != nil {
			p.logger.Warn("Failed to notify booking result",
				zap.String("item_id", item.ID),
				zap.Error(err),
			)
			// Don't fail the whole operation if notification fails
		}

		// Update metrics
		if p.metrics != nil {
			p.metrics.ItemsProcessed.Inc()
		}

		p.logger.Info("Successfully processed queue item",
			zap.String("item_id", item.ID),
			zap.String("booking_id", bookingID),
		)

		return nil
	}

	// All retries failed - move to DLQ
	if p.metrics != nil {
		p.metrics.Errors.Inc()
	}

	// Determine DLQ reason based on error
	dlqReason := queue.DLQReasonMaxRetries
	if errors.Is(lastErr, ErrCircuitOpen) {
		dlqReason = queue.DLQReasonCircuitOpen
	}

	// Add to DLQ for later processing/investigation
	if p.dlq != nil {
		if dlqErr := p.dlq.AddToDLQ(ctx, item, dlqReason, lastErr.Error(), p.config.Worker.MaxRetries); dlqErr != nil {
			p.logger.Error("Failed to add item to DLQ",
				zap.String("item_id", item.ID),
				zap.Error(dlqErr),
			)
		} else {
			p.logger.Info("Item moved to DLQ",
				zap.String("item_id", item.ID),
				zap.String("reason", string(dlqReason)),
			)
		}
	}

	// Notify realtime-service of failure
	if err := p.realtimeClient.NotifyBookingResult(ctx, item.UserID, "", false, fmt.Sprintf("Failed to create booking: %v", lastErr)); err != nil {
		p.logger.Warn("Failed to notify booking failure",
			zap.String("item_id", item.ID),
			zap.Error(err),
		)
	}

	return fmt.Errorf("failed to process item after %d attempts: %w", p.config.Worker.MaxRetries, lastErr)
}

// GetDLQManager returns the DLQ manager
func (p *Processor) GetDLQManager() *queue.DLQManager {
	return p.dlq
}



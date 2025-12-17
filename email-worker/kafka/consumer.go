package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// BookingEvent represents a booking event from Kafka
type BookingEvent struct {
	BookingID        string    `json:"bookingId"`
	BookingReference string    `json:"bookingReference"`
	UserID           string    `json:"userId"`
	EventID          string    `json:"eventId"`
	Status           string    `json:"status"`
	PaymentStatus    string    `json:"paymentStatus"`
	TotalAmount      float64   `json:"totalAmount"`
	Currency         string    `json:"currency"`
	SeatCount        int       `json:"seatCount"`
	SeatNumbers      []string  `json:"seatNumbers"`
	PaymentReference string    `json:"paymentReference,omitempty"`
	ConfirmedAt      string    `json:"confirmedAt,omitempty"`
	CancelledAt      string    `json:"cancelledAt,omitempty"`
	CancellationReason string  `json:"cancellationReason,omitempty"`
	FailureReason    string    `json:"failureReason,omitempty"`
	Timestamp        string    `json:"timestamp"`
}

// BookingEventHandler defines the interface for handling booking events
type BookingEventHandler interface {
	HandleBookingConfirmed(ctx context.Context, event *BookingEvent) error
	HandleBookingCancelled(ctx context.Context, event *BookingEvent) error
	HandleBookingFailed(ctx context.Context, event *BookingEvent) error
}

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	Brokers         []string
	GroupID         string
	Topic           string
	AutoOffsetReset string
}

// BookingEventConsumer consumes booking events from Kafka
type BookingEventConsumer struct {
	config   *ConsumerConfig
	handler  BookingEventHandler
	logger   *zap.Logger
	client   sarama.ConsumerGroup
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	ready    chan bool
}

// NewBookingEventConsumer creates a new booking event consumer
func NewBookingEventConsumer(cfg *ConsumerConfig, handler BookingEventHandler, logger *zap.Logger) (*BookingEventConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	if cfg.AutoOffsetReset == "latest" {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BookingEventConsumer{
		config:  cfg,
		handler: handler,
		logger:  logger,
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
		ready:   make(chan bool),
	}, nil
}

// Start starts consuming booking events
func (c *BookingEventConsumer) Start() error {
	c.logger.Info("Starting booking event consumer",
		zap.Strings("brokers", c.config.Brokers),
		zap.String("topic", c.config.Topic),
		zap.String("group_id", c.config.GroupID),
	)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			// Consume should be called inside an infinite loop
			if err := c.client.Consume(c.ctx, []string{c.config.Topic}, c); err != nil {
				c.logger.Error("Error from consumer", zap.Error(err))
				if c.ctx.Err() != nil {
					return
				}
				time.Sleep(1 * time.Second) // Backoff before retry
			}
			// Check if context was cancelled
			if c.ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	// Wait for consumer to be ready
	<-c.ready
	c.logger.Info("Booking event consumer started and ready")

	return nil
}

// Stop stops the consumer
func (c *BookingEventConsumer) Stop() error {
	c.logger.Info("Stopping booking event consumer")
	c.cancel()
	c.wg.Wait()
	return c.client.Close()
}

// Setup is run at the beginning of a new session
func (c *BookingEventConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session
func (c *BookingEventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim handles messages from a partition
func (c *BookingEventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				c.logger.Info("Message channel closed")
				return nil
			}
			c.processMessage(session, message)
		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a single Kafka message
func (c *BookingEventConsumer) processMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	startTime := time.Now()
	ctx := context.Background()

	c.logger.Debug("Received booking event",
		zap.String("topic", message.Topic),
		zap.Int32("partition", message.Partition),
		zap.Int64("offset", message.Offset),
		zap.String("key", string(message.Key)),
	)

	// Parse event
	var event BookingEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		c.logger.Error("Failed to parse booking event",
			zap.Error(err),
			zap.String("value", string(message.Value)),
		)
		// Mark as processed to avoid infinite retry on malformed messages
		session.MarkMessage(message, "")
		return
	}

	// Determine event type from status
	var err error
	switch event.Status {
	case "CONFIRMED":
		err = c.handler.HandleBookingConfirmed(ctx, &event)
	case "CANCELLED":
		err = c.handler.HandleBookingCancelled(ctx, &event)
	case "FAILED":
		err = c.handler.HandleBookingFailed(ctx, &event)
	default:
		c.logger.Debug("Ignoring booking event with status",
			zap.String("status", event.Status),
			zap.String("booking_id", event.BookingID),
		)
	}

	if err != nil {
		c.logger.Error("Failed to handle booking event",
			zap.Error(err),
			zap.String("status", event.Status),
			zap.String("booking_id", event.BookingID),
		)
		// Still mark as processed - failures should be handled via DLQ or retry
	}

	// Mark message as processed
	session.MarkMessage(message, "")

	c.logger.Debug("Processed booking event",
		zap.String("booking_id", event.BookingID),
		zap.String("status", event.Status),
		zap.Duration("duration", time.Since(startTime)),
	)
}

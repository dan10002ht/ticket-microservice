package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"realtime-service/pkg/logger"
)

// Publisher handles Redis Pub/Sub publishing
type Publisher struct {
	client *redis.Client
}

// NewPublisher creates a new Redis Pub/Sub publisher
func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{
		client: client,
	}
}

// PublishToRoom publishes a message to a specific room channel
func (p *Publisher) PublishToRoom(ctx context.Context, room string, msgType string, payload interface{}) error {
	return p.publish(ctx, ChannelBroadcast, PubSubMessage{
		Type:    msgType,
		Room:    room,
		Payload: p.marshalPayload(payload),
	})
}

// PublishToUser publishes a message to a specific user
func (p *Publisher) PublishToUser(ctx context.Context, userID string, msgType string, payload interface{}) error {
	return p.publish(ctx, ChannelBroadcast, PubSubMessage{
		Type:    msgType,
		UserID:  userID,
		Payload: p.marshalPayload(payload),
	})
}

// PublishBookingEvent publishes a booking-related event
func (p *Publisher) PublishBookingEvent(ctx context.Context, userID string, msgType string, payload interface{}) error {
	return p.publish(ctx, ChannelBookingEvents, PubSubMessage{
		Type:    msgType,
		UserID:  userID,
		Payload: p.marshalPayload(payload),
	})
}

// PublishPaymentEvent publishes a payment-related event
func (p *Publisher) PublishPaymentEvent(ctx context.Context, userID string, msgType string, payload interface{}) error {
	return p.publish(ctx, ChannelPaymentEvents, PubSubMessage{
		Type:    msgType,
		UserID:  userID,
		Payload: p.marshalPayload(payload),
	})
}

// PublishTicketEvent publishes a ticket-related event
func (p *Publisher) PublishTicketEvent(ctx context.Context, room string, msgType string, payload interface{}) error {
	return p.publish(ctx, ChannelTicketEvents, PubSubMessage{
		Type:    msgType,
		Room:    room,
		Payload: p.marshalPayload(payload),
	})
}

// Broadcast publishes a message to all connected clients
func (p *Publisher) Broadcast(ctx context.Context, msgType string, payload interface{}) error {
	return p.publish(ctx, ChannelBroadcast, PubSubMessage{
		Type:    msgType,
		Payload: p.marshalPayload(payload),
	})
}

// publish sends a message to a Redis channel
func (p *Publisher) publish(ctx context.Context, channel string, msg PubSubMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.client.Publish(ctx, channel, data).Err(); err != nil {
		logger.Error("Failed to publish message",
			zap.String("channel", channel),
			zap.String("type", msg.Type),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish to %s: %w", channel, err)
	}

	logger.Debug("Published message",
		zap.String("channel", channel),
		zap.String("type", msg.Type),
		zap.String("room", msg.Room),
		zap.String("user_id", msg.UserID),
	)

	return nil
}

// marshalPayload converts payload to json.RawMessage
func (p *Publisher) marshalPayload(payload interface{}) json.RawMessage {
	if payload == nil {
		return json.RawMessage("{}")
	}

	// If already json.RawMessage, return as-is
	if raw, ok := payload.(json.RawMessage); ok {
		return raw
	}

	// If string (assumed to be JSON), convert directly
	if str, ok := payload.(string); ok {
		return json.RawMessage(str)
	}

	// Otherwise marshal the payload
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal payload", zap.Error(err))
		return json.RawMessage("{}")
	}
	return data
}

package pubsub

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"realtime-service/internal/websocket"
	"realtime-service/pkg/logger"
)

// Subscriber handles Redis Pub/Sub subscriptions
type Subscriber struct {
	client  *redis.Client
	hub     *websocket.Hub
	pubsub  *redis.PubSub
	ctx     context.Context
	cancel  context.CancelFunc
}

// Channel names for different event types
const (
	ChannelBookingEvents  = "booking:events"
	ChannelPaymentEvents  = "payment:events"
	ChannelTicketEvents   = "ticket:events"
	ChannelBroadcast      = "broadcast:all"
)

// PubSubMessage represents a message received from Redis Pub/Sub
type PubSubMessage struct {
	Type    string          `json:"type"`
	Room    string          `json:"room,omitempty"`
	UserID  string          `json:"user_id,omitempty"`
	Payload json.RawMessage `json:"payload"`
}

// NewSubscriber creates a new Redis Pub/Sub subscriber
func NewSubscriber(client *redis.Client, hub *websocket.Hub) *Subscriber {
	ctx, cancel := context.WithCancel(context.Background())
	return &Subscriber{
		client: client,
		hub:    hub,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start subscribes to Redis channels and processes messages
func (s *Subscriber) Start() error {
	// Subscribe to all relevant channels
	s.pubsub = s.client.Subscribe(s.ctx,
		ChannelBookingEvents,
		ChannelPaymentEvents,
		ChannelTicketEvents,
		ChannelBroadcast,
	)

	// Wait for confirmation of subscription
	_, err := s.pubsub.Receive(s.ctx)
	if err != nil {
		return err
	}

	logger.Info("Redis Pub/Sub subscriber started",
		zap.Strings("channels", []string{
			ChannelBookingEvents,
			ChannelPaymentEvents,
			ChannelTicketEvents,
			ChannelBroadcast,
		}),
	)

	// Start message processing goroutine
	go s.processMessages()

	return nil
}

// processMessages handles incoming Redis Pub/Sub messages
func (s *Subscriber) processMessages() {
	ch := s.pubsub.Channel()

	for {
		select {
		case <-s.ctx.Done():
			logger.Info("Stopping Redis Pub/Sub message processor")
			return
		case msg, ok := <-ch:
			if !ok {
				logger.Warn("Redis Pub/Sub channel closed")
				return
			}
			s.handleMessage(msg)
		}
	}
}

// handleMessage processes a single Redis Pub/Sub message
func (s *Subscriber) handleMessage(msg *redis.Message) {
	var pubsubMsg PubSubMessage
	if err := json.Unmarshal([]byte(msg.Payload), &pubsubMsg); err != nil {
		logger.Error("Failed to unmarshal pubsub message",
			zap.String("channel", msg.Channel),
			zap.Error(err),
		)
		return
	}

	logger.Debug("Received pubsub message",
		zap.String("channel", msg.Channel),
		zap.String("type", pubsubMsg.Type),
		zap.String("room", pubsubMsg.Room),
		zap.String("user_id", pubsubMsg.UserID),
	)

	// Create WebSocket message and serialize to bytes
	wsMsg := &websocket.Message{
		Type:    pubsubMsg.Type,
		Payload: pubsubMsg.Payload,
	}
	msgBytes, err := json.Marshal(wsMsg)
	if err != nil {
		logger.Error("Failed to marshal websocket message", zap.Error(err))
		return
	}

	// Route message based on target
	switch {
	case pubsubMsg.UserID != "":
		// Send to specific user
		s.hub.SendToUser(pubsubMsg.UserID, msgBytes)
	case pubsubMsg.Room != "":
		// Broadcast to room
		s.hub.BroadcastToRoom(pubsubMsg.Room, msgBytes)
	default:
		// Broadcast to all
		s.hub.Broadcast(msgBytes)
	}
}

// Stop gracefully stops the subscriber
func (s *Subscriber) Stop() error {
	s.cancel()
	if s.pubsub != nil {
		return s.pubsub.Close()
	}
	return nil
}

// SubscribeToChannel adds a subscription to a specific channel pattern
func (s *Subscriber) SubscribeToChannel(pattern string) error {
	return s.pubsub.PSubscribe(s.ctx, pattern)
}

// UnsubscribeFromChannel removes a subscription from a specific channel pattern
func (s *Subscriber) UnsubscribeFromChannel(pattern string) error {
	return s.pubsub.PUnsubscribe(s.ctx, pattern)
}

package grpcclient

import (
	"context"
	"time"

	"booking-worker/config"
	realtimepb "booking-worker/internal/protos/realtime"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// RealtimeServiceClient wraps the gRPC client for Realtime Service
type RealtimeServiceClient struct {
	conn   *grpc.ClientConn
	client realtimepb.RealtimeServiceClient
	config *config.Config
	logger *zap.Logger
}

// NewRealtimeServiceClient creates a new Realtime Service gRPC client
func NewRealtimeServiceClient(cfg *config.Config, logger *zap.Logger) (*RealtimeServiceClient, error) {
	endpoint := cfg.GRPC.RealtimeServiceEndpoint

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		// Log warning but don't fail - realtime service is optional
		logger.Warn("Failed to connect to Realtime Service (service may not be available)",
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		return &RealtimeServiceClient{
			conn:   nil,
			client: nil,
			config: cfg,
			logger: logger,
		}, nil
	}

	client := realtimepb.NewRealtimeServiceClient(conn)

	logger.Info("Connected to Realtime Service",
		zap.String("endpoint", endpoint),
	)

	return &RealtimeServiceClient{
		conn:   conn,
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// NotifyQueuePosition notifies a client of their queue position
func (c *RealtimeServiceClient) NotifyQueuePosition(ctx context.Context, userID, eventID string, position int) error {
	if c.client == nil {
		c.logger.Debug("Realtime service not connected, skipping queue position notification",
			zap.String("user_id", userID),
			zap.Int("position", position),
		)
		return nil
	}

	req := &realtimepb.NotifyQueuePositionRequest{
		UserId:   userID,
		EventId:  eventID,
		Position: int32(position),
	}

	resp, err := c.client.NotifyQueuePosition(ctx, req)
	if err != nil {
		c.logger.Error("Failed to notify queue position",
			zap.String("user_id", userID),
			zap.String("event_id", eventID),
			zap.Int("position", position),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Queue position notification sent",
		zap.String("user_id", userID),
		zap.String("event_id", eventID),
		zap.Int("position", position),
		zap.Bool("delivered", resp.GetDelivered()),
	)

	return nil
}

// NotifyBookingResult notifies a client of booking result
func (c *RealtimeServiceClient) NotifyBookingResult(ctx context.Context, userID, bookingID string, success bool, message string) error {
	if c.client == nil {
		c.logger.Debug("Realtime service not connected, skipping booking result notification",
			zap.String("user_id", userID),
			zap.Bool("success", success),
		)
		return nil
	}

	req := &realtimepb.NotifyBookingResultRequest{
		UserId:    userID,
		BookingId: bookingID,
		Success:   success,
		Message:   message,
	}

	resp, err := c.client.NotifyBookingResult(ctx, req)
	if err != nil {
		c.logger.Error("Failed to notify booking result",
			zap.String("user_id", userID),
			zap.String("booking_id", bookingID),
			zap.Bool("success", success),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Booking result notification sent",
		zap.String("user_id", userID),
		zap.String("booking_id", bookingID),
		zap.Bool("success", success),
		zap.Bool("delivered", resp.GetDelivered()),
		zap.Int32("active_connections", resp.GetActiveConnections()),
	)

	return nil
}

// NotifyBookingResultWithDetails notifies a client of booking result with full details
func (c *RealtimeServiceClient) NotifyBookingResultWithDetails(
	ctx context.Context,
	userID, bookingID string,
	success bool,
	message string,
	bookingReference, eventID string,
	seatNumbers []string,
	totalAmount, currency string,
) error {
	if c.client == nil {
		c.logger.Debug("Realtime service not connected, skipping booking result notification",
			zap.String("user_id", userID),
			zap.Bool("success", success),
		)
		return nil
	}

	req := &realtimepb.NotifyBookingResultRequest{
		UserId:           userID,
		BookingId:        bookingID,
		Success:          success,
		Message:          message,
		BookingReference: bookingReference,
		EventId:          eventID,
		SeatNumbers:      seatNumbers,
		TotalAmount:      totalAmount,
		Currency:         currency,
	}

	resp, err := c.client.NotifyBookingResult(ctx, req)
	if err != nil {
		c.logger.Error("Failed to notify booking result",
			zap.String("user_id", userID),
			zap.String("booking_id", bookingID),
			zap.Bool("success", success),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Booking result notification sent",
		zap.String("user_id", userID),
		zap.String("booking_id", bookingID),
		zap.Bool("success", success),
		zap.Bool("delivered", resp.GetDelivered()),
		zap.Int32("active_connections", resp.GetActiveConnections()),
	)

	return nil
}

// NotifyPaymentStatus notifies a client of payment status update
func (c *RealtimeServiceClient) NotifyPaymentStatus(
	ctx context.Context,
	userID, bookingID, paymentID string,
	status realtimepb.PaymentStatus,
	message, amount, currency string,
) error {
	if c.client == nil {
		c.logger.Debug("Realtime service not connected, skipping payment status notification",
			zap.String("user_id", userID),
			zap.String("payment_id", paymentID),
		)
		return nil
	}

	req := &realtimepb.NotifyPaymentStatusRequest{
		UserId:    userID,
		BookingId: bookingID,
		PaymentId: paymentID,
		Status:    status,
		Message:   message,
		Amount:    amount,
		Currency:  currency,
	}

	resp, err := c.client.NotifyPaymentStatus(ctx, req)
	if err != nil {
		c.logger.Error("Failed to notify payment status",
			zap.String("user_id", userID),
			zap.String("payment_id", paymentID),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Payment status notification sent",
		zap.String("user_id", userID),
		zap.String("payment_id", paymentID),
		zap.Bool("delivered", resp.GetDelivered()),
	)

	return nil
}

// BroadcastEvent broadcasts an event to a room
func (c *RealtimeServiceClient) BroadcastEvent(ctx context.Context, eventType, room, payload string) error {
	if c.client == nil {
		c.logger.Debug("Realtime service not connected, skipping broadcast",
			zap.String("event_type", eventType),
			zap.String("room", room),
		)
		return nil
	}

	req := &realtimepb.BroadcastEventRequest{
		EventType: eventType,
		Room:      room,
		Payload:   payload,
	}

	resp, err := c.client.BroadcastEvent(ctx, req)
	if err != nil {
		c.logger.Error("Failed to broadcast event",
			zap.String("event_type", eventType),
			zap.String("room", room),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Event broadcasted",
		zap.String("event_type", eventType),
		zap.String("room", room),
		zap.Int32("recipients", resp.GetRecipients()),
	)

	return nil
}

// SendToUser sends a notification to a specific user
func (c *RealtimeServiceClient) SendToUser(ctx context.Context, userID, eventType, payload string) error {
	if c.client == nil {
		c.logger.Debug("Realtime service not connected, skipping send to user",
			zap.String("user_id", userID),
			zap.String("event_type", eventType),
		)
		return nil
	}

	req := &realtimepb.SendToUserRequest{
		UserId:    userID,
		EventType: eventType,
		Payload:   payload,
	}

	resp, err := c.client.SendToUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to send to user",
			zap.String("user_id", userID),
			zap.String("event_type", eventType),
			zap.Error(err),
		)
		return err
	}

	c.logger.Info("Message sent to user",
		zap.String("user_id", userID),
		zap.String("event_type", eventType),
		zap.Bool("delivered", resp.GetDelivered()),
		zap.Int32("connections_notified", resp.GetConnectionsNotified()),
	)

	return nil
}

// GetConnectionStats retrieves connection statistics from realtime service
func (c *RealtimeServiceClient) GetConnectionStats(ctx context.Context, room string) (*realtimepb.GetConnectionStatsResponse, error) {
	if c.client == nil {
		return nil, nil
	}

	req := &realtimepb.GetConnectionStatsRequest{
		Room: room,
	}

	return c.client.GetConnectionStats(ctx, req)
}

// Close closes the gRPC connection
func (c *RealtimeServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

package grpcclient

import (
	"context"
	"time"

	"booking-worker/config"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// RealtimeServiceClient wraps the gRPC client for Realtime Service
// NOTE: This is a stub implementation. Will be updated when realtime.proto is available.
type RealtimeServiceClient struct {
	conn   *grpc.ClientConn
	config *config.Config
	logger *zap.Logger
	// TODO: Add generated protobuf client when realtime.proto is created
	// realtimepb.RealtimeServiceClient
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
			config: cfg,
			logger: logger,
		}, nil
	}

	logger.Info("Connected to Realtime Service",
		zap.String("endpoint", endpoint),
	)

	return &RealtimeServiceClient{
		conn:   conn,
		config: cfg,
		logger: logger,
	}, nil
}

// NotifyQueuePosition notifies a client of their queue position
// TODO: Implement when realtime.proto is generated
func (c *RealtimeServiceClient) NotifyQueuePosition(ctx context.Context, userID, itemID string, position int) error {
	if c.conn == nil {
		c.logger.Debug("Realtime service not connected, skipping queue position notification",
			zap.String("user_id", userID),
			zap.Int("position", position),
		)
		return nil
	}

	// TODO: Call realtimepb.RealtimeServiceClient.NotifyQueuePosition when proto is available
	c.logger.Info("Queue position notification (stub)",
		zap.String("user_id", userID),
		zap.String("item_id", itemID),
		zap.Int("position", position),
	)

	return nil
}

// NotifyBookingResult notifies a client of booking result
// TODO: Implement when realtime.proto is generated
func (c *RealtimeServiceClient) NotifyBookingResult(ctx context.Context, userID, bookingID string, success bool, message string) error {
	if c.conn == nil {
		c.logger.Debug("Realtime service not connected, skipping booking result notification",
			zap.String("user_id", userID),
			zap.Bool("success", success),
		)
		return nil
	}

	// TODO: Call realtimepb.RealtimeServiceClient.NotifyBookingResult when proto is available
	c.logger.Info("Booking result notification (stub)",
		zap.String("user_id", userID),
		zap.String("booking_id", bookingID),
		zap.Bool("success", success),
		zap.String("message", message),
	)

	return nil
}

// Close closes the gRPC connection
func (c *RealtimeServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

package grpcclient

import (
	"context"
	"fmt"
	"ticket-service/config"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	eventpb "shared-lib/protos/event"
)

// EventServiceClient handles communication with Event Service
type EventServiceClient struct {
	conn   *grpc.ClientConn
	client eventpb.EventServiceClient
	logger *zap.Logger
}

// NewEventServiceClient creates a new Event Service gRPC client
func NewEventServiceClient(config config.EventServiceConfig, logger *zap.Logger) (*EventServiceClient, error) {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Event Service: %w", err)
	}

	client := eventpb.NewEventServiceClient(conn)

	logger.Info("Connected to Event Service",
		zap.String("address", address),
	)

	return &EventServiceClient{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// GetEvent retrieves event information
func (c *EventServiceClient) GetEvent(ctx context.Context, eventID string) (*eventpb.Event, error) {
	req := &eventpb.GetEventRequest{
		Id: eventID,
	}

	resp, err := c.client.GetEvent(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get event",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, err
	}

	return resp.Event, nil
}

// GetSeatAvailability retrieves seat availability
func (c *EventServiceClient) GetSeatAvailability(ctx context.Context, eventID, seatID string) (*eventpb.GetSeatAvailabilityResponse, error) {
	req := &eventpb.GetSeatAvailabilityRequest{
		EventId: eventID,
		SeatId:  seatID,
	}

	resp, err := c.client.GetSeatAvailability(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get seat availability",
			zap.String("event_id", eventID),
			zap.String("seat_id", seatID),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// UpdateSeatAvailability updates seat availability status
func (c *EventServiceClient) UpdateSeatAvailability(ctx context.Context, eventID, seatID, status, userID, bookingID string) error {
	req := &eventpb.UpdateSeatAvailabilityRequest{
		EventId:   eventID,
		SeatId:    seatID,
		Status:    status,
		UserId:    userID,
		BookingId: bookingID,
	}

	resp, err := c.client.UpdateSeatAvailability(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update seat availability",
			zap.String("event_id", eventID),
			zap.String("seat_id", seatID),
			zap.String("status", status),
			zap.Error(err),
		)
		return err
	}

	if !resp.Success {
		return fmt.Errorf("failed to update seat availability: %s", resp.Error)
	}

	return nil
}

// BlockSeats blocks multiple seats
func (c *EventServiceClient) BlockSeats(ctx context.Context, eventID string, seatIDs []string, userID, bookingID string, expiresAt time.Time) (*eventpb.BlockSeatsResponse, error) {
	req := &eventpb.BlockSeatsRequest{
		EventId:   eventID,
		SeatIds:   seatIDs,
		UserId:    userID,
		BookingId: bookingID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}

	resp, err := c.client.BlockSeats(ctx, req)
	if err != nil {
		c.logger.Error("Failed to block seats",
			zap.String("event_id", eventID),
			zap.Strings("seat_ids", seatIDs),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// ReleaseSeats releases blocked seats
func (c *EventServiceClient) ReleaseSeats(ctx context.Context, eventID string, seatIDs []string, userID, reason string) (*eventpb.ReleaseSeatsResponse, error) {
	req := &eventpb.ReleaseSeatsRequest{
		EventId: eventID,
		SeatIds: seatIDs,
		UserId:  userID,
		Reason:  reason,
	}

	resp, err := c.client.ReleaseSeats(ctx, req)
	if err != nil {
		c.logger.Error("Failed to release seats",
			zap.String("event_id", eventID),
			zap.Strings("seat_ids", seatIDs),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// Close closes the gRPC connection
func (c *EventServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// HealthCheck performs a health check on the Event Service
func (c *EventServiceClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.GetEvent(ctx, "health-check")
	if err != nil {
		return fmt.Errorf("Event Service health check failed: %w", err)
	}

	return nil
}

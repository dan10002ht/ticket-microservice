package grpcclient

import (
	"context"
	"fmt"
	"time"

	"booking-worker/config"
	bookingpb "booking-worker/internal/protos/booking"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// BookingServiceClient wraps the gRPC client for Booking Service
type BookingServiceClient struct {
	conn   *grpc.ClientConn
	client bookingpb.BookingServiceClient
	config *config.Config
	logger *zap.Logger
}

// NewBookingServiceClient creates a new Booking Service gRPC client
func NewBookingServiceClient(cfg *config.Config, logger *zap.Logger) (*BookingServiceClient, error) {
	endpoint := cfg.GRPC.BookingServiceEndpoint

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
		return nil, fmt.Errorf("failed to connect to booking service: %w", err)
	}

	client := bookingpb.NewBookingServiceClient(conn)

	logger.Info("Connected to Booking Service",
		zap.String("endpoint", endpoint),
	)

	return &BookingServiceClient{
		conn:   conn,
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// CreateBookingRequest contains all parameters for creating a booking
type CreateBookingRequest struct {
	UserID         string
	EventID        string
	SeatNumbers    []string
	SeatCount      int
	TotalAmount    float64
	Currency       string
	IdempotencyKey string // Queue item ID for idempotency
}

// CreateBooking creates a booking via gRPC with idempotency support
func (c *BookingServiceClient) CreateBooking(ctx context.Context, req *CreateBookingRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	grpcReq := &bookingpb.CreateBookingRequest{
		UserId:         req.UserID,
		EventId:        req.EventID,
		TicketQuantity: int32(req.SeatCount),
		SeatNumbers:    req.SeatNumbers,
		IdempotencyKey: req.IdempotencyKey, // Pass idempotency key to prevent duplicates
		Metadata: map[string]string{
			"total_amount": fmt.Sprintf("%.2f", req.TotalAmount),
			"currency":     req.Currency,
		},
	}

	resp, err := c.client.CreateBooking(ctx, grpcReq)
	if err != nil {
		c.logger.Error("Failed to create booking",
			zap.String("user_id", req.UserID),
			zap.String("event_id", req.EventID),
			zap.String("idempotency_key", req.IdempotencyKey),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to create booking: %w", err)
	}

	if !resp.Success {
		return "", fmt.Errorf("booking creation failed: %s", resp.Message)
	}

	c.logger.Info("Booking created successfully",
		zap.String("booking_id", resp.Booking.Id),
		zap.String("booking_reference", resp.Booking.BookingReference),
		zap.String("idempotency_key", req.IdempotencyKey),
	)

	return resp.Booking.Id, nil
}

// GetBooking gets booking details by ID
func (c *BookingServiceClient) GetBooking(ctx context.Context, bookingID string) (*bookingpb.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &bookingpb.GetBookingRequest{
		BookingId: bookingID,
	}

	resp, err := c.client.GetBooking(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("get booking failed: %s", resp.Message)
	}

	return resp.Booking, nil
}

// ConfirmBooking confirms a booking after payment
func (c *BookingServiceClient) ConfirmBooking(ctx context.Context, bookingID, paymentReference string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := &bookingpb.ConfirmBookingRequest{
		BookingId:        bookingID,
		PaymentReference: paymentReference,
	}

	resp, err := c.client.ConfirmBooking(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to confirm booking: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("confirm booking failed: %s", resp.Message)
	}

	c.logger.Info("Booking confirmed",
		zap.String("booking_id", bookingID),
	)

	return nil
}

// CancelBooking cancels a booking
func (c *BookingServiceClient) CancelBooking(ctx context.Context, bookingID, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := &bookingpb.CancelBookingRequest{
		BookingId: bookingID,
		Reason:    reason,
	}

	resp, err := c.client.CancelBooking(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("cancel booking failed: %s", resp.Message)
	}

	c.logger.Info("Booking cancelled",
		zap.String("booking_id", bookingID),
		zap.String("reason", reason),
	)

	return nil
}

// Close closes the gRPC connection
func (c *BookingServiceClient) Close() error {
	return c.conn.Close()
}

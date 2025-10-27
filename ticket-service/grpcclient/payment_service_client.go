package grpcclient

import (
	"context"
	"fmt"
	"ticket-service/config"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentpb "shared-lib/protos/payment"
)

// PaymentServiceClient handles communication with Payment Service
type PaymentServiceClient struct {
	conn   *grpc.ClientConn
	client paymentpb.PaymentServiceClient
	logger *zap.Logger
}

// NewPaymentServiceClient creates a new Payment Service gRPC client
func NewPaymentServiceClient(config config.PaymentServiceConfig, logger *zap.Logger) (*PaymentServiceClient, error) {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Payment Service: %w", err)
	}

	client := paymentpb.NewPaymentServiceClient(conn)

	logger.Info("Connected to Payment Service",
		zap.String("address", address),
	)

	return &PaymentServiceClient{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// ProcessPayment processes a payment
func (c *PaymentServiceClient) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	resp, err := c.client.ProcessPayment(ctx, req)
	if err != nil {
		c.logger.Error("Failed to process payment",
			zap.String("booking_id", req.BookingId),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// RefundPayment processes a refund
func (c *PaymentServiceClient) RefundPayment(ctx context.Context, req *paymentpb.RefundPaymentRequest) (*paymentpb.RefundPaymentResponse, error) {
	resp, err := c.client.RefundPayment(ctx, req)
	if err != nil {
		c.logger.Error("Failed to refund payment",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// GetPaymentStatus retrieves payment status
func (c *PaymentServiceClient) GetPaymentStatus(ctx context.Context, paymentID string) (*paymentpb.GetPaymentStatusResponse, error) {
	req := &paymentpb.GetPaymentStatusRequest{
		PaymentId: paymentID,
	}

	resp, err := c.client.GetPaymentStatus(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get payment status",
			zap.String("payment_id", paymentID),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// Close closes the gRPC connection
func (c *PaymentServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// HealthCheck performs a health check on the Payment Service
func (c *PaymentServiceClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &paymentpb.GetPaymentStatusRequest{
		PaymentId: "health-check",
	}

	_, err := c.client.GetPaymentStatus(ctx, req)
	if err != nil {
		return fmt.Errorf("Payment Service health check failed: %w", err)
	}

	return nil
}

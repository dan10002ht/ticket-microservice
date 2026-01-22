package grpcclient

import (
	"context"
	"fmt"
	"ticket-service/config"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentpb "ticket-service/internal/protos/payment"
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

// CreatePayment creates a new payment
func (c *PaymentServiceClient) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.PaymentResponse, error) {
	resp, err := c.client.CreatePayment(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create payment",
			zap.String("booking_id", req.BookingId),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// CreateRefund creates a refund for a payment
func (c *PaymentServiceClient) CreateRefund(ctx context.Context, req *paymentpb.CreateRefundRequest) (*paymentpb.RefundResponse, error) {
	resp, err := c.client.CreateRefund(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create refund",
			zap.String("payment_id", req.PaymentId),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// GetPayment retrieves payment details
func (c *PaymentServiceClient) GetPayment(ctx context.Context, paymentID string) (*paymentpb.PaymentResponse, error) {
	req := &paymentpb.GetPaymentRequest{
		PaymentId: paymentID,
	}

	resp, err := c.client.GetPayment(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get payment",
			zap.String("payment_id", paymentID),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// MarkPaymentSuccess marks a payment as successful
func (c *PaymentServiceClient) MarkPaymentSuccess(ctx context.Context, paymentID, externalRef string) (*paymentpb.PaymentResponse, error) {
	req := &paymentpb.MarkPaymentSuccessRequest{
		PaymentId:         paymentID,
		ExternalReference: externalRef,
	}

	resp, err := c.client.MarkPaymentSuccess(ctx, req)
	if err != nil {
		c.logger.Error("Failed to mark payment success",
			zap.String("payment_id", paymentID),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// MarkPaymentFailed marks a payment as failed
func (c *PaymentServiceClient) MarkPaymentFailed(ctx context.Context, paymentID, reason string) (*paymentpb.PaymentResponse, error) {
	req := &paymentpb.MarkPaymentFailedRequest{
		PaymentId:     paymentID,
		FailureReason: reason,
	}

	resp, err := c.client.MarkPaymentFailed(ctx, req)
	if err != nil {
		c.logger.Error("Failed to mark payment failed",
			zap.String("payment_id", paymentID),
			zap.Error(err),
		)
		return nil, err
	}

	return resp, nil
}

// CancelPayment cancels a payment
func (c *PaymentServiceClient) CancelPayment(ctx context.Context, paymentID, reason string) (*paymentpb.PaymentResponse, error) {
	req := &paymentpb.CancelPaymentRequest{
		PaymentId: paymentID,
		Reason:    reason,
	}

	resp, err := c.client.CancelPayment(ctx, req)
	if err != nil {
		c.logger.Error("Failed to cancel payment",
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

	req := &paymentpb.HealthRequest{
		Service: "ticket-service",
	}

	_, err := c.client.Health(ctx, req)
	if err != nil {
		return fmt.Errorf("Payment Service health check failed: %w", err)
	}

	return nil
}

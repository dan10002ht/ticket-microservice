package kafka

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// EmailJobCreator interface for creating email jobs
type EmailJobCreator interface {
	CreateEmailJob(ctx context.Context, templateName, recipientID, recipientEmail string, data map[string]interface{}) error
}

// UserFetcher interface for fetching user details
type UserFetcher interface {
	GetUserEmail(ctx context.Context, userID string) (string, error)
}

// BookingEmailHandler handles booking events and creates email jobs
type BookingEmailHandler struct {
	logger       *zap.Logger
	emailCreator EmailJobCreator
	userFetcher  UserFetcher
}

// NewBookingEmailHandler creates a new booking email handler
func NewBookingEmailHandler(logger *zap.Logger, emailCreator EmailJobCreator, userFetcher UserFetcher) *BookingEmailHandler {
	return &BookingEmailHandler{
		logger:       logger,
		emailCreator: emailCreator,
		userFetcher:  userFetcher,
	}
}

// HandleBookingConfirmed handles BOOKING_CONFIRMED events
func (h *BookingEmailHandler) HandleBookingConfirmed(ctx context.Context, event *BookingEvent) error {
	h.logger.Info("Handling booking confirmed event",
		zap.String("booking_id", event.BookingID),
		zap.String("user_id", event.UserID),
		zap.String("booking_reference", event.BookingReference),
	)

	// Get user email
	userEmail, err := h.userFetcher.GetUserEmail(ctx, event.UserID)
	if err != nil {
		h.logger.Error("Failed to get user email",
			zap.Error(err),
			zap.String("user_id", event.UserID),
		)
		return fmt.Errorf("failed to get user email: %w", err)
	}

	// Prepare email data
	data := map[string]interface{}{
		"booking_reference": event.BookingReference,
		"booking_id":        event.BookingID,
		"event_id":          event.EventID,
		"total_amount":      event.TotalAmount,
		"currency":          event.Currency,
		"seat_count":        event.SeatCount,
		"seat_numbers":      event.SeatNumbers,
		"payment_reference": event.PaymentReference,
		"confirmed_at":      event.ConfirmedAt,
	}

	// Create email job
	if err := h.emailCreator.CreateEmailJob(ctx, "booking_confirmation", event.UserID, userEmail, data); err != nil {
		h.logger.Error("Failed to create booking confirmation email job",
			zap.Error(err),
			zap.String("booking_id", event.BookingID),
		)
		return fmt.Errorf("failed to create email job: %w", err)
	}

	h.logger.Info("Created booking confirmation email job",
		zap.String("booking_id", event.BookingID),
		zap.String("user_email", userEmail),
	)

	return nil
}

// HandleBookingCancelled handles BOOKING_CANCELLED events
func (h *BookingEmailHandler) HandleBookingCancelled(ctx context.Context, event *BookingEvent) error {
	h.logger.Info("Handling booking cancelled event",
		zap.String("booking_id", event.BookingID),
		zap.String("user_id", event.UserID),
		zap.String("reason", event.CancellationReason),
	)

	// Get user email
	userEmail, err := h.userFetcher.GetUserEmail(ctx, event.UserID)
	if err != nil {
		h.logger.Error("Failed to get user email",
			zap.Error(err),
			zap.String("user_id", event.UserID),
		)
		return fmt.Errorf("failed to get user email: %w", err)
	}

	// Prepare email data
	data := map[string]interface{}{
		"booking_reference":   event.BookingReference,
		"booking_id":          event.BookingID,
		"event_id":            event.EventID,
		"cancellation_reason": event.CancellationReason,
		"cancelled_at":        event.CancelledAt,
	}

	// Create email job
	if err := h.emailCreator.CreateEmailJob(ctx, "booking_cancellation", event.UserID, userEmail, data); err != nil {
		h.logger.Error("Failed to create booking cancellation email job",
			zap.Error(err),
			zap.String("booking_id", event.BookingID),
		)
		return fmt.Errorf("failed to create email job: %w", err)
	}

	h.logger.Info("Created booking cancellation email job",
		zap.String("booking_id", event.BookingID),
		zap.String("user_email", userEmail),
	)

	return nil
}

// HandleBookingFailed handles BOOKING_FAILED events
func (h *BookingEmailHandler) HandleBookingFailed(ctx context.Context, event *BookingEvent) error {
	h.logger.Info("Handling booking failed event",
		zap.String("booking_id", event.BookingID),
		zap.String("user_id", event.UserID),
		zap.String("reason", event.FailureReason),
	)

	// Get user email
	userEmail, err := h.userFetcher.GetUserEmail(ctx, event.UserID)
	if err != nil {
		h.logger.Error("Failed to get user email",
			zap.Error(err),
			zap.String("user_id", event.UserID),
		)
		return fmt.Errorf("failed to get user email: %w", err)
	}

	// Prepare email data
	data := map[string]interface{}{
		"booking_reference": event.BookingReference,
		"booking_id":        event.BookingID,
		"event_id":          event.EventID,
		"failure_reason":    event.FailureReason,
	}

	// Create email job
	if err := h.emailCreator.CreateEmailJob(ctx, "booking_failed", event.UserID, userEmail, data); err != nil {
		h.logger.Error("Failed to create booking failed email job",
			zap.Error(err),
			zap.String("booking_id", event.BookingID),
		)
		return fmt.Errorf("failed to create email job: %w", err)
	}

	h.logger.Info("Created booking failed email job",
		zap.String("booking_id", event.BookingID),
		zap.String("user_email", userEmail),
	)

	return nil
}

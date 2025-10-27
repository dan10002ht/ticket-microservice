package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"ticket-service/grpcclient"
	"ticket-service/metrics"
	"ticket-service/models"
	"ticket-service/repositories"
)

// TicketService handles ticket business logic
type TicketService struct {
	ticketRepo    *repositories.TicketRepository
	eventClient   *grpcclient.EventServiceClient
	paymentClient *grpcclient.PaymentServiceClient
	logger        *zap.Logger
}

// NewTicketService creates a new ticket service
func NewTicketService(
	ticketRepo *repositories.TicketRepository,
	eventClient *grpcclient.EventServiceClient,
	paymentClient *grpcclient.PaymentServiceClient,
	logger *zap.Logger,
) *TicketService {
	return &TicketService{
		ticketRepo:    ticketRepo,
		eventClient:   eventClient,
		paymentClient: paymentClient,
		logger:        logger,
	}
}

// CreateTicket creates a new ticket
func (s *TicketService) CreateTicket(ctx context.Context, req *CreateTicketRequest) (*models.Ticket, error) {
	// Validate request
	if err := s.validateCreateTicketRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check seat availability via Event Service
	if s.eventClient != nil {
		available, err := s.checkSeatAvailability(ctx, req.EventID, req.SeatID)
		if err != nil {
			return nil, fmt.Errorf("failed to check seat availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("seat %s is not available", req.SeatID)
		}
	}

	// Generate ticket number
	ticketNumber := s.generateTicketNumber(req.EventID, req.SeatID)

	// Create ticket model
	ticket := models.NewTicket(
		req.EventID, req.SeatID, req.ZoneID, req.UserID,
		ticketNumber, req.PricingCategory, req.BasePrice, req.FinalPrice, req.Currency,
	)

	// Set additional fields
	if req.BookingSessionID != "" {
		ticket.BookingSessionID = &req.BookingSessionID
	}
	if req.TicketType != "" {
		ticket.TicketType = req.TicketType
	}
	if req.DiscountReason != "" {
		ticket.DiscountReason = &req.DiscountReason
	}
	if req.ValidUntil != nil {
		ticket.ValidUntil = req.ValidUntil
	}
	if req.CreatedBy != "" {
		ticket.CreatedBy = &req.CreatedBy
	}

	// Validate ticket
	if err := ticket.Validate(); err != nil {
		return nil, fmt.Errorf("ticket validation failed: %w", err)
	}

	// Create ticket in database
	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Update seat status in Event Service
	if s.eventClient != nil {
		err := s.eventClient.UpdateSeatAvailability(ctx, req.EventID, req.SeatID, "sold", req.UserID, ticket.ID)
		if err != nil {
			s.logger.Warn("Failed to update seat status in Event Service",
				zap.String("event_id", req.EventID),
				zap.String("seat_id", req.SeatID),
				zap.Error(err),
			)
		}
	}

	// Increment metrics
	metrics.IncrementTicketCreated(req.EventID, ticket.TicketType, ticket.Status)

	s.logger.Info("Ticket created successfully",
		zap.String("ticket_id", ticket.ID),
		zap.String("event_id", req.EventID),
		zap.String("user_id", req.UserID),
	)

	return ticket, nil
}

// GetTicket retrieves a ticket by ID
func (s *TicketService) GetTicket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}

// GetTicketByNumber retrieves a ticket by ticket number
func (s *TicketService) GetTicketByNumber(ctx context.Context, ticketNumber string) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetByTicketNumber(ctx, ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}

// GetUserTickets retrieves tickets for a user
func (s *TicketService) GetUserTickets(ctx context.Context, userID string, page, limit int) ([]*models.Ticket, int, error) {
	tickets, total, err := s.ticketRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user tickets: %w", err)
	}

	return tickets, total, nil
}

// GetEventTickets retrieves tickets for an event
func (s *TicketService) GetEventTickets(ctx context.Context, eventID string, page, limit int) ([]*models.Ticket, int, error) {
	tickets, total, err := s.ticketRepo.GetByEventID(ctx, eventID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get event tickets: %w", err)
	}

	return tickets, total, nil
}

// UpdateTicketStatus updates ticket status
func (s *TicketService) UpdateTicketStatus(ctx context.Context, ticketID, status, updatedBy string) error {
	// Get current ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Validate status transition
	if !s.isValidStatusTransition(ticket.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", ticket.Status, status)
	}

	// Update status
	if err := s.ticketRepo.UpdateStatus(ctx, ticketID, status, updatedBy); err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	// Handle status-specific logic
	switch status {
	case models.TicketStatusConfirmed:
		s.handleTicketConfirmation(ctx, ticket)
	case models.TicketStatusCancelled:
		s.handleTicketCancellation(ctx, ticket)
	case models.TicketStatusUsed:
		s.handleTicketUsage(ctx, ticket)
	}

	s.logger.Info("Ticket status updated",
		zap.String("ticket_id", ticketID),
		zap.String("old_status", ticket.Status),
		zap.String("new_status", status),
		zap.String("updated_by", updatedBy),
	)

	return nil
}

// ProcessPayment processes payment for a ticket
func (s *TicketService) ProcessPayment(ctx context.Context, req *ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Validate ticket can be paid
	if !ticket.CanBeCancelled() {
		return nil, fmt.Errorf("ticket cannot be paid, current status: %s", ticket.Status)
	}

	// Process payment via Payment Service
	if s.paymentClient != nil {
		paymentReq := &grpcclient.ProcessPaymentRequest{
			TicketId:      req.TicketID,
			Amount:        ticket.FinalPrice,
			Currency:      ticket.Currency,
			PaymentMethod: req.PaymentMethod,
			UserId:        ticket.UserID,
			EventId:       ticket.EventID,
		}

		paymentResp, err := s.paymentClient.ProcessPayment(ctx, paymentReq)
		if err != nil {
			return nil, fmt.Errorf("payment processing failed: %w", err)
		}

		// Update ticket payment status
		err = s.ticketRepo.UpdatePaymentStatus(ctx, req.TicketID,
			models.PaymentStatusPaid, req.PaymentMethod, paymentResp.PaymentId, req.UpdatedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to update payment status: %w", err)
		}

		// Update ticket status to confirmed
		err = s.ticketRepo.UpdateStatus(ctx, req.TicketID, models.TicketStatusConfirmed, req.UpdatedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm ticket: %w", err)
		}

		// Generate QR code and barcode
		if err := s.generateTicketCodes(ctx, ticket); err != nil {
			s.logger.Warn("Failed to generate ticket codes",
				zap.String("ticket_id", req.TicketID),
				zap.Error(err),
			)
		}

		// Increment metrics
		metrics.IncrementPaymentProcessed(ticket.EventID, "success", req.PaymentMethod)

		return &ProcessPaymentResponse{
			Success:       true,
			PaymentID:     paymentResp.PaymentId,
			TicketStatus:  models.TicketStatusConfirmed,
			PaymentStatus: models.PaymentStatusPaid,
		}, nil
	}

	return nil, fmt.Errorf("payment service not available")
}

// CancelTicket cancels a ticket
func (s *TicketService) CancelTicket(ctx context.Context, req *CancelTicketRequest) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Validate ticket can be cancelled
	if !ticket.CanBeCancelled() {
		return fmt.Errorf("ticket cannot be cancelled, current status: %s", ticket.Status)
	}

	// Cancel ticket
	err = s.ticketRepo.Delete(ctx, req.TicketID, req.Reason, req.CancelledBy)
	if err != nil {
		return fmt.Errorf("failed to cancel ticket: %w", err)
	}

	// Release seat in Event Service
	if s.eventClient != nil {
		err = s.eventClient.UpdateSeatAvailability(ctx, ticket.EventID, ticket.SeatID, "available", "", "")
		if err != nil {
			s.logger.Warn("Failed to release seat in Event Service",
				zap.String("event_id", ticket.EventID),
				zap.String("seat_id", ticket.SeatID),
				zap.Error(err),
			)
		}
	}

	// Process refund if ticket was paid
	if ticket.IsPaid() {
		if err := s.processRefund(ctx, ticket, req.Reason); err != nil {
			s.logger.Error("Failed to process refund",
				zap.String("ticket_id", req.TicketID),
				zap.Error(err),
			)
		}
	}

	// Increment metrics
	metrics.IncrementTicketCancelled(ticket.EventID, req.Reason)

	s.logger.Info("Ticket cancelled successfully",
		zap.String("ticket_id", req.TicketID),
		zap.String("reason", req.Reason),
		zap.String("cancelled_by", req.CancelledBy),
	)

	return nil
}

// SearchTickets searches tickets with filters
func (s *TicketService) SearchTickets(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*models.Ticket, int, error) {
	tickets, total, err := s.ticketRepo.Search(ctx, filters, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search tickets: %w", err)
	}

	return tickets, total, nil
}

// Helper methods

func (s *TicketService) validateCreateTicketRequest(req *CreateTicketRequest) error {
	if req.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if req.SeatID == "" {
		return fmt.Errorf("seat_id is required")
	}
	if req.ZoneID == "" {
		return fmt.Errorf("zone_id is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.BasePrice < 0 {
		return fmt.Errorf("base_price cannot be negative")
	}
	if req.FinalPrice < 0 {
		return fmt.Errorf("final_price cannot be negative")
	}
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}

func (s *TicketService) checkSeatAvailability(ctx context.Context, eventID, seatID string) (bool, error) {
	resp, err := s.eventClient.GetSeatAvailability(ctx, eventID, seatID)
	if err != nil {
		return false, err
	}
	return resp.Available, nil
}

func (s *TicketService) generateTicketNumber(eventID, seatID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("TKT-%s-%s-%d", eventID[:8], seatID[:8], timestamp)
}

func (s *TicketService) isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		models.TicketStatusPending:   {models.TicketStatusConfirmed, models.TicketStatusCancelled},
		models.TicketStatusConfirmed: {models.TicketStatusUsed, models.TicketStatusCancelled, models.TicketStatusRefunded},
		models.TicketStatusUsed:      {}, // Terminal state
		models.TicketStatusCancelled: {}, // Terminal state
		models.TicketStatusRefunded:  {}, // Terminal state
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}

func (s *TicketService) handleTicketConfirmation(ctx context.Context, ticket *models.Ticket) {
	// Generate QR code and barcode
	if err := s.generateTicketCodes(ctx, ticket); err != nil {
		s.logger.Error("Failed to generate ticket codes",
			zap.String("ticket_id", ticket.ID),
			zap.Error(err),
		)
	}
}

func (s *TicketService) handleTicketCancellation(ctx context.Context, ticket *models.Ticket) {
	// Release seat in Event Service
	if s.eventClient != nil {
		err := s.eventClient.UpdateSeatAvailability(ctx, ticket.EventID, ticket.SeatID, "available", "", "")
		if err != nil {
			s.logger.Warn("Failed to release seat in Event Service",
				zap.String("event_id", ticket.EventID),
				zap.String("seat_id", ticket.SeatID),
				zap.Error(err),
			)
		}
	}
}

func (s *TicketService) handleTicketUsage(ctx context.Context, ticket *models.Ticket) {
	// Mark ticket as used
	now := time.Now()
	ticket.UsedAt = &now
	// Update in database would be handled by the calling method
}

func (s *TicketService) generateTicketCodes(ctx context.Context, ticket *models.Ticket) error {
	// Generate QR code
	qrCode := fmt.Sprintf("TICKET:%s:%s:%s", ticket.ID, ticket.EventID, ticket.SeatID)
	ticket.QRCode = &qrCode

	// Generate barcode
	barcode := fmt.Sprintf("%s%s%s", ticket.ID[:8], ticket.EventID[:8], ticket.SeatID[:8])
	ticket.Barcode = &barcode

	// Update ticket in database
	return s.ticketRepo.Update(ctx, ticket)
}

func (s *TicketService) processRefund(ctx context.Context, ticket *models.Ticket, reason string) error {
	if s.paymentClient == nil {
		return fmt.Errorf("payment service not available")
	}

	refundReq := &grpcclient.RefundPaymentRequest{
		PaymentId: *ticket.PaymentReference,
		Amount:    ticket.FinalPrice,
		Reason:    reason,
	}

	refundResp, err := s.paymentClient.RefundPayment(ctx, refundReq)
	if err != nil {
		return fmt.Errorf("refund processing failed: %w", err)
	}

	if refundResp.Success {
		// Update ticket status to refunded
		err = s.ticketRepo.UpdateStatus(ctx, ticket.ID, models.TicketStatusRefunded, "system")
		if err != nil {
			return fmt.Errorf("failed to update ticket status: %w", err)
		}

		// Update payment status
		err = s.ticketRepo.UpdatePaymentStatus(ctx, ticket.ID, models.PaymentStatusRefunded, "", "", "system")
		if err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		// Increment metrics
		metrics.IncrementTicketRefunded(ticket.EventID, reason)
	}

	return nil
}

// Request/Response types

type CreateTicketRequest struct {
	EventID          string     `json:"event_id"`
	SeatID           string     `json:"seat_id"`
	ZoneID           string     `json:"zone_id"`
	UserID           string     `json:"user_id"`
	BookingSessionID string     `json:"booking_session_id,omitempty"`
	TicketType       string     `json:"ticket_type,omitempty"`
	PricingCategory  string     `json:"pricing_category"`
	BasePrice        float64    `json:"base_price"`
	FinalPrice       float64    `json:"final_price"`
	Currency         string     `json:"currency"`
	DiscountReason   string     `json:"discount_reason,omitempty"`
	ValidUntil       *time.Time `json:"valid_until,omitempty"`
	CreatedBy        string     `json:"created_by,omitempty"`
}

type ProcessPaymentRequest struct {
	TicketID      string `json:"ticket_id"`
	PaymentMethod string `json:"payment_method"`
	UpdatedBy     string `json:"updated_by"`
}

type ProcessPaymentResponse struct {
	Success       bool   `json:"success"`
	PaymentID     string `json:"payment_id"`
	TicketStatus  string `json:"ticket_status"`
	PaymentStatus string `json:"payment_status"`
}

type CancelTicketRequest struct {
	TicketID    string `json:"ticket_id"`
	Reason      string `json:"reason"`
	CancelledBy string `json:"cancelled_by"`
}

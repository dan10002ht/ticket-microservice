package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ticketpb "checkin-service/internal/protos/ticket"
	"checkin-service/grpcclient"
	"checkin-service/models"
	"checkin-service/repositories"
)

// CheckInError represents domain-level check-in failures.
type CheckInError struct {
	Code    string
	Message string
}

func (e *CheckInError) Error() string { return e.Message }

var (
	ErrInvalidTicket       = &CheckInError{"INVALID_TICKET", "ticket not found or invalid"}
	ErrAlreadyCheckedIn    = &CheckInError{"ALREADY_CHECKED_IN", "ticket has already been checked in"}
	ErrCancelledTicket     = &CheckInError{"CANCELLED_TICKET", "ticket has been cancelled"}
	ErrTicketEventMismatch = &CheckInError{"TICKET_EVENT_MISMATCH", "ticket does not belong to this event"}
)

type CheckinService struct {
	repo    *repositories.CheckinRepository
	clients *grpcclient.Clients
	logger  *zap.Logger
}

func NewCheckinService(repo *repositories.CheckinRepository, clients *grpcclient.Clients, logger *zap.Logger) *CheckinService {
	return &CheckinService{repo: repo, clients: clients, logger: logger}
}

// CheckIn validates and records an event check-in.
//
// Flow:
//  1. Validate ticket via ticket-service gRPC
//  2. Ensure ticket belongs to the given event
//  3. Ensure ticket status is "sold"
//  4. Idempotency check — reject duplicates
//  5. Mark ticket as "used" via ticket-service gRPC
//  6. Persist check-in record
func (s *CheckinService) CheckIn(ctx context.Context, req *models.CheckIn) (*models.CheckIn, error) {
	s.logger.Info("Processing check-in",
		zap.String("ticket_id", req.TicketID),
		zap.String("event_id", req.EventID),
	)

	// --- Step 1: Validate ticket via ticket-service ---
	ticketResp, err := s.clients.Ticket.GetTicket(ctx, &ticketpb.GetTicketRequest{
		TicketId: req.TicketID,
	})
	if err != nil {
		grpcStatus, _ := status.FromError(err)
		if grpcStatus.Code() == codes.NotFound {
			return nil, ErrInvalidTicket
		}
		return nil, fmt.Errorf("GetTicket RPC: %w", err)
	}
	if !ticketResp.Success || ticketResp.Ticket == nil {
		return nil, ErrInvalidTicket
	}

	ticket := ticketResp.Ticket

	// --- Step 2: Check ticket belongs to the event ---
	if ticket.EventId != req.EventID {
		return nil, ErrTicketEventMismatch
	}

	// --- Step 3: Check ticket status ---
	switch ticket.Status {
	case "cancelled", "refunded":
		return nil, ErrCancelledTicket
	case "sold", "confirmed":
		// OK to check in
	default:
		return nil, ErrInvalidTicket
	}

	// --- Step 4: Idempotency check ---
	exists, err := s.repo.ExistsByTicketID(req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("check duplicate: %w", err)
	}
	if exists {
		s.logger.Warn("Duplicate check-in attempt", zap.String("ticket_id", req.TicketID))
		return nil, ErrAlreadyCheckedIn
	}

	// --- Step 5: Mark ticket as used ---
	_, err = s.clients.Ticket.UpdateTicket(ctx, &ticketpb.UpdateTicketRequest{
		TicketId: req.TicketID,
		Status:   "used",
	})
	if err != nil {
		s.logger.Error("Failed to mark ticket as used", zap.String("ticket_id", req.TicketID), zap.Error(err))
		return nil, fmt.Errorf("UpdateTicket RPC: %w", err)
	}

	// Fill user_id from ticket if not provided
	if req.UserID == "" {
		req.UserID = ticket.UserId
	}

	// --- Step 6: Persist check-in record ---
	req.Status = models.CheckInStatusSuccess
	record, err := s.repo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("persist check-in: %w", err)
	}

	s.logger.Info("Check-in recorded",
		zap.String("checkin_id", record.ID),
		zap.String("ticket_id", req.TicketID),
	)
	return record, nil
}

// GetCheckIn retrieves a single check-in record by ID.
func (s *CheckinService) GetCheckIn(ctx context.Context, id string) (*models.CheckIn, error) {
	return s.repo.GetByID(id)
}

// ListCheckIns returns paginated check-in records for an event.
func (s *CheckinService) ListCheckIns(ctx context.Context, eventID, userID, gate string, page, limit int) ([]*models.CheckIn, int, error) {
	return s.repo.ListByEvent(eventID, userID, gate, page, limit)
}

// GetEventStats returns aggregated check-in stats for an event.
func (s *CheckinService) GetEventStats(ctx context.Context, eventID string) (*models.EventStats, error) {
	return s.repo.GetEventStats(eventID)
}

package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	paymentpb "shared-lib/protos/payment"
	"ticket-service/grpcclient"
	"ticket-service/metrics"
	"ticket-service/models"
	"ticket-service/repositories"
)

// TicketBookingSessionService handles booking session business logic within
// ticket-service. It is intentionally named to distinguish it from the Booking
// Service microservice that orchestrates cross-service workflows.
type TicketBookingSessionService struct {
	bookingRepo     *repositories.BookingSessionRepository
	reservationRepo *repositories.SeatReservationRepository
	eventClient     *grpcclient.EventServiceClient
	paymentClient   *grpcclient.PaymentServiceClient
	logger          *zap.Logger
}

// NewTicketBookingSessionService creates a new ticket booking session service
func NewTicketBookingSessionService(
	bookingRepo *repositories.BookingSessionRepository,
	reservationRepo *repositories.SeatReservationRepository,
	eventClient *grpcclient.EventServiceClient,
	paymentClient *grpcclient.PaymentServiceClient,
	logger *zap.Logger,
) *TicketBookingSessionService {
	return &TicketBookingSessionService{
		bookingRepo:     bookingRepo,
		reservationRepo: reservationRepo,
		eventClient:     eventClient,
		paymentClient:   paymentClient,
		logger:          logger,
	}
}

// CreateBookingSession creates a new booking session
func (s *TicketBookingSessionService) CreateBookingSession(ctx context.Context, req *BookingSessionCreateCommand) (*models.BookingSession, error) {
	// Validate request
	if err := s.validateCreateBookingSessionRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Generate session token
	sessionToken := s.generateSessionToken(req.UserID)

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(req.TimeoutMinutes) * time.Minute)

	// Create booking session
	session := models.NewBookingSession(
		req.UserID, req.EventID, sessionToken, 0, 0, req.Currency, expiresAt,
	)

	// Set additional fields
	if req.IPAddress != "" {
		session.IPAddress = &req.IPAddress
	}
	if req.UserAgent != "" {
		session.UserAgent = &req.UserAgent
	}
	if req.CreatedBy != "" {
		session.CreatedBy = &req.CreatedBy
	}

	// Validate session
	if err := session.Validate(); err != nil {
		return nil, fmt.Errorf("session validation failed: %w", err)
	}

	// Create session in database
	if err := s.bookingRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create booking session: %w", err)
	}

	// Increment metrics
	metrics.IncrementBookingSessionCreated(req.EventID, session.Status)

	s.logger.Info("Booking session created successfully",
		zap.String("session_id", session.ID),
		zap.String("event_id", req.EventID),
		zap.String("user_id", req.UserID),
	)

	return session, nil
}

// GetBookingSession retrieves a booking session by ID
func (s *TicketBookingSessionService) GetBookingSession(ctx context.Context, sessionID string) (*models.BookingSession, error) {
	session, err := s.bookingRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking session: %w", err)
	}

	return session, nil
}

// GetBookingSessionByToken retrieves a booking session by token
func (s *TicketBookingSessionService) GetBookingSessionByToken(ctx context.Context, sessionToken string) (*models.BookingSession, error) {
	session, err := s.bookingRepo.GetBySessionToken(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking session: %w", err)
	}

	return session, nil
}

// AddSeatToSession adds a seat to the booking session
func (s *TicketBookingSessionService) AddSeatToSession(ctx context.Context, req *BookingSessionAddSeatCommand) error {
	// Get booking session
	session, err := s.bookingRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get booking session: %w", err)
	}

	// Validate session is active
	if !session.IsActive() {
		return fmt.Errorf("booking session is not active, status: %s", session.Status)
	}

	// Check seat availability
	if s.eventClient != nil {
		available, err := s.checkSeatAvailability(ctx, req.EventID, req.SeatID)
		if err != nil {
			return fmt.Errorf("failed to check seat availability: %w", err)
		}
		if !available {
			return fmt.Errorf("seat %s is not available", req.SeatID)
		}
	}

	// Create seat reservation
	reservation := models.NewSeatReservation(
		req.SessionID, req.EventID, req.SeatID, req.ZoneID,
		s.generateReservationToken(req.SessionID, req.SeatID),
		req.PricingCategory, req.BasePrice, req.FinalPrice, req.Currency,
		session.ExpiresAt,
	)

	if req.CreatedBy != "" {
		reservation.CreatedBy = &req.CreatedBy
	}

	// Validate reservation
	if err := reservation.Validate(); err != nil {
		return fmt.Errorf("reservation validation failed: %w", err)
	}

	// Create reservation in database
	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return fmt.Errorf("failed to create seat reservation: %w", err)
	}

	// Block seat in Event Service
	if s.eventClient != nil {
		_, err := s.eventClient.BlockSeats(ctx, req.EventID, []string{req.SeatID},
			session.UserID, req.SessionID, session.ExpiresAt)
		if err != nil {
			s.logger.Warn("Failed to block seat in Event Service",
				zap.String("event_id", req.EventID),
				zap.String("seat_id", req.SeatID),
				zap.Error(err),
			)
		}
	}

	// Update session totals
	session.SeatCount++
	session.TotalAmount += req.FinalPrice

	if err := s.bookingRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update booking session: %w", err)
	}

	// Increment metrics
	metrics.IncrementSeatReservationCreated(req.EventID, req.ZoneID)

	s.logger.Info("Seat added to booking session",
		zap.String("session_id", req.SessionID),
		zap.String("seat_id", req.SeatID),
		zap.Int("seat_count", session.SeatCount),
		zap.Float64("total_amount", session.TotalAmount),
	)

	return nil
}

// RemoveSeatFromSession removes a seat from the booking session
func (s *TicketBookingSessionService) RemoveSeatFromSession(ctx context.Context, req *BookingSessionRemoveSeatCommand) error {
	// Get booking session
	session, err := s.bookingRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get booking session: %w", err)
	}

	// Validate session is active
	if !session.IsActive() {
		return fmt.Errorf("booking session is not active, status: %s", session.Status)
	}

	// Get seat reservation
	reservations, err := s.reservationRepo.GetByBookingSessionID(ctx, req.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get seat reservations: %w", err)
	}

	var targetReservation *models.SeatReservation
	for _, reservation := range reservations {
		if reservation.SeatID == req.SeatID {
			targetReservation = reservation
			break
		}
	}

	if targetReservation == nil {
		return fmt.Errorf("seat %s not found in session", req.SeatID)
	}

	// Release seat reservation
	err = s.reservationRepo.Release(ctx, targetReservation.ID, req.Reason, req.RemovedBy)
	if err != nil {
		return fmt.Errorf("failed to release seat reservation: %w", err)
	}

	// Release seat in Event Service
	if s.eventClient != nil {
		_, err := s.eventClient.ReleaseSeats(ctx, targetReservation.EventID,
			[]string{req.SeatID}, session.UserID, req.Reason)
		if err != nil {
			s.logger.Warn("Failed to release seat in Event Service",
				zap.String("event_id", targetReservation.EventID),
				zap.String("seat_id", req.SeatID),
				zap.Error(err),
			)
		}
	}

	// Update session totals
	session.SeatCount--
	session.TotalAmount -= targetReservation.FinalPrice

	if err := s.bookingRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update booking session: %w", err)
	}

	// Increment metrics
	metrics.IncrementSeatReservationReleased(targetReservation.EventID, req.Reason)

	s.logger.Info("Seat removed from booking session",
		zap.String("session_id", req.SessionID),
		zap.String("seat_id", req.SeatID),
		zap.Int("seat_count", session.SeatCount),
		zap.Float64("total_amount", session.TotalAmount),
	)

	return nil
}

// CompleteBookingSession completes a booking session
func (s *TicketBookingSessionService) CompleteBookingSession(ctx context.Context, req *BookingSessionCompleteCommand) (*BookingSessionCompleteResult, error) {
	// Get booking session
	session, err := s.bookingRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking session: %w", err)
	}

	// Validate session is active
	if !session.IsActive() {
		return nil, fmt.Errorf("booking session is not active, status: %s", session.Status)
	}

	// Get seat reservations
	reservations, err := s.reservationRepo.GetByBookingSessionID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	if len(reservations) == 0 {
		return nil, fmt.Errorf("no seats reserved in session")
	}

	// Process payment if required
	var paymentID string
	if session.TotalAmount > 0 && s.paymentClient != nil {
		paymentReq := &paymentpb.ProcessPaymentRequest{
			BookingId:     req.SessionID,
			Amount:        session.TotalAmount,
			Currency:      session.Currency,
			PaymentMethod: req.PaymentMethod,
			UserId:        session.UserID,
			EventId:       session.EventID,
		}

		paymentResp, err := s.paymentClient.ProcessPayment(ctx, paymentReq)
		if err != nil {
			return nil, fmt.Errorf("payment processing failed: %w", err)
		}

		paymentID = paymentResp.PaymentId
	}

	// Confirm all seat reservations
	for _, reservation := range reservations {
		err = s.reservationRepo.Confirm(ctx, reservation.ID, req.CompletedBy)
		if err != nil {
			s.logger.Error("Failed to confirm seat reservation",
				zap.String("reservation_id", reservation.ID),
				zap.Error(err),
			)
		}
	}

	// Complete booking session
	err = s.bookingRepo.Complete(ctx, req.SessionID, req.CompletedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to complete booking session: %w", err)
	}

	// Increment metrics
	metrics.IncrementBookingSessionCompleted(session.EventID)
	metrics.IncrementSeatReservationConfirmed(session.EventID, reservations[0].ZoneID)

	s.logger.Info("Booking session completed successfully",
		zap.String("session_id", req.SessionID),
		zap.String("payment_id", paymentID),
		zap.Int("seat_count", len(reservations)),
	)

	return &BookingSessionCompleteResult{
		Success:     true,
		PaymentID:   paymentID,
		SeatCount:   len(reservations),
		TotalAmount: session.TotalAmount,
	}, nil
}

// CancelBookingSession cancels a booking session
func (s *TicketBookingSessionService) CancelBookingSession(ctx context.Context, req *BookingSessionCancelCommand) error {
	// Get booking session
	session, err := s.bookingRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get booking session: %w", err)
	}

	// Validate session can be cancelled
	if session.IsCompleted() {
		return fmt.Errorf("booking session is already completed")
	}

	// Release all seat reservations
	err = s.reservationRepo.ReleaseByBookingSession(ctx, req.SessionID, req.Reason, req.CancelledBy)
	if err != nil {
		s.logger.Error("Failed to release seat reservations",
			zap.String("session_id", req.SessionID),
			zap.Error(err),
		)
	}

	// Release seats in Event Service
	if s.eventClient != nil {
		reservations, err := s.reservationRepo.GetByBookingSessionID(ctx, req.SessionID)
		if err == nil && len(reservations) > 0 {
			seatIDs := make([]string, len(reservations))
			for i, reservation := range reservations {
				seatIDs[i] = reservation.SeatID
			}

			_, err := s.eventClient.ReleaseSeats(ctx, session.EventID, seatIDs, session.UserID, req.Reason)
			if err != nil {
				s.logger.Warn("Failed to release seats in Event Service",
					zap.String("event_id", session.EventID),
					zap.Error(err),
				)
			}
		}
	}

	// Cancel booking session
	err = s.bookingRepo.Cancel(ctx, req.SessionID, req.Reason, req.CancelledBy)
	if err != nil {
		return fmt.Errorf("failed to cancel booking session: %w", err)
	}

	// Increment metrics
	metrics.IncrementBookingSessionExpired(session.EventID)

	s.logger.Info("Booking session cancelled successfully",
		zap.String("session_id", req.SessionID),
		zap.String("reason", req.Reason),
		zap.String("cancelled_by", req.CancelledBy),
	)

	return nil
}

// GetSessionReservations gets seat reservations for a session
func (s *TicketBookingSessionService) GetSessionReservations(ctx context.Context, sessionID string) ([]*models.SeatReservation, error) {
	reservations, err := s.reservationRepo.GetByBookingSessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session reservations: %w", err)
	}

	return reservations, nil
}

// CleanupExpiredSessions cleans up expired booking sessions
func (s *TicketBookingSessionService) CleanupExpiredSessions(ctx context.Context) error {
	expiredSessions, err := s.bookingRepo.GetExpiredSessions(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get expired sessions: %w", err)
	}

	for _, session := range expiredSessions {
		req := &BookingSessionCancelCommand{
			SessionID:   session.ID,
			Reason:      "Session expired",
			CancelledBy: "system",
		}

		if err := s.CancelBookingSession(ctx, req); err != nil {
			s.logger.Error("Failed to cleanup expired session",
				zap.String("session_id", session.ID),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("Cleanup expired sessions completed",
		zap.Int("expired_count", len(expiredSessions)),
	)

	return nil
}

// Helper methods

func (s *TicketBookingSessionService) validateCreateBookingSessionRequest(req *BookingSessionCreateCommand) error {
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if req.TimeoutMinutes <= 0 {
		return fmt.Errorf("timeout_minutes must be positive")
	}
	return nil
}

func (s *TicketBookingSessionService) generateSessionToken(userID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("BKG-%s-%d", userID[:8], timestamp)
}

func (s *TicketBookingSessionService) generateReservationToken(sessionID, seatID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("RSV-%s-%s-%d", sessionID[:8], seatID[:8], timestamp)
}

func (s *TicketBookingSessionService) checkSeatAvailability(ctx context.Context, eventID, seatID string) (bool, error) {
	resp, err := s.eventClient.GetSeatAvailability(ctx, eventID, seatID)
	if err != nil {
		return false, err
	}
	return resp.Available, nil
}

// Request/Response types

type BookingSessionCreateCommand struct {
	UserID         string `json:"user_id"`
	EventID        string `json:"event_id"`
	Currency       string `json:"currency"`
	TimeoutMinutes int    `json:"timeout_minutes"`
	IPAddress      string `json:"ip_address,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
}

type BookingSessionAddSeatCommand struct {
	SessionID       string  `json:"session_id"`
	EventID         string  `json:"event_id"`
	SeatID          string  `json:"seat_id"`
	ZoneID          string  `json:"zone_id"`
	PricingCategory string  `json:"pricing_category"`
	BasePrice       float64 `json:"base_price"`
	FinalPrice      float64 `json:"final_price"`
	Currency        string  `json:"currency"`
	CreatedBy       string  `json:"created_by,omitempty"`
}

type BookingSessionRemoveSeatCommand struct {
	SessionID string `json:"session_id"`
	SeatID    string `json:"seat_id"`
	Reason    string `json:"reason"`
	RemovedBy string `json:"removed_by"`
}

type BookingSessionCompleteCommand struct {
	SessionID     string `json:"session_id"`
	PaymentMethod string `json:"payment_method"`
	CompletedBy   string `json:"completed_by"`
}

type BookingSessionCompleteResult struct {
	Success     bool    `json:"success"`
	PaymentID   string  `json:"payment_id"`
	SeatCount   int     `json:"seat_count"`
	TotalAmount float64 `json:"total_amount"`
}

type BookingSessionCancelCommand struct {
	SessionID   string `json:"session_id"`
	Reason      string `json:"reason"`
	CancelledBy string `json:"cancelled_by"`
}

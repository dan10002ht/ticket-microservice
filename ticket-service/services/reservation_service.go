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

// ReservationService handles seat reservation business logic
type ReservationService struct {
	reservationRepo *repositories.SeatReservationRepository
	eventClient     *grpcclient.EventServiceClient
	logger          *zap.Logger
}

// NewReservationService creates a new reservation service
func NewReservationService(
	reservationRepo *repositories.SeatReservationRepository,
	eventClient *grpcclient.EventServiceClient,
	logger *zap.Logger,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		eventClient:     eventClient,
		logger:          logger,
	}
}

// CreateReservation creates a new seat reservation
func (s *ReservationService) CreateReservation(ctx context.Context, req *CreateReservationRequest) (*models.SeatReservation, error) {
	// Validate request
	if err := s.validateCreateReservationRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check seat availability
	if s.eventClient != nil {
		available, err := s.checkSeatAvailability(ctx, req.EventID, req.SeatID)
		if err != nil {
			return nil, fmt.Errorf("failed to check seat availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("seat %s is not available", req.SeatID)
		}
	}

	// Generate reservation token
	reservationToken := s.generateReservationToken(req.BookingSessionID, req.SeatID)

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(req.TimeoutMinutes) * time.Minute)

	// Create seat reservation
	reservation := models.NewSeatReservation(
		req.BookingSessionID, req.EventID, req.SeatID, req.ZoneID,
		reservationToken, req.PricingCategory, req.BasePrice, req.FinalPrice, req.Currency, expiresAt,
	)

	// Set additional fields
	if req.CreatedBy != "" {
		reservation.CreatedBy = &req.CreatedBy
	}

	// Validate reservation
	if err := reservation.Validate(); err != nil {
		return nil, fmt.Errorf("reservation validation failed: %w", err)
	}

	// Create reservation in database
	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to create seat reservation: %w", err)
	}

	// Block seat in Event Service
	if s.eventClient != nil {
		blockedReason := fmt.Sprintf("Reserved for session %s by user %s", req.BookingSessionID, req.UserID)
		_, err := s.eventClient.BlockSeats(ctx, req.EventID, []string{req.SeatID},
			blockedReason, expiresAt)
		if err != nil {
			s.logger.Warn("Failed to block seat in Event Service",
				zap.String("event_id", req.EventID),
				zap.String("seat_id", req.SeatID),
				zap.Error(err),
			)
		}
	}

	// Increment metrics
	metrics.IncrementSeatReservationCreated(req.EventID, req.ZoneID)

	s.logger.Info("Seat reservation created successfully",
		zap.String("reservation_id", reservation.ID),
		zap.String("event_id", req.EventID),
		zap.String("seat_id", req.SeatID),
	)

	return reservation, nil
}

// GetReservation retrieves a seat reservation by ID
func (s *ReservationService) GetReservation(ctx context.Context, reservationID string) (*models.SeatReservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat reservation: %w", err)
	}

	return reservation, nil
}

// GetReservationByToken retrieves a seat reservation by token
func (s *ReservationService) GetReservationByToken(ctx context.Context, reservationToken string) (*models.SeatReservation, error) {
	reservation, err := s.reservationRepo.GetByReservationToken(ctx, reservationToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat reservation: %w", err)
	}

	return reservation, nil
}

// GetReservationsBySession gets seat reservations for a booking session
func (s *ReservationService) GetReservationsBySession(ctx context.Context, sessionID string) ([]*models.SeatReservation, error) {
	reservations, err := s.reservationRepo.GetByBookingSessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session reservations: %w", err)
	}

	return reservations, nil
}

// GetReservationsByEvent gets seat reservations for an event
func (s *ReservationService) GetReservationsByEvent(ctx context.Context, eventID string, page, limit int) ([]*models.SeatReservation, int, error) {
	reservations, total, err := s.reservationRepo.GetByEventID(ctx, eventID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get event reservations: %w", err)
	}

	return reservations, total, nil
}

// GetReservationsBySeat gets seat reservations for a specific seat
func (s *ReservationService) GetReservationsBySeat(ctx context.Context, seatID string) ([]*models.SeatReservation, error) {
	reservations, err := s.reservationRepo.GetBySeatID(ctx, seatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	return reservations, nil
}

// GetActiveReservations gets active seat reservations for an event
func (s *ReservationService) GetActiveReservations(ctx context.Context, eventID string) ([]*models.SeatReservation, error) {
	reservations, err := s.reservationRepo.GetActiveReservations(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reservations: %w", err)
	}

	return reservations, nil
}

// ConfirmReservation confirms a seat reservation
func (s *ReservationService) ConfirmReservation(ctx context.Context, req *ConfirmReservationRequest) error {
	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, req.ReservationID)
	if err != nil {
		return fmt.Errorf("failed to get seat reservation: %w", err)
	}

	// Validate reservation can be confirmed
	if !reservation.IsReserved() {
		return fmt.Errorf("reservation cannot be confirmed, current status: %s", reservation.Status)
	}

	// Confirm reservation
	err = s.reservationRepo.Confirm(ctx, req.ReservationID, req.ConfirmedBy)
	if err != nil {
		return fmt.Errorf("failed to confirm reservation: %w", err)
	}

	// Increment metrics
	metrics.IncrementSeatReservationConfirmed(reservation.EventID, reservation.ZoneID)

	s.logger.Info("Seat reservation confirmed successfully",
		zap.String("reservation_id", req.ReservationID),
		zap.String("confirmed_by", req.ConfirmedBy),
	)

	return nil
}

// ReleaseReservation releases a seat reservation
func (s *ReservationService) ReleaseReservation(ctx context.Context, req *ReleaseReservationRequest) error {
	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, req.ReservationID)
	if err != nil {
		return fmt.Errorf("failed to get seat reservation: %w", err)
	}

	// Validate reservation can be released
	if reservation.IsReleased() {
		return fmt.Errorf("reservation is already released")
	}

	// Release reservation
	err = s.reservationRepo.Release(ctx, req.ReservationID, req.Reason, req.ReleasedBy)
	if err != nil {
		return fmt.Errorf("failed to release reservation: %w", err)
	}

	// Release seat in Event Service
	if s.eventClient != nil {
		_, err := s.eventClient.ReleaseSeats(ctx, reservation.EventID,
			[]string{reservation.SeatID})
		if err != nil {
			s.logger.Warn("Failed to release seat in Event Service",
				zap.String("event_id", reservation.EventID),
				zap.String("seat_id", reservation.SeatID),
				zap.Error(err),
			)
		}
	}

	// Increment metrics
	metrics.IncrementSeatReservationReleased(reservation.EventID, req.Reason)

	s.logger.Info("Seat reservation released successfully",
		zap.String("reservation_id", req.ReservationID),
		zap.String("reason", req.Reason),
		zap.String("released_by", req.ReleasedBy),
	)

	return nil
}

// ReleaseReservationsBySession releases all reservations for a booking session
func (s *ReservationService) ReleaseReservationsBySession(ctx context.Context, req *ReleaseReservationsBySessionRequest) error {
	// Get reservations
	reservations, err := s.reservationRepo.GetByBookingSessionID(ctx, req.SessionID)
	if err != nil {
		return fmt.Errorf("failed to get session reservations: %w", err)
	}

	if len(reservations) == 0 {
		return nil // No reservations to release
	}

	// Release all reservations
	err = s.reservationRepo.ReleaseByBookingSession(ctx, req.SessionID, req.Reason, req.ReleasedBy)
	if err != nil {
		return fmt.Errorf("failed to release session reservations: %w", err)
	}

	// Release seats in Event Service
	if s.eventClient != nil && len(reservations) > 0 {
		seatIDs := make([]string, len(reservations))
		for i, reservation := range reservations {
			seatIDs[i] = reservation.SeatID
		}

		_, err := s.eventClient.ReleaseSeats(ctx, reservations[0].EventID, seatIDs)
		if err != nil {
			s.logger.Warn("Failed to release seats in Event Service",
				zap.String("event_id", reservations[0].EventID),
				zap.Error(err),
			)
		}
	}

	// Increment metrics
	for _, reservation := range reservations {
		metrics.IncrementSeatReservationReleased(reservation.EventID, req.Reason)
	}

	s.logger.Info("Session reservations released successfully",
		zap.String("session_id", req.SessionID),
		zap.String("reason", req.Reason),
		zap.Int("reservation_count", len(reservations)),
	)

	return nil
}

// CheckSeatAvailability checks if a seat is available for reservation
func (s *ReservationService) CheckSeatAvailability(ctx context.Context, seatID string) (bool, error) {
	available, err := s.reservationRepo.CheckSeatAvailability(ctx, seatID)
	if err != nil {
		return false, fmt.Errorf("failed to check seat availability: %w", err)
	}

	return available, nil
}

// GetReservationStats gets reservation statistics for an event
func (s *ReservationService) GetReservationStats(ctx context.Context, eventID string) (map[string]int, error) {
	stats, err := s.reservationRepo.GetReservationStats(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation stats: %w", err)
	}

	return stats, nil
}

// CleanupExpiredReservations cleans up expired seat reservations
func (s *ReservationService) CleanupExpiredReservations(ctx context.Context) error {
	expiredReservations, err := s.reservationRepo.GetExpiredReservations(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get expired reservations: %w", err)
	}

	for _, reservation := range expiredReservations {
		req := &ReleaseReservationRequest{
			ReservationID: reservation.ID,
			Reason:        "Reservation expired",
			ReleasedBy:    "system",
		}

		if err := s.ReleaseReservation(ctx, req); err != nil {
			s.logger.Error("Failed to cleanup expired reservation",
				zap.String("reservation_id", reservation.ID),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("Cleanup expired reservations completed",
		zap.Int("expired_count", len(expiredReservations)),
	)

	return nil
}

// ExtendReservation extends the expiration time of a reservation
func (s *ReservationService) ExtendReservation(ctx context.Context, req *ExtendReservationRequest) error {
	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, req.ReservationID)
	if err != nil {
		return fmt.Errorf("failed to get seat reservation: %w", err)
	}

	// Validate reservation can be extended
	if !reservation.IsReserved() {
		return fmt.Errorf("reservation cannot be extended, current status: %s", reservation.Status)
	}

	// Calculate new expiration time
	newExpiresAt := time.Now().Add(time.Duration(req.ExtensionMinutes) * time.Minute)

	// Update reservation expiration
	reservation.ExpiresAt = newExpiresAt
	reservation.UpdatedBy = &req.ExtendedBy

	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("failed to extend reservation: %w", err)
	}

	// Update seat block in Event Service
	if s.eventClient != nil {
		blockedReason := fmt.Sprintf("Extended reservation for session %s", reservation.BookingSessionID)
		_, err := s.eventClient.BlockSeats(ctx, reservation.EventID, []string{reservation.SeatID},
			blockedReason, newExpiresAt)
		if err != nil {
			s.logger.Warn("Failed to update seat block in Event Service",
				zap.String("event_id", reservation.EventID),
				zap.String("seat_id", reservation.SeatID),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("Seat reservation extended successfully",
		zap.String("reservation_id", req.ReservationID),
		zap.Int("extension_minutes", req.ExtensionMinutes),
		zap.String("extended_by", req.ExtendedBy),
	)

	return nil
}

// Helper methods

func (s *ReservationService) validateCreateReservationRequest(req *CreateReservationRequest) error {
	if req.BookingSessionID == "" {
		return fmt.Errorf("booking_session_id is required")
	}
	if req.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if req.SeatID == "" {
		return fmt.Errorf("seat_id is required")
	}
	if req.ZoneID == "" {
		return fmt.Errorf("zone_id is required")
	}
	if req.PricingCategory == "" {
		return fmt.Errorf("pricing_category is required")
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
	if req.TimeoutMinutes <= 0 {
		return fmt.Errorf("timeout_minutes must be positive")
	}
	return nil
}

func (s *ReservationService) generateReservationToken(sessionID, seatID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("RSV-%s-%s-%d", sessionID[:8], seatID[:8], timestamp)
}

func (s *ReservationService) checkSeatAvailability(ctx context.Context, eventID, seatID string) (bool, error) {
	resp, err := s.eventClient.GetSeatAvailability(ctx, eventID, seatID)
	if err != nil {
		return false, err
	}
	if resp.Availability == nil {
		return false, nil
	}
	return resp.Availability.AvailabilityStatus == "available", nil
}

// Request/Response types

type CreateReservationRequest struct {
	BookingSessionID string  `json:"booking_session_id"`
	EventID          string  `json:"event_id"`
	SeatID           string  `json:"seat_id"`
	ZoneID           string  `json:"zone_id"`
	UserID           string  `json:"user_id"`
	PricingCategory  string  `json:"pricing_category"`
	BasePrice        float64 `json:"base_price"`
	FinalPrice       float64 `json:"final_price"`
	Currency         string  `json:"currency"`
	TimeoutMinutes   int     `json:"timeout_minutes"`
	CreatedBy        string  `json:"created_by,omitempty"`
}

type ConfirmReservationRequest struct {
	ReservationID string `json:"reservation_id"`
	ConfirmedBy   string `json:"confirmed_by"`
}

type ReleaseReservationRequest struct {
	ReservationID string `json:"reservation_id"`
	Reason        string `json:"reason"`
	ReleasedBy    string `json:"released_by"`
}

type ReleaseReservationsBySessionRequest struct {
	SessionID  string `json:"session_id"`
	Reason     string `json:"reason"`
	ReleasedBy string `json:"released_by"`
}

type ExtendReservationRequest struct {
	ReservationID    string `json:"reservation_id"`
	ExtensionMinutes int    `json:"extension_minutes"`
	ExtendedBy       string `json:"extended_by"`
}

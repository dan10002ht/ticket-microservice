package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ticket-service/metrics"
	"ticket-service/models"
	"ticket-service/services"

	ticketpb "shared-lib/protos/ticket"
)

// BookingController handles gRPC requests for booking operations
type BookingController struct {
	ticketpb.UnimplementedBookingServiceServer
	bookingService *services.BookingService
	logger         *zap.Logger
}

// NewBookingController creates a new booking controller
func NewBookingController(bookingService *services.BookingService, logger *zap.Logger) *BookingController {
	return &BookingController{
		bookingService: bookingService,
		logger:         logger,
	}
}

// CreateBookingSession creates a new booking session
func (c *BookingController) CreateBookingSession(ctx context.Context, req *ticketpb.CreateBookingSessionRequest) (*ticketpb.CreateBookingSessionResponse, error) {
	c.logger.Info("CreateBookingSession request received",
		zap.String("user_id", req.UserId),
		zap.String("event_id", req.EventId),
		zap.Int32("timeout_minutes", req.TimeoutMinutes),
	)

	serviceReq := &services.CreateBookingSessionRequest{
		UserID:         req.UserId,
		EventID:        req.EventId,
		Currency:       req.Currency,
		TimeoutMinutes: int(req.TimeoutMinutes),
		IPAddress:      req.IpAddress,
		UserAgent:      req.UserAgent,
		CreatedBy:      req.CreatedBy,
	}

	session, err := c.bookingService.CreateBookingSession(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to create booking session",
			zap.String("user_id", req.UserId),
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "CreateBookingSession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to create booking session: %v", err)
	}

	response := &ticketpb.CreateBookingSessionResponse{
		Success: true,
		Session: c.convertBookingSessionToProto(session),
	}

	c.logger.Info("Booking session created successfully",
		zap.String("session_id", session.ID),
		zap.String("event_id", req.EventId),
	)

	return response, nil
}

// GetBookingSession retrieves a booking session by ID
func (c *BookingController) GetBookingSession(ctx context.Context, req *ticketpb.GetBookingSessionRequest) (*ticketpb.GetBookingSessionResponse, error) {
	c.logger.Info("GetBookingSession request received",
		zap.String("session_id", req.SessionId),
	)

	session, err := c.bookingService.GetBookingSession(ctx, req.SessionId)
	if err != nil {
		c.logger.Error("Failed to get booking session",
			zap.String("session_id", req.SessionId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "GetBookingSession", "not_found")
		return nil, status.Errorf(codes.NotFound, "booking session not found: %v", err)
	}

	response := &ticketpb.GetBookingSessionResponse{
		Success: true,
		Session: c.convertBookingSessionToProto(session),
	}

	return response, nil
}

// GetBookingSessionByToken retrieves a booking session by token
func (c *BookingController) GetBookingSessionByToken(ctx context.Context, req *ticketpb.GetBookingSessionByTokenRequest) (*ticketpb.GetBookingSessionByTokenResponse, error) {
	c.logger.Info("GetBookingSessionByToken request received",
		zap.String("session_token", req.SessionToken),
	)

	session, err := c.bookingService.GetBookingSessionByToken(ctx, req.SessionToken)
	if err != nil {
		c.logger.Error("Failed to get booking session by token",
			zap.String("session_token", req.SessionToken),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "GetBookingSessionByToken", "not_found")
		return nil, status.Errorf(codes.NotFound, "booking session not found: %v", err)
	}

	response := &ticketpb.GetBookingSessionByTokenResponse{
		Success: true,
		Session: c.convertBookingSessionToProto(session),
	}

	return response, nil
}

// AddSeatToSession adds a seat to the booking session
func (c *BookingController) AddSeatToSession(ctx context.Context, req *ticketpb.AddSeatToSessionRequest) (*ticketpb.AddSeatToSessionResponse, error) {
	c.logger.Info("AddSeatToSession request received",
		zap.String("session_id", req.SessionId),
		zap.String("seat_id", req.SeatId),
	)

	serviceReq := &services.AddSeatToSessionRequest{
		SessionID:       req.SessionId,
		EventID:         req.EventId,
		SeatID:          req.SeatId,
		ZoneID:          req.ZoneId,
		PricingCategory: req.PricingCategory,
		BasePrice:       req.BasePrice,
		FinalPrice:      req.FinalPrice,
		Currency:        req.Currency,
		CreatedBy:       req.CreatedBy,
	}

	err := c.bookingService.AddSeatToSession(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to add seat to session",
			zap.String("session_id", req.SessionId),
			zap.String("seat_id", req.SeatId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "AddSeatToSession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to add seat to session: %v", err)
	}

	response := &ticketpb.AddSeatToSessionResponse{
		Success: true,
		Message: "Seat added to session successfully",
	}

	return response, nil
}

// RemoveSeatFromSession removes a seat from the booking session
func (c *BookingController) RemoveSeatFromSession(ctx context.Context, req *ticketpb.RemoveSeatFromSessionRequest) (*ticketpb.RemoveSeatFromSessionResponse, error) {
	c.logger.Info("RemoveSeatFromSession request received",
		zap.String("session_id", req.SessionId),
		zap.String("seat_id", req.SeatId),
	)

	serviceReq := &services.RemoveSeatFromSessionRequest{
		SessionID: req.SessionId,
		SeatID:    req.SeatId,
		Reason:    req.Reason,
		RemovedBy: req.RemovedBy,
	}

	err := c.bookingService.RemoveSeatFromSession(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to remove seat from session",
			zap.String("session_id", req.SessionId),
			zap.String("seat_id", req.SeatId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "RemoveSeatFromSession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to remove seat from session: %v", err)
	}

	response := &ticketpb.RemoveSeatFromSessionResponse{
		Success: true,
		Message: "Seat removed from session successfully",
	}

	return response, nil
}

// CompleteBookingSession completes a booking session
func (c *BookingController) CompleteBookingSession(ctx context.Context, req *ticketpb.CompleteBookingSessionRequest) (*ticketpb.CompleteBookingSessionResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		metrics.ObserveBookingDuration("", "complete_session", duration)
	}()

	c.logger.Info("CompleteBookingSession request received",
		zap.String("session_id", req.SessionId),
		zap.String("payment_method", req.PaymentMethod),
	)

	serviceReq := &services.CompleteBookingSessionRequest{
		SessionID:     req.SessionId,
		PaymentMethod: req.PaymentMethod,
		CompletedBy:   req.CompletedBy,
	}

	result, err := c.bookingService.CompleteBookingSession(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to complete booking session",
			zap.String("session_id", req.SessionId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "CompleteBookingSession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to complete booking session: %v", err)
	}

	response := &ticketpb.CompleteBookingSessionResponse{
		Success:     result.Success,
		PaymentId:   result.PaymentID,
		SeatCount:   int32(result.SeatCount),
		TotalAmount: result.TotalAmount,
		Message:     "Booking session completed successfully",
	}

	return response, nil
}

// CancelBookingSession cancels a booking session
func (c *BookingController) CancelBookingSession(ctx context.Context, req *ticketpb.CancelBookingSessionRequest) (*ticketpb.CancelBookingSessionResponse, error) {
	c.logger.Info("CancelBookingSession request received",
		zap.String("session_id", req.SessionId),
		zap.String("reason", req.Reason),
	)

	serviceReq := &services.CancelBookingSessionRequest{
		SessionID:   req.SessionId,
		Reason:      req.Reason,
		CancelledBy: req.CancelledBy,
	}

	err := c.bookingService.CancelBookingSession(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to cancel booking session",
			zap.String("session_id", req.SessionId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "CancelBookingSession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to cancel booking session: %v", err)
	}

	response := &ticketpb.CancelBookingSessionResponse{
		Success: true,
		Message: "Booking session cancelled successfully",
	}

	return response, nil
}

// GetSessionReservations gets seat reservations for a session
func (c *BookingController) GetSessionReservations(ctx context.Context, req *ticketpb.GetSessionReservationsRequest) (*ticketpb.GetSessionReservationsResponse, error) {
	c.logger.Info("GetSessionReservations request received",
		zap.String("session_id", req.SessionId),
	)

	reservations, err := c.bookingService.GetSessionReservations(ctx, req.SessionId)
	if err != nil {
		c.logger.Error("Failed to get session reservations",
			zap.String("session_id", req.SessionId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "GetSessionReservations", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get session reservations: %v", err)
	}

	// Convert reservations to proto
	protoReservations := make([]*ticketpb.SeatReservation, len(reservations))
	for i, reservation := range reservations {
		protoReservations[i] = c.convertSeatReservationToProto(reservation)
	}

	response := &ticketpb.GetSessionReservationsResponse{
		Success:      true,
		Reservations: protoReservations,
		Count:        int32(len(reservations)),
	}

	return response, nil
}

// CleanupExpiredSessions cleans up expired booking sessions
func (c *BookingController) CleanupExpiredSessions(ctx context.Context, req *ticketpb.CleanupExpiredSessionsRequest) (*ticketpb.CleanupExpiredSessionsResponse, error) {
	c.logger.Info("CleanupExpiredSessions request received")

	err := c.bookingService.CleanupExpiredSessions(ctx)
	if err != nil {
		c.logger.Error("Failed to cleanup expired sessions",
			zap.Error(err),
		)
		metrics.IncrementGRPCError("booking", "CleanupExpiredSessions", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to cleanup expired sessions: %v", err)
	}

	response := &ticketpb.CleanupExpiredSessionsResponse{
		Success: true,
		Message: "Expired sessions cleaned up successfully",
	}

	return response, nil
}

// Helper methods

func (c *BookingController) convertBookingSessionToProto(session *models.BookingSession) *ticketpb.BookingSession {
	protoSession := &ticketpb.BookingSession{
		Id:           session.ID,
		UserId:       session.UserID,
		EventId:      session.EventID,
		SessionToken: session.SessionToken,
		Status:       session.Status,
		SeatCount:    int32(session.SeatCount),
		TotalAmount:  session.TotalAmount,
		Currency:     session.Currency,
		ExpiresAt:    session.ExpiresAt.Unix(),
		CreatedAt:    session.CreatedAt.Unix(),
		UpdatedAt:    session.UpdatedAt.Unix(),
		// Computed fields
		RemainingTime: int64(session.GetRemainingTime().Seconds()),
		AveragePrice:  session.CalculateAveragePrice(),
		IsActive:      session.IsActive(),
		IsExpired:     session.IsExpired(),
	}

	// Set optional fields
	if session.CompletedAt != nil {
		protoSession.CompletedAt = session.CompletedAt.Unix()
	}
	if session.CancelledAt != nil {
		protoSession.CancelledAt = session.CancelledAt.Unix()
	}
	if session.CancelledReason != nil {
		protoSession.CancelledReason = *session.CancelledReason
	}
	if session.IPAddress != nil {
		protoSession.IpAddress = *session.IPAddress
	}
	if session.UserAgent != nil {
		protoSession.UserAgent = *session.UserAgent
	}
	if session.Metadata != nil {
		protoSession.Metadata = *session.Metadata
	}

	return protoSession
}

func (c *BookingController) convertSeatReservationToProto(reservation *models.SeatReservation) *ticketpb.SeatReservation {
	protoReservation := &ticketpb.SeatReservation{
		Id:               reservation.ID,
		BookingSessionId: reservation.BookingSessionID,
		EventId:          reservation.EventID,
		SeatId:           reservation.SeatID,
		ZoneId:           reservation.ZoneID,
		ReservationToken: reservation.ReservationToken,
		Status:           reservation.Status,
		ReservedAt:       reservation.ReservedAt.Unix(),
		ExpiresAt:        reservation.ExpiresAt.Unix(),
		PricingCategory:  reservation.PricingCategory,
		BasePrice:        reservation.BasePrice,
		FinalPrice:       reservation.FinalPrice,
		Currency:         reservation.Currency,
		CreatedAt:        reservation.CreatedAt.Unix(),
		UpdatedAt:        reservation.UpdatedAt.Unix(),
		// Computed fields
		RemainingTime:      int64(reservation.GetRemainingTime().Seconds()),
		DiscountAmount:     reservation.CalculateDiscount(),
		DiscountPercentage: reservation.GetDiscountPercentage(),
		IsReserved:         reservation.IsReserved(),
		IsExpired:          reservation.IsExpired(),
	}

	// Set optional fields
	if reservation.ReleasedAt != nil {
		protoReservation.ReleasedAt = reservation.ReleasedAt.Unix()
	}
	if reservation.ReleasedReason != nil {
		protoReservation.ReleasedReason = *reservation.ReleasedReason
	}
	if reservation.Metadata != nil {
		protoReservation.Metadata = *reservation.Metadata
	}

	return protoReservation
}

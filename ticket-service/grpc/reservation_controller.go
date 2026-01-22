package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ticket-service/metrics"
	"ticket-service/models"
	"ticket-service/services"

	ticketpb "ticket-service/internal/protos/ticket"
)

// ReservationController handles gRPC requests for reservation operations
type ReservationController struct {
	ticketpb.UnimplementedReservationServiceServer
	reservationService *services.ReservationService
	logger             *zap.Logger
}

// NewReservationController creates a new reservation controller
func NewReservationController(reservationService *services.ReservationService, logger *zap.Logger) *ReservationController {
	return &ReservationController{
		reservationService: reservationService,
		logger:             logger,
	}
}

// CreateReservation creates a new seat reservation
func (c *ReservationController) CreateReservation(ctx context.Context, req *ticketpb.CreateReservationRequest) (*ticketpb.CreateReservationResponse, error) {
	c.logger.Info("CreateReservation request received",
		zap.String("booking_session_id", req.BookingSessionId),
		zap.String("event_id", req.EventId),
		zap.String("seat_id", req.SeatId),
	)

	serviceReq := &services.CreateReservationRequest{
		BookingSessionID: req.BookingSessionId,
		EventID:          req.EventId,
		SeatID:           req.SeatId,
		ZoneID:           req.ZoneId,
		UserID:           req.UserId,
		PricingCategory:  req.PricingCategory,
		BasePrice:        req.BasePrice,
		FinalPrice:       req.FinalPrice,
		Currency:         req.Currency,
		TimeoutMinutes:   int(req.TimeoutMinutes),
		CreatedBy:        req.CreatedBy,
	}

	reservation, err := c.reservationService.CreateReservation(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to create reservation",
			zap.String("booking_session_id", req.BookingSessionId),
			zap.String("seat_id", req.SeatId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "CreateReservation", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to create reservation: %v", err)
	}

	response := &ticketpb.CreateReservationResponse{
		Success:     true,
		Reservation: c.convertSeatReservationToProto(reservation),
	}

	c.logger.Info("Reservation created successfully",
		zap.String("reservation_id", reservation.ID),
		zap.String("seat_id", req.SeatId),
	)

	return response, nil
}

// GetReservation retrieves a seat reservation by ID
func (c *ReservationController) GetReservation(ctx context.Context, req *ticketpb.GetReservationRequest) (*ticketpb.GetReservationResponse, error) {
	c.logger.Info("GetReservation request received",
		zap.String("reservation_id", req.ReservationId),
	)

	reservation, err := c.reservationService.GetReservation(ctx, req.ReservationId)
	if err != nil {
		c.logger.Error("Failed to get reservation",
			zap.String("reservation_id", req.ReservationId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetReservation", "not_found")
		return nil, status.Errorf(codes.NotFound, "reservation not found: %v", err)
	}

	response := &ticketpb.GetReservationResponse{
		Success:     true,
		Reservation: c.convertSeatReservationToProto(reservation),
	}

	return response, nil
}

// GetReservationByToken retrieves a seat reservation by token
func (c *ReservationController) GetReservationByToken(ctx context.Context, req *ticketpb.GetReservationByTokenRequest) (*ticketpb.GetReservationByTokenResponse, error) {
	c.logger.Info("GetReservationByToken request received",
		zap.String("reservation_token", req.ReservationToken),
	)

	reservation, err := c.reservationService.GetReservationByToken(ctx, req.ReservationToken)
	if err != nil {
		c.logger.Error("Failed to get reservation by token",
			zap.String("reservation_token", req.ReservationToken),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetReservationByToken", "not_found")
		return nil, status.Errorf(codes.NotFound, "reservation not found: %v", err)
	}

	response := &ticketpb.GetReservationByTokenResponse{
		Success:     true,
		Reservation: c.convertSeatReservationToProto(reservation),
	}

	return response, nil
}

// GetReservationsBySession gets seat reservations for a booking session
func (c *ReservationController) GetReservationsBySession(ctx context.Context, req *ticketpb.GetReservationsBySessionRequest) (*ticketpb.GetReservationsBySessionResponse, error) {
	c.logger.Info("GetReservationsBySession request received",
		zap.String("session_id", req.SessionId),
	)

	reservations, err := c.reservationService.GetReservationsBySession(ctx, req.SessionId)
	if err != nil {
		c.logger.Error("Failed to get session reservations",
			zap.String("session_id", req.SessionId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetReservationsBySession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get session reservations: %v", err)
	}

	// Convert reservations to proto
	protoReservations := make([]*ticketpb.SeatReservation, len(reservations))
	for i, reservation := range reservations {
		protoReservations[i] = c.convertSeatReservationToProto(reservation)
	}

	response := &ticketpb.GetReservationsBySessionResponse{
		Success:      true,
		Reservations: protoReservations,
		Count:        int32(len(reservations)),
	}

	return response, nil
}

// GetReservationsByEvent gets seat reservations for an event
func (c *ReservationController) GetReservationsByEvent(ctx context.Context, req *ticketpb.GetReservationsByEventRequest) (*ticketpb.GetReservationsByEventResponse, error) {
	c.logger.Info("GetReservationsByEvent request received",
		zap.String("event_id", req.EventId),
		zap.Int32("page", req.Page),
		zap.Int32("limit", req.Limit),
	)

	reservations, total, err := c.reservationService.GetReservationsByEvent(ctx, req.EventId, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get event reservations",
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetReservationsByEvent", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get event reservations: %v", err)
	}

	// Convert reservations to proto
	protoReservations := make([]*ticketpb.SeatReservation, len(reservations))
	for i, reservation := range reservations {
		protoReservations[i] = c.convertSeatReservationToProto(reservation)
	}

	response := &ticketpb.GetReservationsByEventResponse{
		Success:      true,
		Reservations: protoReservations,
		Total:        int32(total),
		Page:         req.Page,
		Limit:        req.Limit,
		HasMore:      int(req.Page)*int(req.Limit) < total,
	}

	return response, nil
}

// GetReservationsBySeat gets seat reservations for a specific seat
func (c *ReservationController) GetReservationsBySeat(ctx context.Context, req *ticketpb.GetReservationsBySeatRequest) (*ticketpb.GetReservationsBySeatResponse, error) {
	c.logger.Info("GetReservationsBySeat request received",
		zap.String("seat_id", req.SeatId),
	)

	reservations, err := c.reservationService.GetReservationsBySeat(ctx, req.SeatId)
	if err != nil {
		c.logger.Error("Failed to get seat reservations",
			zap.String("seat_id", req.SeatId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetReservationsBySeat", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get seat reservations: %v", err)
	}

	// Convert reservations to proto
	protoReservations := make([]*ticketpb.SeatReservation, len(reservations))
	for i, reservation := range reservations {
		protoReservations[i] = c.convertSeatReservationToProto(reservation)
	}

	response := &ticketpb.GetReservationsBySeatResponse{
		Success:      true,
		Reservations: protoReservations,
		Count:        int32(len(reservations)),
	}

	return response, nil
}

// GetActiveReservations gets active seat reservations for an event
func (c *ReservationController) GetActiveReservations(ctx context.Context, req *ticketpb.GetActiveReservationsRequest) (*ticketpb.GetActiveReservationsResponse, error) {
	c.logger.Info("GetActiveReservations request received",
		zap.String("event_id", req.EventId),
	)

	reservations, err := c.reservationService.GetActiveReservations(ctx, req.EventId)
	if err != nil {
		c.logger.Error("Failed to get active reservations",
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetActiveReservations", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get active reservations: %v", err)
	}

	// Convert reservations to proto
	protoReservations := make([]*ticketpb.SeatReservation, len(reservations))
	for i, reservation := range reservations {
		protoReservations[i] = c.convertSeatReservationToProto(reservation)
	}

	response := &ticketpb.GetActiveReservationsResponse{
		Success:      true,
		Reservations: protoReservations,
		Count:        int32(len(reservations)),
	}

	return response, nil
}

// ConfirmReservation confirms a seat reservation
func (c *ReservationController) ConfirmReservation(ctx context.Context, req *ticketpb.ConfirmReservationRequest) (*ticketpb.ConfirmReservationResponse, error) {
	c.logger.Info("ConfirmReservation request received",
		zap.String("reservation_id", req.ReservationId),
		zap.String("confirmed_by", req.ConfirmedBy),
	)

	serviceReq := &services.ConfirmReservationRequest{
		ReservationID: req.ReservationId,
		ConfirmedBy:   req.ConfirmedBy,
	}

	err := c.reservationService.ConfirmReservation(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to confirm reservation",
			zap.String("reservation_id", req.ReservationId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "ConfirmReservation", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to confirm reservation: %v", err)
	}

	response := &ticketpb.ConfirmReservationResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
	}

	return response, nil
}

// ReleaseReservation releases a seat reservation
func (c *ReservationController) ReleaseReservation(ctx context.Context, req *ticketpb.ReleaseReservationRequest) (*ticketpb.ReleaseReservationResponse, error) {
	c.logger.Info("ReleaseReservation request received",
		zap.String("reservation_id", req.ReservationId),
		zap.String("reason", req.Reason),
	)

	serviceReq := &services.ReleaseReservationRequest{
		ReservationID: req.ReservationId,
		Reason:        req.Reason,
		ReleasedBy:    req.ReleasedBy,
	}

	err := c.reservationService.ReleaseReservation(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to release reservation",
			zap.String("reservation_id", req.ReservationId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "ReleaseReservation", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to release reservation: %v", err)
	}

	response := &ticketpb.ReleaseReservationResponse{
		Success: true,
		Message: "Reservation released successfully",
	}

	return response, nil
}

// ReleaseReservationsBySession releases all reservations for a booking session
func (c *ReservationController) ReleaseReservationsBySession(ctx context.Context, req *ticketpb.ReleaseReservationsBySessionRequest) (*ticketpb.ReleaseReservationsBySessionResponse, error) {
	c.logger.Info("ReleaseReservationsBySession request received",
		zap.String("session_id", req.SessionId),
		zap.String("reason", req.Reason),
	)

	serviceReq := &services.ReleaseReservationsBySessionRequest{
		SessionID:  req.SessionId,
		Reason:     req.Reason,
		ReleasedBy: req.ReleasedBy,
	}

	err := c.reservationService.ReleaseReservationsBySession(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to release session reservations",
			zap.String("session_id", req.SessionId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "ReleaseReservationsBySession", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to release session reservations: %v", err)
	}

	response := &ticketpb.ReleaseReservationsBySessionResponse{
		Success: true,
		Message: "Session reservations released successfully",
	}

	return response, nil
}

// CheckSeatAvailability checks if a seat is available for reservation
func (c *ReservationController) CheckSeatAvailability(ctx context.Context, req *ticketpb.CheckSeatAvailabilityRequest) (*ticketpb.CheckSeatAvailabilityResponse, error) {
	c.logger.Info("CheckSeatAvailability request received",
		zap.String("seat_id", req.SeatId),
	)

	available, err := c.reservationService.CheckSeatAvailability(ctx, req.SeatId)
	if err != nil {
		c.logger.Error("Failed to check seat availability",
			zap.String("seat_id", req.SeatId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "CheckSeatAvailability", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to check seat availability: %v", err)
	}

	response := &ticketpb.CheckSeatAvailabilityResponse{
		Success:   true,
		Available: available,
	}

	return response, nil
}

// GetReservationStats gets reservation statistics for an event
func (c *ReservationController) GetReservationStats(ctx context.Context, req *ticketpb.GetReservationStatsRequest) (*ticketpb.GetReservationStatsResponse, error) {
	c.logger.Info("GetReservationStats request received",
		zap.String("event_id", req.EventId),
	)

	stats, err := c.reservationService.GetReservationStats(ctx, req.EventId)
	if err != nil {
		c.logger.Error("Failed to get reservation stats",
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "GetReservationStats", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get reservation stats: %v", err)
	}

	// Convert map[string]int to map[string]int32
	stats32 := make(map[string]int32, len(stats))
	for k, v := range stats {
		stats32[k] = int32(v)
	}

	response := &ticketpb.GetReservationStatsResponse{
		Success: true,
		Stats:   stats32,
	}

	return response, nil
}

// ExtendReservation extends the expiration time of a reservation
func (c *ReservationController) ExtendReservation(ctx context.Context, req *ticketpb.ExtendReservationRequest) (*ticketpb.ExtendReservationResponse, error) {
	c.logger.Info("ExtendReservation request received",
		zap.String("reservation_id", req.ReservationId),
		zap.Int32("extension_minutes", req.ExtensionMinutes),
	)

	serviceReq := &services.ExtendReservationRequest{
		ReservationID:    req.ReservationId,
		ExtensionMinutes: int(req.ExtensionMinutes),
		ExtendedBy:       req.ExtendedBy,
	}

	err := c.reservationService.ExtendReservation(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to extend reservation",
			zap.String("reservation_id", req.ReservationId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "ExtendReservation", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to extend reservation: %v", err)
	}

	response := &ticketpb.ExtendReservationResponse{
		Success: true,
		Message: "Reservation extended successfully",
	}

	return response, nil
}

// CleanupExpiredReservations cleans up expired seat reservations
func (c *ReservationController) CleanupExpiredReservations(ctx context.Context, req *ticketpb.CleanupExpiredReservationsRequest) (*ticketpb.CleanupExpiredReservationsResponse, error) {
	c.logger.Info("CleanupExpiredReservations request received")

	err := c.reservationService.CleanupExpiredReservations(ctx)
	if err != nil {
		c.logger.Error("Failed to cleanup expired reservations",
			zap.Error(err),
		)
		metrics.IncrementGRPCError("reservation", "CleanupExpiredReservations", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to cleanup expired reservations: %v", err)
	}

	response := &ticketpb.CleanupExpiredReservationsResponse{
		Success: true,
		Message: "Expired reservations cleaned up successfully",
	}

	return response, nil
}

// Helper methods

func (c *ReservationController) convertSeatReservationToProto(reservation *models.SeatReservation) *ticketpb.SeatReservation {
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

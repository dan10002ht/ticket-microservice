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

// TicketController handles gRPC requests for ticket operations
type TicketController struct {
	ticketpb.UnimplementedTicketServiceServer
	ticketService *services.TicketService
	logger        *zap.Logger
}

// NewTicketController creates a new ticket controller
func NewTicketController(ticketService *services.TicketService, logger *zap.Logger) *TicketController {
	return &TicketController{
		ticketService: ticketService,
		logger:        logger,
	}
}

// CreateTicket creates a new ticket
func (c *TicketController) CreateTicket(ctx context.Context, req *ticketpb.CreateTicketRequest) (*ticketpb.CreateTicketResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		metrics.ObserveBookingDuration(req.EventId, "create_ticket", duration)
	}()

	c.logger.Info("CreateTicket request received",
		zap.String("event_id", req.EventId),
		zap.String("user_id", req.UserId),
	)

	// Convert request to service request
	serviceReq := &services.CreateTicketRequest{
		EventID:          req.EventId,
		SeatID:           req.SeatId,
		ZoneID:           req.ZoneId,
		UserID:           req.UserId,
		BookingSessionID: req.BookingSessionId,
		TicketType:       req.TicketType,
		PricingCategory:  req.PricingCategory,
		BasePrice:        req.BasePrice,
		FinalPrice:       req.FinalPrice,
		Currency:         req.Currency,
		DiscountReason:   req.DiscountReason,
		CreatedBy:        req.CreatedBy,
	}

	// Set valid until if provided
	if req.ValidUntil != nil {
		validUntil := req.ValidUntil.AsTime()
		serviceReq.ValidUntil = &validUntil
	}

	// Create ticket
	ticket, err := c.ticketService.CreateTicket(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to create ticket",
			zap.String("event_id", req.EventId),
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "CreateTicket", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to create ticket: %v", err)
	}

	// Convert ticket to response
	response := &ticketpb.CreateTicketResponse{
		Success: true,
		Ticket:  c.convertTicketToProto(ticket),
	}

	c.logger.Info("Ticket created successfully",
		zap.String("ticket_id", ticket.ID),
		zap.String("event_id", req.EventId),
	)

	return response, nil
}

// GetTicket retrieves a ticket by ID
func (c *TicketController) GetTicket(ctx context.Context, req *ticketpb.GetTicketRequest) (*ticketpb.GetTicketResponse, error) {
	c.logger.Info("GetTicket request received",
		zap.String("ticket_id", req.TicketId),
	)

	ticket, err := c.ticketService.GetTicket(ctx, req.TicketId)
	if err != nil {
		c.logger.Error("Failed to get ticket",
			zap.String("ticket_id", req.TicketId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "GetTicket", "not_found")
		return nil, status.Errorf(codes.NotFound, "ticket not found: %v", err)
	}

	response := &ticketpb.GetTicketResponse{
		Success: true,
		Ticket:  c.convertTicketToProto(ticket),
	}

	return response, nil
}

// GetTicketByNumber retrieves a ticket by ticket number
func (c *TicketController) GetTicketByNumber(ctx context.Context, req *ticketpb.GetTicketByNumberRequest) (*ticketpb.GetTicketByNumberResponse, error) {
	c.logger.Info("GetTicketByNumber request received",
		zap.String("ticket_number", req.TicketNumber),
	)

	ticket, err := c.ticketService.GetTicketByNumber(ctx, req.TicketNumber)
	if err != nil {
		c.logger.Error("Failed to get ticket by number",
			zap.String("ticket_number", req.TicketNumber),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "GetTicketByNumber", "not_found")
		return nil, status.Errorf(codes.NotFound, "ticket not found: %v", err)
	}

	response := &ticketpb.GetTicketByNumberResponse{
		Success: true,
		Ticket:  c.convertTicketToProto(ticket),
	}

	return response, nil
}

// GetUserTickets retrieves tickets for a user
func (c *TicketController) GetUserTickets(ctx context.Context, req *ticketpb.GetUserTicketsRequest) (*ticketpb.GetUserTicketsResponse, error) {
	c.logger.Info("GetUserTickets request received",
		zap.String("user_id", req.UserId),
		zap.Int32("page", req.Page),
		zap.Int32("limit", req.Limit),
	)

	tickets, total, err := c.ticketService.GetUserTickets(ctx, req.UserId, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get user tickets",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "GetUserTickets", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get user tickets: %v", err)
	}

	// Convert tickets to proto
	protoTickets := make([]*ticketpb.Ticket, len(tickets))
	for i, ticket := range tickets {
		protoTickets[i] = c.convertTicketToProto(ticket)
	}

	response := &ticketpb.GetUserTicketsResponse{
		Success: true,
		Tickets: protoTickets,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
		HasMore: int(req.Page)*int(req.Limit) < total,
	}

	return response, nil
}

// GetEventTickets retrieves tickets for an event
func (c *TicketController) GetEventTickets(ctx context.Context, req *ticketpb.GetEventTicketsRequest) (*ticketpb.GetEventTicketsResponse, error) {
	c.logger.Info("GetEventTickets request received",
		zap.String("event_id", req.EventId),
		zap.Int32("page", req.Page),
		zap.Int32("limit", req.Limit),
	)

	tickets, total, err := c.ticketService.GetEventTickets(ctx, req.EventId, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get event tickets",
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "GetEventTickets", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to get event tickets: %v", err)
	}

	// Convert tickets to proto
	protoTickets := make([]*ticketpb.Ticket, len(tickets))
	for i, ticket := range tickets {
		protoTickets[i] = c.convertTicketToProto(ticket)
	}

	response := &ticketpb.GetEventTicketsResponse{
		Success: true,
		Tickets: protoTickets,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
		HasMore: int(req.Page)*int(req.Limit) < total,
	}

	return response, nil
}

// UpdateTicketStatus updates ticket status
func (c *TicketController) UpdateTicketStatus(ctx context.Context, req *ticketpb.UpdateTicketStatusRequest) (*ticketpb.UpdateTicketStatusResponse, error) {
	c.logger.Info("UpdateTicketStatus request received",
		zap.String("ticket_id", req.TicketId),
		zap.String("status", req.Status),
		zap.String("updated_by", req.UpdatedBy),
	)

	err := c.ticketService.UpdateTicketStatus(ctx, req.TicketId, req.Status, req.UpdatedBy)
	if err != nil {
		c.logger.Error("Failed to update ticket status",
			zap.String("ticket_id", req.TicketId),
			zap.String("status", req.Status),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "UpdateTicketStatus", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to update ticket status: %v", err)
	}

	response := &ticketpb.UpdateTicketStatusResponse{
		Success: true,
		Message: "Ticket status updated successfully",
	}

	return response, nil
}

// ProcessPayment processes payment for a ticket
func (c *TicketController) ProcessPayment(ctx context.Context, req *ticketpb.ProcessPaymentRequest) (*ticketpb.ProcessPaymentResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		metrics.ObservePaymentDuration("", req.PaymentMethod, duration)
	}()

	c.logger.Info("ProcessPayment request received",
		zap.String("ticket_id", req.TicketId),
		zap.String("payment_method", req.PaymentMethod),
	)

	serviceReq := &services.ProcessPaymentRequest{
		TicketID:      req.TicketId,
		PaymentMethod: req.PaymentMethod,
		UpdatedBy:     req.UpdatedBy,
	}

	result, err := c.ticketService.ProcessPayment(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to process payment",
			zap.String("ticket_id", req.TicketId),
			zap.Error(err),
		)
		metrics.IncrementPaymentFailed("", "payment_error")
		metrics.IncrementGRPCError("ticket", "ProcessPayment", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	response := &ticketpb.ProcessPaymentResponse{
		Success:       result.Success,
		PaymentId:     result.PaymentID,
		TicketStatus:  result.TicketStatus,
		PaymentStatus: result.PaymentStatus,
	}

	return response, nil
}

// CancelTicket cancels a ticket
func (c *TicketController) CancelTicket(ctx context.Context, req *ticketpb.CancelTicketRequest) (*ticketpb.CancelTicketResponse, error) {
	c.logger.Info("CancelTicket request received",
		zap.String("ticket_id", req.TicketId),
		zap.String("reason", req.Reason),
		zap.String("cancelled_by", req.CancelledBy),
	)

	serviceReq := &services.CancelTicketRequest{
		TicketID:    req.TicketId,
		Reason:      req.Reason,
		CancelledBy: req.CancelledBy,
	}

	err := c.ticketService.CancelTicket(ctx, serviceReq)
	if err != nil {
		c.logger.Error("Failed to cancel ticket",
			zap.String("ticket_id", req.TicketId),
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "CancelTicket", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to cancel ticket: %v", err)
	}

	response := &ticketpb.CancelTicketResponse{
		Success: true,
		Message: "Ticket cancelled successfully",
	}

	return response, nil
}

// SearchTickets searches tickets with filters
func (c *TicketController) SearchTickets(ctx context.Context, req *ticketpb.SearchTicketsRequest) (*ticketpb.SearchTicketsResponse, error) {
	c.logger.Info("SearchTickets request received",
		zap.Any("filters", req.Filters),
		zap.Int32("page", req.Page),
		zap.Int32("limit", req.Limit),
	)

	// Convert filters
	filters := make(map[string]interface{})
	for key, value := range req.Filters {
		filters[key] = value
	}

	tickets, total, err := c.ticketService.SearchTickets(ctx, filters, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to search tickets",
			zap.Error(err),
		)
		metrics.IncrementGRPCError("ticket", "SearchTickets", "service_error")
		return nil, status.Errorf(codes.Internal, "failed to search tickets: %v", err)
	}

	// Convert tickets to proto
	protoTickets := make([]*ticketpb.Ticket, len(tickets))
	for i, ticket := range tickets {
		protoTickets[i] = c.convertTicketToProto(ticket)
	}

	response := &ticketpb.SearchTicketsResponse{
		Success: true,
		Tickets: protoTickets,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
		HasMore: int(req.Page)*int(req.Limit) < total,
	}

	return response, nil
}

// Helper methods

func (c *TicketController) convertTicketToProto(ticket *models.Ticket) *ticketpb.Ticket {
	protoTicket := &ticketpb.Ticket{
		Id:              ticket.ID,
		EventId:         ticket.EventID,
		SeatId:          ticket.SeatID,
		ZoneId:          ticket.ZoneID,
		UserId:          ticket.UserID,
		TicketNumber:    ticket.TicketNumber,
		TicketType:      ticket.TicketType,
		PricingCategory: ticket.PricingCategory,
		BasePrice:       ticket.BasePrice,
		FinalPrice:      ticket.FinalPrice,
		Currency:        ticket.Currency,
		DiscountAmount:  ticket.DiscountAmount,
		Status:          ticket.Status,
		PaymentStatus:   ticket.PaymentStatus,
		ValidFrom:       ticket.ValidFrom.Unix(),
		CreatedAt:       ticket.CreatedAt.Unix(),
		UpdatedAt:       ticket.UpdatedAt.Unix(),
	}

	// Set optional fields
	if ticket.BookingSessionID != nil {
		protoTicket.BookingSessionId = *ticket.BookingSessionID
	}
	if ticket.DiscountReason != nil {
		protoTicket.DiscountReason = *ticket.DiscountReason
	}
	if ticket.PaymentMethod != nil {
		protoTicket.PaymentMethod = *ticket.PaymentMethod
	}
	if ticket.PaymentReference != nil {
		protoTicket.PaymentReference = *ticket.PaymentReference
	}
	if ticket.QRCode != nil {
		protoTicket.QrCode = *ticket.QRCode
	}
	if ticket.Barcode != nil {
		protoTicket.Barcode = *ticket.Barcode
	}
	if ticket.ValidUntil != nil {
		protoTicket.ValidUntil = ticket.ValidUntil.Unix()
	}
	if ticket.UsedAt != nil {
		protoTicket.UsedAt = ticket.UsedAt.Unix()
	}
	if ticket.CancelledAt != nil {
		protoTicket.CancelledAt = ticket.CancelledAt.Unix()
	}
	if ticket.CancelledReason != nil {
		protoTicket.CancelledReason = *ticket.CancelledReason
	}
	if ticket.RefundedAt != nil {
		protoTicket.RefundedAt = ticket.RefundedAt.Unix()
	}
	if ticket.RefundedAmount != nil {
		protoTicket.RefundedAmount = *ticket.RefundedAmount
	}
	if ticket.Metadata != nil {
		protoTicket.Metadata = *ticket.Metadata
	}

	return protoTicket
}

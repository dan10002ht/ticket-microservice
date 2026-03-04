package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"checkin-service/config"
	checkinpb "checkin-service/internal/protos/checkin"
	"checkin-service/models"
	"checkin-service/services"
)

// Server wraps the gRPC server and registers all service implementations.
type Server struct {
	services *services.Services
	config   *config.Config
	logger   *zap.Logger
	server   *grpc.Server
}

func NewServer(svc *services.Services, cfg *config.Config, logger *zap.Logger) *Server {
	return &Server{services: svc, config: cfg, logger: logger}
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen on :%d: %w", port, err)
	}

	s.server = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 15 * time.Second,
			Time:              5 * time.Second,
			Timeout:           1 * time.Second,
		}),
	)

	handler := &checkinHandler{svc: s.services, logger: s.logger}
	checkinpb.RegisterCheckinServiceServer(s.server, handler)
	reflection.Register(s.server)

	s.logger.Info("gRPC server listening", zap.Int("port", port))
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// checkinHandler implements CheckinServiceServer, delegating to the service layer.
type checkinHandler struct {
	checkinpb.UnimplementedCheckinServiceServer
	svc    *services.Services
	logger *zap.Logger
}

// CheckIn processes a ticket check-in request.
func (h *checkinHandler) CheckIn(ctx context.Context, req *checkinpb.CheckInRequest) (*checkinpb.CheckInResponse, error) {
	var staffID *string
	if req.StaffId != "" {
		s := req.StaffId
		staffID = &s
	}
	var deviceID *string
	if req.DeviceId != "" {
		d := req.DeviceId
		deviceID = &d
	}
	var gate *string
	if req.Gate != "" {
		g := req.Gate
		gate = &g
	}

	m := &models.CheckIn{
		TicketID: req.TicketId,
		EventID:  req.EventId,
		QRCode:   req.QrCode,
		StaffID:  staffID,
		DeviceID: deviceID,
		Gate:     gate,
	}

	record, err := h.svc.Checkin.CheckIn(ctx, m)
	if err != nil {
		var ce *services.CheckInError
		if errors.As(err, &ce) {
			return &checkinpb.CheckInResponse{
				Success:   false,
				ErrorCode: domainErrToProto(ce),
				Message:   ce.Message,
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "check-in failed: %v", err)
	}

	return &checkinpb.CheckInResponse{
		Success:   true,
		Record:    modelToProto(record),
		ErrorCode: checkinpb.CheckInError_CHECKIN_SUCCESS,
	}, nil
}

// GetCheckIn retrieves a single check-in record by ID.
func (h *checkinHandler) GetCheckIn(ctx context.Context, req *checkinpb.GetCheckInRequest) (*checkinpb.GetCheckInResponse, error) {
	record, err := h.svc.Checkin.GetCheckIn(ctx, req.CheckinId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "check-in not found: %v", err)
	}
	return &checkinpb.GetCheckInResponse{
		Success: true,
		Record:  modelToProto(record),
	}, nil
}

// ListCheckIns returns paginated check-in records for an event.
func (h *checkinHandler) ListCheckIns(ctx context.Context, req *checkinpb.ListCheckInsRequest) (*checkinpb.ListCheckInsResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	records, total, err := h.svc.Checkin.ListCheckIns(ctx, req.EventId, req.UserId, req.Gate, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list check-ins failed: %v", err)
	}

	pbRecords := make([]*checkinpb.CheckInRecord, 0, len(records))
	for _, r := range records {
		pbRecords = append(pbRecords, modelToProto(r))
	}

	return &checkinpb.ListCheckInsResponse{
		Success: true,
		Records: pbRecords,
		Total:   int32(total),
		Page:    int32(page),
		Limit:   int32(limit),
		HasMore: int32(page*limit) < int32(total),
	}, nil
}

// GetEventStats returns aggregated check-in statistics for an event.
func (h *checkinHandler) GetEventStats(ctx context.Context, req *checkinpb.GetEventStatsRequest) (*checkinpb.GetEventStatsResponse, error) {
	stats, err := h.svc.Checkin.GetEventStats(ctx, req.EventId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get event stats failed: %v", err)
	}

	byGate := make(map[string]int32, len(stats.ByGate))
	for g, count := range stats.ByGate {
		byGate[g] = int32(count)
	}

	var lastCheckinAt int64
	if stats.LastCheckinAt != nil {
		lastCheckinAt = stats.LastCheckinAt.UnixMilli()
	}

	return &checkinpb.GetEventStatsResponse{
		Success:       true,
		TotalCheckins: int32(stats.TotalCheckins),
		UniqueTickets: int32(stats.UniqueTickets),
		ByGate:        byGate,
		LastCheckinAt: lastCheckinAt,
	}, nil
}

// Health returns the service health status.
func (h *checkinHandler) Health(_ context.Context, _ *checkinpb.HealthRequest) (*checkinpb.HealthResponse, error) {
	return &checkinpb.HealthResponse{
		Status:  "ok",
		Message: "checkin-service is healthy",
	}, nil
}

// modelToProto converts a models.CheckIn to its proto representation.
func modelToProto(m *models.CheckIn) *checkinpb.CheckInRecord {
	r := &checkinpb.CheckInRecord{
		Id:          m.ID,
		TicketId:    m.TicketID,
		EventId:     m.EventID,
		UserId:      m.UserID,
		QrCode:      m.QRCode,
		Status:      m.Status,
		CheckInTime: m.CheckInTime.UnixMilli(),
		CreatedAt:   m.CreatedAt.UnixMilli(),
	}
	if m.StaffID != nil {
		r.StaffId = *m.StaffID
	}
	if m.DeviceID != nil {
		r.DeviceId = *m.DeviceID
	}
	if m.Gate != nil {
		r.Gate = *m.Gate
	}
	if m.Notes != nil {
		r.Notes = *m.Notes
	}
	return r
}

// domainErrToProto maps a services.CheckInError to its proto enum equivalent.
func domainErrToProto(e *services.CheckInError) checkinpb.CheckInError {
	switch e.Code {
	case "INVALID_TICKET":
		return checkinpb.CheckInError_INVALID_TICKET
	case "ALREADY_CHECKED_IN":
		return checkinpb.CheckInError_ALREADY_CHECKED_IN
	case "CANCELLED_TICKET":
		return checkinpb.CheckInError_CANCELLED_TICKET
	case "TICKET_EVENT_MISMATCH":
		return checkinpb.CheckInError_TICKET_EVENT_MISMATCH
	default:
		return checkinpb.CheckInError_CHECKIN_ERROR_UNSPECIFIED
	}
}

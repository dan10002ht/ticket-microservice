package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "realtime-service/internal/protos"
	"realtime-service/internal/service"
	"realtime-service/pkg/logger"

	"go.uber.org/zap"
)

// NotificationHandler handles gRPC requests for realtime notifications
type NotificationHandler struct {
	pb.UnimplementedRealtimeServiceServer
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// NotifyBookingResult handles booking result notifications
func (h *NotificationHandler) NotifyBookingResult(ctx context.Context, req *pb.NotifyBookingResultRequest) (*pb.NotifyBookingResultResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.GetBookingId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "booking_id is required")
	}

	delivered, activeConnections := h.notificationService.NotifyBookingResult(
		req.GetUserId(),
		req.GetBookingId(),
		req.GetSuccess(),
		req.GetMessage(),
		req.GetBookingReference(),
		req.GetEventId(),
		req.GetSeatNumbers(),
		req.GetTotalAmount(),
		req.GetCurrency(),
	)

	logger.Info("NotifyBookingResult called",
		zap.String("user_id", req.GetUserId()),
		zap.String("booking_id", req.GetBookingId()),
		zap.Bool("success", req.GetSuccess()),
		zap.Bool("delivered", delivered),
		zap.Int("active_connections", activeConnections),
	)

	return &pb.NotifyBookingResultResponse{
		Delivered:         delivered,
		ActiveConnections: int32(activeConnections),
	}, nil
}

// NotifyQueuePosition handles queue position notifications
func (h *NotificationHandler) NotifyQueuePosition(ctx context.Context, req *pb.NotifyQueuePositionRequest) (*pb.NotifyQueuePositionResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.GetEventId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "event_id is required")
	}

	delivered := h.notificationService.NotifyQueuePosition(
		req.GetUserId(),
		req.GetEventId(),
		int(req.GetPosition()),
		int(req.GetEstimatedWaitSeconds()),
		int(req.GetTotalInQueue()),
	)

	return &pb.NotifyQueuePositionResponse{
		Delivered: delivered,
	}, nil
}

// NotifyPaymentStatus handles payment status notifications
func (h *NotificationHandler) NotifyPaymentStatus(ctx context.Context, req *pb.NotifyPaymentStatusRequest) (*pb.NotifyPaymentStatusResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.GetBookingId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "booking_id is required")
	}

	statusStr := paymentStatusToString(req.GetStatus())

	delivered := h.notificationService.NotifyPaymentStatus(
		req.GetUserId(),
		req.GetBookingId(),
		req.GetPaymentId(),
		statusStr,
		req.GetMessage(),
		req.GetAmount(),
		req.GetCurrency(),
	)

	logger.Info("NotifyPaymentStatus called",
		zap.String("user_id", req.GetUserId()),
		zap.String("booking_id", req.GetBookingId()),
		zap.String("status", statusStr),
		zap.Bool("delivered", delivered),
	)

	return &pb.NotifyPaymentStatusResponse{
		Delivered: delivered,
	}, nil
}

// BroadcastEvent handles event broadcasting
func (h *NotificationHandler) BroadcastEvent(ctx context.Context, req *pb.BroadcastEventRequest) (*pb.BroadcastEventResponse, error) {
	if req.GetEventType() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "event_type is required")
	}
	if req.GetRoom() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "room is required")
	}

	recipients := h.notificationService.BroadcastToRoom(
		req.GetEventType(),
		req.GetRoom(),
		req.GetPayload(),
	)

	logger.Info("BroadcastEvent called",
		zap.String("event_type", req.GetEventType()),
		zap.String("room", req.GetRoom()),
		zap.Int("recipients", recipients),
	)

	return &pb.BroadcastEventResponse{
		Recipients: int32(recipients),
	}, nil
}

// SendToUser sends a message to a specific user
func (h *NotificationHandler) SendToUser(ctx context.Context, req *pb.SendToUserRequest) (*pb.SendToUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.GetEventType() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "event_type is required")
	}

	delivered := h.notificationService.SendToUser(
		req.GetUserId(),
		req.GetEventType(),
		req.GetPayload(),
	)

	return &pb.SendToUserResponse{
		Delivered:            delivered,
		ConnectionsNotified: 1, // TODO: Return actual count
	}, nil
}

// GetConnectionStats returns connection statistics
func (h *NotificationHandler) GetConnectionStats(ctx context.Context, req *pb.GetConnectionStatsRequest) (*pb.GetConnectionStatsResponse, error) {
	stats := h.notificationService.GetConnectionStats(req.GetRoom())

	roomsMap := make(map[string]int32)
	if rooms, ok := stats["rooms"].(map[string]int); ok {
		for room, count := range rooms {
			roomsMap[room] = int32(count)
		}
	}

	return &pb.GetConnectionStatsResponse{
		TotalConnections:         int32(stats["total_connections"].(int64)),
		AuthenticatedConnections: int32(stats["authenticated_connections"].(int64)),
		AnonymousConnections:     int32(stats["anonymous_connections"].(int64)),
		ConnectionsPerRoom:       roomsMap,
		UptimeSeconds:            stats["uptime_seconds"].(int64),
		MessagesSentTotal:        stats["messages_sent_total"].(int64),
		MessagesReceivedTotal:    stats["messages_received_total"].(int64),
	}, nil
}

// Helper function to convert PaymentStatus enum to string
func paymentStatusToString(status pb.PaymentStatus) string {
	switch status {
	case pb.PaymentStatus_PAYMENT_STATUS_PROCESSING:
		return "processing"
	case pb.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		return "success"
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return "failed"
	case pb.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return "refunded"
	default:
		return "unknown"
	}
}

package grpc

import (
	"context"
	"fmt"
	"time"

	"booking-worker/config"
	pb "booking-worker/internal/protos/booking_worker"
	"booking-worker/internal/queue"
	"booking-worker/internal/worker"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BookingWorkerService implements the BookingWorkerService gRPC service
type BookingWorkerService struct {
	pb.UnimplementedBookingWorkerServiceServer
	queue     queue.QueueManager
	processor *worker.Processor
	config    *config.Config
	logger    *zap.Logger
}

// NewBookingWorkerService creates a new BookingWorkerService
func NewBookingWorkerService(q queue.QueueManager, p *worker.Processor, cfg *config.Config, logger *zap.Logger) *BookingWorkerService {
	return &BookingWorkerService{
		queue:     q,
		processor: p,
		config:    cfg,
		logger:    logger,
	}
}

// EnqueueBooking enqueues a booking request
func (s *BookingWorkerService) EnqueueBooking(ctx context.Context, req *pb.EnqueueBookingRequest) (*pb.EnqueueBookingResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.EventId == "" {
		return nil, status.Error(codes.InvalidArgument, "event_id is required")
	}
	if len(req.SeatNumbers) == 0 && req.SeatCount == 0 {
		return nil, status.Error(codes.InvalidArgument, "seat_numbers or seat_count is required")
	}

	// Generate queue item ID
	itemID := uuid.New().String()

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Duration(s.config.Queue.TimeoutSeconds) * time.Second)

	// Create queue item
	item := &queue.QueueItem{
		ID:          itemID,
		EventID:     req.EventId,
		UserID:      req.UserId,
		SeatNumbers: req.SeatNumbers,
		SeatCount:   int(req.SeatCount),
		TotalAmount: req.TotalAmount,
		Currency:    req.Currency,
		Metadata:    req.Metadata,
		EnqueuedAt:  time.Now(),
		ExpiresAt:   expiresAt,
	}

	// Enqueue
	if err := s.queue.Enqueue(ctx, item); err != nil {
		s.logger.Error("Failed to enqueue booking",
			zap.String("user_id", req.UserId),
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to enqueue booking: %v", err))
	}

	s.logger.Info("Booking enqueued successfully",
		zap.String("queue_item_id", itemID),
		zap.String("user_id", req.UserId),
		zap.String("event_id", req.EventId),
	)

	// Get initial position
	position, err := s.queue.GetPosition(ctx, itemID)
	if err != nil {
		s.logger.Warn("Failed to get queue position",
			zap.String("item_id", itemID),
			zap.Error(err),
		)
		position = -1 // Unknown position
	}

	// Estimate wait time (rough estimate: 5 seconds per item)
	estimatedWait := position * 5
	if position < 0 {
		estimatedWait = 0
	}

	return &pb.EnqueueBookingResponse{
		Success:              true,
		QueueItemId:          itemID,
		QueuePosition:        int32(position),
		EstimatedWaitSeconds: int32(estimatedWait),
		Message:              "Booking request enqueued successfully",
	}, nil
}

// GetQueuePosition gets the current position of a queue item
func (s *BookingWorkerService) GetQueuePosition(ctx context.Context, req *pb.GetQueuePositionRequest) (*pb.GetQueuePositionResponse, error) {
	// Validate request
	if req.QueueItemId == "" {
		return nil, status.Error(codes.InvalidArgument, "queue_item_id is required")
	}

	position, err := s.queue.GetPosition(ctx, req.QueueItemId)
	if err != nil {
		s.logger.Warn("Queue item not found",
			zap.String("queue_item_id", req.QueueItemId),
			zap.Error(err),
		)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("queue item not found: %v", err))
	}

	// Determine status based on position
	queueStatus := "waiting"
	if position == 0 {
		queueStatus = "processing"
	} else if position < 0 {
		queueStatus = "completed"
	}

	// Estimate wait time
	estimatedWait := position * 5
	if position < 0 {
		estimatedWait = 0
	}

	return &pb.GetQueuePositionResponse{
		Success:              true,
		Position:             int32(position),
		EstimatedWaitSeconds: int32(estimatedWait),
		Status:               queueStatus,
		Message:              "Queue position retrieved",
	}, nil
}

// GetQueueStatus gets the status of a queue for an event
func (s *BookingWorkerService) GetQueueStatus(ctx context.Context, req *pb.GetQueueStatusRequest) (*pb.GetQueueStatusResponse, error) {
	// Validate request
	if req.EventId == "" {
		return nil, status.Error(codes.InvalidArgument, "event_id is required")
	}

	queueLength, err := s.queue.GetQueueLength(ctx, req.EventId)
	if err != nil {
		s.logger.Error("Failed to get queue length",
			zap.String("event_id", req.EventId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get queue length: %v", err))
	}

	// Get active workers count from processor
	activeWorkers := s.processor.GetWorkerCount()

	// Calculate processing rate from metrics (simplified: assume 1 item per 5 seconds per worker)
	processingRate := float64(activeWorkers) / 5.0

	return &pb.GetQueueStatusResponse{
		Success:        true,
		QueueLength:    int32(queueLength),
		ActiveWorkers:  int32(activeWorkers),
		ProcessingRate: processingRate,
		Message:        "Queue status retrieved",
	}, nil
}

// CancelQueueItem cancels a queue item with authorization check
func (s *BookingWorkerService) CancelQueueItem(ctx context.Context, req *pb.CancelQueueItemRequest) (*pb.CancelQueueItemResponse, error) {
	// Validate request
	if req.QueueItemId == "" {
		return nil, status.Error(codes.InvalidArgument, "queue_item_id is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Remove with authorization check - only the owner can cancel their queue item
	if err := s.queue.RemoveWithAuthorization(ctx, req.QueueItemId, req.UserId); err != nil {
		s.logger.Warn("Failed to cancel queue item",
			zap.String("queue_item_id", req.QueueItemId),
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)

		// Check if it's an authorization error
		if err.Error() == "unauthorized: user does not own this queue item" {
			return nil, status.Error(codes.PermissionDenied, "you are not authorized to cancel this queue item")
		}

		return nil, status.Error(codes.NotFound, fmt.Sprintf("queue item not found: %v", err))
	}

	s.logger.Info("Queue item cancelled",
		zap.String("queue_item_id", req.QueueItemId),
		zap.String("user_id", req.UserId),
		zap.String("reason", req.Reason),
	)

	return &pb.CancelQueueItemResponse{
		Success: true,
		Message: "Queue item cancelled successfully",
	}, nil
}

// Health performs a health check
func (s *BookingWorkerService) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	// Check Redis connection by getting queue length
	_, err := s.queue.GetQueueLength(ctx, "health-check")
	redisStatus := "healthy"
	if err != nil {
		redisStatus = "unhealthy"
	}

	details := map[string]string{
		"service":        "booking-worker",
		"version":        "1.0.0",
		"redis":          redisStatus,
		"active_workers": fmt.Sprintf("%d", s.processor.GetWorkerCount()),
	}

	overallStatus := "healthy"
	if redisStatus == "unhealthy" {
		overallStatus = "degraded"
	}

	return &pb.HealthResponse{
		Status:  overallStatus,
		Message: "Booking worker service health check",
		Details: details,
	}, nil
}

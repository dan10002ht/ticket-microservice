package service

import (
	"realtime-service/internal/websocket"
	"realtime-service/pkg/logger"

	"go.uber.org/zap"
)

// NotificationService handles sending notifications through WebSocket
type NotificationService struct {
	hub *websocket.Hub
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(hub *websocket.Hub) *NotificationService {
	return &NotificationService{hub: hub}
}

// NotifyBookingResult sends booking result to a user
func (s *NotificationService) NotifyBookingResult(userID, bookingID string, success bool, message, bookingRef, eventID string, seatNumbers []string, totalAmount, currency string) (bool, int) {
	payload := websocket.BookingResultPayload{
		BookingID:        bookingID,
		BookingReference: bookingRef,
		Success:          success,
		Message:          message,
		EventID:          eventID,
		SeatNumbers:      seatNumbers,
		TotalAmount:      totalAmount,
		Currency:         currency,
	}

	var msgType string
	if success {
		msgType = websocket.TypeBookingConfirmed
	} else {
		msgType = websocket.TypeBookingFailed
	}

	msg, err := websocket.NewMessage(msgType, payload)
	if err != nil {
		logger.Error("failed to create booking result message", zap.Error(err))
		return false, 0
	}

	data, err := msg.Bytes()
	if err != nil {
		logger.Error("failed to serialize message", zap.Error(err))
		return false, 0
	}

	count := s.hub.SendToUser(userID, data)
	delivered := count > 0

	logger.Info("booking result notification sent",
		zap.String("user_id", userID),
		zap.String("booking_id", bookingID),
		zap.Bool("success", success),
		zap.Bool("delivered", delivered),
		zap.Int("connections", count),
	)

	return delivered, count
}

// NotifyQueuePosition sends queue position update to a user
func (s *NotificationService) NotifyQueuePosition(userID, eventID string, position, estimatedWait, totalInQueue int) bool {
	payload := websocket.QueuePositionPayload{
		Position:      position,
		EstimatedWait: estimatedWait,
		EventID:       eventID,
		TotalInQueue:  totalInQueue,
	}

	msg, err := websocket.NewMessage(websocket.TypeBookingQueuePosition, payload)
	if err != nil {
		logger.Error("failed to create queue position message", zap.Error(err))
		return false
	}

	data, err := msg.Bytes()
	if err != nil {
		logger.Error("failed to serialize message", zap.Error(err))
		return false
	}

	count := s.hub.SendToUser(userID, data)
	return count > 0
}

// NotifyPaymentStatus sends payment status update to a user
func (s *NotificationService) NotifyPaymentStatus(userID, bookingID, paymentID, status, message, amount, currency string) bool {
	payload := websocket.PaymentStatusPayload{
		PaymentID: paymentID,
		BookingID: bookingID,
		Status:    status,
		Message:   message,
		Amount:    amount,
		Currency:  currency,
	}

	var msgType string
	switch status {
	case "processing":
		msgType = websocket.TypePaymentProcessing
	case "success":
		msgType = websocket.TypePaymentSuccess
	case "failed":
		msgType = websocket.TypePaymentFailed
	default:
		msgType = websocket.TypePaymentProcessing
	}

	msg, err := websocket.NewMessage(msgType, payload)
	if err != nil {
		logger.Error("failed to create payment status message", zap.Error(err))
		return false
	}

	data, err := msg.Bytes()
	if err != nil {
		logger.Error("failed to serialize message", zap.Error(err))
		return false
	}

	count := s.hub.SendToUser(userID, data)

	logger.Info("payment status notification sent",
		zap.String("user_id", userID),
		zap.String("payment_id", paymentID),
		zap.String("status", status),
		zap.Int("connections", count),
	)

	return count > 0
}

// BroadcastToRoom broadcasts a message to all clients in a room
func (s *NotificationService) BroadcastToRoom(eventType, room, payload string) int {
	msg := &websocket.Message{
		Type:    eventType,
		Payload: []byte(payload),
		Room:    room,
	}

	data, err := msg.Bytes()
	if err != nil {
		logger.Error("failed to serialize broadcast message", zap.Error(err))
		return 0
	}

	count := s.hub.BroadcastToRoom(room, data)

	logger.Info("broadcast sent",
		zap.String("event_type", eventType),
		zap.String("room", room),
		zap.Int("recipients", count),
	)

	return count
}

// SendToUser sends a message to a specific user
func (s *NotificationService) SendToUser(userID, eventType, payload string) bool {
	msg := &websocket.Message{
		Type:    eventType,
		Payload: []byte(payload),
		UserID:  userID,
	}

	data, err := msg.Bytes()
	if err != nil {
		logger.Error("failed to serialize user message", zap.Error(err))
		return false
	}

	count := s.hub.SendToUser(userID, data)
	return count > 0
}

// GetConnectionStats returns connection statistics
func (s *NotificationService) GetConnectionStats(room string) map[string]interface{} {
	stats := s.hub.GetStats()

	if room != "" {
		stats["room_connections"] = s.hub.GetRoomConnectionCount(room)
	}

	return stats
}

package websocket

import (
	"encoding/json"
)

// Message types
const (
	// Booking events
	TypeBookingQueuePosition = "booking:queue_position"
	TypeBookingProcessing    = "booking:processing"
	TypeBookingConfirmed     = "booking:confirmed"
	TypeBookingFailed        = "booking:failed"
	TypeBookingCancelled     = "booking:cancelled"

	// Payment events
	TypePaymentProcessing = "payment:processing"
	TypePaymentSuccess    = "payment:success"
	TypePaymentFailed     = "payment:failed"

	// Ticket events
	TypeTicketAvailability = "ticket:availability"
	TypeTicketReserved     = "ticket:reserved"
	TypeTicketReleased     = "ticket:released"

	// System events
	TypeSystemConnected = "system:connected"
	TypeSystemError     = "system:error"
	TypeSystemPing      = "system:ping"
	TypeSystemPong      = "system:pong"

	// Client actions
	TypeRoomJoin      = "room:join"
	TypeRoomLeave     = "room:leave"
	TypeUserSubscribe = "user:subscribe"
)

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Room    string          `json:"room,omitempty"`
	UserID  string          `json:"user_id,omitempty"`
}

// NewMessage creates a new message with the given type and payload
func NewMessage(msgType string, payload interface{}) (*Message, error) {
	var payloadBytes json.RawMessage
	if payload != nil {
		var err error
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:    msgType,
		Payload: payloadBytes,
	}, nil
}

// Parse parses a raw JSON message into a Message struct
func Parse(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Bytes serializes the message to JSON bytes
func (m *Message) Bytes() ([]byte, error) {
	return json.Marshal(m)
}

// ConnectedPayload represents the payload for system:connected event
type ConnectedPayload struct {
	ConnectionID string `json:"connection_id"`
	ServerTime   int64  `json:"server_time"`
}

// QueuePositionPayload represents the payload for booking:queue_position event
type QueuePositionPayload struct {
	Position      int    `json:"position"`
	EstimatedWait int    `json:"estimated_wait"`
	EventID       string `json:"event_id"`
	TotalInQueue  int    `json:"total_in_queue"`
}

// BookingResultPayload represents the payload for booking result events
type BookingResultPayload struct {
	BookingID        string   `json:"booking_id"`
	BookingReference string   `json:"booking_reference"`
	Success          bool     `json:"success"`
	Message          string   `json:"message"`
	EventID          string   `json:"event_id,omitempty"`
	SeatNumbers      []string `json:"seat_numbers,omitempty"`
	TotalAmount      string   `json:"total_amount,omitempty"`
	Currency         string   `json:"currency,omitempty"`
}

// PaymentStatusPayload represents the payload for payment status events
type PaymentStatusPayload struct {
	PaymentID string `json:"payment_id"`
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Amount    string `json:"amount,omitempty"`
	Currency  string `json:"currency,omitempty"`
}

// ErrorPayload represents the payload for error events
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RoomJoinPayload represents the payload for room:join events
type RoomJoinPayload struct {
	EventID string `json:"event_id"`
}

// UserSubscribePayload represents the payload for user:subscribe events
type UserSubscribePayload struct {
	UserID string `json:"user_id"`
}

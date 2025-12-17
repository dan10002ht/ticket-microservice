package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"realtime-service/config"
	"realtime-service/pkg/logger"

	"go.uber.org/zap"
)

// Client represents a WebSocket client connection
type Client struct {
	ID     string
	UserID string // Empty if not authenticated
	conn   *websocket.Conn
	hub    *Hub
	send   chan []byte
	rooms  map[string]bool
	mu     sync.RWMutex
	cfg    *config.WebSocketConfig
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, userID string, cfg *config.WebSocketConfig) *Client {
	return &Client{
		ID:     uuid.New().String(),
		UserID: userID,
		conn:   conn,
		hub:    hub,
		send:   make(chan []byte, 256),
		rooms:  make(map[string]bool),
		cfg:    cfg,
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(c.cfg.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("websocket error", zap.String("client_id", c.ID), zap.Error(err))
			}
			break
		}

		c.handleMessage(message)
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(c.cfg.PingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(data []byte) {
	msg, err := Parse(data)
	if err != nil {
		logger.Error("failed to parse message", zap.String("client_id", c.ID), zap.Error(err))
		c.sendError("PARSE_ERROR", "Invalid message format")
		return
	}

	switch msg.Type {
	case TypeRoomJoin:
		c.handleRoomJoin(msg)
	case TypeRoomLeave:
		c.handleRoomLeave(msg)
	case TypeUserSubscribe:
		c.handleUserSubscribe(msg)
	case TypeSystemPing:
		c.handlePing()
	default:
		logger.Debug("unknown message type", zap.String("type", msg.Type), zap.String("client_id", c.ID))
	}
}

// handleRoomJoin handles room join requests
func (c *Client) handleRoomJoin(msg *Message) {
	var payload RoomJoinPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		c.sendError("INVALID_PAYLOAD", "Invalid room join payload")
		return
	}

	if payload.EventID == "" {
		c.sendError("MISSING_EVENT_ID", "Event ID is required")
		return
	}

	roomID := "event:" + payload.EventID
	c.JoinRoom(roomID)
	c.hub.JoinRoom(c, roomID)

	logger.Info("client joined room",
		zap.String("client_id", c.ID),
		zap.String("room", roomID),
		zap.String("user_id", c.UserID),
	)
}

// handleRoomLeave handles room leave requests
func (c *Client) handleRoomLeave(msg *Message) {
	var payload RoomJoinPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		c.sendError("INVALID_PAYLOAD", "Invalid room leave payload")
		return
	}

	roomID := "event:" + payload.EventID
	c.LeaveRoom(roomID)
	c.hub.LeaveRoom(c, roomID)

	logger.Info("client left room",
		zap.String("client_id", c.ID),
		zap.String("room", roomID),
	)
}

// handleUserSubscribe handles user subscription requests
func (c *Client) handleUserSubscribe(msg *Message) {
	var payload UserSubscribePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		c.sendError("INVALID_PAYLOAD", "Invalid subscribe payload")
		return
	}

	// Only allow subscribing to own user channel if authenticated
	if c.UserID != "" && c.UserID != payload.UserID {
		c.sendError("UNAUTHORIZED", "Cannot subscribe to other user's channel")
		return
	}

	roomID := "user:" + payload.UserID
	c.JoinRoom(roomID)
	c.hub.JoinRoom(c, roomID)

	logger.Info("client subscribed to user channel",
		zap.String("client_id", c.ID),
		zap.String("user_id", payload.UserID),
	)
}

// handlePing handles ping messages
func (c *Client) handlePing() {
	msg, _ := NewMessage(TypeSystemPong, nil)
	data, _ := msg.Bytes()
	c.Send(data)
}

// Send sends a message to the client
func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
		// Channel is full, client is slow
		logger.Warn("client send buffer full", zap.String("client_id", c.ID))
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(code, message string) {
	msg, _ := NewMessage(TypeSystemError, ErrorPayload{
		Code:    code,
		Message: message,
	})
	data, _ := msg.Bytes()
	c.Send(data)
}

// JoinRoom adds the client to a room
func (c *Client) JoinRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rooms[roomID] = true
}

// LeaveRoom removes the client from a room
func (c *Client) LeaveRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.rooms, roomID)
}

// GetRooms returns the list of rooms the client is in
func (c *Client) GetRooms() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	rooms := make([]string, 0, len(c.rooms))
	for room := range c.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// IsInRoom checks if the client is in a specific room
func (c *Client) IsInRoom(roomID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rooms[roomID]
}

// Close closes the client connection
func (c *Client) Close() {
	close(c.send)
}

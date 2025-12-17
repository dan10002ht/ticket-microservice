package websocket

import (
	"sync"
	"sync/atomic"
	"time"

	"realtime-service/pkg/logger"

	"go.uber.org/zap"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Clients by user ID for direct messaging
	userClients map[string]map[*Client]bool

	// Room memberships
	rooms map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast to all clients
	broadcast chan []byte

	// Room-specific broadcast
	roomBroadcast chan *RoomMessage

	// User-specific message
	userMessage chan *UserMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Statistics
	stats *Stats

	// Start time for uptime calculation
	startTime time.Time
}

// RoomMessage represents a message to be broadcast to a room
type RoomMessage struct {
	Room    string
	Message []byte
}

// UserMessage represents a message to be sent to a specific user
type UserMessage struct {
	UserID  string
	Message []byte
}

// Stats holds connection statistics
type Stats struct {
	TotalConnections       int64
	AuthenticatedConns     int64
	AnonymousConns         int64
	MessagesSentTotal      int64
	MessagesReceivedTotal  int64
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		userClients:   make(map[string]map[*Client]bool),
		rooms:         make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan []byte),
		roomBroadcast: make(chan *RoomMessage, 256),
		userMessage:   make(chan *UserMessage, 256),
		stats:         &Stats{},
		startTime:     time.Now(),
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToAll(message)

		case roomMsg := <-h.roomBroadcast:
			h.broadcastToRoom(roomMsg)

		case userMsg := <-h.userMessage:
			h.sendToUser(userMsg)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	atomic.AddInt64(&h.stats.TotalConnections, 1)

	if client.UserID != "" {
		atomic.AddInt64(&h.stats.AuthenticatedConns, 1)
		if h.userClients[client.UserID] == nil {
			h.userClients[client.UserID] = make(map[*Client]bool)
		}
		h.userClients[client.UserID][client] = true

		// Auto-subscribe authenticated users to their user channel
		userRoom := "user:" + client.UserID
		if h.rooms[userRoom] == nil {
			h.rooms[userRoom] = make(map[*Client]bool)
		}
		h.rooms[userRoom][client] = true
		client.JoinRoom(userRoom)
	} else {
		atomic.AddInt64(&h.stats.AnonymousConns, 1)
	}

	logger.Info("client registered",
		zap.String("client_id", client.ID),
		zap.String("user_id", client.UserID),
		zap.Int64("total_connections", atomic.LoadInt64(&h.stats.TotalConnections)),
	)
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; !ok {
		return
	}

	delete(h.clients, client)
	atomic.AddInt64(&h.stats.TotalConnections, -1)

	if client.UserID != "" {
		atomic.AddInt64(&h.stats.AuthenticatedConns, -1)
		if userClients, ok := h.userClients[client.UserID]; ok {
			delete(userClients, client)
			if len(userClients) == 0 {
				delete(h.userClients, client.UserID)
			}
		}
	} else {
		atomic.AddInt64(&h.stats.AnonymousConns, -1)
	}

	// Remove from all rooms
	for room := range client.rooms {
		if roomClients, ok := h.rooms[room]; ok {
			delete(roomClients, client)
			if len(roomClients) == 0 {
				delete(h.rooms, room)
			}
		}
	}

	client.Close()

	logger.Info("client unregistered",
		zap.String("client_id", client.ID),
		zap.String("user_id", client.UserID),
		zap.Int64("total_connections", atomic.LoadInt64(&h.stats.TotalConnections)),
	)
}

// broadcastToAll sends a message to all connected clients
func (h *Hub) broadcastToAll(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		client.Send(message)
	}

	atomic.AddInt64(&h.stats.MessagesSentTotal, int64(len(h.clients)))
}

// broadcastToRoom sends a message to all clients in a room
func (h *Hub) broadcastToRoom(roomMsg *RoomMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomClients, ok := h.rooms[roomMsg.Room]; ok {
		for client := range roomClients {
			client.Send(roomMsg.Message)
		}
		atomic.AddInt64(&h.stats.MessagesSentTotal, int64(len(roomClients)))
	}
}

// sendToUser sends a message to all connections of a specific user
func (h *Hub) sendToUser(userMsg *UserMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if userClients, ok := h.userClients[userMsg.UserID]; ok {
		for client := range userClients {
			client.Send(userMsg.Message)
		}
		atomic.AddInt64(&h.stats.MessagesSentTotal, int64(len(userClients)))
	}
}

// Register registers a client with the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast broadcasts a message to all clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// BroadcastToRoom broadcasts a message to all clients in a room
func (h *Hub) BroadcastToRoom(room string, message []byte) int {
	h.roomBroadcast <- &RoomMessage{Room: room, Message: message}

	h.mu.RLock()
	count := len(h.rooms[room])
	h.mu.RUnlock()
	return count
}

// SendToUser sends a message to all connections of a specific user
func (h *Hub) SendToUser(userID string, message []byte) int {
	h.userMessage <- &UserMessage{UserID: userID, Message: message}

	h.mu.RLock()
	count := len(h.userClients[userID])
	h.mu.RUnlock()
	return count
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if roomClients, ok := h.rooms[roomID]; ok {
		delete(roomClients, client)
		if len(roomClients) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// GetStats returns current connection statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	roomCounts := make(map[string]int)
	for room, clients := range h.rooms {
		roomCounts[room] = len(clients)
	}

	return map[string]interface{}{
		"total_connections":        atomic.LoadInt64(&h.stats.TotalConnections),
		"authenticated_connections": atomic.LoadInt64(&h.stats.AuthenticatedConns),
		"anonymous_connections":     atomic.LoadInt64(&h.stats.AnonymousConns),
		"messages_sent_total":       atomic.LoadInt64(&h.stats.MessagesSentTotal),
		"messages_received_total":   atomic.LoadInt64(&h.stats.MessagesReceivedTotal),
		"uptime_seconds":           int64(time.Since(h.startTime).Seconds()),
		"rooms":                    roomCounts,
	}
}

// GetConnectionCount returns the total number of connections
func (h *Hub) GetConnectionCount() int64 {
	return atomic.LoadInt64(&h.stats.TotalConnections)
}

// GetUserConnectionCount returns the number of connections for a user
func (h *Hub) GetUserConnectionCount(userID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients[userID])
}

// GetRoomConnectionCount returns the number of connections in a room
func (h *Hub) GetRoomConnectionCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}

// Stop gracefully stops the hub by closing all client connections
func (h *Hub) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	logger.Info("Stopping hub, closing all connections",
		zap.Int("client_count", len(h.clients)),
	)

	for client := range h.clients {
		client.Close()
	}
}

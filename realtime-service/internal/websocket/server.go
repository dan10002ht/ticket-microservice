package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"realtime-service/config"
	"realtime-service/internal/middleware"
	"realtime-service/pkg/logger"

	"go.uber.org/zap"
)

// Server represents the WebSocket server
type Server struct {
	hub      *Hub
	upgrader websocket.Upgrader
	auth     *middleware.AuthMiddleware
	cfg      *config.Config
}

// NewServer creates a new WebSocket server
func NewServer(hub *Hub, auth *middleware.AuthMiddleware, cfg *config.Config) *Server {
	return &Server{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:    cfg.WebSocket.ReadBufferSize,
			WriteBufferSize:   cfg.WebSocket.WriteBufferSize,
			EnableCompression: cfg.WebSocket.EnableCompression,
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implement proper origin checking in production
				return true
			},
		},
		auth: auth,
		cfg:  cfg,
	}
}

// HandleWebSocket handles WebSocket connection upgrade requests
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request
	userID, err := s.auth.ValidateRequest(r)
	if err != nil {
		logger.Warn("websocket authentication failed", zap.Error(err))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("failed to upgrade connection", zap.Error(err))
		return
	}

	// Create new client
	client := NewClient(conn, s.hub, userID, &s.cfg.WebSocket)

	// Register client with hub
	s.hub.Register(client)

	// Send connected event
	connectedMsg, _ := NewMessage(TypeSystemConnected, ConnectedPayload{
		ConnectionID: client.ID,
		ServerTime:   time.Now().Unix(),
	})
	data, _ := connectedMsg.Bytes()
	client.Send(data)

	logger.Info("new websocket connection",
		zap.String("client_id", client.ID),
		zap.String("user_id", userID),
		zap.String("remote_addr", r.RemoteAddr),
	)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// GetHub returns the hub instance
func (s *Server) GetHub() *Hub {
	return s.hub
}

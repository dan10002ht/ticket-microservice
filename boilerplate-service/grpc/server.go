package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"boilerplate-service/internal/config"
	"boilerplate-service/services"
)

type Server struct {
	services *services.Services
	config   *config.Config
	logger   *zap.Logger
	server   *grpc.Server
}

func NewServer(services *services.Services, cfg *config.Config, logger *zap.Logger) *Server {
	return &Server{
		services: services,
		config:   cfg,
		logger:   logger,
	}
}

func (s *Server) Start(port int) error {
	// Create gRPC server
	s.server = grpc.NewServer()

	// Register services here
	// Example: pb.RegisterUserServiceServer(s.server, s)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.logger.Info("Starting gRPC server", zap.Int("port", port))

	// Start server
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
} 
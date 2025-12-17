package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"realtime-service/config"
	"realtime-service/internal/grpc/handlers"
	pb "realtime-service/internal/protos"
	"realtime-service/internal/service"
	"realtime-service/pkg/logger"

	"go.uber.org/zap"
)

// Server wraps the gRPC server
type Server struct {
	grpcServer          *grpc.Server
	config              *config.Config
	notificationService *service.NotificationService
}

// NewServer creates a new gRPC server
func NewServer(cfg *config.Config, notificationService *service.NotificationService) *Server {
	// Create gRPC server with options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(4 * 1024 * 1024), // 4MB
		grpc.MaxSendMsgSize(4 * 1024 * 1024), // 4MB
	}

	grpcServer := grpc.NewServer(opts...)

	// Register notification handler
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	pb.RegisterRealtimeServiceServer(grpcServer, notificationHandler)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("realtime.RealtimeService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	return &Server{
		grpcServer:          grpcServer,
		config:              cfg,
		notificationService: notificationService,
	}
}

// Start starts the gRPC server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	logger.Info("Starting gRPC server", zap.String("address", addr))

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	logger.Info("Stopping gRPC server gracefully")
	s.grpcServer.GracefulStop()
}

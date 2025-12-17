package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"user-service/config"
	"user-service/internal/grpc/handlers"
	pb "user-service/internal/protos"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Server represents the gRPC server
type Server struct {
	cfg           *config.Config
	db            *pgxpool.Pool
	grpcServer    *grpc.Server
	healthServer  *health.Server
	metricsServer *http.Server
}

// NewServer creates a new gRPC server
func NewServer(cfg *config.Config, db *pgxpool.Pool) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

// Start starts the gRPC server
func (s *Server) Start(ctx context.Context) error {
	// Create repositories
	profileRepo := repository.NewProfileRepository(s.db)
	addressRepo := repository.NewAddressRepository(s.db)

	// Create services
	profileService := service.NewProfileService(profileRepo)
	addressService := service.NewAddressService(addressRepo)

	// Create handlers
	profileHandler := handlers.NewProfileHandler(profileService, addressService)

	// Create gRPC server with interceptors
	grpcprom.EnableHandlingTimeHistogram()
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcprom.UnaryServerInterceptor,
			loggingInterceptor,
		),
		grpc.ChainStreamInterceptor(
			grpcprom.StreamServerInterceptor,
		),
	)

	// Register services
	pb.RegisterUserServiceServer(s.grpcServer, profileHandler)

	// Register health service
	s.healthServer = health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)
	s.healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for development
	reflection.Register(s.grpcServer)

	// Start metrics server if enabled
	if s.cfg.Metrics.Enabled {
		go s.startMetricsServer()
	}

	// Start gRPC server
	listener, err := net.Listen("tcp", ":"+s.cfg.GRPC.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.cfg.GRPC.Port, err)
	}

	logger.Info("starting gRPC server",
		zap.String("port", s.cfg.GRPC.Port),
		zap.Bool("metrics_enabled", s.cfg.Metrics.Enabled),
	)

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		logger.Info("shutting down gRPC server...")
		s.healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.grpcServer.GracefulStop()
	}()

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// startMetricsServer starts the Prometheus metrics server
func (s *Server) startMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := s.db.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	s.metricsServer = &http.Server{
		Addr:    ":" + s.cfg.Metrics.Port,
		Handler: mux,
	}

	logger.Info("starting metrics server", zap.String("port", s.cfg.Metrics.Port))

	if err := s.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("metrics server error", zap.Error(err))
	}
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	if s.metricsServer != nil {
		s.metricsServer.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	logger.Debug("gRPC request",
		zap.String("method", info.FullMethod),
	)

	resp, err := handler(ctx, req)

	if err != nil {
		logger.Error("gRPC request failed",
			zap.String("method", info.FullMethod),
			zap.Error(err),
		)
	}

	return resp, err
}

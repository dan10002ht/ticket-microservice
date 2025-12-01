package app

import (
	"fmt"
	"net"

	"booking-worker/config"
	"booking-worker/internal/queue"
	"booking-worker/internal/worker"
	"booking-worker/metrics"

	"go.uber.org/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App represents the main application
type App struct {
	logger   *zap.Logger
	config   *config.Config
	queue    queue.QueueManager
	processor *worker.Processor
	grpcServer *grpc.Server
	metrics   *metrics.Exporter
}

// NewApp creates a new application instance
func NewApp(logger *zap.Logger, cfg *config.Config) *App {
	return &App{
		logger: logger,
		config: cfg,
	}
}

// Initialize initializes all application components
func (a *App) Initialize() error {
	a.logger.Info("Initializing application components...")

	// Initialize metrics
	if a.config.Metrics.Enabled {
		metricsExporter, err := metrics.NewExporter(a.config.Metrics.Port, a.logger)
		if err != nil {
			return fmt.Errorf("failed to initialize metrics: %w", err)
		}
		a.metrics = metricsExporter
		a.logger.Info("Metrics initialized", zap.String("port", a.config.Metrics.Port))
	}

	// Initialize queue manager
	queueManager, err := queue.NewRedisQueueManager(a.config, a.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize queue manager: %w", err)
	}
	a.queue = queueManager
	a.logger.Info("Queue manager initialized")

	// Initialize worker processor
	processor, err := worker.NewProcessor(a.config, a.queue, a.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize processor: %w", err)
	}
	a.processor = processor
	a.logger.Info("Worker processor initialized")

	// Initialize gRPC server
	grpcServer, err := a.initGRPCServer()
	if err != nil {
		return fmt.Errorf("failed to initialize gRPC server: %w", err)
	}
	a.grpcServer = grpcServer
	a.logger.Info("gRPC server initialized", zap.String("port", a.config.GRPC.Port))

	return nil
}

// Run starts the application
func (a *App) Run() error {
	a.logger.Info("Starting application...")

	// Start worker processor
	if err := a.processor.Start(); err != nil {
		return fmt.Errorf("failed to start processor: %w", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+a.config.GRPC.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", a.config.GRPC.Port, err)
	}

	a.logger.Info("gRPC server listening", zap.String("port", a.config.GRPC.Port))
	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server failed: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() error {
	a.logger.Info("Shutting down application...")

	// Stop processor
	if a.processor != nil {
		if err := a.processor.Stop(); err != nil {
			a.logger.Error("Error stopping processor", zap.Error(err))
		}
	}

	// Stop gRPC server
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	// Close queue connections
	if a.queue != nil {
		if err := a.queue.Close(); err != nil {
			a.logger.Error("Error closing queue", zap.Error(err))
		}
	}

	a.logger.Info("Application shutdown complete")
	return nil
}

// initGRPCServer initializes the gRPC server
func (a *App) initGRPCServer() (*grpc.Server, error) {
	opts := []grpc.ServerOption{}

	// Add Prometheus interceptors if metrics enabled
	if a.config.Metrics.Enabled {
		opts = append(opts, grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
		opts = append(opts, grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor))
	}

	s := grpc.NewServer(opts...)

	// Register gRPC services
	// TODO: Register BookingWorkerService when protobuf is defined
	// bookingWorkerService := grpcservice.NewBookingWorkerService(a.queue, a.processor, a.logger)
	// bookingpb.RegisterBookingWorkerServiceServer(s, bookingWorkerService)

	// Enable reflection for development
	reflection.Register(s)

	return s, nil
}


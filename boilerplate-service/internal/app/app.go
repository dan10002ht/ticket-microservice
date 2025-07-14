package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"boilerplate-service/database"
	"boilerplate-service/grpcclient"
	"boilerplate-service/internal/config"
	"boilerplate-service/queue"
	"boilerplate-service/services"
)

type App struct {
	logger   *zap.Logger
	config   *config.Config
	services *services.Services
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewApp(logger *zap.Logger, cfg *config.Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		logger: logger,
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Initialize() error {
	a.logger.Info("Initializing application components")

	// Initialize database
	db, err := database.NewConnection(a.config.Database)
	if err != nil {
		return err
	}

	// Initialize Redis
	redis, err := queue.NewRedisClient(a.config.Redis)
	if err != nil {
		return err
	}

	// Initialize gRPC clients
	grpcClients, err := grpcclient.NewClients(a.config.GRPC)
	if err != nil {
		return err
	}

	// Initialize services
	a.services = services.NewServices(db, redis, grpcClients)

	a.logger.Info("Application components initialized successfully")
	return nil
}

func (a *App) GetServices() *services.Services {
	return a.services
}

func (a *App) Run() error {
	a.logger.Info("Starting application")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down application...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close services
	if a.services != nil {
		if err := a.services.Close(); err != nil {
			a.logger.Error("Error closing services", zap.Error(err))
		}
	}

	// Cancel context
	a.cancel()

	// Wait for context to be done
	<-ctx.Done()

	a.logger.Info("Application shutdown complete")
	return nil
} 
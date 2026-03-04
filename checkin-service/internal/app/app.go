package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"checkin-service/config"
	"checkin-service/database"
	"checkin-service/grpcclient"
	"checkin-service/repositories"
	"checkin-service/services"
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
	a.logger.Info("Initializing checkin-service components")

	// Connect to database
	db, err := database.NewConnection(a.config.Database)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	// Run migrations
	migrator := database.NewMigrator(db.DB.DB, a.logger)
	if err := migrator.RunMigrations("./migrations"); err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}

	// Connect to downstream gRPC services
	clients, err := grpcclient.NewClients(a.config)
	if err != nil {
		a.logger.Warn("gRPC clients unavailable at startup — will retry on first request", zap.Error(err))
	}

	// Wire up repositories and services
	checkinRepo := repositories.NewCheckinRepository(db)
	a.services = services.NewServices(checkinRepo, clients, a.logger)

	a.logger.Info("Checkin service initialized successfully")
	return nil
}

func (a *App) GetServices() *services.Services {
	return a.services
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down checkin-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if a.services != nil {
		a.services.Close()
	}
	a.cancel()
	<-ctx.Done()

	a.logger.Info("Shutdown complete")
	return nil
}

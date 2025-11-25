package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"ticket-service/config"
	"ticket-service/database"
	"ticket-service/grpc"
	"ticket-service/grpcclient"
	"ticket-service/metrics"
	"ticket-service/repositories"
	"ticket-service/services"
)

// App represents the main application
type App struct {
	logger             *zap.Logger
	config             *config.Config
	db                 *database.Connection
	ticketService      *services.TicketService
	bookingService     *services.TicketBookingSessionService
	reservationService *services.ReservationService
	eventClient        *grpcclient.EventServiceClient
	paymentClient      *grpcclient.PaymentServiceClient
	grpcServer         *grpc.Server
}

// NewApp creates a new application instance
func NewApp(logger *zap.Logger, cfg *config.Config) *App {
	return &App{
		logger: logger,
		config: cfg,
	}
}

// Initialize sets up all application components
func (a *App) Initialize() error {
	a.logger.Info("Initializing Ticket Service")

	// Initialize database
	dbConn, err := database.NewConnection(a.config.Database, a.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.db = dbConn

	// Run database migrations
	a.logger.Info("Running database migrations")
	migrator := database.NewMigrator(a.db.GetDB().DB, a.logger)
	if err := migrator.RunMigrations("./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	a.logger.Info("Database migrations completed")

	// Initialize repositories
	ticketRepo := repositories.NewTicketRepository(a.db.GetDB(), a.logger)
	bookingRepo := repositories.NewBookingSessionRepository(a.db.GetDB(), a.logger)
	reservationRepo := repositories.NewSeatReservationRepository(a.db.GetDB(), a.logger)

	// Initialize gRPC clients
	eventClient, err := grpcclient.NewEventServiceClient(a.config.Event, a.logger)
	if err != nil {
		a.logger.Warn("Failed to create Event Service client", zap.Error(err))
		eventClient = nil
	}
	a.eventClient = eventClient

	paymentClient, err := grpcclient.NewPaymentServiceClient(a.config.Payment, a.logger)
	if err != nil {
		a.logger.Warn("Failed to create Payment Service client", zap.Error(err))
		paymentClient = nil
	}
	a.paymentClient = paymentClient

	// Initialize services
	ticketService := services.NewTicketService(ticketRepo, eventClient, a.logger)
	bookingService := services.NewTicketBookingSessionService(bookingRepo, reservationRepo, eventClient, a.paymentClient, a.logger)
	reservationService := services.NewReservationService(reservationRepo, eventClient, a.logger)

	a.ticketService = ticketService
	a.bookingService = bookingService
	a.reservationService = reservationService

	// Initialize gRPC server
	grpcServer := grpc.NewServer(a.ticketService, a.bookingService, a.reservationService, a.logger)
	a.grpcServer = grpcServer

	// Initialize Prometheus metrics
	metrics.Init()

	a.logger.Info("Ticket Service initialized successfully")
	return nil
}

// Run starts the application and waits for shutdown signal
func (a *App) Run() error {
	// Start gRPC server
	go func() {
		if err := a.grpcServer.Start(a.config.GRPC.Port); err != nil {
			a.logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	a.logger.Info("Ticket Service started successfully",
		zap.String("grpc_port", a.config.GRPC.Port),
	)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down Ticket Service")

	// Stop gRPC server
	if err := a.grpcServer.Stop(); err != nil {
		a.logger.Error("Error stopping gRPC server", zap.Error(err))
	}

	// Close gRPC clients
	if a.eventClient != nil {
		a.eventClient.Close()
	}
	if a.paymentClient != nil {
		a.paymentClient.Close()
	}

	// Close database
	if err := a.db.Close(); err != nil {
		a.logger.Error("Error closing database", zap.Error(err))
	}

	a.logger.Info("Ticket Service stopped")
	return nil
}

// GetTicketService returns the ticket service instance
func (a *App) GetTicketService() *services.TicketService {
	return a.ticketService
}

// GetBookingService returns the booking service instance
func (a *App) GetBookingService() *services.TicketBookingSessionService {
	return a.bookingService
}

// GetReservationService returns the reservation service instance
func (a *App) GetReservationService() *services.ReservationService {
	return a.reservationService
}

// GetLogger returns the logger instance
func (a *App) GetLogger() *zap.Logger {
	return a.logger
}

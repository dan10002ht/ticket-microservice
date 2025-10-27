package grpc

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	ticketpb "shared-lib/protos/ticket"
	"ticket-service/services"
)

// Server represents the gRPC server
type Server struct {
	server             *grpc.Server
	ticketService      *services.TicketService
	bookingService     *services.BookingService
	reservationService *services.ReservationService
	logger             *zap.Logger
}

// NewServer creates a new gRPC server
func NewServer(
	ticketService *services.TicketService,
	bookingService *services.BookingService,
	reservationService *services.ReservationService,
	logger *zap.Logger,
) *Server {
	// Configure gRPC server options
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    30 * time.Second,
			Timeout: 5 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	server := grpc.NewServer(opts...)

	return &Server{
		server:             server,
		ticketService:      ticketService,
		bookingService:     bookingService,
		reservationService: reservationService,
		logger:             logger,
	}
}

// Start starts the gRPC server
func (s *Server) Start(port string) error {
	// Register services
	s.registerServices()

	// Enable reflection for debugging
	reflection.Register(s.server)

	// Start listening
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	s.logger.Info("Starting gRPC server", zap.String("port", port))

	// Start serving
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop stops the gRPC server gracefully
func (s *Server) Stop() error {
	s.logger.Info("Stopping gRPC server...")

	// Graceful stop with timeout
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	// Wait for graceful stop or timeout
	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
	case <-time.After(30 * time.Second):
		s.logger.Warn("gRPC server stop timeout, forcing stop")
		s.server.Stop()
	}

	return nil
}

// registerServices registers all gRPC services
func (s *Server) registerServices() {
	// Register Ticket Service
	ticketController := NewTicketController(s.ticketService, s.logger)
	ticketpb.RegisterTicketServiceServer(s.server, ticketController)

	// Register Booking Service
	bookingController := NewBookingController(s.bookingService, s.logger)
	ticketpb.RegisterBookingServiceServer(s.server, bookingController)

	// Register Reservation Service
	reservationController := NewReservationController(s.reservationService, s.logger)
	ticketpb.RegisterReservationServiceServer(s.server, reservationController)

	s.logger.Info("gRPC services registered successfully")
}

package main

import (
	"fmt"
	"log"
	"net"

	"ticket-service/config"
	"ticket-service/database"
	"ticket-service/grpc"
	"ticket-service/grpcclient"
	"ticket-service/repositories"
	"ticket-service/services"

	"google.golang.org/grpc"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect DB
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	// Init repositories
	ticketRepo := repositories.NewTicketRepository(db)
	seatReservationRepo := repositories.NewSeatReservationRepository(db)
	bookingSessionRepo := repositories.NewBookingSessionRepository(db)

	// Init services
	ticketService := services.NewTicketService(ticketRepo)
	seatReservationService := services.NewSeatReservationService(seatReservationRepo)
	bookingSessionService := services.NewBookingSessionService(bookingSessionRepo)

	// Init gRPC client for event-service
	// TODO: Lấy address event-service từ config/env
	eventConn, err := grpc.Dial("event-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to event-service: %v", err)
	}
	defer eventConn.Close()
	eventClient := grpcclient.NewEventClient(eventConn)

	// Init gRPC server
	grpcServer := grpc.NewServer()

	// Register controllers
	grpc.RegisterTicketServiceServer(grpcServer, grpc.NewTicketController(ticketService))
	grpc.RegisterBookingServiceServer(grpcServer, grpc.NewBookingController(bookingSessionService))
	grpc.RegisterSeatServiceServer(grpcServer, grpc.NewSeatController(seatReservationService))

	// Start gRPC server
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Ticket Service gRPC server started on port %s\n", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
} 
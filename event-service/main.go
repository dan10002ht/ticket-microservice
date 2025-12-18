package main

import (
	"event-service/config"
	grpcapi "event-service/grpc"
	"event-service/internal/app"
	"net"

	eventpb "event-service/internal/protos/event"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		logger.Fatal("Failed to connect DB", zap.Error(err))
	}

	appInstance := app.NewApp(logger, db)

	// Event Service controllers
	eventController := grpcapi.NewEventController(appInstance.GetEventService())
	zoneController := grpcapi.NewZoneController(appInstance.GetEventSeatingZoneService())
	seatController := grpcapi.NewEventSeatController(appInstance.GetEventSeatService())
	pricingController := grpcapi.NewAdvancedPricingController(appInstance.GetPricingService())
	availabilityController := grpcapi.NewAvailabilityController(appInstance.GetAvailabilityService())

	grpcServer := grpc.NewServer()

	// Register all Event Service gRPC services
	eventpb.RegisterEventServiceServer(grpcServer, eventController)
	eventpb.RegisterEventSeatingZoneServiceServer(grpcServer, zoneController)
	eventpb.RegisterEventSeatServiceServer(grpcServer, seatController)
	eventpb.RegisterPricingServiceServer(grpcServer, pricingController)
	eventpb.RegisterAvailabilityServiceServer(grpcServer, availabilityController)

	logger.Info("Starting Event Service gRPC server", zap.String("port", cfg.GRPC.Port))
	ln, err := listenOn(cfg.GRPC.Port)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}
	if err := grpcServer.Serve(ln); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}

func listenOn(addr string) (ln net.Listener, err error) {
	ln, err = net.Listen("tcp", addr)
	return
}

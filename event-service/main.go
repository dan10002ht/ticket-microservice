package main

import (
	"database/sql"
	"event-service/config"
	grpcapi "event-service/grpc"
	"event-service/internal/app"
	"net"

	eventpb "shared-lib/protos/event"

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

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		logger.Fatal("Failed to connect DB", zap.Error(err))
	}

	appInstance := app.NewApp(logger, db)

	eventController := grpcapi.NewEventController(appInstance.GetEventService())
	zoneController := grpcapi.NewEventSeatingZoneController(appInstance.GetEventSeatingZoneService())
	seatController := grpcapi.NewEventSeatController(appInstance.GetEventSeatService())

	grpcServer := grpc.NewServer()
	eventpb.RegisterEventServiceServer(grpcServer, eventController)
	eventpb.RegisterEventSeatingZoneServiceServer(grpcServer, zoneController)
	eventpb.RegisterEventSeatServiceServer(grpcServer, seatController)

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

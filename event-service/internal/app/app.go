package app

import (
	"event-service/repositories"
	"event-service/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type App struct {
	logger                   *zap.Logger
	db                       *sqlx.DB
	eventService             *services.EventService
	pricingService           *services.PricingService
	availabilityService      *services.AvailabilityService
	scheduleService          *services.ScheduleService
	eventSeatingZoneService  *services.EventSeatingZoneService
	eventSeatService         *services.EventSeatService
}

func NewApp(logger *zap.Logger, db *sqlx.DB) *App {
	// Event repository and service
	eventRepo := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo)

	// Pricing repository and service
	pricingRepo := repositories.NewEventPricingRepository(db)
	pricingService := services.NewPricingService(pricingRepo)

	// Availability repository and service
	availabilityRepo := repositories.NewEventSeatAvailabilityRepository(db)
	availabilityService := services.NewAvailabilityService(availabilityRepo)

	// Schedule repository and service
	scheduleRepo := repositories.NewEventScheduleRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo)

	// Zone repository and service
	zoneRepo := repositories.NewEventSeatingZoneRepository(db)
	eventSeatingZoneService := services.NewEventSeatingZoneService(zoneRepo)

	// Seat repository and service
	seatRepo := repositories.NewEventSeatRepository(db)
	eventSeatService := services.NewEventSeatService(seatRepo)

	return &App{
		logger:                   logger,
		db:                       db,
		eventService:             eventService,
		pricingService:           pricingService,
		availabilityService:      availabilityService,
		scheduleService:          scheduleService,
		eventSeatingZoneService:  eventSeatingZoneService,
		eventSeatService:         eventSeatService,
	}
}

func (a *App) GetEventService() *services.EventService {
	return a.eventService
}
func (a *App) GetPricingService() *services.PricingService {
	return a.pricingService
}
func (a *App) GetAvailabilityService() *services.AvailabilityService {
	return a.availabilityService
}
func (a *App) GetScheduleService() *services.ScheduleService {
	return a.scheduleService
}
func (a *App) GetEventSeatingZoneService() *services.EventSeatingZoneService {
	return a.eventSeatingZoneService
}
func (a *App) GetEventSeatService() *services.EventSeatService {
	return a.eventSeatService
}
func (a *App) GetLogger() *zap.Logger {
	return a.logger
}
func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.logger.Info("Shutting down Event Service")
	if err := a.db.Close(); err != nil {
		a.logger.Error("Error closing database", zap.Error(err))
	}
	return nil
} 
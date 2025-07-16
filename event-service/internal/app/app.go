package app

import (
	"database/sql"
	"event-service/repositories"
	"event-service/services"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	logger         *zap.Logger
	db             *sql.DB
	eventService   *services.EventService
	pricingService *services.PricingService
	availabilityService *services.AvailabilityService
	scheduleService *services.ScheduleService
}

func NewApp(logger *zap.Logger, db *sql.DB) *App {
	eventRepo := repositories.NewEventRepository(db)
	pricingRepo := repositories.NewEventPricingRepository(db)
	availabilityRepo := repositories.NewEventSeatAvailabilityRepository(db)
	scheduleRepo := repositories.NewEventScheduleRepository(db)

	eventService := services.NewEventService(eventRepo)
	pricingService := services.NewPricingService(pricingRepo)
	availabilityService := services.NewAvailabilityService(availabilityRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)

	return &App{
		logger: logger,
		db: db,
		eventService: eventService,
		pricingService: pricingService,
		availabilityService: availabilityService,
		scheduleService: scheduleService,
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
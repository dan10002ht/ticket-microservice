package main

import (
	"log"

	"ticket-service/config"
	"ticket-service/internal/app"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting Ticket Service...")

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded successfully")

	// Create application instance
	appInstance := app.NewApp(logger, cfg)

	// Initialize all components
	if err := appInstance.Initialize(); err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Run the application
	if err := appInstance.Run(); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}

	logger.Info("Ticket Service stopped")
}

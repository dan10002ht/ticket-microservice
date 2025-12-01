package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"booking-worker/config"
	"booking-worker/internal/app"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting Booking Worker Service...")

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded successfully",
		zap.String("redis_host", cfg.Redis.Host),
		zap.Int("redis_port", cfg.Redis.Port),
		zap.String("grpc_port", cfg.GRPC.Port),
	)

	// Create application instance
	appInstance := app.NewApp(logger, cfg)

	// Initialize all components
	if err := appInstance.Initialize(); err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run the application in a goroutine
	go func() {
		if err := appInstance.Run(); err != nil {
			logger.Fatal("Application failed", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutdown signal received, stopping application...")

	// Graceful shutdown
	if err := appInstance.Shutdown(); err != nil {
		logger.Error("Error during shutdown", zap.Error(err))
	}

	logger.Info("Booking Worker Service stopped")
}


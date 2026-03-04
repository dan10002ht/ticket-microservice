package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"checkin-service/config"
	"checkin-service/grpcserver"
	"checkin-service/internal/app"
	"checkin-service/internal/logger"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Checkin Service")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	appInstance := app.NewApp(log, cfg)
	if err := appInstance.Initialize(); err != nil {
		log.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Start gRPC server
	grpcSrv := grpcserver.NewServer(appInstance.GetServices(), cfg, log)
	go func() {
		if err := grpcSrv.Start(cfg.Server.GRPCPort); err != nil {
			log.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf(":%d", cfg.Server.MetricsPort)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Error("Failed to start metrics server", zap.Error(err))
		}
	}()

	// Block until shutdown signal
	if err := appInstance.Run(); err != nil {
		log.Fatal("Application error", zap.Error(err))
	}
}

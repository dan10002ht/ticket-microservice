package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"realtime-service/config"
	grpcserver "realtime-service/internal/grpc"
	"realtime-service/internal/middleware"
	"realtime-service/internal/pubsub"
	"realtime-service/internal/service"
	"realtime-service/internal/websocket"
	"realtime-service/metrics"
	"realtime-service/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Realtime Service",
		zap.String("http_port", fmt.Sprintf("%d", cfg.HTTP.Port)),
		zap.String("grpc_port", fmt.Sprintf("%d", cfg.GRPC.Port)),
	)

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Connected to Redis", zap.String("addr", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)))

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Notification Service
	notificationService := service.NewNotificationService(hub)

	// Initialize Redis Pub/Sub Subscriber
	subscriber := pubsub.NewSubscriber(redisClient, hub)
	if err := subscriber.Start(); err != nil {
		logger.Fatal("Failed to start Redis subscriber", zap.Error(err))
	}

	// Initialize Auth Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Initialize WebSocket Server
	wsServer := websocket.NewServer(hub, authMiddleware, cfg)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/ws", wsServer.HandleWebSocket)

	// Health endpoints
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler(redisClient))

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		mux.Handle("/metrics", metrics.Handler())
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      mux,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// Start gRPC server in goroutine
	grpcServer := grpcserver.NewServer(cfg, notificationService)
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	// Start HTTP server in goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Start metrics server if separate port configured
	var metricsServer *http.Server
	if cfg.Metrics.Enabled && cfg.Metrics.Port != cfg.HTTP.Port {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", metrics.Handler())
		metricsServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
			Handler: metricsMux,
		}
		go func() {
			logger.Info("Starting metrics server", zap.Int("port", cfg.Metrics.Port))
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("Metrics server failed", zap.Error(err))
			}
		}()
	}

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop accepting new connections
	hub.Stop()

	// Stop Redis subscriber
	if err := subscriber.Stop(); err != nil {
		logger.Error("Failed to stop Redis subscriber", zap.Error(err))
	}

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Shutdown metrics server
	if metricsServer != nil {
		if err := metricsServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("Metrics server shutdown error", zap.Error(err))
		}
	}

	// Stop gRPC server
	grpcServer.Stop()

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		logger.Error("Failed to close Redis connection", zap.Error(err))
	}

	logger.Info("Server shutdown complete")
}

// healthHandler returns basic health status
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// readyHandler returns readiness status (includes dependency checks)
func readyHandler(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Check Redis connection
		if err := redisClient.Ping(ctx).Err(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready","reason":"redis connection failed"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	}
}

package metrics

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Exporter handles Prometheus metrics
type Exporter struct {
	logger *zap.Logger
	server *http.Server

	// Metrics
	QueueLength    prometheus.Gauge
	ItemsProcessed prometheus.Counter
	ProcessingDuration prometheus.Histogram
	Errors         prometheus.Counter
}

// NewExporter creates a new metrics exporter
func NewExporter(port string, logger *zap.Logger) (*Exporter, error) {
	exporter := &Exporter{
		logger: logger,
	}

	// Register metrics
	exporter.QueueLength = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "booking_worker_queue_length",
		Help: "Current length of the booking queue",
	})

	exporter.ItemsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "booking_worker_items_processed_total",
		Help: "Total number of queue items processed",
	})

	exporter.ProcessingDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "booking_worker_processing_duration_seconds",
		Help: "Duration of queue item processing in seconds",
	})

	exporter.Errors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "booking_worker_errors_total",
		Help: "Total number of processing errors",
	})

	// Register with Prometheus
	if err := prometheus.Register(exporter.QueueLength); err != nil {
		return nil, fmt.Errorf("failed to register queue_length metric: %w", err)
	}
	if err := prometheus.Register(exporter.ItemsProcessed); err != nil {
		return nil, fmt.Errorf("failed to register items_processed metric: %w", err)
	}
	if err := prometheus.Register(exporter.ProcessingDuration); err != nil {
		return nil, fmt.Errorf("failed to register processing_duration metric: %w", err)
	}
	if err := prometheus.Register(exporter.Errors); err != nil {
		return nil, fmt.Errorf("failed to register errors metric: %w", err)
	}

	// Setup HTTP server for metrics endpoint
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	exporter.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start metrics server in goroutine
	go func() {
		if err := exporter.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	logger.Info("Metrics exporter started", zap.String("port", port))

	return exporter, nil
}


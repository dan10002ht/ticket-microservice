package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Email processing metrics
	EmailJobsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_jobs_processed_total",
			Help: "Total number of email jobs processed",
		},
		[]string{"status", "job_type", "template"},
	)

	EmailJobProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "email_job_processing_duration_seconds",
			Help:    "Time spent processing email jobs",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"job_type", "template"},
	)

	EmailQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "email_queue_size",
			Help: "Current number of emails in queue",
		},
	)

	EmailDeliverySuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_delivery_success_total",
			Help: "Total number of successful email deliveries",
		},
		[]string{"provider", "template"},
	)

	EmailDeliveryFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_delivery_failure_total",
			Help: "Total number of failed email deliveries",
		},
		[]string{"provider", "template", "error_type"},
	)

	// HTTP metrics
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 3, 5, 7, 10},
		},
		[]string{"method", "endpoint"},
	)

	// gRPC metrics
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GrpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method"},
	)

	// Database metrics
	DbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"operation", "table"},
	)

	DbConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	// Cache metrics
	CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	// Business metrics
	BusinessEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_events_total",
			Help: "Total number of business events",
		},
		[]string{"event_type"},
	)
)

func Init() {
	// Register all metrics
	prometheus.MustRegister(EmailJobsProcessed)
	prometheus.MustRegister(EmailJobProcessingDuration)
	prometheus.MustRegister(EmailQueueSize)
	prometheus.MustRegister(EmailDeliverySuccess)
	prometheus.MustRegister(EmailDeliveryFailure)
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(GrpcRequestsTotal)
	prometheus.MustRegister(GrpcRequestDuration)
	prometheus.MustRegister(DbQueryDuration)
	prometheus.MustRegister(DbConnectionsActive)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	prometheus.MustRegister(BusinessEventsTotal)

	// Initialize auto metrics
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "email_worker_info",
		Help: "Information about the email worker",
	}, func() float64 {
		return 1
	})
} 
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// WebSocket connection metrics
	WebSocketConnectionsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "realtime_websocket_connections_total",
		Help: "Total number of active WebSocket connections",
	})

	WebSocketConnectionsAuthenticated = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "realtime_websocket_connections_authenticated",
		Help: "Number of authenticated WebSocket connections",
	})

	WebSocketConnectionsAnonymous = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "realtime_websocket_connections_anonymous",
		Help: "Number of anonymous WebSocket connections",
	})

	WebSocketConnectionsPerRoom = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "realtime_websocket_connections_per_room",
		Help: "Number of WebSocket connections per room",
	}, []string{"room"})

	// Message metrics
	MessagesSentTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_messages_sent_total",
		Help: "Total number of messages sent to clients",
	}, []string{"type"})

	MessagesReceivedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_messages_received_total",
		Help: "Total number of messages received from clients",
	}, []string{"type"})

	MessageBroadcastTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_message_broadcast_total",
		Help: "Total number of broadcast messages",
	}, []string{"room"})

	// gRPC metrics
	GRPCRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_grpc_requests_total",
		Help: "Total number of gRPC requests",
	}, []string{"method", "status"})

	GRPCRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "realtime_grpc_request_duration_seconds",
		Help:    "Duration of gRPC requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})

	// Notification delivery metrics
	NotificationsDelivered = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_notifications_delivered_total",
		Help: "Total number of notifications delivered",
	}, []string{"type", "delivered"})

	// Redis Pub/Sub metrics
	PubSubMessagesReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_pubsub_messages_received_total",
		Help: "Total number of messages received from Redis Pub/Sub",
	}, []string{"channel"})

	PubSubMessagesPublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_pubsub_messages_published_total",
		Help: "Total number of messages published to Redis Pub/Sub",
	}, []string{"channel"})

	// Error metrics
	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "realtime_errors_total",
		Help: "Total number of errors",
	}, []string{"type"})

	// Uptime
	startTime time.Time
)

func init() {
	startTime = time.Now()
}

// Handler returns the Prometheus metrics HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}

// GetUptimeSeconds returns the service uptime in seconds
func GetUptimeSeconds() int64 {
	return int64(time.Since(startTime).Seconds())
}

// RecordNotificationDelivery records a notification delivery attempt
func RecordNotificationDelivery(notificationType string, delivered bool) {
	NotificationsDelivered.WithLabelValues(notificationType, strconv.FormatBool(delivered)).Inc()
}

// RecordGRPCRequest records a gRPC request
func RecordGRPCRequest(method string, status string, duration time.Duration) {
	GRPCRequestsTotal.WithLabelValues(method, status).Inc()
	GRPCRequestDuration.WithLabelValues(method).Observe(duration.Seconds())
}

// RecordError records an error
func RecordError(errorType string) {
	ErrorsTotal.WithLabelValues(errorType).Inc()
}

// RecordMessageSent records a message sent to a client
func RecordMessageSent(msgType string) {
	MessagesSentTotal.WithLabelValues(msgType).Inc()
}

// RecordMessageReceived records a message received from a client
func RecordMessageReceived(msgType string) {
	MessagesReceivedTotal.WithLabelValues(msgType).Inc()
}

// RecordBroadcast records a broadcast message
func RecordBroadcast(room string) {
	MessageBroadcastTotal.WithLabelValues(room).Inc()
}

// UpdateConnectionCount updates the connection count metrics
func UpdateConnectionCount(total, authenticated, anonymous int64) {
	WebSocketConnectionsTotal.Set(float64(total))
	WebSocketConnectionsAuthenticated.Set(float64(authenticated))
	WebSocketConnectionsAnonymous.Set(float64(anonymous))
}

// UpdateRoomConnectionCount updates the connection count for a specific room
func UpdateRoomConnectionCount(room string, count int) {
	WebSocketConnectionsPerRoom.WithLabelValues(room).Set(float64(count))
}

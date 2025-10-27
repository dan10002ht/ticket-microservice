package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for Ticket Service
var (
	// Ticket metrics
	TicketsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tickets_created_total",
			Help: "Total number of tickets created",
		},
		[]string{"event_id", "ticket_type", "status"},
	)

	TicketsCancelled = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tickets_cancelled_total",
			Help: "Total number of tickets cancelled",
		},
		[]string{"event_id", "reason"},
	)

	TicketsRefunded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tickets_refunded_total",
			Help: "Total number of tickets refunded",
		},
		[]string{"event_id", "reason"},
	)

	// Booking metrics
	BookingSessionsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_sessions_created_total",
			Help: "Total number of booking sessions created",
		},
		[]string{"event_id", "status"},
	)

	BookingSessionsExpired = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_sessions_expired_total",
			Help: "Total number of booking sessions expired",
		},
		[]string{"event_id"},
	)

	BookingSessionsCompleted = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_sessions_completed_total",
			Help: "Total number of booking sessions completed",
		},
		[]string{"event_id"},
	)

	// Reservation metrics
	SeatReservationsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "seat_reservations_created_total",
			Help: "Total number of seat reservations created",
		},
		[]string{"event_id", "zone_id"},
	)

	SeatReservationsReleased = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "seat_reservations_released_total",
			Help: "Total number of seat reservations released",
		},
		[]string{"event_id", "reason"},
	)

	SeatReservationsConfirmed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "seat_reservations_confirmed_total",
			Help: "Total number of seat reservations confirmed",
		},
		[]string{"event_id", "zone_id"},
	)

	// Payment metrics
	PaymentsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_processed_total",
			Help: "Total number of payments processed",
		},
		[]string{"event_id", "status", "payment_method"},
	)

	PaymentsFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_failed_total",
			Help: "Total number of payments failed",
		},
		[]string{"event_id", "reason"},
	)

	// Duration metrics
	BookingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "booking_duration_seconds",
			Help:    "Duration of booking process in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_id", "status"},
	)

	PaymentDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_duration_seconds",
			Help:    "Duration of payment process in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_id", "payment_method"},
	)

	// Active metrics
	ActiveBookingSessions = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_booking_sessions",
			Help: "Number of active booking sessions",
		},
		[]string{"event_id"},
	)

	ActiveSeatReservations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_seat_reservations",
			Help: "Number of active seat reservations",
		},
		[]string{"event_id", "zone_id"},
	)

	// Error metrics
	GRPCErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_errors_total",
			Help: "Total number of gRPC errors",
		},
		[]string{"service", "method", "error_type"},
	)

	DatabaseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"operation", "table", "error_type"},
	)
)

// Init initializes the metrics
func Init() {
	// Metrics are automatically registered when created
	// This function can be used for any additional initialization
}

// IncrementTicketCreated increments the tickets created counter
func IncrementTicketCreated(eventID, ticketType, status string) {
	TicketsCreated.WithLabelValues(eventID, ticketType, status).Inc()
}

// IncrementTicketCancelled increments the tickets cancelled counter
func IncrementTicketCancelled(eventID, reason string) {
	TicketsCancelled.WithLabelValues(eventID, reason).Inc()
}

// IncrementTicketRefunded increments the tickets refunded counter
func IncrementTicketRefunded(eventID, reason string) {
	TicketsRefunded.WithLabelValues(eventID, reason).Inc()
}

// IncrementBookingSessionCreated increments the booking sessions created counter
func IncrementBookingSessionCreated(eventID, status string) {
	BookingSessionsCreated.WithLabelValues(eventID, status).Inc()
}

// IncrementBookingSessionExpired increments the booking sessions expired counter
func IncrementBookingSessionExpired(eventID string) {
	BookingSessionsExpired.WithLabelValues(eventID).Inc()
}

// IncrementBookingSessionCompleted increments the booking sessions completed counter
func IncrementBookingSessionCompleted(eventID string) {
	BookingSessionsCompleted.WithLabelValues(eventID).Inc()
}

// IncrementSeatReservationCreated increments the seat reservations created counter
func IncrementSeatReservationCreated(eventID, zoneID string) {
	SeatReservationsCreated.WithLabelValues(eventID, zoneID).Inc()
}

// IncrementSeatReservationReleased increments the seat reservations released counter
func IncrementSeatReservationReleased(eventID, reason string) {
	SeatReservationsReleased.WithLabelValues(eventID, reason).Inc()
}

// IncrementSeatReservationConfirmed increments the seat reservations confirmed counter
func IncrementSeatReservationConfirmed(eventID, zoneID string) {
	SeatReservationsConfirmed.WithLabelValues(eventID, zoneID).Inc()
}

// IncrementPaymentProcessed increments the payments processed counter
func IncrementPaymentProcessed(eventID, status, paymentMethod string) {
	PaymentsProcessed.WithLabelValues(eventID, status, paymentMethod).Inc()
}

// IncrementPaymentFailed increments the payments failed counter
func IncrementPaymentFailed(eventID, reason string) {
	PaymentsFailed.WithLabelValues(eventID, reason).Inc()
}

// ObserveBookingDuration records the booking duration
func ObserveBookingDuration(eventID, status string, duration float64) {
	BookingDuration.WithLabelValues(eventID, status).Observe(duration)
}

// ObservePaymentDuration records the payment duration
func ObservePaymentDuration(eventID, paymentMethod string, duration float64) {
	PaymentDuration.WithLabelValues(eventID, paymentMethod).Observe(duration)
}

// SetActiveBookingSessions sets the number of active booking sessions
func SetActiveBookingSessions(eventID string, count float64) {
	ActiveBookingSessions.WithLabelValues(eventID).Set(count)
}

// SetActiveSeatReservations sets the number of active seat reservations
func SetActiveSeatReservations(eventID, zoneID string, count float64) {
	ActiveSeatReservations.WithLabelValues(eventID, zoneID).Set(count)
}

// IncrementGRPCError increments the gRPC errors counter
func IncrementGRPCError(service, method, errorType string) {
	GRPCErrors.WithLabelValues(service, method, errorType).Inc()
}

// IncrementDatabaseError increments the database errors counter
func IncrementDatabaseError(operation, table, errorType string) {
	DatabaseErrors.WithLabelValues(operation, table, errorType).Inc()
}

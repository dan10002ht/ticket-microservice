package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// Config holds circuit breaker configuration
type Config struct {
	// Name of the circuit breaker (for logging)
	Name string

	// MaxRequests is the maximum number of requests allowed to pass through
	// when the circuit breaker is half-open
	MaxRequests uint32

	// Interval is the cyclic period of the closed state for the circuit breaker
	// to clear the internal counts. If 0, the circuit breaker doesn't clear internal counts during the closed state.
	Interval time.Duration

	// Timeout is the period of the open state, after which the state becomes half-open.
	// Default: 60 seconds
	Timeout time.Duration

	// ReadyToTrip is called whenever a request fails in the closed state.
	// If it returns true, the circuit breaker will trip to open state.
	// Default: trips when consecutive failures >= 5
	ConsecutiveFailuresToTrip uint32

	// OnStateChange is called when the circuit breaker state changes
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// CircuitBreaker wraps gobreaker with logging
type CircuitBreaker struct {
	breaker *gobreaker.CircuitBreaker
	logger  *zap.Logger
	name    string
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig(name string, logger *zap.Logger) Config {
	return Config{
		Name:                      name,
		MaxRequests:               3,           // Allow 3 requests in half-open state
		Interval:                  0,           // Don't reset counts in closed state
		Timeout:                   60 * time.Second, // Open state duration
		ConsecutiveFailuresToTrip: 5,           // Trip after 5 consecutive failures
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("Circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	}
}

// New creates a new circuit breaker with the given configuration
func New(cfg Config, logger *zap.Logger) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.ConsecutiveFailuresToTrip
		},
		OnStateChange: cfg.OnStateChange,
	}

	return &CircuitBreaker{
		breaker: gobreaker.NewCircuitBreaker(settings),
		logger:  logger,
		name:    cfg.Name,
	}
}

// Execute runs the given function if the circuit breaker is closed or half-open.
// If the circuit is open, it returns an error immediately.
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	result, err := cb.breaker.Execute(fn)
	if err != nil {
		cb.logger.Debug("Circuit breaker execution failed",
			zap.String("name", cb.name),
			zap.String("state", cb.breaker.State().String()),
			zap.Error(err),
		)
	}
	return result, err
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.breaker.State()
}

// IsOpen returns true if the circuit breaker is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.breaker.State() == gobreaker.StateOpen
}

// Counts returns the internal counts of the circuit breaker
func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	return cb.breaker.Counts()
}

// Name returns the name of the circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

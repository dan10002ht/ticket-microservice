package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the booking worker
type Config struct {
	Redis     RedisConfig
	GRPC      GRPCConfig
	Queue     QueueConfig
	Worker    WorkerConfig
	Metrics   MetricsConfig
	Logging   LoggingConfig
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// GRPCConfig holds gRPC server settings
type GRPCConfig struct {
	Port string
	// Client endpoints
	BookingServiceEndpoint  string
	RealtimeServiceEndpoint string
}

// QueueConfig holds queue settings
type QueueConfig struct {
	MaxSize        int
	TimeoutSeconds int
	CleanupInterval time.Duration
}

// WorkerConfig holds worker pool settings
type WorkerConfig struct {
	PoolSize      int
	MaxRetries    int
	RetryInterval time.Duration
}

// MetricsConfig holds metrics settings
type MetricsConfig struct {
	Enabled bool
	Port    string
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string
	Format string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DATABASE", 0),
		},
		GRPC: GRPCConfig{
			Port:                   getEnv("GRPC_PORT", "50059"),
			BookingServiceEndpoint: getEnv("GRPC_BOOKING_SERVICE_ENDPOINT", "localhost:50058"),
			RealtimeServiceEndpoint: getEnv("GRPC_REALTIME_SERVICE_ENDPOINT", "localhost:50057"),
		},
		Queue: QueueConfig{
			MaxSize:        getEnvAsInt("QUEUE_MAX_SIZE", 100000),
			TimeoutSeconds: getEnvAsInt("QUEUE_TIMEOUT_SECONDS", 900), // 15 minutes
			CleanupInterval: getEnvAsDuration("QUEUE_CLEANUP_INTERVAL", 30*time.Second),
		},
		Worker: WorkerConfig{
			PoolSize:      getEnvAsInt("WORKER_POOL_SIZE", 10),
			MaxRetries:    getEnvAsInt("WORKER_MAX_RETRIES", 3),
			RetryInterval: getEnvAsDuration("WORKER_RETRY_INTERVAL", 5*time.Second),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvAsBool("METRICS_ENABLED", true),
			Port:     getEnv("METRICS_PORT", "9091"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}



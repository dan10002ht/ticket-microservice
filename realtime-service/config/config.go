package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the realtime service
type Config struct {
	HTTP      HTTPConfig
	GRPC      GRPCConfig
	WebSocket WebSocketConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Metrics   MetricsConfig
	Logging   LoggingConfig
}

// HTTPConfig holds HTTP server settings
type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// GRPCConfig holds gRPC server settings
type GRPCConfig struct {
	Port int
}

// WebSocketConfig holds WebSocket settings
type WebSocketConfig struct {
	ReadBufferSize    int
	WriteBufferSize   int
	PingInterval      time.Duration
	PongWait          time.Duration
	WriteWait         time.Duration
	MaxMessageSize    int64
	EnableCompression bool
	AllowAnonymous    bool
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// JWTConfig holds JWT validation settings
type JWTConfig struct {
	Secret string
}

// MetricsConfig holds metrics settings
type MetricsConfig struct {
	Enabled bool
	Port    int
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	return &Config{
		HTTP: HTTPConfig{
			Port:         getEnvAsInt("HTTP_PORT", 3003),
			ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
		GRPC: GRPCConfig{
			Port: getEnvAsInt("GRPC_PORT", 50057),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:    getEnvAsInt("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize:   getEnvAsInt("WS_WRITE_BUFFER_SIZE", 1024),
			PingInterval:      getEnvAsDuration("WS_PING_INTERVAL", 30*time.Second),
			PongWait:          getEnvAsDuration("WS_PONG_WAIT", 60*time.Second),
			WriteWait:         getEnvAsDuration("WS_WRITE_WAIT", 10*time.Second),
			MaxMessageSize:    getEnvAsInt64("WS_MAX_MESSAGE_SIZE", 65536),
			EnableCompression: getEnvAsBool("WS_ENABLE_COMPRESSION", true),
			AllowAnonymous:    getEnvAsBool("WS_ALLOW_ANONYMOUS", true),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 100),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 10),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-jwt-secret-key-change-in-production"),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvAsBool("METRICS_ENABLED", true),
			Port:    getEnvAsInt("METRICS_PORT", 9057),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}, nil
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

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
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

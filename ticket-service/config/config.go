package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the ticket service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	GRPC     GRPCConfig
	Event    EventServiceConfig
	Payment  PaymentServiceConfig
	Logging  LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// GRPCConfig holds gRPC configuration
type GRPCConfig struct {
	Port                         string
	MaxReceiveMsgSize            int
	MaxSendMsgSize               int
	KeepAliveTime                time.Duration
	KeepAliveTimeout             time.Duration
	KeepAlivePermitWithoutStream bool
}

// EventServiceConfig holds Event Service gRPC configuration
type EventServiceConfig struct {
	Host string
	Port string
}

// PaymentServiceConfig holds Payment Service gRPC configuration
type PaymentServiceConfig struct {
	Host string
	Port string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "3003"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", "30s"),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", "120s"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getIntEnv("DB_PORT", 5432),
			Database:        getEnv("DB_NAME", "booking_system_ticket"),
			Username:        getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres_password"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", "5m"),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", "1m"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getIntEnv("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
			PoolSize: getIntEnv("REDIS_POOL_SIZE", 10),
		},
		GRPC: GRPCConfig{
			Port:                         getEnv("GRPC_PORT", "50053"),
			MaxReceiveMsgSize:            getIntEnv("GRPC_MAX_RECEIVE_MSG_SIZE", 4194304), // 4MB
			MaxSendMsgSize:               getIntEnv("GRPC_MAX_SEND_MSG_SIZE", 4194304),    // 4MB
			KeepAliveTime:                getDurationEnv("GRPC_KEEPALIVE_TIME", "30s"),
			KeepAliveTimeout:             getDurationEnv("GRPC_KEEPALIVE_TIMEOUT", "5s"),
			KeepAlivePermitWithoutStream: getBoolEnv("GRPC_KEEPALIVE_PERMIT_WITHOUT_STREAM", true),
		},
		Event: EventServiceConfig{
			Host: getEnv("EVENT_SERVICE_HOST", "localhost"),
			Port: getEnv("EVENT_SERVICE_PORT", "50051"),
		},
		Payment: PaymentServiceConfig{
			Host: getEnv("PAYMENT_SERVICE_HOST", "localhost"),
			Port: getEnv("PAYMENT_SERVICE_PORT", "50054"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	GRPC     GRPCConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port     int
	GRPCPort int
	Timeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Timeout  time.Duration
}

type GRPCConfig struct {
	AuthServiceURL    string
	UserServiceURL    string
	BookingServiceURL string
	Timeout           time.Duration
}

type LoggingConfig struct {
	Level      string
	Format     string
	OutputPath string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:     getEnvAsInt("SERVER_PORT", 8080),
			GRPCPort: getEnvAsInt("GRPC_PORT", 50051),
			Timeout:  getEnvAsDuration("SERVER_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "boilerplate_db"),
			User:     getEnv("DB_USER", "boilerplate_user"),
			Password: getEnv("DB_PASSWORD", "boilerplate_password"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			MaxConns: getEnvAsInt("DB_MAX_CONNS", 10),
			Timeout:  getEnvAsDuration("DB_TIMEOUT", 5*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Timeout:  getEnvAsDuration("REDIS_TIMEOUT", 5*time.Second),
		},
		GRPC: GRPCConfig{
			AuthServiceURL:    getEnv("GRPC_AUTH_SERVICE_URL", "auth-service:50051"),
			UserServiceURL:    getEnv("GRPC_USER_SERVICE_URL", "user-service:50052"),
			BookingServiceURL: getEnv("GRPC_BOOKING_SERVICE_URL", "booking-service:50053"),
			Timeout:           getEnvAsDuration("GRPC_TIMEOUT", 10*time.Second),
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", ""),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
} 
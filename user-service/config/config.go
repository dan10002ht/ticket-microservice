package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the user service
type Config struct {
	Database DatabaseConfig
	GRPC     GRPCConfig
	Redis    RedisConfig
	Metrics  MetricsConfig
	Logging  LoggingConfig
}

// DatabaseConfig holds PostgreSQL connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

// GRPCConfig holds gRPC server settings
type GRPCConfig struct {
	Port string
}

// RedisConfig holds Redis connection settings (for caching)
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
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
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5433),
			User:     getEnv("DB_USER", "booking_user"),
			Password: getEnv("DB_PASSWORD", "booking_pass"),
			DBName:   getEnv("DB_NAME", "booking_system"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			MaxConns: getEnvAsInt("DB_MAX_CONNS", 25),
			MinConns: getEnvAsInt("DB_MIN_CONNS", 5),
		},
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50052"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DATABASE", 0),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvAsBool("METRICS_ENABLED", true),
			Port:    getEnv("METRICS_PORT", "9092"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

// GetDSN returns PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.DBName + "?sslmode=" + c.SSLMode
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

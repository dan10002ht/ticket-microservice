package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	GRPC          GRPCConfig
	TicketService TicketServiceConfig
	Logging       LoggingConfig
}

type ServerConfig struct {
	GRPCPort    int
	MetricsPort int
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type GRPCConfig struct {
	Timeout time.Duration
}

type TicketServiceConfig struct {
	Host string
	Port int
}

type LoggingConfig struct {
	Level  string
	Format string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			GRPCPort:    getEnvAsInt("GRPC_PORT", 50059),
			MetricsPort: getEnvAsInt("METRICS_PORT", 2112),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 50433),
			Name:            getEnv("DB_NAME", "booking_system"),
			User:            getEnv("DB_USER", "booking_user"),
			Password:        getEnv("DB_PASSWORD", "booking_pass"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		GRPC: GRPCConfig{
			Timeout: getEnvAsDuration("GRPC_TIMEOUT", 10*time.Second),
		},
		TicketService: TicketServiceConfig{
			Host: getEnv("TICKET_SERVICE_HOST", "ticket-service"),
			Port: getEnvAsInt("TICKET_SERVICE_PORT", 50054),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
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

package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type GRPCConfig struct {
	Port string
}

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	GRPC     GRPCConfig
	Env      string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			URL: getEnv("EVENT_DB_URL", "postgres://user:password@localhost:5432/event_db?sslmode=disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		GRPC: GRPCConfig{
			Port: getEnv("EVENT_GRPC_PORT", ":50052"),
		},
		Env: getEnv("ENV", "development"),
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var i int
		_, err := fmt.Sscanf(val, "%d", &i)
		if err == nil {
			return i
		}
	}
	return defaultVal
} 
package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	Env         string
}

func Load() (*Config, error) {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres_password@pgpool-ticket:5432/booking_system_ticket?sslmode=disable"),
		Port:        getEnv("PORT", ":50053"),
		Env:         getEnv("ENV", "development"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

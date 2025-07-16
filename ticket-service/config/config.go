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
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		Env:         os.Getenv("ENV"),
	}, nil
} 
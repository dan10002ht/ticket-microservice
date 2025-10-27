package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	connMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))
	connMaxIdleTime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_IDLE_TIME", "1m"))

	return &DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            port,
		Database:        getEnv("DB_NAME", "booking_system_ticket"),
		Username:        getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres_password"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
	}
}

// ConnectDatabase establishes database connection
func ConnectDatabase(config *DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.String("database", config.Database),
		zap.Int("max_open_conns", config.MaxOpenConns),
		zap.Int("max_idle_conns", config.MaxIdleConns),
	)

	return db, nil
}

// CloseDatabase closes database connection
func CloseDatabase(db *sqlx.DB, logger *zap.Logger) {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		} else {
			logger.Info("Database connection closed")
		}
	}
}

// HealthCheck performs database health check
func HealthCheck(db *sqlx.DB) error {
	var result int
	err := db.Get(&result, "SELECT 1")
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

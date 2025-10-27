package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Connection handles database connection and operations
type Connection struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewConnection creates a new database connection
func NewConnection(dsn string, logger *zap.Logger) (*Connection, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database")

	return &Connection{
		db:     db,
		logger: logger,
	}, nil
}

// GetDB returns the underlying sqlx.DB instance
func (c *Connection) GetDB() *sqlx.DB {
	return c.db
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			c.logger.Error("Failed to close database connection", zap.Error(err))
			return err
		}
		c.logger.Info("Database connection closed")
	}
	return nil
}

// HealthCheck performs a database health check
func (c *Connection) HealthCheck() error {
	var result int
	err := c.db.Get(&result, "SELECT 1")
	if err != nil {
		c.logger.Error("Database health check failed", zap.Error(err))
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// Stats returns database connection statistics
func (c *Connection) Stats() sql.DBStats {
	return c.db.Stats()
}

// Begin starts a new transaction
func (c *Connection) Begin() (*sqlx.Tx, error) {
	return c.db.Beginx()
}

// MustBegin starts a new transaction and panics on error
func (c *Connection) MustBegin() *sqlx.Tx {
	return c.db.MustBegin()
}

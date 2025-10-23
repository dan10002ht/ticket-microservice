package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"booking-system/email-worker/config"
)

// DB wraps sqlx.DB for convenience with PgPool-II support
type DB struct {
	Connection *sqlx.DB
}

// NewConnection creates a new database connection using sqlx with PgPool-II support
func NewConnection(cfg config.DatabaseConfig) (*DB, error) {
	var connection *sqlx.DB
	var err error

	// PgPool-II configuration - handles master/slave routing automatically
	if cfg.Host != "" {
		// Use PgPool-II endpoint
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
		)

		connection, err = sqlx.Connect("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PgPool-II: %w", err)
		}

		// Configure connection pool
		configureConnectionPool(connection, cfg)

		// Test the connection
		if err := connection.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping PgPool-II: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no database configuration provided")
	}

	return &DB{
		Connection: connection,
	}, nil
}

// configureConnectionPool configures the connection pool settings
func configureConnectionPool(db *sqlx.DB, cfg config.DatabaseConfig) {
	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}

	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 5
	}

	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = 5 * time.Minute
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}

// Close closes database connection
func (db *DB) Close() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}
	return nil
}

// GetDB returns the database connection - PgPool-II handles routing
func (db *DB) GetDB() *sqlx.DB {
	return db.Connection
}

// GetSQLDB returns the underlying sql.DB for compatibility
func (db *DB) GetSQLDB() *sql.DB {
	return db.Connection.DB
}

// QueryRowContext executes a query that returns a single row
// PgPool-II automatically routes SELECT queries to slaves
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.Connection.QueryRowContext(ctx, query, args...)
}

// GetContext executes a query that returns a single row
// PgPool-II automatically routes SELECT queries to slaves
func (db *DB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return db.Connection.GetContext(ctx, dest, query, args...)
}

// SelectContext executes a query that returns multiple rows
// PgPool-II automatically routes SELECT queries to slaves
func (db *DB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return db.Connection.SelectContext(ctx, dest, query, args...)
}

// ExecContext executes a query that doesn't return rows
// PgPool-II automatically routes write queries to master
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.Connection.ExecContext(ctx, query, args...)
}

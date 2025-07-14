package database

import (
	"database/sql"
	"fmt"
	"time"

	"venue-service/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewConnection(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxConns / 2)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(time.Minute * 30)
	
	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return &Database{db}, nil
}

func (db *Database) Ping() error {
	return db.DB.Ping()
}

func (db *Database) Close() error {
	return db.DB.Close()
} 
package repositories

import (
	"database/sql"

	"go.uber.org/zap"
)

type VenueLayoutRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewVenueLayoutRepository(db *sql.DB, logger *zap.Logger) *VenueLayoutRepository {
	return &VenueLayoutRepository{db: db, logger: logger}
}

package repositories

import (
	"database/sql"

	"go.uber.org/zap"
)

type SeatingZoneRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSeatingZoneRepository(db *sql.DB, logger *zap.Logger) *SeatingZoneRepository {
	return &SeatingZoneRepository{db: db, logger: logger}
}

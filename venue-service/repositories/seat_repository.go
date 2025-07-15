package repositories

import (
	"context"
	"database/sql"
	"venue-service/models"

	"go.uber.org/zap"
)

type SeatRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSeatRepository(db *sql.DB, logger *zap.Logger) *SeatRepository {
	return &SeatRepository{db: db, logger: logger}
}

func (r *SeatRepository) Create(ctx context.Context, seat *models.Seat) error {

}

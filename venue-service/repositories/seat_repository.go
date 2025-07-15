package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"venue-service/models"

	"github.com/google/uuid"
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
	query := `INSERT INTO seats (public_id, zone_id, seat_number, row_number, seat_type, coordinates, properties, is_active, is_available, display_order, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`
	return r.db.QueryRowContext(ctx, query, seat.PublicID, seat.ZoneID, seat.SeatNumber, seat.RowNumber, seat.SeatType, seat.Coordinates, seat.Properties, seat.IsActive, seat.IsAvailable, seat.DisplayOrder, seat.CreatedAt, seat.UpdatedAt).Scan(&seat.ID)
}

func (r *SeatRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.Seat, error) {
	query := `SELECT id, public_id, zone_id, seat_number, row_number, seat_type, coordinates, properties, is_active, is_available, display_order, created_at, updated_at FROM seats WHERE public_id = $1`
	var seat models.Seat
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&seat.ID, &seat.PublicID, &seat.ZoneID, &seat.SeatNumber, &seat.RowNumber, &seat.SeatType, &seat.Coordinates, &seat.Properties, &seat.IsActive, &seat.IsAvailable, &seat.DisplayOrder, &seat.CreatedAt, &seat.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepository) Update(ctx context.Context, seat *models.Seat) error {
	query := `UPDATE seats SET seat_number=$1, row_number=$2, seat_type=$3, coordinates=$4, properties=$5, is_active=$6, is_available=$7, display_order=$8, updated_at=$9 WHERE public_id=$10`
	_, err := r.db.ExecContext(ctx, query, seat.SeatNumber, seat.RowNumber, seat.SeatType, seat.Coordinates, seat.Properties, seat.IsActive, seat.IsAvailable, seat.DisplayOrder, seat.UpdatedAt, seat.PublicID)
	return err
}

func (r *SeatRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM seats WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *SeatRepository) ListByZoneID(ctx context.Context, zoneID int64) ([]*models.Seat, error) {
	query := `SELECT id, public_id, zone_id, seat_number, row_number, seat_type, coordinates, properties, is_active, is_available, display_order, created_at, updated_at FROM seats WHERE zone_id = $1`
	rows, err := r.db.QueryContext(ctx, query, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []*models.Seat
	for rows.Next() {
		var seat models.Seat
		err := rows.Scan(&seat.ID, &seat.PublicID, &seat.ZoneID, &seat.SeatNumber, &seat.RowNumber, &seat.SeatType, &seat.Coordinates, &seat.Properties, &seat.IsActive, &seat.IsAvailable, &seat.DisplayOrder, &seat.CreatedAt, &seat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		seats = append(seats, &seat)
	}
	return seats, nil
}

func (r *SeatRepository) BulkCreate(ctx context.Context, seats []*models.Seat) error {
	if len(seats) == 0 {
		return nil
	}
	query := `INSERT INTO seats (public_id, zone_id, seat_number, row_number, seat_type, coordinates, properties, is_active, is_available, display_order, created_at, updated_at) VALUES `
	vals := []interface{}{}
	placeholders := []string{}
	for i, seat := range seats {
		offset := i * 12
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8, offset+9, offset+10, offset+11, offset+12))
		vals = append(vals, seat.PublicID, seat.ZoneID, seat.SeatNumber, seat.RowNumber, seat.SeatType, seat.Coordinates, seat.Properties, seat.IsActive, seat.IsAvailable, seat.DisplayOrder, seat.CreatedAt, seat.UpdatedAt)
	}
	query += strings.Join(placeholders, ",")
	_, err := r.db.ExecContext(ctx, query, vals...)
	return err
}

package repositories

import (
	"context"
	"database/sql"
	"venue-service/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SeatingZoneRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSeatingZoneRepository(db *sql.DB, logger *zap.Logger) *SeatingZoneRepository {
	return &SeatingZoneRepository{db: db, logger: logger}
}

func (r *SeatingZoneRepository) Create(ctx context.Context, zone *models.SeatingZone) error {
	query := `INSERT INTO seating_zones (public_id, layout_id, name, description, zone_type, color, coordinates, seat_count, row_count, seats_per_row, pricing_category, is_active, display_order, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`
	return r.db.QueryRowContext(ctx, query, zone.PublicID, zone.LayoutID, zone.Name, zone.Description, zone.ZoneType, zone.Color, zone.Coordinates, zone.SeatCount, zone.RowCount, zone.SeatsPerRow, zone.PricingCategory, zone.IsActive, zone.DisplayOrder, zone.CreatedAt, zone.UpdatedAt).Scan(&zone.ID)
}

func (r *SeatingZoneRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.SeatingZone, error) {
	query := `SELECT id, public_id, layout_id, name, description, zone_type, color, coordinates, seat_count, row_count, seats_per_row, pricing_category, is_active, display_order, created_at, updated_at FROM seating_zones WHERE public_id = $1`
	var zone models.SeatingZone
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&zone.ID, &zone.PublicID, &zone.LayoutID, &zone.Name, &zone.Description, &zone.ZoneType, &zone.Color, &zone.Coordinates, &zone.SeatCount, &zone.RowCount, &zone.SeatsPerRow, &zone.PricingCategory, &zone.IsActive, &zone.DisplayOrder, &zone.CreatedAt, &zone.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *SeatingZoneRepository) Update(ctx context.Context, zone *models.SeatingZone) error {
	query := `UPDATE seating_zones SET name=$1, description=$2, zone_type=$3, color=$4, coordinates=$5, seat_count=$6, row_count=$7, seats_per_row=$8, pricing_category=$9, is_active=$10, display_order=$11, updated_at=$12 WHERE public_id=$13`
	_, err := r.db.ExecContext(ctx, query, zone.Name, zone.Description, zone.ZoneType, zone.Color, zone.Coordinates, zone.SeatCount, zone.RowCount, zone.SeatsPerRow, zone.PricingCategory, zone.IsActive, zone.DisplayOrder, zone.UpdatedAt, zone.PublicID)
	return err
}

func (r *SeatingZoneRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM seating_zones WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *SeatingZoneRepository) ListByLayoutID(ctx context.Context, layoutID int64) ([]*models.SeatingZone, error) {
	query := `SELECT id, public_id, layout_id, name, description, zone_type, color, coordinates, seat_count, row_count, seats_per_row, pricing_category, is_active, display_order, created_at, updated_at FROM seating_zones WHERE layout_id = $1`
	rows, err := r.db.QueryContext(ctx, query, layoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var zones []*models.SeatingZone
	for rows.Next() {
		var zone models.SeatingZone
		err := rows.Scan(&zone.ID, &zone.PublicID, &zone.LayoutID, &zone.Name, &zone.Description, &zone.ZoneType, &zone.Color, &zone.Coordinates, &zone.SeatCount, &zone.RowCount, &zone.SeatsPerRow, &zone.PricingCategory, &zone.IsActive, &zone.DisplayOrder, &zone.CreatedAt, &zone.UpdatedAt)
		if err != nil {
			return nil, err
		}
		zones = append(zones, &zone)
	}
	return zones, nil
}

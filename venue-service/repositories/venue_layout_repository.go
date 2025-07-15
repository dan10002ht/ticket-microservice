package repositories

import (
	"context"
	"database/sql"
	"venue-service/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VenueLayoutRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewVenueLayoutRepository(db *sql.DB, logger *zap.Logger) *VenueLayoutRepository {
	return &VenueLayoutRepository{db: db, logger: logger}
}

func (r *VenueLayoutRepository) Create(ctx context.Context, layout *models.VenueLayout) error {
	query := `INSERT INTO venue_layouts (public_id, venue_id, name, description, layout_type, canvas_config, seat_count, seating_config, is_active, is_default, version, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`
	return r.db.QueryRowContext(ctx, query, layout.PublicID, layout.VenueID, layout.Name, layout.Description, layout.LayoutType, layout.CanvasConfig, layout.SeatCount, layout.SeatingConfig, layout.IsActive, layout.IsDefault, layout.Version, layout.CreatedAt, layout.UpdatedAt).Scan(&layout.ID)
}

func (r *VenueLayoutRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.VenueLayout, error) {
	query := `SELECT id, public_id, venue_id, name, description, layout_type, canvas_config, seat_count, seating_config, is_active, is_default, version, created_at, updated_at FROM venue_layouts WHERE public_id = $1`
	var layout models.VenueLayout
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&layout.ID, &layout.PublicID, &layout.VenueID, &layout.Name, &layout.Description, &layout.LayoutType, &layout.CanvasConfig, &layout.SeatCount, &layout.SeatingConfig, &layout.IsActive, &layout.IsDefault, &layout.Version, &layout.CreatedAt, &layout.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &layout, nil
}

func (r *VenueLayoutRepository) Update(ctx context.Context, layout *models.VenueLayout) error {
	query := `UPDATE venue_layouts SET name=$1, description=$2, layout_type=$3, canvas_config=$4, seat_count=$5, seating_config=$6, is_active=$7, is_default=$8, version=$9, updated_at=$10 WHERE public_id=$11`
	_, err := r.db.ExecContext(ctx, query, layout.Name, layout.Description, layout.LayoutType, layout.CanvasConfig, layout.SeatCount, layout.SeatingConfig, layout.IsActive, layout.IsDefault, layout.Version, layout.UpdatedAt, layout.PublicID)
	return err
}

func (r *VenueLayoutRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM venue_layouts WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *VenueLayoutRepository) ListByVenueID(ctx context.Context, venueID int64) ([]*models.VenueLayout, error) {
	query := `SELECT id, public_id, venue_id, name, description, layout_type, canvas_config, seat_count, seating_config, is_active, is_default, version, created_at, updated_at FROM venue_layouts WHERE venue_id = $1`
	rows, err := r.db.QueryContext(ctx, query, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var layouts []*models.VenueLayout
	for rows.Next() {
		var layout models.VenueLayout
		err := rows.Scan(&layout.ID, &layout.PublicID, &layout.VenueID, &layout.Name, &layout.Description, &layout.LayoutType, &layout.CanvasConfig, &layout.SeatCount, &layout.SeatingConfig, &layout.IsActive, &layout.IsDefault, &layout.Version, &layout.CreatedAt, &layout.UpdatedAt)
		if err != nil {
			return nil, err
		}
		layouts = append(layouts, &layout)
	}
	return layouts, nil
}

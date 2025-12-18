package repositories

import (
	"context"
	"event-service/models"

	"github.com/jmoiron/sqlx"
)

type EventSeatingZoneRepository struct {
	db *sqlx.DB
}

func NewEventSeatingZoneRepository(db *sqlx.DB) *EventSeatingZoneRepository {
	return &EventSeatingZoneRepository{db: db}
}

func (r *EventSeatingZoneRepository) Create(ctx context.Context, zone *models.EventSeatingZone) error {
	query := `INSERT INTO event_seating_zones (public_id, event_id, name, zone_type, coordinates, seat_count, color, created_at, updated_at)
		VALUES (:public_id, :event_id, :name, :zone_type, :coordinates, :seat_count, :color, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, zone)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&zone.ID, &zone.CreatedAt, &zone.UpdatedAt)
	}
	return err
}

func (r *EventSeatingZoneRepository) GetByPublicID(ctx context.Context, publicID string) (*models.EventSeatingZone, error) {
	var zone models.EventSeatingZone
	query := `SELECT * FROM event_seating_zones WHERE public_id = $1`
	err := r.db.GetContext(ctx, &zone, query, publicID)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *EventSeatingZoneRepository) Update(ctx context.Context, zone *models.EventSeatingZone) error {
	query := `UPDATE event_seating_zones SET name=:name, zone_type=:zone_type, coordinates=:coordinates, seat_count=:seat_count, color=:color, updated_at=NOW() WHERE public_id=:public_id`
	_, err := r.db.NamedExecContext(ctx, query, zone)
	return err
}

func (r *EventSeatingZoneRepository) Delete(ctx context.Context, publicID string) error {
	query := `DELETE FROM event_seating_zones WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventSeatingZoneRepository) ListByEventID(ctx context.Context, eventID string) ([]*models.EventSeatingZone, error) {
	var zones []*models.EventSeatingZone
	query := `SELECT * FROM event_seating_zones WHERE event_id = $1 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &zones, query, eventID)
	return zones, err
} 
package repositories

import (
	"context"
	"event-service/models"

	"github.com/jmoiron/sqlx"
)

type EventSeatRepository struct {
	db *sqlx.DB
}

func NewEventSeatRepository(db *sqlx.DB) *EventSeatRepository {
	return &EventSeatRepository{db: db}
}

func (r *EventSeatRepository) Create(ctx context.Context, seat *models.EventSeat) error {
	query := `INSERT INTO event_seats (public_id, event_id, zone_id, seat_number, row_number, coordinates, created_at, updated_at)
		VALUES (:public_id, :event_id, :zone_id, :seat_number, :row_number, :coordinates, NOW(), NOW())`
	_, err := r.db.NamedExecContext(ctx, query, seat)
	return err
}

func (r *EventSeatRepository) GetByPublicID(ctx context.Context, publicID string) (*models.EventSeat, error) {
	var seat models.EventSeat
	query := `SELECT * FROM event_seats WHERE public_id = $1`
	err := r.db.GetContext(ctx, &seat, query, publicID)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *EventSeatRepository) Update(ctx context.Context, seat *models.EventSeat) error {
	query := `UPDATE event_seats SET seat_number=:seat_number, row_number=:row_number, coordinates=:coordinates, updated_at=NOW() WHERE public_id=:public_id`
	_, err := r.db.NamedExecContext(ctx, query, seat)
	return err
}

func (r *EventSeatRepository) Delete(ctx context.Context, publicID string) error {
	query := `DELETE FROM event_seats WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventSeatRepository) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSeat, error) {
	var seats []*models.EventSeat
	query := `SELECT * FROM event_seats WHERE event_id = $1 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &seats, query, eventID)
	return seats, err
} 
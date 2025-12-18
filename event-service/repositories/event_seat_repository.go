package repositories

import (
	"context"
	"event-service/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type EventSeatRepository struct {
	db *sqlx.DB
}

func NewEventSeatRepository(db *sqlx.DB) *EventSeatRepository {
	return &EventSeatRepository{db: db}
}

func (r *EventSeatRepository) Create(ctx context.Context, seat *models.EventSeat) error {
	query := `INSERT INTO event_seats (public_id, event_id, zone_id, seat_number, row_number, coordinates, status, pricing_category, base_price, final_price, currency, version, created_at, updated_at)
		VALUES (:public_id, :event_id, :zone_id, :seat_number, :row_number, :coordinates, :status, :pricing_category, :base_price, :final_price, :currency, :version, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, seat)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&seat.ID, &seat.CreatedAt, &seat.UpdatedAt)
	}
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
	query := `UPDATE event_seats SET seat_number=:seat_number, row_number=:row_number, coordinates=:coordinates,
		pricing_category=:pricing_category, base_price=:base_price, final_price=:final_price, currency=:currency,
		updated_at=NOW() WHERE public_id=:public_id RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, seat)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&seat.UpdatedAt)
	}
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

func (r *EventSeatRepository) ListSeatsByEvent(ctx context.Context, eventID, zoneID, status string, page, limit int32) ([]*models.EventSeat, int, error) {
	var seats []*models.EventSeat
	var total int

	// Build dynamic query
	baseQuery := `SELECT * FROM event_seats WHERE event_id = $1`
	countQuery := `SELECT COUNT(*) FROM event_seats WHERE event_id = $1`
	args := []interface{}{eventID}
	argIndex := 2

	if zoneID != "" {
		baseQuery += fmt.Sprintf(` AND zone_id = $%d`, argIndex)
		countQuery += fmt.Sprintf(` AND zone_id = $%d`, argIndex)
		args = append(args, zoneID)
		argIndex++
	}

	if status != "" {
		baseQuery += fmt.Sprintf(` AND status = $%d`, argIndex)
		countQuery += fmt.Sprintf(` AND status = $%d`, argIndex)
		args = append(args, status)
		argIndex++
	}

	// Get total count
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	baseQuery += fmt.Sprintf(` ORDER BY row_number, seat_number LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	err = r.db.SelectContext(ctx, &seats, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return seats, total, nil
} 
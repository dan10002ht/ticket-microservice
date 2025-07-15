package repositories

import (
	"context"
	"database/sql"
	"event-service/models"

	"github.com/google/uuid"
)

type EventSeatAvailabilityRepository struct {
	db *sql.DB
}

func NewEventSeatAvailabilityRepository(db *sql.DB) *EventSeatAvailabilityRepository {
	return &EventSeatAvailabilityRepository{db: db}
}

func (r *EventSeatAvailabilityRepository) Create(ctx context.Context, avail *models.EventSeatAvailability) error {
	query := `INSERT INTO event_seat_availability (public_id, event_id, seat_id, status, reserved_until, booking_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	return r.db.QueryRowContext(ctx, query, avail.PublicID, avail.EventID, avail.SeatID, avail.Status, avail.ReservedUntil, avail.BookingID, avail.CreatedAt, avail.UpdatedAt).Scan(&avail.ID)
}

func (r *EventSeatAvailabilityRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.EventSeatAvailability, error) {
	query := `SELECT id, public_id, event_id, seat_id, status, reserved_until, booking_id, created_at, updated_at FROM event_seat_availability WHERE public_id = $1`
	var avail models.EventSeatAvailability
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&avail.ID, &avail.PublicID, &avail.EventID, &avail.SeatID, &avail.Status, &avail.ReservedUntil, &avail.BookingID, &avail.CreatedAt, &avail.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &avail, nil
}

func (r *EventSeatAvailabilityRepository) Update(ctx context.Context, avail *models.EventSeatAvailability) error {
	query := `UPDATE event_seat_availability SET event_id=$1, seat_id=$2, status=$3, reserved_until=$4, booking_id=$5, updated_at=$6 WHERE public_id=$7`
	_, err := r.db.ExecContext(ctx, query, avail.EventID, avail.SeatID, avail.Status, avail.ReservedUntil, avail.BookingID, avail.UpdatedAt, avail.PublicID)
	return err
}

func (r *EventSeatAvailabilityRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM event_seat_availability WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventSeatAvailabilityRepository) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSeatAvailability, error) {
	query := `SELECT id, public_id, event_id, seat_id, status, reserved_until, booking_id, created_at, updated_at FROM event_seat_availability WHERE event_id = $1`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var avails []*models.EventSeatAvailability
	for rows.Next() {
		var avail models.EventSeatAvailability
		err := rows.Scan(&avail.ID, &avail.PublicID, &avail.EventID, &avail.SeatID, &avail.Status, &avail.ReservedUntil, &avail.BookingID, &avail.CreatedAt, &avail.UpdatedAt)
		if err != nil {
			return nil, err
		}
		avails = append(avails, &avail)
	}
	return avails, nil
} 
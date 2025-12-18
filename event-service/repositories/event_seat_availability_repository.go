package repositories

import (
	"context"
	"event-service/models"

	"github.com/jmoiron/sqlx"
)

type EventSeatAvailabilityRepository struct {
	db *sqlx.DB
}

func NewEventSeatAvailabilityRepository(db *sqlx.DB) *EventSeatAvailabilityRepository {
	return &EventSeatAvailabilityRepository{db: db}
}

func (r *EventSeatAvailabilityRepository) Create(ctx context.Context, avail *models.EventSeatAvailability) error {
	query := `INSERT INTO event_seat_availability (public_id, event_id, seat_id, zone_id, availability_status, reservation_id, blocked_reason, blocked_until, last_updated, created_at, updated_at)
		VALUES (:public_id, :event_id, :seat_id, :zone_id, :availability_status, :reservation_id, :blocked_reason, :blocked_until, NOW(), NOW(), NOW())
		RETURNING id, last_updated, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, avail)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&avail.ID, &avail.LastUpdated, &avail.CreatedAt, &avail.UpdatedAt)
	}
	return err
}

func (r *EventSeatAvailabilityRepository) GetByPublicID(ctx context.Context, publicID string) (*models.EventSeatAvailability, error) {
	var avail models.EventSeatAvailability
	query := `SELECT * FROM event_seat_availability WHERE public_id = $1`
	err := r.db.GetContext(ctx, &avail, query, publicID)
	if err != nil {
		return nil, err
	}
	return &avail, nil
}

func (r *EventSeatAvailabilityRepository) GetByEventID(ctx context.Context, eventID string) ([]*models.EventSeatAvailability, error) {
	var avails []*models.EventSeatAvailability
	query := `SELECT * FROM event_seat_availability WHERE event_id = $1 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &avails, query, eventID)
	return avails, err
}

func (r *EventSeatAvailabilityRepository) GetByEventAndZone(ctx context.Context, eventID, zoneID string) ([]*models.EventSeatAvailability, error) {
	var avails []*models.EventSeatAvailability
	query := `SELECT * FROM event_seat_availability WHERE event_id = $1 AND zone_id = $2 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &avails, query, eventID, zoneID)
	return avails, err
}

func (r *EventSeatAvailabilityRepository) GetBySeatID(ctx context.Context, eventID, seatID string) (*models.EventSeatAvailability, error) {
	var avail models.EventSeatAvailability
	query := `SELECT * FROM event_seat_availability WHERE event_id = $1 AND seat_id = $2`
	err := r.db.GetContext(ctx, &avail, query, eventID, seatID)
	if err != nil {
		return nil, err
	}
	return &avail, nil
}

func (r *EventSeatAvailabilityRepository) Update(ctx context.Context, avail *models.EventSeatAvailability) error {
	query := `UPDATE event_seat_availability SET availability_status=:availability_status, reservation_id=:reservation_id,
		blocked_reason=:blocked_reason, blocked_until=:blocked_until, last_updated=NOW(), updated_at=NOW()
		WHERE public_id=:public_id RETURNING last_updated, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, avail)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&avail.LastUpdated, &avail.UpdatedAt)
	}
	return err
}

func (r *EventSeatAvailabilityRepository) UpdateStatus(ctx context.Context, eventID, seatID, status, reservationID, blockedReason, blockedUntil string) error {
	query := `UPDATE event_seat_availability SET availability_status=$1, reservation_id=$2,
		blocked_reason=$3, blocked_until=$4, last_updated=NOW(), updated_at=NOW()
		WHERE event_id=$5 AND seat_id=$6`
	_, err := r.db.ExecContext(ctx, query, status, reservationID, blockedReason, blockedUntil, eventID, seatID)
	return err
}

func (r *EventSeatAvailabilityRepository) Delete(ctx context.Context, publicID string) error {
	query := `DELETE FROM event_seat_availability WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

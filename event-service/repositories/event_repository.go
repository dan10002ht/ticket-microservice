package repositories

import (
	"context"
	"event-service/models"

	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	query := `INSERT INTO events (public_id, organization_id, name, description, start_date, end_date, venue_name, venue_address, venue_city, venue_country, venue_capacity, canvas_config, created_at, updated_at)
		VALUES (:public_id, :organization_id, :name, :description, :start_date, :end_date, :venue_name, :venue_address, :venue_city, :venue_country, :venue_capacity, :canvas_config, NOW(), NOW())`
	_, err := r.db.NamedExecContext(ctx, query, event)
	return err
}

func (r *EventRepository) GetByPublicID(ctx context.Context, publicID string) (*models.Event, error) {
	var event models.Event
	query := `SELECT * FROM events WHERE public_id = $1`
	err := r.db.GetContext(ctx, &event, query, publicID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(ctx context.Context, event *models.Event) error {
	query := `UPDATE events SET name=:name, description=:description, start_date=:start_date, end_date=:end_date, venue_name=:venue_name, venue_address=:venue_address, venue_city=:venue_city, venue_country=:venue_country, venue_capacity=:venue_capacity, canvas_config=:canvas_config, updated_at=NOW() WHERE public_id=:public_id`
	_, err := r.db.NamedExecContext(ctx, query, event)
	return err
}

func (r *EventRepository) Delete(ctx context.Context, publicID string) error {
	query := `DELETE FROM events WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventRepository) List(ctx context.Context, organizationID int64) ([]*models.Event, error) {
	var events []*models.Event
	query := `SELECT * FROM events WHERE organization_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &events, query, organizationID)
	return events, err
} 
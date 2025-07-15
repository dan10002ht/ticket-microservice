package repositories

import (
	"context"
	"database/sql"
	"event-service/models"

	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	query := `INSERT INTO events (public_id, venue_id, layout_id, organizer_id, title, description, event_type, start_date, end_date, doors_open, status, max_capacity, current_capacity, images, tags, metadata, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id`
	return r.db.QueryRowContext(ctx, query, event.PublicID, event.VenueID, event.LayoutID, event.OrganizerID, event.Title, event.Description, event.EventType, event.StartDate, event.EndDate, event.DoorsOpen, event.Status, event.MaxCapacity, event.CurrentCapacity, event.Images, event.Tags, event.Metadata, event.IsActive, event.CreatedAt, event.UpdatedAt).Scan(&event.ID)
}

func (r *EventRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.Event, error) {
	query := `SELECT id, public_id, venue_id, layout_id, organizer_id, title, description, event_type, start_date, end_date, doors_open, status, max_capacity, current_capacity, images, tags, metadata, is_active, created_at, updated_at FROM events WHERE public_id = $1`
	var event models.Event
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&event.ID, &event.PublicID, &event.VenueID, &event.LayoutID, &event.OrganizerID, &event.Title, &event.Description, &event.EventType, &event.StartDate, &event.EndDate, &event.DoorsOpen, &event.Status, &event.MaxCapacity, &event.CurrentCapacity, &event.Images, &event.Tags, &event.Metadata, &event.IsActive, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(ctx context.Context, event *models.Event) error {
	query := `UPDATE events SET venue_id=$1, layout_id=$2, organizer_id=$3, title=$4, description=$5, event_type=$6, start_date=$7, end_date=$8, doors_open=$9, status=$10, max_capacity=$11, current_capacity=$12, images=$13, tags=$14, metadata=$15, is_active=$16, updated_at=$17 WHERE public_id=$18`
	_, err := r.db.ExecContext(ctx, query, event.VenueID, event.LayoutID, event.OrganizerID, event.Title, event.Description, event.EventType, event.StartDate, event.EndDate, event.DoorsOpen, event.Status, event.MaxCapacity, event.CurrentCapacity, event.Images, event.Tags, event.Metadata, event.IsActive, event.UpdatedAt, event.PublicID)
	return err
}

func (r *EventRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM events WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventRepository) List(ctx context.Context) ([]*models.Event, error) {
	query := `SELECT id, public_id, venue_id, layout_id, organizer_id, title, description, event_type, start_date, end_date, doors_open, status, max_capacity, current_capacity, images, tags, metadata, is_active, created_at, updated_at FROM events`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.PublicID, &event.VenueID, &event.LayoutID, &event.OrganizerID, &event.Title, &event.Description, &event.EventType, &event.StartDate, &event.EndDate, &event.DoorsOpen, &event.Status, &event.MaxCapacity, &event.CurrentCapacity, &event.Images, &event.Tags, &event.Metadata, &event.IsActive, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
} 
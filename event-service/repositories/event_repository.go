package repositories

import (
	"context"
	"event-service/models"
	"fmt"
	"strings"

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

func (r *EventRepository) ListByOrganizationID(ctx context.Context, organizationID string) ([]*models.Event, error) {
	var events []*models.Event
	query := `SELECT * FROM events WHERE organization_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &events, query, organizationID)
	return events, err
}

func (r *EventRepository) ListAll(ctx context.Context, page, limit int32) ([]*models.Event, error) {
	var events []*models.Event
	offset := (page - 1) * limit
	if page <= 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	query := `SELECT * FROM events ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &events, query, limit, offset)
	return events, err
}

// Advanced search and filtering methods
func (r *EventRepository) SearchEvents(ctx context.Context, query, eventType, category string, page, limit int32) ([]*models.Event, int, error) {
	var events []*models.Event
	var total int

	// Build search query
	searchQuery := `SELECT * FROM events WHERE 
		(name ILIKE $1 OR description ILIKE $1 OR venue_name ILIKE $1)`
	params := []interface{}{"%" + query + "%"}
	paramCount := 2

	if eventType != "" {
		searchQuery += ` AND event_type = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, eventType)
		paramCount++
	}

	if category != "" {
		searchQuery += ` AND category = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, category)
		paramCount++
	}

	// Count total
	countQuery := strings.Replace(searchQuery, "SELECT *", "SELECT COUNT(*)", 1)
	err := r.db.GetContext(ctx, &total, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	searchQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", paramCount) + ` OFFSET $` + fmt.Sprintf("%d", paramCount+1)
	params = append(params, limit, (page-1)*limit)

	err = r.db.SelectContext(ctx, &events, searchQuery, params...)
	return events, total, err
}

func (r *EventRepository) GetEventsByVenue(ctx context.Context, venueID, status string, page, limit int32) ([]*models.Event, int, error) {
	var events []*models.Event
	var total int

	query := `SELECT * FROM events WHERE venue_id = $1`
	params := []interface{}{venueID}
	paramCount := 2

	if status != "" {
		query += ` AND status = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, status)
		paramCount++
	}

	// Count total
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)
	err := r.db.GetContext(ctx, &total, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", paramCount) + ` OFFSET $` + fmt.Sprintf("%d", paramCount+1)
	params = append(params, limit, (page-1)*limit)

	err = r.db.SelectContext(ctx, &events, query, params...)
	return events, total, err
}

func (r *EventRepository) GetUpcomingEvents(ctx context.Context, days int32, eventType, category string, page, limit int32) ([]*models.Event, int, error) {
	var events []*models.Event
	var total int

	query := `SELECT * FROM events WHERE start_date > NOW() AND start_date <= NOW() + INTERVAL '%d days'`
	params := []interface{}{days}
	paramCount := 2

	if eventType != "" {
		query += ` AND event_type = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, eventType)
		paramCount++
	}

	if category != "" {
		query += ` AND category = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, category)
		paramCount++
	}

	// Count total
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)
	err := r.db.GetContext(ctx, &total, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += ` ORDER BY start_date ASC LIMIT $` + fmt.Sprintf("%d", paramCount) + ` OFFSET $` + fmt.Sprintf("%d", paramCount+1)
	params = append(params, limit, (page-1)*limit)

	err = r.db.SelectContext(ctx, &events, query, params...)
	return events, total, err
}

func (r *EventRepository) GetFeaturedEvents(ctx context.Context, eventType, category string, page, limit int32) ([]*models.Event, int, error) {
	var events []*models.Event
	var total int

	query := `SELECT * FROM events WHERE is_featured = true`
	params := []interface{}{}
	paramCount := 1

	if eventType != "" {
		query += ` AND event_type = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, eventType)
		paramCount++
	}

	if category != "" {
		query += ` AND category = $` + fmt.Sprintf("%d", paramCount)
		params = append(params, category)
		paramCount++
	}

	// Count total
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)
	err := r.db.GetContext(ctx, &total, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", paramCount) + ` OFFSET $` + fmt.Sprintf("%d", paramCount+1)
	params = append(params, limit, (page-1)*limit)

	err = r.db.SelectContext(ctx, &events, query, params...)
	return events, total, err
}

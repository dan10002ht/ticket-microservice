package repositories

import (
	"context"
	"event-service/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type EventPricingRepository struct {
	db *sqlx.DB
}

func NewEventPricingRepository(db *sqlx.DB) *EventPricingRepository {
	return &EventPricingRepository{db: db}
}

func (r *EventPricingRepository) Create(ctx context.Context, pricing *models.EventPricing) error {
	query := `INSERT INTO event_pricing (public_id, event_id, zone_id, pricing_category, base_price, currency, pricing_rules, discount_rules, is_active, valid_from, valid_until, created_by, created_at, updated_at)
		VALUES (:public_id, :event_id, :zone_id, :pricing_category, :base_price, :currency, :pricing_rules, :discount_rules, :is_active, :valid_from, :valid_until, :created_by, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, pricing)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&pricing.ID, &pricing.CreatedAt, &pricing.UpdatedAt)
	}
	return err
}

func (r *EventPricingRepository) GetByPublicID(ctx context.Context, publicID string) (*models.EventPricing, error) {
	var pricing models.EventPricing
	query := `SELECT * FROM event_pricing WHERE public_id = $1`
	err := r.db.GetContext(ctx, &pricing, query, publicID)
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *EventPricingRepository) Update(ctx context.Context, pricing *models.EventPricing) error {
	query := `UPDATE event_pricing SET base_price=:base_price, currency=:currency, pricing_rules=:pricing_rules,
		discount_rules=:discount_rules, is_active=:is_active, valid_from=:valid_from, valid_until=:valid_until,
		updated_by=:updated_by, updated_at=NOW() WHERE public_id=:public_id RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, pricing)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&pricing.UpdatedAt)
	}
	return err
}

func (r *EventPricingRepository) Delete(ctx context.Context, publicID string) error {
	query := `DELETE FROM event_pricing WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventPricingRepository) ListPricing(ctx context.Context, eventID string, isActive bool, page, limit int32) ([]*models.EventPricing, int, error) {
	var pricings []*models.EventPricing
	var total int

	countQuery := `SELECT COUNT(*) FROM event_pricing WHERE event_id = $1`
	baseQuery := `SELECT * FROM event_pricing WHERE event_id = $1`
	args := []interface{}{eventID}

	if isActive {
		countQuery += ` AND is_active = true`
		baseQuery += ` AND is_active = true`
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

	baseQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $2 OFFSET $3`)
	args = append(args, limit, offset)

	err = r.db.SelectContext(ctx, &pricings, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return pricings, total, nil
}

func (r *EventPricingRepository) GetPricingByEvent(ctx context.Context, eventID string, isActive bool) ([]*models.EventPricing, error) {
	var pricings []*models.EventPricing
	query := `SELECT * FROM event_pricing WHERE event_id = $1`
	if isActive {
		query += ` AND is_active = true`
	}
	query += ` ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &pricings, query, eventID)
	return pricings, err
}

func (r *EventPricingRepository) GetPricingByZone(ctx context.Context, eventID, zoneID string) ([]*models.EventPricing, error) {
	var pricings []*models.EventPricing
	query := `SELECT * FROM event_pricing WHERE event_id = $1 AND zone_id = $2 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &pricings, query, eventID, zoneID)
	return pricings, err
}

func (r *EventPricingRepository) GetActivePricingByZone(ctx context.Context, eventID, zoneID string) ([]*models.EventPricing, error) {
	var pricings []*models.EventPricing
	query := `SELECT * FROM event_pricing WHERE event_id = $1 AND zone_id = $2 AND is_active = true ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &pricings, query, eventID, zoneID)
	return pricings, err
}

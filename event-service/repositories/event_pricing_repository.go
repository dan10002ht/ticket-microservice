package repositories

import (
	"context"
	"database/sql"
	"event-service/models"

	"github.com/google/uuid"
)

type EventPricingRepository struct {
	db *sql.DB
}

func NewEventPricingRepository(db *sql.DB) *EventPricingRepository {
	return &EventPricingRepository{db: db}
}

func (r *EventPricingRepository) Create(ctx context.Context, pricing *models.EventPricing) error {
	query := `INSERT INTO event_pricing (public_id, event_id, zone_id, price, currency, pricing_type, pricing_rules, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`
	return r.db.QueryRowContext(ctx, query, pricing.PublicID, pricing.EventID, pricing.ZoneID, pricing.Price, pricing.Currency, pricing.PricingType, pricing.PricingRules, pricing.IsActive, pricing.CreatedAt, pricing.UpdatedAt).Scan(&pricing.ID)
}

func (r *EventPricingRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.EventPricing, error) {
	query := `SELECT id, public_id, event_id, zone_id, price, currency, pricing_type, pricing_rules, is_active, created_at, updated_at FROM event_pricing WHERE public_id = $1`
	var pricing models.EventPricing
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&pricing.ID, &pricing.PublicID, &pricing.EventID, &pricing.ZoneID, &pricing.Price, &pricing.Currency, &pricing.PricingType, &pricing.PricingRules, &pricing.IsActive, &pricing.CreatedAt, &pricing.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &pricing, nil
}

func (r *EventPricingRepository) Update(ctx context.Context, pricing *models.EventPricing) error {
	query := `UPDATE event_pricing SET event_id=$1, zone_id=$2, price=$3, currency=$4, pricing_type=$5, pricing_rules=$6, is_active=$7, updated_at=$8 WHERE public_id=$9`
	_, err := r.db.ExecContext(ctx, query, pricing.EventID, pricing.ZoneID, pricing.Price, pricing.Currency, pricing.PricingType, pricing.PricingRules, pricing.IsActive, pricing.UpdatedAt, pricing.PublicID)
	return err
}

func (r *EventPricingRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM event_pricing WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventPricingRepository) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventPricing, error) {
	query := `SELECT id, public_id, event_id, zone_id, price, currency, pricing_type, pricing_rules, is_active, created_at, updated_at FROM event_pricing WHERE event_id = $1`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pricings []*models.EventPricing
	for rows.Next() {
		var pricing models.EventPricing
		err := rows.Scan(&pricing.ID, &pricing.PublicID, &pricing.EventID, &pricing.ZoneID, &pricing.Price, &pricing.Currency, &pricing.PricingType, &pricing.PricingRules, &pricing.IsActive, &pricing.CreatedAt, &pricing.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pricings = append(pricings, &pricing)
	}
	return pricings, nil
} 
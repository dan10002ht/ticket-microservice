package models

type EventPricing struct {
	ID              int64   `db:"id" json:"-"`
	PublicID        string  `db:"public_id" json:"id"`
	EventID         string  `db:"event_id" json:"event_id"`
	ZoneID          string  `db:"zone_id" json:"zone_id"`
	PricingCategory string  `db:"pricing_category" json:"pricing_category"`
	BasePrice       float64 `db:"base_price" json:"base_price"`
	Currency        string  `db:"currency" json:"currency"`
	PricingRules    string  `db:"pricing_rules" json:"pricing_rules"`
	DiscountRules   string  `db:"discount_rules" json:"discount_rules"`
	IsActive        bool    `db:"is_active" json:"is_active"`
	ValidFrom       string  `db:"valid_from" json:"valid_from"`
	ValidUntil      string  `db:"valid_until" json:"valid_until"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
	CreatedBy       string  `db:"created_by" json:"created_by"`
	UpdatedBy       string  `db:"updated_by" json:"updated_by"`
}

package models

type EventPricing struct {
	ID            int64   `db:"id" json:"id"`
	PublicID      string  `db:"public_id" json:"public_id"`
	EventID       int64   `db:"event_id" json:"event_id"`
	ZoneID        int64   `db:"zone_id" json:"zone_id"`
	Price         float64 `db:"price" json:"price" validate:"required,gte=0"`
	Currency      string  `db:"currency" json:"currency"`
	PricingType   string  `db:"pricing_type" json:"pricing_type"`
	PricingRules  string  `db:"pricing_rules" json:"pricing_rules"`
	DiscountRules string  `db:"discount_rules" json:"discount_rules"`
	IsActive      bool    `db:"is_active" json:"is_active"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
	UpdatedAt     string  `db:"updated_at" json:"updated_at"`
} 
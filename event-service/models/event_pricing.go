package models

import (
	"time"

	"github.com/google/uuid"
)

type EventPricing struct {
	ID           int64     `db:"id" json:"-"`
	PublicID     uuid.UUID `db:"public_id" json:"id"`
	EventID      int64     `db:"event_id" json:"event_id"`
	ZoneID       int64     `db:"zone_id" json:"zone_id"`
	Price        float64   `db:"price" json:"price"`
	Currency     string    `db:"currency" json:"currency"`
	PricingType  string    `db:"pricing_type" json:"pricing_type"`
	PricingRules string    `db:"pricing_rules" json:"pricing_rules"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
} 
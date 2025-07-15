package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SeatingZone struct {
	ID              int64           `db:"id" json:"-"`
	PublicID        uuid.UUID       `db:"public_id" json:"id"`
	LayoutID        int64           `db:"layout_id" json:"layout_id"`
	Name            string          `db:"name" json:"name"`
	Description     string          `db:"description" json:"description"`
	ZoneType        string          `db:"zone_type" json:"zone_type"`
	Color           string          `db:"color" json:"color"`
	Coordinates     json.RawMessage `db:"coordinates" json:"coordinates"`
	SeatCount       int             `db:"seat_count" json:"seat_count"`
	RowCount        int             `db:"row_count" json:"row_count"`
	SeatsPerRow     int             `db:"seats_per_row" json:"seats_per_row"`
	PricingCategory string          `db:"pricing_category" json:"pricing_category"`
	IsActive        bool            `db:"is_active" json:"is_active"`
	DisplayOrder    int             `db:"display_order" json:"display_order"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

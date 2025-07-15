package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type VenueLayout struct {
	ID            int64           `db:"id" json:"-"`
	PublicID      uuid.UUID       `db:"public_id" json:"id"`
	VenueID       int64           `db:"venue_id" json:"venue_id"`
	Name          string          `db:"name" json:"name"`
	Description   string          `db:"description" json:"description"`
	LayoutType    string          `db:"layout_type" json:"layout_type"`
	CanvasConfig  json.RawMessage `db:"canvas_config" json:"canvas_config"`
	SeatCount     int             `db:"seat_count" json:"seat_count"`
	SeatingConfig json.RawMessage `db:"seating_config" json:"seating_config"`
	IsActive      bool            `db:"is_active" json:"is_active"`
	IsDefault     bool            `db:"is_default" json:"is_default"`
	Version       int64           `db:"version" json:"version"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updated_at"`
}

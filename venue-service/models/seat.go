package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Seat struct {
	ID           int64           `db:"id" json:"-"`
	PublicID     uuid.UUID       `db:"public_id" json:"id"`
	ZoneID       int64           `db:"zone_id" json:"zone_id"`
	SeatNumber   string          `db:"seat_number" json:"seat_number"`
	RowNumber    string          `db:"row_number" json:"row_number"`
	SeatType     string          `db:"seat_type" json:"seat_type"`
	Coordinates  json.RawMessage `db:"coordinates" json:"coordinates"`
	Properties   json.RawMessage `db:"properties" json:"properties"`
	IsActive     bool            `db:"is_active" json:"is_active"`
	IsAvailable  bool            `db:"is_available" json:"is_available"`
	DisplayOrder int             `db:"display_order" json:"display_order"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

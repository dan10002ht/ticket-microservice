package models

type EventSeat struct {
	ID              int64   `db:"id" json:"-"`
	PublicID        string  `db:"public_id" json:"id"`
	EventID         string  `db:"event_id" json:"event_id"`
	ZoneID          string  `db:"zone_id" json:"zone_id"`
	SeatNumber      string  `db:"seat_number" json:"seat_number"`
	RowNumber       string  `db:"row_number" json:"row_number"`
	Coordinates     string  `db:"coordinates" json:"coordinates"`
	Status          string  `db:"status" json:"status"`
	PricingCategory string  `db:"pricing_category" json:"pricing_category"`
	BasePrice       float64 `db:"base_price" json:"base_price"`
	FinalPrice      float64 `db:"final_price" json:"final_price"`
	Currency        string  `db:"currency" json:"currency"`
	Version         int     `db:"version" json:"version"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
} 
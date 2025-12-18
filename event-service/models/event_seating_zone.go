package models

type EventSeatingZone struct {
	ID          int64  `db:"id" json:"-"`
	PublicID    string `db:"public_id" json:"id"`
	EventID     string `db:"event_id" json:"event_id"`
	Name        string `db:"name" json:"name"`
	ZoneType    string `db:"zone_type" json:"zone_type"`
	Coordinates string `db:"coordinates" json:"coordinates"`
	SeatCount   int    `db:"seat_count" json:"seat_count"`
	Color       string `db:"color" json:"color"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
} 
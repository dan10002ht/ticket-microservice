package models

type EventSeat struct {
	ID         int64  `db:"id" json:"-"`
	PublicID   string `db:"public_id" json:"id"`
	EventID    int64  `db:"event_id" json:"event_id"`
	ZoneID     int64  `db:"zone_id" json:"zone_id"`
	SeatNumber string `db:"seat_number" json:"seat_number"`
	RowNumber  string `db:"row_number" json:"row_number"`
	Coordinates string `db:"coordinates" json:"coordinates"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	UpdatedAt  string `db:"updated_at" json:"updated_at"`
} 
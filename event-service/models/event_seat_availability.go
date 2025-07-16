package models

type EventSeatAvailability struct {
	ID         int64  `db:"id" json:"id"`
	PublicID   string `db:"public_id" json:"public_id"`
	EventID    int64  `db:"event_id" json:"event_id"`
	SeatID     int64  `db:"seat_id" json:"seat_id"`
	Status     string `db:"status" json:"status"`
	ReservedUntil string `db:"reserved_until" json:"reserved_until"`
	BookingID  int64  `db:"booking_id" json:"booking_id"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	UpdatedAt  string `db:"updated_at" json:"updated_at"`
} 
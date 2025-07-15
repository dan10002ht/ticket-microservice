package models

import (
	"time"

	"github.com/google/uuid"
)

type EventSeatAvailability struct {
	ID           int64     `db:"id" json:"-"`
	PublicID     uuid.UUID `db:"public_id" json:"id"`
	EventID      int64     `db:"event_id" json:"event_id"`
	SeatID       int64     `db:"seat_id" json:"seat_id"`
	Status       string    `db:"status" json:"status"`
	ReservedUntil *time.Time `db:"reserved_until" json:"reserved_until,omitempty"`
	BookingID    *int64    `db:"booking_id" json:"booking_id,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
} 
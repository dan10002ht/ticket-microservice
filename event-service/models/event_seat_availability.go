package models

type EventSeatAvailability struct {
	ID                 int64  `db:"id" json:"-"`
	PublicID           string `db:"public_id" json:"id"`
	EventID            string `db:"event_id" json:"event_id"`
	SeatID             string `db:"seat_id" json:"seat_id"`
	ZoneID             string `db:"zone_id" json:"zone_id"`
	AvailabilityStatus string `db:"availability_status" json:"availability_status"`
	ReservationID      string `db:"reservation_id" json:"reservation_id"`
	BlockedReason      string `db:"blocked_reason" json:"blocked_reason"`
	BlockedUntil       string `db:"blocked_until" json:"blocked_until"`
	LastUpdated        string `db:"last_updated" json:"last_updated"`
	CreatedAt          string `db:"created_at" json:"created_at"`
	UpdatedAt          string `db:"updated_at" json:"updated_at"`
}

// EventAvailabilitySummary - Summary of availability for an event
type EventAvailabilitySummary struct {
	TotalSeats     int32
	AvailableSeats int32
	ReservedSeats  int32
	BookedSeats    int32
	BlockedSeats   int32
}

// BlockSeatsResult - Result of blocking seats
type BlockSeatsResult struct {
	BlockedCount   int32
	BlockedSeatIDs []string
}

// ReleaseSeatsResult - Result of releasing seats
type ReleaseSeatsResult struct {
	ReleasedCount   int32
	ReleasedSeatIDs []string
}

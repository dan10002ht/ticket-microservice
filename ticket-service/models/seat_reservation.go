package models

type SeatReservation struct {
	ID                string  `db:"id" json:"id"`
	BookingSessionID  string  `db:"booking_session_id" json:"booking_session_id"`
	EventID           string  `db:"event_id" json:"event_id"`
	SeatID            string  `db:"seat_id" json:"seat_id"`
	ZoneID            string  `db:"zone_id" json:"zone_id"`
	ReservationToken  string  `db:"reservation_token" json:"reservation_token"`
	Status            string  `db:"status" json:"status"`
	ReservedAt        string  `db:"reserved_at" json:"reserved_at"`
	ExpiresAt         string  `db:"expires_at" json:"expires_at"`
	ReleasedAt        string  `db:"released_at" json:"released_at"`
	ReleasedReason    string  `db:"released_reason" json:"released_reason"`
	PricingCategory   string  `db:"pricing_category" json:"pricing_category"`
	BasePrice         float64 `db:"base_price" json:"base_price"`
	FinalPrice        float64 `db:"final_price" json:"final_price"`
	Currency          string  `db:"currency" json:"currency"`
	Metadata          string  `db:"metadata" json:"metadata"`
	CreatedAt         string  `db:"created_at" json:"created_at"`
	UpdatedAt         string  `db:"updated_at" json:"updated_at"`
	CreatedBy         string  `db:"created_by" json:"created_by"`
	UpdatedBy         string  `db:"updated_by" json:"updated_by"`
} 
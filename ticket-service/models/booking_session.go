package models

type BookingSession struct {
	ID              string  `db:"id" json:"id"`
	UserID          string  `db:"user_id" json:"user_id"`
	EventID         string  `db:"event_id" json:"event_id"`
	SessionToken    string  `db:"session_token" json:"session_token"`
	Status          string  `db:"status" json:"status"`
	SeatCount       int     `db:"seat_count" json:"seat_count"`
	TotalAmount     float64 `db:"total_amount" json:"total_amount"`
	Currency        string  `db:"currency" json:"currency"`
	ExpiresAt       string  `db:"expires_at" json:"expires_at"`
	CompletedAt     string  `db:"completed_at" json:"completed_at"`
	CancelledAt     string  `db:"cancelled_at" json:"cancelled_at"`
	CancelledReason string  `db:"cancelled_reason" json:"cancelled_reason"`
	IPAddress       string  `db:"ip_address" json:"ip_address"`
	UserAgent       string  `db:"user_agent" json:"user_agent"`
	Metadata        string  `db:"metadata" json:"metadata"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
	CreatedBy       string  `db:"created_by" json:"created_by"`
	UpdatedBy       string  `db:"updated_by" json:"updated_by"`
} 
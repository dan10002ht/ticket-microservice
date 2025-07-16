package models

type Ticket struct {
	ID                string  `db:"id" json:"id"`
	EventID           string  `db:"event_id" json:"event_id"`
	SeatID            string  `db:"seat_id" json:"seat_id"`
	ZoneID            string  `db:"zone_id" json:"zone_id"`
	UserID            string  `db:"user_id" json:"user_id"`
	BookingSessionID  string  `db:"booking_session_id" json:"booking_session_id"`
	TicketNumber      string  `db:"ticket_number" json:"ticket_number"`
	TicketType        string  `db:"ticket_type" json:"ticket_type"`
	PricingCategory   string  `db:"pricing_category" json:"pricing_category"`
	BasePrice         float64 `db:"base_price" json:"base_price"`
	FinalPrice        float64 `db:"final_price" json:"final_price"`
	Currency          string  `db:"currency" json:"currency"`
	DiscountAmount    float64 `db:"discount_amount" json:"discount_amount"`
	DiscountReason    string  `db:"discount_reason" json:"discount_reason"`
	Status            string  `db:"status" json:"status"`
	PaymentStatus     string  `db:"payment_status" json:"payment_status"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	PaymentReference  string  `db:"payment_reference" json:"payment_reference"`
	QRCode            string  `db:"qr_code" json:"qr_code"`
	Barcode           string  `db:"barcode" json:"barcode"`
	ValidFrom         string  `db:"valid_from" json:"valid_from"`
	ValidUntil        string  `db:"valid_until" json:"valid_until"`
	UsedAt            string  `db:"used_at" json:"used_at"`
	CancelledAt       string  `db:"cancelled_at" json:"cancelled_at"`
	CancelledReason   string  `db:"cancelled_reason" json:"cancelled_reason"`
	RefundedAt        string  `db:"refunded_at" json:"refunded_at"`
	RefundedAmount    float64 `db:"refunded_amount" json:"refunded_amount"`
	Metadata          string  `db:"metadata" json:"metadata"`
	CreatedAt         string  `db:"created_at" json:"created_at"`
	UpdatedAt         string  `db:"updated_at" json:"updated_at"`
	CreatedBy         string  `db:"created_by" json:"created_by"`
	UpdatedBy         string  `db:"updated_by" json:"updated_by"`
} 
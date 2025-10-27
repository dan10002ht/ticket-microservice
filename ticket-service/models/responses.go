package models

import "time"

// Response Models for API

// TicketResponse represents ticket data for API responses
type TicketResponse struct {
	ID               string     `json:"id"`
	EventID          string     `json:"event_id"`
	SeatID           string     `json:"seat_id"`
	ZoneID           string     `json:"zone_id"`
	UserID           string     `json:"user_id"`
	BookingSessionID *string    `json:"booking_session_id,omitempty"`
	TicketNumber     string     `json:"ticket_number"`
	TicketType       string     `json:"ticket_type"`
	PricingCategory  string     `json:"pricing_category"`
	BasePrice        float64    `json:"base_price"`
	FinalPrice       float64    `json:"final_price"`
	Currency         string     `json:"currency"`
	DiscountAmount   float64    `json:"discount_amount"`
	DiscountReason   *string    `json:"discount_reason,omitempty"`
	Status           string     `json:"status"`
	PaymentStatus    string     `json:"payment_status"`
	PaymentMethod    *string    `json:"payment_method,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty"`
	QRCode           *string    `json:"qr_code,omitempty"`
	Barcode          *string    `json:"barcode,omitempty"`
	ValidFrom        time.Time  `json:"valid_from"`
	ValidUntil       *time.Time `json:"valid_until,omitempty"`
	UsedAt           *time.Time `json:"used_at,omitempty"`
	CancelledAt      *time.Time `json:"cancelled_at,omitempty"`
	CancelledReason  *string    `json:"cancelled_reason,omitempty"`
	RefundedAt       *time.Time `json:"refunded_at,omitempty"`
	RefundedAmount   *float64   `json:"refunded_amount,omitempty"`
	Metadata         *string    `json:"metadata,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// BookingSessionResponse represents booking session data for API responses
type BookingSessionResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	EventID         string     `json:"event_id"`
	SessionToken    string     `json:"session_token"`
	Status          string     `json:"status"`
	SeatCount       int        `json:"seat_count"`
	TotalAmount     float64    `json:"total_amount"`
	Currency        string     `json:"currency"`
	ExpiresAt       time.Time  `json:"expires_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	CancelledReason *string    `json:"cancelled_reason,omitempty"`
	IPAddress       *string    `json:"ip_address,omitempty"`
	UserAgent       *string    `json:"user_agent,omitempty"`
	Metadata        *string    `json:"metadata,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	// Computed fields
	RemainingTime int64   `json:"remaining_time"` // seconds
	AveragePrice  float64 `json:"average_price"`
	IsActive      bool    `json:"is_active"`
	IsExpired     bool    `json:"is_expired"`
}

// SeatReservationResponse represents seat reservation data for API responses
type SeatReservationResponse struct {
	ID               string     `json:"id"`
	BookingSessionID string     `json:"booking_session_id"`
	EventID          string     `json:"event_id"`
	SeatID           string     `json:"seat_id"`
	ZoneID           string     `json:"zone_id"`
	ReservationToken string     `json:"reservation_token"`
	Status           string     `json:"status"`
	ReservedAt       time.Time  `json:"reserved_at"`
	ExpiresAt        time.Time  `json:"expires_at"`
	ReleasedAt       *time.Time `json:"released_at,omitempty"`
	ReleasedReason   *string    `json:"released_reason,omitempty"`
	PricingCategory  string     `json:"pricing_category"`
	BasePrice        float64    `json:"base_price"`
	FinalPrice       float64    `json:"final_price"`
	Currency         string     `json:"currency"`
	Metadata         *string    `json:"metadata,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Computed fields
	RemainingTime      int64   `json:"remaining_time"` // seconds
	DiscountAmount     float64 `json:"discount_amount"`
	DiscountPercentage float64 `json:"discount_percentage"`
	IsReserved         bool    `json:"is_reserved"`
	IsExpired          bool    `json:"is_expired"`
}

// Conversion Methods

// ToResponse converts Ticket to TicketResponse
func (t *Ticket) ToResponse() *TicketResponse {
	return &TicketResponse{
		ID:               t.ID,
		EventID:          t.EventID,
		SeatID:           t.SeatID,
		ZoneID:           t.ZoneID,
		UserID:           t.UserID,
		BookingSessionID: t.BookingSessionID,
		TicketNumber:     t.TicketNumber,
		TicketType:       t.TicketType,
		PricingCategory:  t.PricingCategory,
		BasePrice:        t.BasePrice,
		FinalPrice:       t.FinalPrice,
		Currency:         t.Currency,
		DiscountAmount:   t.DiscountAmount,
		DiscountReason:   t.DiscountReason,
		Status:           t.Status,
		PaymentStatus:    t.PaymentStatus,
		PaymentMethod:    t.PaymentMethod,
		PaymentReference: t.PaymentReference,
		QRCode:           t.QRCode,
		Barcode:          t.Barcode,
		ValidFrom:        t.ValidFrom,
		ValidUntil:       t.ValidUntil,
		UsedAt:           t.UsedAt,
		CancelledAt:      t.CancelledAt,
		CancelledReason:  t.CancelledReason,
		RefundedAt:       t.RefundedAt,
		RefundedAmount:   t.RefundedAmount,
		Metadata:         t.Metadata,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

// ToResponse converts BookingSession to BookingSessionResponse
func (bs *BookingSession) ToResponse() *BookingSessionResponse {
	return &BookingSessionResponse{
		ID:              bs.ID,
		UserID:          bs.UserID,
		EventID:         bs.EventID,
		SessionToken:    bs.SessionToken,
		Status:          bs.Status,
		SeatCount:       bs.SeatCount,
		TotalAmount:     bs.TotalAmount,
		Currency:        bs.Currency,
		ExpiresAt:       bs.ExpiresAt,
		CompletedAt:     bs.CompletedAt,
		CancelledAt:     bs.CancelledAt,
		CancelledReason: bs.CancelledReason,
		IPAddress:       bs.IPAddress,
		UserAgent:       bs.UserAgent,
		Metadata:        bs.Metadata,
		CreatedAt:       bs.CreatedAt,
		UpdatedAt:       bs.UpdatedAt,
		// Computed fields
		RemainingTime: int64(bs.GetRemainingTime().Seconds()),
		AveragePrice:  bs.CalculateAveragePrice(),
		IsActive:      bs.IsActive(),
		IsExpired:     bs.IsExpired(),
	}
}

// ToResponse converts SeatReservation to SeatReservationResponse
func (sr *SeatReservation) ToResponse() *SeatReservationResponse {
	return &SeatReservationResponse{
		ID:               sr.ID,
		BookingSessionID: sr.BookingSessionID,
		EventID:          sr.EventID,
		SeatID:           sr.SeatID,
		ZoneID:           sr.ZoneID,
		ReservationToken: sr.ReservationToken,
		Status:           sr.Status,
		ReservedAt:       sr.ReservedAt,
		ExpiresAt:        sr.ExpiresAt,
		ReleasedAt:       sr.ReleasedAt,
		ReleasedReason:   sr.ReleasedReason,
		PricingCategory:  sr.PricingCategory,
		BasePrice:        sr.BasePrice,
		FinalPrice:       sr.FinalPrice,
		Currency:         sr.Currency,
		Metadata:         sr.Metadata,
		CreatedAt:        sr.CreatedAt,
		UpdatedAt:        sr.UpdatedAt,
		// Computed fields
		RemainingTime:      int64(sr.GetRemainingTime().Seconds()),
		DiscountAmount:     sr.CalculateDiscount(),
		DiscountPercentage: sr.GetDiscountPercentage(),
		IsReserved:         sr.IsReserved(),
		IsExpired:          sr.IsExpired(),
	}
}

// Batch Response Models

// TicketListResponse represents a list of tickets
type TicketListResponse struct {
	Tickets []*TicketResponse `json:"tickets"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	HasMore bool              `json:"has_more"`
}

// BookingSessionListResponse represents a list of booking sessions
type BookingSessionListResponse struct {
	Sessions []*BookingSessionResponse `json:"sessions"`
	Total    int                       `json:"total"`
	Page     int                       `json:"page"`
	Limit    int                       `json:"limit"`
	HasMore  bool                      `json:"has_more"`
}

// SeatReservationListResponse represents a list of seat reservations
type SeatReservationListResponse struct {
	Reservations []*SeatReservationResponse `json:"reservations"`
	Total        int                        `json:"total"`
	Page         int                        `json:"page"`
	Limit        int                        `json:"limit"`
	HasMore      bool                       `json:"has_more"`
}

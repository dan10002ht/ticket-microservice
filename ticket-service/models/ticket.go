package models

import (
	"fmt"
	"time"
)

// Ticket represents a ticket in the system
type Ticket struct {
	ID               string     `json:"id" db:"id"`
	EventID          string     `json:"event_id" db:"event_id"`
	SeatID           string     `json:"seat_id" db:"seat_id"`
	ZoneID           string     `json:"zone_id" db:"zone_id"`
	UserID           string     `json:"user_id" db:"user_id"`
	BookingSessionID *string    `json:"booking_session_id" db:"booking_session_id"`
	TicketNumber     string     `json:"ticket_number" db:"ticket_number"`
	TicketType       string     `json:"ticket_type" db:"ticket_type"`
	PricingCategory  string     `json:"pricing_category" db:"pricing_category"`
	BasePrice        float64    `json:"base_price" db:"base_price"`
	FinalPrice       float64    `json:"final_price" db:"final_price"`
	Currency         string     `json:"currency" db:"currency"`
	DiscountAmount   float64    `json:"discount_amount" db:"discount_amount"`
	DiscountReason   *string    `json:"discount_reason" db:"discount_reason"`
	Status           string     `json:"status" db:"status"`
	PaymentStatus    string     `json:"payment_status" db:"payment_status"`
	PaymentMethod    *string    `json:"payment_method" db:"payment_method"`
	PaymentReference *string    `json:"payment_reference" db:"payment_reference"`
	QRCode           *string    `json:"qr_code" db:"qr_code"`
	Barcode          *string    `json:"barcode" db:"barcode"`
	ValidFrom        time.Time  `json:"valid_from" db:"valid_from"`
	ValidUntil       *time.Time `json:"valid_until" db:"valid_until"`
	UsedAt           *time.Time `json:"used_at" db:"used_at"`
	CancelledAt      *time.Time `json:"cancelled_at" db:"cancelled_at"`
	CancelledReason  *string    `json:"cancelled_reason" db:"cancelled_reason"`
	RefundedAt       *time.Time `json:"refunded_at" db:"refunded_at"`
	RefundedAmount   *float64   `json:"refunded_amount" db:"refunded_amount"`
	Metadata         *string    `json:"metadata" db:"metadata"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy        *string    `json:"created_by" db:"created_by"`
	UpdatedBy        *string    `json:"updated_by" db:"updated_by"`
}

// BookingSession represents a booking session
type BookingSession struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	EventID         string     `json:"event_id" db:"event_id"`
	SessionToken    string     `json:"session_token" db:"session_token"`
	Status          string     `json:"status" db:"status"`
	SeatCount       int        `json:"seat_count" db:"seat_count"`
	TotalAmount     float64    `json:"total_amount" db:"total_amount"`
	Currency        string     `json:"currency" db:"currency"`
	ExpiresAt       time.Time  `json:"expires_at" db:"expires_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
	CancelledAt     *time.Time `json:"cancelled_at" db:"cancelled_at"`
	CancelledReason *string    `json:"cancelled_reason" db:"cancelled_reason"`
	IPAddress       *string    `json:"ip_address" db:"ip_address"`
	UserAgent       *string    `json:"user_agent" db:"user_agent"`
	Metadata        *string    `json:"metadata" db:"metadata"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy       *string    `json:"created_by" db:"created_by"`
	UpdatedBy       *string    `json:"updated_by" db:"updated_by"`
}

// SeatReservation represents a temporary seat reservation
type SeatReservation struct {
	ID               string     `json:"id" db:"id"`
	BookingSessionID string     `json:"booking_session_id" db:"booking_session_id"`
	EventID          string     `json:"event_id" db:"event_id"`
	SeatID           string     `json:"seat_id" db:"seat_id"`
	ZoneID           string     `json:"zone_id" db:"zone_id"`
	ReservationToken string     `json:"reservation_token" db:"reservation_token"`
	Status           string     `json:"status" db:"status"`
	ReservedAt       time.Time  `json:"reserved_at" db:"reserved_at"`
	ExpiresAt        time.Time  `json:"expires_at" db:"expires_at"`
	ReleasedAt       *time.Time `json:"released_at" db:"released_at"`
	ReleasedReason   *string    `json:"released_reason" db:"released_reason"`
	PricingCategory  string     `json:"pricing_category" db:"pricing_category"`
	BasePrice        float64    `json:"base_price" db:"base_price"`
	FinalPrice       float64    `json:"final_price" db:"final_price"`
	Currency         string     `json:"currency" db:"currency"`
	Metadata         *string    `json:"metadata" db:"metadata"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy        *string    `json:"created_by" db:"created_by"`
	UpdatedBy        *string    `json:"updated_by" db:"updated_by"`
}

// Ticket Status Constants
const (
	TicketStatusPending   = "pending"
	TicketStatusConfirmed = "confirmed"
	TicketStatusCancelled = "cancelled"
	TicketStatusRefunded  = "refunded"
	TicketStatusUsed      = "used"
)

// Payment Status Constants
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

// Booking Session Status Constants
const (
	BookingSessionStatusActive    = "active"
	BookingSessionStatusCompleted = "completed"
	BookingSessionStatusExpired   = "expired"
	BookingSessionStatusCancelled = "cancelled"
)

// Reservation Status Constants
const (
	ReservationStatusReserved  = "reserved"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusReleased  = "released"
	ReservationStatusExpired   = "expired"
)

// Ticket Type Constants
const (
	TicketTypeStandard   = "standard"
	TicketTypeVIP        = "vip"
	TicketTypeWheelchair = "wheelchair"
	TicketTypeCompanion  = "companion"
)

// Pricing Category Constants
const (
	PricingCategoryPremium  = "premium"
	PricingCategoryStandard = "standard"
	PricingCategoryEconomy  = "economy"
	PricingCategoryVIP      = "vip"
)

// Validation Methods

// Validate validates ticket data
func (t *Ticket) Validate() error {
	if t.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if t.SeatID == "" {
		return fmt.Errorf("seat_id is required")
	}
	if t.ZoneID == "" {
		return fmt.Errorf("zone_id is required")
	}
	if t.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if t.TicketNumber == "" {
		return fmt.Errorf("ticket_number is required")
	}
	if !isValidTicketType(t.TicketType) {
		return fmt.Errorf("invalid ticket_type: %s", t.TicketType)
	}
	if !isValidPricingCategory(t.PricingCategory) {
		return fmt.Errorf("invalid pricing_category: %s", t.PricingCategory)
	}
	if !isValidTicketStatus(t.Status) {
		return fmt.Errorf("invalid status: %s", t.Status)
	}
	if !isValidPaymentStatus(t.PaymentStatus) {
		return fmt.Errorf("invalid payment_status: %s", t.PaymentStatus)
	}
	if t.BasePrice < 0 {
		return fmt.Errorf("base_price cannot be negative")
	}
	if t.FinalPrice < 0 {
		return fmt.Errorf("final_price cannot be negative")
	}
	if t.DiscountAmount < 0 {
		return fmt.Errorf("discount_amount cannot be negative")
	}
	if t.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}

// Validate validates booking session data
func (bs *BookingSession) Validate() error {
	if bs.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if bs.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if bs.SessionToken == "" {
		return fmt.Errorf("session_token is required")
	}
	if !isValidBookingSessionStatus(bs.Status) {
		return fmt.Errorf("invalid status: %s", bs.Status)
	}
	if bs.SeatCount <= 0 {
		return fmt.Errorf("seat_count must be positive")
	}
	if bs.TotalAmount < 0 {
		return fmt.Errorf("total_amount cannot be negative")
	}
	if bs.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if bs.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}
	if time.Now().After(bs.ExpiresAt) {
		return fmt.Errorf("expires_at must be in the future")
	}
	return nil
}

// Validate validates seat reservation data
func (sr *SeatReservation) Validate() error {
	if sr.BookingSessionID == "" {
		return fmt.Errorf("booking_session_id is required")
	}
	if sr.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if sr.SeatID == "" {
		return fmt.Errorf("seat_id is required")
	}
	if sr.ZoneID == "" {
		return fmt.Errorf("zone_id is required")
	}
	if sr.ReservationToken == "" {
		return fmt.Errorf("reservation_token is required")
	}
	if !isValidReservationStatus(sr.Status) {
		return fmt.Errorf("invalid status: %s", sr.Status)
	}
	if !isValidPricingCategory(sr.PricingCategory) {
		return fmt.Errorf("invalid pricing_category: %s", sr.PricingCategory)
	}
	if sr.BasePrice < 0 {
		return fmt.Errorf("base_price cannot be negative")
	}
	if sr.FinalPrice < 0 {
		return fmt.Errorf("final_price cannot be negative")
	}
	if sr.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if sr.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}
	if time.Now().After(sr.ExpiresAt) {
		return fmt.Errorf("expires_at must be in the future")
	}
	return nil
}

// Helper validation functions
func isValidTicketType(ticketType string) bool {
	validTypes := []string{TicketTypeStandard, TicketTypeVIP, TicketTypeWheelchair, TicketTypeCompanion}
	for _, validType := range validTypes {
		if ticketType == validType {
			return true
		}
	}
	return false
}

func isValidPricingCategory(category string) bool {
	validCategories := []string{PricingCategoryPremium, PricingCategoryStandard, PricingCategoryEconomy, PricingCategoryVIP}
	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}

func isValidTicketStatus(status string) bool {
	validStatuses := []string{TicketStatusPending, TicketStatusConfirmed, TicketStatusCancelled, TicketStatusRefunded, TicketStatusUsed}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func isValidPaymentStatus(status string) bool {
	validStatuses := []string{PaymentStatusPending, PaymentStatusPaid, PaymentStatusFailed, PaymentStatusRefunded}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func isValidBookingSessionStatus(status string) bool {
	validStatuses := []string{BookingSessionStatusActive, BookingSessionStatusCompleted, BookingSessionStatusExpired, BookingSessionStatusCancelled}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func isValidReservationStatus(status string) bool {
	validStatuses := []string{ReservationStatusReserved, ReservationStatusConfirmed, ReservationStatusReleased, ReservationStatusExpired}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// Business Logic Methods

// Ticket Methods
func (t *Ticket) IsConfirmed() bool {
	return t.Status == TicketStatusConfirmed
}

func (t *Ticket) IsCancelled() bool {
	return t.Status == TicketStatusCancelled
}

func (t *Ticket) IsRefunded() bool {
	return t.Status == TicketStatusRefunded
}

func (t *Ticket) IsUsed() bool {
	return t.Status == TicketStatusUsed
}

func (t *Ticket) IsPaid() bool {
	return t.PaymentStatus == PaymentStatusPaid
}

func (t *Ticket) CanBeCancelled() bool {
	return t.Status == TicketStatusPending || t.Status == TicketStatusConfirmed
}

func (t *Ticket) CanBeRefunded() bool {
	return t.Status == TicketStatusConfirmed && t.PaymentStatus == PaymentStatusPaid
}

func (t *Ticket) CalculateDiscount() float64 {
	return t.BasePrice - t.FinalPrice
}

func (t *Ticket) GetDiscountPercentage() float64 {
	if t.BasePrice == 0 {
		return 0
	}
	return (t.CalculateDiscount() / t.BasePrice) * 100
}

// BookingSession Methods
func (bs *BookingSession) IsActive() bool {
	return bs.Status == BookingSessionStatusActive && time.Now().Before(bs.ExpiresAt)
}

func (bs *BookingSession) IsExpired() bool {
	return time.Now().After(bs.ExpiresAt)
}

func (bs *BookingSession) IsCompleted() bool {
	return bs.Status == BookingSessionStatusCompleted
}

func (bs *BookingSession) IsCancelled() bool {
	return bs.Status == BookingSessionStatusCancelled
}

func (bs *BookingSession) GetRemainingTime() time.Duration {
	if bs.IsExpired() {
		return 0
	}
	return time.Until(bs.ExpiresAt)
}

func (bs *BookingSession) CalculateAveragePrice() float64 {
	if bs.SeatCount == 0 {
		return 0
	}
	return bs.TotalAmount / float64(bs.SeatCount)
}

// SeatReservation Methods
func (sr *SeatReservation) IsReserved() bool {
	return sr.Status == ReservationStatusReserved && time.Now().Before(sr.ExpiresAt)
}

func (sr *SeatReservation) IsExpired() bool {
	return time.Now().After(sr.ExpiresAt)
}

func (sr *SeatReservation) IsConfirmed() bool {
	return sr.Status == ReservationStatusConfirmed
}

func (sr *SeatReservation) IsReleased() bool {
	return sr.Status == ReservationStatusReleased
}

func (sr *SeatReservation) GetRemainingTime() time.Duration {
	if sr.IsExpired() {
		return 0
	}
	return time.Until(sr.ExpiresAt)
}

func (sr *SeatReservation) CalculateDiscount() float64 {
	return sr.BasePrice - sr.FinalPrice
}

func (sr *SeatReservation) GetDiscountPercentage() float64 {
	if sr.BasePrice == 0 {
		return 0
	}
	return (sr.CalculateDiscount() / sr.BasePrice) * 100
}

// Factory Methods

// NewTicket creates a new ticket with default values
func NewTicket(eventID, seatID, zoneID, userID, ticketNumber, pricingCategory string, basePrice, finalPrice float64, currency string) *Ticket {
	now := time.Now()
	return &Ticket{
		ID:              "", // Will be set by database
		EventID:         eventID,
		SeatID:          seatID,
		ZoneID:          zoneID,
		UserID:          userID,
		TicketNumber:    ticketNumber,
		TicketType:      TicketTypeStandard,
		PricingCategory: pricingCategory,
		BasePrice:       basePrice,
		FinalPrice:      finalPrice,
		Currency:        currency,
		DiscountAmount:  basePrice - finalPrice,
		Status:          TicketStatusPending,
		PaymentStatus:   PaymentStatusPending,
		ValidFrom:       now,
		ValidUntil:      nil, // Will be set based on event
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewBookingSession creates a new booking session
func NewBookingSession(userID, eventID, sessionToken string, seatCount int, totalAmount float64, currency string, expiresAt time.Time) *BookingSession {
	now := time.Now()
	return &BookingSession{
		ID:           "", // Will be set by database
		UserID:       userID,
		EventID:      eventID,
		SessionToken: sessionToken,
		Status:       BookingSessionStatusActive,
		SeatCount:    seatCount,
		TotalAmount:  totalAmount,
		Currency:     currency,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewSeatReservation creates a new seat reservation
func NewSeatReservation(bookingSessionID, eventID, seatID, zoneID, reservationToken, pricingCategory string, basePrice, finalPrice float64, currency string, expiresAt time.Time) *SeatReservation {
	now := time.Now()
	return &SeatReservation{
		ID:               "", // Will be set by database
		BookingSessionID: bookingSessionID,
		EventID:          eventID,
		SeatID:           seatID,
		ZoneID:           zoneID,
		ReservationToken: reservationToken,
		Status:           ReservationStatusReserved,
		ReservedAt:       now,
		ExpiresAt:        expiresAt,
		PricingCategory:  pricingCategory,
		BasePrice:        basePrice,
		FinalPrice:       finalPrice,
		Currency:         currency,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

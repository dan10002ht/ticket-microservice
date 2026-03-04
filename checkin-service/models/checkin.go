package models

import (
	"errors"
	"time"
)

// CheckIn statuses
const (
	CheckInStatusSuccess    = "success"
	CheckInStatusInvalid    = "invalid"
	CheckInStatusAlreadyUsed = "already_used"
	CheckInStatusCancelled  = "cancelled"
)

type CheckIn struct {
	ID          string     `db:"id"`
	TicketID    string     `db:"ticket_id"`
	EventID     string     `db:"event_id"`
	UserID      string     `db:"user_id"`
	StaffID     *string    `db:"staff_id"`
	QRCode      string     `db:"qr_code"`
	Status      string     `db:"status"`
	CheckInTime time.Time  `db:"check_in_time"`
	DeviceID    *string    `db:"device_id"`
	Gate        *string    `db:"gate"`
	Notes       *string    `db:"notes"`
	CreatedAt   time.Time  `db:"created_at"`
}

func (c *CheckIn) Validate() error {
	if c.TicketID == "" {
		return errors.New("ticket_id is required")
	}
	if c.EventID == "" {
		return errors.New("event_id is required")
	}
	if c.QRCode == "" {
		return errors.New("qr_code is required")
	}
	return nil
}

// EventStats holds aggregated check-in statistics for an event
type EventStats struct {
	EventID        string
	TotalCheckins  int
	UniqueTickets  int
	ByGate         map[string]int
	LastCheckinAt  *time.Time
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type EventSchedule struct {
	ID            int64     `db:"id" json:"-"`
	PublicID      uuid.UUID `db:"public_id" json:"id"`
	EventID       int64     `db:"event_id" json:"event_id"`
	ScheduleType  string    `db:"schedule_type" json:"schedule_type"`
	StartDate     time.Time `db:"start_date" json:"start_date"`
	EndDate       time.Time `db:"end_date" json:"end_date"`
	RecurrenceRule string   `db:"recurrence_rule" json:"recurrence_rule"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
} 
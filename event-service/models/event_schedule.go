package models

type EventSchedule struct {
	ID            int64  `db:"id" json:"id"`
	PublicID      string `db:"public_id" json:"public_id"`
	EventID       int64  `db:"event_id" json:"event_id"`
	ScheduleType  string `db:"schedule_type" json:"schedule_type"`
	StartDate     string `db:"start_date" json:"start_date"`
	EndDate       string `db:"end_date" json:"end_date"`
	RecurrenceRule string `db:"recurrence_rule" json:"recurrence_rule"`
	IsActive      bool   `db:"is_active" json:"is_active"`
	CreatedAt     string `db:"created_at" json:"created_at"`
	UpdatedAt     string `db:"updated_at" json:"updated_at"`
} 
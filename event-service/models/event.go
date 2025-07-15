package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID              int64       `db:"id" json:"-"`
	PublicID        uuid.UUID   `db:"public_id" json:"id"`
	VenueID         int64       `db:"venue_id" json:"venue_id"`
	LayoutID        int64       `db:"layout_id" json:"layout_id"`
	OrganizerID     int64       `db:"organizer_id" json:"organizer_id"`
	Title           string      `db:"title" json:"title"`
	Description     string      `db:"description" json:"description"`
	EventType       string      `db:"event_type" json:"event_type"`
	StartDate       time.Time   `db:"start_date" json:"start_date"`
	EndDate         time.Time   `db:"end_date" json:"end_date"`
	DoorsOpen       *time.Time  `db:"doors_open" json:"doors_open,omitempty"`
	Status          string      `db:"status" json:"status"`
	MaxCapacity     int         `db:"max_capacity" json:"max_capacity"`
	CurrentCapacity int         `db:"current_capacity" json:"current_capacity"`
	Images          []string    `db:"images" json:"images"`
	Tags            []string    `db:"tags" json:"tags"`
	Metadata        string      `db:"metadata" json:"metadata"`
	IsActive        bool        `db:"is_active" json:"is_active"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
} 
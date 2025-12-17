package model

import (
	"time"

	"github.com/google/uuid"
)

// Address represents a user address in the database
type Address struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Label      string    `json:"label"` // 'home', 'work', 'other'
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateAddressInput represents input for creating an address
type CreateAddressInput struct {
	UserID     uuid.UUID
	Label      string
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
	IsDefault  bool
}

// UpdateAddressInput represents input for updating an address
type UpdateAddressInput struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Label      *string
	Street     *string
	City       *string
	State      *string
	PostalCode *string
	Country    *string
	IsDefault  *bool
}

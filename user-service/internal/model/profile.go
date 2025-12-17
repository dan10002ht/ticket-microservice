package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Profile represents a user profile in the database
type Profile struct {
	ID          uuid.UUID         `json:"id"`
	UserID      uuid.UUID         `json:"user_id"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	AvatarURL   string            `json:"avatar_url"`
	DateOfBirth *time.Time        `json:"date_of_birth,omitempty"`
	Preferences map[string]string `json:"preferences"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// PreferencesJSON returns preferences as JSON string for database storage
func (p *Profile) PreferencesJSON() ([]byte, error) {
	if p.Preferences == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(p.Preferences)
}

// SetPreferencesFromJSON sets preferences from JSON bytes
func (p *Profile) SetPreferencesFromJSON(data []byte) error {
	if len(data) == 0 {
		p.Preferences = make(map[string]string)
		return nil
	}
	return json.Unmarshal(data, &p.Preferences)
}

// CreateProfileInput represents input for creating a profile
type CreateProfileInput struct {
	UserID      uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	AvatarURL   string
	DateOfBirth *time.Time
	Preferences map[string]string
}

// UpdateProfileInput represents input for updating a profile
type UpdateProfileInput struct {
	UserID      uuid.UUID
	FirstName   *string
	LastName    *string
	Phone       *string
	AvatarURL   *string
	DateOfBirth *time.Time
	Preferences map[string]string
}

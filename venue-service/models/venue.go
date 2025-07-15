package models

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type StringArray []string

type Venue struct {
	ID          int64       `db:"id" json:"-"`
	PublicID    uuid.UUID   `db:"public_id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description string      `db:"description" json:"description"`
	Address     string      `db:"address" json:"address"`
	City        string      `db:"city" json:"city"`
	State       string      `db:"state" json:"state"`
	Country     string      `db:"country" json:"country"`
	PostalCode  string      `db:"postal_code" json:"postal_code"`
	Phone       string      `db:"phone" json:"phone"`
	Email       string      `db:"email" json:"email"`
	Website     string      `db:"website" json:"website"`
	Capacity    int         `db:"capacity" json:"capacity"`
	VenueType   string      `db:"venue_type" json:"venue_type"`
	Amenities   StringArray `db:"amenities" json:"amenities"`
	Images      StringArray `db:"images" json:"images"`
	Coordinates StringArray `db:"coordinates" json:"coordinates"`
	Status      string      `db:"status" json:"status"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
}

// Value implements driver.Valuer for StringArray
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	// Convert to PostgreSQL array format: {"item1","item2","item3"}
	if len(s) == 0 {
		return "{}", nil
	}

	// Escape quotes and format as PostgreSQL array
	result := "{"
	for i, item := range s {
		if i > 0 {
			result += ","
		}
		// Escape double quotes by doubling them
		escaped := strings.ReplaceAll(item, `"`, `""`)
		result += `"` + escaped + `"`
	}
	result += "}"

	return result, nil
}

// Scan implements sql.Scanner for StringArray
func (s *StringArray) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Handle PostgreSQL array format: {"item1","item2","item3"}
		str := string(v)
		if str == "{}" {
			*s = StringArray{}
			return nil
		}

		// Remove outer braces
		if len(str) < 2 || str[0] != '{' || str[len(str)-1] != '}' {
			return errors.New("invalid PostgreSQL array format")
		}
		str = str[1 : len(str)-1]

		if str == "" {
			*s = StringArray{}
			return nil
		}

		// Parse array elements
		var result []string
		var current strings.Builder
		inQuotes := false
		escapeNext := false

		for i := 0; i < len(str); i++ {
			char := str[i]

			if escapeNext {
				current.WriteByte(char)
				escapeNext = false
				continue
			}

			if char == '\\' {
				escapeNext = true
				continue
			}

			if char == '"' {
				inQuotes = !inQuotes
				continue
			}

			if char == ',' && !inQuotes {
				result = append(result, current.String())
				current.Reset()
				continue
			}

			current.WriteByte(char)
		}

		// Add the last element
		if current.Len() > 0 || len(result) > 0 {
			result = append(result, current.String())
		}

		*s = StringArray(result)
		return nil

	case string:
		// Handle string representation
		return s.Scan([]byte(v))

	default:
		return errors.New("unsupported type for StringArray scan")
	}
}

func NewVenue(name, description, address, city, state, country, postalCode, phone, email, website string, capacity int, venueType string, amenities, images, coordinates []string) *Venue {
	return &Venue{
		PublicID:    uuid.New(),
		Name:        name,
		Description: description,
		Address:     address,
		City:        city,
		State:       state,
		Country:     country,
		PostalCode:  postalCode,
		Phone:       phone,
		Email:       email,
		Website:     website,
		Capacity:    capacity,
		VenueType:   venueType,
		Amenities:   amenities,
		Images:      images,
		Coordinates: coordinates,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

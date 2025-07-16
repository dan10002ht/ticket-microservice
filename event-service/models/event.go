package models

type Event struct {
	ID             int64   `db:"id" json:"-"`
	PublicID       string  `db:"public_id" json:"id"`
	OrganizationID int64   `db:"organization_id" json:"organization_id"`
	Name           string  `db:"name" json:"name"`
	Description    string  `db:"description" json:"description"`
	StartDate      string  `db:"start_date" json:"start_date"`
	EndDate        string  `db:"end_date" json:"end_date"`
	VenueName      string  `db:"venue_name" json:"venue_name"`
	VenueAddress   string  `db:"venue_address" json:"venue_address"`
	VenueCity      string  `db:"venue_city" json:"venue_city"`
	VenueCountry   string  `db:"venue_country" json:"venue_country"`
	VenueCapacity  int     `db:"venue_capacity" json:"venue_capacity"`
	CanvasConfig   string  `db:"canvas_config" json:"canvas_config"`
	CreatedAt      string  `db:"created_at" json:"created_at"`
	UpdatedAt      string  `db:"updated_at" json:"updated_at"`
} 
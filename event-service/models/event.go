package models

type Event struct {
	ID             int64  `db:"id" json:"-"`
	PublicID       string `db:"public_id" json:"id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	Name           string `db:"name" json:"name"`
	Description    string `db:"description" json:"description"`
	StartDate      string `db:"start_date" json:"start_date"`
	EndDate        string `db:"end_date" json:"end_date"`
	VenueName      string `db:"venue_name" json:"venue_name"`
	VenueAddress   string `db:"venue_address" json:"venue_address"`
	VenueCity      string `db:"venue_city" json:"venue_city"`
	VenueCountry   string `db:"venue_country" json:"venue_country"`
	VenueCapacity  int    `db:"venue_capacity" json:"venue_capacity"`
	CanvasConfig   string `db:"canvas_config" json:"canvas_config"`
	Status         string `db:"status" json:"status"`
	EventType      string `db:"event_type" json:"event_type"`
	Category       string `db:"category" json:"category"`
	SaleStartDate  string `db:"sale_start_date" json:"sale_start_date"`
	SaleEndDate    string `db:"sale_end_date" json:"sale_end_date"`
	MinAge         int    `db:"min_age" json:"min_age"`
	IsFeatured     bool   `db:"is_featured" json:"is_featured"`
	Images         string `db:"images" json:"images"`
	Tags           string `db:"tags" json:"tags"`
	Metadata       string `db:"metadata" json:"metadata"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
} 
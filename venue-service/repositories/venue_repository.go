package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"venue-service/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VenueRepository struct {
	db *sql.DB
	logger *zap.Logger
}

func NewVenueRepository(db *sql.DB, logger *zap.Logger) *VenueRepository {
	return &VenueRepository{db: db, logger: logger}
}

func (r *VenueRepository) Create(ctx context.Context, venue *models.Venue) error {
	query := `
		INSERT INTO venues (name, description, address, city, state, country, postal_code, phone, email, website, capacity, venue_type, amenities, images, coordinates, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		venue.Name, venue.Description, venue.Address, venue.City, venue.State, venue.Country, venue.PostalCode, venue.Phone, venue.Email, venue.Website, venue.Capacity, venue.VenueType, venue.Amenities, venue.Images, venue.Coordinates, venue.Status, venue.CreatedAt, venue.UpdatedAt,
	).Scan(&venue.ID)

	if err != nil {
		return fmt.Errorf("failed to create venue: %w", err)
	}
	r.logger.Info("Venue created", zap.String("venue_id", venue.PublicID.String()))
	return nil
}

func (r *VenueRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.Venue, error) {
	query := `
		SELECT id, public_id, name, description, address, city, state, country, postal_code, phone, email, website, capacity, venue_type, amenities, images, coordinates, status, created_at, updated_at
		FROM venues
		WHERE public_id = $1
	`

	var venue models.Venue
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(
		&venue.ID, &venue.PublicID, &venue.Name, &venue.Description, &venue.Address, &venue.City, &venue.State, &venue.Country, &venue.PostalCode, &venue.Phone, &venue.Email, &venue.Website, &venue.Capacity, &venue.VenueType, &venue.Amenities, &venue.Images, &venue.Coordinates, &venue.Status, &venue.CreatedAt, &venue.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("venue not found: %s", publicID)
		}
		return nil, fmt.Errorf("failed to get venue: %w", err)
	}
	return &venue, nil
}


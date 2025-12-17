package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"user-service/internal/model"
)

var (
	ErrAddressNotFound = errors.New("address not found")
)

// AddressRepository handles database operations for addresses
type AddressRepository struct {
	db *pgxpool.Pool
}

// NewAddressRepository creates a new AddressRepository
func NewAddressRepository(db *pgxpool.Pool) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(ctx context.Context, input model.CreateAddressInput) (*model.Address, error) {
	address := &model.Address{
		ID:         uuid.New(),
		UserID:     input.UserID,
		Label:      input.Label,
		Street:     input.Street,
		City:       input.City,
		State:      input.State,
		PostalCode: input.PostalCode,
		Country:    input.Country,
		IsDefault:  input.IsDefault,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// If this is set as default, unset other defaults
	if address.IsDefault {
		if err := r.unsetOtherDefaults(ctx, input.UserID, uuid.Nil); err != nil {
			return nil, err
		}
	}

	query := `
		INSERT INTO user_addresses (id, user_id, label, street, city, state, postal_code, country, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, label, street, city, state, postal_code, country, is_default, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		address.ID, address.UserID, address.Label, address.Street, address.City,
		address.State, address.PostalCode, address.Country, address.IsDefault,
		address.CreatedAt, address.UpdatedAt,
	).Scan(
		&address.ID, &address.UserID, &address.Label, &address.Street, &address.City,
		&address.State, &address.PostalCode, &address.Country, &address.IsDefault,
		&address.CreatedAt, &address.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return address, nil
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Address, error) {
	query := `
		SELECT id, user_id, label, street, city, state, postal_code, country, is_default, created_at, updated_at
		FROM user_addresses
		WHERE id = $1
	`

	address := &model.Address{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&address.ID, &address.UserID, &address.Label, &address.Street, &address.City,
		&address.State, &address.PostalCode, &address.Country, &address.IsDefault,
		&address.CreatedAt, &address.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}

	return address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Address, error) {
	query := `
		SELECT id, user_id, label, street, city, state, postal_code, country, is_default, created_at, updated_at
		FROM user_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*model.Address
	for rows.Next() {
		address := &model.Address{}
		err := rows.Scan(
			&address.ID, &address.UserID, &address.Label, &address.Street, &address.City,
			&address.State, &address.PostalCode, &address.Country, &address.IsDefault,
			&address.CreatedAt, &address.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Update updates an existing address
func (r *AddressRepository) Update(ctx context.Context, input model.UpdateAddressInput) (*model.Address, error) {
	// Get existing address first
	existing, err := r.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if existing.UserID != input.UserID {
		return nil, ErrAddressNotFound
	}

	// Apply updates
	if input.Label != nil {
		existing.Label = *input.Label
	}
	if input.Street != nil {
		existing.Street = *input.Street
	}
	if input.City != nil {
		existing.City = *input.City
	}
	if input.State != nil {
		existing.State = *input.State
	}
	if input.PostalCode != nil {
		existing.PostalCode = *input.PostalCode
	}
	if input.Country != nil {
		existing.Country = *input.Country
	}
	if input.IsDefault != nil {
		existing.IsDefault = *input.IsDefault
	}
	existing.UpdatedAt = time.Now()

	// If this is set as default, unset other defaults
	if existing.IsDefault {
		if err := r.unsetOtherDefaults(ctx, input.UserID, input.ID); err != nil {
			return nil, err
		}
	}

	query := `
		UPDATE user_addresses
		SET label = $1, street = $2, city = $3, state = $4, postal_code = $5, country = $6, is_default = $7, updated_at = $8
		WHERE id = $9 AND user_id = $10
		RETURNING id, user_id, label, street, city, state, postal_code, country, is_default, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query,
		existing.Label, existing.Street, existing.City, existing.State,
		existing.PostalCode, existing.Country, existing.IsDefault, existing.UpdatedAt,
		input.ID, input.UserID,
	).Scan(
		&existing.ID, &existing.UserID, &existing.Label, &existing.Street, &existing.City,
		&existing.State, &existing.PostalCode, &existing.Country, &existing.IsDefault,
		&existing.CreatedAt, &existing.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete deletes an address by ID
func (r *AddressRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM user_addresses WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrAddressNotFound
	}
	return nil
}

// unsetOtherDefaults unsets the default flag on all other addresses for a user
func (r *AddressRepository) unsetOtherDefaults(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID) error {
	query := `UPDATE user_addresses SET is_default = false WHERE user_id = $1 AND id != $2 AND is_default = true`
	_, err := r.db.Exec(ctx, query, userID, exceptID)
	return err
}

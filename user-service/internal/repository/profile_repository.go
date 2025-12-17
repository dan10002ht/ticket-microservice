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
	ErrProfileNotFound = errors.New("profile not found")
	ErrProfileExists   = errors.New("profile already exists for this user")
)

// ProfileRepository handles database operations for profiles
type ProfileRepository struct {
	db *pgxpool.Pool
}

// NewProfileRepository creates a new ProfileRepository
func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create creates a new profile
func (r *ProfileRepository) Create(ctx context.Context, input model.CreateProfileInput) (*model.Profile, error) {
	profile := &model.Profile{
		ID:          uuid.New(),
		UserID:      input.UserID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Phone:       input.Phone,
		AvatarURL:   input.AvatarURL,
		DateOfBirth: input.DateOfBirth,
		Preferences: input.Preferences,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	prefsJSON, err := profile.PreferencesJSON()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_profiles (id, user_id, first_name, last_name, email, phone, avatar_url, date_of_birth, preferences, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, first_name, last_name, email, phone, avatar_url, date_of_birth, preferences, created_at, updated_at
	`

	var prefsBytes []byte
	err = r.db.QueryRow(ctx, query,
		profile.ID, profile.UserID, profile.FirstName, profile.LastName, profile.Email,
		profile.Phone, profile.AvatarURL, profile.DateOfBirth, prefsJSON,
		profile.CreatedAt, profile.UpdatedAt,
	).Scan(
		&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Email,
		&profile.Phone, &profile.AvatarURL, &profile.DateOfBirth, &prefsBytes,
		&profile.CreatedAt, &profile.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := profile.SetPreferencesFromJSON(prefsBytes); err != nil {
		return nil, err
	}

	return profile, nil
}

// GetByUserID retrieves a profile by user ID
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, avatar_url, date_of_birth, preferences, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	profile := &model.Profile{}
	var prefsBytes []byte

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Email,
		&profile.Phone, &profile.AvatarURL, &profile.DateOfBirth, &prefsBytes,
		&profile.CreatedAt, &profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}

	if err := profile.SetPreferencesFromJSON(prefsBytes); err != nil {
		return nil, err
	}

	return profile, nil
}

// Update updates an existing profile
func (r *ProfileRepository) Update(ctx context.Context, input model.UpdateProfileInput) (*model.Profile, error) {
	// Get existing profile first
	existing, err := r.GetByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.FirstName != nil {
		existing.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		existing.LastName = *input.LastName
	}
	if input.Phone != nil {
		existing.Phone = *input.Phone
	}
	if input.AvatarURL != nil {
		existing.AvatarURL = *input.AvatarURL
	}
	if input.DateOfBirth != nil {
		existing.DateOfBirth = input.DateOfBirth
	}
	if input.Preferences != nil {
		existing.Preferences = input.Preferences
	}
	existing.UpdatedAt = time.Now()

	prefsJSON, err := existing.PreferencesJSON()
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE user_profiles
		SET first_name = $1, last_name = $2, phone = $3, avatar_url = $4, date_of_birth = $5, preferences = $6, updated_at = $7
		WHERE user_id = $8
		RETURNING id, user_id, first_name, last_name, email, phone, avatar_url, date_of_birth, preferences, created_at, updated_at
	`

	var prefsBytes []byte
	err = r.db.QueryRow(ctx, query,
		existing.FirstName, existing.LastName, existing.Phone, existing.AvatarURL,
		existing.DateOfBirth, prefsJSON, existing.UpdatedAt, input.UserID,
	).Scan(
		&existing.ID, &existing.UserID, &existing.FirstName, &existing.LastName, &existing.Email,
		&existing.Phone, &existing.AvatarURL, &existing.DateOfBirth, &prefsBytes,
		&existing.CreatedAt, &existing.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := existing.SetPreferencesFromJSON(prefsBytes); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete deletes a profile by user ID
func (r *ProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_profiles WHERE user_id = $1`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrProfileNotFound
	}
	return nil
}

// Exists checks if a profile exists for a user
func (r *ProfileRepository) Exists(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}

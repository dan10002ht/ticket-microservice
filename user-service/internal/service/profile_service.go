package service

import (
	"context"

	"github.com/google/uuid"

	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/pkg/logger"

	"go.uber.org/zap"
)

// ProfileService handles business logic for profiles
type ProfileService struct {
	repo *repository.ProfileRepository
}

// NewProfileService creates a new ProfileService
func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// CreateProfile creates a new user profile
func (s *ProfileService) CreateProfile(ctx context.Context, input model.CreateProfileInput) (*model.Profile, error) {
	// Check if profile already exists
	exists, err := s.repo.Exists(ctx, input.UserID)
	if err != nil {
		logger.Error("failed to check profile existence", zap.Error(err), zap.String("user_id", input.UserID.String()))
		return nil, err
	}
	if exists {
		return nil, repository.ErrProfileExists
	}

	profile, err := s.repo.Create(ctx, input)
	if err != nil {
		logger.Error("failed to create profile", zap.Error(err), zap.String("user_id", input.UserID.String()))
		return nil, err
	}

	logger.Info("profile created", zap.String("user_id", input.UserID.String()), zap.String("profile_id", profile.ID.String()))
	return profile, nil
}

// GetProfile retrieves a user profile by user ID
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if err == repository.ErrProfileNotFound {
			logger.Debug("profile not found", zap.String("user_id", userID.String()))
		} else {
			logger.Error("failed to get profile", zap.Error(err), zap.String("user_id", userID.String()))
		}
		return nil, err
	}
	return profile, nil
}

// UpdateProfile updates an existing user profile
func (s *ProfileService) UpdateProfile(ctx context.Context, input model.UpdateProfileInput) (*model.Profile, error) {
	profile, err := s.repo.Update(ctx, input)
	if err != nil {
		if err == repository.ErrProfileNotFound {
			logger.Debug("profile not found for update", zap.String("user_id", input.UserID.String()))
		} else {
			logger.Error("failed to update profile", zap.Error(err), zap.String("user_id", input.UserID.String()))
		}
		return nil, err
	}

	logger.Info("profile updated", zap.String("user_id", input.UserID.String()))
	return profile, nil
}

// DeleteProfile deletes a user profile
func (s *ProfileService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, userID)
	if err != nil {
		if err == repository.ErrProfileNotFound {
			logger.Debug("profile not found for delete", zap.String("user_id", userID.String()))
		} else {
			logger.Error("failed to delete profile", zap.Error(err), zap.String("user_id", userID.String()))
		}
		return err
	}

	logger.Info("profile deleted", zap.String("user_id", userID.String()))
	return nil
}

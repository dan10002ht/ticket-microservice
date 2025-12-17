package service

import (
	"context"

	"github.com/google/uuid"

	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/pkg/logger"

	"go.uber.org/zap"
)

// AddressService handles business logic for addresses
type AddressService struct {
	repo *repository.AddressRepository
}

// NewAddressService creates a new AddressService
func NewAddressService(repo *repository.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

// AddAddress creates a new address for a user
func (s *AddressService) AddAddress(ctx context.Context, input model.CreateAddressInput) (*model.Address, error) {
	address, err := s.repo.Create(ctx, input)
	if err != nil {
		logger.Error("failed to create address",
			zap.Error(err),
			zap.String("user_id", input.UserID.String()),
		)
		return nil, err
	}

	logger.Info("address created",
		zap.String("user_id", input.UserID.String()),
		zap.String("address_id", address.ID.String()),
	)
	return address, nil
}

// GetAddresses retrieves all addresses for a user
func (s *AddressService) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*model.Address, error) {
	addresses, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("failed to get addresses", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, err
	}
	return addresses, nil
}

// GetAddress retrieves a single address by ID
func (s *AddressService) GetAddress(ctx context.Context, addressID uuid.UUID) (*model.Address, error) {
	address, err := s.repo.GetByID(ctx, addressID)
	if err != nil {
		if err == repository.ErrAddressNotFound {
			logger.Debug("address not found", zap.String("address_id", addressID.String()))
		} else {
			logger.Error("failed to get address", zap.Error(err), zap.String("address_id", addressID.String()))
		}
		return nil, err
	}
	return address, nil
}

// UpdateAddress updates an existing address
func (s *AddressService) UpdateAddress(ctx context.Context, input model.UpdateAddressInput) (*model.Address, error) {
	address, err := s.repo.Update(ctx, input)
	if err != nil {
		if err == repository.ErrAddressNotFound {
			logger.Debug("address not found for update", zap.String("address_id", input.ID.String()))
		} else {
			logger.Error("failed to update address",
				zap.Error(err),
				zap.String("address_id", input.ID.String()),
				zap.String("user_id", input.UserID.String()),
			)
		}
		return nil, err
	}

	logger.Info("address updated",
		zap.String("address_id", input.ID.String()),
		zap.String("user_id", input.UserID.String()),
	)
	return address, nil
}

// DeleteAddress deletes an address
func (s *AddressService) DeleteAddress(ctx context.Context, addressID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, addressID, userID)
	if err != nil {
		if err == repository.ErrAddressNotFound {
			logger.Debug("address not found for delete", zap.String("address_id", addressID.String()))
		} else {
			logger.Error("failed to delete address",
				zap.Error(err),
				zap.String("address_id", addressID.String()),
				zap.String("user_id", userID.String()),
			)
		}
		return err
	}

	logger.Info("address deleted",
		zap.String("address_id", addressID.String()),
		zap.String("user_id", userID.String()),
	)
	return nil
}

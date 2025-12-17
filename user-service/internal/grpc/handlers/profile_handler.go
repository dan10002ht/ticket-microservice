package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"user-service/internal/model"
	pb "user-service/internal/protos"
	"user-service/internal/repository"
	"user-service/internal/service"
)

// ProfileHandler handles profile-related gRPC requests
type ProfileHandler struct {
	pb.UnimplementedUserServiceServer
	profileService *service.ProfileService
	addressService *service.AddressService
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(profileService *service.ProfileService, addressService *service.AddressService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		addressService: addressService,
	}
}

// GetProfile retrieves a user profile
func (h *ProfileHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	profile, err := h.profileService.GetProfile(ctx, userID)
	if err != nil {
		if err == repository.ErrProfileNotFound {
			return nil, status.Errorf(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get profile: %v", err)
	}

	return &pb.GetProfileResponse{
		Profile: modelToProtoProfile(profile),
	}, nil
}

// CreateProfile creates a new user profile
func (h *ProfileHandler) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	var dob *time.Time
	if req.GetDateOfBirth() != "" {
		parsed, err := time.Parse("2006-01-02", req.GetDateOfBirth())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_of_birth format (expected YYYY-MM-DD): %v", err)
		}
		dob = &parsed
	}

	input := model.CreateProfileInput{
		UserID:      userID,
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Email:       req.GetEmail(),
		Phone:       req.GetPhone(),
		AvatarURL:   req.GetAvatarUrl(),
		DateOfBirth: dob,
		Preferences: req.GetPreferences(),
	}

	profile, err := h.profileService.CreateProfile(ctx, input)
	if err != nil {
		if err == repository.ErrProfileExists {
			return nil, status.Errorf(codes.AlreadyExists, "profile already exists for this user")
		}
		return nil, status.Errorf(codes.Internal, "failed to create profile: %v", err)
	}

	return &pb.CreateProfileResponse{
		Profile: modelToProtoProfile(profile),
	}, nil
}

// UpdateProfile updates an existing user profile
func (h *ProfileHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	input := model.UpdateProfileInput{
		UserID:      userID,
		Preferences: req.GetPreferences(),
	}

	// Only set fields that are provided
	if req.GetFirstName() != "" {
		firstName := req.GetFirstName()
		input.FirstName = &firstName
	}
	if req.GetLastName() != "" {
		lastName := req.GetLastName()
		input.LastName = &lastName
	}
	if req.GetPhone() != "" {
		phone := req.GetPhone()
		input.Phone = &phone
	}
	if req.GetAvatarUrl() != "" {
		avatarURL := req.GetAvatarUrl()
		input.AvatarURL = &avatarURL
	}
	if req.GetDateOfBirth() != "" {
		parsed, err := time.Parse("2006-01-02", req.GetDateOfBirth())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_of_birth format (expected YYYY-MM-DD): %v", err)
		}
		input.DateOfBirth = &parsed
	}

	profile, err := h.profileService.UpdateProfile(ctx, input)
	if err != nil {
		if err == repository.ErrProfileNotFound {
			return nil, status.Errorf(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &pb.UpdateProfileResponse{
		Profile: modelToProtoProfile(profile),
	}, nil
}

// GetAddresses retrieves all addresses for a user
func (h *ProfileHandler) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	addresses, err := h.addressService.GetAddresses(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get addresses: %v", err)
	}

	protoAddresses := make([]*pb.Address, len(addresses))
	for i, addr := range addresses {
		protoAddresses[i] = modelToProtoAddress(addr)
	}

	return &pb.GetAddressesResponse{
		Addresses: protoAddresses,
	}, nil
}

// AddAddress adds a new address for a user
func (h *ProfileHandler) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	input := model.CreateAddressInput{
		UserID:     userID,
		Label:      req.GetLabel(),
		Street:     req.GetStreet(),
		City:       req.GetCity(),
		State:      req.GetState(),
		PostalCode: req.GetPostalCode(),
		Country:    req.GetCountry(),
		IsDefault:  req.GetIsDefault(),
	}

	address, err := h.addressService.AddAddress(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add address: %v", err)
	}

	return &pb.AddAddressResponse{
		Address: modelToProtoAddress(address),
	}, nil
}

// UpdateAddress updates an existing address
func (h *ProfileHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	addressID, err := uuid.Parse(req.GetAddressId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address_id: %v", err)
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	input := model.UpdateAddressInput{
		ID:     addressID,
		UserID: userID,
	}

	// Only set fields that are provided
	if req.GetLabel() != "" {
		label := req.GetLabel()
		input.Label = &label
	}
	if req.GetStreet() != "" {
		street := req.GetStreet()
		input.Street = &street
	}
	if req.GetCity() != "" {
		city := req.GetCity()
		input.City = &city
	}
	if req.GetState() != "" {
		state := req.GetState()
		input.State = &state
	}
	if req.GetPostalCode() != "" {
		postalCode := req.GetPostalCode()
		input.PostalCode = &postalCode
	}
	if req.GetCountry() != "" {
		country := req.GetCountry()
		input.Country = &country
	}
	// IsDefault can be explicitly set to false, so we always include it
	isDefault := req.GetIsDefault()
	input.IsDefault = &isDefault

	address, err := h.addressService.UpdateAddress(ctx, input)
	if err != nil {
		if err == repository.ErrAddressNotFound {
			return nil, status.Errorf(codes.NotFound, "address not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update address: %v", err)
	}

	return &pb.UpdateAddressResponse{
		Address: modelToProtoAddress(address),
	}, nil
}

// DeleteAddress deletes an address
func (h *ProfileHandler) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	addressID, err := uuid.Parse(req.GetAddressId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address_id: %v", err)
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	err = h.addressService.DeleteAddress(ctx, addressID, userID)
	if err != nil {
		if err == repository.ErrAddressNotFound {
			return nil, status.Errorf(codes.NotFound, "address not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete address: %v", err)
	}

	return &pb.DeleteAddressResponse{
		Success: true,
		Message: "Address deleted successfully",
	}, nil
}

// Legacy methods for backward compatibility

// GetUser retrieves a user (legacy)
func (h *ProfileHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	profile, err := h.profileService.GetProfile(ctx, userID)
	if err != nil {
		if err == repository.ErrProfileNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.GetUserResponse{
		Id:    profile.UserID.String(),
		Name:  profile.FirstName + " " + profile.LastName,
		Email: profile.Email,
	}, nil
}

// CreateUser creates a user (legacy)
func (h *ProfileHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Generate a new user ID for legacy compatibility
	userID := uuid.New()

	input := model.CreateProfileInput{
		UserID:    userID,
		FirstName: req.GetName(),
		Email:     req.GetEmail(),
	}

	profile, err := h.profileService.CreateProfile(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUserResponse{
		Id:    profile.UserID.String(),
		Name:  profile.FirstName,
		Email: profile.Email,
	}, nil
}

// ListUsers lists users (legacy - returns empty for now)
func (h *ProfileHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// This would require pagination support in the repository
	// For now, return empty list
	return &pb.ListUsersResponse{
		Users: []*pb.GetUserResponse{},
		Total: 0,
	}, nil
}

// Helper functions

func modelToProtoProfile(m *model.Profile) *pb.Profile {
	if m == nil {
		return nil
	}

	var dob string
	if m.DateOfBirth != nil {
		dob = m.DateOfBirth.Format("2006-01-02")
	}

	return &pb.Profile{
		UserId:      m.UserID.String(),
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email,
		Phone:       m.Phone,
		AvatarUrl:   m.AvatarURL,
		DateOfBirth: dob,
		Preferences: m.Preferences,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

func modelToProtoAddress(m *model.Address) *pb.Address {
	if m == nil {
		return nil
	}

	return &pb.Address{
		Id:         m.ID.String(),
		UserId:     m.UserID.String(),
		Label:      m.Label,
		Street:     m.Street,
		City:       m.City,
		State:      m.State,
		PostalCode: m.PostalCode,
		Country:    m.Country,
		IsDefault:  m.IsDefault,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  m.UpdatedAt.Format(time.RFC3339),
	}
}

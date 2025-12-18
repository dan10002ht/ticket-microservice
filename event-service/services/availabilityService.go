package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
	"fmt"
)

type AvailabilityService struct {
	repo *repositories.EventSeatAvailabilityRepository
}

func NewAvailabilityService(repo *repositories.EventSeatAvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{repo: repo}
}

// GetEventAvailability - Get all seat availability for an event
func (s *AvailabilityService) GetEventAvailability(ctx context.Context, eventID string) ([]*models.EventSeatAvailability, *models.EventAvailabilitySummary, error) {
	availability, err := s.repo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}

	// Calculate summary
	summary := &models.EventAvailabilitySummary{
		TotalSeats: int32(len(availability)),
	}
	for _, a := range availability {
		switch a.AvailabilityStatus {
		case "available":
			summary.AvailableSeats++
		case "reserved":
			summary.ReservedSeats++
		case "booked":
			summary.BookedSeats++
		case "blocked":
			summary.BlockedSeats++
		}
	}

	return availability, summary, nil
}

// GetZoneAvailability - Get seat availability for a specific zone
func (s *AvailabilityService) GetZoneAvailability(ctx context.Context, eventID, zoneID string) ([]*models.EventSeatAvailability, *models.EventAvailabilitySummary, error) {
	availability, err := s.repo.GetByEventAndZone(ctx, eventID, zoneID)
	if err != nil {
		return nil, nil, err
	}

	// Calculate summary
	summary := &models.EventAvailabilitySummary{
		TotalSeats: int32(len(availability)),
	}
	for _, a := range availability {
		switch a.AvailabilityStatus {
		case "available":
			summary.AvailableSeats++
		case "reserved":
			summary.ReservedSeats++
		case "booked":
			summary.BookedSeats++
		case "blocked":
			summary.BlockedSeats++
		}
	}

	return availability, summary, nil
}

// GetSeatAvailability - Get availability for a specific seat
func (s *AvailabilityService) GetSeatAvailability(ctx context.Context, eventID, seatID string) (*models.EventSeatAvailability, error) {
	return s.repo.GetBySeatID(ctx, eventID, seatID)
}

// UpdateSeatAvailability - Update seat availability status
func (s *AvailabilityService) UpdateSeatAvailability(ctx context.Context, eventID, seatID, status, reservationID, blockedReason, blockedUntil string) error {
	if eventID == "" || seatID == "" || status == "" {
		return fmt.Errorf("invalid update parameters")
	}
	return s.repo.UpdateStatus(ctx, eventID, seatID, status, reservationID, blockedReason, blockedUntil)
}

// BlockSeats - Block multiple seats
func (s *AvailabilityService) BlockSeats(ctx context.Context, eventID string, seatIDs []string, blockedReason, blockedUntil string) (*models.BlockSeatsResult, error) {
	if eventID == "" || len(seatIDs) == 0 {
		return nil, fmt.Errorf("invalid block parameters")
	}

	result := &models.BlockSeatsResult{
		BlockedSeatIDs: make([]string, 0),
	}

	for _, seatID := range seatIDs {
		err := s.repo.UpdateStatus(ctx, eventID, seatID, "blocked", "", blockedReason, blockedUntil)
		if err == nil {
			result.BlockedSeatIDs = append(result.BlockedSeatIDs, seatID)
			result.BlockedCount++
		}
	}

	return result, nil
}

// ReleaseSeats - Release blocked seats
func (s *AvailabilityService) ReleaseSeats(ctx context.Context, eventID string, seatIDs []string) (*models.ReleaseSeatsResult, error) {
	if eventID == "" || len(seatIDs) == 0 {
		return nil, fmt.Errorf("invalid release parameters")
	}

	result := &models.ReleaseSeatsResult{
		ReleasedSeatIDs: make([]string, 0),
	}

	for _, seatID := range seatIDs {
		err := s.repo.UpdateStatus(ctx, eventID, seatID, "available", "", "", "")
		if err == nil {
			result.ReleasedSeatIDs = append(result.ReleasedSeatIDs, seatID)
			result.ReleasedCount++
		}
	}

	return result, nil
}

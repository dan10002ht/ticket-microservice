package grpc

import (
	"context"
	"event-service/services"
	eventpb "shared-lib/protos/event"
)

type AvailabilityController struct {
	service *services.AvailabilityService
	eventpb.UnimplementedAvailabilityServiceServer
}

func NewAvailabilityController(service *services.AvailabilityService) *AvailabilityController {
	return &AvailabilityController{service: service}
}

// GetEventAvailability - Get availability for entire event
func (c *AvailabilityController) GetEventAvailability(ctx context.Context, req *eventpb.GetEventAvailabilityRequest) (*eventpb.GetEventAvailabilityResponse, error) {
	availability, err := c.service.GetEventAvailability(ctx, req.EventId)
	if err != nil {
		return &eventpb.GetEventAvailabilityResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.GetEventAvailabilityResponse{
		Success:        true,
		EventId:        req.EventId,
		TotalSeats:     availability.TotalSeats,
		AvailableSeats: availability.AvailableSeats,
		ReservedSeats:  availability.ReservedSeats,
		SoldSeats:      availability.SoldSeats,
		BlockedSeats:   availability.BlockedSeats,
		Zones:          availability.Zones,
	}, nil
}

// GetZoneAvailability - Get availability for specific zone
func (c *AvailabilityController) GetZoneAvailability(ctx context.Context, req *eventpb.GetZoneAvailabilityRequest) (*eventpb.GetZoneAvailabilityResponse, error) {
	availability, err := c.service.GetZoneAvailability(ctx, req.EventId, req.ZoneId)
	if err != nil {
		return &eventpb.GetZoneAvailabilityResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.GetZoneAvailabilityResponse{
		Success:        true,
		EventId:        req.EventId,
		ZoneId:         req.ZoneId,
		TotalSeats:     availability.TotalSeats,
		AvailableSeats: availability.AvailableSeats,
		ReservedSeats:  availability.ReservedSeats,
		SoldSeats:      availability.SoldSeats,
		BlockedSeats:   availability.BlockedSeats,
		Seats:          availability.Seats,
	}, nil
}

// GetSeatAvailability - Get availability for specific seat
func (c *AvailabilityController) GetSeatAvailability(ctx context.Context, req *eventpb.GetSeatAvailabilityRequest) (*eventpb.GetSeatAvailabilityResponse, error) {
	availability, err := c.service.GetSeatAvailability(ctx, req.EventId, req.SeatId)
	if err != nil {
		return &eventpb.GetSeatAvailabilityResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.GetSeatAvailabilityResponse{
		Success:       true,
		EventId:       req.EventId,
		SeatId:        req.SeatId,
		Status:        availability.Status,
		ReservedBy:    availability.ReservedBy,
		ReservedUntil: availability.ReservedUntil,
		BookedBy:      availability.BookedBy,
		BookingId:     availability.BookingId,
		Price:         availability.Price,
		Currency:      availability.Currency,
	}, nil
}

// UpdateSeatAvailability - Update seat availability status
func (c *AvailabilityController) UpdateSeatAvailability(ctx context.Context, req *eventpb.UpdateSeatAvailabilityRequest) (*eventpb.UpdateSeatAvailabilityResponse, error) {
	err := c.service.UpdateSeatAvailability(ctx, req.EventId, req.SeatId, req.Status, req.UserId, req.BookingId)
	if err != nil {
		return &eventpb.UpdateSeatAvailabilityResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.UpdateSeatAvailabilityResponse{
		Success: true,
		Message: "Seat availability updated successfully",
	}, nil
}

// BlockSeats - Block multiple seats for booking
func (c *AvailabilityController) BlockSeats(ctx context.Context, req *eventpb.BlockSeatsRequest) (*eventpb.BlockSeatsResponse, error) {
	result, err := c.service.BlockSeats(ctx, req.EventId, req.SeatIds, req.UserId, req.BookingId, req.ExpiresAt)
	if err != nil {
		return &eventpb.BlockSeatsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.BlockSeatsResponse{
		Success:      true,
		BlockedSeats: result.BlockedSeats,
		FailedSeats:  result.FailedSeats,
		Message:      "Seats blocked successfully",
	}, nil
}

// ReleaseSeats - Release blocked seats
func (c *AvailabilityController) ReleaseSeats(ctx context.Context, req *eventpb.ReleaseSeatsRequest) (*eventpb.ReleaseSeatsResponse, error) {
	result, err := c.service.ReleaseSeats(ctx, req.EventId, req.SeatIds, req.UserId, req.Reason)
	if err != nil {
		return &eventpb.ReleaseSeatsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.ReleaseSeatsResponse{
		Success:       true,
		ReleasedSeats: result.ReleasedSeats,
		FailedSeats:   result.FailedSeats,
		Message:       "Seats released successfully",
	}, nil
}

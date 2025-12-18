package grpc

import (
	"context"
	"event-service/services"
	eventpb "event-service/internal/protos/event"
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
	availability, summary, err := c.service.GetEventAvailability(ctx, req.EventId)
	if err != nil {
		return &eventpb.GetEventAvailabilityResponse{
			Error: err.Error(),
		}, nil
	}

	var pbAvailability []*eventpb.SeatAvailability
	for _, a := range availability {
		pbAvailability = append(pbAvailability, &eventpb.SeatAvailability{
			Id:                 a.PublicID,
			EventId:            a.EventID,
			SeatId:             a.SeatID,
			ZoneId:             a.ZoneID,
			AvailabilityStatus: a.AvailabilityStatus,
			ReservationId:      a.ReservationID,
			BlockedReason:      a.BlockedReason,
			BlockedUntil:       a.BlockedUntil,
			LastUpdated:        a.LastUpdated,
			CreatedAt:          a.CreatedAt,
			UpdatedAt:          a.UpdatedAt,
		})
	}

	return &eventpb.GetEventAvailabilityResponse{
		Availability:   pbAvailability,
		TotalSeats:     summary.TotalSeats,
		AvailableSeats: summary.AvailableSeats,
		ReservedSeats:  summary.ReservedSeats,
		BookedSeats:    summary.BookedSeats,
		BlockedSeats:   summary.BlockedSeats,
	}, nil
}

// GetZoneAvailability - Get availability for specific zone
func (c *AvailabilityController) GetZoneAvailability(ctx context.Context, req *eventpb.GetZoneAvailabilityRequest) (*eventpb.GetZoneAvailabilityResponse, error) {
	availability, summary, err := c.service.GetZoneAvailability(ctx, req.EventId, req.ZoneId)
	if err != nil {
		return &eventpb.GetZoneAvailabilityResponse{
			Error: err.Error(),
		}, nil
	}

	var pbAvailability []*eventpb.SeatAvailability
	for _, a := range availability {
		pbAvailability = append(pbAvailability, &eventpb.SeatAvailability{
			Id:                 a.PublicID,
			EventId:            a.EventID,
			SeatId:             a.SeatID,
			ZoneId:             a.ZoneID,
			AvailabilityStatus: a.AvailabilityStatus,
			ReservationId:      a.ReservationID,
			BlockedReason:      a.BlockedReason,
			BlockedUntil:       a.BlockedUntil,
			LastUpdated:        a.LastUpdated,
			CreatedAt:          a.CreatedAt,
			UpdatedAt:          a.UpdatedAt,
		})
	}

	return &eventpb.GetZoneAvailabilityResponse{
		Availability:   pbAvailability,
		TotalSeats:     summary.TotalSeats,
		AvailableSeats: summary.AvailableSeats,
		ReservedSeats:  summary.ReservedSeats,
		BookedSeats:    summary.BookedSeats,
		BlockedSeats:   summary.BlockedSeats,
	}, nil
}

// GetSeatAvailability - Get availability for specific seat
func (c *AvailabilityController) GetSeatAvailability(ctx context.Context, req *eventpb.GetSeatAvailabilityRequest) (*eventpb.GetSeatAvailabilityResponse, error) {
	avail, err := c.service.GetSeatAvailability(ctx, req.EventId, req.SeatId)
	if err != nil {
		return &eventpb.GetSeatAvailabilityResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.GetSeatAvailabilityResponse{
		Availability: &eventpb.SeatAvailability{
			Id:                 avail.PublicID,
			EventId:            avail.EventID,
			SeatId:             avail.SeatID,
			ZoneId:             avail.ZoneID,
			AvailabilityStatus: avail.AvailabilityStatus,
			ReservationId:      avail.ReservationID,
			BlockedReason:      avail.BlockedReason,
			BlockedUntil:       avail.BlockedUntil,
			LastUpdated:        avail.LastUpdated,
			CreatedAt:          avail.CreatedAt,
			UpdatedAt:          avail.UpdatedAt,
		},
	}, nil
}

// UpdateSeatAvailability - Update seat availability status
func (c *AvailabilityController) UpdateSeatAvailability(ctx context.Context, req *eventpb.UpdateSeatAvailabilityRequest) (*eventpb.UpdateSeatAvailabilityResponse, error) {
	err := c.service.UpdateSeatAvailability(ctx, req.EventId, req.SeatId, req.AvailabilityStatus, req.ReservationId, req.BlockedReason, req.BlockedUntil)
	if err != nil {
		return &eventpb.UpdateSeatAvailabilityResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.UpdateSeatAvailabilityResponse{
		Success: true,
	}, nil
}

// BlockSeats - Block multiple seats for booking
func (c *AvailabilityController) BlockSeats(ctx context.Context, req *eventpb.BlockSeatsRequest) (*eventpb.BlockSeatsResponse, error) {
	result, err := c.service.BlockSeats(ctx, req.EventId, req.SeatIds, req.BlockedReason, req.BlockedUntil)
	if err != nil {
		return &eventpb.BlockSeatsResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.BlockSeatsResponse{
		BlockedCount:   result.BlockedCount,
		BlockedSeatIds: result.BlockedSeatIDs,
	}, nil
}

// ReleaseSeats - Release blocked seats
func (c *AvailabilityController) ReleaseSeats(ctx context.Context, req *eventpb.ReleaseSeatsRequest) (*eventpb.ReleaseSeatsResponse, error) {
	result, err := c.service.ReleaseSeats(ctx, req.EventId, req.SeatIds)
	if err != nil {
		return &eventpb.ReleaseSeatsResponse{
			Error: err.Error(),
		}, nil
	}

	return &eventpb.ReleaseSeatsResponse{
		ReleasedCount:   result.ReleasedCount,
		ReleasedSeatIds: result.ReleasedSeatIDs,
	}, nil
}

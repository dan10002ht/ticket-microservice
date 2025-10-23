package grpc

import (
	"context"
	"event-service/services"
	eventpb "shared-lib/protos/event"
)

type EventSeatController struct {
	service *services.EventSeatService
	eventpb.UnimplementedEventSeatServiceServer
}

func NewEventSeatController(service *services.EventSeatService) *EventSeatController {
	return &EventSeatController{service: service}
}

// CreateSeat - Create new seat
func (c *EventSeatController) CreateSeat(ctx context.Context, req *eventpb.CreateSeatRequest) (*eventpb.CreateSeatResponse, error) {
	seat, err := c.service.CreateSeat(ctx, req.EventId, req.ZoneId, req.SeatNumber, req.RowNumber, req.Coordinates, req.PricingCategory, req.BasePrice, req.FinalPrice, req.Currency)
	if err != nil {
		return &eventpb.CreateSeatResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.CreateSeatResponse{
		Success: true,
		Seat: &eventpb.EventSeat{
			Id:              seat.PublicID,
			EventId:         seat.EventID,
			ZoneId:          seat.ZoneID,
			SeatNumber:      seat.SeatNumber,
			RowNumber:       seat.RowNumber,
			Coordinates:     seat.Coordinates,
			Status:          seat.Status,
			PricingCategory: seat.PricingCategory,
			BasePrice:       seat.BasePrice,
			FinalPrice:      seat.FinalPrice,
			Currency:        seat.Currency,
			CreatedAt:       seat.CreatedAt,
			UpdatedAt:       seat.UpdatedAt,
		},
		Message: "Seat created successfully",
	}, nil
}

// GetSeat - Get seat by ID
func (c *EventSeatController) GetSeat(ctx context.Context, req *eventpb.GetSeatRequest) (*eventpb.GetSeatResponse, error) {
	seat, err := c.service.GetSeat(ctx, req.SeatId)
	if err != nil {
		return &eventpb.GetSeatResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.GetSeatResponse{
		Success: true,
		Seat: &eventpb.EventSeat{
			Id:              seat.PublicID,
			EventId:         seat.EventID,
			ZoneId:          seat.ZoneID,
			SeatNumber:      seat.SeatNumber,
			RowNumber:       seat.RowNumber,
			Coordinates:     seat.Coordinates,
			Status:          seat.Status,
			PricingCategory: seat.PricingCategory,
			BasePrice:       seat.BasePrice,
			FinalPrice:      seat.FinalPrice,
			Currency:        seat.Currency,
			CreatedAt:       seat.CreatedAt,
			UpdatedAt:       seat.UpdatedAt,
		},
	}, nil
}

// UpdateSeat - Update seat information
func (c *EventSeatController) UpdateSeat(ctx context.Context, req *eventpb.UpdateSeatRequest) (*eventpb.UpdateSeatResponse, error) {
	seat, err := c.service.UpdateSeat(ctx, req.SeatId, req.SeatNumber, req.RowNumber, req.Coordinates, req.PricingCategory, req.BasePrice, req.FinalPrice, req.Currency)
	if err != nil {
		return &eventpb.UpdateSeatResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.UpdateSeatResponse{
		Success: true,
		Seat: &eventpb.EventSeat{
			Id:              seat.PublicID,
			EventId:         seat.EventID,
			ZoneId:          seat.ZoneID,
			SeatNumber:      seat.SeatNumber,
			RowNumber:       seat.RowNumber,
			Coordinates:     seat.Coordinates,
			Status:          seat.Status,
			PricingCategory: seat.PricingCategory,
			BasePrice:       seat.BasePrice,
			FinalPrice:      seat.FinalPrice,
			Currency:        seat.Currency,
			CreatedAt:       seat.CreatedAt,
			UpdatedAt:       seat.UpdatedAt,
		},
		Message: "Seat updated successfully",
	}, nil
}

// DeleteSeat - Delete seat
func (c *EventSeatController) DeleteSeat(ctx context.Context, req *eventpb.DeleteSeatRequest) (*eventpb.DeleteSeatResponse, error) {
	err := c.service.DeleteSeat(ctx, req.SeatId)
	if err != nil {
		return &eventpb.DeleteSeatResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.DeleteSeatResponse{
		Success: true,
		Message: "Seat deleted successfully",
	}, nil
}

// ListSeatsByEvent - List all seats for an event
func (c *EventSeatController) ListSeatsByEvent(ctx context.Context, req *eventpb.ListSeatsByEventRequest) (*eventpb.ListSeatsByEventResponse, error) {
	seats, total, err := c.service.ListSeatsByEvent(ctx, req.EventId, req.ZoneId, req.Status, req.Page, req.Limit)
	if err != nil {
		return &eventpb.ListSeatsByEventResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	var pbSeats []*eventpb.EventSeat
	for _, seat := range seats {
		pbSeats = append(pbSeats, &eventpb.EventSeat{
			Id:              seat.PublicID,
			EventId:         seat.EventID,
			ZoneId:          seat.ZoneID,
			SeatNumber:      seat.SeatNumber,
			RowNumber:       seat.RowNumber,
			Coordinates:     seat.Coordinates,
			Status:          seat.Status,
			PricingCategory: seat.PricingCategory,
			BasePrice:       seat.BasePrice,
			FinalPrice:      seat.FinalPrice,
			Currency:        seat.Currency,
			CreatedAt:       seat.CreatedAt,
			UpdatedAt:       seat.UpdatedAt,
		})
	}

	return &eventpb.ListSeatsByEventResponse{
		Success: true,
		Seats:   pbSeats,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

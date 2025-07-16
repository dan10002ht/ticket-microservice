package grpc

import (
	"context"
	"event-service/models"
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

func (c *EventSeatController) CreateSeat(ctx context.Context, req *eventpb.CreateEventSeatRequest) (*eventpb.CreateEventSeatResponse, error) {
	seat := &models.EventSeat{
		PublicID:    req.Seat.Id,
		EventID:     req.Seat.EventId,
		ZoneID:      req.Seat.ZoneId,
		SeatNumber:  req.Seat.SeatNumber,
		RowNumber:   req.Seat.RowNumber,
		Coordinates: req.Seat.Coordinates,
	}
	err := c.service.CreateSeat(ctx, seat)
	if err != nil {
		return nil, err
	}
	return &eventpb.CreateEventSeatResponse{Seat: req.Seat}, nil
}

func (c *EventSeatController) GetSeat(ctx context.Context, req *eventpb.GetEventSeatRequest) (*eventpb.GetEventSeatResponse, error) {
	seat, err := c.service.GetSeat(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.GetEventSeatResponse{Seat: &eventpb.EventSeat{
		Id:         seat.PublicID,
		EventId:    seat.EventID,
		ZoneId:     seat.ZoneID,
		SeatNumber: seat.SeatNumber,
		RowNumber:  seat.RowNumber,
		Coordinates: seat.Coordinates,
		CreatedAt:  seat.CreatedAt,
		UpdatedAt:  seat.UpdatedAt,
	}}, nil
}

func (c *EventSeatController) UpdateSeat(ctx context.Context, req *eventpb.UpdateEventSeatRequest) (*eventpb.UpdateEventSeatResponse, error) {
	seat := &models.EventSeat{
		PublicID:    req.Seat.Id,
		EventID:     req.Seat.EventId,
		ZoneID:      req.Seat.ZoneId,
		SeatNumber:  req.Seat.SeatNumber,
		RowNumber:   req.Seat.RowNumber,
		Coordinates: req.Seat.Coordinates,
	}
	err := c.service.UpdateSeat(ctx, seat)
	if err != nil {
		return nil, err
	}
	return &eventpb.UpdateEventSeatResponse{Seat: req.Seat}, nil
}

func (c *EventSeatController) DeleteSeat(ctx context.Context, req *eventpb.DeleteEventSeatRequest) (*eventpb.DeleteEventSeatResponse, error) {
	err := c.service.DeleteSeat(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.DeleteEventSeatResponse{Success: true}, nil
}

func (c *EventSeatController) ListSeats(ctx context.Context, req *eventpb.ListEventSeatsRequest) (*eventpb.ListEventSeatsResponse, error) {
	seats, err := c.service.ListByEventID(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	var pbSeats []*eventpb.EventSeat
	for _, seat := range seats {
		pbSeats = append(pbSeats, &eventpb.EventSeat{
			Id:         seat.PublicID,
			EventId:    seat.EventID,
			ZoneId:     seat.ZoneID,
			SeatNumber: seat.SeatNumber,
			RowNumber:  seat.RowNumber,
			Coordinates: seat.Coordinates,
			CreatedAt:  seat.CreatedAt,
			UpdatedAt:  seat.UpdatedAt,
		})
	}
	return &eventpb.ListEventSeatsResponse{Seats: pbSeats}, nil
} 
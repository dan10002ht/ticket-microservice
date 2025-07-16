package grpc

import (
	"context"
	seatpb "shared-lib/protos/ticket_booking"
	"ticket-service/services"
)

type SeatController struct {
	seatpb.UnimplementedSeatServiceServer
	service *services.SeatReservationService
}

func NewSeatController(service *services.SeatReservationService) *SeatController {
	return &SeatController{service: service}
}

func (c *SeatController) ReserveSeats(ctx context.Context, req *seatpb.ReserveSeatsRequest) (*seatpb.ReserveSeatsResponse, error) {
	// TODO: Implement logic
	return &seatpb.ReserveSeatsResponse{}, nil
}

func (c *SeatController) ReleaseSeats(ctx context.Context, req *seatpb.ReleaseSeatsRequest) (*seatpb.ReleaseSeatsResponse, error) {
	// TODO: Implement logic
	return &seatpb.ReleaseSeatsResponse{}, nil
}

// ... Các method khác: ConfirmSeats, GetReservedSeats, GetAvailableSeats, GetSeatReservation, ReserveSeatsBulk, ReleaseSeatsBulk ... 
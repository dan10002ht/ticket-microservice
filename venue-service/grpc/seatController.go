package grpc

import (
	"venue-service/services"
)

type SeatController struct {
	seatService *services.SeatService
}

func NewSeatController(seatService *services.SeatService) *SeatController {
	return &SeatController{seatService: seatService}
}

// TODO: Implement gRPC methods for Seat CRUD, ValidateSeat, BulkCreateSeats, etc. 
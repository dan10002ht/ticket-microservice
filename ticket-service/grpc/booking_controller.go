package grpc

import (
	"context"
	bookingpb "shared-lib/protos/ticket_booking"
	"ticket-service/services"
)

type BookingController struct {
	bookingpb.UnimplementedBookingServiceServer
	service *services.BookingSessionService
}

func NewBookingController(service *services.BookingSessionService) *BookingController {
	return &BookingController{service: service}
}

func (c *BookingController) CreateBookingSession(ctx context.Context, req *bookingpb.CreateBookingSessionRequest) (*bookingpb.CreateBookingSessionResponse, error) {
	// TODO: Implement logic
	return &bookingpb.CreateBookingSessionResponse{}, nil
}

func (c *BookingController) GetBookingSession(ctx context.Context, req *bookingpb.GetBookingSessionRequest) (*bookingpb.GetBookingSessionResponse, error) {
	// TODO: Implement logic
	return &bookingpb.GetBookingSessionResponse{}, nil
}

// ... Các method khác: UpdateBookingSession, CancelBookingSession, StartBooking, CompleteBooking, ExtendBookingSession, GetUserBookings, GetActiveBookings ... 
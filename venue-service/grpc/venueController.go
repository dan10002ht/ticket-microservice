package grpc

import (
	"venue-service/services"
)

type VenueController struct {
	venueService *services.VenueService
}

func NewVenueController(venueService *services.VenueService) *VenueController {
	return &VenueController{venueService: venueService}
}

// TODO: Implement gRPC methods for Venue CRUD, GetVenue, etc. 
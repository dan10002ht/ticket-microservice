package grpc

import (
	"venue-service/services"
)

type LayoutController struct {
	layoutService *services.LayoutService
}

func NewLayoutController(layoutService *services.LayoutService) *LayoutController {
	return &LayoutController{layoutService: layoutService}
}

// TODO: Implement gRPC methods for Layout CRUD, ValidateLayout, etc. 
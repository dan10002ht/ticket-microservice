package grpc

import (
	"venue-service/services"
)

type ZoneController struct {
	zoneService *services.ZoneService
}

func NewZoneController(zoneService *services.ZoneService) *ZoneController {
	return &ZoneController{zoneService: zoneService}
}

// TODO: Implement gRPC methods for Zone CRUD, ValidateZone, etc. 
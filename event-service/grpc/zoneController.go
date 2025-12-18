package grpc

import (
	"context"
	"event-service/services"
	eventpb "event-service/internal/protos/event"
)

type ZoneController struct {
	service *services.EventSeatingZoneService
	eventpb.UnimplementedEventSeatingZoneServiceServer
}

func NewZoneController(service *services.EventSeatingZoneService) *ZoneController {
	return &ZoneController{service: service}
}

// CreateZone - Create new seating zone
func (c *ZoneController) CreateZone(ctx context.Context, req *eventpb.CreateZoneRequest) (*eventpb.CreateZoneResponse, error) {
	zone, err := c.service.CreateZone(ctx, req.EventId, req.Name, req.ZoneType, req.Coordinates, req.Color)
	if err != nil {
		return &eventpb.CreateZoneResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.CreateZoneResponse{
		Success: true,
		Zone: &eventpb.EventSeatingZone{
			Id:          zone.PublicID,
			EventId:     zone.EventID,
			Name:        zone.Name,
			ZoneType:    zone.ZoneType,
			Coordinates: zone.Coordinates,
			SeatCount:   int32(zone.SeatCount),
			Color:       zone.Color,
			CreatedAt:   zone.CreatedAt,
			UpdatedAt:   zone.UpdatedAt,
		},
		Message: "Zone created successfully",
	}, nil
}

// GetZone - Get zone by ID
func (c *ZoneController) GetZone(ctx context.Context, req *eventpb.GetZoneRequest) (*eventpb.GetZoneResponse, error) {
	zone, err := c.service.GetZone(ctx, req.ZoneId)
	if err != nil {
		return &eventpb.GetZoneResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.GetZoneResponse{
		Success: true,
		Zone: &eventpb.EventSeatingZone{
			Id:          zone.PublicID,
			EventId:     zone.EventID,
			Name:        zone.Name,
			ZoneType:    zone.ZoneType,
			Coordinates: zone.Coordinates,
			SeatCount:   int32(zone.SeatCount),
			Color:       zone.Color,
			CreatedAt:   zone.CreatedAt,
			UpdatedAt:   zone.UpdatedAt,
		},
	}, nil
}

// UpdateZone - Update zone information
func (c *ZoneController) UpdateZone(ctx context.Context, req *eventpb.UpdateZoneRequest) (*eventpb.UpdateZoneResponse, error) {
	zone, err := c.service.UpdateZone(ctx, req.ZoneId, req.Name, req.ZoneType, req.Coordinates, req.Color)
	if err != nil {
		return &eventpb.UpdateZoneResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.UpdateZoneResponse{
		Success: true,
		Zone: &eventpb.EventSeatingZone{
			Id:          zone.PublicID,
			EventId:     zone.EventID,
			Name:        zone.Name,
			ZoneType:    zone.ZoneType,
			Coordinates: zone.Coordinates,
			SeatCount:   int32(zone.SeatCount),
			Color:       zone.Color,
			CreatedAt:   zone.CreatedAt,
			UpdatedAt:   zone.UpdatedAt,
		},
		Message: "Zone updated successfully",
	}, nil
}

// DeleteZone - Delete zone
func (c *ZoneController) DeleteZone(ctx context.Context, req *eventpb.DeleteZoneRequest) (*eventpb.DeleteZoneResponse, error) {
	err := c.service.DeleteZone(ctx, req.ZoneId)
	if err != nil {
		return &eventpb.DeleteZoneResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &eventpb.DeleteZoneResponse{
		Success: true,
		Message: "Zone deleted successfully",
	}, nil
}

// ListZonesByEvent - List all zones for an event
func (c *ZoneController) ListZonesByEvent(ctx context.Context, req *eventpb.ListZonesByEventRequest) (*eventpb.ListZonesByEventResponse, error) {
	zones, err := c.service.ListZonesByEvent(ctx, req.EventId)
	if err != nil {
		return &eventpb.ListZonesByEventResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	var pbZones []*eventpb.EventSeatingZone
	for _, zone := range zones {
		pbZones = append(pbZones, &eventpb.EventSeatingZone{
			Id:          zone.PublicID,
			EventId:     zone.EventID,
			Name:        zone.Name,
			ZoneType:    zone.ZoneType,
			Coordinates: zone.Coordinates,
			SeatCount:   int32(zone.SeatCount),
			Color:       zone.Color,
			CreatedAt:   zone.CreatedAt,
			UpdatedAt:   zone.UpdatedAt,
		})
	}

	return &eventpb.ListZonesByEventResponse{
		Success: true,
		Zones:   pbZones,
		Total:   int32(len(pbZones)),
	}, nil
}

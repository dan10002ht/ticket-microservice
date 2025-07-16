package grpc

import (
	"context"
	"event-service/models"
	"event-service/services"
	eventpb "shared-lib/protos/event"
)

type EventSeatingZoneController struct {
	service *services.EventSeatingZoneService
	eventpb.UnimplementedEventSeatingZoneServiceServer
}

func NewEventSeatingZoneController(service *services.EventSeatingZoneService) *EventSeatingZoneController {
	return &EventSeatingZoneController{service: service}
}

func (c *EventSeatingZoneController) CreateZone(ctx context.Context, req *eventpb.CreateEventSeatingZoneRequest) (*eventpb.CreateEventSeatingZoneResponse, error) {
	zone := &models.EventSeatingZone{
		PublicID:    req.Zone.Id,
		EventID:     req.Zone.EventId,
		Name:        req.Zone.Name,
		ZoneType:    req.Zone.ZoneType,
		Coordinates: req.Zone.Coordinates,
		SeatCount:   int(req.Zone.SeatCount),
		Color:       req.Zone.Color,
	}
	err := c.service.CreateZone(ctx, zone)
	if err != nil {
		return nil, err
	}
	return &eventpb.CreateEventSeatingZoneResponse{Zone: req.Zone}, nil
}

func (c *EventSeatingZoneController) GetZone(ctx context.Context, req *eventpb.GetEventSeatingZoneRequest) (*eventpb.GetEventSeatingZoneResponse, error) {
	zone, err := c.service.GetZone(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.GetEventSeatingZoneResponse{Zone: &eventpb.EventSeatingZone{
		Id:          zone.PublicID,
		EventId:     zone.EventID,
		Name:        zone.Name,
		ZoneType:    zone.ZoneType,
		Coordinates: zone.Coordinates,
		SeatCount:   int32(zone.SeatCount),
		Color:       zone.Color,
		CreatedAt:   zone.CreatedAt,
		UpdatedAt:   zone.UpdatedAt,
	}}, nil
}

func (c *EventSeatingZoneController) UpdateZone(ctx context.Context, req *eventpb.UpdateEventSeatingZoneRequest) (*eventpb.UpdateEventSeatingZoneResponse, error) {
	zone := &models.EventSeatingZone{
		PublicID:    req.Zone.Id,
		EventID:     req.Zone.EventId,
		Name:        req.Zone.Name,
		ZoneType:    req.Zone.ZoneType,
		Coordinates: req.Zone.Coordinates,
		SeatCount:   int(req.Zone.SeatCount),
		Color:       req.Zone.Color,
	}
	err := c.service.UpdateZone(ctx, zone)
	if err != nil {
		return nil, err
	}
	return &eventpb.UpdateEventSeatingZoneResponse{Zone: req.Zone}, nil
}

func (c *EventSeatingZoneController) DeleteZone(ctx context.Context, req *eventpb.DeleteEventSeatingZoneRequest) (*eventpb.DeleteEventSeatingZoneResponse, error) {
	err := c.service.DeleteZone(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.DeleteEventSeatingZoneResponse{Success: true}, nil
}

func (c *EventSeatingZoneController) ListZones(ctx context.Context, req *eventpb.ListEventSeatingZonesRequest) (*eventpb.ListEventSeatingZonesResponse, error) {
	zones, err := c.service.ListByEventID(ctx, req.EventId)
	if err != nil {
		return nil, err
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
	return &eventpb.ListEventSeatingZonesResponse{Zones: pbZones}, nil
} 
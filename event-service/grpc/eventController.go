package grpc

import (
	"context"
	"event-service/models"
	"event-service/services"
	eventpb "shared-lib/protos/event"
)

type EventController struct {
	service *services.EventService
	eventpb.UnimplementedEventServiceServer
}

func NewEventController(service *services.EventService) *EventController {
	return &EventController{service: service}
}

func (c *EventController) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	event := &models.Event{
		PublicID:       req.Event.Id,
		OrganizationID: req.Event.OrganizationId,
		Name:           req.Event.Name,
		Description:    req.Event.Description,
		StartDate:      req.Event.StartDate,
		EndDate:        req.Event.EndDate,
		VenueName:      req.Event.VenueName,
		VenueAddress:   req.Event.VenueAddress,
		VenueCity:      req.Event.VenueCity,
		VenueCountry:   req.Event.VenueCountry,
		VenueCapacity:  int(req.Event.VenueCapacity),
		CanvasConfig:   req.Event.CanvasConfig,
	}
	err := c.service.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &eventpb.CreateEventResponse{Event: req.Event}, nil
}

func (c *EventController) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	event, err := c.service.GetEvent(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.GetEventResponse{Event: &eventpb.Event{
		Id:             event.PublicID,
		OrganizationId: event.OrganizationID,
		Name:           event.Name,
		Description:    event.Description,
		StartDate:      event.StartDate,
		EndDate:        event.EndDate,
		VenueName:      event.VenueName,
		VenueAddress:   event.VenueAddress,
		VenueCity:      event.VenueCity,
		VenueCountry:   event.VenueCountry,
		VenueCapacity:  int32(event.VenueCapacity),
		CanvasConfig:   event.CanvasConfig,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}}, nil
}

func (c *EventController) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	event := &models.Event{
		PublicID:       req.Event.Id,
		OrganizationID: req.Event.OrganizationId,
		Name:           req.Event.Name,
		Description:    req.Event.Description,
		StartDate:      req.Event.StartDate,
		EndDate:        req.Event.EndDate,
		VenueName:      req.Event.VenueName,
		VenueAddress:   req.Event.VenueAddress,
		VenueCity:      req.Event.VenueCity,
		VenueCountry:   req.Event.VenueCountry,
		VenueCapacity:  int(req.Event.VenueCapacity),
		CanvasConfig:   req.Event.CanvasConfig,
	}
	err := c.service.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return &eventpb.UpdateEventResponse{Event: req.Event}, nil
}

func (c *EventController) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	err := c.service.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.DeleteEventResponse{Success: true}, nil
}

func (c *EventController) ListEvents(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) {
	events, err := c.service.ListByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	var pbEvents []*eventpb.Event
	for _, event := range events {
		pbEvents = append(pbEvents, &eventpb.Event{
			Id:             event.PublicID,
			OrganizationId: event.OrganizationID,
			Name:           event.Name,
			Description:    event.Description,
			StartDate:      event.StartDate,
			EndDate:        event.EndDate,
			VenueName:      event.VenueName,
			VenueAddress:   event.VenueAddress,
			VenueCity:      event.VenueCity,
			VenueCountry:   event.VenueCountry,
			VenueCapacity:  int32(event.VenueCapacity),
			CanvasConfig:   event.CanvasConfig,
			CreatedAt:      event.CreatedAt,
			UpdatedAt:      event.UpdatedAt,
		})
	}
	return &eventpb.ListEventsResponse{Events: pbEvents}, nil
} 
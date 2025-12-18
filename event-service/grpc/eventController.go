package grpc

import (
	"context"
	eventpb "event-service/internal/protos/event"
	"event-service/models"
	"event-service/services"
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
		Name:          req.Name,
		Description:   req.Description,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		VenueCapacity: int(req.MaxCapacity),
	}
	err := c.service.CreateEvent(ctx, event)
	if err != nil {
		return &eventpb.CreateEventResponse{Error: err.Error()}, nil
	}
	return &eventpb.CreateEventResponse{
		Event: &eventpb.Event{
			Id:            event.PublicID,
			Name:          event.Name,
			Description:   event.Description,
			StartDate:     event.StartDate,
			EndDate:       event.EndDate,
			VenueCapacity: int32(event.VenueCapacity),
			CreatedAt:     event.CreatedAt,
			UpdatedAt:     event.UpdatedAt,
		},
	}, nil
}

func (c *EventController) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	event, err := c.service.GetEvent(ctx, req.Id)
	if err != nil {
		return &eventpb.GetEventResponse{Error: err.Error()}, nil
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
		PublicID:    req.Id,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
	err := c.service.UpdateEvent(ctx, event)
	if err != nil {
		return &eventpb.UpdateEventResponse{Error: err.Error()}, nil
	}
	return &eventpb.UpdateEventResponse{
		Event: &eventpb.Event{
			Id:          event.PublicID,
			Name:        event.Name,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		},
	}, nil
}

func (c *EventController) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	err := c.service.DeleteEvent(ctx, req.Id)
	if err != nil {
		return &eventpb.DeleteEventResponse{Success: false, Error: err.Error()}, nil
	}
	return &eventpb.DeleteEventResponse{Success: true}, nil
}

func (c *EventController) ListEvents(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) {
	events, err := c.service.List(ctx, req.Page, req.Limit)
	if err != nil {
		return &eventpb.ListEventsResponse{Error: err.Error()}, nil
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
	return &eventpb.ListEventsResponse{Events: pbEvents, Total: int32(len(pbEvents))}, nil
}

func (c *EventController) SearchEvents(ctx context.Context, req *eventpb.SearchEventsRequest) (*eventpb.SearchEventsResponse, error) {
	// TODO: Implement search
	return &eventpb.SearchEventsResponse{Error: "not implemented"}, nil
}

func (c *EventController) GetEventsByVenue(ctx context.Context, req *eventpb.GetEventsByVenueRequest) (*eventpb.GetEventsByVenueResponse, error) {
	// TODO: Implement
	return &eventpb.GetEventsByVenueResponse{Error: "not implemented"}, nil
}

func (c *EventController) GetUpcomingEvents(ctx context.Context, req *eventpb.GetUpcomingEventsRequest) (*eventpb.GetUpcomingEventsResponse, error) {
	// TODO: Implement
	return &eventpb.GetUpcomingEventsResponse{Error: "not implemented"}, nil
}

func (c *EventController) GetFeaturedEvents(ctx context.Context, req *eventpb.GetFeaturedEventsRequest) (*eventpb.GetFeaturedEventsResponse, error) {
	// TODO: Implement
	return &eventpb.GetFeaturedEventsResponse{Error: "not implemented"}, nil
}

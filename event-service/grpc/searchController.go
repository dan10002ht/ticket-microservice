package grpc

import (
	"context"
	"event-service/services"
	eventpb "shared-lib/protos/event"
)

type SearchController struct {
	service *services.EventService
	eventpb.UnimplementedEventServiceServer
}

func NewSearchController(service *services.EventService) *SearchController {
	return &SearchController{service: service}
}

// SearchEvents - Search events by query
func (c *SearchController) SearchEvents(ctx context.Context, req *eventpb.SearchEventsRequest) (*eventpb.SearchEventsResponse, error) {
	events, total, err := c.service.SearchEvents(ctx, req.Query, req.EventType, req.Category, req.Page, req.Limit)
	if err != nil {
		return &eventpb.SearchEventsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
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

	return &eventpb.SearchEventsResponse{
		Success: true,
		Events:  pbEvents,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// GetEventsByVenue - Get events by venue
func (c *SearchController) GetEventsByVenue(ctx context.Context, req *eventpb.GetEventsByVenueRequest) (*eventpb.GetEventsByVenueResponse, error) {
	events, total, err := c.service.GetEventsByVenue(ctx, req.VenueId, req.Status, req.Page, req.Limit)
	if err != nil {
		return &eventpb.GetEventsByVenueResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
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

	return &eventpb.GetEventsByVenueResponse{
		Success: true,
		Events:  pbEvents,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// GetUpcomingEvents - Get upcoming events
func (c *SearchController) GetUpcomingEvents(ctx context.Context, req *eventpb.GetUpcomingEventsRequest) (*eventpb.GetUpcomingEventsResponse, error) {
	events, total, err := c.service.GetUpcomingEvents(ctx, req.Days, req.EventType, req.Category, req.Page, req.Limit)
	if err != nil {
		return &eventpb.GetUpcomingEventsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
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

	return &eventpb.GetUpcomingEventsResponse{
		Success: true,
		Events:  pbEvents,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// GetFeaturedEvents - Get featured events
func (c *SearchController) GetFeaturedEvents(ctx context.Context, req *eventpb.GetFeaturedEventsRequest) (*eventpb.GetFeaturedEventsResponse, error) {
	events, total, err := c.service.GetFeaturedEvents(ctx, req.EventType, req.Category, req.Page, req.Limit)
	if err != nil {
		return &eventpb.GetFeaturedEventsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
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

	return &eventpb.GetFeaturedEventsResponse{
		Success: true,
		Events:  pbEvents,
		Total:   int32(total),
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

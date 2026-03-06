package grpc

import (
	"context"
	eventpb "event-service/internal/protos/event"
	"event-service/models"
	"event-service/services"

	"github.com/google/uuid"
)

type EventController struct {
	service *services.EventService
	eventpb.UnimplementedEventServiceServer
}

func NewEventController(service *services.EventService) *EventController {
	return &EventController{service: service}
}

func eventToProto(event *models.Event) *eventpb.Event {
	return &eventpb.Event{
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
		Status:         event.Status,
		EventType:      event.EventType,
		Category:       event.Category,
		SaleStartDate:  event.SaleStartDate,
		SaleEndDate:    event.SaleEndDate,
		MinAge:         int32(event.MinAge),
		IsFeatured:     event.IsFeatured,
		Images:         event.Images,
		Tags:           event.Tags,
		Metadata:       event.Metadata,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func (c *EventController) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	canvasConfig := req.CanvasConfig
	if canvasConfig == "" {
		canvasConfig = "{}"
	}
	images := req.Images
	if images == "" {
		images = "[]"
	}
	tags := req.Tags
	if tags == "" {
		tags = "[]"
	}
	metadata := req.Metadata
	if metadata == "" {
		metadata = "{}"
	}

	event := &models.Event{
		PublicID:       uuid.New().String(),
		OrganizationID: req.CreatedBy,
		Name:           req.Name,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		VenueName:      req.VenueName,
		VenueAddress:   req.VenueAddress,
		VenueCity:      req.VenueCity,
		VenueCountry:   req.VenueCountry,
		VenueCapacity:  int(req.VenueCapacity),
		CanvasConfig:   canvasConfig,
		Status:         "draft",
		EventType:      req.EventType,
		Category:       req.Category,
		SaleStartDate:  req.SaleStartDate,
		SaleEndDate:    req.SaleEndDate,
		MinAge:         int(req.MinAge),
		IsFeatured:     req.IsFeatured,
		Images:         images,
		Tags:           tags,
		Metadata:       metadata,
	}
	err := c.service.CreateEvent(ctx, event)
	if err != nil {
		return &eventpb.CreateEventResponse{Error: err.Error()}, nil
	}
	return &eventpb.CreateEventResponse{Event: eventToProto(event)}, nil
}

func (c *EventController) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	event, err := c.service.GetEvent(ctx, req.Id)
	if err != nil {
		return &eventpb.GetEventResponse{Error: err.Error()}, nil
	}
	return &eventpb.GetEventResponse{Event: eventToProto(event)}, nil
}

func (c *EventController) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	event := &models.Event{
		PublicID:      req.Id,
		Name:          req.Name,
		Description:   req.Description,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		VenueCapacity: int(req.VenueCapacity),
		Status:        req.Status,
		EventType:     req.EventType,
		Category:      req.Category,
		SaleStartDate: req.SaleStartDate,
		SaleEndDate:   req.SaleEndDate,
		MinAge:        int(req.MinAge),
		IsFeatured:    req.IsFeatured,
		Images:        req.Images,
		Tags:          req.Tags,
		Metadata:      req.Metadata,
	}
	err := c.service.UpdateEvent(ctx, event)
	if err != nil {
		return &eventpb.UpdateEventResponse{Error: err.Error()}, nil
	}
	return &eventpb.UpdateEventResponse{Event: eventToProto(event)}, nil
}

func (c *EventController) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	err := c.service.DeleteEvent(ctx, req.Id)
	if err != nil {
		return &eventpb.DeleteEventResponse{Success: false, Error: err.Error()}, nil
	}
	return &eventpb.DeleteEventResponse{Success: true}, nil
}

func (c *EventController) ListEvents(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) {
	events, total, err := c.service.List(ctx, req.Page, req.Limit, req.Status, req.EventType, req.Category, req.StartDateFrom, req.StartDateTo)
	if err != nil {
		return &eventpb.ListEventsResponse{Error: err.Error()}, nil
	}
	var pbEvents []*eventpb.Event
	for _, event := range events {
		pbEvents = append(pbEvents, eventToProto(event))
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	return &eventpb.ListEventsResponse{
		Events: pbEvents,
		Total:  int32(total),
		Page:   page,
		Limit:  limit,
	}, nil
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

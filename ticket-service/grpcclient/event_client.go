package grpcclient

import (
	context "context"
	eventpb "shared-lib/protos/event"

	"google.golang.org/grpc"
)

type EventClient struct {
	EventService        eventpb.EventServiceClient
	AvailabilityService eventpb.AvailabilityServiceClient
}

func NewEventClient(conn *grpc.ClientConn) *EventClient {
	return &EventClient{
		EventService:        eventpb.NewEventServiceClient(conn),
		AvailabilityService: eventpb.NewAvailabilityServiceClient(conn),
	}
}

func (c *EventClient) GetEvent(ctx context.Context, eventID string) (*eventpb.Event, error) {
	// TODO: Implement call to GetEvent
	return nil, nil
}

func (c *EventClient) GetEventAvailability(ctx context.Context, eventID string) (*eventpb.GetEventAvailabilityResponse, error) {
	// TODO: Implement call to GetEventAvailability
	return nil, nil
}

// ... Các method khác: GetZoneAvailability, GetSeatAvailability ... 